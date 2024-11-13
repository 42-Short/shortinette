# Overview of the learning structure
If you ask yourself why we decided to structure the Rust Short how we did it this is the right place!

We differentiate between mandatory and bonus on an
* `exercise`
  * Mandatory exercises (00-04) basic usage of a concept
  * Bonus exercises (05-07) encourage deeper understanding through either bigger scope or logically harder problems
* `module`
  * Mandatory modules (00-04) cover fundamental Rust language concepts
  * Bonus modules (06-07) explore more advanced, complex topics.

level. This structure is helping us find the right exercises for each area and communicates well to the participants what we would consider a `pass` or `bonus` in our system.

We mainly used the following sources to create the exercises:
* [The Rust Programming Language (The Book)](https://doc.rust-lang.org/book/)
* [Rustlings](https://github.com/rust-lang/rustlings)
* [Rust By Example](https://doc.rust-lang.org/stable/rust-by-example/)

Based on those sources we found these topics which we made into subjects:
### Module 00 - Basic Syntax
* ex00 - print
* ex01 - no return
* ex02 - loops
* ex03 - match
* ex04 - cargo

### Module 01 - Lifetimes & Slices
* ex00 - borrow
* ex01 - lifetimes
* ex02 - match with ranges
* ex03 - iter with lifetimes
* ex04 - slices

### Module 02 - Structs
* ex00 - drop/free
* ex01 - impl 
* ex02 - derived
* ex03 - Result + Option
* ex04 - NOT_YET_MADE

### Module 03 - Traits
* ex00 - annotation
* ex01 - function that needs trait
* ex02 - create Trait and implemnt it
* ex03 - Iter implmentation on colatz
* ex04 - Implement String parsing for own Struct

### Module 04 - Side Effects
* ex00 - TO_BE_DONE
* ex01 - file creation
* ex02 - directory size
* ex03 - process
* ex04 - processes

### Module 05 - Multithreading

### Module 06 - Unsafe


# Ideas
Here we put topics that maybe one day go into the Rust short
* Playing uni in 06 with channels
* String types and which one does what
* Vectors
* Complex Lifetime problems
* Chaining of map iter and collect with anonymus functions
* Benchmarking of stuff (iter vs. for loop)