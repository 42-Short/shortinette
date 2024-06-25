package functioncheck

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func writeDummyLib(allowedFunctions []string) error {
	filePath := "allowedfunctions/src/lib.rs"
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating %s: %s", filePath, err)
	}
	defer file.Close()

	_, err = file.WriteString(("#[no_std]\n\n"))
	if err != nil {
		return fmt.Errorf("error writing to %s: %s", filePath, err)
	}

	for _, functionName := range allowedFunctions {
		if functionName == "println" {

			_, err = file.WriteString(`
		#[macro_export]
		macro_rules! println {
		($($arg:tt)*) => {{
				// Dummy Macro
			}}
		}
				`)
			if err != nil {
				return fmt.Errorf("error writing to %s: %s", filePath, err)
			}
		}
	}
	return nil
}

func compileWithDummyLib(studentProjectPath string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '[dependencies]\nallowed_functions = { path = \"../allowed_functions\" }\n' >> %s/Cargo.toml", studentProjectPath))
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("cargo", "build", "--manifest-path", fmt.Sprintf("%s/Cargo.toml", studentProjectPath))
	if err := cmd.Run(); err != nil {
		return err
	}

	log.Println("Compilation successful")
	return nil
}

func Execute(allowedFunctions []string) error {
	if err := writeDummyLib(allowedFunctions); err != nil {
		return err
	}

	if err := compileWithDummyLib("ex00"); err != nil {
		return err
	}

	return nil
}
