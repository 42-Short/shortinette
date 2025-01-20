use std::path;

use crate::{repository_path, testable::Testable};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise05;

impl Testable for Exercise05 {
    fn path(&self) -> path::PathBuf {
        repository_path().join("ex05")
    }
}
