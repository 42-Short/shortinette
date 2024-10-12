# Module 01: Indirections

## Foreword

```rust
fn punch_card() {
    let rust = (
        ..=..=.. ..    .. .. .. ..    .. .. .. ..    .. ..=.. ..
        ..=.. ..=..    .. .. .. ..    .. .. .. ..    ..=..=..=..
        ..=.. ..=..    ..=.. ..=..    .. ..=..=..    .. ..=.. ..
        ..=..=.. ..    ..=.. ..=..    ..=.. .. ..    .. ..=.. ..
        ..=.. ..=..    ..=.. ..=..    .. ..=.. ..    .. ..=.. ..
        ..=.. ..=..    ..=.. ..=..    .. .. ..=..    .. ..=.. ..
        ..=.. ..=..    .. ..=..=..    ..=..=.. ..    .. ..=.. ..
    );
    println!("{rust:?}");
}
```

*Extracted from `rustc`'s [unit tests](https://github.com/rust-lang/rust/blob/131f0c6df6777800aa884963bdba0739299cd31f/tests/ui/weird-exprs.rs#L126-L134).*

## General Rules

* Any exercise you turn in must compile using the `cargo` package manager, either with `cargo run`
if the subject requires a _program_, or with `cargo test` otherwise. Only dependencies specified
in the allowed dependencies section are allowed. Only symbols specified in the `allowed symbols`
section are allowed.

* Every exercise must be part of a virtual Cargo workspace, a single `workspace.members` table must
be declared for the whole module.

* Everything must compile _without warnings_ with the `rustc` compiler available on the school's
machines without additional options.  You are _not_ allowed to use `unsafe` code anywere in your
code.

* You are generally not authorized to modify lint levels - either using `#[attributes]`,
`#![global_attributes]` or with command-line arguments. You may optionally allow the `dead_code`
lint to silence warnings about unused variables, functions, etc.

* For exercises managed with cargo, the command `cargo clippy -- -D warnings` must run with no errors!

* You are _strongly_ encouraged to write extensive tests for the functions and programs you turn in.
 Tests (when not specifically required by the subject) can use the symbols you want, even if
they are not specified in the `allowed symbols` section. **However**, tests should **not** introduce **any additional external dependencies** beyond those already required by the subject.

## Exercise 00: Reference me daddy

```txt
turn-in directory:
    ex00/

files to turn-in:
    src/lib.rs  Cargo.toml
```

Create two **functions**. Both must add two integers together.

```rust
fn add(a: &i32, b: i32) -> i32;
fn add_assign(a: &mut i32, b: i32);
```

* `add` must return the result of the operation.
* `add_assign` must store the result of the operation in `a`.

## Exercise 01: Point Of No Return (v2)

```txt
turn-in directory:
    ex01/

files to turn in:
    src/lib.rs  Cargo.toml
```

Write a **function** that returns the smallest value among two numbers.

```rust
fn min(a: &i32, b: &i32) -> &i32;
```

* Note that you may have to add some *lifetime annotations* to the function in order to make it
compile.
* The `return` keyword is still disallowed.

## Exercise 02: It's getting GOOD

```txt
turn-in directory:
    ex02/

files to turn in:
    src/lib.rs  Cargo.toml
```

Create a **function** that maps three color components to a name.

The name of a color is determined using the following rules, applied in order. The first rule that
`match`es the input color must be selected.

* The color `[0, 0, 0]` is "pure black".
* The color `[255, 255, 255]` is "pure white".
* The color `[255, 0, 0]` is "pure red".
* The color `[0, 255, 0]` is "pure green".
* The color `[0, 0, 255]` is "pure blue".
* The color `[128, 128, 128]` is "perfect grey".
* Any color whose components are all bellow 31 is "almost black".
* Any color whose red component is above 128, whose green and blue components are between 0 and 127 (inclusive),
is "redish".
* Any color whose green component is above 128, whose red and blue components are between 0 and 127 (inclusive),
is "greenish".
* Any color whose blue component is above 128, whose red and green components are between 0 and 127 (inclusive),
is "blueish".
* Any other color is named "unknown".

The `if` keyword is **_not_** allowed!

```rust
const fn color_name(color: &[u8; 3]) -> &str;
```

You might need to add *lifetime* annotations to the function to make it compile. Specifically, the
following test must compile and run:

```rust
#[cfg(test)]
#[test]
fn test_lifetimes() {
    let name_of_the_best_color;

    {
        let the_best_color = [42, 42, 42];
        name_of_the_best_color = color_name(&the_best_color);
    }

    assert_eq!(name_of_the_best_color, "unknown");
}
```

## Exercise 03: This module is fun!

```txt
turn-in directory:
    ex03/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    <[u32]>::{len, is_empty, contains}
    std::iter*
```

