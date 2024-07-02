use std::thread;
use std::time::Duration;
use std::sync::mpsc;

fn main() {
    let (tx, rx) = mpsc::channel();
    let handle = thread::spawn(move || {
        collatz(1);
        tx.send(()).unwrap();
    });
    let timeout_duration = Duration::from_secs(1);
    if rx.recv_timeout(timeout_duration).is_err() {
        panic!("Timeout!");
    }
    handle.join().unwrap();
}
