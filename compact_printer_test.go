package jl

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompactPrinter_PrintNonJSON(t *testing.T) {
	raw := "hello world"
	buf := &bytes.Buffer{}
	printer := NewCompactPrinter(buf)
	entry := &Entry{
		Raw: []byte(raw),
	}
	printer.Print(entry)
	assert.Equal(t, raw+"\n", buf.String())
}

func TestCompactPrinter_Print(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		formatted string
	}{{
		name:      "basic",
		json:      `{"timestamp":"2019-01-01 15:23:45","level":"INFO","thread":"truck-manager","logger":"TruckRepairServiceOverlordManager","message":"There are 7 more trucks in the garage to fix. Get to work."}`,
		formatted: "INFO 2019-01-01 15:23:45 [truck-manager]    TruckRepa…ordManager| There are 7 more trucks in the garage to fix. Get to work.\n",
	}, {
		name: "exception",
		json: `{"timestamp":"2019-01-01 15:34:45","level":"ERROR","thread":"repair-worker-2","logger":"TruckRepairMinion","message":"Truck 5 has is really broken! I'm need parts, waiting till they come.","exceptions":[{"frames_omitted":8,"message":"I can't do anything","module":"com.truckstop","type":"MustStandAroundException","stack_trace":[{"file":"TruckRepairMinion.java","func":"repair","line":84,"module":"com.truckstop"},{"file":"TruckRepairDepot.java","func":"makeMinionDoWork","line":1,"module":"com.truckstop"}]},{"frames_omitted":5,"message":"missing 3 parts","module":"com.truckstop","type":"MissingPartsException","stack_trace":[{"file":"TruckRepairChecklist.java","func":"checkParts","line":42,"module":"com.truckstop"},{"file":"TruckRepairMinion.java","func":"runChecklist","line":84,"module":"com.truckstop"}]}]}`,
		formatted: `ERRO 2019-01-01 15:34:45 [repair-worker-2]     TruckRepairMinion| Truck 5 has is really broken! I'm need parts, waiting till they come.
  com.truckstop.MustStandAroundException: I can't do anything
    at com.truckstop.repair(TruckRepairMinion.java:84)
    at com.truckstop.makeMinionDoWork(TruckRepairDepot.java:1)
    ...8 frames omitted...
  Caused by: com.truckstop.MissingPartsException: missing 3 parts
    at com.truckstop.checkParts(TruckRepairChecklist.java:42)
    at com.truckstop.runChecklist(TruckRepairMinion.java:84)
    ...5 frames omitted...
`}, {
		name: "logrus_pgk_error",
		json: `{"timestamp":"2019-01-01 15:23:45","level":"error","message":"an error occurred","error":"BOOM!","stack":"github.com/pkg/errors_test.fn\n\t/home/dfc/src/github.com/pkg/errors/example_test.go:47\ngithub.com/pkg/errors_test.Example_stackTrace\n\t/home/dfc/src/github.com/pkg/errors/example_test.go:127\n"}`,
		formatted: `ERRO 2019-01-01 15:23:45 an error occurred
  BOOM!
	github.com/pkg/errors_test.fn
		/home/dfc/src/github.com/pkg/errors/example_test.go:47
	github.com/pkg/errors_test.Example_stackTrace
		/home/dfc/src/github.com/pkg/errors/example_test.go:127
`,
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			printer := NewCompactPrinter(buf)
			printer.DisableColor = true
			entry := &Entry{
				Raw: []byte(test.json),
			}
			require.NoError(t, json.Unmarshal([]byte(test.json), &entry.Partials))
			printer.Print(entry)
			assert.Equal(t, test.formatted, buf.String())
		})
	}
}

func TestCompactPrinter_PrintWithColors(t *testing.T) {
	var logs = []string{
		`{"timestamp":"2019-01-01 15:24:45","level":"info","thread":"repair-worker-1","logger":"truckrepairminion","message":"fixing truck 1, it's got a broken axle"}`,
		`{"timestamp":"2019-01-01 15:25:45","level":"info","thread":"repair-worker-2","logger":"truckrepairminion","message":"fixing truck 2, it's got a broken axle"}`,
	}
	var formatted = []string{
		"\x1b[32mINFO\x1b[0m 2019-01-01 15:24:45 \x1b[32m[repair-worker-1] \x1b[0m \x1b[32m   truckrepairminion|\x1b[0m fixing truck 1, it's got a broken axle\n",
		"\x1b[32mINFO\x1b[0m 2019-01-01 15:25:45 \x1b[33m[repair-worker-2] \x1b[0m \x1b[32m   truckrepairminion|\x1b[0m fixing truck 2, it's got a broken axle\n",
	}
	printer := NewCompactPrinter(nil)
	for i, log := range logs {
		buf := &bytes.Buffer{}
		printer.Out = buf
		entry := &Entry{
			Raw: []byte(log),
		}
		require.NoError(t, json.Unmarshal([]byte(log), &entry.Partials))
		printer.Print(entry)
		assert.Equal(t, formatted[i], buf.String())
	}
}
