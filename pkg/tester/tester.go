package shortinette

import (
    "os/exec"
    "bytes"
    "log"
)

func RunUserCode(scriptPath string) (string, error) {
    cmd := exec.Command("python", scriptPath)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
        return "", err
    }
    return out.String(), nil
}
