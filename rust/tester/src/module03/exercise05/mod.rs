use crate::{testable::Testable, result::TestResult};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise05;

impl Testable for Exercise05 {
    fn run_test(&self) -> TestResult {
        todo!()
    }
}
