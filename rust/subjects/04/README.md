# Module 04: Side Effects

## Foreword

(Intro)

My friend with the gift of gab? Ferris Crab.

(Verse 1)

One of my crates got a lot of fly traits
Twenty “i” eight edition? My decision: time to migrate
I’m getting irate at all the excess `unsafe`
wait — backtrace

We got a cute crab, which is the best crate?
That’s up for grabs. GitHub or Phab-
ricator, review my pull now or later
Hit @bors with the r+ and you’ll be my saviour

And when I’m coming through, I got a cargo too
Reaction to wasm? Domain working group
If you need a `regex`, BurntSushi is your dude
But if you need a `Future` well we also got a few

Popping off this Vec like a pimple
And you know that the block I’m from is an impl
So if I talk about an IR, no it’s not GIMPLE
Only `rustc` MIR, just that simple

(Chorus)

Thought there’d never be a Rust Rap?
Turns out this is just that
impl newsletter #RustFacts
Ferris Crab, that’s a must have
Data race, we gon’ bust that
Mem unsafe, we gon’ bust that
This the first and only Rust Rap
Ferris Crab, that’s a must have

(Verse 2)

If you never borrow check, then you’re gonna get wrecked
Pull out `gdb` cause you need to inspect out-of-bounds index
Oh guess what’s next?
Use after free turns out it’s gonna be

Or… just use the `rustc`
And you’ll be flushing all of these bugs down the drain
Gushing super fast code from your brain
No dusting: quite easy to maintain

What’s the secret sauce? It’s all zero cost
Couldn’t do it better if your boss
Demand you try to do it all by hand, but why?
Hate to be that guy, but generics monomorphize

Don’t use a `while` loop, `i < n`
Use an `Iterator`: much better by ten
And when you have a dozen eggs don’t start counting hens
But me and Ferris Crab: best friends to the end

(Chorus)

Thought there’d never be a Rust Rap?
Turns out this is just that
impl newsletter #RustFacts
Ferris Crab, that’s a must have
Data race, we gon’ bust that
Mem unsafe, we gon’ bust that
This the first and only Rust Rap
Ferris Crab, that’s a must have

(Outro)

My friend with the gift of gab? Ferris Crab.

*"[Ferris Crab](https://fitzgeraldnick.com/2018/12/13/rust-raps.html)"*

```rust
struct 🦀;
```

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

## Exercise 00: Wait that's it?

```rust
// no allowed symbols

const allowed_dependencies = [""];
const turn_in_directory = "ex00/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Create the `Outcome` and `Maybe` types, which should mimic `Result` and `Option` such that below test compiles and runs. For this exercise, you are exceptionally (_and obviously_) **not** allowed to use the `Option` and `Result` types.

```rust
#[cfg(test)]
mod tests{
    use super::*;

