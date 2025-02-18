mod collatz;
mod print_bytes;
mod yes;

use base64::{engine::general_purpose::STANDARD, Engine as _};
use std::env;

#[cfg(test)]
mod shortinette_tests_0002 {
    use base64::{engine::general_purpose::STANDARD, Engine as _};
    use rand::{distributions::Alphanumeric, Rng};
    use std::io::Read;
    use std::process::{Command, Stdio};
    use std::time::Duration;
    use wait_timeout::ChildExt;

    fn format_non_printable(s: &str) -> String {
        s.chars()
            .map(|c| match c {
                '\0' => format!("\\0"),
                '\n' => format!("\\n"),
                _ => c.to_string(),
            })
            .collect()
    }

    fn print_bytes_helper(input: &str) {
        let encoded = STANDARD.encode(input);

        let run_output = Command::new("target/release/shortinette-test-module00-ex02")
            .arg("print_bytes")
            .arg(encoded)
            .output()
            .expect("Failed to execute program with print_bytes() function");

        if !run_output.stderr.is_empty() {
            panic!(
                "Unexpected content on stderr: {}",
                String::from_utf8_lossy(&run_output.stderr)
            );
        }

        let output = String::from_utf8_lossy(&run_output.stdout);
        let lines: Vec<_> = output
            .lines()
            .map(|line| line.parse::<u8>().expect("Unable to parse line into u8"))
            .collect();

        assert_eq!(
            input.as_bytes(),
            lines,
            "Output doesn't match with input '{}'",
            format_non_printable(input)
        );
    }

    #[test]
    fn test_print_bytes_subject() {
        print_bytes_helper("Déjà Vu\n");
    }

    #[test]
    fn test_print_bytes_empty() {
        print_bytes_helper("");
    }

    #[test]
    fn test_print_bytes_0_byte() {
        print_bytes_helper("\0");
    }

    fn random_string(size: usize) -> String {
        let mut rng = rand::thread_rng();
        (0..size)
            .map(|_| rng.sample(Alphanumeric))
            .map(char::from)
            .collect()
    }

    #[test]
    fn test_print_bytes_random() {
        let randstring = random_string(32);

        print_bytes_helper(&randstring);
    }

    #[test]
    fn test_print_bytes_random_with_0_byte() {
        let mut randstring = random_string(32);
        randstring.push('\0');

        print_bytes_helper(&randstring);
    }

    #[test]
    fn test_print_bytes_random_unicode() {
        let mut randstring = random_string(32);
        randstring.push('🦀');

        print_bytes_helper(&randstring);
    }

    fn execute_with_timeout(arguments: &[&str], timeout_expected: bool) -> String {
        let mut child = Command::new("target/release/shortinette-test-module00-ex02")
            .args(arguments)
            .stdout(Stdio::piped())
            .stderr(Stdio::piped())
            .spawn()
            .expect("Unable to execute program");

        let result = child.wait_timeout(Duration::from_secs(3)).unwrap();

        if timeout_expected && result.is_some() {
            panic!("Program exited too early");
        }

        if result.is_none() {
            child.kill().expect("Unable to kill process");
            child.wait().expect("Unable to wait for process");
        }

        if !timeout_expected && result.is_none() {
            panic!("Program didn't finish in time");
        }

        let mut stderr = String::new();
        child
            .stderr
            .take()
            .expect("Unable to get stderr from program")
            .read_to_string(&mut stderr);

        if !stderr.is_empty() {
            panic!("Unexpected content on stderr: {}", stderr);
        }

        let mut output = String::new();
        child
            .stdout
            .take()
            .expect("Unable to get stdout from program")
            .read_to_string(&mut output);

        output
    }

    fn validate_collatz_output(input: u32, output: String) {
        let mut expected = input;

        for (index, line) in output.lines().enumerate() {
            let received = line
                .parse::<u32>()
                .expect("Unable to parse output into u32");

            assert_eq!(
                received,
                expected,
                "Incorrect output in Line {}, Expected {}, Got {}",
                index + 1,
                expected,
                received
            );

            if expected % 2 == 0 {
                expected /= 2;
            } else if expected != 1 {
                expected = expected * 3 + 1;
            }
        }

        assert_eq!(expected, 1, "Empty or incomplete output");
    }

    #[test]
    fn test_yes() {
        let output = execute_with_timeout(&["yes"], true);

        if output.is_empty() {
            panic!("Empty output from yes() function");
        }

        for (index, line) in output.lines().enumerate() {
            if line != "y" {
                panic!("Invalid content on line {}: {}", index + 1, line);
            }
        }

        if output.lines().count() < 10000 {
            panic!("Expected more lines of output");
        }
    }

    #[test]
    fn test_collatz_0_endless_loop() {
        let output = execute_with_timeout(&["collatz", "0"], false);

        assert!(output.is_empty(), "Invalid output with collatz(0)");
    }

    #[test]
    fn test_collatz_1() {
        let output = execute_with_timeout(&["collatz", "1"], false);

        assert_eq!(output, "1\n", "Invalid output with collatz(1)");
    }

    #[test]
    fn test_collatz_random() {
        let mut rng = rand::thread_rng();

        let input = rng.gen_range(1000..10000);
        let arg = format!("{}", input);
        let output = execute_with_timeout(&["collatz", &arg], false);

        validate_collatz_output(input, output);
    }
}

fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() < 2 {
        panic!("Not enough arguments");
    }

    match args[1].as_str() {
        "yes" => yes::yes(),
        "print_bytes" => {
            if args.len() < 3 {
                panic!("Not enough arguments");
            }

            let decoded = STANDARD.decode(&args[2]).expect("Unable to decode base64");
            print_bytes::print_bytes(&String::from_utf8_lossy(&decoded));
        }
        "collatz" => {
            if args.len() < 3 {
                panic!("Not enough arguments");
            }
            collatz::collatz(
                args[2]
                    .parse::<u32>()
                    .expect("Unable to parse collatz value into u32"),
            );
        }
        _ => panic!("Invalid command"),
    }
}
