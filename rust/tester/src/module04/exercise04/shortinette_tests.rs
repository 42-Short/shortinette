#[cfg(test)]
mod shortinette_tests {
    use ex04::*;
    use std::time::Instant;

    #[test]
    fn parallel_execution() {
        let cli1 = &[String::from("sleep"), String::from("1")];
        let cli2 = &[String::from("sleep"), String::from("2")];
        let command_lines: Vec<&[String]> = vec![cli1, cli2];

        let start_time = Instant::now();

        let mut output = Vec::new();

        if let Err(err) = multiplexer(&mut output, &command_lines) {
            panic!("call to multiplexer failed with error: {}.", err);
        }

        let elapsed = start_time.elapsed();

        assert!(elapsed.as_secs_f64() >= 2.0 && elapsed.as_secs_f64() <= 2.5,
        "Expected execution time ~2 seconds, got {:?} - are you sure the commands are running in parallel?", elapsed);

        let output_str = String::from_utf8(output).expect("Failed to parse output as UTF-8.");

        assert!(output_str.contains("sleep 1"));
        assert!(output_str.contains("sleep 2"));
    }

    #[test]
    fn subject() {
        let cli1: &[String] = &[String::from("echo"), String::from("a"), String::from("b")];
        let cli2: &[String] = &[String::from("sleep"), String::from("1")];
        let cli3: &[String] = &[String::from("cat"), String::from("Cargo.toml")];
        let command_lines = vec![cli1, cli2, cli3];

        let mut output = Vec::new();

        if let Err(err) = multiplexer(&mut output, &command_lines) {
            panic!("Call to multiplexer failed with error: {}", err);
        }

        let output_str = String::from_utf8(output).expect("Failed to parse output as UTF-8.");

        assert!(output_str.contains("echo a b"));
        assert!(output_str.contains("sleep 1"));
        assert!(output_str.contains("cat Cargo.toml"));
        assert!(output_str.contains("[package]"));
    }

    #[test]
    fn non_existing_command() {
        let cli: &[String] = &[String::from("DONOTPANIC")];
        let command_lines = vec![cli];

        let mut output = Vec::new();

        if let Err(err) = multiplexer(&mut output, &command_lines) {
            panic!("Call to multiplexer failed with error: {}.", err);
        }
    }

    #[test]
    fn failing_command() {
        let cli: &[String] = &[String::from("cat"), String::from("idonotexist.txt")];
        let command_lines = vec![cli];

        let mut output = Vec::new();

        if let Err(err) = multiplexer(&mut output, &command_lines) {
            panic!("Call to multiplexer failed with error: {}.", err);
        }

        let output_str = String::from_utf8(output).expect("Failed to parse output as UTF-8.");
        eprintln!("{}", output_str);

        assert!(output_str.contains("cat idonotexist.txt"));
    }
}
