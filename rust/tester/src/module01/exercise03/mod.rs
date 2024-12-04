use crate::{testable::Testable, result::TestResult};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise03;

impl Testable for Exercise03 {
    fn run_test(&self) -> TestResult {
        todo!()
    }
}
