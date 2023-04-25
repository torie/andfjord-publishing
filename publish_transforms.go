package main

import (
	"strings"

	"github.com/clarify/clarify-go/views"
)

// In this file, you can define your own set of publishSignal transform
// functions. The provided functions serve as an example only.

// prettifyEngUnit is an example transform rule for prettifying the engUnit of
// temperature types.
func prettifyEngUnit(item *views.ItemSave) {
	switch strings.ToUpper(item.EngUnit) {
	case "DEG_C", "C", "CELSIUS", "℃":
		item.EngUnit = "°C"
	case "DEG_F", "F", "FAHRENHEIT":
		item.EngUnit = "°F"
	case "DEG_K", "KELVIN":
		item.EngUnit = "K"
	}
}

// addISQLabels is an example transform rule that sets labels according to the
// quantity type we detect for items. The example only handles thermodynamic
// temperature, and requires the engUnit to be prettified.
//
// Read more about ISQ (International System of Quantities) here:
//   - https://en.wikipedia.org/wiki/International_System_of_Quantities#Base_quantities
func addISQLabels(item *views.ItemSave) {
	switch item.EngUnit {
	case "°C", "°F", "K":
		item.Labels.Add("isq-quantity", "Thermodynamic temperature")
		item.Labels.Add("isq-quantity-symbol", "T")
		item.Labels.Add("isq-dimension-symbol", "Θ")
	}
}
