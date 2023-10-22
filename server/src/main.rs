use axum::{extract::ConnectInfo, routing::get, Router};
use axum_server::tls_rustls::RustlsConfig;
use clap::Parser;
use std::{
	net::{Ipv4Addr, Ipv6Addr, SocketAddr},
	path::PathBuf
};
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

/// A HTTP server that echoes your remote IP-address
#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
	/// IPv4 address to bind to
	#[arg(short = '4', long = "ipv4")]
	ipv4_addr: String,

	/// IPv6 address to bind to
	#[arg(short = '6', long = "ipv6")]
	ipv6_addr: String,

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

#[tokio::main]
async fn main() {
	let args = Args::parse();

	tracing_subscriber::registry()
		.with(
			tracing_subscriber::EnvFilter::try_from_default_env()
				.unwrap_or_else(|_| "example_tls_rustls=debug".into()),
		)
		.with(tracing_subscriber::fmt::layer())
		.init();

	// configure certificate and private key used by https
	let config = RustlsConfig::from_pem_file(
		PathBuf::from(args.cert_path),
		PathBuf::from(args.key_path),
	)
		.await
		.unwrap();

	let addr_v4 = args.ipv4_addr.parse::<Ipv4Addr>().unwrap();
	let socket_v4 = SocketAddr::new(addr_v4.into(), args.port);

	let addr_v6 = args.ipv6_addr.parse::<Ipv6Addr>().unwrap();
	let socket_v6 = SocketAddr::new(addr_v6.into(), args.port);

	let app = Router::new().route("/", get(handler));

	// run https server
	println!("listening on {}, {}, port {}", addr_v4, addr_v6, args.port);

	let server_handle_v4 = tokio::spawn(
		axum_server::bind_rustls(socket_v4, config.clone())
			.serve(
				app.clone().into_make_service_with_connect_info::<SocketAddr>()
			)
	);

	let server_handle_v6 = tokio::spawn(
		axum_server::bind_rustls(socket_v6, config)
			.serve(
				app.into_make_service_with_connect_info::<SocketAddr>()
			)
	);

	server_handle_v4.await.unwrap().unwrap();
	server_handle_v6.await.unwrap().unwrap();
}

async fn handler(ConnectInfo(addr): ConnectInfo<SocketAddr>) -> String {
	format!("{}", addr.ip())
}
