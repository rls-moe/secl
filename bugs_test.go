package secl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.rls.moe/secl/query"
	"go.rls.moe/secl/types"
)

// This file tests various bugs and issues found over time, each is
// a minimal testcase

type BugIssue4 struct {
	BugIssue4Inner `secl:"bug4"`
	Decoy          int
}

type BugIssue4Inner struct {
	mock.Mock
}

func (m *BugIssue4Inner) UnmarshalSECL(v types.Value) (err error) {
	args := m.Called(v)
	return args.Error(0)
}

func TestBugIssue4(t *testing.T) {
	assert := assert.New(t)

	testObj := new(BugIssue4)
	secl := "(bug4: test)"
	exptV := &types.MapList{
		Map: map[types.String]types.Value{
			types.String{Value: "bug4"}: &types.String{Value: "test"},
		},
		List: []types.Value{},
	}

	ml, err := LoadConfig(secl)
	if err != nil {
		assert.FailNow(err.Error())
	}
	testObj.BugIssue4Inner.On("UnmarshalSECL", exptV).Return(nil)
	if err := query.SimpleStructUnmarshal(ml, testObj); err != nil {
		assert.FailNow(err.Error())
	}
	testObj.AssertExpectations(t)
}
