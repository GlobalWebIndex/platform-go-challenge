package services

import (
	"errors"
	"platform-go-challenge/consts"
	"platform-go-challenge/domain"
	"platform-go-challenge/persistence/fakes"
	"reflect"
	"testing"
	"time"
)

func TestAddToFavourites(t *testing.T) {
	fr := &fakes.FakeFavouriteRepository{}
	ur := &fakes.FakeUserRepository{}
	ar := &fakes.FakeAssetRepository{}

	favourite := domain.Favourite{
		From:      "user/1",
		To:        "asset/2",
		CreatedAt: time.Now(),
	}
	ur.GetUserReturns("user/1", nil)
	ar.GetAssetReturns("asset/2", nil)
	fr.AddToFavouritesReturns(favourite, nil)

	favouriteService := FavouriteService{fr: fr, ur: ur, ar: ar}

	actual, err := favouriteService.AddToFavourites("1", "2")
	if err != nil {
		t.Error(err)
		return
	}

	if actual.From != favourite.From {
		t.Errorf("It was expected to have %v as output,\n but we got %v instead", favourite.From, actual.From)
	}
	if actual.To != favourite.To {
		t.Errorf("It was expected to have %v as output,\n but we got %v instead", favourite.To, actual.To)
	}
}

func TestAddToFavourites_FAIL(t *testing.T) {
	fr := &fakes.FakeFavouriteRepository{}
	ur := &fakes.FakeUserRepository{}
	ar := &fakes.FakeAssetRepository{}

	favourite := domain.Favourite{}
	ur.GetUserReturns("", errors.New("document not found"))
	fr.AddToFavouritesReturns(favourite, errors.New("document not found"))

	favouriteService := FavouriteService{fr: fr, ur: ur, ar: ar}

	_, err := favouriteService.AddToFavourites("1", "2")

	if err.Error() != "document not found" {
		t.Errorf("Expected document not found error, got %v instead", err)
	}
}

func TestGetFavourites(t *testing.T) {
	fr := &fakes.FakeFavouriteRepository{}
	ur := &fakes.FakeUserRepository{}
	ar := &fakes.FakeAssetRepository{}

	user := consts.Users[0]
	favouriteAssets := consts.Charts

	ur.GetUserReturns("user/"+user.Key, nil)
	fr.GetFavouritesReturns(favouriteAssets, nil)

	favouriteService := FavouriteService{fr: fr, ur: ur, ar: ar}

	actualAssets, err := favouriteService.GetFavourites("user" + user.Key)
	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(actualAssets, favouriteAssets) {
		t.Errorf("It was expected to have %v as output,\n but we got %v instead", favouriteAssets, actualAssets)
		return
	}
}
