package presentation

import "platform-go-challenge/domain"

func CreateRoutes(fs domain.FavouriteService) map[string]domain.RouteDefinition {
	fh := NewFavouriteHandler(fs)

	return map[string]domain.RouteDefinition{
		"/assets/{asset_id}/starred": {
			Methods:     []string{"POST"},
			HandlerFunc: fh.AddToFavourites,
		},
		"/assets/starred": {
			Methods:     []string{"GET"},
			HandlerFunc: fh.GetFavourites,
		},
		"/assets/{asset_id}/edit": {
			Methods:     []string{"PUT"},
			HandlerFunc: fh.EditFavourites,
		},
		"/assets/{asset_id}/delete": {
			Methods:     []string{"DELETE"},
			HandlerFunc: fh.DeleteFavourite,
		},
	}
}
