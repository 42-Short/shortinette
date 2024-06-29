package functioncheck

import (
	"fmt"
	"os"
	"bufio"
	"strings"

	"github.com/42-Short/shortinette/internal/datastructures"
)

func scanForForbiddenKeywords(scanner *bufio.Scanner, forbiddenKeywords []string) (err error) {
	foundKeywords := make([]string, 0, len(forbiddenKeywords))
	
    for scanner.Scan() {
        word := scanner.Text()
        for _, keyword := range forbiddenKeywords {
            if word == keyword {
                foundKeywords = append(foundKeywords, keyword)
            }
        }
    }
	if len(foundKeywords) > 0 {
		return fmt.Errorf("found forbidden keywords: %s", strings.Join(foundKeywords, ", "))
	}
    return scanner.Err()
}

//Lints the student code and checks for forbidden keywords
func LintStudentCode(exercisePath string, exerciseConfig datastructures.Exercise) (err error) {
	file, err := os.Open(exercisePath)
    if err != nil {
        return fmt.Errorf("could not open %s: %w", exercisePath, err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanWords)

    err = scanForForbiddenKeywords(scanner, exerciseConfig.ForbiddenKeywords)
	if err != nil {
        return err
    }
    return nil
}