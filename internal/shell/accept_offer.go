package shell

import (
	"fmt"

	"github.com/mateom/vaultsh/internal/command"
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

func (a AcceptOffer) Usage() string {
	return a.name
}

func (a AcceptOffer) Execute(args []string, _ command.Input) command.Result {
	if len(args) != 0 {
		return command.Result{
			Output:   fmt.Sprintf("usage: %s", a.name),
			ExitCode: command.ExitUsage,
		}
	}

	counterOffer, found := a.negotiation.Accept()
	if !found {
		return command.Result{
			Output:   fmt.Sprintf("%s: no pending counter-offer", a.name),
			ExitCode: command.ExitFailure,
		}
	}

	return command.Result{
		Output: fmt.Sprintf(
			"counter-offer accepted: %.2f\npaperwork has entered the chat.",
			counterOffer,
		),
		ExitCode: command.ExitSuccess,
	}
}
