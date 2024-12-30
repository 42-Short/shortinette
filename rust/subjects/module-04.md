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

## Exercise 00: Wait thats it?

```txt
turn-in directories:
    ex00/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    none
```

Create `Outcome` and `Maybe` which should mimic `Result` and `Option` so that test compiles and runs

```rust
#[cfg(test)]
mod test {
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

```txt
turn-in directory:
    ex01/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    std::io::{Write, Read, stdin, stdout}
    std::io::{Stdout, StdoutLock, Stdin, StdinLock}
    std::io::{Error, Result}
    std::fs::File 
    std::vec::Vec
    std::string::String
    std::iter::*
    std::{print, println, eprintln}
```

Write a **function** that copies `input` to `writer` and all filenames provided as arguments.

Your function must have the following signature:

```rust
pub fn tee<R: std::io::Read, W: std::io::Write>(input: &mut R, writer: &mut W, filenames: Vec<String>);
```

Example:

```rust
fn main() {
    let args: Vec<String> = std::env::args().skip(1).collect();
    tee(&mut std::io::stdin(), &mut std::io::stdout().lock(), &args);
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

You program must not panic when interacting with the file system. All errors must be handled properly. You are free to choose what to do in that case, but you must *not* crash/panic.

## Exercise 02: Duh

```txt
turn-in directory:
    ex02/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    std::fs::{metadata, Metadata, read_dir, DirEntry, ReadDir}
    std::path::Path
    std::{print, println, eprintln}
```

Create a **function** that computes the total size of a directory or file. Once computed, write the size followed by a newline, and formatted like in the examples below.

Your function must have the following signature:
```rust
pub fn duh<W: std::io::Write>(writer: &mut W, basedir: &str);
```

Example Usage:

```rust
fn main() {
    let arg: String = std::env::args().nth(1);
    duh(&mut std::io::stdout(), &arg);
}

// Output:
// 1.2 gigabytes
```

 * If a size is less than a kilobyte, it is written in bytes. (e.g. 245 bytes)
 * If a size is more than a kilobyte, it is written in kilobytes, with one decimal (e.g. `12.2 kilobytes`).
 * If a size is more than a megabyte, it is written in megabytes, with one decimal (e.g. `100.4 megabytes`).
 * If a size is more than a gigabyte, it is written in gigabytes, with one decimal (e.g. `23.9 gigabytes`).
 * For simplicity, you will assume that a kilobyte is `1000 bytes`, a megabyte is `1000 kilobytes`,
   etc.

Your function must never panic when interacting with the file system. Errors must be handled properly.

## Exercise 03: Pipe-Line

```txt
turn-in directory:
    ex03/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    std::process::Command
    std::os::unix::process::CommandExt
    std::io::{Read, stdin}
    std::vec::Vec
    std::iter::*
```

Create a **function** with the following signature:

```rust
pub fn pipeline<R: std::io::Read, W: std::io::Write>(input: &mut R, writer: &mut W, args: Vec<String>);
```

It must spawn a process using the path (`args[0]`, arguments (`args[1..]`)), and input (`input`) passed as arguments,
and write its output to `writer`.

Example Usage:

```rust
fn main() {
    pipeline(&mut std::io::stdin(), &mut std::io::stdout(), vec![String::from("echo"), String::from("-n")]);
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

```txt
turn-in directory:
    ex04/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    std::iter::*
    std::process::{Command, Stdio, Child}
    std::vec::Vec
    std::io::{stdout, Write, Read}
    std::{write, writeln}
    std::eprintln
```

Create a **function** with the following signature:

```rust
pub fn multiplexer<W: std::io::Write>(writer: &mut W, command_lines: Vec<Vec<String>>);
```

It must start multiple command lines passed as arguments, and write each of them followed by its output to `stdout`, separated by empty lines, to `writer`. 
**The different commands' outputs must _not_ be mixed up.**

 * Commands must be executed in parallel. You must spawn a process for each command.
 * The standard error must be ignored.
 * Any error occuring when interacting with the system must be handled properly. Your program must never panic.
 * The output of a child must be displayed entirely as soon as it finishes execution, even when other commands are still in progress.

Example Usage:

```rust
fn main() {
    let command_lines = vec![
        vec![String::from("echo"), String::from("a"), String::from("b")],
        vec![String::from("sleep"), String::from("1")],
        vec![String::from("cat"), String::from("Cargo.toml")],
    ];

    multiplexer(&mut std::io::stdout(), command_lines);
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

```txt
turn-in directory:
    ex05/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    std::net::{TcpStream, SocketAddr, ToSocketAddrs}
    std::io::{Write, Read, stdout}
```

Create a **function** with the following signature:

```rust
pub fn get<W: std::io::Write>(writer: &mut W, address: &str);
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

```txt
turn-in directory:
    ex06/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    std::fs::read
    std::str::{from_utf8, Utf8Error}
```

Create a **function** with the following signature:
```rust
pub fn strings<W: std::io::Write>(writer: &mut W, path: &str, z: bool, min: Option<usize>, max: Option<usize>);
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

* A *printable UTF-8 string* is only composed of non-control characters.

The function must have the following options, passed to it as arguments:

* `z` filters out strings that are not null-terminated.
* `min` filters out strings that are strictly smaller than `min`.
* `max` filters out strings that are strictly larger than `max`.

Your function must never panic when interacting with the file system. Handle errors properly.

Test your function with different binary files and option combinations to verify functionality.

## Exercise 07: Pretty Bad Privacy

```txt
turn-in directory:
    ex07/

files to turn in:
    src/lib.rs Cargo.toml

allowed dependencies:
    rug(v1.19.0)
    rand(v0.8.5)

allowed symbols:
    std::vec::Vec
    std::io::{stdin, stdout, stderr, Write, Read}
    std::fs::File
    rand::*
    rug::*
```

Write a **library** with the 3 following functions:

```rust
pub fn gen_keys(pub_key_path: &str, priv_key_path: &str);
pub fn encrypt<R: std::io::Read, W: std::io::Write>(input: &mut R, writer: &mut W, pub_key_path: &str);
pub fn decrypt<R: std::io::Read, W: std::io::Write>(input: &mut R, writer: &mut W, priv_key_path: &str);
```

### `gen_keys`

In order to generate keys, your program must perform the following steps:

1. Generate two random prime numbers (`p` and `q`).
2. Calculate `M = p * q`.
2. Calculate `PHI = (p - 1) * (q - 1)`.
4. Pick a random number `E`, such that:
    * `E < Phi`
    * `E` and `Phi` are coprime
    * `E` and `M` are coprime
5. Calculate `D`, as the multiplicative inverse of `E` modulo `Phi`.

The resulting keys are:

* Private key: `(D, M)`
* Public key: `(E, M)`

### `encrypt` and `decrypt`

* Encryption: `encrypt(m) = m^E % M`
* Decryption: `encrypt(m') = m'^D % M`

For any `m < M`, `decrypt(encrypt(m)) == m` should hold true.

### Key File Format

When saving keys to files, use the following format:

```plaintext
E/D
M
```

Where `E/D` is the encryption or decryption component, and `M` is the modulus.

### Encoding

To handle messages of arbitrary length:

1. Let `C` be the largest integer such that ``255^C < M`
2. **For encryption**:
    * Read `C` bytes at a time from the input.
    * Treat these bytes as a base-256 number.
    * Encrypt yhis number using the encryption function.
    * Encode the result into `B + 1` bytes in the output.
3. **For decryption**:
    * Read `B + 1` bytes at a time from the input.
    * Treat these bytes as a base-256 number.
    * Decrypt this number using the decryption function.
    * Encode the result into `C` bytes in the output.

_Note: Choose appropriate sizes for your numbers. The `rug` crate provides many integer sizes._

---
**License Notice:**
This file contains content licensed under two different terms:
- The MIT License applies to the original content (see `LICENSES/MIT-rust-subjects.txt`).
- The Creative Commons Attribution-ShareAlike 4.0 International (CC BY-SA 4.0 applies to any modifications or additions (see `LICENSES/CC-BY-SA-4.0.txt`).

When distributing modified versions, you must comply with both the MIT License and the CC BY-SA 4.0.
For complete details, refer to the main licensing file of Shortinette.
---
