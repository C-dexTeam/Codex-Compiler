package services

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/C-dexTeam/codex-compiler/internal/domains"
	serviceErrors "github.com/C-dexTeam/codex-compiler/internal/errors"
	dto "github.com/C-dexTeam/codex-compiler/internal/http/dtos"
	"github.com/C-dexTeam/codex-compiler/pkg/file"
	typechecker "github.com/C-dexTeam/codex-compiler/pkg/type_checker"
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
	checks := s.createChecks(chapter.CheckTmp, tests)
	chapter.DockerTmp = strings.Replace(chapter.DockerTmp, "$checks$", checks, -1)
	chapter.DockerTmp = strings.Replace(chapter.DockerTmp, "$res$", fmt.Sprint(len(tests)-1), -1)
	chapter.DockerTmp = strings.Replace(chapter.DockerTmp, "$usercode$", chapter.UserCode, -1)
	chapter.DockerTmp = strings.Replace(chapter.DockerTmp, "$funcname$", chapter.FuncName, -1)
	chapter.DockerTmp = strings.Replace(chapter.DockerTmp, "$success$", "Test Passed", -1)

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

func (s *runnerService) BuildCode(build, userAuthID, ChapterID, defaultFileName string) error {
	fmt.Println("Building Code")

	binaryPath := s.generateUserBinaryPath(userAuthID)
	userCodePath := s.generateUserCodePath(userAuthID, ChapterID, defaultFileName)
	buildCode := fmt.Sprintf(build, binaryPath, userCodePath)

	cmd := exec.Command("sh", "-c", buildCode)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return domains.NewCodeResponse(string(output), err, false)
	}

	return domains.NewCodeResponse("", nil, true)
}

func (s *runnerService) RunCode(name string, tests []dto.QuestTest) {
	// fmt.Println("RunCode Runned")
}

func (s *runnerService) createChecks(check string, tests []dto.QuestTest) string {
	var checks strings.Builder

	for i, test := range tests {
		tmp := check
		tmp = strings.Replace(tmp, "$rnd$", fmt.Sprintf("%v", i), -1)

		// Split input and output by "|"
		inputs := strings.Split(test.Input, "|")
		outputs := strings.Split(test.Output, "|")

		var inputValues []string
		for _, in := range inputs {
			in = strings.TrimSpace(in)
			// Identify type of input element
			if typechecker.IsString(in) {
				inputValues = append(inputValues, in) // Already string so we can append it normaly
			} else if typechecker.IsBool(in) {
				inputValues = append(inputValues, fmt.Sprintf("%v", typechecker.ParseBool(in)))
			} else if typechecker.IsNumber(in) {
				inputValues = append(inputValues, fmt.Sprintf("%v", typechecker.ParseNumber(in)))
			} else {
				inputValues = append(inputValues, fmt.Sprintf("%v", in))
			}
		}
		tmp = strings.Replace(tmp, "$input$", strings.Join(inputValues, ", "), -1)

		// Handle output replacement
		var outputValues []string
		var fails []string
		for _, out := range outputs {
			out = strings.TrimSpace(out)
			// Identify type of output element
			if typechecker.IsString(out) {
				outputValues = append(outputValues, out)
				fails = append(fails, fmt.Sprintf("%v", out))
			} else if typechecker.IsBool(out) {
				outputValues = append(outputValues, fmt.Sprintf("%v", typechecker.ParseBool(out)))
				fails = append(fails, fmt.Sprintf("%v", typechecker.ParseBool(out)))
			} else if typechecker.IsNumber(out) {
				outputValues = append(outputValues, fmt.Sprintf("%v", typechecker.ParseNumber(out)))
				fails = append(fails, fmt.Sprintf("%v", typechecker.ParseNumber(out)))
			} else {
				outputValues = append(outputValues, fmt.Sprintf("%v", out))
				fails = append(fails, fmt.Sprintf("%v", out))
			}
		}
		tmp = strings.Replace(tmp, "$output$", strings.Join(outputValues, ", "), -1)

		// Handle $out$ if exists
		if strings.Contains(check, "$out$") {
			tmp = strings.Replace(tmp, "$out$", strings.Join(fails, ", "), -1)
		}

		checks.WriteString(tmp + "\n")
	}

	return checks.String()
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
