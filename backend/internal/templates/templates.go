package templates

const (
	AllowedMacroTemplate = `
#[cfg(not(feature = "allowed_%s"))]
#[macro_export]
macro_rules! %s {
	($($arg:tt)*) => {{}}
}
`
	AllowedFunctionTemplate = `
#[cfg(not(feature = "allowed_%s"))]
pub fn %s() {}
`
	AllowedItemsLibHeader = "pub mod %s { "
	CargoTomlTemplate     = `[package]
name = "%s"
version = "0.1.0"
edition = "2021"

[dependencies]
allowedfunctions = { path = "allowedfunctions" }

[[bin]]
name = "compile-environment"
path = "src/%s/temp.rs"

[workspace]
`
	AllowedItemsCargoToml = `[package]
name = "allowedfunctions"
version = "0.1.0"
edition = "2021"
`
	StudentCodePrefix = `#![no_std]
#[macro_use]
extern crate allowedfunctions;
use allowedfunctions::%s::*;
`
	DummyMain = `
fn main() {
	%s;
}`
)
