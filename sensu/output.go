package sensu

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/apex/log"
	"gopkg.in/yaml.v2"
)

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
			translateStatus(status.Status),
			status.Total,
			status.Silenced,
			status.Critical,
			status.Warning,
			status.Unknown,
			status.Ok,
		)
	}
	w.Flush()
}

// PrintJSONResult : Export data in JSON format
func PrintJSONResult(statusMap map[string]EntityStatus) {
	jsonString, _ := json.MarshalIndent(statusMap, "", "\t")
	fmt.Println(string(jsonString))
}

// PrintYAMLResult : Export data in JSON format
func PrintYAMLResult(statusMap map[string]EntityStatus) {
	data, err := yaml.Marshal(statusMap)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	fmt.Println(string(data))
}
