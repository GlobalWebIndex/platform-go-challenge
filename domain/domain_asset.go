package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/pariz/gountries"
)

func NewDomain(db IDBRepository) *Domain {
	return &Domain{
		validate: validator.New(),
		repo:     db,
	}
}

func (d *Domain) validateAsset(asset Asset) error {
	switch v := asset.Data.(type) {
	case *Insight:
		err := d.validate.Struct(v)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrWrongAssetInput, err)
		}

	case *Chart:
		err := d.validate.Struct(v)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrWrongAssetInput, err)
		}
		if len(v.Data.X) == 0 || len(v.Data.Y) == 0 {
			err := errors.New("data are empty")
			return fmt.Errorf("%w: %v", ErrWrongAssetInput, err)
		}

		if len(v.Data.X) != len(v.Data.Y) {
			err := errors.New("data are not equal")
			return fmt.Errorf("%w: %v", ErrWrongAssetInput, err)
		}
	case *Audience:
		err := d.validate.Struct(v)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrWrongAssetInput, err)
		}
		query := gountries.New()
		_, err = query.FindCountryByName(v.Country)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrWrongAssetInput, err)
		}
		if v.Gender != MaleGenderType && v.Gender != FemaleGenderType {
			err := errors.New("gender is not correct")
			return fmt.Errorf("%w: %v", ErrWrongAssetInput, err)
		}

	}
	return nil
}

func (d *Domain) AddAsset(ctx context.Context, user *User, asset Asset) (*Asset, error) {
	err := d.validateAsset(asset)
	if err != nil {
		return nil, err
	}
	if !user.IsAdmin {
		return nil, fmt.Errorf("%w: %v", ErrUnauthorized, errors.New("only administrators are authorized"))
	}

	newAsset, err := d.repo.AddAsset(ctx, asset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalDBFailure, err)
	}
	return newAsset, nil
}

func (d *Domain) UpdateAsset(ctx context.Context, user *User, asset Asset) (*Asset, error) {
	err := d.validateAsset(asset)
	if err != nil {
		return nil, err
	}
	if !user.IsAdmin {
		return nil, fmt.Errorf("%w: %v", ErrUnauthorized, errors.New("only administrators are authorized"))
	}
	newAsset, err := d.repo.UpdateAsset(ctx, asset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalDBFailure, err)
	}
	return newAsset, nil
}

func (d *Domain) DeleteAsset(ctx context.Context, user *User, assetID uint, assetType AssetType) error {
	if !user.IsAdmin {
		return fmt.Errorf("%w: %v", ErrUnauthorized, errors.New("only administrators are authorized"))
	}
	err := d.repo.DeleteAsset(ctx, assetType, assetID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInternalDBFailure, err)
	}
	return nil
}

func (d *Domain) GetAsset(ctx context.Context, user *User, assetID uint, assetType AssetType) (*Asset, error) {
	asset, err := d.repo.GetAsset(ctx, assetType, assetID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalDBFailure, err)
	}
	return asset, nil
}

func (d *Domain) ListAssets(ctx context.Context, user *User, query QueryAssets, favQuery *QueryFavouriteAssets) (*ListedAssets, error) {
	err := d.validate.Struct(query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWrongQueryInput, err)
	}

	if favQuery == nil {
		ls, err := d.repo.ListAssets(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInternalDBFailure, err)
		}
		return ls, nil
	}
	return nil, nil
}

func (d *Domain) FavouriteAsset(ctx context.Context, user *User, assetID uint, isFavourite bool) error {
	return nil
}
