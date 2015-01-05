package example

import (
	"fmt"

	"github.com/progrium/go-extpoints/_examples/tool/extpoints"
	"github.com/progrium/go-extpoints/_examples/tool/types"
)

func init() {
	extpoints.Register(new(exampleExtension))
}

type exampleExtension struct{}

func (h *exampleExtension) Commands() []*types.Command {
	return []*types.Command{cmdExample}
}

func (h *exampleExtension) CommandStarts(commandName string) error {
	return nil
}

func (h *exampleExtension) CommandFinished(commandName string) {
	if commandName == "hello" {
		fmt.Println("Example extension says hello, too!")
	}
}

var cmdExample = &types.Command{
	Run:   runExample,
	Usage: "example",
	Short: "command from example extension",
	Long:  "command from example extension",
}

func runExample(cmd *types.Command, args []string) {
	fmt.Println("This command was added by the example extension!")
}
