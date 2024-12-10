use std::{fs, io::Write, path};

use rand::seq::SliceRandom;

use crate::{cargo::Cargo, result::TestResult};

pub trait Testable {
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

        TestResult::Passed
    }

    fn run_cargo_tests(&self) -> Result<(), String> {
        // We do not want to create a new Cargo project for every test
        let cargo = self.prepare_cargo_tests();

        let mut cargo_test_list = cargo.test_list(&["shortinette_tests"]);
        assert!(!cargo_test_list.is_empty(), "No shortinette tests found");

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
        let cargo = Cargo::new("shortinette-test-module", true);
        cargo
            .add_local_dependency(
                self.path()
                    .to_str()
                    .expect("Path did not contain valid unicode"),
            )
            .expect("Failed to add exercise as dependency to test project");

        let mut lib_file = fs::File::create(cargo.path().join("src/lib.rs"))
            .expect("Failed to open src/lib.rs of test module");
        lib_file
            .write_all(self.cargo_test_mod().as_bytes())
            .expect("Failed to write test module into src/lib.rs of test module");

        cargo
    }

    fn cargo_test_mod(&self) -> &'static str {
        include_str!("./default-cargo-test.rs")
    }

    fn path(&self) -> path::PathBuf;

    fn ensure_path(&self) -> path::PathBuf {
        let path = self.path();

        // Assert is fine here, since the tester will only be run on existing folders
        assert!(path.exists(), "Folder does not exist");

        path
    }

    fn compile(&self) -> Result<Option<path::PathBuf>, TestResult> {
        Cargo::copy_from(&self.ensure_path())
            .compile()
            .map_err(|_| TestResult::CompilationError)
    }

    fn clippy_config(&self) -> &'static str {
        include_str!("./default-clippy-rules.toml")
    }

    fn check_clippy(&self) -> bool {
        let config = self.clippy_config();

        self.get_cargo().check_clippy(config)
    }

    fn get_cargo(&self) -> Cargo {
        Cargo::copy_from(&self.ensure_path())
    }
}
