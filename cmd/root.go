package cmd

import (
	"fmt"

	"github.com/42-Short/shortinette/pkg/tester"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "grader",
		Short: "Interactive CLI for grading coding exercises",
		Run: func(cmd *cobra.Command, args []string) {
			startInteractiveMenu()
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}

func startInteractiveMenu() {
	mainMenu := []string{"Grade Exercise", "Exit"}
	var selectedMenu string

	survey.AskOne(&survey.Select{
		Message: "Choose an option:",
		Options: mainMenu,
	}, &selectedMenu)

	switch selectedMenu {
	case "Grade Exercise":
		gradeExercise()
	case "Exit":
		fmt.Println("Goodbye!")
	}
}

func gradeExercise() {
	var filePath string
	survey.AskOne(&survey.Input{
		Message: "Enter the path to the exercise file:",
	}, &filePath)

	result, err := shortinette.RunUserCode(filePath)
	if err != nil {
		fmt.Println("Error running the code:", err)
		return
	}

	fmt.Println("Result:", result)
	startInteractiveMenu()
}
