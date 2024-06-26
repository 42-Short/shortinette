package functioncheck

const (
	allowedMacroTemplate = `
#[cfg(not(feature = "allowed_%s"))]
#[macro_export]
macro_rules! %s {
	($($arg:tt)*) => {{}}
}
`
	allowedFunctionTemplate = `
#[cfg(not(feature = "allowed_%s"))]
pub fn %s() {}
`
	allowedItemsLibHeader = "pub mod ex%s { "
	cargoTomlTemplate     = `[package]
name = "%s"
version = "0.1.0"
edition = "2021"

[dependencies]
allowedfunctions = { path = "allowedfunctions" }

[[bin]]
name = "internal"
path = "src/%s/temp.rs"

[workspace]
`
	allowedItemsCargoToml = `[package]
name = "allowedfunctions"
version = "0.1.0"
edition = "2021"
`
	studentCodePrefix = `#![no_std]
#[macro_use]
extern crate allowedfunctions;
use allowedfunctions::ex%s::*;
`
)
