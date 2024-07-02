
use std::thread;
use std::time::Duration;
#[cfg(test)]
mod tests {
    use super::*;
    use std::sync::mpsc;
    
    #[test]
    fn test_yes_in_thread() {
        let (tx, rx) = mpsc::channel();

        let handle = thread::spawn(move || {
            tx.send(()).unwrap();
            yes();
        });

        rx.recv().unwrap();

        thread::sleep(Duration::from_secs(1));

        handle.thread().unpark();
    }
}