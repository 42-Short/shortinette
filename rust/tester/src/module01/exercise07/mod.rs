use crate::{testable::Testable, result::TestResult};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise07;

impl Testable for Exercise07 {
    fn run_test(&self) -> TestResult {
        todo!()
    }
}
