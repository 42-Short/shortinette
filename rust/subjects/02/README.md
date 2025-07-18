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

"Fascinating," said Farbold, nodding. The neophyte turned back to her work.

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
fn my_unused_function() {}x
```

* For exercises managed with cargo, the command `cargo clippy -- -D warnings` must run with no errors!

* You are _strongly_ encouraged to write extensive tests for the functions and programs you turn in. Tests can use the symbols you want, even if
they are not specified in the `allowed symbols` section. **However**, tests should not introduce **any additional external dependencies** beyond
those already required by the subject.

* All primitive types, i.e the ones you are able to use without importing them, are allowed.

* A type being allowed implies that its methods and attributes are allowed to be used as well, including the attributes of its implemented traits.

* You are **always** allowed to use `std::eprintln` for error handling.

* These rules may be overridden by specific exercises.

### Exercise 00: Bruh
```rust
// no allowed symbols

const allowed_dependencies = [""];
const turn_in_directory = "ex00/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Given this struct
```rust
struct ComplexStruct {
    name: String,
    optional_value: Option<Box<i32>>,
    values: Vec<i32>,
    some_other: Vec<u128>,
    metadata: std::collections::HashMap<String, Vec<u8>>,
    nested: Box<NestedStruct>,
}

struct NestedStruct {
    number: Box<i32>,
    optional_floats: Vec<Option<Box<f64>>>,
    data: std::collections::HashMap<String, Option<Box<i32>>>,
}
```

Implement the `free` method, such that the following code does *not* compile:
```rust
impl ComplexStruct {
    pub fn free(&mut self) {
        // Your implementation
    }
}

pub fn main() {
    let bruh: ComplexStruct = ComplexStruct {
        name: "hey".to_string(),
        optional_value: Some(Box::new(42)),
        values: vec![1377; 5],
        some_other: vec![137700000000; 5],
        metadata: std::collections::HashMap::new(),
        nested: Box::new(NestedStruct {
            number: Box::new(42),
            optional_floats: vec![Some(Box::new(42 as f64)), None, Some(Box::new(42 as f64 / 2.0))],
            data: std::collections::HashMap::new(),
        }),
    };

    bruh.free();
    
    println!("{}", bruh.name); 
}
```

Note that removing the call to `bruh.free()` should allow the code to compile successfully!

