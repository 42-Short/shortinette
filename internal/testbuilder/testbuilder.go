package testbuilder

type Result struct {
	Passed bool
}

type Test struct {
	Name             string
	TurnInDirectory  string
	TurnInFile       string
	Type             string
	Prototype        string
	AllowedMacros    []string
	AllowedFunctions []string
	AllowedKeywords  map[string]int
	Executer         func(test *Test) bool
}

type TestBuilder interface {
	SetName(name string) TestBuilder
	SetTurnInDirectory(turnInDirectory string) TestBuilder
	SetTurnInFile(turnInFile string) TestBuilder
	SetExerciseType(exerciseType string) TestBuilder
	SetPrototype(prototype string) TestBuilder
	SetAllowedMacros(allowedMacros []string) TestBuilder
	SetAllowedFunctions(allowedFunctions []string) TestBuilder
	SetAllowedKeywords(allowedKeywords map[string]int) TestBuilder
	SetExecuter(executer func(test *Test) bool) TestBuilder
	Build() Test
	Run() Result
}

type TestBuilderImpl struct {
	name             string
	turnInDirectory  string
	turnInFile       string
	exerciseType     string
	prototype        string
	allowedMacros    []string
	allowedFunctions []string
	allowedKeywords  map[string]int
	executer         func(test *Test) bool
}

func NewTestBuilder() TestBuilder {
	return &TestBuilderImpl{}
}

func (b *TestBuilderImpl) SetName(name string) TestBuilder {
	b.name = name
	return b
}

func (b *TestBuilderImpl) SetTurnInDirectory(turnInDirectory string) TestBuilder {
	b.turnInDirectory = turnInDirectory
	return b
}

func (b *TestBuilderImpl) SetTurnInFile(turnInFile string) TestBuilder {
	b.turnInFile = turnInFile
	return b
}

func (b *TestBuilderImpl) SetExerciseType(exerciseType string) TestBuilder {
	b.exerciseType = exerciseType
	return b
}

func (b *TestBuilderImpl) SetPrototype(prototype string) TestBuilder {
	b.prototype = prototype
	return b
}

func (b *TestBuilderImpl) SetAllowedMacros(allowedMacros []string) TestBuilder {
	b.allowedMacros = allowedMacros
	return b
}

func (b *TestBuilderImpl) SetAllowedFunctions(allowedFunctions []string) TestBuilder {
	b.allowedFunctions = allowedFunctions
	return b
}

func (b *TestBuilderImpl) SetAllowedKeywords(allowedKeywords map[string]int) TestBuilder {
	b.allowedKeywords = allowedKeywords
	return b
}

func (b *TestBuilderImpl) SetExecuter(executer func(test *Test) bool) TestBuilder {
	b.executer = executer
	return b
}

func (b *TestBuilderImpl) Build() Test {
	return Test{
		Name:             b.name,
		TurnInDirectory:  b.turnInDirectory,
		TurnInFile:       b.turnInFile,
		AllowedMacros:    b.allowedMacros,
		AllowedFunctions: b.allowedFunctions,
		AllowedKeywords:  b.allowedKeywords,
		Executer:         b.executer,
	}
}

func (b *TestBuilderImpl) Run() Result {
	test := b.Build()
	if b.executer != nil {
		return Result{Passed: b.executer(&test)}
	}
	return Result{Passed: false}
}
