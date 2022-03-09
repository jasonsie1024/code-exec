package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type Language struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	SourceFile     string `json:"source_file"`
	CompileCommand string `json:"compile_command"`
	RunCommand     string `json:"run_command"`
}

var languages map[int]Language

func loadLanguages() map[int]Language {
	ents, err := os.ReadDir("languages")
	if err != nil {
		log.Fatalln(err)
	}

	languages = map[int]Language{}

	for _, ent := range ents {
		if !ent.Type().IsRegular() || !strings.HasSuffix(ent.Name(), ".json") {
			continue
		}

		content, err := os.ReadFile("languages/" + ent.Name())
		if err != nil {
			log.Fatalln(err)
		}

		var lang Language
		if err := json.Unmarshal(content, &lang); err != nil {
			log.Fatalln(err)
		}

		if l, exist := languages[lang.Id]; exist {
			log.Fatalf("language id collision: %s and %s", lang.Name, l.Name)
		}

		languages[lang.Id] = lang
	}

	return languages
}

func GetLanguages() map[int]Language {
	if languages == nil {
		languages = loadLanguages()
	}

	return languages
}
