use std::collections::HashSet;
use std::io::{self, BufRead, BufReader, Read};
use std::path::{self, PathBuf};
use std::process::{Command, Stdio};

use crate::{repository_path, result::TestResult, testable::Testable};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise06;

fn validate_first_line<R: io::Read>(reader: &mut BufReader<R>) -> Result<(), String> {
    let mut output = String::new();

    reader
        .read_line(&mut output)
        .map_err(|e| format!("Unable to read line: {}", e.to_string()))?;

    if output != "Me and my infinite wisdom have found an appropriate secret you shall yearn for.\n"
    {
        return Err(format!("Incorrect content on first line: {}", output));
    }

    Ok(())
}

fn play_guessing_game<R: io::Read, W: io::Write>(
    reader: &mut BufReader<R>,
    writer: &mut W,
) -> Result<i64, String> {
    let mut output = String::new();
    let mut min: i64 = i32::MIN.into();
    let mut max: i64 = i32::MAX.into();

    validate_first_line(reader)?;

    loop {
        let current = (min + max) / 2;

        writer
            .write_all(format!("{}\n", current).as_bytes())
            .map_err(|err| format!("Unable to write into stdin: {}", err.to_string()))?;

        writer.flush().unwrap();

        reader
            .read_line(&mut output)
            .map_err(|err| format!("Unable to read line from stdout: {}", err.to_string()))?;

        let correct_str = format!("That is right! The secret was indeed the number {}, which you have brilliantly discovered!\n", current);

        if correct_str == output {
            return Ok(current);
        } else if min == max {
            return Err(format!(
                "Incorrect output\nExpected: {}Got: {}",
                correct_str, output
            ));
        }

        match output.as_str() {
            "This student might not be as smart as I was told. This answer is obviously too weak.\n" => min = current + 1,
            "Sometimes I wonder whether I should retire. I would have guessed higher.\n" => max = current - 1,
            _ => return Err(format!("Unexpected output: {}", output)),
        }

        output.clear();
    }
}

fn run_guessing_game(path: &PathBuf) -> Result<i64, String> {
    let mut child = Command::new("cargo")
        .current_dir(path)
        .arg("run")
        .arg("--quiet")
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()
        .expect("Unable to execute cargo run");

    let mut stdin = child
        .stdin
        .take()
        .expect("Unable to get stdin from command");
    let stdout = child
        .stdout
        .take()
        .expect("Unable to get stdout from command");
    let stderr = child
        .stderr
        .take()
        .expect("Unable to get stderr from command");

    let mut reader = BufReader::new(stdout);
    let mut err_reader = BufReader::new(stderr);
    let result = play_guessing_game(&mut reader, &mut stdin);

    child.kill().ok();
    child
        .wait()
        .map_err(|err| format!("Unable to wait for child process: {}", err.to_string()))?;

    let mut stderr_content = String::new();
    err_reader
        .read_to_string(&mut stderr_content)
        .map_err(|err| format!("Unable to read stderr: {}", err.to_string()))?;

    if !stderr_content.is_empty() {
        return Err(format!("Unexpected content on stderr: {}", stderr_content));
    }

    result
}

impl Testable for Exercise06 {
    fn path(&self) -> path::PathBuf {
        repository_path().join("ex06")
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

        let mut results = HashSet::new();
        for _ in 0..5 {
            match run_guessing_game(&self.path()) {
                Ok(number) => results.insert(number),
                Err(e) => {
                    eprintln!("{}", e);
                    return TestResult::Failed;
                }
            };
        }

        if results.len() == 1 {
            eprintln!("Number doesn't seem to be random");
            return TestResult::Failed;
        }

        TestResult::Passed
    }
}
