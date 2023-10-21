use std::io;
use clap::Parser;
use may_minihttp::{HttpServer, HttpService, Request, Response};

/// A HTTP server that echoes your remote IP-address
#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
    /// Address to bind to
    #[arg(short, long, default_value = "[::]:11313")]
    address: String,
}

#[derive(Clone)]
struct WhatServer;

impl HttpService for WhatServer {
    fn call(&mut self, _req: Request, res: &mut Response) -> io::Result<()> {
        res.body("<remote-ip-address>");
        Ok(())
    }
}

fn main() {
    let args = Args::parse();

    println!("Serving at address {}", args.address);
    let server = HttpServer(WhatServer).start(args.address).unwrap();
    server.join().unwrap();
}
