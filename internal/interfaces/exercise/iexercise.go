package exercisebuilder

type Result struct {
	Passed bool
}

type Test struct {
	Name             string
	TurnInDirectory  string
	TurnInFile       string
	ExerciseType     string
	Prototype        string
	AllowedMacros    []string
	AllowedFunctions []string
	AllowedKeywords  map[string]int
	Executer         func(test *Test) bool
}

type ExerciseBuilder interface {
	SetName(name string) ExerciseBuilder
	SetTurnInDirectory(turnInDirectory string) ExerciseBuilder
	SetTurnInFile(turnInFile string) ExerciseBuilder
	SetExerciseType(exerciseType string) ExerciseBuilder
	SetPrototype(prototype string) ExerciseBuilder
	SetAllowedMacros(allowedMacros []string) ExerciseBuilder
	SetAllowedFunctions(allowedFunctions []string) ExerciseBuilder
	SetAllowedKeywords(allowedKeywords map[string]int) ExerciseBuilder
	SetExecuter(executer func(test *Test) bool) ExerciseBuilder
	Build() Test
	Run() Result
}

type ExerciseBuilderImpl struct {
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

func NewExerciseBuilder() ExerciseBuilder {
	return &ExerciseBuilderImpl{}
}

func (b *ExerciseBuilderImpl) SetName(name string) ExerciseBuilder {
	b.name = name
	return b
}

func (b *ExerciseBuilderImpl) SetTurnInDirectory(turnInDirectory string) ExerciseBuilder {
	b.turnInDirectory = turnInDirectory
	return b
}

func (b *ExerciseBuilderImpl) SetTurnInFile(turnInFile string) ExerciseBuilder {
	b.turnInFile = turnInFile
	return b
}

func (b *ExerciseBuilderImpl) SetExerciseType(exerciseType string) ExerciseBuilder {
	b.exerciseType = exerciseType
	return b
}

func (b *ExerciseBuilderImpl) SetPrototype(prototype string) ExerciseBuilder {
	b.prototype = prototype
	return b
}

func (b *ExerciseBuilderImpl) SetAllowedMacros(allowedMacros []string) ExerciseBuilder {
	b.allowedMacros = allowedMacros
	return b
}

func (b *ExerciseBuilderImpl) SetAllowedFunctions(allowedFunctions []string) ExerciseBuilder {
	b.allowedFunctions = allowedFunctions
	return b
}

func (b *ExerciseBuilderImpl) SetAllowedKeywords(allowedKeywords map[string]int) ExerciseBuilder {
	b.allowedKeywords = allowedKeywords
	return b
}

func (b *ExerciseBuilderImpl) SetExecuter(executer func(test *Test) bool) ExerciseBuilder {
	b.executer = executer
	return b
}

func (b *ExerciseBuilderImpl) Build() Test {
	return Test{
		Name:             b.name,
		TurnInDirectory:  b.turnInDirectory,
		TurnInFile:       b.turnInFile,
		Prototype:        b.prototype,
		ExerciseType:     b.exerciseType,
		AllowedMacros:    b.allowedMacros,
		AllowedFunctions: b.allowedFunctions,
		AllowedKeywords:  b.allowedKeywords,
		Executer:         b.executer,
	}
}

func (b *ExerciseBuilderImpl) Run() Result {
	test := b.Build()
	if b.executer != nil {
		return Result{Passed: b.executer(&test)}
	}
	return Result{Passed: false}
}
