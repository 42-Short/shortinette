use std::{env, process::ExitCode};

use module::Module;

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

fn main() -> ExitCode {
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

    let test_result = module.run_test();

    if test_result.is_success() {
        ExitCode::SUCCESS
    } else {
        ExitCode::FAILURE
    }
}
