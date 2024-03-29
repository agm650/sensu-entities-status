package sensu

import (
	"testing"

	corev2 "github.com/sensu/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/stretchr/testify/assert"
)

// GetEntitiesFromEvents : Get a list of entities based on a list of event
func TestGetEntitiesFromEvents(t *testing.T) {
	assert := assert.New(t)

	var eventList []corev2.Event = []corev2.Event{}
	eventList = append(eventList, *corev2.FixtureEvent("localhost", "dummy-check1"))
	eventList = append(eventList, *corev2.FixtureEvent("localhost", "dummy-check2"))
	eventList = append(eventList, *corev2.FixtureEvent("localhost2", "dummy-check1"))
	eventList = append(eventList, *corev2.FixtureEvent("localhost2", "dummy-check2"))
	eventList = append(eventList, *corev2.FixtureEvent("localhost3", "dummy-check3"))

	entitiesList := GetEntitiesFromEvents(eventList)

	assert.NotNil(entitiesList)

	// Duplicated entities have been removed
	assert.Len(entitiesList, 3)
	assert.Contains(entitiesList, "localhost")
	assert.Contains(entitiesList, "localhost2")
	assert.Contains(entitiesList, "localhost3")
	assert.NotContains(entitiesList, "localhost4")

	entitiesList = GetEntitiesFromEvents([]corev2.Event{})
	assert.Len(entitiesList, 0)
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
	assert := assert.New(t)

	var eventList []corev2.Event = []corev2.Event{}
	evt1 := *corev2.FixtureEvent("localhost", "dummy-check1")
	evt2 := *corev2.FixtureEvent("localhost", "dummy-check2")
	evt3 := *corev2.FixtureEvent("localhost2", "dummy-check1")
	evt4 := *corev2.FixtureEvent("localhost2", "dummy-check2")
	evt5 := *corev2.FixtureEvent("localhost3", "dummy-check3")
	evt6 := *corev2.FixtureEvent("localhost4", "dummy-check4")

	evt1.Check.Status = sensu.CheckStateOK
	evt2.Check.Status = sensu.CheckStateOK
	evt3.Check.Status = sensu.CheckStateOK
	evt4.Check.Status = sensu.CheckStateWarning
	evt5.Check.Status = sensu.CheckStateCritical
	evt6.Check.Status = sensu.CheckStateUnknown
	evt6.Check.Silenced = []string{"test"}

	eventList = append(eventList, evt1, evt2, evt3, evt4, evt5, evt6)

	// Check status for localhost
	ent1Status := GetEntityStatus("localhost", eventList)
	assert.Equal(ent1Status.Status, sensu.CheckStateOK)
	assert.Equal(ent1Status.Silenced, 0)
	assert.Equal(ent1Status.Ok, 2)
	assert.Equal(ent1Status.Warning, 0)
	assert.Equal(ent1Status.Critical, 0)
	assert.Equal(ent1Status.Unknown, 0)

	// Check status for localhost2
	ent2Status := GetEntityStatus("localhost2", eventList)
	assert.Equal(ent2Status.Status, sensu.CheckStateWarning)
	assert.Equal(ent2Status.Silenced, 0)
	assert.Equal(ent2Status.Ok, 1)
	assert.Equal(ent2Status.Warning, 1)
	assert.Equal(ent2Status.Critical, 0)
	assert.Equal(ent2Status.Unknown, 0)

	// Check status for localhost3
	ent3Status := GetEntityStatus("localhost3", eventList)
	assert.Equal(ent3Status.Status, sensu.CheckStateCritical)
	assert.Equal(ent3Status.Silenced, 0)
	assert.Equal(ent3Status.Ok, 0)
	assert.Equal(ent3Status.Warning, 0)
	assert.Equal(ent3Status.Critical, 1)
	assert.Equal(ent3Status.Unknown, 0)

	// Check status for localhost4
	ent4Status := GetEntityStatus("localhost4", eventList)
	assert.Equal(ent4Status.Status, sensu.CheckStateOK)
	assert.Equal(ent4Status.Silenced, 1)
	assert.Equal(ent4Status.Ok, 0)
	assert.Equal(ent4Status.Warning, 0)
	assert.Equal(ent4Status.Critical, 0)
	assert.Equal(ent4Status.Unknown, 1)
}

