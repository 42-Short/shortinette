#[cfg(test)]
mod shortinette_tests_rust_0206 {
    use super::*;

    fn test_next_token(s: &str, expected_tokens: &[Option<Token>]) {
        let mut input = s;
        for expected in expected_tokens {
            let token = next_token(&mut input);
            assert_eq!(token, *expected, "Failed to parse mixed tokens as expected");
        }
    }

    #[test]
    fn test_next_token_no_tokens() {
        let mut input = "";
        assert!(next_token(&mut input).is_none());
    }

    #[test]
    fn test_next_token_only_spaces() {
        let mut input = "                               ";
        assert!(next_token(&mut input).is_none());
    }

    #[test]
    fn test_next_token_one_word() {
        test_next_token("WORD", &[Some(Token::Word("WORD")), None]);
    }

    #[test]
    fn test_next_token_one_redirect() {
        test_next_token("<", &[Some(Token::RedirectStdin), None]);
    }

    #[test]
    fn test_next_token_mixed() {
        let expected_tokens = [
            Some(Token::Word("+++#]#-")),
            Some(Token::RedirectStdout),
            Some(Token::Word("WORD")),
            Some(Token::Pipe),
            Some(Token::RedirectStdin),
        ];
        test_next_token("+++#]#->WORD|<", &expected_tokens)
    }

    #[test]
    fn test_next_token_long_input() {
        let long_input = "+++#]#- ".repeat(1000);
        let mut input = &long_input[..];
        for _ in 0..1000 {
            assert_eq!(
                next_token(&mut input),
                Some(Token::Word("+++#]#-")),
                "Failed to parse long input correctly"
            );
        }
    }
}
