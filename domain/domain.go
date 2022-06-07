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

func (d *Domain) AddAsset(ctx context.Context, asset Asset) error {
	err := d.validateAsset(asset)
	if err != nil {
		return err
	}
	return nil
}

func (d *Domain) UpdateAsset(ctx context.Context, assetID uint, asset Asset) error {
	err := d.validateAsset(asset)
	if err != nil {
		return err
	}
	return nil
}

func (d *Domain) DeleteAsset(ctx context.Context, assetID uint) error {

	return nil
}

func (d *Domain) ListAssets(ctx context.Context, userID uint, query QueryAssets) (*ListedAssets, error) {
	err := d.validate.Struct(query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWrongQueryInput, err)
	}
	return nil, nil
}
func (d *Domain) FavourAnAsset(ctx context.Context, userID, assetID uint) error {
	return nil
}
func (d *Domain) CreateUser(ctx context.Context, user User) (*User, error) {
	err := d.validate.Struct(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWrongUserInput, err)
	}
	return nil, nil
}
func (d *Domain) LoginUser(ctx context.Context, cred LoginCredentials) error {
	err := d.validate.Struct(cred)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWrongLoginInput, err)
	}
	return nil
}
