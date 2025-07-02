package api

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/42-Short/shortinette/config"
	"github.com/42-Short/shortinette/dao"
	"github.com/42-Short/shortinette/git"
	"github.com/42-Short/shortinette/logger"
	"github.com/42-Short/shortinette/tester"
)

//todo: scheduler in api for repo creation

type moduleGrader struct {
	moduleDao      *dao.DAO[dao.Module]
	participantDao *dao.DAO[dao.Participant]
	ctx            context.Context
	config         config.Config
	gitService     *git.GithubService
}

func newModuleGrader(moduleDao *dao.DAO[dao.Module], participantDao *dao.DAO[dao.Participant], ctx context.Context, config config.Config) *moduleGrader {
	return &moduleGrader{
		moduleDao:      moduleDao,
		participantDao: participantDao,
		ctx:            ctx,
		config:         config,
		gitService:     git.NewGithubService(config.TokenGithub, config.OrgaGithub, "../"),
	}
}

func (mg *moduleGrader) process(intraLogin string, moduleId int) error {
	module, err := mg.moduleDao.Get(mg.ctx, moduleId, intraLogin)
	if err != nil {
		return err
	}
	participant, err := mg.participantDao.Get(mg.ctx, module.IntraLogin)
	if err != nil {
		return err
	}

	result, err := mg.grade(*module, *participant)
	if err != nil {
		return err
	}

	err = mg.updateModuleState(module, *result)
	if err != nil {
		return err
	}
	err = mg.updateParticipantState(participant, *result)
	if err != nil {
		return err
	}

	return nil
}

func (mg moduleGrader) isValidGradingAttempt(module dao.Module, participant dao.Participant) bool {
	remainingWaitTime := module.WaitTime - time.Since(module.LastGraded)
	if remainingWaitTime > 0 {
		logger.File.Printf("grading attempt too early. Please wait %s before trying again", remainingWaitTime)
		return false
	}

	if participant.CurrentModuleId < module.Id {
		logger.File.Printf("complete the previous module/modules before attempting this one.")
		return false
	}

	return true
}

func (mg *moduleGrader) updateModuleState(module *dao.Module, result tester.GradingResult) error {
	module.LastGraded = time.Now()
	module.WaitTime = 0
	// module.WaitTime = time.Duration(1<<module.Attempts) * time.Minute
	module.Attempts++
	module.Score = result.Score
	return mg.moduleDao.Update(mg.ctx, *module)
}

func (mg *moduleGrader) updateParticipantState(participant *dao.Participant, result tester.GradingResult) error {
	if !result.Passed {
		return nil
	}
	participant.CurrentModuleId++
	return mg.participantDao.Update(mg.ctx, *participant)
}

func (mg moduleGrader) grade(module dao.Module, participant dao.Participant) (*tester.GradingResult, error) {
	traceFile := filepath.Join("traces", fmt.Sprintf("%s%d_%s.log", module.IntraLogin, module.Id, time.Now().Format("20060102_150405")))
	if err := logger.InitializeTraceLogger(traceFile); err != nil {
		return nil, fmt.Errorf("trace logger could not be initialized: %v", err)
	}

	defer os.Remove(traceFile)

	if !mg.isValidGradingAttempt(module, participant) {
		return nil, fmt.Errorf("invalid grading attempt")
	}

	repoName := fmt.Sprintf("%s-%02d", module.IntraLogin, module.Id)
	if err := mg.gitService.Clone(repoName); err != nil {
		return nil, fmt.Errorf("could not clone repo '%s': %v", repoName, err)
	}

	defer func() {
		if err := os.RemoveAll(repoName); err != nil {
			logger.Error.Printf("could not tear down cloned repo: %v", err)
		}
	}()

	result, err := tester.GradeModule(mg.config.Modules[module.Id], repoName, "../testenv/Dockerfile")
	if err != nil {
		return nil, err
	}
	logger.File.Print(result.Trace)
	mg.uploadTraces(traceFile, module)
	return result, nil
}

func (mg moduleGrader) uploadTraces(traceFile string, module dao.Module) {
	commitMessage := fmt.Sprintf("chore: automated upload of trace logs for module %d (user: %s)", module.Id, module.IntraLogin)
	repoName := fmt.Sprintf("%s-%02d", module.IntraLogin, module.Id)

	if err := mg.gitService.UploadFiles(repoName, commitMessage, "traces", false, traceFile); err != nil {
		logger.Error.Printf("could not upload traces for user %s, module %d: %v", module.IntraLogin, module.Id, err)
	}
}
