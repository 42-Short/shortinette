use std::path;

use crate::{repository_path, testable::Testable};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise04;

impl Testable for Exercise04 {
    fn path(&self) -> path::PathBuf {
        repository_path().join("ex04")
    }
}