## Exercise 01: A Point In Space
```rust
// allowed symbols
const allowed_dependencies = [""];
const turn_in_directory = "ex01/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
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

## Exercise 02: Derive

```rust
// allowed symbols
use std::{
    clone::Clone,
    cmp::{PartialOrd, PartialEq},
    default::default,
    fmt::Debug,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex00/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Create a type, may it be a `struct` or an `enum`. You simply have to name it `MyType`.

You are **not** allowed to use the `impl` keyword!

```rust
#[cfg(test)]
mod tests{
    use super::*;

    #[test]
    fn test_my_type() {
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
}
```

Copy the above `test` function and make it compile.

## Exercise 03 Money money money

```rust
const allowed_dependencies = [""];
const turn_in_directory = "ex03/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Define the following types:

```rust
#[derive(PartialEq, Debug)]
enum BuyError {
    NotEnoughCoins,
    TooManyItems,
}
#[derive(PartialEq, Debug)]
enum SellError {
    TooManyCoins,
    NoItemToSell,
}

#[repr(u8)]
#[derive(Copy, Clone,Debug,PartialEq)]
enum Item {
    Sword = 10,        
    Shield = 15,       
    HealthPotion = 5,  
    UpgradeStone = 25, 
    Ring = 50,         
}

#[derive(PartialEq, Debug)]
struct Player {
    coins: u8,
    item: Option<Item>
}
```

Implement these functions
```rust
impl Player {
    pub fn buy(&mut self, item: Item) -> Result<(), BuyError>;
    pub fn sell(&mut self) -> Result<(), SellError>;
}
```
Ensure they operate as follows:

`buy` verifies that:
1. The player has enough coins
2. The player has enough room to store the item.

If either condition is unmet, return the appropriate error.

`sell` verifies that:
1. The player has an item to sell
2. The player can store the received coins without his pocket `overflowing`

If either condition is unmet, return the appropriate error.

If both errors apply simultaneously, return any of them.

The following test must compile and execute successfully:
```rust
#[cfg(test)]
mod tests{
    use super::*;

    #[test]
    fn test_player() {
        let mut player = Player { coins: 0, item: None};


        assert_eq!(player.buy(Item::HealthPotion), Err(BuyError::NotEnoughCoins));
        player.coins = 250;
        assert_eq!(player.buy(Item::Ring), Ok(()));
        assert_eq!(player.buy(Item::UpgradeStone), Err(BuyError::TooManyItems));

        player.coins = 242;
        assert_eq!(player.sell(), Err(SellError::TooManyCoins));
        player.coins = 0;
        assert_eq!(player.sell(), Ok(()));
        assert_eq!(player.sell(), Err(SellError::NoItemToSell));
    }
}
```

## Exercise 04: Swipe Left

```rust
// allowed symbols
use std::{
    iter::Iterator,
    default::Default,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex04/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Define the following struct:

```rust
struct Fibonacci {
    // Add whatever data you need here
}
```

Your `Fibonacci` struct must:
1. Be instantiable by calling `Fibonacci::default()`
2. Be iterable, always returning the next number in the Fibonacci sequence **up until (including) its 42nd element** ($F_{41}$ with 0-based indexing)

To be precise, the following test must pass:
```rust
#[test]
fn fib() {
    let len = Fibonacci::default()
    .take(1000000)
    .collect::<Vec<u32>>()
    .len();

    assert_eq!(len, 42)
}
```

The following code must compile and output the sequence's first 5 values:
```rust
fn fibonacci() {
    let fib = Fibonacci::default();

    for f in fib.take(5) {
        println!("{f}");
    }
}
```

Definition of the Fibonacci sequence:

$F_0 = 0, F_1 = 1, F_n = (F_{n - 1} + F_{n - 2})$


## Exercise 05: Lexical Analysis

```rust
// allowed symbols
use ftkit::ARGS;
use std::{
    fmt::Debug,
    println,
};

const allowed_dependencies = ["ftkit"];
const allowed_dependencies = [""];
const turn_in_directory = "ex05/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
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
                println!("{:?}", token);
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

```rust
// allowed symbols
use std::{
    collections::HashMap,
    fmt::Debug,
};

const allowed_dependencies = [""];
const turn_in_directory = "ex05/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
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
* `get_item`: Returns a reference to an item by name, or `None` if it doesn't exist.
* `list_items`: Returns a vector of references to all items in the inventory.
* `total_value`: Calculates and returns the total value of all items in the inventory (price * quantity for each item).

Implement a `Discountable` trait for `Item`:
```rust
trait Discountable {
    fn apply_discount(&mut self, percentage: f32);
}
```
Implement this trait for `Item`, such that it reduces the price of the item by the given percentage. Invalid percentage values (`< 0 || > 100`) are considered **undefined behavior**. You are free to handle them as you please.

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

```rust
// allowed symbols
use std::{
    println,
    print,
    thread::sleep,
    time::Duration,
    vec::Vec,
    marker::Copy,
    clone::Clone,
    cmp::PartialEq,
};

const allowed_dependencies = ["ftkit"];
const allowed_dependencies = [""];
const turn_in_directory = "ex07/";
const files_to_turn_in = ["src/main.rs", "Cargo.toml"];
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
* `new` must generate a random board of size (`width` by `height`), with approximately
`percentage`% live cells in it.
* `from_args` must parse the command-line arguments passed to the application and use them to
create a `Board` instance. Errors are communicated through the `ParseError` enumeration.
* `step` must simulate an entire step of the simulation. We will assume that the board repeats
itself infinitely in both directions. The cell at coordinate `width + 1` is the cell at coordinate
`1`. Similarly, the cell at coordinate `-1` is the cell at coordinate `width - 1`.
* `print` must print the board to the terminal. When `clear` is `true`, the function must also
clear a previously displayed board. Try not to clear the whole terminal! Just the board.

**Hint:** you might want to look at *ANSI Escape Codes* if you don't know where to start!

* Finally, write a **main** function that uses above function to simulate the game of life. At each simulation step, the previous board must be replaced by the one in the terminal.

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

---
**License Notice:**
This file contains content licensed under two different terms:
- The MIT License applies to the original content (see `LICENSES/MIT-rust-subjects.txt`).
- The Creative Commons Attribution-ShareAlike 4.0 International (CC BY-SA 4.0 applies to any modifications or additions (see `LICENSES/CC-BY-SA-4.0.txt`).

When distributing modified versions, you must comply with both the MIT License and the CC BY-SA 4.0.
For complete details, refer to the main licensing file of Shortinette.
---
