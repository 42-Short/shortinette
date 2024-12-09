use std::path;

use crate::{repository_path, testable::Testable};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise02;

impl Testable for Exercise02 {
    fn path(&self) -> path::PathBuf {
        repository_path().join("ex02")
    }
}
