package domains

import "fmt"

type Language struct {
	Name        string
	Run         string
	Execute     string
	Build       string
	DefaultName string
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
			"go build -o main .", // Go derleme komutu
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
