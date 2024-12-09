use std::path;

use crate::{repository_path, testable::Testable};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise07;

impl Testable for Exercise07 {
    fn path(&self) -> path::PathBuf {
        repository_path().join("ex07")
    }

    fn cargo_test_mod(&self) -> &'static str {
        include_str!("./shortinette_tests.rs")
    }
}
