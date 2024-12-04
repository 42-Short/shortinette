use crate::{testable::Testable, result::TestResult};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise04;

impl Testable for Exercise04 {
    fn run_test(&self) -> TestResult {
        todo!()
    }
}
