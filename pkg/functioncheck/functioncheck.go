package functioncheck

import (
	"fmt"
	"os"
)

func Execute(allowedFunctions []string) error {
	filePath := "allowed_functions/src/lib.rs"
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
