use std::{
    path::{self, PathBuf},
    process::Command,
};

use crate::{repository_path, result::TestResult, testable::Testable};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise07;

impl Testable for Exercise07 {
    fn path(&self) -> path::PathBuf {
        repository_path().join("ex07")
    }

    fn compile(&self) -> Result<Option<path::PathBuf>, crate::result::TestResult> {
        self.ensure_path();

        let source_file = self.path().join("ft_putchar.rs");

        let _ = Command::new("rustc")
            .arg("-C")
            .arg("panic=abort")
            .arg("-C")
            .arg("link-args=-nostartfiles")
            .arg("-o")
            .arg("ft_putchar")
            .arg(&source_file)
            .output()
            .expect("Failed to compile ft_putchar.rs");

        Ok(Some(PathBuf::from("./ft_putchar")))
    }

    fn run_test(&self) -> crate::result::TestResult {
        let executable_path = self.compile().expect("Compilation failed.").unwrap();

        println!("{}", executable_path.to_str().unwrap());

        let output = Command::new(executable_path)
            .output()
            .expect("Failed to execute 'ft_putchar'");

        if output.status.code() != Some(42) {
            eprintln(
                "Incorrect exit code, expected: 42, got: {}",
                output.status.code().unwrap(),
            );
            return TestResult::Failed;
        }

        if output.stdout != b"42\n" {
            eprintln(
                "Incorrect output, expected: '42\\n', got: '{}'",
                output.stdout,
            );
            return TestResult::Failed;
        }

        TestResult::Passed
    }
}
