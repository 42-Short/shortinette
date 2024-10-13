#[cfg(test)]
mod shortinette_rust_test_module05_ex06_0001 {
    use std::{
        env, ffi, io,
        net::TcpListener,
        panic::{self, AssertUnwindSafe},
        path::PathBuf,
        process::{self, Command, Output},
        thread,
        time::{self, Duration},
    };

    use crate::*;

    // TODO: This could maybe be moved to it's own module
    // But since this is a todo it will never happen
    struct Exercise;
    #[allow(dead_code)]
    impl Exercise {
        const EXERCISE: &'static str = env!("CARGO_PKG_NAME");

        fn new() -> Self {
            Self::compile();

            Self
        }

        fn cmd(&self) -> Command {
            Command::new(self.path())
        }

        fn path(&self) -> PathBuf {
            let mut path = PathBuf::new();
            path.push("./target/release/");
            path.push(Self::EXERCISE);

            path
        }

        fn spawn_child_args<I, S>(&self, args: I) -> process::Child
        where
            I: IntoIterator<Item = S>,
            S: AsRef<ffi::OsStr>,
        {
            self.cmd()
                .args(args)
                .stdin(process::Stdio::piped())
                .stdout(process::Stdio::piped())
                .spawn()
                .expect("Failed to execute command")
        }

        // TODO: This creates a side effect which could interfear with other tests
        // Maybe should instead create a directory in /tmp
        fn compile() {
            let mut path = PathBuf::new();
            path.push("./target/release/");
            path.push(Self::EXERCISE);

            if path.exists() {
                return;
            }

            let output = Command::new("cargo")
                .args(["build", "--release", "--target-dir", "./target"])
                .output()
                .expect("Failed to build exercise");

            if !output.status.success() {
                panic!("Failed to build exercise");
            }
        }
    }

    trait CommandOutputTimeout {
        fn output_with_timeout(&mut self, timeout: time::Duration) -> io::Result<Output>;
    }

    impl CommandOutputTimeout for Command {
        fn output_with_timeout(&mut self, timeout: time::Duration) -> io::Result<Output> {
            let start = time::Instant::now();
            let child = self
                .stdout(process::Stdio::piped())
                .stderr(process::Stdio::piped())
                .spawn();

            let mut child = match child {
                Ok(child) => child,
                Err(err) => return Err(err),
            };

            loop {
                match child.try_wait() {
                    Ok(Some(_)) => return child.wait_with_output(),
                    Ok(None) => {
                        if start.elapsed() >= timeout {
                            _ = child.kill();

                            return Err(io::Error::new(
                                io::ErrorKind::TimedOut,
                                "Command timed out",
                            ));
                        }

                        thread::sleep(time::Duration::from_millis(10));
                    }
                    Err(err) => return Err(err),
                }
            }
        }
    }

    fn get_random_port() -> u16 {
        let listener = TcpListener::bind("127.0.0.1:0").expect("Failed to start port discovery.");

        listener
            .local_addr()
            .expect("Failed to get port discovery local address")
            .port()
    }

    #[test]
    fn no_args() {
        let ex = Exercise::new();

        let output = ex.cmd().output_with_timeout(time::Duration::from_millis(1));
        assert!(output.is_ok());
    }

    // #[test]
    // fn is_multithreaded() {
    //     let ex = Exercise::new();

    //     let mut child = Command::new("strace")
    //         // `--kill-on-exit` is not available on strace 6.1, which is used by bookworm
    //         // Because of that, when we kill the `strace`, it will not kill the http server
    //         // Which will cause the test to run forever
    //         // And attaching `strace` to a `pid` is not viable either
    //         .args(["--kill-on-exit", "-f", "-e", "trace=none"])
    //         .arg(ex.path())
    //         .arg(format!("127.0.0.1:{}", get_random_port()))
    //         .stderr(process::Stdio::piped())
    //         .spawn()
    //         .expect("Failed to run server");

    //     thread::sleep(time::Duration::from_millis(500));
    //     child.kill().unwrap();

    //     let out = child.wait_with_output().unwrap();
    //     let thread_count = String::from_utf8_lossy(&out.stderr)
    //         .lines()
    //         .filter(|line| line.starts_with("strace: Process ") && line.ends_with(" attached"))
    //         .count();

    //     assert!(
    //         thread_count >= 2,
    //         "Your http server should spawn at least 2 threads"
    //     );
    // }

    #[test]
    fn basic_curl() {
        let ex = Exercise::new();

        let address = format!("127.0.0.1:{}", get_random_port());
        let mut server = ex
            .cmd()
            .arg(&address)
            .spawn()
            .expect("Failed to run server");

        let curl = Command::new("curl")
            .args(["--http1.1", "-I", "-X", "GET", &address])
            .output_with_timeout(Duration::from_millis(100));

        server.kill().expect("Could not stop server");

        // Expect needs to be here.
        // We should only panic, after the child got killed.
        // Before could make the child a zombie, we don't want that.
        let curl = curl.expect("Failed to run curl");
        assert!(curl.status.success());

        let out = String::from_utf8_lossy(&curl.stdout);
        let mut lines = out.lines();

        assert_eq!("HTTP/1.1 404 Not Found", lines.next().unwrap());
    }

