# Module 06: Feeling Unsafe

## Foreword

For the eighth day in a row, Neophyte Col found himself standing before the Two Great Guards of the
temple. They stood before the large entrance of the temple, clad in simple robes. Nevertheless, they
were imposing, and feared. He strode toward the first guard, confident, and handed the parchment
bearing his program.

The First Guard read through it carefully. This step was but a formality; yesterday, he had only
failed to gain the assent of the Second Guard. He was certain he had addressed all outstanding
complaints.

The First Guard handed the parchment back to Col. Then, in a blinding motion, slapped Col across the
face with his bare hand. In a measured voice, the First Guard spoke to him: "mismatched types:
expected &Target, found &<T as Deref>::Target", then fell silent.

Col took his parchment and retreated back to a nearby bench, close to tears. Eight days. It was not
as though his program was particularly complicated, and yet he could not convince the Two Great
Guards to permit him entry to the temple. He had not had this much difficulty at other temples!

At another bench, he saw a fellow neophyte. They had spoken two days previous, when he had learned
his compatriot had been toiling outside the temple for close to two weeks to get his program
accepted.

It was the guard's fault. Col knew his program would work as intended. All they ever seemed to do
was pick on minor errors, denying him for the most petty of reasons. They were looking for reasons
to deny him entry; that was it!

He was beginning to seethe with resentment.

Thereupon, he noticed a monk speaking to the other neophyte. The conversation seemed quite animated
when, all of a sudden, the neophyte whooped, leapt up, and rushed toward the temple. As he ran, he
seemed to be frantically modifying his program.

However, rather than face the Two Great Guards, he instead walked over to a small, dingy part of the
wall. To Col's surprise, the wall opened into what appeared to be a secret entrance. The neophyte
passed through, and was gone from sight.

Col sat stunned. A secret entrance? Then... then the Two Great Guards might merely be a prank!
Something the other monks put neophyte through to teach them humility. Or... resilience. Or... or...
perhaps it was just to laugh at them secretly.

"Do you wish to know what I said to him?" a voice asked. Col turned to see the monk standing beside
his bench. "You told him there was another way in, did you not?"

"Yes," she replied. "I told him of the unsafe door."

"unsafe?" Col asked.

