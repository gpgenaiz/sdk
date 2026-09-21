package broker

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"

	"genaiz.com/genaiz/task"
)

func TestNewDataSourceListTask(t *testing.T) {
	var testTask = NewDataSourceListTask()

	assert.NotEmpty(t, testTask.Name)
	assert.NotNil(t, testTask.OnPrepare)
	assert.NotNil(t, testTask.OnComplete)
	assert.NotNil(t, testTask.OnPretend)
}

func Test_handleDataSourceListComplete(t *testing.T) {
	var expectedInstance = []DataLinkInstance{
		{
			Id: new(int64(37)),
		},
	}
	var testParams = &DataInstanceListParams{
		Broker: Broker{
			AuthFile: "file",
			HostAddr: "hostAddr",
		},
	}
	var testState = &task.State{}
	var restoredFactory = clientFactory.Get

	defer func() {
		clientFactory.Get = restoredFactory
	}()
	clientFactory.Get = func(authFile, addr string) (Client, error) {
		return &stubDataShareClient{
			listDataSources: expectedInstance,
		}, nil
	}

	assert.NoError(t, handleDataSourceListComplete(testParams, testState))
	actual, ok := testState.Internal.([]DataLinkInstance)
	assert.True(t, ok)
	assert.Equal(t, expectedInstance, actual)
}

func Test_handleDataSourceListComplete_FilterOem(t *testing.T) {
	var expectedOem = "expectedOem"
	var expectedInstances = []DataLinkInstance{
		{
			Id: new(int64(37)),
			DataLink: DataLink{
				Oem: expectedOem,
			},
		},
		{
			Id: new(int64(42)),
			DataLink: DataLink{
				Oem: "notExpected",
			},
		},
	}
	var testParams = &DataInstanceListParams{
		Broker: Broker{
			AuthFile: "file",
			HostAddr: "hostAddr",
		},
		DataLink: &DataLink{
			Oem: expectedOem,
		},
	}
	var testState = &task.State{
		Logger: logrus.New(),
	}
	var restoredFactory = clientFactory.Get

	defer func() {
		clientFactory.Get = restoredFactory
	}()
	clientFactory.Get = func(authFile, addr string) (Client, error) {
		return &stubDataShareClient{
			listDataSources: expectedInstances,
		}, nil
	}

	assert.NoError(t, handleDataSourceListComplete(testParams, testState))
	actual, ok := testState.Internal.([]DataLinkInstance)
	assert.True(t, ok)
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, expectedInstances[0], actual[0])
}

func Test_handleDataSourceListComplete_FilterOemHandle(t *testing.T) {
	var expectedOem = "expectedOem"
	var expectedHandle = "expectedHandle"
	var expectedInstances = []DataLinkInstance{
		{
			Id: new(int64(37)),
			DataLink: DataLink{
				Oem:    expectedOem,
				Handle: "notExpected",
			},
		},
		{
			Id: new(int64(42)),
			DataLink: DataLink{
				Oem: "notExpected",
			},
		},
		{
			Id: new(int64(69)),
			DataLink: DataLink{
				Oem:    expectedOem,
				Handle: expectedHandle,
			},
		},
	}
	var testParams = &DataInstanceListParams{
		Broker: Broker{
			AuthFile: "file",
			HostAddr: "hostAddr",
		},
		DataLink: &DataLink{
			Oem:    expectedOem,
			Handle: expectedHandle,
		},
	}
	var testState = &task.State{
		Logger: logrus.New(),
	}
	var restoredFactory = clientFactory.Get

	defer func() {
		clientFactory.Get = restoredFactory
	}()
	clientFactory.Get = func(authFile, addr string) (Client, error) {
		return &stubDataShareClient{
			listDataSources: expectedInstances,
		}, nil
	}

	assert.NoError(t, handleDataSourceListComplete(testParams, testState))
	actual, ok := testState.Internal.([]DataLinkInstance)
	assert.True(t, ok)
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, expectedInstances[2], actual[0])
}

func Test_handleDataSourceListComplete_FilterOemHandleVersion(t *testing.T) {
	var expectedOem = "expectedOem"
	var expectedHandle = "expectedHandle"
	var expectedVersion = "expectedVersion"
	var expectedSequence = 1
	var expectedInstances = []DataLinkInstance{
		{
			Id: new(int64(37)),
			DataLink: DataLink{
				Oem:    expectedOem,
				Handle: "notExpected",
			},
		},
		{
			Id: new(int64(42)),
			DataLink: DataLink{
				Oem: "notExpected",
			},
		},
		{
			Id: new(int64(69)),
			DataLink: DataLink{
				Oem:     expectedOem,
				Handle:  expectedHandle,
				Version: expectedVersion,
				Seq:     &expectedSequence,
			},
		},
	}
	var testParams = &DataInstanceListParams{
		Broker: Broker{
			AuthFile: "file",
			HostAddr: "hostAddr",
		},
		DataLink: &DataLink{
			Oem:    expectedOem,
			Handle: expectedHandle,
		},
	}
	var testState = &task.State{
		Logger: logrus.New(),
	}
	var restoredFactory = clientFactory.Get

	defer func() {
		clientFactory.Get = restoredFactory
	}()
	clientFactory.Get = func(authFile, addr string) (Client, error) {
		return &stubDataShareClient{
			listDataSources: expectedInstances,
		}, nil
	}

	assert.NoError(t, handleDataSourceListComplete(testParams, testState))
	actual, ok := testState.Internal.([]DataLinkInstance)
	assert.True(t, ok)
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, expectedInstances[2], actual[0])
}

