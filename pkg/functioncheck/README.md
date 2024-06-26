# Documentation for `functioncheck` Package
## Overview
This package is designed to validate that no forbidden functions are used.

## Structure
* **alloweditems.go**: Defines data structures and functions for setting up the Rust library containing allowed items.
* **functioncheck.go**: Contains the main logic for initializing the compilation environment, adding necessary compiler hints to the student code, and compiling it.
* **templates.go**: Contains templates for generating code and configuration files.

## Key Components
### Data Structures
```go
type AllowedItem struct {
	Name string
	Type string // "macro" or "function"
}
```
The `AllowedItem` struct represents an item (function or macro) that is allowed in student code.

### Functions
#### alloweditems.go
* **writeTemplateToFile(template, itemName string, file \*os.File) error**: Writes a formatted template to a file.
* **writeAllowedItemsLib(allowedItems []AllowedItem, file \*os.File) error**: Writes the allowed items to a Rust library file.
* **createFileWithDirs(filePath string) (*os.File, error)**: Creates a file along with necessary directories.
* **writeCargoToml(filePath, content string) error**: Writes content to a Cargo.toml file.

#### functioncheck.go
* **initCompilingEnvironment(allowedItems []AllowedItem, exercise string) error**: Sets up the compilation environment.
* **prependHeadersToStudentCode(filePath, exercise string) error**: Adds headers to the student's code.
* **compileWithDummyLib(sourceDir string) (string, error)**: Compiles the code with a dummy library and returns the output.
* **setToSlice(forbiddenFunctionSet map[string]bool) []string**: Converts a set of forbidden functions to a slice.
* **parseForbiddenFunctions(compilerOutput string) ([]string, error)**: Parses compiler output to extract forbidden functions.
* **handleCompileError(output string) error**: Handles compilation errors by parsing forbidden functions.
* **Execute(allowedItems []AllowedItem, exercise string) (err error)**: The main function to execute the check process.

## Usage
1. **Initialize Compilation Environment**
	* The **initCompilingEnvironment** function creates the required Rust files and directories, writing the allowed items into the Rust library.
2. **Modify Student Code**
	* The **prependHeadersToStudentCode** function reads the student's code, adds necessary headers, and writes it to a temporary file.
3. **Compile Code**
	* The **compileWithDummyLib** function compiles the student's code using the dummy library that contains the allowed items.
4. **Parse and Handle Errors**
	* If the compilation fails, the **handleCompileError** function parses the output to identify any forbidden functions used.
5. **Execute the Process**
	* The **Execute** function orchestrates the entire process: initializing the environment, modifying the code, compiling it, and handling any errors.

## Templates
* **allowedMacroTemplate**: Template for allowed macros.
* **allowedFunctionTemplate**: Template for allowed functions.
* **allowedItemsLibHeader**: Header for the allowed items library.
* **cargoTomlTemplate**: Template for the Cargo.toml file.
* **allowedItemsCargoToml**: Cargo.toml content for the allowed items library.
* **studentCodePrefix**: Prefix headers for the student's code.
