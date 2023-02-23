package sensu

import (
	"github.com/apex/log"
	v2 "github.com/sensu/core/v2"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// EntityStatus : Structure used to sumarize an Entity current State
type EntityStatus struct {
	Status   int `json:"status" yaml:"status"`
	Silenced int `json:"silenced" yaml:"silenced"`
	Critical int `json:"critical" yaml:"critical"`
	Warning  int `json:"warning" yaml:"warning"`
	Unknown  int `json:"unknown" yaml:"unknown"`
	Ok       int `json:"ok" yaml:"ok"`
	Total    int `json:"total" yaml:"total"`
}

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

// calculateStatus : This function is used to calculate the resulting status when comparing two status
// The comparison matrix will be the following one:
//
//	|      | Crit | Warn | Unkn |  OK  | => New events
//	|------|------|------|------|------|
//	| Crit | Crit | Crit | Crit | Crit |
//	| Warn | Crit | Warn | Warn | Warn |
//	| Unkn | Crit | Warn | Unkn | Unkn |
//	|  OK  | Crit | Warn | Unkn |  Ok  |
//	   ⬇︎
//
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
		Status:   0,
		Silenced: 0,
		Critical: 0,
		Warning:  0,
		Unknown:  0,
		Ok:       0,
		Total:    0,
	}

	for _, evt := range events {
		if evt.Entity.Name != entityName {
			// Only look into our entity events
			continue
		}

		if evt.IsSilenced() {
			gstatus.Silenced++
		}
		if evt.Check.Status == sensu.CheckStateCritical {
			gstatus.Critical++
		} else if evt.Check.Status == sensu.CheckStateWarning {
			gstatus.Warning++
		} else if evt.Check.Status == sensu.CheckStateUnknown {
			gstatus.Unknown++
		} else {
			gstatus.Ok++
		}

		if !evt.IsSilenced() {
			gstatus.Status = calculateStatus(gstatus.Status, int(evt.Check.Status))
		}
	}

	ctx.Errorf("Status for %s is %d", entityName, gstatus.Status)
	ctx.Debugf("\tnb OK %d\tWarning %d\tCritical %d\tSilenced %d", gstatus.Ok, gstatus.Warning, gstatus.Critical, gstatus.Silenced)

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
			estatus.Silenced++
		}
		if evt.Check.Status == sensu.CheckStateCritical {
			estatus.Critical++
		} else if evt.Check.Status == sensu.CheckStateWarning {
			estatus.Warning++
		} else if evt.Check.Status == sensu.CheckStateUnknown {
			estatus.Unknown++
		} else {
			estatus.Ok++
		}

		if !evt.IsSilenced() {
			estatus.Status = calculateStatus(estatus.Status, int(evt.Check.Status))
		}
		set[evt.Entity.Name] = estatus
		ctx.Errorf("Setting %s status to %d", evt.Entity.Name, estatus.Status)
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
	} else if status == sensu.CheckStateOK {
		return "OK"
	} else {
		return "UNKN"
	}
}
