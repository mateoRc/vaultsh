package command

import (
	"fmt"

	"github.com/mateom/vaultsh/internal/negotiation"
)

type AcceptOffer struct {
	name        string
	negotiation *negotiation.State
}

func NewAcceptOffer(name string, negotiation *negotiation.State) AcceptOffer {
	return AcceptOffer{
		name:        name,
		negotiation: negotiation,
	}
}

func (a AcceptOffer) Name() string {
	return a.name
}

func (AcceptOffer) Description() string {
	return "Accept a pending counter-offer"
}

func (AcceptOffer) Hidden() bool {
	return true
}

func (a AcceptOffer) Execute(args []string, _ Input) Result {
	if len(args) != 0 {
		return Result{
			Output:   fmt.Sprintf("usage: %s", a.name),
			ExitCode: ExitUsage,
		}
	}

	counterOffer, found := a.negotiation.Accept()
	if !found {
		return Result{
			Output:   fmt.Sprintf("%s: no pending counter-offer", a.name),
			ExitCode: ExitFailure,
		}
	}

	return Result{
		Output: fmt.Sprintf(
			"counter-offer accepted: %.2f\nwelcome aboard. paperwork has entered the chat.",
			counterOffer,
		),
		ExitCode: ExitSuccess,
	}
}
