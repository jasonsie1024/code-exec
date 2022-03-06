package isolate

import (
	"os"
	"strings"
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
		dict["status"] = "Accepted"
	case "RE":
		dict["status"] = "Runtime Error"
	case "TO":
		dict["status"] = "Time Limit Exceed"
	default:
		dict["status"] = "Internal Error"
	}

	return dict
}
