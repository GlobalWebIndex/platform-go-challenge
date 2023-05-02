package cache

import (
	"testing"
	"time"
)

// TestCacheStore test the possible outcomes of the Store operation.
func TestCacheStore(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"resource error":                        testCacheStoreResourceError,
		"entry with custom expiration success":  testCacheStoreEntryWithCustomExpiration,
		"entry with default expiration success": testCacheStoreEntryWithDefaultExpiration,
		"entry without invalidator success":     testCacheStoreEntryWithoutInvalidator,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testCacheStoreResourceError (test scenario)
func testCacheStoreResourceError(t *testing.T) {
	ConnectToCache("memory")
	Cache.Release()
	err := Cache.Store("test", []byte("exampleResponse"), false)
	if err == nil {
		t.Fatalf("cache resource is still available")
	}
}

// testCacheStoreEntryWithCustomExpiration (test scenario)
func testCacheStoreEntryWithCustomExpiration(t *testing.T) {
	ConnectToCache("memory")
	expiration := time.Now().Add(20 * time.Second)
	err := Cache.Store("test", []byte("exampleResponse"), true, expiration)
	if err != nil {
		t.Fatalf("failed during store operation")
	}
	Cache.Release()
}

// testCacheStoreEntryWithDefaultExpiration (test scenario)
func testCacheStoreEntryWithDefaultExpiration(t *testing.T) {
	ConnectToCache("memory")
	err := Cache.Store("test", []byte("exampleResponse"), true)
	if err != nil {
		t.Fatalf("failed during store operation")
	}
	Cache.Release()
}

// testCacheStoreEntryWithoutInvalidator (test scenario)
func testCacheStoreEntryWithoutInvalidator(t *testing.T) {
	ConnectToCache("memory")
	err := Cache.Store("test", []byte("exampleResponse"), false)
	if err != nil {
		t.Fatalf("failed during store operation")
	}
	Cache.Release()
}

// TestCacheLoad test the possible outcomes of the Load operation.
func TestCacheLoad(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"resource error":      testCacheLoadResourceError,
		"cache entry not set": testCacheLoadNotSet,
		"cache entry expired": testCacheLoadExpired,
		"cache load success":  testCacheLoadSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testCacheLoadResourceError (test scenario)
func testCacheLoadResourceError(t *testing.T) {
	ConnectToCache("memory")
	Cache.Release()
	_, _, err := Cache.Load("test")
	if err == nil {
		t.Fatalf("cache resource is still available")
	}
}

// testCacheLoadNotSet (test scenario)
func testCacheLoadNotSet(t *testing.T) {
	ConnectToCache("memory")
	set, _, err := Cache.Load("test")
	if err != nil {
		t.Fatalf("cache resource is not available available")
	}
	if set {
		t.Fatalf("cache entry is set")
	}
	Cache.Release()
}

// testCacheLoadExpired (test scenario)
func testCacheLoadExpired(t *testing.T) {
	ConnectToCache("memory")
	data := MockCacheEntryData()
	err := Cache.Store(data.Token, data.Response, false, time.Now())
	if err != nil {
		t.Fatalf("failed during store operation")
	}
	set, _, err := Cache.Load(data.Token)
	if err != nil {
		t.Fatalf("failed during load operation")
	}
	if set {
		t.Fatalf("cache entry is not yet expired")
	}
	Cache.Release()
}

// testCacheLoadNotSuccess (test scenario)
func testCacheLoadSuccess(t *testing.T) {
	ConnectToCache("memory")
	data := MockCacheEntryData()
	err := Cache.Store(data.Token, data.Response, false)
	if err != nil {
		t.Fatalf("failed during store operation")
	}
	set, resp, err := Cache.Load(data.Token)
	if err != nil {
		t.Fatalf("failed during load operation")
	}
	if !set {
		t.Fatalf("cache entry is not set")
	}
	if string(resp) != string(data.Response) {
		t.Fatalf("data corrupted during operations")
	}
	Cache.Release()
}

// TestCacheEvict test the possible outcomes of the Evict operation.
func TestCacheEvict(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"evict cache entry":    testCacheEvictSuccess,
		"cache resource error": testCacheEvictResourceError,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testCacheEvicSuccess (test scenario)
func testCacheEvictSuccess(t *testing.T) {
	ConnectToCache("memory")
	data := MockCacheEntryData()
	err := Cache.Store(data.Token, data.Response, false)
	if err != nil {
		t.Fatalf("Failed during store operation")
	}
	err = Cache.Evict(data.Token)
	if err != nil {
		t.Fatalf("Failed during evict operation")
	}
	set, resp, e := Cache.Load(data.Token)
	if e != nil {
		t.Fatalf("Failed during load operation")
	}
	if set {
		t.Errorf("cache entry still available: got Set=%v and Data=%v want Set=%v and Data=%v",
			set, string(resp), false, nil)
	}
	Cache.Release()
}

// testCacheEvictResourceError
func testCacheEvictResourceError(t *testing.T) {
	ConnectToCache("memory")
	Cache.Release()
	err := Cache.Evict("test")
	if err == nil {
		t.Fatalf("cache resource is still available")
	}
}

// TestCacheRelease test the possible outcomes of the Release operation.
func TestCacheRelease(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"cache release success": testCacheReleaseSuccess,
		"cache resource error":  testCacheReleaseError,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testCacheReleaseSuccess (test scenario)
func testCacheReleaseSuccess(t *testing.T) {
	ConnectToCache("memory")
	err := Cache.Release()
	if err != nil {
		t.Fatalf("cache resource failed to be released")
	}
	_, _, err = Cache.Load("test")
	if err == nil {
		t.Fatalf("cache resource is still allocated")
	}
}

// testCacheReleaseSuccess (test scenario)
func testCacheReleaseError(t *testing.T) {
	ConnectToCache("memory")
	err := Cache.Release()
	// Release the resource again to trigger the error.
	err = Cache.Release()
	if err == nil {
		t.Fatalf("cache resource is already released and should trigger an error")
	}
}

// TestCacheEntryInvalidator tests the possible outcomes of of the check for
// cache entry expiry and invalidation.
// it expires.
func TestCacheEntryInvalidator(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"cache entry invalidation success": testCacheEntryInvalidatorSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testCacheEntryInvalidatorSuccess (test scenario)
func testCacheEntryInvalidatorSuccess(t *testing.T) {
	ConnectToCache("memory")
	// Create cache entry with invalidator
	data := MockCacheEntryData()
	expiration := time.Now()
	err := Cache.Store(data.Token, data.Response, true, expiration)
	if err != nil {
		t.Fatalf("failed during store operation")
	}
	set, resp, err := Cache.Load(data.Token)
	if err != nil {
		t.Fatalf("failed during load operation")
	}
	if set {
		t.Fatalf("Invalidator did not operate correctly: got Set=%v and Data=%v want Set=%v and Data=%v",
			set, string(resp), false, nil)
	}
	err = Cache.Release()
	if err != nil {
		t.Fatalf("Failed to release cache resource")
	}
}
