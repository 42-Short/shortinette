//go:build ignore
package R00

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"github.com/42-Short/shortinette/rust/cargo"
	"strconv"
	"strings"
	"time"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
	"github.com/42-Short/shortinette/rust/alloweditems"
)

const TestExistenceHelper = `
use std::os::unix::net::UnixStream;
use std::io::prelude::*;

fn is_leap_year_helper(year: u32) {
	let text = format!("is_leap_year({});\n", year);
	if let Ok(mut stream) = UnixStream::connect("/var/run/shortinette.sock") {
		stream.write_all(text.as_bytes()).ok();
	}
}

fn num_days_in_month_helper(year: u32, month: u32) {
	let text = format!("num_days_in_month({}, {});\n", year, month);
	if let Ok(mut stream) = UnixStream::connect("/var/run/shortinette.sock") {
		stream.write_all(text.as_bytes()).ok();
	}
}
`

const CargoTestEx05 = `
#[cfg(test)]
mod shortinette_tests_rust_0005 {
	use super::*;

	#[test]
	fn test_is_leap_year() {
		assert!(is_leap_year(1600));
		assert!(!is_leap_year(1500));
		assert!(is_leap_year(2004));
		assert!(!is_leap_year(2003));
	}

	#[test]
	fn test_num_days_in_month() {
		assert_eq!(num_days_in_month(2001, 1), 31);
		assert_eq!(num_days_in_month(2000, 1), 31);
		assert_eq!(num_days_in_month(2000, 2), 29);
		assert_eq!(num_days_in_month(2001, 2), 28);
		assert_eq!(num_days_in_month(2004, 2), 29);
		assert_eq!(num_days_in_month(2100, 2), 28);
		assert_eq!(num_days_in_month(2001, 3), 31);
		assert_eq!(num_days_in_month(2000, 3), 31);
		assert_eq!(num_days_in_month(2001, 4), 30);
		assert_eq!(num_days_in_month(2000, 4), 30);
		assert_eq!(num_days_in_month(2001, 5), 31);
		assert_eq!(num_days_in_month(2000, 5), 31);
		assert_eq!(num_days_in_month(2001, 6), 30);
		assert_eq!(num_days_in_month(2000, 6), 30);
		assert_eq!(num_days_in_month(2001, 7), 31);
		assert_eq!(num_days_in_month(2000, 7), 31);
		assert_eq!(num_days_in_month(2001, 8), 31);
		assert_eq!(num_days_in_month(2000, 8), 31);
		assert_eq!(num_days_in_month(2001, 9), 30);
		assert_eq!(num_days_in_month(2000, 9), 30);
		assert_eq!(num_days_in_month(2001, 10), 31);
		assert_eq!(num_days_in_month(2000, 10), 31);
		assert_eq!(num_days_in_month(2001, 11), 30);
		assert_eq!(num_days_in_month(2000, 11), 30);
		assert_eq!(num_days_in_month(2001, 12), 31);
		assert_eq!(num_days_in_month(2000, 12), 31);
	}

	#[test]
	#[should_panic]
	fn test_invalid_year() {
		is_leap_year(0);
	}

	#[test]
	#[should_panic]
	fn test_invalid_month() {
		num_days_in_month(2000, 13);
	}

	#[test]
	#[should_panic]
	fn test_invalid_month_2() {
		num_days_in_month(2001, 13);
	}

	#[test]
	#[should_panic]
	fn test_invalid_month_3() {
		num_days_in_month(0, 2);
	}
}
`

func runUnixSocket(ctx context.Context, output *strings.Builder, ready chan<- bool) {
	socketPath := `/var/run/shortinette.sock`
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return
	}
	defer os.Remove(socketPath)
	defer listener.Close()
	ready <- true
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				continue
			}
			go func(conn net.Conn) {
				defer conn.Close()
				reader := bufio.NewReader(conn)
				for {
					select {
					case <-ctx.Done():
						return
					default:
						line, err := reader.ReadString('\n')
						if err != nil {
							if ctx.Err() != nil {
								return
							}
							break
						}
						output.WriteString(line)
					}
				}
			}(conn)
		}
	}
}

func addHelperToFunction(exercise *Exercise.Exercise, function string, helper string) error {
	content, err := os.ReadFile(exercise.TurnInFiles[0])
	if err != nil {
		return err
	}
	pattern := `fn\s+` + function + `\s*\(\s*(\w+)\s*:\s*\w+\s*(?:,\s*(\w+)\s*:\s*\w+\s*)?\)`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	total := 0
	for {
		result := re.FindAllSubmatchIndex(content, -1)
		if len(result) == total {
			if total == 0 {
				return nil
			}
			return os.WriteFile(exercise.TurnInFiles[0], content, 0644)
		}
		match := result[total]
		var arguments []string
		for i := 2; i < len(match); i += 2 {
			start := match[i]
			end := match[i+1]
			if start != -1 && end != -1 {
				arguments = append(arguments, string(content[start:end]))
			}
		}
		for pos := match[1]; pos < len(content); pos++ {
			if content[pos] == '{' {
				pos++
				argumentsString := strings.Join(arguments, ", ")
				content = append(content[:pos], append([]byte(fmt.Sprintf(helper, argumentsString)), content[pos:]...)...)
				total++
				break
			}
		}
	}
}

