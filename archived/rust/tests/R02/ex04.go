package R02

import (
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
)

var cargoTestModAsString04 = `

#[cfg(test)]
mod shortinette_tests_rust_0204 {
    use super::*;

    fn get_filled_todo_list(todos_ammount: u32) -> TodoList {
        let mut todo_list = TodoList::new();
        for i in 0..todos_ammount {
            todo_list.add(format!("dubidubidu{}", i));
        }
        todo_list
    }

    fn set_todo_list_dones(is_done: usize, todo_list: &mut TodoList) {
        for _ in 0..is_done {
            todo_list.done(0);
        }
    }

    #[test]
    fn test_todo_list_add() {
        assert_eq!(get_filled_todo_list(2).todos, vec!["dubidubidu0", "dubidubidu1"])
    }

    #[test]
    fn test_todo_list_done_valid_index() {
        let mut todo_list = get_filled_todo_list(3);
        set_todo_list_dones(2, &mut todo_list);
        assert_eq!(todo_list.todos, vec!["dubidubidu2"]);
        assert_eq!(todo_list.dones, vec!["dubidubidu0", "dubidubidu1"]);
    }

    #[test]
    fn test_todo_list_done_invalid_index() {
        let mut todo_list = get_filled_todo_list(1);
        todo_list.done(1);
        assert_eq!(todo_list.todos, vec!["dubidubidu0"]);
        assert!(todo_list.dones.is_empty());
    }

    #[test]
    fn test_todo_list_purge_with_done_items() {
        let mut todo_list = get_filled_todo_list(10);
        set_todo_list_dones(5, &mut todo_list);
        todo_list.purge();
        assert!(todo_list.dones.is_empty());
        assert_eq!(todo_list.todos.len(), 5);
    }

    #[test]
    fn test_todo_list_purge_without_done_itmes() {
        let mut todo_list = get_filled_todo_list(10);
        todo_list.purge();
        assert!(todo_list.dones.is_empty());
        assert_eq!(todo_list.todos.len(), 10);
    }
}

`

var clippyTomlAsString04 = ``

func ex04Test(exercise *Exercise.Exercise) Exercise.Result {
    return runDefaultTest(exercise, cargoTestModAsString04, clippyTomlAsString04, map[string]int{"unsafe": 0})
}

func ex04() Exercise.Exercise {
	return Exercise.NewExercise("04", "ex04", []string{"src/main.rs", "Cargo.toml"}, 10, ex04Test)
}