Write a **function** that returns the largest subslice of `haystack` that contains *all* numbers in
`needle`.

```rust
fn largest_group(haystack: &[u32], needle: &[u32]) -> &[u32];
```

* When multiple groups match the `needle`, the largest one is returned.
* When multiple largest groups are found, the first one is returned.

Example:

```rust
assert_eq!(largest_group(&[1, 3, 4, 3, 5, 5, 4], &[5, 3]), &[3, 5, 5]);
assert_eq!(largest_group(&[1, 3, 4, 3, 5, 5, 4], &[5]), &[5, 5]);
assert_eq!(largest_group(&[1, 3, 4, 3, 5, 5, 4], &[]), &[]);
assert_eq!(largest_group(&[1, 3, 4, 3, 5, 5, 4], &[4, 1]), &[]);
```

Once again, you may need to specify some *lifetime annotations* for the function. To check whether
your annotations are correct for that case, you can use this pre-defined `test_lifetimes` test.
It must compile and run.

```rust
#[test]
#[cfg(test)]
fn test_lifetimes() {
    let haystack = [1, 2, 3, 2, 1];
    let result;

    {
        let needle = [2, 3];
        result = largest_group(&haystack, &needle);
    }

    assert_eq!(result, &[2, 3, 2]);
}
```

## Exercise 04: Wait no...

```txt
turn-in directory:
    ex04/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    <[i32]>::{len, is_empty, swap}  std::{assert, assert_eq, panic}
    std::iter*
```

You are given a list of boxes (`[width, height]`). Sort that list of boxes in a way for every box
to be *contained* in the previous one. If the operation is not possible, the function must panic.

You are **not** allowed to flip the boxes to make them fit.

```rust
fn sort_boxes(boxes: &mut [[u32; 2]]);
```

Example:

```rust
let mut boxes = [[3, 3], [4, 3], [1, 0], [5, 7], [3, 3]];
sort_boxes(&mut boxes);
assert_eq!(boxes, [[5, 7], [4, 3], [3, 3], [3, 3], [1, 0]]);
```

## Exercise 05: One DSA a day keeps the doctor away

```txt
turn-in directory:
    ex05/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    std::vec::Vec::{remove, len, is_empty}
    std::iter*
```

Write a **function** that removes all repeated elements of a list, preserving its initial ordering.

```rust
fn deduplicate(list: &mut Vec<i32>);
```

Example:

```rust
let mut v = vec![1, 2, 2, 3, 2, 4, 3];
deduplicate(&mut v);
assert_eq!(v, [1, 2, 3, 4]);
```

## Exercise 06: Do you _really_ want this job?

```txt
turn-in directory:
    ex06/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    <[i32]>::{is_empty, len}
    std::vec::Vec::{push, len, is_empty, new, reverse}
    u8::is_ascii_digit
    std::assert
    std::iter*
```

Write a **function** that adds two numbers together. The numbers are given as a list of decimal
digits and may be arbitrarily large.

```rust
fn big_add(a: &[u8], b: &[u8]) -> Vec<u8>;
```

* `a` and `b` must only contain digits (`b'0'` to `b'9'` included). If anything else is found, the
function must panic.
* If either `a` or `b` is empty, the function panics.
* Input numbers may contain leading zeros, but the result must not have any.

Example:

```rust
assert_eq!(big_add(b"2", b"4"), b"6");
assert_eq!(big_add(b"0010", b"0200"), b"210");
```

## Exercise 07: Can I get into Google now?

```txt
turn-in directory:
    ex07/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    std::Vec::*
    std::iter*
```
Leonardo has `n` tasks, which he needs to prioritize. He organized them into a vector of tasks. One task is defined as follows:

```rust
struct Task{
    start_time: u32,
    end_time: u32,
    cookies: u32,
}
```

For a task `i`:
* `tasks[i].start_time` is the start time for `task[i]`
* `tasks[i].end_time` is the end time for `task[i]`
* `tasks[i].cookies` is how many cookies he will get from students for finishing `task[i]` 

Unfortunately, he sucks at multitasking. Write a **function** which returns the maximum amount of cookies he can get without any tasks overlapping. 
Your function must have this signature:

```rust
fn time_manager(tasks: &mut Vec<Task>) -> u32
```

**Constraints**

_note_: If Leonardo chooses a task ending at time `t`, he will be able to start another task that starts at time `t` right away.

You do not need to perform any input checks. You _may_ assume the following: 
* `task[i].start_time < task[i].end_time`
* `task[i].start_time >= 0`
* `task[i].end_time >= 1`

What you _may not_ assume is our tester not having a timeout ಠ_ಠ, so **_don't be a brute_**.

```
MIT License

Copyright (c) 2024 Nils Mathieu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
