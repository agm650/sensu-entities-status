package main

import (
	"errors"
	"fmt"
	"os"

	customSensu "las/accs/entities-status/sensu"

	"github.com/apex/log"
	"github.com/sensu/sensu-go/types"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Namespace        string
	Timeout          string
	RuntimeAssets    string
	SensuAPIUrl      string
	SensuAccessToken string
	SensuFormat      string
	Debug            bool
}

var (
	config = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-entities-status",
			Short:    "Sensu Entities Status. Get entities status.",
			Keyspace: "las/accs/entities-status/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "namespace",
			Env:       "SENSU_NAMESPACE", // provided by the sensuctl command plugin execution environment
			Argument:  "namespace",
			Shorthand: "n",
			Default:   "",
			Usage:     "Sensu Namespace to perform the runbook automation (defaults to $SENSU_NAMESPACE)",
			Value:     &config.Namespace,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "sensu-api-url",
			Env:       "SENSU_API_URL", // provided by the sensuctl command plugin execution environment
			Argument:  "sensu-api-url",
			Shorthand: "",
			Default:   "",
			Usage:     "Sensu API URL (defaults to $SENSU_API_URL)",
			Value:     &config.SensuAPIUrl,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "sensu-access-token",
			Env:       "SENSU_ACCESS_TOKEN", // provided by the sensuctl command plugin execution environment
			Argument:  "sensu-access-token",
			Shorthand: "",
			Default:   "",
			Usage:     "Sensu API Access Token (defaults to $SENSU_ACCESS_TOKEN)",
			Value:     &config.SensuAccessToken,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "sensu-format",
			Env:       "SENSU_FORMAT", // provided by the sensuctl command plugin execution environment
			Argument:  "sensu-format",
			Shorthand: "",
			Default:   "tabular",
			Usage:     "Sensu Format (defaults to $SENSU_FORMAT). Authorized values: tabular, yaml wrapped-json",
			Value:     &config.SensuFormat,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "sensu-debug",
			Env:       "SENSU_DEBUG", // provided by the sensuctl command plugin execution environment
			Argument:  "sensu-debug",
			Shorthand: "",
			Default:   false,
			Usage:     "Activate debug logs",
			Value:     &config.Debug,
		},
	}
)

func main() {
	plugin := sensu.NewGoCheck(&config.PluginConfig, options, checkArgs, executeCheck, false)
	plugin.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	if len(config.SensuAPIUrl) == 0 {
		return sensu.CheckStateCritical, errors.New("--sensu-api-url flag or $SENSU_API_URL environment variable must be set")
	} else if len(config.Namespace) == 0 {
		return sensu.CheckStateCritical, errors.New("--namespace flag or $SENSU_NAMESPACE environment variable must be set")
	}
	return sensu.CheckStateOK, nil
}

func printResult(statusMap map[string]customSensu.EntityStatus) {
	// Depending on format different output is possible
	if config.SensuFormat == "tabular" {
		customSensu.PrintTabularResult(statusMap)
	} else if config.SensuFormat == "yaml" {
		customSensu.PrintYAMLResult(statusMap)
	} else if config.SensuFormat == "wrapped-json" {
		customSensu.PrintJSONResult(statusMap)
	} else {
		fmt.Fprintln(os.Stderr, "Invalid format output")
	}
}

func executeCheck(event *types.Event) (int, error) {

	if config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.FatalLevel)
	}

	endpointURL := fmt.Sprintf("%s/api/core/v2/namespaces/%s/events",
		config.SensuAPIUrl,
		config.Namespace,
	)

	header := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", config.SensuAccessToken),
	}
	evts, err := customSensu.EventExtractJSONWithHeader(endpointURL, header, nil)
	if err != nil {
		return sensu.CheckStateCritical, err
	}

	entitiesStatus := customSensu.GetEntitiesStatus(evts)

	printResult(entitiesStatus)

	return sensu.CheckStateOK, nil
}
