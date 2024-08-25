package short

import (
	"os"
	"testing"
)

func TestUnsetConfigPath(t *testing.T) {
	os.Unsetenv("CONFIG_PATH")
	if _, err := GetConfig(); err == nil {
		t.Fatal("missing error when calling GetConfig with unset CONFIG_PATH")
	}
}
