package main

import (
	"strings"

	"github.com/clarify/clarify-go/automation"
	"github.com/clarify/clarify-go/query"
	"github.com/clarify/clarify-go/views"
)

// publishRules defines a map of named publishing rules. To be able to easily
// reference named rules on the command-line, you should avoid spaces, commas
// and special characters in the rule names.
var publishRules = map[string]automation.PublishSignals{
	"example-rule": {
		// List of integration IDs to publish signals from using this rule set.
		Integrations: []string{},
		// Filter to apply to signals before publishing (optional).
		SignalsFilter:    query.Field("annotations.clarify/template-clarify-automation/publish", query.Equal("true")),
		TransformVersion: "v0",
		// List of transforms to apply.
		Transforms: []func(item *views.ItemSave){
			// Transform functions can be inlined. This could be useful if you
			// want to create a transform that's not reusable by other
			// configurations.
			func(item *views.ItemSave) {
				// Detect names of format "device-id/measurement-name/eng-unit".
				// Note that production code may want to use regular expressions
				// to specify more precise rules.
				comps := strings.Split(item.Name, "/")
				if len(comps) == 3 {
					item.Name = comps[1]
					item.EngUnit = comps[2]
					item.Labels.Add("device-id", comps[0])
				}

			},
			// You can also refer to transforms you have written elsewhere. The
			// following transform functions are defined in the
			// publish_transforms.go file.
			prettifyEngUnit,
			addISQLabels,
		},
	},
}
