# Documentation for `tester` package
The `tester` package is a wrapper around all the testing and code analysis that is made on the submitted code.

## Package Structure
* **tester.go**: Containes the necessary functions for running all the tests.

## Setting up Tests
Currently, the `tester` package supports 2 kinds of testing:
1. Output assertion for programs which do _not_ take input.
**Configuration**: In your test configuration (`.yaml` file in the `testconfig` directory), specify the type as `program`. Then, under `tests`, you can add `assert_eq` and `assert_ne` tests, which will be run automatically.
2. Output assertion for functions.
**Configuration**: In `tests/R{project_number}/ex{exercise_number}/` add a `tests.rs` file. This will be appended to the student's file, and the tests will be run. Please have a look at the [Rust Doc](https://doc.rust-lang.org/book/ch11-01-writing-tests.html) for more details and how to set up tests.

Tests for programs with command-line input coming soon.
