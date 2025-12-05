package main

// This is a map of maps containing the lists of favourites for all users of the system
// UserID -> map of (AssetID -> Asset)
var allUserFavourites = make(map[string]map[string]Asset)