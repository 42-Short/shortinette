use crate::{testable::Testable, result::TestResult};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise01;

impl Testable for Exercise01 {
    fn run_test(&self) -> TestResult {
        todo!()
    }
}
