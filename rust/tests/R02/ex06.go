package R02

import (
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"

	"os"
	"path/filepath"

	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var cargoTestModAsString06 = `

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

`

var cargoFuzzAsSting = `

#![no_main]
use libfuzzer_sys::fuzz_target;
use ex06::*;

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
                let token_len = s.find(|c: char| c.is_whitespace() || c == '<' || c == '>' || c == '|')
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

`

var clippyTomlAsString06 = ``

func writeStringToFile(source string, destFilePath string) error {
	destFile, err := os.OpenFile(destFilePath, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err = destFile.WriteString(source); err != nil {
		return err
	}
	return nil
}

func runCargoFuzz(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	if _, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"fuzz", "init"}); err != nil {
		logger.Exercise.Printf("could not initialize cargo fuzz: %v", err)
		return Exercise.InternalError(err.Error())
	}
	if _, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"fuzz", "add", "next_token_fuzz"}); err != nil {
		logger.Exercise.Printf("could not add target to fuzz targets: %v", err)
		return Exercise.InternalError(err.Error())
	}
	if err := writeStringToFile(cargoFuzzAsSting, filepath.Join(workingDirectory, "fuzz/fuzz_targets/next_token_fuzz.rs")); err != nil {
		logger.Exercise.Printf("could not write to fuzz file: %v", err)
		return Exercise.InternalError(err.Error())
	}
	if _, err := testutils.RunCommandLine(workingDirectory, "rustup", []string{"override", "set", "nightly"}); err != nil {
		logger.Exercise.Printf("cant configure nightly toolchain: %v", err)
		return Exercise.InternalError(err.Error())
	}
	if _, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"fuzz", "run", "next_token_fuzz", "--", "-max_total_time=10"}); err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex06Test(exercise *Exercise.Exercise) Exercise.Result {
	if result := runDefaultTest(exercise, cargoTestModAsString06, clippyTomlAsString06, map[string]int{"unsafe": 0}); result.Passed {
		return result
	}
	return runCargoFuzz(exercise)
}

func ex06() Exercise.Exercise {
	return Exercise.NewExercise("06", "ex06", []string{"src/main.rs", "src/lib.rs", "Cargo.toml"}, 15, ex06Test)
}
