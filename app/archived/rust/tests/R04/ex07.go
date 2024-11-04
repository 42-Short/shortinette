//go:build ignore
package R04

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/42-Short/shortinette/rust/alloweditems"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var clippyTomlAsString07 = ``

func filterEmpty(slice []string) []string {
	var result []string
	for _, str := range slice {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}

func testKeyGen(workingDirectory string) (result Exercise.Result) {
	commandLine := "cargo run -- gen-keys key.prv key.pub"
	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", commandLine}); err != nil {
		return Exercise.RuntimeError(err.Error(), commandLine)
	}

	if _, err := os.OpenFile(filepath.Join(workingDirectory, "key.prv"), os.O_RDONLY, 0644); os.IsNotExist(err) {
		return Exercise.RuntimeError("did not create file 'key.prv'", commandLine)
	}
	if _, err := os.OpenFile(filepath.Join(workingDirectory, "key.pub"), os.O_RDONLY, 0644); os.IsNotExist(err) {
		return Exercise.RuntimeError("did not create file 'key.pub'", commandLine)
	}

	prvKeyBytes, err := os.ReadFile(filepath.Join(workingDirectory, "key.prv"))
	if err != nil {
		return Exercise.InternalError(err.Error())
	}
	pubKeyBytes, err := os.ReadFile(filepath.Join(workingDirectory, "key.pub"))
	if err != nil {
		return Exercise.InternalError(err.Error())
	}

	prvKeySliced := filterEmpty(strings.Split(string(prvKeyBytes), "\n"))
	pubKeySliced := filterEmpty(strings.Split(string(pubKeyBytes), "\n"))
	if len(prvKeySliced) != 2 || len(pubKeySliced) != 2 {
		return Exercise.RuntimeError(fmt.Sprintf("private and public key must be 2 lines each (line 1 E/D, line 2 M)\nkey.prv (%d lines, exp. 2):\n%s\nkey.pub (%d lines, exp. 2):\n%s", len(pubKeySliced), string(pubKeyBytes), len(prvKeySliced), string(prvKeyBytes)), commandLine)
	}
	if prvKeySliced[1] != pubKeySliced[1] {
		return Exercise.RuntimeError("M(key.prv) != M(key.pub)", commandLine)
	}
	return Exercise.Passed("OK")
}

func testEncryption(workingDirectory string) (result Exercise.Result) {
	encryptedMessage := "encrypt me"
	commandLine := fmt.Sprintf("echo '%s' | cargo run -- encrypt key.pub > encrypted", encryptedMessage)
	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", commandLine}); err != nil {
		return Exercise.RuntimeError(err.Error())
	}

	encryptedBytes, err := os.ReadFile(filepath.Join(workingDirectory, "encrypted"))
	if err != nil {
		return Exercise.InternalError(err.Error())
	}

	if string(encryptedBytes) == encryptedMessage {
		return Exercise.RuntimeError("the encrypted message should not be the same as the original, nice try", "cargo run -- gen-keys key.prv key.pub", commandLine)
	}

	return Exercise.Passed("OK")
}

func testDecryption(workingDirectory string) (result Exercise.Result) {
	expectedMessage := "encrypt me"
	commandLine := "cat encrypted | cargo run -- decrypt key.prv"
	output, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", commandLine})
	if err != nil {
		return Exercise.RuntimeError(err.Error(), "cargo run -- gen-keys key.prv key.pub", "echo 'encrypt me' | cargo run -- encrypt key.pub > encrypted", commandLine)
	}
	if output[:len(expectedMessage)] != expectedMessage {
		return Exercise.AssertionError(expectedMessage, output[:len(expectedMessage)])
	}
	return Exercise.Passed("OK")
}

func ex07Test(exercise *Exercise.Exercise) (result Exercise.Result) {
	if err := alloweditems.Check(*exercise, clippyTomlAsString07, nil, "#![allow(clippy::slow_vector_initialization)]"); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if result = testNoInput(workingDirectory); !result.Passed {
		return result
	}
	if result = testKeyGen(workingDirectory); !result.Passed {
		return result
	}
	if result = testEncryption(workingDirectory); !result.Passed {
		return result
	}
	if result = testDecryption(workingDirectory); !result.Passed {
		return result
	}

	return Exercise.Passed("OK")
}

func ex07() Exercise.Exercise {
	return Exercise.NewExercise("07", "ex07", []string{"Cargo.toml", "src/main.rs"}, 20, ex07Test)
}
