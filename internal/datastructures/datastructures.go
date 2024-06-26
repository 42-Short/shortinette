package datastructures

type Test struct {
	Output struct {
		AssertEq []string `yaml:"assert_eq"`
		AssertNe []string `yaml:"assert_ne"`
	}
}

type Exercise struct {
	TurnInDirectory string   `yaml:"turn_in_directory"`
	TurnInFiles     []string `yaml:"turn_in_files"`
	AllowedItems    []string `yaml:"allowed_items"`
	Tests           Test     `yaml:"tests"`
}

type Config struct {
	Ex00 Exercise `yaml:"ex00"`
}

type AllowedItem struct {
	Name string
	Type string
}
