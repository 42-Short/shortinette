package R02

import (
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
)

var cargoTestModAsString07 = `

#[cfg(test)]
mod shortinette_tests_rust_0207 {
    use super::*;

    #[test]
    fn test_cell_is_alive() {
        assert_eq!(Cell::Alive.is_alive(), true);
        assert_eq!(Cell::Dead.is_alive(), false);
    }

    #[test]
    fn test_cell_is_dead() {
        assert_eq!(Cell::Alive.is_dead(), false);
        assert_eq!(Cell::Dead.is_dead(), true);
    }

    fn count_alive_cells(board: &Board) -> usize {
        board
            .cells
            .iter()
            .filter(|&cell| *cell == Cell::Alive)
            .count()
    }

    fn assert_board(
        board: &Board,
        expected_width: usize,
        expected_height: usize,
        expected_percantage_alive: f64,
    ) {
        let actual_cells_alive =
            (count_alive_cells(&board) as f64 / board.cells.len() as f64) * 100.0;
        let margin_of_error = 10.0;

        assert_eq!(board.width, expected_width);
        assert_eq!(board.height, expected_height);
        assert_eq!(board.cells.len(), expected_width * expected_height);
        assert!(
            (expected_percantage_alive - margin_of_error) <= actual_cells_alive
                && actual_cells_alive <= (expected_percantage_alive + margin_of_error),
            "Expected around {}% alive cells, got {}%",
            expected_percantage_alive,
            actual_cells_alive
        );
    }

    #[test]
    fn test_board_new() {
        let board = Board::new(10, 20, 30);
        assert_board(&board, 10, 20, 30.0);
    }

    fn set_multiple_cells(board: &mut Board, cells_to_set: &[(usize, usize, Cell)]) {
        for &(x, y, ref new_cell) in cells_to_set {
            board.cells[y * board.width + x] = new_cell.clone();
        }
    }

    fn get_cell(board: &Board, x: usize, y: usize) -> Cell {
        board.cells[y * board.width + x]
    }

    #[test]
    fn test_board_step_underpopulation() {
        let mut board = Board::new(10, 10, 0);
        let alive_cells = [(1, 1, Cell::Alive), (0, 0, Cell::Alive)];

        set_multiple_cells(&mut board, &alive_cells);
        board.step();
        assert!(
            get_cell(&board, 1, 1) == Cell::Dead,
            "Cell should die due to underpopulation"
        );
        assert!(
            get_cell(&board, 0, 0) == Cell::Dead,
            "Cell should die due to underpopulation"
        );
    }

    #[test]
    fn test_board_step_survival() {
        let mut board = Board::new(10, 10, 0);
        let alive_cells = [
            (1, 1, Cell::Alive),
            (0, 1, Cell::Alive),
            (0, 0, Cell::Alive),
        ];

        set_multiple_cells(&mut board, &alive_cells);
        board.step();
        for &(x, y, _) in &alive_cells {
            assert!(
                get_cell(&board, x, y) == Cell::Alive,
                "Cell at ({}, {}) should be Alive",
                x,
                y
            );
        }
    }

    #[test]
    fn test_board_step_overpopulation() {
        let mut board = Board::new(10, 10, 0);
        let alive_cells = [
            (1, 1, Cell::Alive),
            (0, 1, Cell::Alive),
            (0, 0, Cell::Alive),
            (1, 0, Cell::Alive),
            (2, 1, Cell::Alive),
        ];
        set_multiple_cells(&mut board, &alive_cells);
        board.step();
        assert!(
            get_cell(&board, 1, 1) == Cell::Dead,
            "Cell at (1, 1) should be Dead due to overpopulation"
        );
    }

    #[test]
    fn test_board_step_reproduction() {
        let mut board = Board::new(10, 10, 0);
        let cells = [
            (1, 1, Cell::Dead),
            (0, 1, Cell::Alive),
            (0, 0, Cell::Alive),
            (1, 0, Cell::Alive),
        ];
        set_multiple_cells(&mut board, &cells);
        board.step();
        assert!(
            get_cell(&board, 1, 1) == Cell::Alive,
            "Cell at (1, 1) should be Alive due to reproduction"
        );
    }
}
	
`

var clippyTomlAsString07 = ``

func ex07Test(exercise *Exercise.Exercise) Exercise.Result {
    return runDefaultTest(exercise, cargoTestModAsString07, clippyTomlAsString07, map[string]int{"unsafe": 0})
}

func ex07() Exercise.Exercise {
	return Exercise.NewExercise("07", "ex07", []string{"src/main.rs", "Cargo.toml"}, 20, ex07Test)
}