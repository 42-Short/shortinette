package Exercise

type Result struct {
	Passed bool
}

type Exercise struct {
	Name             string
	RepoDirectory    string
	TurnInDirectory  string
	TurnInFiles      []string
	ExerciseType     string
	Prototype        string
	AllowedMacros    []string
	AllowedFunctions []string
	AllowedKeywords  map[string]int
	Executer         func(test *Exercise) bool
}

// NewExercise initializes and returns an Exercise struct
func NewExercise(name string, repoDirectory string, turnInDirectory string, turnInFiles []string, exerciseType string, prototype string,
	allowedMacros []string, allowedFunctions []string, allowedKeywords map[string]int,
	executer func(test *Exercise) bool) Exercise {

	return Exercise{
		Name:             name,
		RepoDirectory:    repoDirectory,
		TurnInDirectory:  turnInDirectory,
		TurnInFiles:      turnInFiles,
		ExerciseType:     exerciseType,
		Prototype:        prototype,
		AllowedMacros:    allowedMacros,
		AllowedFunctions: allowedFunctions,
		AllowedKeywords:  allowedKeywords,
		Executer:         executer,
	}
}

// Run executes the exercise and returns the result
func (e *Exercise) Run() Result {
	if e.Executer != nil {
		return Result{Passed: e.Executer(e)}
	}
	return Result{Passed: false}
}