func appendHelperFunctions(exercise *Exercise.Exercise) error {
	if err := addHelperToFunction(exercise, "is_leap_year", `is_leap_year_helper(%s);`); err != nil {
		return err
	}
	if err := addHelperToFunction(exercise, "num_days_in_month", `num_days_in_month_helper(%s);`); err != nil {
		return err
	}
	if err := testutils.AppendStringToFile(TestExistenceHelper, exercise.TurnInFiles[0]); err != nil {
		return err
	}
	return nil
}

func missingTests(message string) Exercise.Result {
	fmt.Println(message)
	return Exercise.Passed("OK")
}

func isLeapYearTest(content string, result *strings.Builder) error {
	isLeapYearFound := make(map[uint32]bool)
	for _, year := range []uint32{1600, 1500, 2004, 2003, 0} {
		isLeapYearFound[year] = false
	}
	pattern := `is_leap_year\((\d+)\)`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) == 2 {
			year, err := strconv.ParseUint(match[1], 10, 32)
			if err != nil {
				return err
			}
			for expectedYear := range isLeapYearFound {
				if uint32(year) == expectedYear {
					isLeapYearFound[uint32(year)] = true
				}
			}
		}
	}
	for year, found := range isLeapYearFound {
		if !found {
			result.WriteString(fmt.Sprintf("Missing test for is_leap_year() with year: %d\n", year))
		}
	}
	return nil
}

func numDaysInMonthTest(content string, result *strings.Builder) error {
	leapYearsFound := make(map[uint32]bool)
	commonYearsFound := make(map[uint32]bool)
	invalidMonthFound := false

	for month := uint32(1); month <= 12; month++ {
		leapYearsFound[month] = false
		commonYearsFound[month] = false
	}
	pattern := `num_days_in_month\((\d+),\s*(\d+)\)`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) == 3 {
			year, err := strconv.ParseUint(match[1], 10, 32)
			if err != nil {
				return err
			}
			month, err := strconv.ParseUint(match[2], 10, 32)
			if err != nil {
				return err
			}
			if month < 1 || month > 12 {
				invalidMonthFound = true
			} else {
				if year > 0 && (year%4 != 0 || (year%100 == 0 && year%400 != 0)) {
					commonYearsFound[uint32(month)] = true
				} else {
					leapYearsFound[uint32(month)] = true
				}
			}
		}
	}
	if !invalidMonthFound {
		result.WriteString("Missing test for num_days_in_month() with invalid month\n")
	}
	for month, found := range leapYearsFound {
		if !found {
			result.WriteString(fmt.Sprintf("Missing test for num_days_in_month() with a leap year and month: %d\n", month))
		}
	}
	for month, found := range commonYearsFound {
		if !found {
			result.WriteString(fmt.Sprintf("Missing test for num_days_in_month() with a common year and month: %d\n", month))
		}
	}
	return nil
}

func testExistenceTest(exercise *Exercise.Exercise) Exercise.Result {
	if err := appendHelperFunctions(exercise); err != nil {
		return Exercise.InternalError(err.Error())
	}
	var output strings.Builder
	ctx, cancel := context.WithCancel(context.Background())
	ready := make(chan bool)
	defer cancel()
	go runUnixSocket(ctx, &output, ready)
	<-ready
	if result := cargo.CargoTest(exercise, 1*time.Second, []string{}); !result.Passed {
		return result
	}
	cancel()
	content := output.String()
	var result strings.Builder
	if err := isLeapYearTest(content, &result); err != nil {
		return Exercise.InternalError(err.Error())
	}
	if err := numDaysInMonthTest(content, &result); err != nil {
		return Exercise.InternalError(err.Error())
	}
	if result.String() != "" {
		return missingTests(result.String())
	}
	return Exercise.Passed("")
}

func expectedMainOutput() string {
	currentDate := time.Now()
	year, month, day := currentDate.Date()

	var result strings.Builder
	for i := 1; i <= year; i++ {
		maxMonth := time.December
		if i == year {
			maxMonth = month
			if day < 13 {
				maxMonth--
			}
		}
		for j := time.January; j <= maxMonth; j++ {
			date := time.Date(i, j, 13, 0, 0, 0, 0, time.UTC)
			if date.Weekday() == time.Friday {
				result.WriteString(fmt.Sprintf("Friday, %s 13, %d\n", j, i))
			}
		}
	}
	return result.String()
}

func ex05Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	if result := cargo.CargoTest(exercise, 1*time.Second, []string{}); !result.Passed {
		return result
	}
	if result := testExistenceTest(exercise); !result.Passed {
		return result
	}
	if err := testutils.AppendStringToFile(CargoTestEx05, exercise.TurnInFiles[0]); err != nil {
		return Exercise.InternalError(err.Error())
	}
	if result := cargo.CargoTest(exercise, 1*time.Second, []string{}); !result.Passed {
		return result
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	output, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"run"})
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	expectedOutput := expectedMainOutput()
	if output != expectedOutput {
		return Exercise.AssertionError(expectedOutput, output)
	}
	return Exercise.Passed("OK")
}

func ex05() Exercise.Exercise {
	return Exercise.NewExercise("05", "ex05", []string{"src/main.rs", "Cargo.toml"}, 15, ex05Test)
}
