use std::path;

use crate::{repository_path, testable::Testable, TestResult};
use std::fmt;
use std::fmt::Write;
use std::process::Command;

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise03;

/*
* Intentionally using if conditions so nobody can just copy the function
* (once we actually have the keyword check)
*/
fn expected_output(number: usize) -> Result<String, fmt::Error> {
    let mut output = String::new();

    if number % 3 == 0 {
        output.push_str("fizz");
    }

    if number % 5 == 0 {
        output.push_str("buzz");
    }

    if output.is_empty() {
        match number % 11 {
            3 => output.push_str("FIZZ"),
            5 => output.push_str("BUZZ"),
            _ => write!(&mut output, "{}", number)?,
        }
    }

    Ok(output)
}

impl Testable for Exercise03 {
    fn path(&self) -> path::PathBuf {
        repository_path().join("ex03")
    }

    fn run_test(&self) -> TestResult {
        let path = match self.compile() {
            Ok(Some(path)) => path,
            _ => {
                eprintln!("Failed to compile");
                return TestResult::CompilationError;
            }
        };

        let output = match Command::new(&path).output() {
            Ok(output) => output,
            Err(_) => {
                eprintln!("Failed to execute ./fizzbuzz");
                return TestResult::CompilationError;
            }
        };

        if !output.status.success() {
            eprintln!("Exited with non-0 exit code");
            return TestResult::Failed;
        }

        if !output.stderr.is_empty() {
            eprintln!(
                "Unexpected output on stderr:\nExpected: \"\"\nGot: \"{}\"",
                String::from_utf8_lossy(&output.stderr)
            );
            return TestResult::Failed;
        }

        let output_string = String::from_utf8_lossy(&output.stdout);
        for (i, line) in output_string.lines().enumerate() {
            if i >= 100 {
                eprintln!("Too many lines in output!");
                return TestResult::Failed;
            }

            let expected = match expected_output(i + 1) {
                Ok(value) => value,
                Err(e) => {
                    eprintln!("Internal error: {}", e);
                    return TestResult::Failed;
                }
            };
            if expected != line {
                eprintln!(
                    "Output differs in line {}\nExpected: {}\nGot: {}",
                    i + 1,
                    expected,
                    line
                );
                return TestResult::Failed;
            }
        }

        TestResult::Passed
    }

    fn compile(&self) -> Result<Option<path::PathBuf>, TestResult> {
        self.ensure_path();

        let source_file = self.path().join("fizzbuzz.rs");

        let output = Command::new("rustc")
            .current_dir(self.path())
            .arg(&source_file)
            .output()
            .expect("Failed to compile fizzbuzz.rs");

        if !output.status.success() {
            return Err(TestResult::CompilationError);
        }

        let path = self.path().join("fizzbuzz");

        Ok(Some(path))
    }
}
