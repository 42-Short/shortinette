#[cfg(test)]
mod shortinette_rust_test_module05_ex01_0001 {
    use super::*;

    use std::{
        collections::HashMap,
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
    fn threads() {
        let ex = Exercise::new();

        let child = Command::new("strace")
            .args(["-f", "-e", "trace=clone"])
            .arg(ex.path())
            .output_with_timeout(time::Duration::from_secs(1))
            .expect("Failed to execute exercise");

        let thread_count = String::from_utf8_lossy(&child.stderr)
            .lines()
            .filter(|line| line.contains("strace: Process ") && line.ends_with(" attached"))
            .count();

        assert_eq!(10, thread_count);
    }

    #[test]
    fn is_the_output_correct() {
        fn check_messages(messages: Vec<&str>) -> bool {
            let mut thread_messages: HashMap<u32, Vec<u32>> = HashMap::new();

            for message in messages {
                if !message.starts_with("hello ") || !message.ends_with('!') {
                    return false;
                }
                let middle_part = &message[6..message.len() - 1];
                let parts: Vec<&str> = middle_part.split(" from thread ").collect();

                if parts.len() != 2 {
                    return false;
                }
                let message_number = match parts[0].parse::<u32>() {
                    Ok(num) => num,
                    Err(_) => return false,
                };
                let thread_number = match parts[1].parse::<u32>() {
                    Ok(num) => num,
                    Err(_) => return false,
                };

                let entry = thread_messages.entry(thread_number).or_default();
                entry.push(message_number);
            }

            if thread_messages.len() != 10 {
                return false;
            }

            if thread_messages.iter().any(|(key, _)| key >= &10) {
                return false;
            }

            for message_numbers in thread_messages.values() {
                if message_numbers.len() != 10 {
                    return false;
                }

                let sum: u32 = message_numbers.iter().sum();
                if sum != 45 {
                    return false;
                }
            }

            true
        }

        let ex = Exercise::new();

        let output = ex
            .cmd()
            .output_with_timeout(time::Duration::from_secs(1))
            .expect("Failed to execute exercise");

        let out = String::from_utf8_lossy(&output.stdout);
        let is_ok = check_messages(out.lines().collect());

        assert!(
            is_ok,
            "The console output of your exercise is not correct:\n{out}"
        );
    }

    #[test]
    fn no_buffer_empty_message() {
        let mut out = Vec::new();
        let mut log = Logger::new(0, &mut out);

        log.log("").unwrap();
        assert_eq!(out, b"\n");
    }

    #[test]
    fn no_buffer_sinlge_letter() {
        let mut out = Vec::new();
        let mut logger = Logger::new(0, &mut out);

        logger.log("h").unwrap();
        assert_eq!(out, b"h\n");
    }

    #[test]
    fn no_buffer_sinle_word() {
        let mut out = Vec::new();
        let mut logger = Logger::new(0, &mut out);

        logger.log("hello").unwrap();
        assert_eq!(out, b"hello\n");
    }

    #[test]
    fn no_buffer_multiple_words() {
        let mut out = Vec::new();
        let mut logger = Logger::new(0, &mut out);

        logger.log("hello").unwrap();
        logger.log("world").unwrap();
        logger.log("testing").unwrap();
        logger.log("h").unwrap();
        logger.log("").unwrap();
        assert_eq!(out, b"hello\nworld\ntesting\nh\n\n");
    }

    #[test]
    fn buffered_empty_message() {
        let mut out = Vec::new();
        let mut logger = Logger::new(12, &mut out);

        logger.log("").unwrap();
        assert_eq!(logger.writer, b"");

        logger.flush().unwrap();
        assert_eq!(logger.writer, b"\n")
    }

    #[test]
    fn buffered_message() {
        let mut out = Vec::new();
        let mut logger = Logger::new(12, &mut out);

        logger.log("hello").unwrap();
        assert_eq!(logger.writer, b"");

        logger.flush().unwrap();
        assert_eq!(logger.writer, b"hello\n");
    }

    #[test]
    fn buffered_messages() {
        let mut out = Vec::new();
        let mut logger = Logger::new(12, &mut out);

        logger.log("hello").unwrap();
        assert_eq!(logger.writer, b"");

        logger.log("world").unwrap();
        assert_eq!(logger.writer, b"hello\nworld\n");
    }

    #[test]
    fn buffer_len_same_as_message() {
        let mut out = Vec::new();
        let mut logger = Logger::new(12, &mut out);

        logger.log("Hello World!").unwrap();
        assert_eq!(logger.writer, b"Hello World!\n");
    }

    #[test]
    fn buffer_len_same_as_message_with_newline() {
        let mut out = Vec::new();
        let mut logger = Logger::new(12, &mut out);

        logger.log("Hello World").unwrap();
        assert_eq!(logger.writer, b"Hello World\n");
    }

    #[test]
    fn buffer_too_long_message() {
        let mut out = Vec::new();
        let mut logger = Logger::new(1024, &mut out);

        logger.log(&"a".repeat(2048)).unwrap();
        assert_eq!(logger.writer, format!("{}\n", "a".repeat(2048)).as_bytes());
    }

    #[test]
    fn empty_buffer_flush() {
        let mut out = Vec::new();
        let mut logger = Logger::new(1024, &mut out);

        logger.flush().unwrap();
        assert!(logger.writer.is_empty());
    }

    #[test]
    fn correct_buffer_size() {
        let logger = Logger::new(0, vec![0; 0]);
        assert_eq!(0, logger.buffer.len());

        let logger = Logger::new(128, vec![0; 0]);
        assert_eq!(128, logger.buffer.len());

        let logger = Logger::new(1024, vec![0; 0]);
        assert_eq!(1024, logger.buffer.len());
    }
}
