package R03

import (
	"path/filepath"
	"time"

	"github.com/42-Short/shortinette/rust/alloweditems"

	"github.com/42-Short/shortinette/pkg/logger"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var Ex04TestMod = `
#[cfg(test)]
mod shortinette_rust_test_module03_ex04_0001 {
    use super::*;

    #[test]
    fn provided_tests() {
        let a: Time = "12:20".parse().unwrap();
        let b: Time = "15:14".parse().unwrap();

        assert_eq!(format!("{a}"), "12 hours, 20 minutes");
        assert_eq!(format!("{b}"), "15 hours, 14 minutes");

        let err1: TimeParseError = "12.20".parse::<Time>().unwrap_err();
        let err2: TimeParseError = "12:2".parse::<Time>().unwrap_err();
        let err3: TimeParseError = "12:2a".parse::<Time>().unwrap_err();
        assert_eq!(format!("error: {err1}"), "error: missing ':'");
        assert_eq!(format!("error: {err2}"), "error: invalid length");
        assert_eq!(format!("error: {err3}"), "error: invalid number");
    }

    fn parse_and_test_time(hours: u32, minutes: u32) {
        let raw = format!("{hours:0>2}:{minutes:0>2}");
        let time: Time = raw.parse().expect(&format!("Failed to parse {raw}"));

        assert_eq!(time.hours, hours, "Invalid hours for time {raw}");
        assert_eq!(time.minutes, minutes, "Invalid minutes for time {raw}");
    }

    #[test]
    fn valid_time() {
        for hour in 0..24 {
            for minute in 0..60 {
                parse_and_test_time(hour, minute);
            }
        }
        parse_and_test_time(24, 00);
    }

    #[test]
    fn time_display_impl() {
        for hour in 0..24 {
            for minute in 0..60 {
                let expected = format!(
                    "{hour} {}, {minute} {}",
                    if hour == 1 { "hour" } else { "hours" },
                    if minute == 1 { "minute" } else { "minutes" }
                );

                let raw = format!("{hour:0>2}:{minute:0>2}");
                let time: Time = raw.parse().expect(&format!("Failed to parse {raw}"));
                assert_eq!(time.to_string(), expected);
            }
        }
    }

    #[test]
    fn error_display_impl() {
        assert_eq!(TimeParseError::MissingColon.to_string(), "missing ':'");
        assert_eq!(TimeParseError::InvalidLength.to_string(), "invalid length");
        assert_eq!(TimeParseError::InvalidNumber.to_string(), "invalid number");
    }

    fn parse_and_assert_error(s: &str, expected_err: TimeParseError) {
        let actual_err = match s.parse::<Time>() {
            Ok(time) => panic!("Expected error: {expected_err}, but got {time} for {s}"),
            Err(err) => err,
        };

        assert_eq!(actual_err.to_string(), expected_err.to_string());
    }

    #[test]
    fn missing_colon() {
        let err = "2331".parse::<Time>().expect_err("Parsing 2331 should fail");
        match err {
            TimeParseError::MissingColon | TimeParseError::InvalidLength => {}
            TimeParseError::InvalidNumber => {
                panic!("Expected either invalid colon or length when parsing 2331")
            }
        }

        parse_and_assert_error("23.31", TimeParseError::MissingColon);
        parse_and_assert_error("23331", TimeParseError::MissingColon);
        parse_and_assert_error("23?31", TimeParseError::MissingColon);
        parse_and_assert_error("23;31", TimeParseError::MissingColon);
    }

    #[test]
    fn invalid_length() {
        parse_and_assert_error("23:231", TimeParseError::InvalidLength);
        parse_and_assert_error("123:23", TimeParseError::InvalidLength);

        let err = "".parse::<Time>().expect_err("Parsing an empty string should fail");
        match err {
            TimeParseError::MissingColon | TimeParseError::InvalidLength => {}
            TimeParseError::InvalidNumber => {
                panic!("Expected either invalid colon or length when parsing empty string")
            }
        }

        let err = "2".parse::<Time>().expect_err("Parsing 2 should fail");
        match err {
            TimeParseError::MissingColon | TimeParseError::InvalidLength => {}
            TimeParseError::InvalidNumber => {
                panic!("Expected either invalid colon or length when parsing 2")
            }
        }

        let err = "23".parse::<Time>().expect_err("Parsing 23 should fail");
        match err {
            TimeParseError::MissingColon | TimeParseError::InvalidLength => {}
            TimeParseError::InvalidNumber => {
                panic!("Expected either invalid colon or length when parsing 23")
            }
        }

        parse_and_assert_error("23:2", TimeParseError::InvalidLength);
        parse_and_assert_error("23:223", TimeParseError::InvalidLength);
        parse_and_assert_error(
            "1823812831283819:812381289123",
            TimeParseError::InvalidLength,
        );
    }

    #[test]
    fn invalid_time() {
        parse_and_assert_error("25:12", TimeParseError::InvalidNumber);
        parse_and_assert_error("12:60", TimeParseError::InvalidNumber);
        parse_and_assert_error("24:01", TimeParseError::InvalidNumber);
    }

    #[test]
    fn invalid_number() {
        parse_and_assert_error("12:2a", TimeParseError::InvalidNumber);
        parse_and_assert_error("1z:23", TimeParseError::InvalidNumber);
        parse_and_assert_error("ft:23", TimeParseError::InvalidNumber);
        parse_and_assert_error("23:sd", TimeParseError::InvalidNumber);
        parse_and_assert_error("12:0/", TimeParseError::InvalidNumber);
        parse_and_assert_error("/2:00", TimeParseError::InvalidNumber);
        parse_and_assert_error("12:-1", TimeParseError::InvalidNumber);
        parse_and_assert_error("-1:12", TimeParseError::InvalidNumber);
    }
}
`

func ex04Test(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}

	if err := testutils.AppendStringToFile(Ex04TestMod, exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}

	_, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"test", "--release", "shortinette_rust_test_module03_ex04_0001"}, testutils.WithTimeout(5*time.Second))
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex04() Exercise.Exercise {
	return Exercise.NewExercise("04", "ex04", []string{"src/main.rs", "Cargo.toml"}, 10, ex04Test)
}
