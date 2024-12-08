use std::path;

use crate::{repository_path, testable::Testable};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise03;

impl Testable for Exercise03 {
    fn path(&self) -> path::PathBuf {
        repository_path().join("ex03")
    }
}
