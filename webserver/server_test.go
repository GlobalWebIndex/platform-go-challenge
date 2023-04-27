package webserver

import "testing"

// TestCreateServer tests all the possible scenarios for creation of a new //
// webserver instance. //
func TestCreateServer(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"configuration file not found": testCreateServerNoConfig,
		"unsupported storage option":   testCreateServerUnsupportedStorage,
		"unsupported cache option":     testCreateServerUnsupportedCache,
		"server creation success":      testCreateServerSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testACreateServerNoConfig  (test scenario)
func testCreateServerNoConfig(t *testing.T) {
	_, err := CreateServer(".nonexistentEnvFile")
	if err == nil {
		t.Errorf("server creator returned wrong result: got %v want %v",
			nil, err)
	}
}

// testACreateServerUnsupportedStorage  (test scenario)
func testCreateServerUnsupportedStorage(t *testing.T) {
	_, err := CreateServer(".env", "db")
	if err == nil {
		t.Errorf("server creator returned wrong result: got %v want %v",
			nil, err)
	}
}

// testACreateServerUnsupportedCache  (test scenario)
func testCreateServerUnsupportedCache(t *testing.T) {
	_, err := CreateServer(".env", "memory", "redis")
	if err == nil {
		t.Errorf("server creator returned wrong result: got %v want %v",
			nil, err)
	}
}

// testACreateServerSuccess  (test scenario)
func testCreateServerSuccess(t *testing.T) {
	_, err := CreateServer("../.env")
	if err != nil {
		t.Errorf("server creator returned wrong result: got %v want %v",
			err, nil)
	}
}
