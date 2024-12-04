use crate::{testable::Testable, result::TestResult};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise00;

impl Testable for Exercise00 {
    fn run_test(&self) -> TestResult {
        todo!()
    }
}
