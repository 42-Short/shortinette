package short

import (
	"fmt"

	"github.com/42-Short/shortinette/internal/logger"
)

type HourlyTestMode struct {
	Delay             int
	FrequenzyDuration int
}

type MainBranchTestMode struct {
}

// TODO: find better name
type TestMode struct {
	Hourly     *HourlyTestMode
	MainBranch *MainBranchTestMode
}

func NewHourlyTestMode(hourly *HourlyTestMode) TestMode {
	return TestMode{
		Hourly:     hourly,
		MainBranch: nil,
	}
}
func NewMainBranchTestMode(mainBranch *MainBranchTestMode) TestMode {
	return TestMode{
		Hourly:     nil,
		MainBranch: mainBranch,
	}
}

// type SubjectSupplyMode enum {

// }

type Short struct {
	Name     string
	TestMode TestMode
	// Is it one per day or are users automatically assigne if they have the previous Subject at XX% --SubjectSupplyMode
	// modules [excercises]
}

func Run() {
	config, err := getConfig()
	if err != nil {
		logger.Error.Printf("internal error: %v", err)
	}
	
	fmt.Println(config)
}
