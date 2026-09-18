package source

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"genaiz.com/genaiz-lib/mock"
)

type stubListExecutor struct {
	listArg   string
	listError error
}

func (sa *stubListExecutor) List(listArg string) error {
	sa.listArg = listArg
	return sa.listError
}

func TestNewListSource(t *testing.T) {
	var expectedArg = "expectedArg"
	var testStubExecutor = &stubListExecutor{}
	var testCmd = NewListSource(newListFactory(testStubExecutor))

	testCmd.SetArgs([]string{expectedArg})
	assert.NoError(t, testCmd.Execute())
	assert.Equal(t, expectedArg, testStubExecutor.listArg)
}

func TestNewListSource_Error(t *testing.T) {
	var expectedError = errors.New("expected")
	var patch = mock.Patches{T: t}.OsExit(func(int) {})
	var testStubExecutor = &stubListExecutor{
		listError: expectedError,
	}
	var testCmd = NewListSource(newListFactory(testStubExecutor))

	defer patch.Unpatch()
	testCmd.SetArgs([]string{})
	assert.NoError(t, testCmd.Execute())
	assert.True(t, patch.Called)
	assert.Equal(t, 1, patch.CalledWith)
}

func newListFactory(executor ListExecutor) ListExecutorFactory {
	return func(command *cobra.Command) ListExecutor {
		return executor
	}
}
