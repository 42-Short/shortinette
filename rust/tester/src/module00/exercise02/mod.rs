use crate::cargo::Cargo;
use crate::TestResult;
use fs_extra::file::CopyOptions;
use rand::prelude::SliceRandom;
use std::fs::File;
use std::io::Write;
use std::path;
use std::process::Command;

use crate::{repository_path, testable::Testable};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise02;

impl Testable for Exercise02 {
    fn path(&self) -> path::PathBuf {
        repository_path().join("ex02")
    }

    fn cargo_test_mod(&self) -> &'static str {
        include_str!("./shortinette_tests.rs")
    }

    fn run_test(&self) -> TestResult {
        if let Err(test_output) = self.run_cargo_tests() {
            eprintln!("{test_output}");

            return TestResult::Failed;
        }

        TestResult::Passed
    }

    fn run_cargo_tests(&self) -> Result<(), String> {
        let cargo = self.prepare_cargo_tests();

        let mut cargo_test_list = cargo.test_list(&["shortinette_tests"]);
        assert!(!cargo_test_list.is_empty(), "No shortinette tests found");

        let command = Command::new("cargo")
            .current_dir(cargo.path())
            .arg("build")
            .arg("--release")
            .output()
            .expect("Unable to run cargo build");

        assert!(
            command.status.success(),
            "Unable to compile program\n{}",
            String::from_utf8_lossy(&command.stderr)
        );

        let mut rng = rand::thread_rng();
        cargo_test_list.shuffle(&mut rng);

        let mut failed_output = Vec::new();
        for test in cargo_test_list {
            let Err(test_output) = cargo.run_test([test.as_str()]) else {
                continue;
            };

            failed_output.push(test_output);
        }

        if failed_output.is_empty() {
            Ok(())
        } else {
            Err(failed_output.join("\n"))
        }
    }

    fn prepare_cargo_tests(&self) -> Cargo {
        let cargo = Cargo::new("shortinette-test-module00-ex02", false);

        cargo
            .add_dependency("rand", "0.8")
            .expect("Failed to add rand as dependency to test project");

        cargo
            .add_dependency("base64", "0.22")
            .expect("Failed to add base64 as dependency to test project");

        cargo
            .add_dependency("wait-timeout", "0.2")
            .expect("Failed to add wait-timeout as dependency to test project");

        let mut main_file = File::create(cargo.path().join("src/main.rs"))
            .expect("Failed to open src/main.rs of test module");
        main_file
            .write_all(self.cargo_test_mod().as_bytes())
            .expect("Failed to write test module into src/main.rs of test module");

        let files = ["yes.rs", "collatz.rs", "print_bytes.rs"];

        let options = CopyOptions::new().overwrite(true);
        for file in files {
            fs_extra::file::copy(
                &self.ensure_path().join(file),
                cargo.path().join("src").join(file),
                &options,
            )
            .expect("Unable to copy file into cargo project");
        }

        cargo
    }
}
