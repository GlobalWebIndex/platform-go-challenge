package cfg

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

// TestReadConfig tests all the possible scenarios for reading configuration //
// from  a file. //
func TestReadConfig(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"valid provided configuration file": testReadConfigValidProvided,
		"valid default configuration file":  testReadConfigValidDefault,
		"invalid configuration file":        testReadConfigInvalid,
		"configuration file not found":      testReadConfigFileNotFound,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testReadConfigValidProvided  (test scenario)
func testReadConfigValidProvided(t *testing.T) {
	_, err := ReadConfig("../.env")
	if err != nil {
		t.Errorf("config reader returned wrong result: got %v want %v",
			err, nil)
	}
}

// testReadConfigValidDefault  (test scenario)
func testReadConfigValidDefault(t *testing.T) {
	_, err := ReadConfig("")
	if err != nil {
		t.Errorf("config reader returned wrong result: got %v want %v",
			err, nil)
	}
}

// testReadConfigInvalid  (test scenario)
func testReadConfigInvalid(t *testing.T) {
	tmpJsonConfig := map[string]any{
		"FirstName": "Mark",
		"LastName":  "Jones",
		"Email":     "mark@gmail.com",
		"Age":       25,
	}
	file, _ := json.MarshalIndent(tmpJsonConfig, "", " ")
	_ = ioutil.WriteFile("cfg.json", file, 0644)
	_, err := ReadConfig("cfg.json")
	if err == nil {
		t.Errorf("config reader returned wrong result: got %v want %v",
			nil, "wrong configuration file type")
	}
	err2 := os.Remove("cfg.json")
	if err2 != nil {
		t.Log("Failed to remove temporary configuration file")
	}
}

// testReadConfigFileNotFOund  (test scenario)
func testReadConfigFileNotFound(t *testing.T) {
	_, err := ReadConfig(".env")
	if err == nil {
		t.Errorf("config reader returned wrong result: got %v want %v",
			nil, "configuration file not found")
	}
}
