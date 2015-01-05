//go:generate go-extpoints
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/progrium/go-extpoints/_examples/tool/extpoints"
	"github.com/progrium/go-extpoints/_examples/tool/types"
)

var (
	lifecycleContributors = extpoints.LifecycleContributors
	commandProviders      = extpoints.CommandProviders

	commands []*types.Command
)

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetFlags(0)

	for _, provider := range commandProviders.All() {
		commands = append(commands, provider.Commands()...)
	}

	// make sure command is specified, disallow global args
	args := os.Args[1:]
	if len(args) < 1 || strings.IndexRune(args[0], '-') == 0 {
		usage()
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() {
				cmd.PrintUsage()
			}
			if err := cmd.Flag.Parse(args[1:]); err != nil {
				os.Exit(2)
			}
			for _, contributor := range lifecycleContributors.All() {
				if err := contributor.CommandStarts(cmd.Name()); err != nil {
					os.Exit(3)
				}
			}
			cmd.Run(cmd, cmd.Flag.Args())
			for _, contributor := range lifecycleContributors.All() {
				contributor.CommandFinished(cmd.Name())
			}
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
	usage()
}
