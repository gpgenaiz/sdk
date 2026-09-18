package mgmt

import (
	"cmp"
	"slices"

	"github.com/sirupsen/logrus"

	"genaiz.com/genaiz/task"
	"genaiz.com/genaiz/task/broker"
)

type UserDataSourceFacade Facade[[]UserLinkInstance, broker.DataInstanceListParams]

type ListDataSourceTaskFactory func() *task.Task[broker.DataInstanceListParams]

type userDataSourcesFacade struct {
	baseLoggingFacade
	params *broker.DataInstanceListParams
}

func (udf userDataSourcesFacade) Filtering(filter string) Provider[[]UserLinkInstance] {
	return &userDataSourcesProvider{
		Plan: task.Plan{
			Logger: udf.logger,
		},
		filter:                    filter,
		params:                    udf.params,
		listDataSourceTaskFactory: broker.NewDataSourceListTask,
	}
}

func (udf userDataSourcesFacade) Provider() Provider[[]UserLinkInstance] {
	return udf.Filtering("")
}

func (udf userDataSourcesFacade) WithLogger(logger *logrus.Logger) Facade[[]UserLinkInstance, broker.DataInstanceListParams] {
	udf.logger = logger
	return udf
}

func (udf userDataSourcesFacade) WithParams(params *broker.DataInstanceListParams) Facade[[]UserLinkInstance, broker.DataInstanceListParams] {
	udf.params = params
	return udf
}

type userDataSourcesProvider struct {
	task.Plan
	filter                    string
	params                    *broker.DataInstanceListParams
	listDataSourceTaskFactory ListDataSourceTaskFactory
}

func (udp *userDataSourcesProvider) Get() ([]UserLinkInstance, task.Error) {
	var dataSources []broker.DataLinkInstance
	var workers []task.Worker
	var failure interface{}

	udp.OnReturn = func(i interface{}) { dataSources = i.([]broker.DataLinkInstance) }
	udp.OnFailure = func(i interface{}) { failure = i }
	workers = append(workers, task.NewWorker(udp.params, udp.listDataSourceTaskFactory()))
	udp.Sequence(workers...)

	if failure == nil {
		var result = make([]UserLinkInstance, 0)

		for _, ds := range dataSources {
			var dataSource = ToUserLinkInstance(&ds)

			result = append(result, *dataSource)
		}

		if len(result) > 1 {
			slices.SortFunc(result, func(a, b UserLinkInstance) int {
				return cmp.Compare(b.Created, a.Created)
			})
		}

		return result, nil
	}

	return nil, task.NewFailure(failure)
}

func NewUserDataSourceFacade() UserDataSourceFacade {
	return &userDataSourcesFacade{
		baseLoggingFacade: baseLoggingFacade{},
	}
}
