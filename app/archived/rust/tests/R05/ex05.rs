#[cfg(test)]
mod shortinette_rust_test_module05_ex05_0001 {
    use std::{
        env, ffi, io,
        path::PathBuf,
        process::{self, Command, Output},
        thread, time,
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

    #[test]
    fn no_args() {
        let ex = Exercise::new();

        let output = ex.cmd().output_with_timeout(time::Duration::from_millis(1));
        // Could also be an execute fail, but more likely a timeout
        assert!(output.is_ok());
    }

    #[test]
    fn zero() {
        let ex = Exercise::new();

        let output = ex
            .cmd()
            .arg("0")
            .output_with_timeout(time::Duration::from_millis(10))
            .expect("Failed to run command");

        let out = String::from_utf8_lossy(&output.stdout);
        let pi = out
            .lines()
            .find(|line| line.starts_with("pi: "))
            .map(|pi| pi[3..].trim());

        if let Some(pi) = pi {
            assert!(!pi.contains("NaN"), "Pi cannot be NaN");
        }
    }

    #[test]
    fn one() {
        let ex = Exercise::new();

        let output = ex
            .cmd()
            .arg("1")
            .output_with_timeout(time::Duration::from_millis(100))
            .expect("Failed to run command");

        let out = String::from_utf8_lossy(&output.stdout);
        let pi = out
            .lines()
            .find(|line| line.starts_with("pi: "))
            .map(|pi| pi[3..].trim())
            .expect("Where is pi?");

        assert!(pi == "0" || pi == "4");
    }

    #[test]
    fn two() {
        let ex = Exercise::new();

        let output = ex
            .cmd()
            .arg("2")
            .output_with_timeout(time::Duration::from_millis(100))
            .expect("Failed to run command");

        let out = String::from_utf8_lossy(&output.stdout);
        let pi = out
            .lines()
            .find(|line| line.starts_with("pi: "))
            .map(|pi| pi[3..].trim())
            .expect("Where is pi?");

        assert!(pi == "0" || pi == "2" || pi == "4");
    }

    #[test]
    fn three() {
        let ex = Exercise::new();

        let output = ex
            .cmd()
            .arg("3")
            .output_with_timeout(time::Duration::from_millis(100))
            .expect("Failed to run command");

        let out = String::from_utf8_lossy(&output.stdout);
        let pi = out
            .lines()
            .find(|line| line.starts_with("pi: "))
            .map(|pi| pi[3..].trim())
            .expect("Where is pi?")
            .parse::<f64>()
            .expect("Pi is not a number?");

        assert!(pi == 0.0 || pi == 1.3333333333333333 || pi == 2.6666666666666665 || pi == 4.0);
    }

    #[test]
    fn million() {
        let ex = Exercise::new();

        let output = ex
            .cmd()
            .arg("1000000")
            .output_with_timeout(time::Duration::from_secs(1))
            .expect("Failed to calculate pi. Is your progam too slow?");

        let out = String::from_utf8_lossy(&output.stdout);
        let pi = out
            .lines()
            .find(|line| line.starts_with("pi: "))
            .map(|pi| pi[3..].trim())
            .expect("Where is pi?")
            .parse::<f64>()
            .expect("Pi is not a number?");

        assert!((3.1..=3.2).contains(&pi));
    }

    #[test]
    #[allow(clippy::approx_constant)] // we do not want to use the precise pi
    fn billion() {
        let ex = Exercise::new();

        let now = time::Instant::now();
        let output = ex
            .cmd()
            .arg("1000000000")
            // Normally it should be even under 2 seconds
            .output_with_timeout(time::Duration::from_secs(5))
            .expect("Failed to calculate pi. Is your progam too slow?");

        // Practically it is impossible for it to be under 800 millis on normal pcs
        if now.elapsed() < time::Duration::from_millis(800) {
            panic!("Imagine returning a hardcoded pi");
        }

        let out = String::from_utf8_lossy(&output.stdout);
        let pi = out
            .lines()
            .find(|line| line.starts_with("pi: "))
            .map(|pi| pi[3..].trim())
            .expect("Where is pi?")
            .parse::<f64>()
            .expect("Pi is not a number?");

        assert!((3.141..=3.142).contains(&pi));
    }

    #[test]
    fn threads() {
        let ex = Exercise::new();

        let child = Command::new("strace")
            .args(["-f", "-e", "trace=none"])
            .arg(ex.path())
            .arg("1")
            .output()
            .expect("Failed to execute command");

        let thread_count = String::from_utf8_lossy(&child.stderr)
            .lines()
            .filter(|line| line.contains("strace: Process ") && line.ends_with(" attached"))
            .count();

        let cores = rayon::current_num_threads();
        assert_eq!(cores, thread_count);
    }
}
