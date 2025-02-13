use crate::TestResult;
use chrono::{Datelike, NaiveDate, Weekday};
use similar_asserts::SimpleDiff;
use std::path;
use std::process::Command;

use crate::{repository_path, testable::Testable};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise05;

fn expected_output(year: u32) -> String {
    (1..=year)
        .flat_map(move |year| {
            (1..=12).filter_map(move |month| {
                let date = NaiveDate::from_ymd_opt(year.try_into().unwrap(), month, 13).unwrap();
                if date.weekday() != Weekday::Fri {
                    None
                } else {
                    Some(format!("Friday, {} 13, {}\n", date.format("%B"), year))
                }
            })
        })
        .collect()
}

impl Testable for Exercise05 {
    fn path(&self) -> path::PathBuf {
        repository_path().join("ex05")
    }

    fn cargo_test_mod(&self) -> &'static str {
        include_str!("./shortinette_tests.rs")
    }

    fn run_test(&self) -> TestResult {
        if !self.check_clippy() {
            eprintln!("`cargo clippy -- -D warnings` failed");

            return TestResult::CompilationError;
        }

        if self.compile().is_err() {
            eprintln!("Failed to compile");

            return TestResult::CompilationError;
        }

        if let Err(test_output) = self.run_cargo_tests() {
            eprintln!("{test_output}");

            return TestResult::Failed;
        }

        let command = Command::new("cargo")
            .current_dir(self.path())
            .arg("run")
            .arg("--release")
            .arg("--quiet")
            .output()
            .expect("Unable to execute cargo run");

        if !command.status.success() {
            eprintln!(
                "Error when executing cargo run: {}",
                String::from_utf8_lossy(&command.stderr)
            );
        }

        if !command.stderr.is_empty() {
            eprintln!(
                "Unexpected output on stderr: {}",
                String::from_utf8_lossy(&command.stderr)
            );
        }

        let expected = expected_output(2025);
        let actual = String::from_utf8_lossy(&command.stdout);

        if expected != actual {
            let diff = SimpleDiff::from_str(&actual, &expected, "got", "expected");
            eprintln!("Incorrect output from cargo run\n{}", diff);
            return TestResult::Failed;
        }

        TestResult::Passed
    }
}
