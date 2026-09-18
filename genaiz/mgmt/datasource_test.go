package mgmt

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"genaiz.com/genaiz/task"
	"genaiz.com/genaiz/task/broker"
)

func TestUserDataSourcesFacade_Filtering(t *testing.T) {
	var expectedFilter = "filter"
	var testProvider = NewUserDataSourceFacade().
		WithLogger(logrus.New()).
		Filtering(expectedFilter)

	// We can't unit test the provider here
	assert.NotNil(t, testProvider)
}

func TestUserDataSourcesFacade_Provider(t *testing.T) {
	var testProvider = NewUserDataSourceFacade().
		WithParams(&broker.DataInstanceListParams{}).
		Provider()

	// We can't unit test the provider here
	assert.NotNil(t, testProvider)
}

func TestUserDataSourcesProvider_Get(t *testing.T) {
	var calledParams broker.DataInstanceListParams
	var testDatalinks = []broker.DataLinkInstance{
		{
			Id:    new(int64(37)),
			Name:  "expected",
			Flags: new(10),
		},
		{
			Id:      new(int64(42)),
			Created: 1,
		},
		{
			Id:      new(int64(73)),
			Created: 2,
		},
		{
			Id:      new(int64(69)),
			Created: 2,
		},
	}
	var testParams = &broker.DataInstanceListParams{}
	var testProvider = &userDataSourcesProvider{
		Plan: task.Plan{
			Logger: logrus.New(),
		},
		params:                    testParams,
		listDataSourceTaskFactory: newListDataSourceTaskCompleteCapture(&calledParams, testDatalinks),
	}
	var actual []UserLinkInstance
	var err error

	if actual, err = testProvider.Get(); err == nil {
		assert.Equal(t, 4, len(actual))
		assert.Equal(t, *testDatalinks[2].Id, *actual[0].Id)
		assert.Equal(t, *testDatalinks[3].Id, *actual[1].Id)
		assert.Equal(t, *testDatalinks[1].Id, *actual[2].Id)
		assert.Equal(t, *testDatalinks[0].Id, *actual[3].Id)
		return
	}

	assert.Fail(t, "expected a list of results")
}

func TestUserDataSourcesProvider_Get_SingleResult(t *testing.T) {
	var calledParams broker.DataInstanceListParams
	var testDatalinks = []broker.DataLinkInstance{
		{
			Id:    new(int64(37)),
			Name:  "expected",
			Flags: new(10),
		},
	}
	var testParams = &broker.DataInstanceListParams{}
	var testProvider = &userDataSourcesProvider{
		Plan: task.Plan{
			Logger: logrus.New(),
		},
		params:                    testParams,
		listDataSourceTaskFactory: newListDataSourceTaskCompleteCapture(&calledParams, testDatalinks),
	}
	var actual []UserLinkInstance
	var err error

	if actual, err = testProvider.Get(); err == nil {
		assert.Equal(t, 1, len(actual))
		assert.Equal(t, *testDatalinks[0].Id, *actual[0].Id)
		assert.Equal(t, testDatalinks[0].Name, actual[0].Name)
		assert.Equal(t, testDatalinks[0].Flags, actual[0].Flags)
		return
	}

	assert.Fail(t, "expected a list of results")
}

func TestUserDataSourcesProvider_Get_Failure(t *testing.T) {
	var expectedError = task.NewError("expected")
	var testParams = &broker.DataInstanceListParams{
		DataLink: &broker.DataLink{
			Oem: "oem",
		},
	}
	var testProvider = &userDataSourcesProvider{
		Plan: task.Plan{
			Logger: logrus.New(),
		},
		params:                    testParams,
		filter:                    "3",
		listDataSourceTaskFactory: newListDataSourceTaskCompleteError(expectedError),
	}
	var err error

	if _, err = testProvider.Get(); err != nil {
		assert.ErrorIs(t, err, expectedError)
		return
	}

	assert.Fail(t, "expected an error")
}

func newListDataSourceTaskCompleteCapture(capture *broker.DataInstanceListParams, seeded []broker.DataLinkInstance) func() *task.Task[broker.DataInstanceListParams] {
	return func() *task.Task[broker.DataInstanceListParams] {
		return &task.Task[broker.DataInstanceListParams]{
			OnPrepare: func(params *broker.DataInstanceListParams, state *task.State) error {
				return nil
			},
			OnComplete: func(params *broker.DataInstanceListParams, state *task.State) error {
				*capture = *params
				state.Internal = seeded
				return nil
			},
		}
	}
}

func newListDataSourceTaskCompleteError(err error) func() *task.Task[broker.DataInstanceListParams] {
	return func() *task.Task[broker.DataInstanceListParams] {
		return &task.Task[broker.DataInstanceListParams]{
			OnPrepare: func(params *broker.DataInstanceListParams, state *task.State) error {
				return nil
			},
			OnComplete: func(params *broker.DataInstanceListParams, state *task.State) error {
				return err
			},
		}
	}
}