// GetEntitiesStatus : Get entities status based on a list of event
func TestGetEntitiesStatus(t *testing.T) {
	assert := assert.New(t)

	var eventList []corev2.Event = []corev2.Event{}
	evt1 := *corev2.FixtureEvent("localhost", "dummy-check1")
	evt2 := *corev2.FixtureEvent("localhost", "dummy-check2")
	evt3 := *corev2.FixtureEvent("localhost2", "dummy-check1")
	evt4 := *corev2.FixtureEvent("localhost2", "dummy-check2")
	evt5 := *corev2.FixtureEvent("localhost3", "dummy-check3")
	evt6 := *corev2.FixtureEvent("localhost4", "dummy-check4")

	evt1.Check.Status = sensu.CheckStateOK
	evt2.Check.Status = sensu.CheckStateOK
	evt3.Check.Status = sensu.CheckStateOK
	evt4.Check.Status = sensu.CheckStateWarning
	evt5.Check.Status = sensu.CheckStateCritical
	evt6.Check.Status = sensu.CheckStateUnknown
	evt6.Check.Silenced = []string{"test"}

	eventList = append(eventList, evt1, evt2, evt3, evt4, evt5, evt6)

	statuses := GetEntitiesStatus(eventList)

	// 3 entities are provided, 4 items in the map are expected
	assert.Len(statuses, 4)

	// Checking all 3 entities
	assert.Equal(statuses["localhost"].Status, sensu.CheckStateOK)
	assert.Equal(statuses["localhost"].Silenced, 0)
	assert.Equal(statuses["localhost"].Ok, 2)
	assert.Equal(statuses["localhost"].Warning, 0)
	assert.Equal(statuses["localhost"].Critical, 0)
	assert.Equal(statuses["localhost"].Unknown, 0)

	assert.Equal(statuses["localhost2"].Status, sensu.CheckStateWarning)
	assert.Equal(statuses["localhost2"].Silenced, 0)
	assert.Equal(statuses["localhost2"].Ok, 1)
	assert.Equal(statuses["localhost2"].Warning, 1)
	assert.Equal(statuses["localhost2"].Critical, 0)
	assert.Equal(statuses["localhost2"].Unknown, 0)

	assert.Equal(statuses["localhost3"].Status, sensu.CheckStateCritical)
	assert.Equal(statuses["localhost3"].Silenced, 0)
	assert.Equal(statuses["localhost3"].Ok, 0)
	assert.Equal(statuses["localhost3"].Warning, 0)
	assert.Equal(statuses["localhost3"].Critical, 1)
	assert.Equal(statuses["localhost3"].Unknown, 0)

	assert.Equal(statuses["localhost4"].Status, sensu.CheckStateOK)
	assert.Equal(statuses["localhost4"].Silenced, 1)
	assert.Equal(statuses["localhost4"].Ok, 0)
	assert.Equal(statuses["localhost4"].Warning, 0)
	assert.Equal(statuses["localhost4"].Critical, 0)
	assert.Equal(statuses["localhost4"].Unknown, 1)
}

func TestTranslateStatus(t *testing.T) {
	assert := assert.New(t)

	const MaxInt = int(^uint(0) >> 1)
	const MinInt = -(MaxInt - 1)

	assert.Equal(translateStatus(0), "OK")
	assert.Equal(translateStatus(1), "WARN")
	assert.Equal(translateStatus(2), "CRIT")
	assert.Equal(translateStatus(3), "UNKN")
	assert.Equal(translateStatus(255), "UNKN")
	assert.Equal(translateStatus(MaxInt), "UNKN")
	assert.Equal(translateStatus(-1), "UNKN")
	assert.Equal(translateStatus(MinInt), "UNKN")
}
