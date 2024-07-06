package toml

import (
	"os"

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

// ReadToml reads a TOML file into a Toml struct, where values can then be easily accessed
func ReadToml(tomlFilePath string) (*Toml, error) {
	var conf Toml

	tomlContentBytes, err := os.ReadFile(tomlFilePath)
	if err != nil {
		return nil, err
	}

	if err = toml.Unmarshal(tomlContentBytes, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
