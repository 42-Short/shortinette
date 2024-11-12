# Module 02: Structure

## Foreword

Upon a warm autumn day, Brother Farbold was walking through the Gardens of Abstraction near the
recently opened Temple of Rust. He passed many curious and elaborate displays of complex
abstractions, when he chanced upon a monk toiling away at her own display. Her robes identified her
as a neophyte, though her creased, worn hands spoke of many, many years of experience. Around her
were scatted dozens of white stones of all shapes and sizes along with several different rakes.
Attached to each was a single, small, paper label with a word written on it.

Intrigued, he walked closer.

"Ho!" he said, "what is your name?"

"I am Neophyte Maran."

"What problem are you working on, if I might enquire?"

"I am attempting to translate a previous design that it might be suitable for the temple," the
neophyte replied, concentrating on a large sheet of paper covered with sketches of boxes, connected
by many lines and arrows.

"I see many labels; what are those for?"

"They define the class of each object."

Brother Farbold paused. "What are classes?"

The neophyte gave him a look of confusion. "You do not know what classes are?"

"No," he replied. "I have heard such a phrase carried on the wind from far away, but in this place,
we do not use it."

"Ah, I see," said Maran. "In the temple where I came from, we use classes to define all the
properties and behaviours of a thing, that many may share a single definition."

"A useful device, then. In such a scheme, as a simple gardener, I would have a tend behaviour,
then?"

The neophyte nodded and picked up a small rock. "This is a Rock," she said, showing Brother Farbold
the label. "It has no methods, for all it does is lie on the ground until operated on, but it has
properties such as weight and volume."

"So a rock cannot do anything at all?"

"Nothing. Unlike this Rake, which possesses the ability to rake one or more Rocks. This allows me to
model a garden and its evolution in its entirety."

"Indeed? This simple model allows you to capture all possible interactions, then?"

"Yes!" exclaimed Neophyte Maran proudly. "I have used this for many years with great success."

"But what if the rock wishes to roll, or shine, or talk? Classes sound restrictive; I prefer to
define smaller behaviours, that each thing may pick and choose what it wishes to do."

The neophyte scoffed. "An object can do only that which it was designed to do, nothing more. You
said you were a gardener; are you tasked to instruct new students, or to cook meals?"

"I am not; I am tasked with gardening."

"With classes, such is expressed by your class at design-time. As such, there can never be any
confusion! To change behaviours is to change the design of the system!"

"Facinating," said Farbold, nodding. The neophyte turned back to her work.

"I must admit I am having trouble finding how to express this concept in a way that the temple will
approve of, but I am su--" the neophyte was interrupted when one of the smaller white rocks struck
her head.

Holding her head, Maran turned around. "Did you throw a rock at me?" she demanded.

"Me?" said Brother Farbold. "Of course not. I am but a simple Gardener; I know nothing of throwing
Rocks."

"The rock could not have thrown itself!"

"And yet, your design does not allow for such a thing to have happened. Surely, you must be
mistaken: no rock was thrown."

Neophyte Maran grumbled and turned back to her sketches.

A moment later, she flinched away as a long wooden pole struck her arm. "Did you...?"

"A Rake can only rake Rocks; it cannot "hit Neophytes". You must be imagining things.

"But you..."

"Ho! Brother Farbold!" another monk cried from the other end of the garden. "Brother Schieb has
taken ill; would you be so kind as to cook the evening meal in his stead?"

"I would be happy to, Sister Doure; I am just now finishing my instruction of this neophyte."

At that moment, Maran was enlightened.

