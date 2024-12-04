use crate::{
    module00::Module00, module01::Module01, module02::Module02, module03::Module03,
    module04::Module04, module05::Module05, module06::Module06, result::TestResult,
    testable::Testable,
};

#[derive(Debug)]
pub enum Module {
    M00(Module00),
    M01(Module01),
    M02(Module02),
    M03(Module03),
    M04(Module04),
    M05(Module05),
    M06(Module06),
}

impl Module {
    pub fn new(module_id: u32, exercise_id: u32) -> anyhow::Result<Self> {
        let module = match module_id {
            0 => Self::M00(Module00::new(exercise_id)?),
            1 => Self::M01(Module01::new(exercise_id)?),
            2 => Self::M02(Module02::new(exercise_id)?),
            3 => Self::M03(Module03::new(exercise_id)?),
            4 => Self::M04(Module04::new(exercise_id)?),
            5 => Self::M05(Module05::new(exercise_id)?),
            6 => Self::M06(Module06::new(exercise_id)?),
            _ => anyhow::bail!("Failed to parse module id"),
        };

        Ok(module)
    }

    pub fn run_test(&self) -> TestResult {
        match self {
            Module::M00(module) => module.run_test(),
            Module::M01(module) => module.run_test(),
            Module::M02(module) => module.run_test(),
            Module::M03(module) => module.run_test(),
            Module::M04(module) => module.run_test(),
            Module::M05(module) => module.run_test(),
            Module::M06(module) => module.run_test(),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn all_valid_module_exercise_configurations() {
        for module_id in 0..=6 {
            for exercise_id in 0..=7 {
                let module = Module::new(module_id, exercise_id);

                assert!(module.is_ok());
            }
        }
    }

    #[test]
    fn too_big_module_id() {
        let module = Module::new(7, 0);

        assert!(module.is_err());
    }

    #[test]
    fn too_big_exercise_id() {
        let module = Module::new(0, 8);

        assert!(module.is_err());
    }
}
