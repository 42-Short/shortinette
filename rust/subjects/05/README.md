# Module 05: Concurrency

## Foreword

```txt
error: pineapple doesn't go on pizza
 --> main.rs:6:18
  |
 6|     let _: Pizza<Pineapple>;
  |            ----- ^^^^^^^^^
  |            |
  |            this is the pizza you ruined
  |
  = note: `#[forbid(bad_taste)]` on by default
  = note: you're a monster
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

## Exercise 00: Cellular

```rust
// allowed symbols
use std::cell::Cell;

const allowed_dependencies = [""];
const turn_in_directory = "ex00/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Write a **function** with the following prototype:

```rust
pub fn swap_u32(a: &Cell<u32>, b: &Cell<u32>);
```

* The `swap_u32` function must swap the integers it is given.

Example:

```rust
#[cfg(test)]
mod tests {
    #[test]    
    fn example() {
        let a = Cell::new(1);
        let b = Cell::new(3);

        swap_u32(&a, &b);

        assert_eq!(a.get(), 3);
        assert_eq!(b.get(), 1);
    }
}
```

Let's complicate things a bit!

```rust
pub fn swap_string(a: &Cell<String>, b: &Cell<String>);
```

* The `swap_string` function must swap the strings it is given.

Example:

```rust
#[cfg(test)]
mod tests {
    #[test]    
    fn example() {
        let a = Cell::new("ABC".into());
        let b = Cell::new("DEF".into());

        swap_string(&a, &b);

        assert_eq!(a.into_inner(), "DEF");
        assert_eq!(b.into_inner(), "ABC");
    }
}
```

## Exercise 01: Atomic

```rust
// allowed symbols
use std::sync::atomic::{AtomicU8, Ordering};

const allowed_dependencies = [""];
const turn_in_directory = "ex01/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Create a type named `Unique`.

```rust
#[derive(Debug, PartialEq, Eq)]
pub struct Unique(u8);

impl Unique {
    pub fn new() -> Option<Self>;
    pub fn id(&self) -> u8;
}
```

* There can be no two `Unique` instances with the same identifier (`u8`).
* `new` must create a new, unique instance of `Unique`.
* `id` must return the unique id of the instance of `Unique`.
* It must be possible to `Clone` a `Unique`, and the created `Unique` must still be unique.
* Trying to create a `Unique` when no more identifiers are available causes the function to return `None`.
* Since `Unique` uses a `u8`, it is only possible to create up to `255` instances of `Unique`. Think about why this is the case.

Example:

```rust
fn main()
{
    let a = Unique::new();
    let b = Unique::new();
    let c = Unique::new();

    println!("{a:?}");
    println!("{b:?}");
    println!("{c:?}");

    let d = a.clone();
    let e = c.clone();

    println!("{d:?}");
    println!("{e:?}");
}
```

Would produce the following output:

```txt
>_ cargo run
Unique(0)
Unique(1)
Unique(2)
Unique(3)
Unique(4)
```

What atomic memory ordering did you use? Why?

## Exercise 02: Last Error

```rust
// allowed symbols
use std::{
    thread_local,
    cell::Cell,
    marker::Copy,
    clone::Clone,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex02/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Create an `Error` enum with the following variants:

```rust
pub enum Error {
    Success,
    FileNotFound,
    IsDirectory,
    WriteFail,
    ReadFail,
}

impl Error {
    pub fn last() -> Self;
    pub fn make_last(self);
}
```

 * `last` must return the calling thread's last `Error` instance. If `make_last` has never been called before, `Error::Success` is returned.
 * `make_last` must set the calling thread's last `Error` instance.

## Exercise 03: A Philosopher's Tiny Brain

```rust
// allowed symbols
use std::{
    sync::Arc,
    sync::mpsc::{sync_channel, SyncSender, Receiver},
    thread::{spawn, sleep},
    time::Duration,
};
use ftkit::ARGS;

