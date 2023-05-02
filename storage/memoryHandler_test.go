package storage

import (
	"challenge/models"
	"testing"
	"time"
)

// TestMemoryRetrieveUserAssets tests the possible outcomes from retrieving all
// the favourite assets of a user.
func TestMemoryRetrieveUserAssets(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"resource error":      testMRUAResourceError,
		"user does not exist": testMRUAUserNotExist,
		"assets retrieved":    testMRUASuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testMRUAResourceError (test scenario)
func testMRUAResourceError(t *testing.T) {
	ConnectToStorage("memory")
	err := Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
	_, _, err = Db.RetrieveUserAssets("testuser")
	if err == nil {
		t.Fatalf("storage resource is still available")
	}
}

// testMRUAUserNotExist (test scenario)
func testMRUAUserNotExist(t *testing.T) {
	ConnectToStorage("memory")
	res, _, err := Db.RetrieveUserAssets("testuser")
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if res {
		t.Fatalf("user exists")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testMRUASuccess (test scenario)
func testMRUASuccess(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	_, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	res, _, err := Db.RetrieveUserAssets(newUser.ID)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !res {
		t.Fatalf("user does not exist")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestMemoryAddAssetToUser tests the possible outcomes from adding an asset to
// user's favourite assets.
func TestMemoryAddAssetToUser(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"resource error":       testMAATUResourceError,
		"user does not exist":  testMAATUUserNotExist,
		"asset already exists": testMAATUAssetExists,
		"asset added":          testMAATUSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testMAATUResourceError (test scenario)
func testMAATUResourceError(t *testing.T) {
	ConnectToStorage("memory")
	err := Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
	newAsset := models.Asset{}
	_, err = Db.AddAssetToUser("testuser", newAsset)
	if err == nil {
		t.Fatalf("storage resource is still available")
	}
}

// test testMAATUUserNotExist (test scenario)
func testMAATUUserNotExist(t *testing.T) {
	ConnectToStorage("memory")
	newAsset := models.Asset{}
	res, err := Db.AddAssetToUser("testuser", newAsset)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if res {
		t.Fatalf("user exists")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testMAATUAssetExists (test scenario)
func testMAATUAssetExists(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	_, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	newAsset := models.Asset{}
	res, err := Db.AddAssetToUser(newUser.ID, newAsset)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !res {
		t.Fatalf("asset already exists")
	}
	res, err = Db.AddAssetToUser(newUser.ID, newAsset)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if res {
		t.Fatalf("asset does not exist")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testMAATUSuccess (test scenario)
func testMAATUSuccess(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	_, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	newAsset := models.Asset{}
	res, err := Db.AddAssetToUser(newUser.ID, newAsset)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !res {
		t.Fatalf("asset already exists")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestMemoryUpdateUserAssetDescription tests the possible outcomes from
// updating the description on a user's asset.
func TestMemoryUpdateUserAssetDescription(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"resource error":            testMUUADResourceError,
		"user does not exist":       testMUUAUserNotExist,
		"asset does not exist":      testMUUAAssetNotExist,
		"asset description updated": testMUUAASuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testMUUADResourceError (test scenario)
func testMUUADResourceError(t *testing.T) {
	ConnectToStorage("memory")
	err := Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
	_, err = Db.RemoveUserAsset("testuser", "testasset")
	if err == nil {
		t.Fatalf("storage resource is still available")
	}
}

// testMUUAUserNotExist (test scenario)
func testMUUAUserNotExist(t *testing.T) {
	ConnectToStorage("memory")
	res, err := Db.UpdateUserAssetDescription("testuser", "testasset", "newDescription")
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if res {
		t.Fatalf("user exists")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testMUUAAssetNotExist (test scenario)
func testMUUAAssetNotExist(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	_, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	res, err := Db.UpdateUserAssetDescription(newUser.ID, "testasset", "newDescription")
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if res {
		t.Fatalf("asset exists")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testMUUAASuccess (test scenario)
func testMUUAASuccess(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	_, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	newAsset := models.Asset{}
	res, err := Db.AddAssetToUser(newUser.ID, newAsset)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !res {
		t.Fatalf("failed to add asset to user")
	}
	res, err = Db.UpdateUserAssetDescription(newUser.ID, newAsset.ID, "newDescription")
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !res {
		t.Fatalf("asset does not exists")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestRemoveUser tests the possible outcomes from removing a user.
func TestRemoveUser(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"resource error":      testRUResourceError,
		"user does not exist": testRUNotExist,
		"user removed":        testRUSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testRUResourceError (test scenario)
func testRUResourceError(t *testing.T) {
	ConnectToStorage("memory")
	err := Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
	_, err = Db.RemoveUser("testuser")
	if err == nil {
		t.Fatalf("storage resource is still available")
	}
}

// testRUNotExist (test scenario)
func testRUNotExist(t *testing.T) {
	ConnectToStorage("memory")
	found, err := Db.RemoveUser("testuser")
	if err != nil {
		t.Fatalf("storage resource is still available")
	}
	if found {
		t.Fatalf("user exists")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testRUSuccess (test scenario)
func testRUSuccess(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{
		ID:         "testuser",
		Password:   "testpass",
		Cookie:     "",
		Privileged: false,
		Created:    time.Now(),
	}
	res, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !res {
		t.Fatalf("failed to add example user")
	}
	found, err := Db.RemoveUser("testuser")
	if err != nil {
		t.Fatalf("storage resource is still available")
	}
	if !found {
		t.Fatalf("user does not exist")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestMemoryRemoveUserAsset tests the possible outcomes from removing an asset
// from a user.
func TestMemoryRemoveUserAsset(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"resource error":       testMRUEResourceError,
		"user does not exist":  testMRUEUserNotExist,
		"asset does not exist": testMRUEAssetNotExist,
		"asset removed":        testMRUEASuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testMRUEResourceError (test scenario)
func testMRUEResourceError(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	_, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	newAsset := models.Asset{}
	res, err := Db.AddAssetToUser(newUser.ID, newAsset)
	if err != nil {
		t.Fatalf("failed during asset add operation")
	}
	if !res {
		t.Fatalf("asset already exists")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
	_, err = Db.RemoveUserAsset(newUser.ID, newAsset.ID)
	if err == nil {
		t.Fatalf("storage resource is still available")
	}
}

// testMRUEUserNotExist (test scenario)
func testMRUEUserNotExist(t *testing.T) {
	ConnectToStorage("memory")
	_, err := Db.RemoveUserAsset("testuser", "testAsset")
	if err == nil {
		t.Fatalf("storage resource is still available")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testMRUEAssetNotExist (test scenario)
func testMRUEAssetNotExist(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	res, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !res {
		t.Fatalf("failed to add user")
	}
	res, err = Db.RemoveUserAsset(newUser.ID, "testAsset")
	if err != nil {
		t.Fatalf("storage resource is still available")
	}
	if res {
		t.Fatalf("asset exists")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testMRUEASuccess (test scenario)
func testMRUEASuccess(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	res, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !res {
		t.Fatalf("failed to add user")
	}
	newAsset := models.Asset{}
	res, err = Db.AddAssetToUser(newUser.ID, newAsset)
	if err != nil {
		t.Fatalf("failed during add asset operation")
	}
	if !res {
		t.Fatalf("asset already exists")
	}
	res, err = Db.RemoveUserAsset(newUser.ID, "testAsset")
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if res {
		t.Fatalf("asset exists")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestMemoryCheckUserExistence tests the possible outcomes from checking
// if user exists.
func TestMemoryCheckUserExistence(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"resource error":      testMCUEResourceError,
		"user does not exist": testMCUENotExist,
		"user exists":         testMCUESuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testMCUEResourceError (test scenario)
func testMCUEResourceError(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	_, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
	exists, err := Db.CheckUserExistence(newUser.ID)
	if err == nil {
		t.Fatalf("storage resource is still available")
	}
	if exists {
		t.Fatalf("user exists")
	}
}

// testMCUENotExist (test scenario)
func testMCUENotExist(t *testing.T) {
	ConnectToStorage("memory")
	exists, err := Db.CheckUserExistence("testuser")
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if exists {
		t.Fatalf("user exists")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testMCUESuccess (test scenario)
func testMCUESuccess(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	_, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	exists, err := Db.CheckUserExistence(newUser.ID)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !exists {
		t.Fatalf("user does not exist")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestMemoryCheckUserAssetExistence tests the possible outcomes from checking
// if an asset is favourited by the specified user.
func TestMemoryCheckUserAssetExistence(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"resource error":       testMCUAssetResourceError,
		"asset does not exist": testMCUAssetAlreadyExists,
		"asset exists":         testMCUAESuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testMCUAssetResourceError (test scenario)
func testMCUAssetResourceError(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	_, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	newAsset := models.Asset{}
	res, err := Db.AddAssetToUser(newUser.ID, newAsset)
	if err != nil {
		t.Fatalf("failed during asset add operation")
	}
	if !res {
		t.Fatalf("asset already exists for this user")
	}
	Db.Release()
	res, err = Db.CheckUserAssetExistence(newUser.ID, newAsset.ID)
	if err == nil {
		t.Fatalf("storage resource is still available")
	}
	if res {
		t.Fatalf("asset added to user")
	}
}

// testMCUAssetAlreadyExists (test scenario)
func testMCUAssetAlreadyExists(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	_, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	newAsset := models.Asset{}
	res, err := Db.AddAssetToUser(newUser.ID, newAsset)
	if err != nil {
		t.Fatalf("failed during asset add operation")
	}
	if !res {
		t.Fatalf("asset already exists for this user")
	}
	res, err = Db.CheckUserAssetExistence(newUser.ID, newAsset.ID)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !res {
		t.Fatalf("asset already exists for this user")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testMCUAESuccess (test scenario)
func testMCUAESuccess(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	_, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	newAsset := models.Asset{}
	res, err := Db.AddAssetToUser(newUser.ID, newAsset)
	if err != nil {
		t.Fatalf("failed during asset add operation")
	}
	if !res {
		t.Fatalf("asset already exists for this user")
	}
	res, err = Db.CheckUserAssetExistence(newUser.ID, newAsset.ID)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !res {
		t.Fatalf("asset already exists for this user")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestMemoryCreateNewUser tests the possible outcomes from creating a new
// user.
func TestMemoryCreateNewUser(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"resource error": testMCNUResourceError,
		"user exists":    testMCNUUserExists,
		"user created":   testMCNUSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testMCNUResourceError (test scenario)
func testMCNUResourceError(t *testing.T) {
	ConnectToStorage("memory")
	Db.Release()
	newUser := models.User{}
	_, err := Db.CreateNewUser(newUser)
	if err == nil {
		t.Fatalf("storage resource is still available")
	}
}

// testMCNUUserExists (test scenario)
func testMCNUUserExists(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	check, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !check {
		t.Fatalf("user does not exist")
	}
	check, err = Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if check {
		t.Fatalf("user does not exist")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testMCNUSuccess (test scenario)
func testMCNUSuccess(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{}
	check, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !check {
		t.Fatalf("user already exists")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestCheckAdmin tests the possible outcomes from checking if a user is
// privileged.
func TestCheckAdmin(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"resource error":      testCAResourceError,
		"user does not exist": testCANotExist,
		"user is not admin":   testCASuccessFalse,
		"user is admin":       testCASuccessTrue,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testCAResourceError (test scenario)
func testCAResourceError(t *testing.T) {
	ConnectToStorage("memory")
	err := Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
	_, _, err = Db.CheckAdmin("testuser")
	if err == nil {
		t.Fatalf("storage resource is still available")
	}
}

// testCANotExist (test scenario)
func testCANotExist(t *testing.T) {
	ConnectToStorage("memory")
	found, _, err := Db.CheckAdmin("testuser")
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if found {
		t.Fatalf("user exists")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testCASuccessFalse (test scenario)
func testCASuccessFalse(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{
		ID:         "testUser",
		Password:   "testPass",
		Cookie:     "",
		Privileged: false,
		Created:    time.Now(),
	}
	res, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !res {
		t.Fatalf("failed to add example user")
	}
	found, check, err := Db.CheckAdmin("testUser")
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !found {
		t.Fatalf("user does not exists")
	}
	if check {
		t.Fatalf("user is privileged")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testCASuccessTrue (test scenario)
func testCASuccessTrue(t *testing.T) {
	ConnectToStorage("memory")
	newUser := models.User{
		ID:         "testUser",
		Password:   "testPass",
		Cookie:     "",
		Privileged: true,
		Created:    time.Now(),
	}
	res, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !res {
		t.Fatalf("failed to add example user")
	}
	found, check, err := Db.CheckAdmin("testUser")
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if !found {
		t.Fatalf("user does not exists")
	}
	if !check {
		t.Fatalf("user is not privileged")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestMemoryUpdateUserToken tests the possible outcomes from updating a
// user's token.
func TestMemoryUpdateUserToken(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"resource error":      testMUTResourceError,
		"user does not exist": testMUTUserNotExist,
		"user token updated":  testMUTUserTokenRetrieve,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testMUTResourceError (test scenario)
func testMUTResourceError(t *testing.T) {
	ConnectToStorage("memory")
	err := Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
	found, err := Db.UpdateUserToken("testuser", "testtoken")
	if err == nil {
		t.Fatalf("storage resource is still available")
	}
	if found {
		t.Fatalf("user exists")
	}
}

// testMUTUserNotExist (test scenario)
func testMUTUserNotExist(t *testing.T) {
	ConnectToStorage("memory")
	found, err := Db.UpdateUserToken("testuser", "testtoken")
	if err != nil {
		t.Fatalf("storage resource is still available")
	}
	if found {
		t.Fatalf("user exists")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testMUTUserTokenRetrieve (test scenario)
func testMUTUserTokenRetrieve(t *testing.T) {
	ConnectToStorage("memory")
	newUser, _ := MockStorageData()
	_, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("failure during user creation operation")
	}
	found, err := Db.UpdateUserToken("testuser", "testtoken")
	if err != nil {
		t.Fatalf("storage resource is still available")
	}
	if !found {
		t.Fatalf("user was not found")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestMemoryFetchUserToken tests the possible outcomes from requesting a
// user's token.
func TestMemoryFetchUserToken(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"resource error":       testFUTResourceError,
		"user does not exist":  testFUTUserNotExist,
		"user token retrieved": testFUTUserTokenRetrieve,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testFUTResourceError (test scenario)
func testFUTResourceError(t *testing.T) {
	ConnectToStorage("memory")
	Db.Release()
	_, _, err := Db.FetchUserToken("testuser")
	if err == nil {
		t.Fatalf("storage resource is still available")
	}
}

// testFUTUserNotExist (test scenario)
func testFUTUserNotExist(t *testing.T) {
	ConnectToStorage("memory")
	found, _, err := Db.FetchUserToken("testuser")
	if err != nil {
		t.Fatalf("failed during token retrieval")
	}
	if found {
		t.Fatalf("user found in storage")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testFUTUserTokenRetrieve (test scenario)
func testFUTUserTokenRetrieve(t *testing.T) {
	ConnectToStorage("memory")
	newUser, _ := MockStorageData()
	_, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("failed during creation of test user")
	}
	found, tok, err := Db.FetchUserToken("testuser")
	if err != nil {
		t.Fatalf("failed during token retrieval")
	}
	if !found {
		t.Fatalf("user not found in storage")
	}
	if tok != newUser.Cookie {
		t.Fatalf("user token retrieval wrong result: got %v, want %v",
			tok, newUser.Cookie)
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestMemoryFetchUserPass tests the possible outcomes from requesting a
// user's password.
func TestMemoryFetchUserPass(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"resource error":      testFUPResourceError,
		"user does not exist": testFUPUserNotExist,
		"user pass retrieved": testFUPUserPassRetrieve,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testFUPResourceError (test scenario)
func testFUPResourceError(t *testing.T) {
	ConnectToStorage("memory")
	Db.Release()
	_, _, err := Db.FetchUserPass("testuser")
	if err == nil {
		t.Fatalf("storage resource is still available")
	}
}

// testFUPUserNotExist (test scenario)
func testFUPUserNotExist(t *testing.T) {
	ConnectToStorage("memory")
	found, _, err := Db.FetchUserPass("testuser")
	if err != nil {
		t.Fatalf("failed during pass retrieval")
	}
	if found {
		t.Fatalf("user found in storage")
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testFUPUserPassRetrieve (test scenario)
func testFUPUserPassRetrieve(t *testing.T) {
	ConnectToStorage("memory")
	newUser, _ := MockStorageData()
	_, err := Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("failed during creation of test user")
	}
	found, pass, err := Db.FetchUserPass(newUser.ID)
	if err != nil {
		t.Fatalf("failed during pass retrieval")
	}
	if !found {
		t.Fatalf("user not found in storage")
	}
	if pass != newUser.Password {
		t.Fatalf("user pass retrieval wrong result: got %v, want %v",
			pass, newUser.Password)
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestMemoryFetchAllUsers tests the possible outcomes from requesting all the
// users.
func TestMemoryFetchAllUsers(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"users retrieved success": testFAUUsersRetrieved,
		"resource error":          testFAUResourceError,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testFAUResourceError (test scenario)
func testFAUResourceError(t *testing.T) {
	ConnectToStorage("memory")
	Db.Release()
	_, err := Db.FetchAllUsers()
	if err == nil {
		t.Fatalf("storage resource is still available")
	}
}

// testFAUUsersRetrieved (test scenario)
func testFAUUsersRetrieved(t *testing.T) {
	ConnectToStorage("memory")
	// Add users
	newUser, _ := MockStorageData()
	want := []models.User{newUser}
	Db.CreateNewUser(newUser)
	res, err := Db.FetchAllUsers()
	if err != nil {
		t.Fatalf("storage resource is not available")
	}
	if len(res) != len(want) {
		t.Fatalf("different number of users retrivied")
	}
	for i := range want {
		if want[i] != res[i] {
			t.Fatalf("fetch users result error: got %v want %v",
				res[i], want[i])
		}
	}
	err = Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestMemoryRelease checks the possible outcomes from releasing the storage
// resource.
func TestMemoryRelease(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"storage memory release success": testStorageReleaseSuccess,
		"storage memory resource error":  testStorageReleaseError,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testStorageReleaseSuccess (test scenario)
func testStorageReleaseSuccess(t *testing.T) {
	ConnectToStorage("memory")
	err := Db.Release()
	if err != nil {
		t.Fatalf("storage resource failed to be released")
	}
	_, _, err = Db.FetchUserPass("testuser")
	if err == nil {
		t.Fatalf("storage resource is still allocated")
	}
}

// testStorageReleaseSuccess (test scenario)
func testStorageReleaseError(t *testing.T) {
	ConnectToStorage("memory")
	err := Db.Release()
	if err != nil {
		t.Fatalf("failed to release resources")
	}
	// Release the resource again to trigger the error.
	err = Db.Release()
	if err == nil {
		t.Fatalf("storage resource is already released and should trigger an error")
	}
}
