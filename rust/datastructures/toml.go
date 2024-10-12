package toml

import (
	"fmt"
	"os"
	"path/filepath"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/pelletier/go-toml/v2"
)

type Profile struct {
	Inherits       string `toml:"inherits"`
	Strip          bool   `toml:"strip"`
	OverflowChecks bool   `toml:"overflow-checks"`
}

type Bin struct {
	Name string `toml:"name"`
	Path string `toml:"path"`
}

type Package struct {
	Name        string   `toml:"name"`
	Version     string   `toml:"version"`
	Edition     string   `toml:"edition"`
	Authors     []string `toml:"authors"`
	Description string   `toml:"description"`
	DefaultRun  string   `toml:"default-run"`
}

type Toml struct {
	Package      Package            `toml:"package"`
	Dependencies map[string]string  `toml:"dependencies"`
	Bins         []Bin              `toml:"bin"`
	Profiles     map[string]Profile `toml:"profile"`
}

func CheckCargoTomlContent(exercise Exercise.Exercise, expectedContent map[string]string) Exercise.Result {
	tomlPath := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory, "Cargo.toml")
	fieldMap, err := ReadToml(tomlPath)
	if err != nil {
		logger.Error.Printf("internal error: %s", err)
		return Exercise.Result{Passed: false, Output: "internal error"}
	}
	var result = Exercise.Result{Passed: true, Output: "OK"}
	for key, expectedValue := range expectedContent {
		value, ok := fieldMap[key]
		if !ok {
			result.Passed = false
			result.Output = result.Output + fmt.Sprintf("\n'%s' not found in Cargo.toml", key)
		} else if value != expectedValue {
			result.Passed = false
			result.Output = result.Output + fmt.Sprintf("\nCargo.toml content mismatch, expected '%s', got '%s'", expectedValue, value)
		}
	}
	return result
}

// createFieldMap creates a map of string keys to struct field values
func createFieldMap(conf *Toml) map[string]interface{} {
	fieldMap := map[string]interface{}{
		"package.name":        conf.Package.Name,
		"package.version":     conf.Package.Version,
		"package.edition":     conf.Package.Edition,
		"package.authors":     conf.Package.Authors,
		"package.description": conf.Package.Description,
		"package.default_run": conf.Package.DefaultRun,
	}

	for key, value := range conf.Dependencies {
		fieldMap[fmt.Sprintf("dependencies.%s", key)] = value
	}

	for i, bin := range conf.Bins {
		fieldMap[fmt.Sprintf("bin[%d].name", i)] = bin.Name
		fieldMap[fmt.Sprintf("bin[%d].path", i)] = bin.Path
	}

	for key, profile := range conf.Profiles {
		fieldMap[fmt.Sprintf("profile.%s.inherits", key)] = profile.Inherits
		fieldMap[fmt.Sprintf("profile.%s.strip", key)] = profile.Strip
		fieldMap[fmt.Sprintf("profile.%s.overflow-checks", key)] = profile.OverflowChecks
	}

	return fieldMap
}

// Read a TOML file and returns a field map allowing
// dynamic access to the config's contents
func ReadToml(tomlFilePath string) (map[string]interface{}, error) {
	var conf Toml

	tomlContentBytes, err := os.ReadFile(tomlFilePath)
	if err != nil {
		return nil, err
	}

	if err = toml.Unmarshal(tomlContentBytes, &conf); err != nil {
		return nil, err
	}

	return createFieldMap(&conf), nil
}
