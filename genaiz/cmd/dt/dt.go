package dt

import (
	"context"

	"github.com/spf13/cobra"

	"genaiz.com/genaiz/cli"
	"genaiz.com/genaiz/config"
)

type BaseExecutor struct {
	Cli     *Cli
	Context context.Context
	Ledger  *config.Ledger
}

type Cli struct {
	cli.BaseCli
}

func NewDt(ledger *config.Ledger) *cobra.Command {
	var dtCmd = &cobra.Command{
		Use:     "data",
		Aliases: []string{"dt"},
		Short:   "GenAIz Data Instance Toolkit",
	}

	dtCmd.AddCommand(NewSource(ledger))
	return dtCmd
}
