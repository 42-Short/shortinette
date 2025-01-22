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
* You **must not** have a `main` present if not specifically requested.

* Any exercise managed by cargo you turn in must compile _without warnings_ using the `cargo test` command. If not managed by cargo, it must compile _without warnings_ with the `rustc` compiler available on the school's
machines without additional options.

* Only dependencies specified in the allowed dependencies section are allowed.


* You are _not_ allowed to use the `unsafe` keyword anywere in your code.

* If not specified otherwise by the task description, you are generally not authorized to modify lint levels - either using `#[attributes]`,
`#![global_attributes]` or with command-line arguments. You may optionally allow the `dead_code`
lint to silence warnings about unused variables, functions, etc.

```rust
// Either globally:
#![allow(dead_code)] 

// Or locally, for a simple item:
#[allow(dead_code)]
fn my_unused_function() {}
```

* For exercises managed with cargo, the command `cargo clippy -- -D warnings` must run with no errors!

* You are _strongly_ encouraged to write extensive tests for the functions and programs you turn in. Tests can use the symbols you want, even if
they are not specified in the `allowed symbols` section. **However**, tests should not introduce **any additional external dependencies** beyond
those already required by the subject.

* When a type is in the allowed symbols, it is **implied** that its methods and attributes are also allowed to be used, including the attributes of its implemented traits.

* You are **always** allowed to use `Option` and `Result` types (either `std::io::Result` or the plain `Result`, up to you and your use case).

* You are **always** allowed to use `std::eprintln` for error handling.

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

allowed symbols:
    none
```

Write a **function** that returns the smallest value among two numbers.

```rust
fn min(a: &i32, b: &i32) -> &i32;
```

* Note that you may have to add some *lifetime annotations* to the function in order to make it
compile.
* The `return` keyword is still disallowed.

## Exercise 02: Don't we all love darkmode?

```txt
turn-in directory:
    ex02/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    none
```

Create a **function** that maps three color components to a name.

```rust
const fn color_name(color: &[u8; 3]) -> &str;
```
The `if` keyword is **_not_** allowed!

You are not allowed to use the `_ => ...` syntax in the match statement.

The name of a color is determined using the following rules, applied in order. The first rule that
`match`es the input color must be selected.

`Legend: [red, green, blue]`

* **"dark gray"**: Any color whose red, green, and blue components are all between 0 and 128 (inclusive) is "Dark Gray/Black".

* **"dark red"**: Any color whose red component is between 128 and 255 (inclusive), and whose green and blue components are both between 0 and 128 (inclusive), is "Dark Red".

* **"dark green"**: Any color whose green component is between 128 and 255 (inclusive), and whose red and blue components are both between 0 and 128 (inclusive), is "Dark Green".

* **"olive"**: Any color whose red and green components are both between 128 and 255 (inclusive), and whose blue component is between 0 and 128 (inclusive), is "Dark Yellow/Olive".

* **"dark blue"**: Any color whose blue component is between 128 and 255 (inclusive), and whose red and green components are both between 0 and 128 (inclusive), is "Dark Blue".

* **"purple"**: Any color whose red and blue components are both between 128 and 255 (inclusive), and whose green component is between 0 and 128 (inclusive), is "Dark Magenta/Purple".

* **"teal"**: Any color whose green and blue components are both between 128 and 255 (inclusive), and whose red component is between 0 and 128 (inclusive), is "Dark Cyan/Teal".

* **"light gray"**: Any color whose red, green, and blue components are all between 128 and 255 (inclusive) is "Light Gray/White".


**You might need to add *lifetime* annotations to the function to make it compile. Specifically, the
following test must compile and run:**

```rust
#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_lifetimes() {
        let name_of_the_best_color;

        {
            let the_best_color = [42, 42, 42];
            name_of_the_best_color = color_name(&the_best_color);
        }

        assert_eq!(name_of_the_best_color, "dark grey");
    }
}
```

## Exercise 03: Where are my damn keys?!

```txt
turn-in directory:
    ex03/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    <[u32]>::{len, is_empty, contains}
    std::iter::*
```

Write a **function** that returns the first occurrence of `needle` in `haystack`.

```rust
fn largest_group(haystack: &[u32], needle: &[u32]) -> &[u32];
```

* Note that you will need to add **lifetime annotations** to the function signature to ensure correct borrowing. The borrow checker must enforce that the resulting slice is borrowed from `haystack`, not from `needle`.
* If `needle` is not in `haystack`, return an empty slice.

Example:

```rust
assert_eq!(largest_group(&[1, 3, 4, 3, 5, 5, 4], &[1, 3]), &[1, 3]);
assert_eq!(largest_group(&[1, 3, 4, 3, 5, 5, 4], &[5]), &[5]);
assert_eq!(largest_group(&[1, 3, 4, 3, 5, 5, 4], &[6, 9]), &[]);
assert_eq!(largest_group(&[1, 3, 4, 3, 5, 5, 4], &[4, 3]), &[4, 3]);
```
This test must compile and run:
```rust
#[test]
#[cfg(test)]
fn test_lifetimes() {
    let haystack = [1, 2, 3, 2, 1];
    let result;

    {
        let needle = [2, 3];
        // The result should be a valid slice of haystack after needle has expired
        result = largest_group(&haystack, &needle);
    }
    
    assert_eq!(result, &[2, 3]);
}
```

## Exercise 04: Wait no...

```txt
turn-in directory:
    ex04/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    <[i32]>::{len, is_empty, swap}
    std::{assert, assert_eq, panic}
    std::iter::*
```

Write a function that sorts a list of boxes in such a way that each box can be "contained" in the previous one without any box being flipped.

The function signature should look like this:
```rust
fn sort_boxes(boxes: &mut [[u32; 2]]);
```

The sorting should follow these criteria:
* **Definition of containment:** A box `[width, height]` can be contained inside another box `[prev_width, prev_height]` if:
    * `prev_width >= width`
    * `prev_height >= height`
* If sorting is not possible, the function should panic.

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

---
**License Notice:**
This file contains content licensed under two different terms:
- The MIT License applies to the original content (see `LICENSES/MIT-rust-subjects.txt`).
- The Creative Commons Attribution-ShareAlike 4.0 International (CC BY-SA 4.0 applies to any modifications or additions (see `LICENSES/CC-BY-SA-4.0.txt`).

When distributing modified versions, you must comply with both the MIT License and the CC BY-SA 4.0.
For complete details, refer to the main licensing file of Shortinette.
---
