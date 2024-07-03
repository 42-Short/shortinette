package R00

import exercisebuilder "github.com/42-Short/shortinette/internal/interfaces/exercise"

const YesMain = `
use std::thread;
use std::time::Duration;
use std::sync::mpsc;

fn main() {
	let (tx, rx) = mpsc::channel();

	let handle = thread::spawn(move || {
		tx.send(()).unwrap();
		yes();
	});

	rx.recv().unwrap();

	thread::sleep(Duration::from_secs(1));

	handle.thread().unpark();
}
`

func yes() bool {
	return true
}

func collatz() bool {
	return true
}

func print_bytes() bool {
	return true
}

func ex02Test(test *exercisebuilder.Test) bool {
	if yes() && collatz() && print_bytes() {
		return true
	}
	return false
}

func ex02() exercisebuilder.ExerciseBuilder {
	return exercisebuilder.NewExerciseBuilder().
		SetName("EX02").
		SetTurnInDirectory("ex02").
		SetTurnInFile("").
		SetExerciseType("").
		SetPrototype("").
		SetAllowedMacros(nil).
		SetAllowedFunctions(nil).
		SetAllowedKeywords(nil).
		SetExecuter(ex02Test)
}
