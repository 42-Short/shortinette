use std::{env, path};

use module::Module;
use result::TestResult;

mod cargo;
mod module;
mod result;
mod testable;

mod module00;
mod module01;
mod module02;
mod module03;
mod module04;
mod module05;
mod module06;

fn main() -> TestResult {
    let module = {
        let module = env::var("MODULE").expect("MODULE env variable not set");
        let exercise = env::var("EXERCISE").expect("EXERCISE env variable not set");

        let module_id = module.parse::<u32>().expect("Failed to parse module id");

        let exercise_id = exercise
            .parse::<u32>()
            .expect("Failed to parse exercise id");

        Module::new(module_id, exercise_id).unwrap_or_else(|_| {
            panic!("Failed to create module ({module_id}) with exercise id {exercise_id}")
        })
    };

    module.run_test()
}

fn repository_path() -> path::PathBuf {
    let repository_path = env::var("REPOSITORY")
        .map(|path| path::PathBuf::from(&path))
        .unwrap_or_else(|_| env::current_dir().expect("Failed to access cwd"));

    // Assert is fine here since the tester should only be run on valid directories
    assert!(repository_path.exists(), "Repository path does not exist");

    repository_path
}
