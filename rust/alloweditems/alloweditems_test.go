package alloweditems

import (
	"io/fs"
	"os"
	"testing"
)

func TestOneLineComment(t *testing.T) {
	var allowedKeywords = map[string]int{"match": 0}
	oneLineComment := `
	// match
	fn main() {
		}
		`
	os.WriteFile("test.rs", []byte(oneLineComment), fs.FileMode(os.O_CREATE))
	defer os.RemoveAll("test.rs")

	if err := allowedKeywordsCheck([]string{"test.rs"}, allowedKeywords); err != nil {
		t.Fatalf("`match` should not throw an error in comments")
	}
}

func TestMultiLineComment(t *testing.T) {
	var allowedKeywords = map[string]int{"match": 0}
	oneLineComment := `
	// hello world
	// match
	// helleas
	fn main() {
		}
		`
	os.WriteFile("test.rs", []byte(oneLineComment), fs.FileMode(os.O_CREATE))
	defer os.RemoveAll("test.rs")

	if err := allowedKeywordsCheck([]string{"test.rs"}, allowedKeywords); err != nil {
		t.Fatalf("`match` should not throw an error in comments")
	}
}

func TestOneLineDoc(t *testing.T) {
	var allowedKeywords = map[string]int{"match": 0}
	oneLineComment := `
	/// match
	fn main() {
		}
		`
	os.WriteFile("test.rs", []byte(oneLineComment), fs.FileMode(os.O_CREATE))
	defer os.RemoveAll("test.rs")

	if err := allowedKeywordsCheck([]string{"test.rs"}, allowedKeywords); err != nil {
		t.Fatalf("`match` should not throw an error in docs")
	}
}

func TestMultiLineDoc(t *testing.T) {
	var allowedKeywords = map[string]int{"match": 0}
	oneLineComment := `
	/// hello world
	/// match
	/// helleas
	fn main() {
		}
		`
	os.WriteFile("test.rs", []byte(oneLineComment), fs.FileMode(os.O_CREATE))
	defer os.RemoveAll("test.rs")

	if err := allowedKeywordsCheck([]string{"test.rs"}, allowedKeywords); err != nil {
		t.Fatalf("`match` should not throw an error in docs")
	}
}

func TestInFunctionNamePrefix(t *testing.T) {
	var allowedKeywords = map[string]int{"match": 0}
	oneLineComment := `
	fn match_to_string() {
	}
		`
	os.WriteFile("test.rs", []byte(oneLineComment), fs.FileMode(os.O_CREATE))
	defer os.RemoveAll("test.rs")

	if err := allowedKeywordsCheck([]string{"test.rs"}, allowedKeywords); err != nil {
		t.Fatalf("`match` should only throw an error as a keyword, not in function names")
	}
}

func TestInFunctionNameSuffix(t *testing.T) {
	var allowedKeywords = map[string]int{"match": 0}
	oneLineComment := `
	fn string_match() {
	}
		`
	os.WriteFile("test.rs", []byte(oneLineComment), fs.FileMode(os.O_CREATE))
	defer os.RemoveAll("test.rs")

	if err := allowedKeywordsCheck([]string{"test.rs"}, allowedKeywords); err != nil {
		t.Fatalf("`match` should only throw an error as a keyword, not in function names")
	}
}

func TestIsUsed(t *testing.T) {
	var allowedKeywords = map[string]int{"match": 0}
	oneLineComment := `
	fn main() {
		match {
			
		}
	}
		`
	os.WriteFile("test.rs", []byte(oneLineComment), fs.FileMode(os.O_CREATE))
	defer os.RemoveAll("test.rs")

	if err := allowedKeywordsCheck([]string{"test.rs"}, allowedKeywords); err == nil {
		t.Fatalf("should throw an error here")
	}
}

func TestIsUsedAllowedAmount(t *testing.T) {
	var allowedKeywords = map[string]int{"match": 5}
	oneLineComment := `
	fn main() {
		match {
			match {
				match {
					match {
						match {}
					}
				}
			}
		}
	}
		`
	os.WriteFile("test.rs", []byte(oneLineComment), fs.FileMode(os.O_CREATE))
	defer os.RemoveAll("test.rs")

	if err := allowedKeywordsCheck([]string{"test.rs"}, allowedKeywords); err != nil {
		t.Fatalf("only used 5 times, allowed 5 times")
	}
}