    #[test]
    fn multiple_curls() {
        let ex = Exercise::new();

        let address = format!("127.0.0.1:{}", get_random_port());
        let mut server = ex
            .cmd()
            .arg(&address)
            .spawn()
            .expect("Failed to run server");

        // Hacky solution but should work.
        // In case any of the test fails, it would result in a panic.
        // A panic would exit the test early causing the child to not get killed properly
        let result = panic::catch_unwind(AssertUnwindSafe(|| {
            let threads = (0..1000).map(|_| {
                let address = address.clone();

                thread::spawn(move || {
                    let curl = Command::new("curl")
                        .args(["--http1.1", "-I", "-X", "GET", &address])
                        .output_with_timeout(Duration::from_millis(100))
                        .expect("Failed to run curl");

                    assert!(curl.status.success());

                    let out = String::from_utf8_lossy(&curl.stdout);
                    let mut lines = out.lines();

                    assert_eq!("HTTP/1.1 404 Not Found", lines.next().unwrap());
                })
            });

            threads.for_each(|thread| {
                thread.join().unwrap();
            });
        }));

        server.kill().expect("Could not stop server");

        if let Err(err) = result {
            panic::resume_unwind(err);
        }
    }

    #[test]
    fn thread_pool_zero_threads() {
        let pool = ThreadPool::new(0);

        assert_eq!(0, pool.threads.len());

        assert!(pool.spawn_task(|| {}).is_err());
    }

    #[test]
    fn thread_pool_one_thread() {
        let pool = ThreadPool::new(1);
        assert_eq!(1, pool.threads.len());

        let num = Arc::new(Mutex::new(0));
        let now = time::Instant::now();

        for _ in 0..2 {
            let num = num.clone();
            pool.spawn_task(move || {
                thread::sleep(time::Duration::from_millis(100));
                *num.lock().unwrap() += 1;
            })
            .unwrap();
        }

        while *num.lock().unwrap() != 2 && now.elapsed() < time::Duration::from_millis(290) {
            thread::sleep(time::Duration::from_millis(1));
        }

        let elapsed = now.elapsed();
        assert!(elapsed < time::Duration::from_millis(250));
        assert!(time::Duration::from_millis(150) < elapsed);
    }

    #[test]
    fn thread_pool_two_threads() {
        let pool = ThreadPool::new(2);
        assert_eq!(2, pool.threads.len());

        let num = Arc::new(Mutex::new(0));
        let now = time::Instant::now();

        for _ in 0..2 {
            let num = num.clone();
            pool.spawn_task(move || {
                thread::sleep(time::Duration::from_millis(100));
                *num.lock().unwrap() += 1;
            })
            .unwrap();
        }

        while *num.lock().unwrap() != 2 && now.elapsed() < time::Duration::from_millis(150) {
            thread::sleep(time::Duration::from_millis(1));
        }
        assert!(now.elapsed() < time::Duration::from_millis(150));
        assert!(now.elapsed() >= time::Duration::from_millis(100));
    }

    #[test]
    fn thread_pool_distribution() {
        let pool = ThreadPool::new(100);
        assert_eq!(100, pool.threads.len());

        let num = Arc::new(Mutex::new(0));
        let now = time::Instant::now();

        for _ in 0..100 {
            let num = num.clone();
            pool.spawn_task(move || {
                thread::sleep(time::Duration::from_millis(100));
                *num.lock().unwrap() += 1;
            })
            .unwrap();
        }

        while *num.lock().unwrap() != 100 && now.elapsed() < time::Duration::from_millis(150) {
            thread::sleep(time::Duration::from_millis(1));
        }
        assert!(now.elapsed() < time::Duration::from_millis(190));
        assert!(now.elapsed() >= time::Duration::from_millis(100));
    }

    #[test]
    fn thread_pool_many_tasks() {
        let pool = ThreadPool::new(10);
        assert_eq!(10, pool.threads.len());

        let num = Arc::new(Mutex::new(0));
        let now = time::Instant::now();

        for _ in 0..100 {
            let num = num.clone();
            pool.spawn_task(move || {
                thread::sleep(time::Duration::from_millis(100));
                *num.lock().unwrap() += 1;
            })
            .unwrap();
        }

        while *num.lock().unwrap() != 100 && now.elapsed() < time::Duration::from_millis(1500) {
            thread::sleep(time::Duration::from_millis(1));
        }
        assert!(now.elapsed() < time::Duration::from_millis(1100));
        assert!(now.elapsed() >= time::Duration::from_millis(1000));
    }

    #[test]
    fn thread_pool_spawn_error() {
        let mut pool = ThreadPool::new(5);

        // Artificially drain the pool this should work with any implementation
        {
            let mut should_stop = pool.should_stop.write().expect("Poisoned lock");
            *should_stop = true;
        }

        // Send a no-op task to every thread, to ensure every thread received the `should_stop`
        for _ in &pool.threads {
            pool.spawn_task(|| {}).ok();
        }

        for thread in pool.threads.drain(..) {
            if thread.join().is_err() {
                panic!("A thread in ThreadPool panicked")
            }
        }

        // Since all threads have been closed, there should no receiver be available
        assert!(pool.spawn_task(|| {}).is_err());
    }

    #[test]
    #[should_panic]
    fn thread_pool_panic_on_drop() {
        let pool = ThreadPool::new(5);

        pool.spawn_task(|| panic!("This task should panic"))
            .expect("Failed to send task");

        // Just to explicitly show that the test is about `drop`ping the pool
        drop(pool);
    }

    #[test]
    fn thread_pool_all_threads_panic() {
        let pool = ThreadPool::new(5);

        for _ in 0..5 {
            pool.spawn_task(|| panic!("This task should panic"))
                .expect("Failed to send task");
        }

        thread::sleep(time::Duration::from_millis(100));

        assert!(pool.spawn_task(|| {}).is_err());

        // Dropping the pool should panic.
        // Hacky way of not letting main thread panic
        thread::spawn(move || drop(pool))
            .join()
            .expect_err("Why did dropping pool not panic?");
    }
}
