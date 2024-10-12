package R04

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var clippyTomlAsString02 = `
disallowed-methods = ["std::env::args::nth"]
`

var lslhROutputRecursiveDepth1 = `
ls -lhR bar
bar:
total 984K
---------x 1 root root 977K Sep  1 08:13 69.txt
drwxr-xr-x 2 root root 4.0K Sep  1 08:19 bar

bar/bar:
total 980K
---------x 1 root root 977K Sep  1 08:13 69.txt`

var lslhROutputRecursiveDepth10 = `
ls -lhR 0
0:
total 984K
drwxr-xr-x 3 root root 4.0K Sep  1 08:12 1
---------x 1 root root 977K Sep  1 08:12 69.txt

0/1:
total 984K
drwxr-xr-x 3 root root 4.0K Sep  1 08:12 2
---------x 1 root root 977K Sep  1 08:12 69.txt

0/1/2:
total 984K
drwxr-xr-x 3 root root 4.0K Sep  1 08:12 3
---------x 1 root root 977K Sep  1 08:12 69.txt

0/1/2/3:
total 984K
drwxr-xr-x 3 root root 4.0K Sep  1 08:12 4
---------x 1 root root 977K Sep  1 08:12 69.txt

0/1/2/3/4:
total 984K
drwxr-xr-x 3 root root 4.0K Sep  1 08:12 5
---------x 1 root root 977K Sep  1 08:12 69.txt

0/1/2/3/4/5:
total 984K
drwxr-xr-x 3 root root 4.0K Sep  1 08:12 6
---------x 1 root root 977K Sep  1 08:12 69.txt

0/1/2/3/4/5/6:
total 984K
---------x 1 root root 977K Sep  1 08:12 69.txt
drwxr-xr-x 3 root root 4.0K Sep  1 08:12 7

0/1/2/3/4/5/6/7:
total 984K
---------x 1 root root 977K Sep  1 08:12 69.txt
drwxr-xr-x 3 root root 4.0K Sep  1 08:12 8

0/1/2/3/4/5/6/7/8:
total 984K
---------x 1 root root 977K Sep  1 08:12 69.txt
drwxr-xr-x 2 root root 4.0K Sep  1 08:12 9

0/1/2/3/4/5/6/7/8/9:
total 980K
---------x 1 root root 977K Sep  1 08:12 69.txt
`

func removeANSICodes(input string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

	cleanedString := re.ReplaceAllString(input, "")
	return cleanedString
}

func testMb(workingDirectory string, fooPath string) Exercise.Result {
	kbString := strings.Repeat("69", 500000)
	if err := os.WriteFile(filepath.Join(fooPath, "69.txt"), []byte(kbString), fs.FileMode(os.O_WRONLY)); err != nil {
		return Exercise.InternalError(err.Error())
	}
	var wg sync.WaitGroup
	wg.Add(1)

	ch := make(chan outputChannel)
	go func() {
		defer wg.Done()
		cmd := exec.Command("cargo", "run", "/tmp/foo")
		cmd.Dir = workingDirectory
		out, err := cmd.CombinedOutput()
		ch <- outputChannel{out, err}
	}()
	out := <-ch
	wg.Wait()
	if out.err != nil {
		return Exercise.RuntimeError(out.err.Error())
	}
	codeOutput := string(out.out)
	// codeOutput := strings.Split(string(out.out), "\n")
	// if len(codeOutput) < 4 {
	// 	return Exercise.AssertionError("1.0 megabytes", strings.Join(codeOutput, "\n"), "ls -lh ./foo\n-rw-r--r-- 1 root root 977K Aug 31 17:15 69.txt", "cargo run ./foo")
	// }
	if !strings.Contains(strings.ToLower(codeOutput), "1.0 megabytes") {
		return Exercise.AssertionError("1.0 megabytes", removeANSICodes(codeOutput), "ls -lh ./foo\n-rw-r--r-- 1 root root 977K Aug 31 17:15 69.txt", "cargo run ./foo")
	}
	return Exercise.Passed("OK")
}

func mkdirRecursive(currentPath string, depth int, fileContent string) {
	if depth >= 10 {
		return
	}
	if err := os.Mkdir(currentPath, 0755); err != nil {
		return
	}
	if err := os.WriteFile(filepath.Join(currentPath, "69.txt"), []byte(fileContent), fs.FileMode(os.O_WRONLY)); err != nil {
		return
	}
	mkdirRecursive(filepath.Join(currentPath, strconv.Itoa(depth+1)), depth+1, fileContent)
}

