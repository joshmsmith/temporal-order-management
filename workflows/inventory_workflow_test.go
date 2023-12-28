package workflows

import (
	"ordermanagement/activities"
	
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"

)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_Workflow() {
	env := s.NewTestWorkflowEnvironment()
	env.RegisterActivity(activities.CheckFraud)
	env.OnActivity(activities.CheckFraud, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(func(ctx context.Context, orderID string, paymentInfo string) (error) {
			
			return activities.CheckFraud(ctx, orderID, paymentInfo)
		})
	env.ExecuteWorkflow(InventoryWorkflow)

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	var result string
	s.NoError(env.GetWorkflowResult(&result))
	s.Equal("Branch 0 done in 1ns.", result)
	env.AssertExpectations(s.T())
}
