package command

import "strings"

// Answer is the Oracle's reply plus the content paths used to ground it.
type Answer struct {
	Answer  string   `json:"answer"`
	Sources []string `json:"sources"`
}

// The Oracle's lore voice fails in character; the plain assistant reports
// a routine outage message.
const (
	oracleWithdrawn  = "The Oracle has wandered off into the smoke. Seek again shortly."
	assistantOffline = "The assistant is unavailable right now. Please try again shortly."
)

func formatAnswer(answer Answer) Result {
	output := answer.Answer
	if len(answer.Sources) > 0 {
		output += "\n\nsources: " + strings.Join(answer.Sources, ", ")
	}
	return Result{Output: output, ExitCode: ExitSuccess}
}
