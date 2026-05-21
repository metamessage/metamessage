/// MetaMessage Rust test harness - parse JSONC file and re-print to JSONC.
use std::env;
use std::fs;
use std::process;

use metamessage::{parse_jsonc, to_jsonc_string};

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        eprintln!("usage: harness <file.jsonc>");
        process::exit(1);
    }

    let data = match fs::read_to_string(&args[1]) {
        Ok(d) => d,
        Err(e) => {
            eprintln!("read error: {}", e);
            process::exit(1);
        }
    };

    let node = match parse_jsonc(&data) {
        Ok(n) => n,
        Err(e) => {
            eprintln!("parse error: {}", e);
            process::exit(1);
        }
    };

    let output = to_jsonc_string(&node);
    print!("{}", output);
}