package modulebuilder

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/internal/git"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/testbuilder"
)

type Module struct {
	Name      string
	Exercises []testbuilder.TestBuilder
}

type ModuleBuilder interface {
	SetName(name string) ModuleBuilder
	SetExercises(exercises []testbuilder.TestBuilder) ModuleBuilder
	SetUp(repoId string, codeDirectory string) error
	Build() Module
	Run() []testbuilder.Result
}

type ModuleBuilderImpl struct {
	name      string
	exercises []testbuilder.TestBuilder
}

func NewModuleBuilder() ModuleBuilder {
	return &ModuleBuilderImpl{}
}

func (b *ModuleBuilderImpl) SetName(name string) ModuleBuilder {
	b.name = name
	return b
}

func (b *ModuleBuilderImpl) SetExercises(exercises []testbuilder.TestBuilder) ModuleBuilder {
	b.exercises = exercises
	return b
}

func (b *ModuleBuilderImpl) SetUp(repoId string, testDirectory string) error {
	repoLink := fmt.Sprintf("https://github.com/%s/%s.git", os.Getenv("GITHUB_ORGANISATION"), repoId)
	if err := git.Get(repoLink, testDirectory); err != nil {
		errorMessage := fmt.Sprintf("failed to clone repository: %v", err)
		return errors.NewInternalError(errors.ErrInternal, errorMessage)
	}
	if err := logger.InitializeTraceLogger(repoId); err != nil {
		errorMessage := fmt.Sprintf("failed to initalize logging system (%v), does the ./traces directory exist?", err)
		return errors.NewInternalError(errors.ErrInternal, errorMessage)
	}
	if err := git.Get(fmt.Sprintf("https://github.com/%s/%s.git", os.Getenv("GITHUB_ORGANISATION"), repoId), "compile-environment/src/"); err != nil {
		return err
	}
	return nil
}

func (b *ModuleBuilderImpl) Build() Module {
	return Module{
		Name:      b.name,
		Exercises: b.exercises,
	}
}

func (b *ModuleBuilderImpl) Run() []testbuilder.Result {
	var results []testbuilder.Result
	if b.exercises != nil {
		for _, exercise := range b.exercises {
			res := exercise.Run()
			logger.File.Printf("[%s]: %t", exercise.Build().Name, res.Passed)
			results = append(results, res)
		}
	}
	return results
}
