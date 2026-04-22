package bambulan

import (
	"testing"
)

func TestTempLimitsRefactor(t *testing.T) {
	// Test cases based on printer_capabilities.json
	tests := []struct {
		model          string
		expectedBed    int
		expectedNozzle int
	}{
		{"BL-P001", 110, 300}, // X1C
		{"N1", 80, 300},       // A1 Mini
		{"C11", 100, 300},     // P1P
		{"C13", 110, 320},     // X1E
		{"Unknown", 100, 300}, // Fallback
	}

	for _, tc := range tests {
		t.Run(tc.model, func(t *testing.T) {
			gotBed := getBedTempLimit(tc.model)
			if gotBed != tc.expectedBed {
				t.Errorf("getBedTempLimit(%s) = %d; want %d", tc.model, gotBed, tc.expectedBed)
			}

			gotNozzle := getNozzleTempLimit(tc.model)
			if gotNozzle != tc.expectedNozzle {
				t.Errorf("getNozzleTempLimit(%s) = %d; want %d", tc.model, gotNozzle, tc.expectedNozzle)
			}
		})
	}
}

func TestBuildStartPrintCommandForGcode3MFIncludesStudioParityFields(t *testing.T) {
	autoBed := 2
	extrudeFlag := 1
	extrudeManual := 0
	nozzleOffset := 2

	cmd := buildStartPrintCommand("20001", "/Cube.gcode.3mf", PrintOptions{
		BedType:               "textured_plate",
		Timelapse:             false,
		BedLeveling:           true,
		FlowCalibration:       true,
		VibrationCalibration:  false,
		LayerInspection:       true,
		UseAMS:                true,
		MD5:                   "ABCDEF",
		AMSMapping:            []int{1, -1},
		AutoBedLevelingMode:   &autoBed,
		ExtrudeCaliFlag:       &extrudeFlag,
		ExtrudeCaliManualMode: &extrudeManual,
		NozzleOffsetCali:      &nozzleOffset,
	})

	printCmd := cmd["print"].(map[string]any)
	if got := printCmd["command"]; got != "project_file" {
		t.Fatalf("command = %v, want project_file", got)
	}
	if got := printCmd["file"]; got != "Cube.gcode.3mf" {
		t.Fatalf("file = %v, want Cube.gcode.3mf", got)
	}
	if got := printCmd["param"]; got != "Metadata/plate_1.gcode" {
		t.Fatalf("param = %v, want Metadata/plate_1.gcode", got)
	}
	if got := printCmd["url"]; got != "ftp:///Cube.gcode.3mf" {
		t.Fatalf("url = %v, want ftp:///Cube.gcode.3mf", got)
	}
	if got := printCmd["subtask_name"]; got != "Cube" {
		t.Fatalf("subtask_name = %v, want Cube", got)
	}
	if got := printCmd["md5"]; got != "ABCDEF" {
		t.Fatalf("md5 = %v, want ABCDEF", got)
	}
	if got := printCmd["cfg"]; got != "0" {
		t.Fatalf("cfg = %v, want 0", got)
	}

	gotMapping := printCmd["ams_mapping"].([]int)
	if len(gotMapping) != 2 || gotMapping[0] != 1 || gotMapping[1] != -1 {
		t.Fatalf("ams_mapping = %#v, want []int{1,-1}", gotMapping)
	}
	gotMapping2 := printCmd["ams_mapping2"].([]AMSSlotMapping)
	if len(gotMapping2) != 2 {
		t.Fatalf("ams_mapping2 len = %d, want 2", len(gotMapping2))
	}
	if gotMapping2[0] != (AMSSlotMapping{AMSID: 0, SlotID: 1}) {
		t.Fatalf("ams_mapping2[0] = %#v", gotMapping2[0])
	}
	if gotMapping2[1] != (AMSSlotMapping{AMSID: 255, SlotID: 255}) {
		t.Fatalf("ams_mapping2[1] = %#v", gotMapping2[1])
	}
}

func TestBuildStartPrintCommandRespectsExplicitOverrides(t *testing.T) {
	cmd := buildStartPrintCommand("42", "job.gcode", PrintOptions{
		PlateGCodePath: "plates/custom.gcode",
		SubtaskName:    "Custom Job",
		AMSMapping2:    []AMSSlotMapping{{AMSID: 3, SlotID: 1}},
	})
	printCmd := cmd["print"].(map[string]any)
	if got := printCmd["param"]; got != "plates/custom.gcode" {
		t.Fatalf("param = %v, want plates/custom.gcode", got)
	}
	if got := printCmd["subtask_name"]; got != "Custom Job" {
		t.Fatalf("subtask_name = %v, want Custom Job", got)
	}
	gotMapping2 := printCmd["ams_mapping2"].([]AMSSlotMapping)
	if len(gotMapping2) != 1 || gotMapping2[0] != (AMSSlotMapping{AMSID: 3, SlotID: 1}) {
		t.Fatalf("ams_mapping2 = %#v", gotMapping2)
	}
}

func TestNormalizeSubtaskName(t *testing.T) {
	tests := map[string]string{
		"Cube.gcode.3mf": "Cube",
		"Cube.3mf":       "Cube",
		"Cube.gcode":     "Cube",
		"Cube.gco":       "Cube",
	}
	for input, want := range tests {
		if got := normalizeSubtaskName(input, ""); got != want {
			t.Fatalf("%s => %s, want %s", input, got, want)
		}
	}
}
