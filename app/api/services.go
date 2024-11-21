package api

import "github.com/42-Short/shortinette/logger"

func processGrading(intra_login string, moduleId int) {
	logger.Info.Printf("grading %s module %d...\n", intra_login, moduleId)
}
