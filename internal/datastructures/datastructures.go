package datastructures

type Test struct {
	Input          string `yaml:"input,omitempty"`
	ExpectedOutput string `yaml:"expected_output,omitempty"`
}

type AllowedItems struct {
	Macros    []string  `yaml:"macros"`
	Functions []string  `yaml:"functions"`
	Keywords  []Keyword `yaml:"keywords"`
}

type Keyword struct {
	Keyword string `yaml:"keyword"`
	Amount  int    `yaml:"amount"`
}

type Exercise struct {
	TurnInDirectory string       `yaml:"turn_in_directory"`
	TurnInFile      string       `yaml:"turn_in_file"`
	AllowedItems    AllowedItems `yaml:"allowed_items"`
	Tests           []Test       `yaml:"tests"`
	TestsPath       string       `yaml:"tests_path,omitempty"`
	MainPath        string       `yaml:"main_path,omitempty"`
	Type            string       `yaml:"type"`
	DummyCall       string       `yaml:"dummy_call,omitempty"`
	AssertEq        string       `yaml:"assert_eq,omitempty"`
	SubExercises    []Exercise   `yaml:"sub_exercises,omitempty"`
}

type Config struct {
	Exercises map[string]Exercise `yaml:"exercises"`
}

type AllowedItem struct {
	Name string
	Type string
}
