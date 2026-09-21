package dt

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"genaiz.com/genaiz/cli"
	"genaiz.com/genaiz/cmd/dk"
	"genaiz.com/genaiz/cmd/dt/source"
	"genaiz.com/genaiz/config"
	"genaiz.com/genaiz/mgmt"
	"genaiz.com/genaiz/schema"
	"genaiz.com/genaiz/task"
	"genaiz.com/genaiz/task/broker"
)

var (
	errorSourceSeqInvalid = task.NewError("sequence number is invalid")
)

type SourceListExecutor struct {
	*SourceListOptions
	ledger *config.Ledger

	accountParams                config.AccountParametric
	printerParams                cli.PrinterParametric
	userDataSourceFacadeProvider func() mgmt.UserDataSourceFacade
}

func (sle SourceListExecutor) List(arg string) error {
	var listParams *broker.DataInstanceListParams
	var err error

	if listParams, err = sle.newSourceListParams(arg); err == nil {
		var dataSourceList = sle.userDataSourceFacadeProvider()
		var printer = sle.printerParams.Printer()
		var sources []mgmt.UserLinkInstance

		sources, err = dataSourceList.WithParams(listParams).
			WithLogger(sle.ledger.Logger).
			Provider().
			Get()

		if err == nil {
			return printer.Print(sources)
		}

		return printer.Error(err)
	}

	return err
}

func (sle SourceListExecutor) newSourceListParams(arg string) (*broker.DataInstanceListParams, error) {
	var oem, handle, vn = dk.ParseDataLinkArgument(arg)
	var result = &broker.DataInstanceListParams{
		Broker: *sle.accountParams.BrokerParams(),
	}
	var sequence *int

	if vn != "" {
		if ts := strings.SplitN(vn, "-rc-", 2); len(ts) > 1 {
			var err error
			var ti int

			if ti, err = strconv.Atoi(ts[1]); err != nil {
				return nil, errorSourceSeqInvalid
			}

			vn = ts[0]
			sequence = &ti
		}
	}

	if oem != "" {
		result.DataLink = &broker.DataLink{
			Oem:     oem,
			Handle:  handle,
			Version: vn,
			Seq:     sequence,
		}
	}

	return result, nil
}

type SourceListOptions struct {
	optionAccount     *config.StringOption
	optionJsonPrinter *config.BoolOption
}

func (slo SourceListOptions) allDefiners() []config.Definer {
	return []config.Definer{
		slo.optionAccount,
		slo.optionJsonPrinter,
	}
}

func NewSource(ledger *config.Ledger) *cobra.Command {
	var listOptions = NewSourceListOptions()
	var listCmd = source.NewListSource(newSourceListExecutorFactory(ledger, listOptions))
	var srcCmd = &cobra.Command{
		Use:     "source",
		Aliases: []string{"src"},
		Short:   "Manages data source for an account",
	}

	srcCmd.AddCommand(listCmd)
	ledger.Register(listCmd, listOptions.allDefiners()...)
	cli.AutoBridge.Accounts().Option(listCmd, ledger, listOptions.optionAccount)
	return srcCmd
}

func NewSourceListExecutor(ledger *config.Ledger, options *SourceListOptions) *SourceListExecutor {
	return &SourceListExecutor{
		SourceListOptions: options,
		ledger:            ledger,

		accountParams:                config.NewAccountParams(ledger, options.optionAccount),
		printerParams:                cli.NewPrinterParam(ledger, options.optionJsonPrinter),
		userDataSourceFacadeProvider: mgmt.NewUserDataSourceFacade,
	}
}

func NewSourceListOptions() *SourceListOptions {
	return &SourceListOptions{
		optionAccount: cli.Options.Sources.Account().
			WithKeys(&schema.Genaiz.Data.Source.List.Account).
			BuildStringOption(),
		optionJsonPrinter: cli.Options.Printer.JsonPrinter().
			WithKeys(&schema.Genaiz.Data.Source.List.Printer).
			BuildBoolOption(),
	}
}

func newSourceListExecutorFactory(ledger *config.Ledger, options *SourceListOptions) source.ListExecutorFactory {
	return func(command *cobra.Command) source.ListExecutor {
		return NewSourceListExecutor(ledger, options)
	}
}
