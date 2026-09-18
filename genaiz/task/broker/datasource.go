package broker

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"

	"genaiz.com/genaiz-lib/lang/slicez"
	"genaiz.com/genaiz/task"
	ver "genaiz.com/genaiz/version"
)

func NewDataSourceListTask() *task.Task[DataInstanceListParams] {
	return &task.Task[DataInstanceListParams]{
		Name:       "data-source-list",
		OnPrepare:  handleDataSourceListContext,
		OnComplete: handleDataSourceListComplete,
		OnPretend:  handleDataSourceListPretend,
	}
}

func handleDataSourceListComplete(params *DataInstanceListParams, state *task.State) error {
	var err error
	var brokerClient Client

	if brokerClient, err = params.GetClient(); err == nil {
		var sources []DataLinkInstance

		handleDataSourceFilterDebugging(params, state.Logger)

		if sources, err = brokerClient.ListDataSources(); err == nil {
			if params.hasDataLinkFilter() {
				var filtered = slicez.Filter(sources, func(instance DataLinkInstance) bool {
					if instance.HasLink(params.DataLink) {
						return params.Seq == nil || (*params.Seq == *instance.DataLink.Seq)
					}

					// Case where we only filter by oem
					if strings.EqualFold(params.Oem, instance.DataLink.Oem) {
						if params.Handle != "" {
							// Case where we filter by oem and handle
							return strings.EqualFold(params.Handle, instance.DataLink.Handle)
						}

						// If params had a version, HasLink would return true
						return true
					}

					return false
				})

				state.Internal = filtered
			} else {
				state.Internal = sources
			}

			return nil
		}
	}

	return err
}

func handleDataSourceListContext(params *DataInstanceListParams, state *task.State) error {
	if state.Output == "" {
		if err := params.validateDataLinkFilter(); err != nil {
			return err
		}

		if id := params.getLinkId(); id != nil {
			state.Logger.Debugf("Listing data sources for data link [%d]", *id)
			state.Logger.Warnf("Filtering by id is not supported in version [%s]", ver.GetVersion())
			state.Output = cast.ToString(*id)
		}
	}

	return nil
}

func handleDataSourceListPretend(params *DataInstanceListParams, state *task.State) error {
	var brokerClient Client
	var err error

	if brokerClient, err = params.GetClient(); err == nil {
		state.Logger.Debugf("Pretending to list data source for account [%s]", brokerClient.GetHostAddr())
		handleDataSourceFilterDebugging(params, state.Logger)
		fmt.Printf("curl -X GET -H \"Content-Type: application/x-www-form-urlencoded\" \\\n")
		fmt.Printf("  --cookie=\"s=%s\"\\\n", brokerClient.GetAuthToken())

		if id := params.getLinkId(); id != nil {
			fmt.Printf("  -G -d id=%d\\\n", *id)
		}

		fmt.Printf("%s\n", brokerClient.ListDataSourcesUrl())
	}

	return err
}

func handleDataSourceFilterDebugging(params *DataInstanceListParams, logger *logrus.Logger) {
	if params.hasDataLinkFilter() {
		logger.Debugf("Filtering data sources for oem [%s]", params.Oem)

		if params.Handle != "" {
			logger.Debugf("With handle [%s]", params.Handle)
		}

		if params.Version != "" {
			logger.Debugf("And version [%s]", params.DataLink.GetVersion())
		}
	}
}
