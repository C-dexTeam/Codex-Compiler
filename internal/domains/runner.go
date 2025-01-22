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
	Correct    bool
	BuildError string
	Err        error
}

func NewCodeResponse(buildError string, err error, correct bool) error {
	return &CodeResponse{
		Correct:    correct,
		BuildError: buildError,
		Err:        err,
	}
}

func (e *CodeResponse) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return ""
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
			"go run main.go",     // Go derlemeden çalıştırma
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
