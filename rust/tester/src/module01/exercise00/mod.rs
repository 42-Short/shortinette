use std::path;

use crate::{repository_path, testable::Testable};

#[derive(Debug, PartialEq, Eq)]
pub struct Exercise00;

impl Testable for Exercise00 {
    fn path(&self) -> path::PathBuf {
        repository_path().join("ex00")
    }
}
