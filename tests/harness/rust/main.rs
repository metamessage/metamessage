/// MetaMessage Rust test harness - parse JSONC file and re-print to JSONC.
use std::env;
use std::fs;
use std::io::Read;
use std::process;

use metamessage::{parse_jsonc, to_jsonc_string, encode_from_jsonc, decode_to_jsonc};

fn bytes_to_hex(bytes: &[u8]) -> String {
    bytes.iter().map(|b| format!("{:02x}", b)).collect()
}

fn hex_to_bytes(hex: &str) -> Result<Vec<u8>, String> {
    let hex = hex.trim();
    if hex.len() % 2 != 0 {
        return Err("hex string must have even length".to_string());
    }
    (0..hex.len())
        .step_by(2)
        .map(|i| {
            u8::from_str_radix(&hex[i..i + 2], 16)
                .map_err(|e| format!("hex decode error: {}", e))
        })
        .collect()
}

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        eprintln!("usage: harness [--encode|--decode] <file.jsonc>");
        process::exit(1);
    }

    if args[1] == "--encode" {
        if args.len() < 3 {
            eprintln!("usage: harness --encode <file.jsonc>");
            process::exit(1);
        }
        let data = match fs::read_to_string(&args[2]) {
            Ok(d) => d,
            Err(e) => {
                eprintln!("read error: {}", e);
                process::exit(1);
            }
        };
        let wire = match encode_from_jsonc(&data) {
            Ok(w) => w,
            Err(e) => {
                eprintln!("encode error: {}", e);
                process::exit(1);
            }
        };
        print!("{}", bytes_to_hex(&wire));
        return;
    }

    if args[1] == "--decode" {
        let mut hex_str = String::new();
        if std::io::stdin().read_to_string(&mut hex_str).is_err() {
            eprintln!("error reading hex from stdin");
            process::exit(1);
        }
        let wire = match hex_to_bytes(&hex_str) {
            Ok(w) => w,
            Err(e) => {
                eprintln!("hex decode error: {}", e);
                process::exit(1);
            }
        };
        let output = match decode_to_jsonc(&wire) {
            Ok(s) => s,
            Err(e) => {
                eprintln!("decode error: {}", e);
                process::exit(1);
            }
        };
        print!("{}", output);
        return;
    }

    // Existing behavior
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