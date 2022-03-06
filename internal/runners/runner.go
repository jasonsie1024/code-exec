package runners

import "github.com/jason-plainlog/code-exec/internal/models"

type (
	IRunner interface {
		Handle(*models.Submission) chan *models.Result

		Compile(*models.Submission) *models.Result
		Execute(*models.Submission, *models.Task) *models.Result
	}
)

var Runners = map[string]IRunner{
	"python3": &python3Runner{},
}
