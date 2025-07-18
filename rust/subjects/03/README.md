# Module 03: Polymorphine

## Foreword

```rust
// Bastion of the Turbofish
// ------------------------
// Beware travellers, lest you venture into waters callous and unforgiving,
// where hope must be abandoned, ere it is cruelly torn from you. For here
// stands the bastion of the Turbofish: an impenetrable fortress holding
// unshaking against those who would dare suggest the supererogation of the
// Turbofish.
//
// Once I was young and foolish and had the impudence to imagine that I could
// shake free from the coils by which that creature had us tightly bound. I
// dared to suggest that there was a better way: a brighter future, in which
// Rustaceans both new and old could be rid of that vile beast. But alas! In
// my foolhardiness my ignorance was unveiled and my dreams were dashed
// unforgivingly against the rock of syntactic ambiguity.
//
// This humble program, small and insignificant though it might seem,
// demonstrates that to which we had previously cast a blind eye: an ambiguity
// in permitting generic arguments to be provided without the consent of the
// Great Turbofish. Should you be so naïve as to try to revolt against its
// mighty clutches, here shall its wrath be indomitably displayed. This
// program must pass for all eternity, fundamentally at odds with an impetuous
// rebellion against the Turbofish.
//
// My heart aches in sorrow, for I know I am defeated. Let this be a warning
// to all those who come after. Here stands the bastion of the Turbofish.

fn main() {
    let (oh, woe, is, me) = ("the", "Turbofish", "remains", "undefeated");
    let _: (bool, bool) = (oh<woe, is>(me));
}
```

*Extracted from `rustc`'s [unit tests](https://github.com/rust-lang/rust/blob/79d8a0fcefa5134db2a94739b1d18daa01fc6e9f/src/test/ui/bastion-of-the-turbofish.rs),
in memory of Anna Harren.*

## General Rules
* You **must not** have a `main` present if not specifically requested.

* Any exercise managed by cargo you turn in must compile _without warnings_ using the `cargo test` command. If not managed by cargo, it must compile _without warnings_ with the `rustc` compiler available on the school's
machines without additional options.

* Only dependencies specified in the allowed dependencies section are allowed.

* You are _not_ allowed to use the `unsafe` keyword anywhere in your code.

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

* All primitive types, i.e the ones you are able to use without importing them, are allowed.

* A type being allowed implies that its methods and attributes are allowed to be used as well, including the attributes of its implemented traits.

* You are **always** allowed to use `std::eprintln` for error handling.

* These rules may be overridden by specific exercises.

## Exercise 00: `choose`

```rust
// allowed symbols
use ftkit::random_number;

const allowed_dependencies = ["ftkit"];
const allowed_dependencies = [""];
const turn_in_directory = "ex00/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Create a **function** that randomly chooses a value among an input slice. If the provided list is
empty, return `None`

```rust
pub fn choose<T>(values: &[T]) -> Option<&T>;
```

## Exercise 01: Point Of No Return (v3)
```rust
// allowed symbols
use std::{
    cmp::PartialOrd,
    std::marker::Sized,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex01/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Again? Yes. Another `min` function! But I promise, this one's the last one.

* Create a public `min` function that takes *any* two values of a type that supports the `<` operator, and
returns the smaller one.

Example:

```rust
assert_eq!(min(12i32, -14i32), -14);
assert_eq!(min(12f32, 14f32), 12f32);
assert_eq!(min("abc", "def"), "abc");
assert_eq!(min(String::from("abc"), String::from("def")), "abc");
```

Still not allowed to use `return`!

## Exercise 02: 42

```rust
// allowed symbols
use std::{
    fmt::Debug,
    println,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex02/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Define the following trait:

```rust
pub trait FortyTwo {
    fn forty_two() -> Self;
}
```

* The `forty_two` associated function must return an instance of the implementer that represents the number 42 in some way.

Implement this trait for some common types, at least `u32` and `String`.

```rust
pub fn print_forty_two<T: Debug + FortyTwo>();
```

* The `print_forty_two` function must create an instance of `T` using the `FortyTwo` trait, and then print it to the standard output using its `Debug` implementation.

Create a `test` function that showcase this function being called for at least two distinct types.

## Exercise 03: Hello again Mr. Collatz
```rust
// allowed symbols
use std::{
    iter::{
        Iterator,
        traits::collect::FromIterator,
    }
};

const allowed_dependencies = [""];
const turn_in_directory = "ex03/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Define the following struct:
```rust
pub struct Collatz {
    value: u32,
}
```

Implement the `Iterator` and `FromIterator` trait, such that this test compiles and runs.

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_collatz() {
        let mut c = Collatz { value: 10 };
        let c_res = vec![5, 16, 8, 4, 2, 1];
        assert_eq!(c.count(), 6);

        c = Collatz { value: 10 };
        for (i, v) in c.enumerate() {
            assert_eq!(c_res[i], v);
        }

        c = Collatz { value: 10 };
        assert_eq!(c.collect::<Vec<u32>>(), c_res);
    }
}
```

## Exercise 04: What Time Is It?

```rust
// allowed symbols
use std::{
    fmt::{Display, Debug, Formatter},
    write,
    println,
    iter::*,
};

