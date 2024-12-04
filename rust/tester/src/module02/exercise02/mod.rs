use crate::{testable::Testable, result::TestResult};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise02;

impl Testable for Exercise02 {
    fn run_test(&self) -> TestResult {
        todo!()
    }
}
