package app

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()
	// Mock activity implementation

	env.OnActivity(ImportUsers, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(ApproveUsers, mock.Anything).Return(nil)
	env.ExecuteWorkflow(OnboardUsers(), mock.Anything)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}
