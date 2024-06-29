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
	TurnInDirectory   string       `yaml:"turn_in_directory"`
	TurnInFile        string       `yaml:"turn_in_file"`
	AllowedItems      AllowedItems `yaml:"allowed_items"`
	ForbiddenKeywords []string     `yaml:"forbidden_keywords,omitempty"`
	Tests             Test         `yaml:"tests"`
	TestsPath         string       `yaml:"tests_path,omitempty"`
	Type              string       `yaml:"type"`
	DummyCall         string       `yaml:"dummy_call,omitempty"`
}

type Config struct {
	Exercises map[string]Exercise `yaml:"exercises"`
}

type AllowedItem struct {
	Name string
	Type string
}