"Indeed. It is a secret known to those who have studied long and hard at the temple. In truth, one
can overcome many of the obstacles posed in writing one's program through the use of the unsafe
arts, as spoken of in the [Rustonomicon](http://doc.rust-lang.org/nightly/nomicon/)."

"Are they powerful?" Col asked in wonder.

"Immensely powerful. With transmute, one can simply re-assign the type of a value, or extend the
lifetime of a pointer. One can even summon pointers from the air itself, or data from nothingness."

Col felt he finally understood how the temple worked. It was this "unsafe" magic that the monks
used! However...

"Then, why are the Two Great Guards employed if one can simply walk through the unsafe door to reach
the temple? Why not..."

At that moment, a blood-curdling scream was heard from within the temple. It echoed across the
courtyard before ending suddenly.

Silence descended. No one moved. No one spoke. The wind stilled. Even the birds halted their
singing.

Col could feel his heart pounding in his ears.

"They are there," the monk said, breaking the spell, "to protect you from the temple, and what lies
within."

Col turned to gaze once more at the hidden doorway. "Then why does that door exist?"

"Because, even they are not infallible. Some times, one must face peril alone." She sighed. "But not
all are so brave or skilled."

With that, the monk took a sheet of parchment from her cloak and walked toward the Two Great Guards.
The First Guard read her program and nodded. The Second Guard read her program, handed it back, and
then struck her across the face.

"Borrowed value does not live long enough."

The monk, rubbing her face, walked back and sat down on one of the benches, muttering curses.

_The [first Rust Koan](https://users.rust-lang.org/t/rust-koans/2408)._

## General Rules
* You **must not** have a `main` present if not specifically requested.

* Any exercise managed by cargo you turn in must compile _without warnings_ using the `cargo test` command. If not managed by cargo, it must compile _without warnings_ with the `rustc` compiler available on the school's
machines without additional options.

* Only dependencies specified in the allowed dependencies section are allowed.

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

## Module Rules

In this module, you will take your first steps in writing dangerous code safely. 
If you keep going until the very end, you will learn stuff no sane person should ever have to
worry about, like cross-compiling languages with **F**oreign **F**unction **I**nterfaces (FFI), or the `#![no_main]` attribute. 

Attentive readers will have noticed a missing line in the general rules of this subject. You are _allowed_ to use `unsafe` code in this module! However, some rules must be followed.

1. You must use the `#![forbid(unsafe_op_in_unsafe_fn)]` global attribute.

2. When an `unsafe fn` is defined, its documentation must contain a `# Safety` section
   describing how to use it correctly.

```rust
/// Returns one of the elements of `slice`, as specified by
/// `index`
///
/// # Safety
///
/// The provided `index` must be in bounds (i.e. it must be
/// **strictly** less than `slice.len()`).
unsafe fn get_unchecked(slice: &[u32], index: usize) -> u32 {
    // SAFETY:
    //  - We have been given a regular `&[u32]` slice, which
    //    ensures that the pointer is valid for reads and is
    //    properly aligned. We can turn it back into a regular
    //    reference.
    //  - The responsibility to ensure that the `index` is in bounds,
    //    is on the caller.
    unsafe { *slice.as_ptr().add(index) }
}
```

3. When an `unsafe trait` is defined, its documentation must contain a `# Safety` section
   describing how to implement it correctly.

```rust
/// Types that can be initialized to zeros.
///
/// # Safety
///
/// Implementers of this trait must allow the "all-zeros" bit pattern.
unsafe trait Zeroable {
    fn zeroed() -> Self {
        // SAFETY:
        //  Implementers of the `Zeroable` trait can be initialized
        //  with the "all-zeros" bit pattern, ensuring that calling
        //  this function won't produce an invalid value.
        unsafe { std::mem::zeroed() }
    }
}
```

4. Every time an `unsafe` block is used, it must be annotated with a `SAFETY:` directive, explaining the 
thinking process behind this code.

```rust
let slice: &[u32] = &[1, 2, 3];
// SAFETY:
//  We know that `slice` has a length of 3, ensuring that accessing
//  the element at index 2 is always valid.
let val = unsafe { get_unchecked(slice, 2) };
```

5. Every time an `unsafe impl` is declared, it must be annotated with a `SAFETY:` directive.

```rust
// SAFETY:
//  The `u64` type allows the "all-zeros" bit pattern - it corresponds
//  to the value `0u64`.
unsafe impl Zeroable for u64 {}
```

To summarize:

- `unsafe fn` means **_'know what you are doing before calling this function'_**.
- `unsafe trait` means **_'know what you are doing before implementing this trait'_**.
- `unsafe {}` and `unsafe impl` both mean **_'I know what I am doing'_**.

## Exercise 00: Libft

```rust
// allowed symbols
use std::ptr::{write, read, add};

const turn_in_directory = "ex00/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Let's start simple.

```rust
pub fn ft_swap<T>(a: &mut T, b: &mut T);
pub unsafe fn ft_strlen(s: *const u8) -> usize;
pub unsafe fn ft_strcpy(dst: *mut u8, src: *const u8);
```

- `ft_swap` must swap any two values of any type. Maybe `T` can be copied, maybe not. Maybe it has a default value. Maybe not.
- `ft_strlen` must count the number of non-null bytes, starting at `s`. You must write an
  appropriate `# Safety` section in the documentation of that function to educate about its users
  about its correct usage.
- `ft_strcpy` must copy the null-terminated string at `src` into `dst`. Just like `ft_strlen`, you
  must _precisely_ describe the requirements of your function within a `# Safety` section in its
  documentation.

Example:

```rust
let mut a = String::from("Hello, World!");
let mut b = String::from("Goodbye, World!");
ft_swap(&mut a, &mut b);
assert_eq!(a, "Goodbye, World!");
assert_eq!(b, "Hello, World!");

let s = b"Hello, World!\0";
// # Safety
//  /* ... */
let len = unsafe { ft_strlen(s.as_ptr()) };
assert_eq!(len, 13);

let mut dst = [0u8; 14];
// # Safety
//  /* ... */
unsafe { ft_strcpy(dst.as_mut_ptr(), s.as_ptr()) };
assert_eq!(&dst, b"Hello, World!\0");
```

## Exercise 01: Philosopher's Stone

```rust
// allowed symbols
use std::{slice::from_raw_parts, mem::transmute};

const turn_in_directory = "ex01/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

```rust
pub type GoldNugget = u16;

pub type Iron = u32;
pub type Mercure = u64;

pub struct PhilosopherStone;

impl PhilosopherStone {
    fn transmute_iron(self, iron: Iron) -> [GoldNugget; 2];
    fn transmute_mercure(self, mercure: Mercure) -> [GoldNugget; 4];
}
```

- The `transmute_iron` function must convert the given `Iron` into a bunch of `GoldNugget`s. The
  bit-pattern of the original iron _must be preserved_, ignoring byte-order.
- The `transmute_mercure` function must convert the given `Mercure` into a bunch of `GoldNugget`s.
  The bit-pattern of the original mercure _must be preserved_, ignoring byte-order.

Example:

```rust
// On a LITTLE-ENDIAN machine!
let iron = 0x12345678;
assert_eq!(PhilosopherStone.transmute_iron(iron), [0x5678, 0x1234]);
let mercure = 0x0123456789ABCDEF;
assert_eq!(
    PhilosopherStone.transmute_mercure(mercure),
    [0xCDEF, 0x89AB, 0x4567, 0x0123],
);
```

Let's generalize a bit.

```rust
pub type Gold = [GoldNugget];

unsafe trait Metal {}
```

- A `Metal` is a type that may be turned into gold by the `PhilosopherStone`.
- Do not forget the `# Safety` comment in the documentation for `Metal`!

```rust
impl PhilosopherStone {
    fn transmute_metal<M: Metal>(self, metal: &M) -> &Gold;
}
```

- The `transmute_metal` function must convert the given `metal` into `&Gold`.

Example:

```rust
let mercure: Mercure = 0x0123456789ABCDEF;
assert_eq!(
    PhilosopherStone.transmute_metal(&mercure),
    &[0xCDEF, 0x89AB, 0x4567, 0x0123],
);
```

## Exercise 02: Carton
```rust
// allowed symbols
use std::{
    ptr::NonNull
    ops::{Deref, DerefMut},
    alloc::{alloc, dealloc, handle_alloc_error, Layout},
};

const turn_in_directory = "ex02/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Create a type named `Carton<T>`, which must manage a heap allocation of a single `T`.

```rust
pub struct Carton<T> {
    data: NonNull<T>,
}

impl<T> Carton<T> {
    fn new(value: T) -> Self;
    fn into_inner(self) -> T;
}
```

- You must make sure that `Carton<T>` has the correct _variance_ over `T`.
- You must make sure that the _drop checker_ makes the correct assumptions about the lifetime of
  a `T` owned by a `Carton<T>`.
- You must make sure that `Carton<T>` properly manages the memory it owns. Allocated memory must be freed later! [`cargo-valgrind`](https://crates.io/crates/cargo-valgrind) can help you track unfreed memory.

Example:

```rust
#[derive(Clone)]
struct Point { x: u32, y: u32 }

let point_in_carton = Carton::new(Point { x: 1, y: 2 });
assert_eq!(point_in_carton.x, 1);
assert_eq!(point_in_carton.y, 2);

let mut another_point = point_in_carton.clone();
another_point.x = 2;
another_point.y = 3;
assert_eq!(another_point.x, 2);
assert_eq!(another_point.y, 3);
assert_eq!(point_in_carton.x, 1);
assert_eq!(point_in_carton.y, 2);
```

## Exercise 03: `Cellule<T>`

```rust
// allowed symbols
use std::{
    clone::Clone,
    marker::Copy,
    cell::UnsafeCell,
    ptr::*,
    mem::*,
};

const turn_in_directory = "ex03/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Let's re-create our own `Cell<T>` named `Cellule<T>`.

You must implement the following inherent methods, as specified in the official documentation for
[`Cell<T>`](https://doc.rust-lang.org/std/cell/#cellt).

```rust
pub struct Cellule<T> {
    cell: UnsafeCell<T>,
}

impl<T> Cellule<T> {
    pub const fn new(value: T) -> Self;

    fn set(&self, value: T);
    fn replace(&self, value: T) -> T;

    fn get(&self) -> T;
    fn get_mut(&mut self) -> &mut T;

    fn into_inner(self) -> T;
}
```

Note that you may need to add trait bounds to some of the above methods to ensure their safety,
and once again, be extra careful of the _variance_ of your type.

## Exercise 04: RAII

```rust
// allowed symbols
use std::ptr::{
    copy::Copy,
    clone::Clone,
    cmp::{PartialEq, Eq, PartialOrd, Ord},
    fmt::{Debug, Display},
    mem::forget,
};
use libc::{__errno_location, c_int, strerror, write, read, open, close};
use cstr::cstr;

const allowed_dependencies = ["libc", "cstr"];
const turn_in_directory = "ex04/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

Create an `Errno` type.

```rust
pub struct Errno(libc::c_int);

impl Errno {
    pub fn last() -> Self;
    fn make_last(self);
    fn description(self) -> &'static str;
}
```

- `last` must return the **calling thread**'s last `errno`.
- `make_last` must make an `Errno` the calling thread's last `errno`.
- `description` must return a description of the error. Don't try to enumerate _every_
  possible error! I'm sure a function from `libc` can do it for you.

Example:

```rust
Errno(12).make_last();
assert_eq!(Errno::last(), Errno(12));

let desc = format!("{}", Errno(1));
assert_eq!(desc, "Operation not permitted");
```

With a robust way to handle errors, we can now start for real:

```rust
pub struct Fd(libc::c_int);

impl Fd {
    const STDIN: Self = /* ... */;
    const STDOUT: Self = /* ... */;
    const STDERR: Self = /* ... */;

    pub fn open(file: &CStr) -> Result<Self, Errno>;
    pub fn create(file: &CStr) -> Result<Self, Errno>;
    fn write(self, data: &[u8]) -> Result<usize, Errno>;
    fn read(self, buffer: &mut [u8]) -> Result<usize, Errno>;
    fn close(self) -> Result<(), Errno>;
}
```

- `open` must open a new file descriptor for reading (only).
- `create` must open a new file descriptor for writing (only). If the file already exists, it must
  be truncated.
- `write` must write the data referenced to by `data` to the file descriptor.
- `read` must read data from the file descriptor into `buffer`.
- `close` must attempt to close the file descriptor.
- In any case, errors must be handled properly.

That's cool, and all, but let's add one more layer!

```rust
pub struct File(Fd);

impl File {
    pub fn open(file: &CStr) -> Result<Self, Errno>;
    pub fn create(file: &CStr) -> Result<Self, Errno>;
    fn write(&self, data: &[u8]) -> Result<usize, Errno>;
    fn read(&self, buffer: &mut [u8]) -> Result<usize, Errno>;
    fn leak(self) -> Fd;
}
```

- `open` and `create` work exactly the same as `Fd::open` and `Fd::create`.
- `write` and `read` work the same way as `Fd::write` and `Fd::read`. Note, however, that they only
  _borrow_ the `File`.
- `leak` must "leak" the file descriptor of the file; returning it and "forgetting" that it had to
  be closed layer.

When a `File` is dropped, it must automatically close its file descriptor.

## Exercise 05: Tableau

```rust
// allowed symbols
use std::{
    alloc::{alloc, dealloc, Layout},
    marker::Copy,
    clone::Clone,
    ops::{Deref, DerefMut},
    ptr::*,
    mem::*,
};

const turn_in_directory = "ex05/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml"];
```

It must implement the following inherent methods, as specified in the official documentation for [`Vec`](https://doc.rust-lang.org/std/vec/struct.Vec.html):

```rust
pub struct Tableau<T> {
    // Up to you!
}

impl<T> Tableau<T> {
    pub const fn new() -> Self;

    fn len(&self) -> usize;
    fn is_empty(&self) -> bool;

    fn push(&mut self, item: T);
    fn pop(&mut self) -> Option<T>;

    fn clear(&mut self);
}
```

It must be possible to do the following:

```rust
let mut tab0 = Tableau::new();
tab0.push(1); tab0.push(2); tab0.push(4);
let tab1 = tab0.clone();

for it in tab1 {
    println!("{it}");
}
// This will print:
// 1
// 2
// 4

let c: &[i32] = &*tab0;
assert_eq!(c, [1, 2, 4]);
```

**Note:** Be careful! _Careful_! Any code that you didn't write can panic. Cloning a value can
panic. Dropping a value can panic. Make sure that your type do _not_ leak memory; even when cloning
or dropping.

If you're feeling like taking a challenge, you can try to write a macro to construct a `Tableau<T>`
automatically:

```rust
let tab: Tableau<i32> = tableau![1, 2, 4];
assert_eq!(tab, [1, 2, 4]);
```

## Exercise 06: Foreign User

```rust
// allowed symbols
use std::{
    mem::MaybeUninit,
    ffi::{CStr, c_int, c_char},
};

const turn_in_directory = "ex06/";
const files_to_turn_in = ["src/lib.rs", "Cargo.toml", "build.rs", "awesome.c"];
```

However sad it may be, Rust is not the only programming language in existence.

Let's create a simple C library.

```c
typedef unsigned int t_id;

typedef struct {
    t_id id;
    char const *name;
} t_user;

typedef struct {
    t_id next_user_id;
    t_user *users;
    size_t count;
    size_t allocated;
} t_database;

typedef enum {
    ERR_SUCCESS,
    ERR_MEMORY,
    ERR_NO_MORE_IDS,
    ERR_UNKNOWN_ID,
} e_result;

e_result create_database(t_database *database);
void delete_database(t_database *database);

e_result create_user(t_database *database, char const *name, t_id *result);
e_result delete_user(t_database *database, t_id id);
e_result get_user(t_database const *database, t_id id, t_user const *result);
```

- `create_database` must initialize the `t_database` instance.
- `delete_database` must destroy the `t_database` instance, freeing any allocated memory.
- `create_user` must insert a new `t_user` instance into the database.
- `delete_user` must remove a `t_user` from the database.
- `get_user` must store a pointer to the user with the provided ID into `result`.

In any case, on success, `ERR_SUCCESS` is returned. When a memory error occurs, `ERR_MEMORY` is
returned. When no more IDs can be allocated, `ERR_NO_MORE_IDS` is returned. When a given ID is
invalid, `ERR_UNKNOWN_ID` is returned.

You now have an awesome C library, but you unfortunately cannot use it in Rust...

Set up your project, such that this C library is automatically compiled into a `.a` static library
when you call `cargo build`. Your Rust library must link against that compiled C library.

```rust
pub enum Error { /* ... */ }

pub type Id = /* ... */;

pub struct User { /* ... */ }

pub struct Database { /* ... */ }

impl Database {
    pub fn new() -> Self;

    fn create_user(&mut self, name: &CStr) -> Result<Id, Error>;
    fn delete_user(&mut self, id: Id) -> Result<(), Error>;
    fn get_user(&self, id: Id) -> Result<&User, Error>;
}
```

- `new` must call the `create_database` function of your C library.
- `create_user` must call `create_user`.
- `delete_user` must call `delete_user`.
- `get_user` must call `get_user`.

When a `Database` goes out of scope, it must automatically call `delete_database`.

## Exercise 07: Bare Metal

```rust
// allowed symbols
use core::arch::asm;

const turn_in_directory = "ex07/";
const files_to_turn_in = ["ft_putchar.rs"];
```

What better way to finish this great journey than by writing your very first `C` function?

```rust
fn ft_putchar(c: u8);
fn ft_exit(code: u8) -> !;
```

- You must use the `#![no_std]` and `#![no_main]` global attributes.
- If you want to make sure nothing gets into your sacred namespace, you can optionally add the
  `#![no_implicit_prelude]` global attribute.
- `ft_putchar` must print the specified character to the standard output of the program.
- `ft_exit` must exit the process with the specified exit code.
- Create a **program** that calls `ft_putchar` three times. Once with `'4'`, once with `'2'`, and once with
  `'\n'`.
- Up to you to find out how to compile this! (Hint: it's not cargo).

Example:

```txt
>_ ./ft_putchar
42
>_ echo $?
42
```

---
**License Notice:**
This file contains content licensed under two different terms:
- The MIT License applies to the original content (see `LICENSES/MIT-rust-subjects.txt`).
- The Creative Commons Attribution-ShareAlike 4.0 International (CC BY-SA 4.0 applies to any modifications or additions (see `LICENSES/CC-BY-SA-4.0.txt`).

When distributing modified versions, you must comply with both the MIT License and the CC BY-SA 4.0.
For complete details, refer to the main licensing file of Shortinette.
---
