package command

import "strings"

type AskService interface {
	Ask(question string) (Answer, error)
}

type Ask struct {
	service AskService
}

func NewAsk(service AskService) Ask {
	return Ask{service: service}
}

func (Ask) Name() string {
	return "ask"
}

func (Ask) Description() string {
	return "Ask the Oracle plainly about Mateo and Backend Lab"
}

func (Ask) Usage() string {
	return "ask <question>"
}

func (Ask) Help() string {
	return `Example: ask "which languages does Mateo use?"`
}

func (a Ask) Execute(args []string, _ Input) Result {
	if len(args) == 0 {
		return Result{Output: "usage: ask <question>", ExitCode: ExitUsage}
	}

	answer, err := a.service.Ask(strings.Join(args, " "))
	if err != nil {
		return Result{Output: assistantOffline, ExitCode: ExitFailure}
	}
	return formatAnswer(answer)
}
