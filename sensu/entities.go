package sensu

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	v2 "github.com/sensu/sensu-go/api/core/v2"
)

// GetEntitiesFromEvents : Get a list of entities based on a list of event
func GetEntitiesFromEvents(events []v2.Event) []string {

	var entities []string
	set := make(map[string]struct{})

	for _, evt := range events {
		set[evt.Entity.Name] = struct{}{}
	}

	for key := range set {
		entities = append(entities, key)
	}

	return entities
}

// EntityStatus : Structure used to sumarize an Entity current State
type EntityStatus struct {
	status   int
	silenced int
	critical int
	warning  int
	unknown  int
	ok       int
	total    int
}

func calculateStatus(oldState int, newState int) int {

	if oldState == sensu.CheckStateCritical || newState == sensu.CheckStateCritical {
		return sensu.CheckStateCritical
	}

	if oldState == sensu.CheckStateWarning || newState == sensu.CheckStateWarning {
		return sensu.CheckStateWarning
	}

	if oldState == sensu.CheckStateUnknown || newState == sensu.CheckStateUnknown {
		return sensu.CheckStateUnknown
	}

	return sensu.CheckStateOK
}

// GetEntityStatus : Get an entity status based on a list of events
func GetEntityStatus(entityName string, events []v2.Event) EntityStatus {

	gstatus := EntityStatus{
		status:   0,
		silenced: 0,
		critical: 0,
		warning:  0,
		unknown:  0,
		ok:       0,
		total:    0,
	}

	for _, evt := range events {
		if evt.IsSilenced() {
			gstatus.silenced++
		}
		if evt.Check.Status == sensu.CheckStateCritical {
			gstatus.critical++
		} else if evt.Check.Status == sensu.CheckStateWarning {
			gstatus.warning++
		} else if evt.Check.Status == sensu.CheckStateUnknown {
			gstatus.unknown++
		} else {
			gstatus.ok++
		}

		if !evt.IsSilenced() {
			gstatus.status = calculateStatus(gstatus.status, int(evt.Check.Status))
		}
	}

	return gstatus
}

// GetEntitiesStatus : Get entities status based on a list of event
func GetEntitiesStatus(events []v2.Event) map[string]EntityStatus {

	set := make(map[string]EntityStatus)

	for _, evt := range events {
		var estatus EntityStatus
		if _, ok := set[evt.Entity.Name]; ok {
			// Update existing status
			estatus = set[evt.Entity.Name]
		} else {
			estatus = EntityStatus{}
		}

		if evt.IsSilenced() {
			estatus.silenced++
		}
		if evt.Check.Status == sensu.CheckStateCritical {
			estatus.critical++
		} else if evt.Check.Status == sensu.CheckStateWarning {
			estatus.warning++
		} else if evt.Check.Status == sensu.CheckStateUnknown {
			estatus.unknown++
		} else {
			estatus.ok++
		}

		if !evt.IsSilenced() {
			estatus.status = calculateStatus(estatus.status, int(evt.Check.Status))
		}
		set[evt.Entity.Name] = estatus
	}

	return set
}

// PrintTabularResult : Print in tabular format the Entities Status result
func PrintTabularResult(statusMap map[string]EntityStatus) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Entity\tStatus\tEvents\tSilenced\tCrtical\tWarning\tUnkown\tOk")
	fmt.Fprintln(w, "------\t------\t------\t--------\t-------\t-------\t------\t--")
	for entity, status := range statusMap {
		fmt.Fprintf(
			w,
			"%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n",
			entity,
			status.status,
			status.total,
			status.silenced,
			status.critical,
			status.warning,
			status.unknown,
			status.ok,
		)
	}
}

func printJSONResult(statusMap map[string]EntityStatus) {

}

func printYAMLResult(statusMap map[string]EntityStatus) {

}
