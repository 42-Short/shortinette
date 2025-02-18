use fs_extra::file::CopyOptions;
use std::fs::File;
use std::io::Write;
use std::path;

use crate::{cargo::Cargo, repository_path, result::TestResult, testable::Testable};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise01;

impl Testable for Exercise01 {
    fn path(&self) -> path::PathBuf {
        repository_path().join("ex01")
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

    fn prepare_cargo_tests(&self) -> Cargo {
        let cargo = Cargo::new("shortinette-test-module00-ex01", true);

        cargo
            .add_dependency("rand", "0.8")
            .expect("Failed to add rand as dependency to test project");

        let mut lib_file = File::create(cargo.path().join("src/lib.rs"))
            .expect("Failed to open src/lib.rs of test module");
        lib_file
            .write_all(self.cargo_test_mod().as_bytes())
            .expect("Failed to write test module into src/lib.rs of test module");

        let files = ["min.rs"];

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
