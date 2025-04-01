package domains

import "fmt"

type Language struct {
	Name        string
	Run         string
	Execute     string
	Build       string
	DefaultName string
}

type CodeResponse struct {
	Correct        bool
	Output         string
	BuildError     string
	Err            string
	CorrectTestsID []string
	WrongTestID    string
}

func NewCodeResponse(buildError, output, wrongTestID string, err error, correct bool, correctTestsID []string) CodeResponse {
	var errStr string
	if err != nil {
		errStr = err.Error()
	} else {
		errStr = ""
	}

	return CodeResponse{
		WrongTestID:    wrongTestID,
		CorrectTestsID: correctTestsID,
		Correct:        correct,
		BuildError:     buildError,
		Output:         output,
		Err:            errStr,
	}
}

func newLanguage(name, run, build, defaultName string) Language {
	return Language{
		Name:        name,
		Execute:     fmt.Sprintf("./%v", name),
		Build:       build,
		Run:         run,
		DefaultName: defaultName,
	}
}

func Languages() []Language {
	return []Language{
		newLanguage(
			"Go",
			"go run %v %v",       // Go derlemeden çalıştırma
			"go build -o %v %v ", // Go derleme komutu
			"main.go",            // Varsayılan dosya adı
		),

		newLanguage(
			"Rust",
			"cargo run",             // Rust derlemeden çalıştırma
			"cargo build --release", // Rust derleme komutu
			"main.rs",               // Varsayılan dosya adı
		),
	}
}

func GetLanguage(name string) *Language {
	langs := Languages()
	for _, lang := range langs {
		if lang.Name == name {
			return &lang
		}
	}

	return nil
}

func GetLanguagesName() []string {
	var names []string
	langs := Languages()
	for _, lang := range langs {
		names = append(names, lang.Name)
	}

	return names
}
