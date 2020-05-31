package sensu

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/apex/log"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	v2 "github.com/sensu/sensu-go/api/core/v2"
)

// GetEntitiesFromEvents : Get a list of entities based on a list of event
func GetEntitiesFromEvents(events []v2.Event) []string {

	ctx := log.WithFields(log.Fields{
		"file":     "sensu/entities.go",
		"function": "GetEntitiesFromEvents",
	})

	var entities []string
	set := make(map[string]struct{})

	for _, evt := range events {
		set[evt.Entity.Name] = struct{}{}
	}

	for key := range set {
		entities = append(entities, key)
		ctx.Errorf("Found entity %s", key)
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

// calculateStatus : This function is used to calculate the resulting status when comparing two status
// The comparison matrix will be the following one:
//   |      | Crit | Warn | Unkn |  OK  | => New events
//   |------|------|------|------|------|
//   | Crit | Crit | Crit | Crit | Crit |
//   | Warn | Crit | Warn | Warn | Warn |
//   | Unkn | Crit | Warn | Unkn | Unkn |
//   |  OK  | Crit | Warn | Unkn |  Ok  |
//      ⬇︎
// Exisiting events
func calculateStatus(oldState int, newState int) int {
	ctx := log.WithFields(log.Fields{
		"file":     "sensu/entities.go",
		"function": "calculateStatus",
	})

	ctx.Debugf("Old State %d // New state %d", oldState, newState)
	if oldState == sensu.CheckStateCritical || newState == sensu.CheckStateCritical {
		ctx.Debugf("Returning Critical")
		return sensu.CheckStateCritical
	}

	if oldState == sensu.CheckStateWarning || newState == sensu.CheckStateWarning {
		ctx.Debugf("Returning Warning")
		return sensu.CheckStateWarning
	}

	if oldState == sensu.CheckStateUnknown || newState == sensu.CheckStateUnknown {
		ctx.Debugf("Returning Unknown")
		return sensu.CheckStateUnknown
	}

	ctx.Debugf("Returning OK")
	return sensu.CheckStateOK
}

// GetEntityStatus : Get an entity status based on a list of events
func GetEntityStatus(entityName string, events []v2.Event) EntityStatus {
	ctx := log.WithFields(log.Fields{
		"file":     "sensu/entities.go",
		"function": "GetEntityStatus",
	})

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
		if evt.Entity.Name != entityName {
			// Only look into our entity events
			continue
		}

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

	ctx.Errorf("Status for %s is %d", entityName, gstatus.status)
	ctx.Debugf("\tnb OK %d\tWarning %d\tCritical %d\tSilenced %d", gstatus.ok, gstatus.warning, gstatus.critical, gstatus.silenced)

	return gstatus
}

// GetEntitiesStatus : Get entities status based on a list of event
func GetEntitiesStatus(events []v2.Event) map[string]EntityStatus {
	ctx := log.WithFields(log.Fields{
		"file":     "sensu/entities.go",
		"function": "GetEntityStatus",
	})

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
		ctx.Errorf("Setting %s status to %d", evt.Entity.Name, estatus.status)
	}

	return set
}

func translateStatus(status int) string {
	if status == sensu.CheckStateUnknown {
		return "UNKN"
	} else if status == sensu.CheckStateWarning {
		return "WARN"
	} else if status == sensu.CheckStateCritical {
		return "CRIT"
	} else {
		return "OK"
	}
}

// PrintTabularResult : Print in tabular format the Entities Status result
func PrintTabularResult(statusMap map[string]EntityStatus) {
	ctx := log.WithFields(log.Fields{
		"file":     "sensu/entities.go",
		"function": "GetEntityStatus",
	})

	ctx.Infof("using PrintTabularResult for %d entities", len(statusMap))
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Entity\tStatus\tEvents\tSilenced\tCrtical\tWarning\tUnkown\tOk")
	fmt.Fprintln(w, "------\t------\t------\t--------\t-------\t-------\t------\t--")
	for entity, status := range statusMap {
		fmt.Fprintf(
			w,
			"%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\n",
			entity,
			translateStatus(status.status),
			status.total,
			status.silenced,
			status.critical,
			status.warning,
			status.unknown,
			status.ok,
		)
	}
	w.Flush()
}

func printJSONResult(statusMap map[string]EntityStatus) {

}

func printYAMLResult(statusMap map[string]EntityStatus) {

}
