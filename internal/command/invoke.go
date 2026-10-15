package command

import "strings"

type InvokeService interface {
	Invoke(question string) (Answer, error)
}

type Invoke struct {
	service InvokeService
}

func NewInvoke(service InvokeService) Invoke {
	return Invoke{service: service}
}

func (Invoke) Name() string {
	return "invoke"
}

func (Invoke) Description() string {
	return "Invoke the Oracle's lore voice about Mateo and Backend Lab"
}

func (Invoke) Usage() string {
	return "invoke <question>"
}

func (Invoke) Help() string {
	return `Example: invoke "how does Atlas search work?"`
}

func (i Invoke) Execute(args []string, _ Input) Result {
	if len(args) == 0 {
		return Result{Output: "usage: invoke <question>", ExitCode: ExitUsage}
	}

	answer, err := i.service.Invoke(strings.Join(args, " "))
	if err != nil {
		return Result{Output: oracleWithdrawn, ExitCode: ExitFailure}
	}
	return formatAnswer(answer)
}
