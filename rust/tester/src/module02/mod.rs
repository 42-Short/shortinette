use exercise00::Exercise00;
use exercise01::Exercise01;
use exercise02::Exercise02;
use exercise03::Exercise03;
use exercise04::Exercise04;
use exercise05::Exercise05;
use exercise06::Exercise06;
use exercise07::Exercise07;

use crate::{result::TestResult, testable::Testable};

mod exercise00;
mod exercise01;
mod exercise02;
mod exercise03;
mod exercise04;
mod exercise05;
mod exercise06;
mod exercise07;

#[derive(Debug, PartialEq, Eq)]
pub enum Module02 {
    Ex00(Exercise00),
    Ex01(Exercise01),
    Ex02(Exercise02),
    Ex03(Exercise03),
    Ex04(Exercise04),
    Ex05(Exercise05),
    Ex06(Exercise06),
    Ex07(Exercise07),
}

impl Module02 {
    pub fn new(exercise_id: u32) -> anyhow::Result<Self> {
        let module = match exercise_id {
            0 => Self::Ex00(Exercise00),
            1 => Self::Ex01(Exercise01),
            2 => Self::Ex02(Exercise02),
            3 => Self::Ex03(Exercise03),
            4 => Self::Ex04(Exercise04),
            5 => Self::Ex05(Exercise05),
            6 => Self::Ex06(Exercise06),
            7 => Self::Ex07(Exercise07),
            _ => anyhow::bail!("Invalid exercise id for module00"),
        };

        Ok(module)
    }
}

impl Testable for Module02 {
    fn run_test(&self) -> TestResult {
        match self {
            Self::Ex00(exercise) => exercise.run_test(),
            Self::Ex01(exercise) => exercise.run_test(),
            Self::Ex02(exercise) => exercise.run_test(),
            Self::Ex03(exercise) => exercise.run_test(),
            Self::Ex04(exercise) => exercise.run_test(),
            Self::Ex05(exercise) => exercise.run_test(),
            Self::Ex06(exercise) => exercise.run_test(),
            Self::Ex07(exercise) => exercise.run_test(),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn exercise00() {
        let module = Module02::new(0).expect("Parsing module02 with exercise id 0 failed");

        assert_eq!(module, Module02::Ex00(Exercise00));
    }

    #[test]
    fn exercise01() {
        let module = Module02::new(1).expect("Parsing module02 with exercise id 1 failed");

        assert_eq!(module, Module02::Ex01(Exercise01));
    }

    #[test]
    fn exercise02() {
        let module = Module02::new(2).expect("Parsing module02 with exercise id 2 failed");

        assert_eq!(module, Module02::Ex02(Exercise02));
    }

    #[test]
    fn exercise03() {
        let module = Module02::new(3).expect("Parsing module02 with exercise id 3 failed");

        assert_eq!(module, Module02::Ex03(Exercise03));
    }

    #[test]
    fn exercise04() {
        let module = Module02::new(4).expect("Parsing module02 with exercise id 4 failed");

        assert_eq!(module, Module02::Ex04(Exercise04));
    }

    #[test]
    fn exercise05() {
        let module = Module02::new(5).expect("Parsing module02 with exercise id 5 failed");

        assert_eq!(module, Module02::Ex05(Exercise05));
    }

    #[test]
    fn exercise06() {
        let module = Module02::new(6).expect("Parsing module02 with exercise id 6 failed");

        assert_eq!(module, Module02::Ex06(Exercise06));
    }

    #[test]
    fn exercise07() {
        let module = Module02::new(7).expect("Parsing module02 with exercise id 7 failed");

        assert_eq!(module, Module02::Ex07(Exercise07));
    }
}
