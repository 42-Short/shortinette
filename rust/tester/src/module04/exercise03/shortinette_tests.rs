#[cfg(test)]
mod shortinette_tests {
    use ex03::*;

    #[test]
    fn basic() {
        let cli: &[String] = &[String::from("cat"), String::from("Cargo.toml")];

        let mut input: &[u8] = b"";
        let mut writer = Vec::new();

        if let Err(err) = pipeline(&mut input, &mut writer, cli) {
            panic!("pipeline() call failed with error: {}.", err)
        }

        let output_str = String::from_utf8(writer).expect("Could not parse output as UTF-8.");

        assert!(output_str.contains("[package]"));
    }

    #[test]
    fn stdin() {
        let cli: &[String] = &[String::from("cat")];

        let mut input: &[u8] = b"Cat me if you can!";
        let mut writer = Vec::new();

        if let Err(err) = pipeline(&mut input, &mut writer, cli) {
            panic!("pipeline() call failed with error: {}.", err)
        }

        let output_str = String::from_utf8(writer).expect("Could not parse output as UTF-8.");

        assert_eq!(
            output_str, "Cat me if you can!\n",
            "Don't forget to pipe the input to your process' standard input!"
        );
    }

    #[test]
    fn panic() {
        let cli: &[String] = &[String::from("idonotthinkthiscommandexists")];

        let mut input: &[u8] = b"";
        let mut writer = Vec::new();

        if let Ok(()) = pipeline(&mut input, &mut writer, cli) {
            panic!("Expected Err on a non-existing command.");
        }
    }
}
