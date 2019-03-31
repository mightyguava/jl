package jl

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogfmtPrinter_Print(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		formatted string
		color     bool
	}{{
		name:      "basic",
		json:      `{"timestamp":"2019-01-01 15:23:45","level":"INFO","thread":"truck-manager","logger":"TruckRepairServiceOverlordManager","message":"There are 7 more trucks in the garage to fix. Get to work."}`,
		formatted: "timestamp=2019-01-01 15:23:45 level=INFO thread=truck-manager logger=TruckRepairServiceOverlordManager message=There are 7 more trucks in the garage to fix. Get to work.\n",
	}, {
		name:      "color",
		json:      `{"timestamp":"2019-01-01 15:23:45","level":"INFO","thread":"truck-manager","logger":"TruckRepairServiceOverlordManager","message":"There are 7 more trucks in the garage to fix. Get to work."}`,
		color:     true,
		formatted: "\x1b[32mtimestamp\x1b[0m=2019-01-01 15:23:45 \x1b[32mlevel\x1b[0m=INFO \x1b[32mthread\x1b[0m=truck-manager \x1b[32mlogger\x1b[0m=TruckRepairServiceOverlordManager \x1b[32mmessage\x1b[0m=There are 7 more trucks in the garage to fix. Get to work.\n",
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			printer := NewLogfmtPrinter(buf)
			printer.DisableColor = !test.color
			entry := &Entry{
				Raw: []byte(test.json),
			}
			require.NoError(t, json.Unmarshal([]byte(test.json), &entry.Partials))
			printer.Print(entry)
			assert.Equal(t, test.formatted, buf.String())
		})
	}
}
