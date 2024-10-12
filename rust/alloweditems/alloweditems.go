package alloweditems

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

func prependLintLevel(filePath string, lintLevelModifications []string) (err error) {
	contentAsBytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	contentAsStringSlice := strings.Split(string(contentAsBytes), "\n")
	for index, line := range contentAsStringSlice {
		if !strings.HasPrefix(line, "#![") {
			for _, modification := range lintLevelModifications {
				contentAsStringSlice = slices.Insert(contentAsStringSlice, index, modification)
			}
			break
		}
	}

	err = os.WriteFile(filePath, []byte(strings.Join(contentAsStringSlice, "\n")), 0644)
	if err != nil {
		return err
	}

	return nil
}

func concatenateFilesIntoString(files []string) (fileContents string, err error) {
	res := ""
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return "", err
		}
		res += string(content)
	}
	return res, nil
}

func wordBoundaryCheck(character byte) bool {
	return unicode.IsLetter(rune(character)) || unicode.IsDigit(rune(character)) || character == '_'
}

func getRegexResults(keywordsSlice []string, cleanFileBytes []byte, allowedKeywords map[string]int) (err error) {
	escapedKeywords := make([]string, len(keywordsSlice))
	for idx, keyword := range keywordsSlice {
		escapedKeywords[idx] = regexp.QuoteMeta(keyword)
		if wordBoundaryCheck(escapedKeywords[idx][0]) {
			escapedKeywords[idx] = `\b` + escapedKeywords[idx]
		}
		length := len(escapedKeywords[idx])
		if wordBoundaryCheck(escapedKeywords[idx][length-1]) {
			escapedKeywords[idx] = escapedKeywords[idx] + `\b`
		}
	}

	keywordExpr, err := regexp.Compile(`(` + strings.Join(escapedKeywords, "|") + `)`)
	if err != nil {
		return err
	}

	preprocessedBytes := []byte(strings.ReplaceAll(string(cleanFileBytes), "->", ""))
	preprocessedBytes = []byte(strings.ReplaceAll(string(preprocessedBytes), "=>", ""))

	if foundKeyWords := keywordExpr.FindAll(preprocessedBytes, -1); foundKeyWords != nil {
		for _, keyword := range foundKeyWords {
			allowedKeywords[string(keyword)] -= 1
		}
	}

	badKeywords := []string{}
	for keyword, amount := range allowedKeywords {
		if amount < 0 {
			badKeywords = append(badKeywords, fmt.Sprintf("%s: %d", keyword, amount*-1))
		}
	}
	if len(badKeywords) != 0 {
		return fmt.Errorf("forbidden/limited keywords found (keyword: amount): %s", strings.Join(badKeywords, ", "))
	}
	return nil
}

func removeTestModules(content []byte) []byte {
	testModuleRegex := regexp.MustCompile(`(?s)#\[cfg\(test\)\][\s\n]*mod[\s\n]+test[\s\n]*\{.*?\}`)
	return testModuleRegex.ReplaceAll(content, []byte{})
}

func allowedKeywordsCheck(filesToCheck []string, allowedKeywords map[string]int) (err error) {
	keywordsSlice := []string{}
	for keyword := range allowedKeywords {
		keywordsSlice = append(keywordsSlice, keyword)
	}

	filesAsString, err := concatenateFilesIntoString(filesToCheck)
	if err != nil {
		return err
	}

	ignoreExpr, _ := regexp.Compile(`(?m)(?s)//.*?$|///.*?$|/\*.*?\*/|"(?:\\.|[^"\\])*"|r#*"(?:.|\n)*?"#*`)
	cleanFileBytes := ignoreExpr.ReplaceAll([]byte(filesAsString), []byte(""))

	cleanFileBytes = removeTestModules(cleanFileBytes)

	return getRegexResults(keywordsSlice, cleanFileBytes, allowedKeywords)
}

func allowedItemsCheck(clippyTomlAsString string, exercise Exercise.Exercise) (err error) {
	file, err := os.Create(filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory, ".clippy.toml"))
	if err != nil {
		return err
	}
	if _, err = file.WriteString(clippyTomlAsString); err != nil {
		return err
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	if _, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"clippy", "--", "-D", "warnings"}); err != nil {
		return err
	}
	return nil
}

func getRustFiles(exercise Exercise.Exercise) []string {
	files := []string{}
	for _, file := range exercise.TurnInFiles {
		if strings.HasSuffix(file, ".rs") {
			files = append(files, file)
		}
	}
	return files
}

// Checks for forbidden methods/macros using `cargo clippy` and for keywords using regex.
// Args:
//
//   - exercise: `Exercise.Exercise` structure containing the exercise metadata
//
//   - clippyTomlAsString: string representation of the `.clippy.toml` file which should dictate the lint rules
//
//   - allowedKeywords: key value pairs -> keyword: allowed amount, nil for none
//
//   - lintLevelModifications (optional): arbitrary amount of lint modifications (#![allow(clippy::doc_lazy_continuation)] and #![allow(dead_code)] are added by default)
//
// Example Usage:
//   - I want to ban `std::ptr::read` and `std::println`
//   - I also want to allow the `match` keyword maximum once
//
// To achieve this, I can call call this function like follows:
//
//	clippyTomlAsString := `
//	disallowed-macros = ["std::println"]
//	disallowed-methods = ["std::ptr::read"]
//	`
//	lintLevelModifications := []string{"#[allow(clippy::doc_lazy_continuation)]"}
//	allowedKeywords := map[string]int{"match": 1}
//	if err := allowedItems.Check(exercise, clippyTomlAsString, allowedKeywords, lintLevelModifications); err != nil {
//		// err != nil -> linting failed, meaning the submission did not pass your static analysis.
//		// err.Error() will contain all necessary information for your trace, such as which line posed an issue,
//		// which disallowed item/keyword(s) was/were found, (...), you can simply handle this as follows:
//		return Exercise.CompilationError(err.Error())
//	}
//
// See https://rust-lang.github.io/rust-clippy/master/index.html for details.
func Check(exercise Exercise.Exercise, clippyTomlAsString string, allowedKeywords map[string]int, lintLevelModifications ...string) (err error) {
	lintLevelModifications = append(lintLevelModifications, "#![allow(clippy::doc_lazy_continuation)]", "#![allow(dead_code)]", "#![allow(clippy::duplicated_attributes)]", "#![allow(clippy::explicit_counter_loop)]")
	filesToCheck := getRustFiles(exercise)

	for _, file := range filesToCheck {
		if err = prependLintLevel(file, lintLevelModifications); err != nil {
			return err
		}
	}

	if err = allowedItemsCheck(clippyTomlAsString, exercise); err != nil {
		return err
	}

	if allowedKeywords != nil {
		if err = allowedKeywordsCheck(filesToCheck, allowedKeywords); err != nil {
			return err
		}
	}
	return nil
}
