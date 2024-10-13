package R01

import (
	"time"

	"github.com/42-Short/shortinette/rust/alloweditems"
	"github.com/42-Short/shortinette/rust/cargo"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var cargoTest06 = `
#[cfg(test)]
mod shortinette_tests_rust_0106 {
	use super::*;

	#[test]
	#[should_panic]
	fn test_both_empty() {
		big_add(b"", b"");
	}

	#[test]
	#[should_panic]
	fn test_a_empty() {
		big_add(b"", b"42");
	}

	#[test]
	#[should_panic]
	fn test_b_empty() {
		big_add(b"42", b"");
	}

	#[test]
	fn test_invalid_bytes() {
		let input_strings = vec![b"/42", b"42/", b":42", b"42:",  b" 42", b"42 ", b"+42", b"-42", &[52, 50, 0]];
		for a in &input_strings {
			for b in &input_strings {
				let result = std::panic::catch_unwind(|| {
					big_add(*a, *b);
				});
				match result {
					Ok(_) => panic!("Invalid input not correctly handled"),
					Err(payload) => {
						let payload_str = match payload.downcast_ref::<String>() {
							Some(s) => s.as_str(),
							None => match payload.downcast_ref::<&str>() {
								Some(s) => s,
								None => continue
							}
						};
						if payload_str.contains("subtract with overflow") {
							panic!("Invalid input not correctly handled");
						}
					}
				}
			}
		}
	}

	#[test]
	fn test_0() {
		assert_eq!(big_add(b"1", b"2"), b"3", r#"Failed for (b"1", b"2")"#);
	}

	#[test]
	fn test_1() {
		assert_eq!(big_add(b"8", b"4"), b"12", r#"Failed for (b"8", b"4")"#);
	}

	#[test]
	fn test_2() {
		assert_eq!(big_add(b"0", b"0"), b"0", r#"Failed for (b"0", b"0")"#);
	}

	#[test]
	fn test_3() {
		assert_eq!(big_add(b"00003674", b"1757"), b"5431", r#"Failed for (b"00003674", b"1757")"#);
	}

	#[test]
	fn test_4() {
		assert_eq!(big_add(b"3319", b"00001259"), b"4578", r#"Failed for (b"3319", b"00001259")"#);
	}

	#[test]
	fn test_5() {
		assert_eq!(big_add(b"42", b"57"), b"99", r#"Failed for (b"42", b"57")"#);
	}

	#[test]
	fn test_6() {
		assert_eq!(big_add(b"42", b"58"), b"100", r#"Failed for (b"42", b"58")"#);
	}

	#[test]
	fn test_7() {
		assert_eq!(big_add(b"42", b"59"), b"101", r#"Failed for (b"42", b"59")"#);
	}

	#[test]
	fn test_8() {
		assert_eq!(big_add(b"50000", b"50000"), b"100000", r#"Failed for (b"50000", b"50000")"#);
	}

	#[test]
	fn test_9() {
		assert_eq!(big_add(b"340282366920938463463374607431768211456", b"18446744073709551616"), b"340282366920938463481821351505477763072", r#"Failed for (b"340282366920938463463374607431768211456", b"18446744073709551616")"#);
	}

	#[test]
	fn test_10() {
		assert_eq!(big_add(b"000000", b"00000"), b"0", r#"Failed for (b"000000", b"00000")"#);
	}
}
`

func ex06Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	if err := testutils.AppendStringToFile(cargoTest06, exercise.TurnInFiles[0]); err != nil {
		return Exercise.InternalError(err.Error())
	}
	return cargo.CargoTest(exercise, 500*time.Millisecond, []string{})
}

func ex06() Exercise.Exercise {
	return Exercise.NewExercise("06", "ex06", []string{"src/lib.rs", "Cargo.toml"}, 15, ex06Test)
}