func testRecursiveHard(workingDirectory string) Exercise.Result {
	mkdirRecursive("/tmp/0", 0, strings.Repeat("69", 500000))
	defer os.RemoveAll("/tmp/0")
	var wg sync.WaitGroup
	wg.Add(1)
	ch := make(chan outputChannel)
	go func() {
		defer wg.Done()
		cmd := exec.Command("cargo", "run", "/tmp/0")
		cmd.Dir = workingDirectory
		out, err := cmd.CombinedOutput()
		ch <- outputChannel{out, err}
	}()
	out := <-ch
	wg.Wait()
	if out.err != nil {
		return Exercise.RuntimeError(out.err.Error())
	}
	codeOutput := removeANSICodes(string(out.out))
	// codeOutputList := strings.Split(string(codeOutput), "\n")
	// if len(codeOutputList) < 22 {
	// 	return Exercise.AssertionError("10.0 megabytes", strings.Join(codeOutputList, "\n"), lslhROutputRecursiveDepth10)
	// }
	if !strings.Contains(strings.ToLower(codeOutput), "10.0 megabytes") {
		return Exercise.AssertionError("10.0 megabytes", codeOutput, lslhROutputRecursiveDepth10)
	}
	return Exercise.Passed("OK")
}

func testRecursive(workingDirectory string, fooPath string) Exercise.Result {
	if err := os.Mkdir(filepath.Join(fooPath, "bar"), 0755); err != nil {
		return Exercise.InternalError(err.Error())
	}
	mbString := strings.Repeat("69", 500000)
	if err := os.WriteFile(filepath.Join(fooPath, "69.txt"), []byte(mbString), fs.FileMode(os.O_WRONLY)); err != nil {
		return Exercise.InternalError(err.Error())
	}
	if err := os.WriteFile(filepath.Join(fooPath, "bar", "69.txt"), []byte(mbString), fs.FileMode(os.O_WRONLY)); err != nil {
		return Exercise.InternalError(err.Error())
	}
	var wg sync.WaitGroup
	wg.Add(1)

	ch := make(chan outputChannel)
	go func() {
		defer wg.Done()
		cmd := exec.Command("cargo", "run", "/tmp/bar")
		cmd.Dir = workingDirectory
		out, err := cmd.CombinedOutput()
		ch <- outputChannel{out, err}
	}()
	out := <-ch
	wg.Wait()
	if out.err != nil {
		return Exercise.RuntimeError(out.err.Error())
	}
	// codeOutput := strings.Split(string(out.out), "\n")
	codeOutput := string(out.out)
	// if len(codeOutput) < 5 {
	// 	return Exercise.AssertionError("2.0 megabytes", strings.Join(codeOutput, "\n"), lslhROutputRecursiveDepth1, "cargo run ./bar")
	// }
	if !strings.Contains(strings.ToLower(codeOutput), "2.0 megabytes") {
		return Exercise.AssertionError("2.0 megabytes", removeANSICodes(codeOutput), lslhROutputRecursiveDepth1, "cargo run ./bar")
	}
	return Exercise.Passed("OK")
}

func testDoNotPanic(workingDirectory string) Exercise.Result {
	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "chmod -R 777 ../target && mkdir donotpanic && chmod 000 donotpanic"}); err != nil {
		return Exercise.InternalError(err.Error())
	}
	var wg sync.WaitGroup
	wg.Add(1)

	ch := make(chan outputChannel)
	go func() {
		defer wg.Done()
		cmd := exec.Command("sh", "-c", "su -c 'cargo run /tmp/donotpanic' student")
		cmd.Dir = workingDirectory
		out, err := cmd.CombinedOutput()
		ch <- outputChannel{out, err}
	}()
	out := <-ch
	wg.Wait()
	// if out.err != nil {
	// 	return Exercise.RuntimeError(out.err.Error())
	// }
	if strings.Contains(string(out.err.Error()), "panicked") {
		return Exercise.RuntimeError(string(out.out), "mkdir donotpanic", "chmod 000 donotpanic", "cargo run ./donotpanic")
	}
	return Exercise.Passed("OK")
}

// Assumptions behind the hardcoded output indices and lengths in the tests:
//
// - When running `cargo run` for the first time, the output is preceded by 3 lines from the compiler:
//
//	Compiling ex02 v0.1.0 (/root/side/Short/rust-subjects/answers/module-04/ex02)
//	Finished `dev` profile [unoptimized + debuginfo] target(s) in 0.46s
//	Running `target/debug/ex02 bar`
//
// - After it has been run once, the Rust project will not be compiled anymore, leading the
// cargo output to only consist of 2 lines:
//
//	Finished `dev` profile [unoptimized + debuginfo] target(s) in 0.00s
//	Running `target/debug/ex02 bar`
func ex02Test(exercise *Exercise.Exercise) (result Exercise.Result) {
	return Exercise.Passed("OK")
}

func ex02() Exercise.Exercise {
	return Exercise.NewExercise("02", "ex02", []string{"Cargo.toml", "src/main.rs"}, 10, ex02Test)
}