const allowed_dependencies = [""];
const allowed_dependencies = [""];
const turn_in_directory = "ex04/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Create a type named `Time` responsible for storing, well, a time.

```rust
pub struct Time {
    hours: u32,
    minutes: u32,
}

pub enum TimeParseError {
    MissingColon,
    InvalidLength,
    InvalidNumber,
}
```

Implement the required traits for `Time`, such that the following test compiles.

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_collatz() {
        let a = "12:20".parse::<Time>();
        let b = "15:14".parse::<Time>();

        assert!(a.is_ok());
        assert!(b.is_ok());

        let c = "12.20".parse::<Time>();
        let d = "12:2".parse::<Time>();
        let e = "12:2a".parse::<Time>();
        assert!(c.is_err());
        assert!(d.is_err());
        assert!(e.is_err());
    }
}
```

## Exercise 05: Quick Maths

```rust
// allowed symbols
use std::{
    fmt::Debug,
    ops::{
        Mul, MulAssign, Div, DivAssign, Add, AddAssign, Sub, SubAssign,
    },
    marker::Copy,
    clone::Clone,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex05/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

```rust
pub struct Vector<T> {
    pub x: T,
    pub y: Y,
}

impl<T> Vector<T> {
    pub fn new(x: T, y: T) -> Self;
}
```

* The `new` function must create a new `Vector<T>` with the specified components.
* Overload the `+`, `-`, `+=` and `-=` operators for `Vector<T>`, for any `T` that itself has support for those operators.
* Overload the `*`, `*=`, `/` and `/=` operators for `Vector<T>`, for any `T` that itself has support for those operators. The second operand of those operations *must not* be `Vector<T>`, but `T` itself, meaning that you must be able to compute `Vector::new(1, 2) * 3` but not `Vector::new(1, 2) * Vector::new(2, 3)`. You can require `T: Copy` when needed.
* Overload the `==` and `!=` operators for any `T` that supports them.
* Implement specifically for both `Vector<f32>` and `Vector<f64>` a `length` function that computes its length. The length of a vector can be computed using this formula: $|(x,y)| = \sqrt{(x^2 + y^2)}$.

The following tests must compile and run properly:

```rust
#[cfg(test)]
mod tests {
    #[test]
    fn test_a() {
        let v = Vector {
            x: String::from("Hello, World!"),
            y: String::from("Hello, Rust!"),
        };

        let w = v.clone();

        assert_eq!(&v, &w);
    }

    #[test]
    fn test_b() {
        let v = Vector::new("Hello, World!", "Hello, Rust!");
        let a = v;
        let b = v;
        assert_eq!(a, b);
    }
}
```

## Exercise 06: A Singly-Linked List

```rust
// allowed symbols
use std::{
    boxed::Box,
    panic,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex06/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

* Create a linked list type named `List<T>` defined as follows.

_Please note that the struct attributes are `pub` due to tester requirements, normally you would leave them private!_
```rust
pub struct Node<T> {
    pub value: T,
    pub next: Option<Box<Node<T>>>,
}

pub struct List<T> {
    pub head: Option<Box<Node<T>>>
}

impl<T> List<T> {
    pub fn new() -> Self;

    pub fn push_front(&mut self, value: T);
    pub fn push_back(&mut self, value: T);

    pub fn count(&self) -> usize;

    pub fn get(&self, i: usize) -> Option<&T>;
    pub fn get_mut(&mut self, i: usize) -> Option<&mut T>;

    pub fn remove_front(&mut self) -> Option<T>;
    pub fn remove_back(&mut self) -> Option<T>;
    pub fn clear(&mut self);
}
```

* `new` must create an empty list.
* `push_back` must append an element to the list.
* `push_front` must prepend an element to the list.
* `count` must return the number of elements present in the list.
* `get` must return a shared reference to the element at index `i`.
* `get_mut` must return an exclusive reference to the element at index `i`.
* `remove_back` must remove the last element of the list and return it.
* `remove_front` must remove the first element of the list and return it.
* `clear` must remove all elements of the list, leaving it empty.

The following tests must compile and pass.

```rust
#[cfg(test)]
mod tests {
    #[test]
    fn default_list_is_empty() {
        let list: List<i32> = Default::default();
        assert_eq!(list.count(), 0);
    }

    #[test]
    fn cloned_list_are_equal() {
        let mut list = List::new();
        list.push_back(String::from("Hello"));
        list.push_back(String::from("World"));

        let cloned = list.clone();
        assert_eq!(cloned.count(), list.count());
        assert_eq!(&cloned[0], &list[0]);
        assert_eq!(&cloned[1], &list[1]);
    }

    #[test]
    #[should_panic(expected = "tried to access out of bound index 10")]
    fn out_of_bound_access_panics() {
        let mut list: List<u32> = List::new();
        list.push_back(1);
        list.push_back(2);
        list.push_back(3);

        assert_eq!(list[10], 42);
    }
}
```

## Exercise 07: Comma-Separated Values

```rust
// allowed symbols
use std::{
    write,
    fmt::{Debug, Display, Formatter, Write},
    cmp::PartialEq,
    marker::Sized
};

const allowed_dependencies = [""];
const turn_in_directory = "ex07/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Let's create a generic CSV Encoder & Decoder. A CSV file is defined like this:

```txt
value1,value1,value1,value1
value2,value2,value2,value2
value3,value3,value3,value3
...
```

Each line corresponds to a *record*, and each column corresponds to a *field*.

* Create a `Field` trait, which describes how to encode or decode a value.

```rust
pub struct EncodingError;
pub struct DecodingError;

pub trait Field: Sized {
    fn encode(&self, target: &mut String) -> Result<(), EncodingError>;
    fn decode(field: &str) -> Result<Self, DecodingError>;
}
```

* Implement the `Field` trait for `String`. Keep in mind that finding a ',' or a '\n' in the string
is an `EncodingError`!
* Implement the `Field` trait for `Option<T>` as long as `T` implements `Field` too. The empty
string maps to `None`, while a non-empty string maps to the `Field` implementation of `T`.
* Implement the `Field` trait for *every possible integer type*. Because this is long, repetitive
and boring, write a *macro* to do it for you.

```rust
// ez
impl_field_for_int!(u8, u16, u32, u64, u128, usize, i8, i16, i32, i64, i128, isize);
```

* Create a `Record` trait, which describes how to encode or decode a collection of `Field`s.

```rust
pub trait Record: Sized {
    fn encode(&self, target: &mut String) -> Result<(), EncodingError>;
    fn decode(line: &str) -> Result<Self, DecodingError>; 
}
```

* Now, you have everything you need to create `decode_csv` and `encode_csv` functions.

```rust
pub fn encode_csv<R: Record>(records: &[R]) -> Result<String, EncodingError>;
pub fn decode_csv<R: Record>(contents: &str) -> Result<Vec<R>, DecodingError>;
```

* `encode_csv` takes a list of records and encode them into a `String`.
* `decode_csv` takes the content of a CSV file and decodes it into a list of records.

Example:

```rust
#[cfg(test)]
mod tests {
    #[derive(Debug, PartialEq)]
    struct User {
        name: String,
        age: u32,
    }

    impl Record for User { /* ... */ }

    #[test]
    fn test_encode() {
        let database = [
            User { name: "aaa".into(), age : 23 },
            User { name: "bb".into(), age: 2 },
        ];

        let csv = encode_csv(&database).unwrap();

        assert_eq!(
            csv,
            "\
            aaa,23\n\
            bb,2\n\
            "
        );
    }

    #[test]
    fn test_decode() {
        let csv = "\
            hello,2\n\
            yes,5\n\
            no,100\n\
        ";

        let database: Vec<User> = decode_csv(csv).unwrap();

        assert_eq!(
            database,
            [
                User { name: "hello".into(), age: 2 },
                User { name: "yes".into(), age: 5 },
                User { name: "no".into(), age: 100 },
            ]
        );
    }

    #[test]
    fn decoding_error() {
        let csv = "\
            hello,2\n\
            yes,6\n\
            no,23,hello\n\
        ";

        decode_csv::<User>(csv).unwrap_err();
    }
}

```

You might have noticed that implementing the `Record` trait is *very* repetitive. As a bonus (a bonus to the bonus, if you will), you can create an `impl_record!` macro to implement it in a single line:

```rust
struct MyType {
    id: u32,
    name: String,
}

// ez
impl_record!(MyType(id, name));

#[cfg(test)]
mod tests {
    #[test]
    fn test_impl_record() {
        let records = [
            MyType { id: 10, name: "Marvin".into() },
            MyType { id: 11, name: "Marvin".into() },
            MyType { id: 12, name: "Marvin".into() },
        ];

        let csv = encode_csv(&records).unwrap();
        assert_eq!(
            csv,
            "\
            10,Marvin\n\
            11,Marvin\n\
            12,Marvin\n\
            "
        );
    }
}
```

---
**License Notice:**
This file contains content licensed under two different terms:
- The MIT License applies to the original content (see `LICENSES/MIT-rust-subjects.txt`).
- The Creative Commons Attribution-ShareAlike 4.0 International (CC BY-SA 4.0 applies to any modifications or additions (see `LICENSES/CC-BY-SA-4.0.txt`).

When distributing modified versions, you must comply with both the MIT License and the CC BY-SA 4.0.
For complete details, refer to the main licensing file of Shortinette.
---
