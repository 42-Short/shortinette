package functioncheck

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type AllowedItem struct {
	Name   string
	Type   string
	Parent string
}

func parseCSV(filePath string) ([]AllowedItem, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file: %s", err)
	}

	var allItems []AllowedItem
	for _, record := range records[1:] { // Skip header
		allItems = append(allItems, AllowedItem{
			Name:   record[0],
			Type:   record[1],
			Parent: record[2],
		})
	}
	return allItems, nil
}

func getAllPossibleItems(allItems []AllowedItem) map[string]AllowedItem {
	// Convert the list of all items to a map for quick lookup
	itemMap := make(map[string]AllowedItem)
	for _, item := range allItems {
		itemMap[item.Name] = item
	}
	return itemMap
}

func getForbiddenItems(allItems map[string]AllowedItem, allowedSet map[string]bool) []AllowedItem {
	forbiddenItems := []AllowedItem{}
	for name, item := range allItems {
		if !allowedSet[name] {
			forbiddenItems = append(forbiddenItems, item)
		}
	}
	return forbiddenItems
}

func writeDummyLib(allowedItems, forbiddenItems []AllowedItem) error {
	libFilePath := "allowedfunctions/lib.rs"
	modFilePath := "allowedfunctions/ex00.rs"

	libContent := `#![no_std]
#![recursion_limit = "512"]
pub mod ex00;
`

	modContent := "pub mod ex00 {\n"

	definedItems := make(map[string]bool)

	for _, item := range allowedItems {
		if definedItems[item.Name] {
			continue
		}
		definedItems[item.Name] = true

		modContent += fmt.Sprintf("#[cfg(feature = \"allowed_%s\")]\n", item.Name)

		switch item.Type {
		case "macro":
			modContent += fmt.Sprintf(`
#[macro_export]
macro_rules! %s {
    ($($arg:tt)*) => {{
        // Dummy Macro
    }}
}
`, item.Name)
		case "function":
			modContent += fmt.Sprintf("pub fn %s() {}\n", item.Name)
		case "module":
			modContent += fmt.Sprintf("pub mod %s {}\n", item.Name)
		}
	}

	for _, item := range forbiddenItems {
		if definedItems[item.Name] {
			continue
		}
		definedItems[item.Name] = true

		// Define forbidden items without causing compile errors in the dummy lib
		modContent += fmt.Sprintf("#[cfg(not(feature = \"allowed_%s\"))]\n", item.Name)

		switch item.Type {
		case "macro":
			modContent += fmt.Sprintf(`
#[macro_export]
macro_rules! %s {
    ($($arg:tt)*) => {{
        // Forbidden macro definition placeholder
    }}
}
`, item.Name)
		case "function":
			modContent += fmt.Sprintf(`
#[cfg(not(feature = "allowed_%s"))]
pub fn %s() {
    // Forbidden function definition placeholder
}
`, item.Name, item.Name)
		case "module":
			modContent += fmt.Sprintf(`
#[cfg(not(feature = "allowed_%s"))]
pub mod %s {
    // Forbidden module definition placeholder
}
`, item.Name, item.Name)
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
	lines = append([]string{"#[macro_use]", "extern crate allowedfunctions;", "use allowedfunctions::ex00::*;"}, lines...)
	return strings.Join(lines, "\n")
}

func compileWithDummyLib(studentProjectPath string, allowedList []string) error {
	featureFlags := []string{}
	for _, name := range allowedList {
		featureFlags = append(featureFlags, fmt.Sprintf("--cfg=feature=\"allowed_%s\"", name))
	}

	cmdArgs := append([]string{"-L", "allowedfunctions/", "--extern", "allowedfunctions=allowedfunctions/liballowedfunctions.rlib"}, featureFlags...)
	cmdArgs = append(cmdArgs, studentProjectPath)
	cmd := exec.Command("rustc", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error compiling code: %s\noutput: %s", err, output)
	}
	fmt.Println("Student code compiled successfully")
	return nil
}

func Execute(csvFilePath string, allowedList []string) error {
	allItems, err := parseCSV(csvFilePath)
	if err != nil {
		return err
	}

	allItemsMap := getAllPossibleItems(allItems)
	allowedSet := make(map[string]bool)
	for _, name := range allowedList {
		allowedSet[name] = true
	}

	allowedItems := []AllowedItem{}
	for _, name := range allowedList {
		if item, exists := allItemsMap[name]; exists {
			allowedItems = append(allowedItems, item)
		}
	}

	forbiddenItems := getForbiddenItems(allItemsMap, allowedSet)

	if err = writeDummyLib(allowedItems, forbiddenItems); err != nil {
		return err
	}

	if err = compileDummyLib(); err != nil {
		return err
	}

	studentCode, err := os.ReadFile("ex00/ex00.rs")
	if err != nil {
		return fmt.Errorf("error reading student code: %s", err)
	}

	modifiedCode := prependExternCrate(string(studentCode))
	if err = os.WriteFile("ex00/tempstudentcode.rs", []byte(modifiedCode), 0644); err != nil {
		return fmt.Errorf("error writing modified student code: %s", err)
	}

	if err = compileWithDummyLib("ex00/tempstudentcode.rs", allowedList); err != nil {
		return err
	}

	fmt.Println("Execution completed successfully")
	return nil
}
