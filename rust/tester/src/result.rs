#[derive(Debug, PartialEq, Eq)]
pub enum TestResult {
    Passed,
    Failed,
    CompilationError,
    Timeout,
    ForbiddenFunction,
}

impl TestResult {
    pub fn is_success(&self) -> bool {
        self == &Self::Passed
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn passed_is_success() {
        assert!(TestResult::Passed.is_success());
    }

    #[test]
    fn failed_is_not_success() {
        assert!(!TestResult::Failed.is_success());
    }

    #[test]
    fn compilation_is_not_success() {
        assert!(!TestResult::CompilationError.is_success());
    }

    #[test]
    fn timeout_is_not_success() {
        assert!(!TestResult::Timeout.is_success());
    }

    #[test]
    fn forbidden_function_is_not_success() {
        assert!(!TestResult::ForbiddenFunction.is_success());
    }
}