func Test_handleDataSourceListComplete_ListError(t *testing.T) {
	var expectedError = errors.New("expected")
	var testParams = &DataInstanceListParams{
		Broker: Broker{
			AuthFile: "file",
			HostAddr: "hostAddr",
		},
	}
	var testState = &task.State{}
	var restoredFactory = clientFactory.Get

	defer func() {
		clientFactory.Get = restoredFactory
	}()
	clientFactory.Get = func(authFile, addr string) (Client, error) {
		return &stubDataShareClient{
			listDataSourcesError: expectedError,
		}, nil
	}

	assert.ErrorIs(t, handleDataSourceListComplete(testParams, testState), expectedError)
	assert.Nil(t, testState.Internal)
}

func Test_handleDataSourceListComplete_SessionError(t *testing.T) {
	var expectedError = errors.New("expected")
	var testParams = &DataInstanceListParams{
		Broker: Broker{
			AuthFile: "file",
			HostAddr: "hostAddr",
		},
	}
	var restoredFactory = clientFactory.Get

	defer func() {
		clientFactory.Get = restoredFactory
	}()
	clientFactory.Get = func(authFile, addr string) (Client, error) {
		return nil, expectedError
	}

	assert.ErrorIs(t, handleDataSourceListComplete(testParams, &task.State{}), expectedError)
}

func Test_handleDataSourceListContext(t *testing.T) {
	assert.NoError(t, handleDataSourceListContext(&DataInstanceListParams{}, &task.State{}))
}

func Test_handleDataSourceListContext_InvalidFilter(t *testing.T) {
	var testParams = &DataInstanceListParams{
		DataLink: &DataLink{
			Handle: "handleWithNoOem",
		},
	}

	assert.Error(t, handleDataSourceListContext(testParams, &task.State{}), errorDataShareOemRequired)
}

func Test_handleDataSourceListContext_OutputKnown(t *testing.T) {
	var testState = &task.State{
		Output: "known",
	}

	assert.NoError(t, handleDataSourceListContext(&DataInstanceListParams{}, testState))
}

func Test_handleDataSourceListContext_WithLink(t *testing.T) {
	var testState = &task.State{
		Logger: logrus.New(),
	}
	var testParams = &DataInstanceListParams{
		DataLink: &DataLink{
			Id:    new(int64(37)),
			Flags: new(DataLinkFlags.Active),
		},
	}

	assert.NoError(t, handleDataSourceListContext(testParams, testState))
	assert.Equal(t, cast.ToString(*testParams.DataLink.Id), testState.Output)
}

func Test_handleDataSourceListPretend(t *testing.T) {
	var testParams = &DataInstanceListParams{
		Broker: Broker{
			AuthFile: "file",
			HostAddr: "hostAddr",
		},
	}
	var testState = &task.State{
		Logger: logrus.New(),
	}
	var restoredFactory = clientFactory.Get
	var stdoutRestore = os.Stdout
	var r, w, _ = os.Pipe()

	os.Stdout = w

	defer func() {
		clientFactory.Get = restoredFactory
		os.Stdout = stdoutRestore

	}()
	clientFactory.Get = func(authFile, addr string) (Client, error) {
		return &stubDataShareClient{
			client: client{
				HostAddr: testParams.Broker.HostAddr,
			},
		}, nil
	}

	assert.NoError(t, handleDataSourceListPretend(testParams, testState))

	_ = w.Close()
	b, _ := io.ReadAll(r)
	output := string(b)
	assert.Contains(t, output, testParams.Broker.HostAddr)
}

func Test_handleDataSourceListPretend_WithId(t *testing.T) {
	var testParams = &DataInstanceListParams{
		Broker: Broker{
			AuthFile: "file",
			HostAddr: "hostAddr",
		},
		DataLink: &DataLink{
			Id: new(int64(37)),
		},
	}
	var testState = &task.State{
		Logger: logrus.New(),
	}
	var restoredFactory = clientFactory.Get
	var stdoutRestore = os.Stdout
	var r, w, _ = os.Pipe()

	os.Stdout = w

	defer func() {
		clientFactory.Get = restoredFactory
		os.Stdout = stdoutRestore

	}()
	clientFactory.Get = func(authFile, addr string) (Client, error) {
		return &stubDataShareClient{
			client: client{
				HostAddr: testParams.Broker.HostAddr,
			},
		}, nil
	}

	assert.NoError(t, handleDataSourceListPretend(testParams, testState))

	_ = w.Close()
	b, _ := io.ReadAll(r)
	output := string(b)
	assert.Contains(t, output, testParams.Broker.HostAddr)
	assert.Contains(t, output, cast.ToString(*testParams.DataLink.Id))
}

func Test_handleDataSourceListPretend_SessionError(t *testing.T) {
	var expectedError = errors.New("expected")
	var testParams = &DataInstanceListParams{
		Broker: Broker{
			AuthFile: "file",
			HostAddr: "hostAddr",
		},
	}
	var restoredFactory = clientFactory.Get

	defer func() {
		clientFactory.Get = restoredFactory
	}()
	clientFactory.Get = func(authFile, addr string) (Client, error) {
		return nil, expectedError
	}

	assert.ErrorIs(t, handleDataSourceListPretend(testParams, &task.State{}), expectedError)
}
