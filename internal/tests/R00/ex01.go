package R00

import "github.com/42-Short/shortinette/internal/testbuilder"

func ex01() testbuilder.TestBuilder {
	return testbuilder.NewTestBuilder().
		SetName("EX01").
		SetTurnInDirectory("ex01").
		SetTurnInFile("min.rs").
		SetAllowedMacros([]string{"println"}).
		SetAllowedFunctions(nil).
		SetAllowedKeywords(map[string]int{"unsafe": 0}).
		SetExecuter(ex00Test)
}
