package services

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"

	"github.com/C-dexTeam/codex-compiler/internal/domains"
	serviceErrors "github.com/C-dexTeam/codex-compiler/internal/errors"
	dto "github.com/C-dexTeam/codex-compiler/internal/http/dtos"
	"github.com/C-dexTeam/codex-compiler/pkg/file"
)

type runnerService struct {
	utilService IUtilService
	mainDir     string
	userCodeDir string
	binaryDir   string
}

func NewRunnerService(utilService IUtilService) *runnerService {
	return &runnerService{
		utilService: utilService,
		mainDir:     "users",
		userCodeDir: "users/usercode",
		binaryDir:   "users/binaries",
	}
}

func (s *runnerService) CreateFiles(userAuthID, defaultFileName string, chapter dto.QuestChapter, tests []dto.QuestTest) error {
	// For now we are checking the outputs of the code. Not need checks_template
	// checks := s.createChecks(chapter.CheckTmp, tests)
	// chapter.DockerTmp = strings.Replace(chapter.DockerTmp, "$checks$", checks, -1)
	chapter.DockerTmp = strings.Replace(chapter.DockerTmp, "$code$", chapter.UserCode, -1)
	chapter.DockerTmp = strings.Replace(chapter.DockerTmp, "$funcname$", chapter.FuncName, -1)

	codePath := s.generateUserCodePath(userAuthID, chapter.ChapterID, defaultFileName)

	err := os.WriteFile(codePath, []byte(chapter.DockerTmp), 0644)
	if err != nil {
		return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusInternalServerError, serviceErrors.ErrUploadingUserCode)
	}

	return nil
}

func (s *runnerService) CreateDirectories(userAuthID string) error {
	if err := file.CheckDir(s.mainDir); err != nil {
		if err := file.CreateDir(s.mainDir); err != nil {
			return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreateDirectoryError)
		}
	}
	if err := file.CheckDir(s.userCodeDir); err != nil {
		if err := file.CreateDir(s.userCodeDir); err != nil {
			return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreateDirectoryError)
		}
	}
	if err := file.CheckDir(s.userCodeDir + "/" + userAuthID); err != nil {
		if err := file.CreateDir(s.userCodeDir + "/" + userAuthID); err != nil {
			return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreateDirectoryError)
		}
	}
	if err := file.CheckDir(s.binaryDir); err != nil {
		if err := file.CreateDir(s.binaryDir); err != nil {
			return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreateDirectoryError)
		}
	}
	if err := file.CheckDir(s.binaryDir + "/" + userAuthID); err != nil {
		if err := file.CreateDir(s.binaryDir + "/" + userAuthID); err != nil {
			return serviceErrors.NewServiceErrorWithMessage(serviceErrors.StatusInternalServerError, serviceErrors.ErrCreateDirectoryError)
		}
	}

	return nil
}

func (s *runnerService) BuildCode(build, userAuthID, chapterID, defaultFileName string) domains.CodeResponse {
	binaryPath := s.generateUserBinaryPath(userAuthID)
	userCodePath := s.generateUserCodePath(userAuthID, chapterID, defaultFileName)
	buildCode := fmt.Sprintf(build, binaryPath, userCodePath)

	cmd := exec.Command("sh", "-c", buildCode)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return domains.NewCodeResponse(string(output), "", "", err, false, nil)
	}

	return domains.NewCodeResponse("", "", "", nil, true, nil)
}

func (s *runnerService) RunCode(userAuthID, chapterID, defaultFileName, run string, tests []dto.QuestTest) domains.CodeResponse {
	// If there is "" in numbers thats a string.
	userCodePath := s.generateUserCodePath(userAuthID, chapterID, defaultFileName)

	testCount := int(math.Ceil(float64(len(tests)) * 0.30))
	var correctTestsID []string // Doğru olan testlerin 3 de 1 ini döndüreceğim.
	var lastOutput string
	for i, test := range tests {
		fmt.Println(test, i)

		testOutput := s.createTestOutput(test)
		runCode := s.createRunWithTest(run, userAuthID, chapterID, defaultFileName, test)

		cmd := exec.Command("sh", "-c", runCode)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return domains.NewCodeResponse(string(output), "", "", err, false, nil)
		}
		cleanOutput := strings.TrimSpace(string(output))

		if testCount > i {
			correctTestsID = append(correctTestsID, test.TestID)
		}

		if testOutput != cleanOutput {
			return domains.NewCodeResponse("", cleanOutput, test.TestID, nil, false, correctTestsID)
		}

		lastOutput = string(output)
	}

	if len(tests) == 0 {
		runCode := fmt.Sprintf(run, userCodePath, "")
		cmd := exec.Command("sh", "-c", runCode)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return domains.NewCodeResponse(string(output), "", "", err, false, nil)
		}
		lastOutput = string(output)
	}

	return domains.NewCodeResponse("", lastOutput, "", nil, true, correctTestsID)
}

func (s *runnerService) createRunWithTest(run, userAuthID, chapterID, defaultFileName string, test dto.QuestTest) string {
	userCodePath := s.generateUserCodePath(userAuthID, chapterID, defaultFileName)

	var inputs string
	argsArr := strings.Split(strings.ReplaceAll(test.Input, "|", ""), " ")
	for _, str := range argsArr {
		inputs += str + " "
	}
	runCode := fmt.Sprintf(run, userCodePath, strings.TrimSpace(inputs))

	return runCode
}

func (s *runnerService) createTestOutput(test dto.QuestTest) string {
	var output string
	argsArr := strings.Split(strings.ReplaceAll(test.Output, "|", ""), " ")
	for _, str := range argsArr {
		output += str + " "
	}
	output = strings.TrimSpace(output)

	return output
}

func (s *runnerService) generateUserCodePath(userID, chapterID, defaultName string) string {
	extention := strings.Split(defaultName, ".")[1]

	userDir := s.userCodeDir + "/" + userID
	fileName := fmt.Sprintf("%v.%v", chapterID, extention)

	return userDir + "/" + fileName
}

func (s *runnerService) generateUserBinaryPath(userID string) string {
	userDir := s.binaryDir + "/" + userID

	return userDir
}
