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

## Exercise 00: Didn't Panic!

```txt
turn-in directories:
    ex00/

files to turn in:
    src/main.rs  Cargo.toml

allowed symbols:
    std::io::{stdout, Write}
    std::writeln
```

Create a **program** that prints integers from 1 to 10.

```txt
>_ cargo run
1
2
3
4
5
6
7
8
9
10
```

The program must *never* panic.

```txt
>_ cargo run | true
>_ 
```

## Exercise 01: Tee-Hee

```txt
turn-in directory:
    ex01/

files to turn in:
    src/main.rs  Cargo.toml

allowed symbols:
    std::io::{Write, Read, stdin, stdout}
    std::io::{Stdout, StdoutLock, Stdin, StdinLock}
    std::io::{Error, Result}
    std::fs::File  std::env::args
    std::vec::Vec  std::string::String
    std::iter::*
    std::{print, println, eprintln}
```

Create a **program** that reads the standard input, and copies it to the standard output, as well as
to every file specified in command-line arguments.

Example:

```txt
>_ echo "Hello, World!" | cargo run -- a b c
Hello, World!
>_ cat a b c
Hello, World!
Hello, World!
Hello, World!
```

You program must not panic when interacting with the file system. All errors must be handled
properly. You are free to choose what to do in that case, but you must *not* crash/panic.

## Exercise 02: Duh

```txt
turn-in directory:
    ex02/

files to turn in:
    src/main.rs  Cargo.toml

allowed symbols:
    std::fs::{metadata, Metadata, read_dir, DirEntry, ReadDir}
    std::path::Path  std::io::{Error, Result}
    std::env::args
    std::{print, println, eprintln}
```

Create a **program** that computes the total size of a directory or file. The program must write the
aggregated size of directories *in real-time*. As more files are taken in account in the count,
the total size must be updated in the terminal.

```txt
>_ cargo run -- ~
1.2 gigabytes
```

 * If a size is less than a kilobyte, it is written in bytes. (e.g. 245 bytes)
 * If a size is more than a kilobyte, it is written in kilobytes, with one decimal (e.g. `12.2 kilobytes`).
 * If a size is more than a megabyte, it is written in megabytes, with one decimal (e.g. `100.4 megabytes`).
 * If a size is more than a gigabyte, it is written in gigabytes, with one decimal (e.g. `23.9 gigabytes`).
 * For simplicty's sake, we'll assume that a kilobyte is `1000 bytes`, a megabyte is `1000 kilobytes`,
   etc.

Your program must not panic when interacting with the file system. Errors must be handled properly.

## Exercise 03: Pipe-Line

```txt
turn-in directory:
    ex03/

files to turn in:
    src/main.rs  Cargo.toml

allowed symbols:
    std::env::args
    std::process::Command
    std::os::unix::process::CommandExt
    std::io::{Read, stdin}
    std::vec::Vector
    std::iter::*
```

Create a **program** that takes a path and some arguments as an input, and spawns that process with:

1. The arguments passed in command-line arguments.
2. Each line of its standard input.

Example:

```rust
>_ << EOF cargo run -- echo -n
hello
test
EOF
hello test>_
```

The program invoked the `echo -n hello test` command.

Your program must not panic when interacting with the system, you must handle errors properly.

## Exercise 04: Command Multiplexer

```txt
turn-in directory:
    ex04/

files to turn in:
    src/main.rs  Cargo.toml

allowed symbols:
    std::env::args  std::iter::*
    std::process::{Command, Stdio, Child}
    std::vec::Vec
    std::io::{stdout, Write, Read}
    std::{write, writeln}
    std::eprintln
```

Create a **program** that starts multiple commands, gathers their output and then prints it to the
standard output. **The different commands' outputs must _not_ be mixed up.**

Example:

_Note: You are free to format this exercise's output as you like, as long as **each command's output is separated from the others by at least one newline**._
```txt
>_ cargo run -- echo a b , sleep 1 , echo b , cat Cargo.toml , cat i-dont-exit.txt
===== cat i-dont-exit.txt ====

===== echo a b =====
a b

===== echo b =====
b

==== cat Cargo.toml =====
[package]
name = "ex03"
version = "0.1.0"
...

==== sleep 1 =====

```

 * Commands must be executed in parallel.
 * The standard error must be ignored.
 * Any error occuring when interacting with the system must be handled properly. Your program must never panic.
 * The output of a child must be displayed entirely as soon as it finishes execution, even when other commands are still in progress.

## Exercise 05: GET

```txt
turn-in directory:
    ex05/

files to turn in:
    src/main.rs  Cargo.toml

allowed symbols:
    std::env::args
    std::net::{TcpStream, SocketAddr, ToSocketAddrs}
    std::io::{Write, Read, stdout}
```

Create a **program** that sends an `HTTP/1.1` request and prints the response.

Example:

_Note: You are free to format this exercise as you like, as long as the HTTP/1.1 status code and the Content-Length header are displayed._

```txt
>_ cargo run -- https://github.com/42-Short
HTTP/1.1 200 OK
Server: tiny-http (Rust)
Date: Sat, 04 Feb 2023 12:40:33 GMT
Content-Length: ...
...
<html>
...
```

 * The program must send *valid* HTTP/1.1 requests.
 * Only the `GET` method is required.

**Note:** you should probably ask the server to `close` the `Connection` instantly to avoid
having to detect the end of the payload.

## Exercise 06: ft_strings

```txt
turn-in directory:
    ex06/

files to turn in:
    src/main.rs  Cargo.toml

allowed symbols:
    std::env::args
    std::fs::read
    std::str::{from_utf8, Utf8Error}
```

Create a **program** that reads an arbitrary binary file and prints printable UTF-8 strings it finds.

Example:

```txt
>_ cargo run -- ./a.out
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

The program must have the following options:

* `-z` filters out strings that are not null-terminated.
* `-m <min>` filters out strings that are strictly smaller than `min`.
* `-M <max>` filters out strings that are strictly larger than `max`.

Implementation requirements:
1. Do not panic when interacting with the file system. Handle errors properly.
2. Use only the allowed symbols listed above.
3. Implement all specified options.
4. Ensure correct handling of various binary file types.

Test your program with different binary files and option combinations to verify functionality.

## Exercise 07: Pretty Bad Privacy

```txt
turn-in directory:
    ex07/

files to turn in:
    src/main.rs Cargo.toml

allowed dependencies:
    rug(v1.19.0)
    rand(v0.8.5)

allowed symbols:
    std::vec::Vec
    std::env::args
    std::io::{stdin, stdout, stderr, Write, Read}
    std::fs::File
    rand::*
    rug::*
```

Write a **program** that behaves in the following way:

```sh
# Generate key pair
>_ cargo run -- gen-keys my-key.pub my-key.priv

# Encrypt a message
>_ << EOF cargo run -- encrypt my-key.pub > encrypted-message
This is a very secret message.
EOF

# Decrypt a message
>_ cat encrypted-message | cargo run -- decrypt my-key.priv
This is a very secret message.
```

### Key Generation

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

### Encryption & Decryption

* Encryption: `encrypt(m) = m^E % M`
* Encryption: `encrypt(m') = m'^D % M`

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
