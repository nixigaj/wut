use axum::{extract::ConnectInfo, routing::get, Router};
use axum_server::tls_rustls::{RustlsConfig};
use clap::Parser;
use std::{
	io::Error,
	net::{Ipv4Addr, Ipv6Addr, SocketAddr},
	path::PathBuf
};

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
	key_path: String
}

//noinspection DuplicatedCode
#[tokio::main]
async fn main() {
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

	let app = Router::new().route("/", get(handler));

	// Run HTTPS server
	println!("listening on port {}", args.port);

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

async fn handler(ConnectInfo(addr): ConnectInfo<SocketAddr>) -> String {
	format!("{}", addr.ip())
}
