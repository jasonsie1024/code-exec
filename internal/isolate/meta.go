package isolate

import (
	"os"
	"strings"

	"github.com/jason-plainlog/code-exec/internal/models"
)

func ParseMetafile(metafile string) map[string]string {
	dict := map[string]string{}

	meta, err := os.ReadFile(metafile)
	if err != nil {
		return dict
	}

	lines := strings.Split(string(meta), "\n")
	for _, line := range lines {
		pair := strings.SplitN(line, ":", 2)
		if len(pair) == 2 {
			dict[pair[0]] = pair[1]
		}
	}

	switch dict["status"] {
	case "":
		dict["status"] = string(models.Accepted)
	case "RE":
		dict["status"] = string(models.RuntimeError)
	case "TO":
		dict["status"] = string(models.TimeLimitExceeded)
	default:
		dict["status"] = string(models.InternalError)
	}

	return dict
}