*The [second Rust Koan](https://users.rust-lang.org/t/rust-koans/2408/2).*

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
 Tests can use the symbols & attributes you want, even if they are not specified in the `allowed symbols` section. **However**, tests should **not** introduce **any additional external dependencies** beyond those already required by the subject.

## Exercise 00: Dimensional Analysis

```txt
turn-in directory:
    ex00/

files to turn in:
    src/main.rs  Cargo.toml
```

Copy/Paste the following code and make it compile by adding type alias definitions.

```rust
fn seconds_to_minutes(seconds: Seconds) -> Minutes {
    seconds / 60.0
}

fn main() {
    let s: Seconds = 120.0;
    let m: Minutes = seconds_to_minutes(s);

    println!("{s} seconds is {m} minutes");
}
```

```txt
>_ cargo run
120 seconds is 2 minutes
```

## Exercise 01: A Point In Space

```txt
turn-in directory:
    ex01/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    f32::sqrt
```

Defines the following type:

```rust
struct Point {
    x: f32,
    y: f32,
}
```

Implement the following inherent functions:

* `new`, which creates a new `Point` with the coordinates passed to it.
* `zero`, which creates a new `Point` at coordinates `(0, 0)`.
* `distance`, which computes the distance between two existing points.
* `translate`, which adds the vector `(dx, dy)` to the coordinates of the point.

```rust
impl Point {
    fn new(x: f32, y: f32) -> Self;
    fn zero() -> Self;
    fn distance(&self, other: &Self) -> f32;
    fn translate(&mut self, dx: f32, dy: f32);
}
```

## Exercise 02: Where's My Pizza?


```txt
turn-in directory:
    ex02/

files to turn in:
    src/lib.rs  Cargo.toml
```

* Once a pizza has been ordered, it takes two days before the cook start working on it.
* Making a pizza takes roughly 5 days.
* Once the pizza is ready, the only delivery man must pick it up. It takes 3 days on average.
* Delivering the pizza always takes a whole week.

Define the following type:

```rust
enum PizzaStatus {
    Ordered,
    Cooking,
    Cooked,
    Delivering,
    Delivered,
}
```

It must have the following inherent methods.

```rust
impl PizzaStatus {
    fn from_delivery_time(ordered_days_ago: u32) -> Self;
    fn get_delivery_time_in_days(&self) -> u32;
}
```

* `from_delivery_time` predicts the status of a pizza that was ordered `ordered_days_ago` days ago.
* `get_delivery_time_in_days` returns the estimated time before the pizza is delivered, in days. The
**worst case** (longest delivery time) is always returned.

## Exercise 03: Dry Boilerplates

```txt
turn-in directory:
    ex03/

files to turn in:
    src/main.rs  Cargo.toml

allowed symbols:
    std::clone::Clone  std::cmp::{PartialOrd, PartialEq}
    std::default::Default  std::fmt::Debug
```

Create a type, may it be a `struct` or an `enum`. You simply have to name it `MyType`.

```rust
fn main() {
    let instance = MyType::default();

    let other_instance = instance.clone();

    println!("the default value of MyType is {instance:?}");
    println!("the clone of `instance` is {other_instance:#?}");
    assert_eq!(
        instance,
        other_instance,
        "the clone isn't the same :/"
    );
    assert!(
        instance == other_instance,
        "why would the clone be less or greater than the original?",
    );
}
```

Copy the above `main` function and make it compile and run. You are not allowed to use the `impl`
keyword!

## Exercise 04: Todo List

```txt
turn-in directory:
    ex04/

files to turn in:
    src/main.rs  Cargo.toml

allowed dependencies:
    ftkit

allowed symbols:
    std::{print, println}
    std::io::stdout
    std::io::Stdout::{flush}
    std::vec::Vec::{new, push, remove, clear, len, is_empty}
    std::string::String::as_str
    str::{to_string, parse, len, is_empty, trim, strip_prefix, strip_suffix}
    ftkit::{read_line, read_number}
    std::result::Result
```

Create a simple TODO-List application.

1. `Command` enum:
Define it as shown below. Add the `from_input` function to parse
commands and return the correct variant.
```rust
enum Command {
    Todo(String),   // Command: "TODO"
    Done(usize),    // Command: "DONE"
    Purge,          // Command: "PURGE"
    Quit,           // Command: "QUIT"
}

impl Command {
    fn from_input(input: &str) -> Result<Self, String> {
        ...
    }
}
```
2. `TodoList` Struct:
Define the struct as shown below.  
```rust
struct TodoList {
    todos: Vec<String>,
    dones: Vec<String>,
}

impl TodoList {
    fn new() -> Self;

    fn display(&self);
    fn add(&mut self, todo: String);
    fn done(&mut self, index: usize);
    fn purge(&mut self);
}
```

* `add` appends a new task.
* `done` removes the task at `index` from `todos` and pushes it to `dones`.
* `purge` clears the `dones` vector.
* `display` prints the content of the todolist to the user.

3. `main` Function:
Write a `main` function, responsible for using both `TodoList` and `Command`. The content of the todolist must be displayed to the user before each prompt.

You may design the interface you want to in this exercise. Here is an example.

```txt
>_ cargo run

TODO star shortinette (https://github.com/42-Short/shortinette)

    0 [ ] star shortinette (https://github.com/42-Short/shortinette)

TODO finish this module

    0 [ ] star shortinette (https://github.com/42-Short/shortinette)
    1 [ ] finish this module

DONE 0

    0 [ ] finish this module
      [x] star shortinette (https://github.com/42-Short/shortinette)

PURGE

    0 [ ] finish this module

QUIT
```

## Exercise 05: Lexical Analysis

```txt
turn-in directory:
    ex06/

files to turn in:
    src/main.rs src/lib.rs  Cargo.toml

allowed dependencies:
    ftkit

allowed symbols:
    ftkit::ARGS
    std::option::Option
    std::fmt::Debug
    str::{strip_prefix, trim_start, char_indices, split_at, is_empty}
    char::is_whitespace
    std::{println, eprintln}
```

Create a simple token parser. It must be able to take an input string, and turn it into a list
of tokens.

Each token must be represented using the following `enum`:

```rust
#[derive(Debug, PartialEq)]
enum Token<'a> {
    Word(&'a str),
    RedirectStdout,
    RedirectStdin,
    Pipe,
}
```

* The character `>` produces a `RedirectStdout` token.
* The character `<` produces a `RedirectStdin` token.
* The character `|` produces a `Pipe` token.
* Any other character is part of a `Word`, unless it's a whitespace.
* Whitespaces are ignored.

Write the following function to your library:
```rust
fn next_token(s: &mut &str) -> Option<Token>;
```

* You may need to add some *lifetime annotations* to make the function compile properly.
* The `next_token` function either produce `Some(_)` value if a token is available, or `None` when
the input string contains no tokens. The part of `s` that has been consumed is stripped from the
original `&str`.

**Note:** You do not have to handle single or double quotes, nor do you have to care about escaping
with `\\`!

`next_token` must be usable in this way:

```rust
fn main() {
    let mut args = std::env::args();

    args.next();
    match args.next() {
        Some(arg) => {
            if args.next().is_some() {
                eprintln!("error: exactly one argument expected");
                return;
            }
            let mut arg_str: &str = &arg;
            while let Some(token) = next_token(&mut arg_str) {
                println!("{:?}", token); // Debug output of tokens
            }
        }
        None => eprintln!("error: exactly one argument expected"),
    }
}
```

```txt
>_ cargo run -- a b
error: exactly one argument expected
>_ cargo run -- "echo hello|cat -e> file.txt"
Word("echo")
Word("hello")
Pipe
Word("cat")
Word("-e")
RedirectStdout
Word("file.txt")
```

## Exercise 06: Inventory Management

```txt
turn-in directory:
    ex05/

files to turn in:
    src/lib.rs  Cargo.toml

allowed symbols:
    std::vec::Vec
    std::string::String
    std::collections::HashMap::*
    std::fmt::Debug
```
Create an inventory management system for a shop. Define the following types:
```rust
#[derive(Debug, Clone)]
struct Item {
    name: String,
    price: f32,
    quantity: u32,
}

struct Inventory {
    items: HashMap<String, Item>,
}
```
Implement the following methods for `Inventory`:
```rust
impl Inventory {
    fn new() -> Self;
    fn add_item(&mut self, item: Item) -> Result<(), String>;
    fn remove_item(&mut self, name: &str) -> Result<(), String>;
    fn update_quantity(&mut self, name: &str, new_quantity: u32) -> Result<(), String>;
    fn get_item(&self, name: &str) -> Option<&Item>;
    fn list_items(&self) -> Vec<&Item>;
    fn total_value(&self) -> f32;
}
```
* `new`: Creates a new empty inventory.
* `add_item`: Adds a new item to the inventory. If an item with the same name already exists, return an error.
* `remove_item`: Removes an item from the inventory by name. If the item doesn't exist, return an error.
* `update_quantity`: Updates the quantity of an existing item. If the item doesn't exist, return an error.
* `get_item`: Returns a reference to an item by name, or None if it doesn't exist.
* `list_items`: Returns a vector of references to all items in the inventory.
* `total_value`: Calculates and returns the total value of all items in the inventory (price * quantity for each item).

Implement a `Discountable` trait for `Item`:
```rust
trait Discountable {
    fn apply_discount(&mut self, percentage: f32);
}
```
Implement this trait for `Item` so that it reduces the price of the item by the given percentage. Invalid percentage values (`< 0 || > 100`) are considered **undefined behavior**. You are free to handle them as you please.

Your implementation should work with the following test function:
```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_inventory() {
        let mut inventory = Inventory::new();
        
        let item1 = Item { name: "Apple".to_string(), price: 0.5, quantity: 100 };
        let item2 = Item { name: "Banana".to_string(), price: 0.3, quantity: 150 };
        
        assert!(inventory.add_item(item1).is_ok());
        assert!(inventory.add_item(item2).is_ok());
        
        assert_eq!(inventory.list_items().len(), 2);
        
        assert!(inventory.update_quantity("Apple", 90).is_ok());
        assert_eq!(inventory.get_item("Apple").unwrap().quantity, 90);
        
        assert!((inventory.total_value() - 90.0).abs() < 0.001);
        
        let mut discounted_item = inventory.get_item("Banana").unwrap().clone();
        discounted_item.apply_discount(10.0);
        assert!((discounted_item.price - 0.27).abs() < 0.001);
        
        assert!(inventory.remove_item("Apple").is_ok());
        assert_eq!(inventory.list_items().len(), 1);
    }
}
```

## Exercise 07: The Game Of Life

```txt
turn-in directory:
    ex07/

files to turn in:
    src/main.rs  Cargo.toml

allowed dependencies:
    ftkit

allowed symbols:
    ftkit::ARGS ftkit::random_number
    std::{println, print}
    std::thread::sleep  std::time::Duration
    std::vec::Vec::{new, push}
    std::result::Result
    std::marker::Copy  std::clone::Clone
    std::cmp::PartialEq
```

Create a **program** that plays [Conway's Game Of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life).

* The board must be represented using a `struct`, and each cell with an `enum`.

```rust
#[derive(Debug, PartialEq)]
enum ParseError {
    InvalidWidth { arg: &'static str },
    InvalidHeight { arg: &'static str },
    InvalidPercentage { arg: &'static str },
    TooManyArguments,
    NotEnoughArguments,
}

enum Cell {
    Dead,
    Alive,
}

impl Cell {
    fn is_alive(self) -> bool;
    fn is_dead(self) -> bool;
}

struct Board {
    width: usize,
    height: usize,
    cells: Vec<Cell>,
}

impl Board {
    fn new(width: usize, height: usize, percentage: u32) -> Self;
    fn from_args() -> Result<Self, ParseError>;
    fn step(&mut self);
    fn print(&self, clear: bool);
}
```

* `is_alive` must return whether the cell is alive or not.
* `is_dead` must return whether the cell is dead or not.
* `new` must generate a random board of size (`width` by `height`), with approximatly
`percentage`% live cells in it.
* `from_args` must parse the command-line arguments passed to the application and use them to
create a `Board` instance. Errors are communicated through the `ParseError` enumeration.
* `step` must simulate an entire step of the simulation. We will assume that the board repeats
itself infinitely in both directions. The cell at coordinate `width + 1` is the cell at coordinate
`1`. Similarly, the cell at coordinate `-1` is the cell at coordinate `width - 1`.
* `print` must print the board to the terminal. When `clear` is `true`, the function must also
clear a previously displayed board. Try not to clear the whole terminal! Just the board.

**Hint:** you might want to look at *ANSI Escape Codes* if you don't know where to start!

* Finally, write a **main** function that uses above function to simulate the game of life. At each
simulation step, the previous board must be replaced by the one in the terminal.

Example:

```txt
>_ cargo run -- 20 10 40
. . . . # . . . . . . . . . . # . . . .
. . . # # . . . . . . . . . . . . . . .
. . # . . # . . # . . . . . # . . . # .
. . . . . . . . . . . . . # . . . . . #
. . . # # . . # # # # . . . . . . . . .
# . . . # . . # . . # . . . . # . . # #
. . # # # # . # # . # # . . . . . . . .
. . # . . # . . # . . . . # # . . . . .
. # # . . # # . . . . . # # . . . . . .
. # . . . # . . . . # . . . . . # . . .
^C
>_ cargo run -- 
error: not enough arguments
>_ cargo run -- a b c
error: `a` is not a valid width
```

Keep in mind that this is only an example. You may use any characters and messages you wish.

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