const allowed_dependencies = ["ftkit"];
const turn_in_directory = "ex03/";
const files_to_turn_in = ["src/main.rs", "Cargo.toml"];
```

Create a **program** that works in the following way:

```txt
>_ cargo run -- 3
cakes
the philosopher is thinking about cakes
code
42
the philosopher is thinking about code
a
b
c
the philosopher's head is full
the philosopher is thinking about 42
the philosopher is thinking about a
the philosopher is thinking about b
^C
>_
```

* The program must ask the user to enter words in the standard input.
* Each time a word is entered, it is saved in the philosopher's brain.
* If the brain is full, an error is displayed and the word is not added to the brain.
* When a word is available in the brain, the philosopher thinks about it for 5 seconds.
* The program runs until it receives `EOF`.
* The size of the philosopher's brain is provided as a command-line argument.
* The nature of `sync_channel` makes it possible for the philosopher to think about a single topic even if the brain size is `0`. Think about it.
* Read the example very carefully.

## Exercise 04: Logger

```rust
// allowed symbols
use std::{
    sync::{Arc, Mutex},
    thread:spawn,
    io::Write,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex04/";
const files_to_turn_in = ["src/main.rs", "Cargo.toml"];
```

Create a `Logger` type.

```rust
pub struct Logger<W> {
    buffer: Box<[u8]>,
    writer: W,
}

impl<W> Logger<W> {
    pub fn new(threshold: usize, writer: W) -> Self;
    pub fn buffer(&self) -> &Box<[u8]>
    pub fn writer(&self) -> &W
}
```

 * `new` must create a new `Logger` with a buffer of size `threshold` and the given `W` instance.

In order to avoid performing too many `write` system calls, you should first write the messages to an internal `buffer`, and THEN, write everything to the given writer.

```rust
impl<W: io::Write> Logger<W> {
    pub fn log(&mut self, message: &str) -> io::Result<()>;
    pub fn flush(&mut self) -> io::Result<()>;
}
```

 * `log` must try to write `message` to its internal buffer. When the buffer is full, everything must be sent to the specified `io::Write` implementation. After that the buffer is cleared for new data to be added. A `\n` is automatically added at the end of the message.
 * `flush` must properly send the content of the buffer and clears it.

Create a `main` function spawning 10 threads. Each thread must try to write to the standard output using the same `Logger<Stdout>` 10 times.

```txt
>_ cargo run
hello 0 from thread 1!
hello 0 from thread 2!
hello 0 from thread 0!
hello 0 from thread 6!
hello 1 from thread 1!
...
```

 * A message from any given thread must not appear within the message of another thread.

## Exercise 05: PI * Rayon * Rayon

```rust
// allowed symbols
use std::{
    iter::*,
    println,
    env::args,
    time::Instant,
};

const allowed_dependencies = ["rand", "rayon"];
const turn_in_directory = "ex05/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Let's look at some popular third-party crates!

First, let's create a single threaded **program** that uses [Monte Carlo's method](https://en.wikipedia.org/wiki/Monte_Carlo_method#Overview) to compute PI. The program takes a single argument: the number of points to sample.

Try to write this algorithm without a `for` loop. Instead, rely on chained iterators. This will make it easier for you in the second part of the exercise.

```txt
>_ RUSTFLAGS="-C opt-level=3 -C target-cpu=native" cargo run -- 1000000
pi: 3.1413
duration: 722ms
```

Even for as little as a million points, the algorithm is already pretty slow. Try to speed it up a little using the [`rayon`](https://crates.io/crates/rayon) crate.

```txt
>_ RUSTFLAGS="-C opt-level=3 -C target-cpu=native" cargo run -- 1000000
pi: 3.144044
duration: 147ms
```

## Exercise 06: 404 Not Found

```rust
// allowed symbols
use std::{
    thread::{spawn, JoinHandle},
    sync::{Arc, RwLock, mpsc::{Sender, Receiver, channel}},
    net::{TcpListener, SockerAddr},
};

const allowed_dependencies = [""];
const turn_in_directory = "ex06/";
const files_to_turn_in = ["src/lib.rs", "src/main.rs", "Cargo.toml"];
```

Create a `ThreadPool` type.

_Please note that the struct attributes are only `pub` due to tester requirements, normally you would keep them private._

```rust
type Task = Box<dyn 'static + Send + FnOnce()>;

pub struct ThreadPool {
    pub threads: Vec<JoinHandle<()>>,
    pub should_stop: Arc<RwLock<bool>>,
    pub task_sender: Sender<Task>,
}

impl ThreadPool {
    fn new(thread_count: usize) -> Self;
    fn spawn_task<F>(task: F) -> Result<(), /* ... */>
    where
        F: 'static + Send + FnOnce();
}
```

 * The `new` function must create a new `ThreadPool` instance by spawning `thread_count` threads.
 * The `spawn_task` function must send the task to a thread in the thread pool.
 * When a `ThreadPool` is dropped, its threads must stop. If any of the threads panicked, en error should be printed to standard error.

When a thread is not executing a task, it waits until one is available and executes it.

Let's create a multithreaded HTTP server!

* Your **program** must listen on an address and port specified in command-line arguments:

```txt
>_ cargo run -- 127.0.0.1:8080
```

* When a connection is received, the server must respond with a "404 Not Found" page.
* Every new connection to the server must be handled on a thread pool.

```txt
>_ curl 127.0.0.1:8080/
This page does not exist :(
>_
```

* This must also work in a regular web browser.

## Exercise 07: Rendez-Vous

```rust
// allowed symbols
use std::{
    sync::{CondVar, Mutex},
    mem::{Replace, swap},
};

const allowed_dependencies = [""];
const turn_in_directory = "ex07/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Let's create a "Rendez-Vous" primitive in Rust.

Example:

```rust
// THREAD 1
let a = rdv.wait(42u32); // if thread2 has not arrived yet, wait.
// We know that thread2 has arrived too!
assert_eq!(a, 21u32);

// THREAD 2
let a = rdv.wait(21u32); // if thread1 has not arrived yet, wait.
// We know that thread1 has arrived too!
assert_eq!(a, 42u32);
```

The "Rendez-Vous" must be defined as follows:

```rust
pub struct RendezVous<T> { /* ... */ }

impl<T> RendezVous<T> {
    const fn new() -> Self;
    pub fn wait(&self, value: T) -> T;
    pub fn try_wait(&self, value: T) -> Result<T, T>;
}
```

 * `new` must create a new `RendezVous`.
 * `wait` must block until it has been called twice. The first call returns the value passed by the
   second call, and the second call returns the value passed by the first call.
 * `try_wait` checks whether someone is waiting using `wait`. If so, the values are exchanged and
   `Ok(_)` is returned. Otherwise, the input value is returned in the `Err(_)`.
 * `RendezVous` must be reusable, allowing multiple dates even after the first exchange.
 * A thread must never [spin](https://en.wikipedia.org/wiki/Busy_waiting) when waiting for
   something to happen!

---
**License Notice:**
This file contains content licensed under two different terms:
- The MIT License applies to the original content (see `LICENSES/MIT-rust-subjects.txt`).
- The Creative Commons Attribution-ShareAlike 4.0 International (CC BY-SA 4.0 applies to any modifications or additions (see `LICENSES/CC-BY-SA-4.0.txt`).

When distributing modified versions, you must comply with both the MIT License and the CC BY-SA 4.0.
For complete details, refer to the main licensing file of Shortinette.
---
