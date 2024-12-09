use crate::result::TestResult;

pub trait Testable {
    fn run_test(&self) -> TestResult;
}
