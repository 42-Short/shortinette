package datastructures

type Test struct {
	AssertEq []string `yaml:"assert_eq"`
	AssertNe []string `yaml:"assert_ne"`
}

type AllowedItems struct {
	Macros    []string `yaml:"macros"`
	Functions []string `yaml:"functions"`
}

type Exercise struct {
	TurnInDirectory string       `yaml:"turn_in_directory"`
	TurnInFile      string       `yaml:"turn_in_file"`
	AllowedItems    AllowedItems `yaml:"allowed_items"`
	Tests           Test         `yaml:"tests"`
}

type Config struct {
	Ex00 Exercise `yaml:"ex00"`
}

type AllowedItem struct {
	Name string
	Type string
}
