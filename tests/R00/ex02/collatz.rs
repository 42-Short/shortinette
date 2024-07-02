use std::thread;
use std::time::Duration;
use std::sync::mpsc;
use std::panic;

fn main() {
    panic::set_hook(Box::new(|panic_info| {
        if let Some(s) = panic_info.payload().downcast_ref::<&str>() {
            println!("Panic occurred: {}", s);
        } else {
            println!("Panic occurred");
        }
    }));

    let (tx, rx) = mpsc::channel();
    let handle = thread::spawn(move || {
        collatz(0);
        tx.send(()).unwrap();
    });
    let timeout_duration = Duration::from_millis(1);
    if rx.recv_timeout(timeout_duration).is_err() {
        eprintln!("timeout");
        panic!("");
    }
    handle.join().unwrap();
}