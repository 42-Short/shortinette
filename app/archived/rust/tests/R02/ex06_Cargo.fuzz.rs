#![no_main]
use ex06::*;
use libfuzzer_sys::fuzz_target;

fn testing_next_token<'a>(s: &mut &'a str) -> Option<Token<'a>> {
    *s = s.trim_start();
    if s.is_empty() {
        return None;
    }

    let mut char_indices = s.char_indices();
    if let Some((_, c)) = char_indices.next() {
        match c {
            '<' => {
                *s = &s[1..];
                return Some(Token::RedirectStdin);
            }
            '>' => {
                *s = &s[1..];
                return Some(Token::RedirectStdout);
            }
            '|' => {
                *s = &s[1..];
                return Some(Token::Pipe);
            }
            _ => {
                let token_len = s
                    .find(|c: char| c.is_whitespace() || c == '<' || c == '>' || c == '|')
                    .unwrap_or_else(|| s.len());
                let token = &s[..token_len];
                *s = &s[token_len..];
                return Some(Token::Word(token));
            }
        }
    }
    None
}

fuzz_target!(|data: &[u8]| {
    let s = match std::str::from_utf8(data) {
        Ok(str) => str,
        Err(_) => return,
    };
    let mut input_str = s;
    let mut input_str_copy = s;
    loop {
        let token = next_token(&mut input_str);
        let correct_token = testing_next_token(&mut input_str_copy);
        assert_eq!(token, correct_token, "Tokens do not match");
        if token.is_none() && correct_token.is_none() {
            break;
        }
        assert_eq!(input_str, input_str_copy, "Input strings do not match");
    }
});
