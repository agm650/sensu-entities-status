package sensu

import (
	"testing"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/stretchr/testify/assert"
)

// GetEntitiesFromEvents : Get a list of entities based on a list of event
func TestGetEntitiesFromEvents(t *testing.T) {
	// assert := assert.New(t)
}

func TestCalculateStatus(t *testing.T) {
	assert := assert.New(t)

	// New and old event are identical, status should be equal to input
	assert.Equal(calculateStatus(sensu.CheckStateOK, sensu.CheckStateOK), sensu.CheckStateOK)
	assert.Equal(calculateStatus(sensu.CheckStateCritical, sensu.CheckStateCritical), sensu.CheckStateCritical)
	assert.Equal(calculateStatus(sensu.CheckStateWarning, sensu.CheckStateWarning), sensu.CheckStateWarning)
	assert.Equal(calculateStatus(sensu.CheckStateUnknown, sensu.CheckStateUnknown), sensu.CheckStateUnknown)

	// Old state is critical. No matter what it has to remain critical
	assert.Equal(calculateStatus(sensu.CheckStateCritical, sensu.CheckStateUnknown), sensu.CheckStateCritical)
	assert.Equal(calculateStatus(sensu.CheckStateCritical, sensu.CheckStateOK), sensu.CheckStateCritical)
	assert.Equal(calculateStatus(sensu.CheckStateCritical, sensu.CheckStateWarning), sensu.CheckStateCritical)

	// New state is critical. No matter what it has to change to critical
	assert.Equal(calculateStatus(sensu.CheckStateUnknown, sensu.CheckStateCritical), sensu.CheckStateCritical)
	assert.Equal(calculateStatus(sensu.CheckStateOK, sensu.CheckStateCritical), sensu.CheckStateCritical)
	assert.Equal(calculateStatus(sensu.CheckStateWarning, sensu.CheckStateCritical), sensu.CheckStateCritical)

	// Old state is warning. It can only change if new state is Critical
	assert.Equal(calculateStatus(sensu.CheckStateWarning, sensu.CheckStateUnknown), sensu.CheckStateWarning)
	assert.Equal(calculateStatus(sensu.CheckStateWarning, sensu.CheckStateOK), sensu.CheckStateWarning)
	assert.Equal(calculateStatus(sensu.CheckStateWarning, sensu.CheckStateWarning), sensu.CheckStateWarning)
	assert.Equal(calculateStatus(sensu.CheckStateWarning, sensu.CheckStateCritical), sensu.CheckStateCritical)

	// New state is warning. It will change result unless old state is Critical
	assert.Equal(calculateStatus(sensu.CheckStateUnknown, sensu.CheckStateWarning), sensu.CheckStateWarning)
	assert.Equal(calculateStatus(sensu.CheckStateOK, sensu.CheckStateWarning), sensu.CheckStateWarning)
	assert.Equal(calculateStatus(sensu.CheckStateWarning, sensu.CheckStateWarning), sensu.CheckStateWarning)
	assert.Equal(calculateStatus(sensu.CheckStateCritical, sensu.CheckStateWarning), sensu.CheckStateCritical)

	// Old state is unknown. It can only change if new state is warning or above
	assert.Equal(calculateStatus(sensu.CheckStateUnknown, sensu.CheckStateUnknown), sensu.CheckStateUnknown)
	assert.Equal(calculateStatus(sensu.CheckStateUnknown, sensu.CheckStateOK), sensu.CheckStateUnknown)
	assert.Equal(calculateStatus(sensu.CheckStateUnknown, sensu.CheckStateWarning), sensu.CheckStateWarning)
	assert.Equal(calculateStatus(sensu.CheckStateUnknown, sensu.CheckStateCritical), sensu.CheckStateCritical)

	// New state is unkown. It will change result unless old state is warning or above
	assert.Equal(calculateStatus(sensu.CheckStateUnknown, sensu.CheckStateUnknown), sensu.CheckStateUnknown)
	assert.Equal(calculateStatus(sensu.CheckStateOK, sensu.CheckStateUnknown), sensu.CheckStateUnknown)
	assert.Equal(calculateStatus(sensu.CheckStateWarning, sensu.CheckStateUnknown), sensu.CheckStateWarning)
	assert.Equal(calculateStatus(sensu.CheckStateCritical, sensu.CheckStateUnknown), sensu.CheckStateCritical)

	// Old state is ok. It will always change unless new state is OK
	assert.Equal(calculateStatus(sensu.CheckStateOK, sensu.CheckStateUnknown), sensu.CheckStateUnknown)
	assert.Equal(calculateStatus(sensu.CheckStateOK, sensu.CheckStateOK), sensu.CheckStateOK)
	assert.Equal(calculateStatus(sensu.CheckStateOK, sensu.CheckStateWarning), sensu.CheckStateWarning)
	assert.Equal(calculateStatus(sensu.CheckStateOK, sensu.CheckStateCritical), sensu.CheckStateCritical)

	// New state is unkown. It will change result unless old state is warning or above
	assert.Equal(calculateStatus(sensu.CheckStateUnknown, sensu.CheckStateOK), sensu.CheckStateUnknown)
	assert.Equal(calculateStatus(sensu.CheckStateOK, sensu.CheckStateOK), sensu.CheckStateOK)
	assert.Equal(calculateStatus(sensu.CheckStateWarning, sensu.CheckStateOK), sensu.CheckStateWarning)
	assert.Equal(calculateStatus(sensu.CheckStateCritical, sensu.CheckStateOK), sensu.CheckStateCritical)
}

// GetEntityStatus : Get an entity status based on a list of events
func TestGetEntityStatus(t *testing.T) {
	// assert := assert.New(t)
}

// GetEntitiesStatus : Get entities status based on a list of event
func TestGetEntitiesStatus(t *testing.T) {
	// assert := assert.New(t)
}
