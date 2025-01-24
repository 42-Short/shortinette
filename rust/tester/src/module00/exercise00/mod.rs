use std::{
    path,
    process::{self},
};

use crate::{repository_path, result::TestResult, testable::Testable};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise00;

impl Testable for Exercise00 {
    fn path(&self) -> path::PathBuf {
        repository_path().join("ex00")
    }

    fn run_test(&self) -> TestResult {
        let path = match self.compile() {
            Ok(Some(path)) => path,
            _ => {
                eprintln!("Failed to compile");
                return TestResult::CompilationError;
            }
        };

        let output = match process::Command::new(&path).output() {
            Ok(output) => output,
            Err(_) => {
                eprintln!("Failed to execute ./hello");
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

        if output.stdout != "Hello, World!\n".as_bytes() {
            eprintln!(
                "Incorrect output on stdout: \nExpected: \"Hello, World!\n\"\nGot: \"{}\"",
                String::from_utf8_lossy(&output.stdout)
            );
            return TestResult::Failed;
        }

        TestResult::Passed
    }

    fn compile(&self) -> Result<Option<path::PathBuf>, TestResult> {
        self.ensure_path();

        let source_file = self.path().join("hello.rs");

        let output = process::Command::new("rustc")
            .current_dir(self.path())
            .arg(&source_file)
            .output()
            .expect("Failed to compile hello.rs");

        if !output.status.success() {
            return Err(TestResult::CompilationError);
        }

        let path = self.path().join("hello");

        Ok(Some(path))
    }
}