    fn outcome() -> Outcome<u32, &'static str> {
        Outcome::Good(42)
    }

    fn maybe() -> Maybe<u8> {
        Maybe::Definitely(42)
    }

    #[test]
    fn test() {
        let o = outcome();
        match o {
            Outcome::Good(n) => assert_eq!(n, 42),
            Outcome::Bad(_) => panic!("should be Good")
        }

        let m = maybe();
        match m {
            Maybe::Definitely(n) => assert_eq!(n, 42),
            Maybe::No => panic!("should be Definitely")
        }

    }
}
```

## Exercise 01: Tee-Hee

```rust
// allowed symbols
use std::{
    io::{Write, Read},
    fs::File,
    iter::*,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex01/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Write a **function** that copies `input` to `writer` and all filenames provided as arguments.

Your function must have the following signature:

```rust
pub fn tee<R: std::io::Read, W: std::io::Write>(input: &mut R, writer: &mut W, filenames: &[String]);
```

Example:

```rust
fn main() {
    let args: Vec<String> = std::env::args().skip(1).collect();
    tee(&mut std::io::stdin(), &mut std::io::stdout(), &args);
}
```

```plaintext
>_ echo "Hello, World!" | cargo run -- a b c
Hello, World!
>_ cat a b c
Hello, World!
Hello, World!
Hello, World!
```

Your function must not panic when interacting with the file system. All errors must be handled properly. You are free to choose what to do in that case, but you must *not* crash/panic.

## Exercise 02: Duh

```rust
// allowed symbols
use std::{
    io::Write,
    fs::{Metadata, metadata, read_dir, DirEntry, ReadDir},
    path::Path,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex02/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Create a **function** that computes the total size of a directory or file. Once computed, write the size followed by a newline to `writer`, formatted like in the examples below.

Your function must have the following signature:
```rust
pub fn duh<W: std::io::Write>(writer: &mut W, basedir: &str) -> Result<(), String>;
```

Example Usage:

```rust
fn main() {
    let arg: String = std::env::args().nth(1).unwrap();
    duh(&mut std::io::stdout(), &arg);
}

// Output:
// 1.2 gigabytes
```

 * If a size is less than a kilobyte, it is written in bytes. (e.g. `245 bytes`)
 * If a size is more than a kilobyte, it is written in kilobytes, with one decimal (e.g. `12.2 kilobytes`).
 * If a size is more than a megabyte, it is written in megabytes, with one decimal (e.g. `100.4 megabytes`).
 * If a size is more than a gigabyte, it is written in gigabytes, with one decimal (e.g. `23.9 gigabytes`).
 * For simplicity, you will assume that a kilobyte is `1000 bytes`, a megabyte is `1000 kilobytes`, etc.

Your function must never panic when interacting with the file system. Errors must be handled properly.

## Exercise 03: Pipeline

```rust
// allowed symbols
use std::{
    io::{Read, BufRead, Write},
    process::{Command, Stdio}.
    iter::*,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex03/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Create a **function** with the following signature:

```rust
pub fn pipeline<R: std::io::Read + BufRead, W: std::io::Write>(input: &mut R, writer: &mut W, args: &[String]) -> Result<(), String>;
```

It must spawn a process using the path (`args[0]`, arguments (`args[1..]`)), and standard input (`input`) passed as arguments, and write its output to `writer`.

Example Usage:

```rust
fn main() {
    pipeline(&mut std::io::stdin(), &mut std::io::stdout(), &[String::from("echo"), String::from("-n")]);
}
```
Expected Output:
```
$ echo "Hello, World!" | cargo run
Hello, World!%
```

The example invoked the `echo -n "Hello, World!"` command.

Your function must never panic when interacting with the system, you must handle errors properly.

## Exercise 04: Command Multiplexer

```rust
// allowed symbols
use std::{
    io::Write,
    process::{Command, Stdio, Child}.
    iter::*,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex04/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Create a **function** with the following signature:

```rust
pub fn multiplexer<W: std::io::Write>(writer: &mut W, command_lines: &[&[String]]) -> Result<(), String>;
```

It must start multiple command lines passed as arguments, and write each of them followed by its output to `writer`, separated by empty lines. 
**The different commands' outputs must _not_ be mixed up.**

 * Commands must be executed in parallel. You must spawn a process for each command.
 * The standard error must be ignored.
 * Any error occurring when interacting with the system must be handled properly. Your program must never panic.
 * The output of a child must be displayed entirely as soon as it finishes execution, even if other commands are still in progress.

Example Usage:

```rust
fn main() {
        let cli1: &[String] = &["echo".to_string(), "a".to_string(), "b".to_string()];
        let cli2: &[String] = &["sleep".to_string(), "1".to_string()];
        let cli3: &[String] = &["cat".to_string(), "Cargo.toml".to_string()];
        let command_lines = vec![cli1, cli2, cli3];

    multiplexer(&mut std::io::stdout(), &command_lines);
}
```

Expected Output:
```txt
>_ cargo run
echo a b
a b

cat Cargo.toml
[package]
name = "ex04"
version = "0.1.0"

sleep 1
```

## Exercise 05: GET

```rust
// allowed symbols
use std::{
    io::Write,
    net::TcpStream,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex05/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Create a **function** with the following signature:

```rust
pub fn get<W: std::io::Write>(writer: &mut W, address: &str) -> Result<(), String>;
```

It must send an `HTTP/1.1` request and write the response to `writer`.

Example Usage:

```rust
fn main() {
    get(&mut std::io::stdout(), "localhost");
}
```

Expected Output:
```txt
$ cargo run
HTTP/1.1 200 OK
Server: nginx/1.24.0 (Ubuntu)
Date: Sun, 29 Dec 2024 15:48:11 GMT
Content-Type: text/html
Content-Length: 615
Last-Modified: Wed, 06 Nov 2024 10:04:23 GMT
Connection: close
ETag: "672b3f27-267"
Accept-Ranges: bytes

<!DOCTYPE html>
<html>
...
</html>
```

 * The function must send *valid* HTTP/1.1 requests.
 * Only the `GET` method is required.

**Note:** you should probably ask the server to `close` the `Connection` instantly to avoid having to detect the end of the payload.

## Exercise 06: ft_strings

```rust
// allowed symbols
use std::{
    io::Write,
    fs::read,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex06/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Create a **function** with the following signature:
```rust
pub fn strings<W: std::io::Write>(writer: &mut W, path: &str, z: bool, min: Option<usize>, max: Option<usize>) -> Result<(), String>;
```

It must read an arbitrary binary file and write printable UTF-8 strings it finds into `writer`.

Example Usage:

```sh
$ echo 'int main() { return 0; }' > test.c && cc test.c
```

```rust
fn main() {
    strings(&mut std::io::stdout, "./a.out", false, None, None);
}
```

Example Output:
```txt
>_ cargo run
ELF
>
М
@
+F
@
8
@
-
,
...
```

* A *printable UTF-8 string* is only composed of non-control characters (`TAB` _is_ a control character!).

The function must have the following options, passed to it as arguments:

* `z` filters out strings that are not null-terminated.
* `min` filters out strings where `string.len() <= min`.
* `max` filters out strings where `string.len() >= max`.

Your function must never panic when interacting with the file system. Handle errors properly.

Some level of input sanitization will be necessary - a `min` value higher than the `max` value does not make much sense.

Test your function with different binary files and option combinations to verify functionality.

## Exercise 07: Pretty Bad Privacy

```rust
// allowed symbols
use std::{
    io::{Read, Write, stdin, stdout, stderr},
    fs::File,
};
use rand::*;
use rug::*;

const allowed_dependencies = ["rug(v1.19.0)", "rand(v0.8.5)"];
const turn_in_directory = "ex07/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Write a **library** with the 3 following functions:

```rust
pub fn gen_keys(pub_key_path: &str, priv_key_path: &str)  -> Result<(), String>;
pub fn encrypt<R: std::io::Read, W: std::io::Write>(input: &mut R, writer: &mut W, pub_key_path: &str) -> Result<(), String>;
pub fn decrypt<R: std::io::Read, W: std::io::Write>(input: &mut R, writer: &mut W, priv_key_path: &str) -> Result<(), String>;
```

### `gen_keys`

In order to generate keys, your program must perform the following steps:

1. Generate two random prime numbers ($p$ and $q$).
2. Calculate $M = p \times q$.
2. Calculate $Phi = (p - 1) \times (q - 1)$.
4. Pick a random number $E$, such that:
    * $E < Phi$
    * $E$ and $Phi$ are coprime
    * $E$ and $M$ are coprime
5. Calculate $D$ as the multiplicative inverse of $E \mod Phi$.

The resulting keys are:

* Private key: $(D, M)$
* Public key: $(E, M)$

### `encrypt` and `decrypt`

* Encryption: $encrypt(m) = m^E \mod M$
* Decryption: $decrypt(m') = m'^D \mod M$

For any $m < M$, $decrypt(encrypt(m)) = m$.

### Key File Format

When saving keys to files, use the following format:

```plaintext
E/D
M
```

Where $E/D$ is the encryption or decryption component, and $M$ is the modulus.

### Encoding

To handle messages of arbitrary length:

1. Let $C$ be the largest integer such that $255^C < M$
2. **For encryption**:
    * Read $C$ bytes at a time from the input.
    * Treat these bytes as a base-256 number.
    * Encrypt this number using the encryption function.
    * Encode the result into $B + 1$ bytes in the output.
3. **For decryption**:
    * Read $B + 1$ bytes at a time from the input.
    * Treat these bytes as a base-256 number.
    * Decrypt this number using the decryption function.
    * Encode the result into $C$ bytes in the output.

_Note: Choose appropriate sizes for your numbers. The `rug` crate provides many integer sizes._

---
**License Notice:**
This file contains content licensed under two different terms:
- The MIT License applies to the original content (see `LICENSES/MIT-rust-subjects.txt`).
- The Creative Commons Attribution-ShareAlike 4.0 International (CC BY-SA 4.0 applies to any modifications or additions (see `LICENSES/CC-BY-SA-4.0.txt`).

When distributing modified versions, you must comply with both the MIT License and the CC BY-SA 4.0.
For complete details, refer to the main licensing file of Shortinette.
---
