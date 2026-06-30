package command

import (
	"fmt"
	"math"
	"strconv"

	"github.com/mateom/vaultsh/internal/negotiation"
)

const (
	salaryMultiplier = 1.5
	minimumSalary    = 70000
	suggestedSalary  = 100000
	immediateHire    = 100000
)

type Sudo struct {
	negotiation *negotiation.State
}

func NewSudo(negotiation *negotiation.State) Sudo {
	return Sudo{negotiation: negotiation}
}

func (Sudo) Name() string {
	return "sudo"
}

func (Sudo) Description() string {
	return "Run a command with elevated access"
}

func (Sudo) Hidden() bool {
	return true
}

func (s Sudo) Execute(args []string, _ Input) Result {
	if len(args) != 4 ||
		args[0] != "hire" ||
		args[1] != "mateo" ||
		args[2] != "-s" {
		return sudoAccessDenied()
	}

	salary, err := strconv.ParseFloat(args[3], 64)
	if err != nil || salary <= 0 || math.IsNaN(salary) || math.IsInf(salary, 0) {
		return Result{
			Output:   fmt.Sprintf("sudo: invalid salary: %s", args[3]),
			ExitCode: ExitUsage,
		}
	}
	if salary < minimumSalary {
		return Result{
			Output:   fmt.Sprintf("did you mean %.0f?", float64(suggestedSalary)),
			ExitCode: ExitFailure,
		}
	}
	if salary > immediateHire {
		return Result{
			Output:   "when do i start?",
			ExitCode: ExitSuccess,
		}
	}

	counterOffer := salary * salaryMultiplier
	s.negotiation.Propose(counterOffer)

	return Result{
		Output: fmt.Sprintf(
			"sudo: access granted\n"+
				"salary offered: %.2f\n"+
				"counter-offer: %.2f\n"+
				"accept counter-offer? [Y/y]",
			salary,
			counterOffer,
		),
		ExitCode: ExitSuccess,
	}
}

func sudoAccessDenied() Result {
	return Result{
		Output: "sudo: access denied\n" +
			"hint: only one privileged workflow is available\n" +
			"hint: try: sudo hire mateo -s <yearly>",
		ExitCode: ExitFailure,
	}
}
