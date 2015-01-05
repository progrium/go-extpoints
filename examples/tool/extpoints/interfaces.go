package extpoints

import (
	"github.com/progrium/go-extpoints/_examples/tool/types"
)

type LifecycleContributor interface {
	CommandStarts(commandName string) error
	CommandFinished(commandName string)
}

type CommandProvider interface {
	Commands() []*types.Command
}
