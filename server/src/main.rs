#[macro_use] extern crate log;

use axum::{extract::ConnectInfo, routing::get, Router, extract::State};
use axum_server::tls_rustls::{RustlsConfig};
use clap::{arg, Parser};
use std::{
	env,
	io::Error,
	net::{Ipv4Addr, Ipv6Addr, SocketAddr},
	path::PathBuf,
	sync::{
		atomic::{AtomicU64, Ordering},
		Arc
	},
};
use std::env::VarError;
use tokio::time::{self, Duration, Instant};


/// A HTTP server that echoes your remote IP-address
#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
	/// IPv4 address to bind to (can be provided multiple times)
	#[arg(short = '4', long = "ipv4", required = true)]
	ipv4_addrs: Vec<String>,

	/// IPv6 address to bind to (can be provided multiple times)
	#[arg(short = '6', long = "ipv6", required = true)]
	ipv6_addrs: Vec<String>,

	/// Port to bind to
	#[arg(short, long, default_value_t = 11313)]
	port: u16,

	/// Certificate file path
	#[arg(short, long)]
	cert_path: String,

	/// Key file path
	#[arg(short, long)]
	key_path: String,

	/// Log interval in seconds
	#[arg(short = 'i', long = "interval", default_value_t = 60)]
	log_interval: u64,
}

#[derive(Clone)]
struct ServerConfig {
	req_counter: Arc<AtomicU64>,
}

//noinspection DuplicatedCode
#[tokio::main]
async fn main() {
	match env::var("RUST_LOG") {
		Err(_) => {
			env::set_var("RUST_LOG", "info");
		}
		Ok(_) => {}
	}

	env_logger::init();

	let args = Args::parse();

	// Configure certificate and private key used by HTTPS
	let config = RustlsConfig::from_pem_file(
		PathBuf::from(args.cert_path),
		PathBuf::from(args.key_path),
	)
		.await
		.unwrap();

	let sockets_v4: Vec<SocketAddr> = args.ipv4_addrs.iter().map(|addr| {
		let addr_v4 = addr.parse::<Ipv4Addr>().unwrap();
		SocketAddr::new(addr_v4.into(), args.port)
	})
		.collect();

	let sockets_v6: Vec<SocketAddr> = args.ipv6_addrs.iter().map(|addr| {
		let addr_v6 = addr.parse::<Ipv6Addr>().unwrap();
		SocketAddr::new(addr_v6.into(), args.port)
	})
		.collect();

	let req_counter = AtomicU64::new(0);
	let req_counter_arc = Arc::new(req_counter);
	let server_config = ServerConfig{
		req_counter: req_counter_arc.clone(),
	};

	let app = Router::new()
		.route("/", get(handler))
		.with_state(server_config);

	// Run HTTPS server
	info!("listening on port {}", args.port);

	let server_handles_v4: Vec<tokio::task::JoinHandle<Result<(), Error>>> = sockets_v4
		.iter()
		.map(|socket_v4| {
			tokio::spawn(
				axum_server::bind_rustls(*socket_v4, config.clone())
					.serve(
						app.clone().into_make_service_with_connect_info::<SocketAddr>()
					)
			)
		})
		.collect();

	let server_handles_v6: Vec<tokio::task::JoinHandle<Result<(), Error>>> = sockets_v6
		.iter()
		.map(|socket_v6| {
			tokio::spawn(
				axum_server::bind_rustls(*socket_v6, config.clone())
					.serve(
						app.clone().into_make_service_with_connect_info::<SocketAddr>()
					)
			)
		})
		.collect();

	let start_time = Instant::now();
	let mut prev_elapsed_time = Duration::new(0,0);
	let mut prev_total_requests = 0;
	let mut interval = time::interval(Duration::from_secs(args.log_interval));
	interval.tick().await;
	loop {
		interval.tick().await;
		let total_requests = req_counter_arc.load(Ordering::Relaxed);
		let total_requests_diff = total_requests - prev_total_requests;
		let elapsed_time = start_time.elapsed() - prev_elapsed_time;
		let rps = total_requests_diff as f64 / elapsed_time.as_secs() as f64;
		let rps_tot = total_requests as f64 / start_time.elapsed().as_secs() as f64;

		info!("\nRequests per second: {:.2}\nTotal requests per second: {:.2}\nTotal requests: {}",
			rps, rps_tot, total_requests);

		prev_elapsed_time = elapsed_time;
		prev_total_requests = total_requests;
	}

	for handle_v4 in server_handles_v4 {
		handle_v4
			.await
			.unwrap()
			.unwrap();
	}

	for handle_v6 in server_handles_v6 {
		handle_v6
			.await
			.unwrap()
			.unwrap();
	}
}

async fn handler(
	ConnectInfo(addr): ConnectInfo<SocketAddr>,
	State(server_config): State<ServerConfig>
) -> String {
	// Increment the request counter
	server_config.req_counter.fetch_add(1, Ordering::SeqCst);

	format!("{}", addr.ip())
}
