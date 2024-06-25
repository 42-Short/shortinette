package functioncheck

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func writeDummyLib(allowedFunctions []string) error {
	libFilePath := "allowedfunctions/lib.rs"
	modFilePath := "allowedfunctions/ex00.rs"

	libContent := `#![no_std]
pub mod ex00;
`

	modContent := "pub mod ex00 {\n"

	for _, functionName := range allowedFunctions {
		if functionName == "println" {
			modContent += `
#[macro_export]
macro_rules! println {
    ($($arg:tt)*) => {{
        // Dummy Macro
    }}
}
`
		} else {
			modContent += fmt.Sprintf("pub fn %s() {}\n", functionName)
		}
	}

	modContent += "}\n"

	if err := os.WriteFile(libFilePath, []byte(libContent), 0644); err != nil {
		return fmt.Errorf("error writing dummy lib file: %s", err)
	}

	if err := os.WriteFile(modFilePath, []byte(modContent), 0644); err != nil {
		return fmt.Errorf("error writing mod file: %s", err)
	}

	fmt.Println("Dummy library created successfully")

	return nil
}

func compileDummyLib() error {
	cmd := exec.Command("rustc", "--crate-type=rlib", "allowedfunctions/lib.rs", "-o", "allowedfunctions/liballowedfunctions.rlib")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error compiling dummy lib: %s\noutput: %s", err, output)
	}
	fmt.Println("Dummy library compiled successfully")
	return nil
}

func prependExternCrate(code string) string {
	lines := strings.Split(code, "\n")
	lines = append([]string{"#[macro_use]", "extern crate allowedfunctions;", "use allowedfunctions::*;"}, lines...)
	return strings.Join(lines, "\n")
}

func compileWithDummyLib(studentProjectPath string) error {
	cmd := exec.Command("rustc", "-L", "allowedfunctions/", "--extern", "allowedfunctions=allowedfunctions/liballowedfunctions.rlib", studentProjectPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error compiling code: %s\noutput: %s", err, output)
	}
	fmt.Println("Student code compiled successfully")
	return nil
}

func Execute(allowedFunctions []string) error {
	var studentCode []byte
	var err error

	if err = writeDummyLib(allowedFunctions); err != nil {
		return err
	}

	if err = compileDummyLib(); err != nil {
		return err
	}

	if studentCode, err = os.ReadFile("ex00/ex00.rs"); err != nil {
		return fmt.Errorf("error reading student code: %s", err)
	}

	modifiedCode := prependExternCrate(string(studentCode))
	if err = os.WriteFile("ex00/tempstudentcode.rs", []byte(modifiedCode), 0644); err != nil {
		return fmt.Errorf("error writing modified student code: %s", err)
	}

	if err = compileWithDummyLib("ex00/tempstudentcode.rs"); err != nil {
		return err
	}

	fmt.Println("Execution completed successfully")
	return nil
}
