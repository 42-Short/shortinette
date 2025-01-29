use rand::seq::SliceRandom;
use serde::Deserialize;
use std::{
    fs,
    path::{self, PathBuf},
    process::Command,
};

use crate::{repository_path, result::TestResult, testable::Testable};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise04;

#[derive(Deserialize)]
struct PackageConfig {
    name: String,
    edition: Option<String>,
    description: Option<String>,
    authors: Option<Vec<String>>,
    publish: Option<bool>,
}

#[derive(Deserialize)]
struct CargoToml {
    package: PackageConfig,
}

fn read_cargo_toml(path: &PathBuf) -> Option<CargoToml> {
    let content = match fs::read_to_string(path) {
        Ok(str) => str,
        Err(err) => {
            eprintln!("Unable to read Cargo.toml: {}", err);
            return None;
        }
    };

    let cargotoml: CargoToml = match toml::from_str(&content) {
        Ok(toml) => toml,
        Err(e) => {
            eprintln!("Unable to deserialize Cargo.toml: {}", e);
            return None;
        }
    };

    Some(cargotoml)
}

fn check_cargo_toml(path: &PathBuf) -> bool {
    let cargotoml = match read_cargo_toml(path) {
        Some(toml) => toml,
        None => return false,
    };

    if cargotoml.package.name != "module00-ex04" {
        eprintln!("Incorrect module name: '{}'", cargotoml.package.name);
        return false;
    }

    if cargotoml
        .package
        .edition
        .is_none_or(|value| value != "2021")
    {
        eprintln!("Edition is not set to 2021");
        return false;
    }

    if cargotoml
        .package
        .authors
        .is_none_or(|authors| authors.len() != 1 || authors[0].is_empty())
    {
        eprintln!("Author check failed");
        return false;
    }

    if cargotoml
        .package
        .description
        .as_ref()
        .is_none_or(|description| {
            description
                != "my answer to the fifth exercise of the first module of 42's Rust Piscine"
        })
    {
        eprintln!(
            "Incorrect description: '{:?}'",
            cargotoml.package.description
        );
        return false;
    }

    if cargotoml.package.publish.is_none_or(|publish| publish) {
        eprintln!("Publish not set to false");
        return false;
    }

    true
}

struct TestOutput<'a> {
    exit_code: i32,
    profile: Option<&'a str>,
    expected_output: &'a str,
}

struct TestConfig<'a> {
    target: Option<&'a str>,
    executable_name: &'a str,
    test_output: &'a TestOutput<'a>,
}

impl<'a> TestConfig<'a> {
    fn new(target: Option<&'a str>, test_output: &'a TestOutput) -> Self {
        let executable_name = target.unwrap_or("module00-ex04");

        Self {
            target,
            executable_name,
            test_output,
        }
    }
}

fn nm_check(path: &PathBuf, testconfig: &TestConfig) -> bool {
    let profile_folder = testconfig.test_output.profile.unwrap_or("debug");
    let path = path
        .join("target")
        .join(profile_folder)
        .join(testconfig.executable_name);
    let release = testconfig
        .test_output
        .profile
        .is_some_and(|profile| profile == "release");

    let command = Command::new("nm").arg(path).output();

    if let Ok(output) = command {
        if !output.status.success() {
            eprintln!("nm didn't execute successfully");
            return false;
        } else if !output.stderr.is_empty()
            && !String::from_utf8_lossy(&output.stderr).contains("no symbols")
        {
            eprintln!(
                "Unexpected content on stderr: {}",
                String::from_utf8_lossy(&output.stderr)
            );
            return false;
        } else if release && !output.stdout.is_empty() {
            eprintln!("Debugging symbols haven't been stripped in Release mode!");
            return false;
        } else if !release && output.stdout.is_empty() {
            eprintln!("Debugging symbols should only be stripped in Release mode!");
            return false;
        }
    }
    true
}

fn run_executable(path: &PathBuf, testconfig: &TestConfig) -> bool {
    let mut binding = Command::new("cargo");
    let mut command = binding.current_dir(path).arg("run");

    if let Some(profile) = testconfig.test_output.profile {
        command = command.arg("--profile").arg(profile);
    }

    if let Some(target) = testconfig.target {
        command = command.arg("--bin").arg(target);
    }

    let output = command.output();
    if let Ok(output) = output {
        if let Some(exitcode) = output.status.code() {
            if exitcode != testconfig.test_output.exit_code {
                eprintln!(
                    "Incorrect exit code, expected {}, got {}",
                    testconfig.test_output.exit_code, exitcode
                );
                return false;
            }
        }

        let received_output = String::from_utf8_lossy(&output.stdout);
        if received_output != testconfig.test_output.expected_output {
            eprintln!(
                "Output differs\nExpected: '{}'\nGot: '{}'",
                testconfig.test_output.expected_output, received_output
            );
            return false;
        }

        if !nm_check(&path, testconfig) {
            return false;
        }
    } else {
        eprintln!("Unable to execute cargo run");
        return false;
    }
    true
}

impl Testable for Exercise04 {
    fn path(&self) -> path::PathBuf {
        repository_path().join("ex04")
    }

    fn run_test(&self) -> TestResult {
        let mut rng = rand::thread_rng();

        if !self.check_clippy() {
            eprintln!("`cargo clippy -- -D warnings` failed");

            return TestResult::CompilationError;
        }

        if self.compile().is_err() {
            eprintln!("Failed to compile");

            return TestResult::CompilationError;
        }

        if !check_cargo_toml(&self.ensure_path().join("Cargo.toml")) {
            return TestResult::Failed;
        }

        let mut tests = [
            TestConfig::new(
                None,
                &TestOutput {
                    exit_code: 0,
                    profile: None,
                    expected_output: "Hello, Cargo!\n",
                },
            ),
            TestConfig::new(
                None,
                &TestOutput {
                    exit_code: 0,
                    profile: Some("release"),
                    expected_output: "Hello, Cargo!\n",
                },
            ),
            TestConfig::new(
                Some("other"),
                &TestOutput {
                    exit_code: 0,
                    profile: None,
                    expected_output: "Hey! I'm the other bin target!\n",
                },
            ),
            TestConfig::new(
                Some("other"),
                &TestOutput {
                    exit_code: 0,
                    profile: Some("release"),
                    expected_output: "Hey! I'm the other bin target!\nI'm in release mode!\n",
                },
            ),
            TestConfig::new(
                Some("test-overflows"),
                &TestOutput {
                    exit_code: 101,
                    profile: None,
                    expected_output: "",
                },
            ),
            TestConfig::new(
                Some("test-overflows"),
                &TestOutput {
                    exit_code: 0,
                    profile: Some("no-overflows"),
                    expected_output: "255u8 + 1u8 == 0\n",
                },
            ),
        ];

        tests.shuffle(&mut rng);

        for test in &tests {
            if !run_executable(&self.ensure_path(), test) {
                return TestResult::Failed;
            }
        }

        TestResult::Passed
    }
}
