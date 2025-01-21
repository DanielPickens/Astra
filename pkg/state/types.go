package state

import (
	"github\.com/danielpickens/astra/pkg/api"
)

type Content struct {
	// PID is the ID of the process to which the state belongs
	PID int `json:"pid"`
	// Platform indicates on which platform the session works
	Platform string `json:"platform"`
	// ForwardedPorts are the ports forwarded during astra dev session
	ForwardedPorts []api.ForwardedPort `json:"forwardedPorts"`
	APIServerPort  int                 `json:"apiServerPort,omitempty"`
}
