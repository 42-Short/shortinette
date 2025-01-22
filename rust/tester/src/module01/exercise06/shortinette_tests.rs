#[cfg(test)]
mod shortinette_tests_rust_0106 {
    use ex06::big_add;
    use rand::Rng;

    #[test]
    #[should_panic]
    fn test_both_empty() {
        big_add(b"", b"");
    }

    #[test]
    #[should_panic]
    fn test_a_empty() {
        let mut rng = rand::thread_rng();

        let number: Vec<u8> = (0..5).map(|_| rng.gen_range(0..10)).collect();
        big_add(b"", &number);
    }

    #[test]
    #[should_panic]
    fn test_b_empty() {
        let mut rng = rand::thread_rng();

        let number: Vec<u8> = (0..5).map(|_| rng.gen_range(0..10)).collect();
        big_add(&number, b"");
    }

    fn format_characters(characters: &[u8]) -> String {
        let mut output = String::from("\"");

        for &char in characters {
            if char.is_ascii_graphic() || char.is_ascii_whitespace() {
                output.push(char as char);
            } else {
                output.push_str(&format!("\\{}", char));
            }
        }
        output.push('"');
        output
    }

    fn big_add_wrapper(a: &[u8], b: &[u8]) {
        let result = std::panic::catch_unwind(|| {
            big_add(a, b);
        });
        match result {
            Ok(_) => panic!(
                "Invalid input not correctly handled\nInput: {} + {}",
                format_characters(a),
                format_characters(b)
            ),
            Err(payload) => {
                let payload_str = match payload.downcast_ref::<String>() {
                    Some(s) => s.as_str(),
                    None => match payload.downcast_ref::<&str>() {
                        Some(s) => s,
                        None => return,
                    },
                };
                if payload_str.contains("subtract with overflow") {
                    panic!(
                        "Invalid input not correctly handled\nInput: {} + {}",
                        format_characters(a),
                        format_characters(b)
                    );
                }
            }
        }
    }

    #[test]
    fn test_invalid_bytes() {
        let mut rng = rand::thread_rng();
        let test_chars = ['/', ':', ' ', '+', '-', '\0'];
        let mut input_strings = Vec::new();

        for character in test_chars {
            let number: Vec<u8> = (0..5).map(|_| rng.gen_range(b'0'..=b'9')).collect();

            let mut prefixed = vec![character as u8];
            prefixed.extend(&number);

            let mut appended = number.clone();
            appended.push(character as u8);

            input_strings.push(prefixed);
            input_strings.push(appended);
        }

        for input in &input_strings {
            let validnbr: Vec<u8> = (0..5).map(|_| rng.gen_range(b'0'..=b'9')).collect();
            big_add_wrapper(&validnbr, input);
            big_add_wrapper(input, &validnbr);
        }
    }

    // Subtract `a` from the received result and check if its equal to `b`, so I can check the result
    // for randomized inputs without having the solution itself in the repo
    fn big_sub(a: &[u8], b: &[u8], output: &[u8]) -> bool {
        if !output.iter().all(|val| val.is_ascii_digit()) {
            return false;
        }

        let mut result = Vec::new();
        let mut rest = 0;
        let max_len = output.len().max(a.len());
        for i in 0..max_len {
            let digit1 = if i < output.len() {
                output[output.len() - 1 - i] - b'0'
            } else {
                0
            };

            let digit2 = if i < a.len() {
                a[a.len() - 1 - i] - b'0'
            } else {
                0
            };

            let mut diff = digit1 as i8 + rest as i8 - digit2 as i8;
            if diff < 0 {
                diff += 10;
                rest = -1;
            } else {
                rest = 0;
            }

            result.insert(0, diff as u8 + b'0');
        }

        while result.len() > 1 && result[0] == b'0' {
            result.remove(0);
        }

        let mut b = b.to_vec();
        while b.len() > 1 && b[0] == b'0' {
            b.remove(0);
        }

        result == b
    }

    #[test]
    fn test_only_zeros() {
        let mut rng = rand::thread_rng();

        let number1: Vec<u8> = (0..rng.gen_range(5..15)).map(|_| b'0').collect();
        let number2: Vec<u8> = (0..rng.gen_range(5..15)).map(|_| b'0').collect();

        assert_eq!(
            big_add(&number1, &number2),
            b"0",
            "Failed for ({}, {})",
            format_characters(&number1),
            format_characters(&number2)
        );
    }

    #[test]
    fn test_leading_zeros() {
        let mut rng = rand::thread_rng();

        let number1: Vec<u8> = (0..rng.gen_range(5..15)).map(|_| b'0').collect();
        let mut number2: Vec<u8> = (0..rng.gen_range(5..15)).map(|_| b'0').collect();
        let additional = rng.gen_range(b'0'..=b'9');
        number2.push(additional);

        assert_eq!(
            big_add(&number1, &number2),
            &[additional],
            "Failed for ({}, {})",
            format_characters(&number1),
            format_characters(&number2)
        );
    }

    #[test]
    fn test_short_numbers() {
        let mut rng = rand::thread_rng();

        let number1: Vec<u8> = (0..rng.gen_range(1..5))
            .map(|_| rng.gen_range(b'0'..=b'9'))
            .collect();
        let number2: Vec<u8> = (0..rng.gen_range(1..5))
            .map(|_| rng.gen_range(b'0'..=b'9'))
            .collect();

        let outcome = big_add(&number1, &number2);
        assert!(
            big_sub(&number1, &number2, &outcome),
            "Failed for ({}, {})",
            format_characters(&number1),
            format_characters(&number2)
        );
    }

    #[test]
    fn test_long_numbers() {
        let mut rng = rand::thread_rng();

        let number1: Vec<u8> = (0..rng.gen_range(50..100))
            .map(|_| rng.gen_range(b'0'..=b'9'))
            .collect();
        let number2: Vec<u8> = (0..rng.gen_range(50..100))
            .map(|_| rng.gen_range(b'0'..=b'9'))
            .collect();

        let outcome = big_add(&number1, &number2);
        assert!(
            big_sub(&number1, &number2, &outcome),
            "Failed for ({}, {})",
            format_characters(&number1),
            format_characters(&number2)
        );
    }
}
