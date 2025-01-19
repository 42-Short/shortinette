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

        let _ = Command::new("rustc")
            .arg("-C")
            .arg("panic=abort")
            .arg("link-args=-nostartfiles")
            .arg("-o")
            .arg("ft_putchar")
            .arg("ft_putchar.rs")
            .output()
            .expect("Failed to compile ft_putchar.rs.");

        Ok(Some(PathBuf::from("ft_putchar.rs")))
    }

    fn run_test(&self) -> crate::result::TestResult {
        let executable_path = self.compile().expect("Compilation failed.").unwrap();

        let output = Command::new(executable_path)
            .output()
            .expect("Failed to execute ft_putchar.");

        if output.status.code() != Some(42) {
            return TestResult::Failed;
        }

        if output.stdout != b"42" {
            return TestResult::Failed;
        }

        TestResult::Passed
    }
}
