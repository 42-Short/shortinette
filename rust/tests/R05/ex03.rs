#[cfg(test)]
mod shortinette_rust_test_module05_ex03_0001 {
    use std::{
        env, ffi,
        io::{self, Write},
        path::PathBuf,
        process::{self, Command, Output},
        thread,
        time::{self, Duration},
    };

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
            // let mut command = Command::new("cargo");
            // command.args(["run", "--release"]);

            // command
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
                // .stderr(process::Stdio::piped())
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

            // fs::rename(format!("./target/release/{}", Self::EXERCISE), path)
            //     .expect("Failed to move executable");

            // fs::remove_dir_all("./target").expect("Could not delete generated files");
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

    trait ChildWaitTimeout {
        fn wait_with_timeout(self, timeout: time::Duration) -> io::Result<Output>;
    }

    impl ChildWaitTimeout for process::Child {
        fn wait_with_timeout(mut self, timeout: time::Duration) -> io::Result<Output> {
            let start = time::Instant::now();

            loop {
                match self.try_wait() {
                    Ok(Some(_)) => return self.wait_with_output(),
                    Ok(None) => {
                        if start.elapsed() >= timeout {
                            _ = self.kill();

                            return Err(io::Error::new(
                                io::ErrorKind::TimedOut,
                                "Process timed out",
                            ));
                        }

                        thread::sleep(time::Duration::from_millis(10));
                    }
                    Err(err) => return Err(err),
                }
            }
        }
    }

    #[test]
    fn no_args() {
        let ex = Exercise::new();

        let output = ex.cmd().output_with_timeout(Duration::from_millis(100));
        // Could also be an execute fail, but more likely a timeout
        assert!(output.is_ok());
    }

    #[test]
    fn never_gonna_end() {
        let ex = Exercise::new();

        let child = ex
            .cmd()
            .arg("1")
            .stdin(process::Stdio::piped())
            .stdout(process::Stdio::piped())
            .spawn()
            .expect("Failed to execute command");

        let output = child.wait_with_timeout(time::Duration::from_secs(1));

        assert!(output.is_err());
        assert!(io::ErrorKind::TimedOut == output.unwrap_err().kind());
    }

    #[test]
    fn threads() {
        let ex = Exercise::new();

        let child = Command::new("strace")
            .args(["-f", "-e", "trace=none"])
            .arg(ex.path())
            .arg("0")
            .output()
            .expect("Failed to execute command");

        let thread_count = String::from_utf8_lossy(&child.stderr)
            .lines()
            .filter(|line| line.contains("strace: Process ") && line.ends_with(" attached"))
            .count();

        assert_eq!(1, thread_count);
    }

    #[test]
    fn zero_brain_size() {
        let ex = Exercise::new();

        // Since a channel is used even size 0 makes the philosopher
        // able to think about at least one topic.
        let mut child = ex.spawn_child_args(["0"]);
        thread::sleep(time::Duration::from_millis(1000));
        let mut stdin = child.stdin.take().expect("Stdin vanished");

        writeln!(&mut stdin, "cakes").unwrap();
        thread::sleep(time::Duration::from_millis(100));

        writeln!(&mut stdin, "code").unwrap();
        thread::sleep(time::Duration::from_millis(100));

        thread::sleep(time::Duration::from_millis(5000));

        drop(stdin);

        match child.wait_with_timeout(time::Duration::from_millis(100)) {
            Ok(out) => {
                let output = String::from_utf8_lossy(&out.stdout);
                let output = output.lines().collect::<Vec<_>>();

                let expected = vec![
                    "the philosopher is thinking about cakes",
                    "the philosopher's head is full",
                ];

                assert_eq!(expected, output);
            }
            Err(_) => {
                panic!("The philosopher died too slow")
            }
        }
    }

    #[test]
    fn one_brain_size() {
        let ex = Exercise::new();

        let mut child = ex.spawn_child_args(["1"]);
        thread::sleep(time::Duration::from_millis(1));
        let mut stdin = child.stdin.take().expect("Stdin vanished");

        writeln!(&mut stdin, "cakes").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        writeln!(&mut stdin, "code").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        writeln!(&mut stdin, "42").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        thread::sleep(time::Duration::from_millis(5000));

        drop(stdin);

        match child.wait_with_timeout(time::Duration::from_millis(1)) {
            Ok(out) => {
                let output = String::from_utf8_lossy(&out.stdout);
                let output = output.lines().collect::<Vec<_>>();

                let expected = vec![
                    "the philosopher is thinking about cakes",
                    "the philosopher's head is full",
                    "the philosopher is thinking about code",
                ];

                assert_eq!(expected, output);
            }
            Err(_) => {
                panic!("The philosopher died too slow")
            }
        }
    }

    #[test]
    fn huge_brain_size() {
        let ex = Exercise::new();

        let mut child = ex.spawn_child_args(["1000"]);
        thread::sleep(time::Duration::from_millis(1));
        let mut stdin = child.stdin.take().expect("Stdin vanished");

        writeln!(&mut stdin, "cakes").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        writeln!(&mut stdin, "code").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        writeln!(&mut stdin, "42").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        writeln!(&mut stdin, "21").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        writeln!(&mut stdin, "leet").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        writeln!(&mut stdin, "spaghetti").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        writeln!(&mut stdin, "who stole his fork").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        writeln!(&mut stdin, "starving").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        thread::sleep(time::Duration::from_secs(40));

        drop(stdin);

        match child.wait_with_timeout(time::Duration::from_millis(1)) {
            Ok(out) => {
                let output = String::from_utf8_lossy(&out.stdout);
                let output = output.lines().collect::<Vec<_>>();

                let expected = vec![
                    "the philosopher is thinking about cakes",
                    "the philosopher is thinking about code",
                    "the philosopher is thinking about 42",
                    "the philosopher is thinking about 21",
                    "the philosopher is thinking about leet",
                    "the philosopher is thinking about spaghetti",
                    "the philosopher is thinking about who stole his fork",
                    "the philosopher is thinking about starving",
                ];

                assert_eq!(expected, output);
            }
            Err(_) => {
                panic!("The philosopher died too slow")
            }
        }
    }

    #[test]
    fn subject() {
        let ex = Exercise::new();

        let mut child = ex.spawn_child_args(["3"]);
        thread::sleep(time::Duration::from_millis(1));
        let mut stdin = child.stdin.take().expect("Stdin vanished");

        writeln!(&mut stdin, "cakes").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        writeln!(&mut stdin, "code").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        writeln!(&mut stdin, "42").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        thread::sleep(time::Duration::from_millis(5000));

        writeln!(&mut stdin, "a").unwrap();
        thread::sleep(time::Duration::from_millis(1));
        writeln!(&mut stdin, "b").unwrap();
        thread::sleep(time::Duration::from_millis(1));
        writeln!(&mut stdin, "c").unwrap();
        thread::sleep(time::Duration::from_millis(1));

        thread::sleep(time::Duration::from_millis(15000));

        drop(stdin);

        match child.wait_with_timeout(time::Duration::from_millis(1)) {
            Ok(out) => {
                let output = String::from_utf8_lossy(&out.stdout);
                let output = output.lines().collect::<Vec<_>>();

                let expected = vec![
                    "the philosopher is thinking about cakes",
                    "the philosopher is thinking about code",
                    "the philosopher's head is full",
                    "the philosopher is thinking about 42",
                    "the philosopher is thinking about a",
                    "the philosopher is thinking about b",
                ];

                assert_eq!(expected, output);
            }
            Err(_) => {
                panic!("The philosopher died too slow")
            }
        }
    }
}
