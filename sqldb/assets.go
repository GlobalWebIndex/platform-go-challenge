package sqldb

import (
	"context"
	"errors"
	"platform-go-challenge/domain"
)

func (d *DB) AddAsset(ctx context.Context, asset domain.Asset) (*domain.Asset, error) {
	newAsset := &domain.Asset{}
	switch v := asset.Data.(type) {
	case *domain.Insight:
		in := &Insight{}
		in.FromDomain(v)
		err := d.db.Create(in).Error
		if err != nil {
			return nil, err
		}
		newAsset.ID = in.ID
		newAsset.Data = in.ToDomain()
	case *domain.Chart:
		ch := &Chart{}
		ch.FromDomain(v)
		err := d.db.Create(ch).Error
		if err != nil {
			return nil, err
		}
		newAsset.ID = ch.ID
		newAsset.Data = ch.ToDomain()
	case *domain.Audience:
		au := &Audience{}
		au.FromDomain(v)
		err := d.db.Create(au).Error
		if err != nil {
			return nil, err
		}
		newAsset.ID = au.ID
		newAsset.Data = au.ToDomain()
	default:
		return nil, errors.New("this asset type does not exist")
	}
	return newAsset, nil
}

func (d *DB) UpdateAsset(ctx context.Context, asset domain.Asset) (*domain.Asset, error) {
	if asset.ID == 0 {
		return nil, errors.New("add id ")
	}
	newAsset := &domain.Asset{}
	switch v := asset.Data.(type) {
	case *domain.Insight:
		in := &Insight{}
		d.db.First(in, asset.ID)
		in.FromDomain(v)
		err := d.db.Save(in).Error
		if err != nil {
			return nil, err
		}
		newAsset.ID = in.ID
		newAsset.Data = in.ToDomain()
	case *domain.Chart:
		ch := &Chart{}
		d.db.First(ch, asset.ID)
		ch.FromDomain(v)
		err := d.db.Save(ch).Error
		if err != nil {
			return nil, err
		}
		newAsset.ID = ch.ID
		newAsset.Data = ch.ToDomain()
	case *domain.Audience:
		au := &Audience{}
		d.db.First(au, asset.ID)
		au.FromDomain(v)
		err := d.db.Save(au).Error
		if err != nil {
			return nil, err
		}
		newAsset.ID = au.ID
		newAsset.Data = au.ToDomain()
	default:
		return nil, errors.New("this asset type does not exist")
	}

	return newAsset, nil
}

func (d *DB) GetAsset(ctx context.Context, at domain.AssetType, assetID uint) (*domain.Asset, error) {
	switch at {
	case domain.InsightAssetType:
	case domain.ChartAssetType:
	case domain.AudienceAssetType:

	}
	return nil, nil
}

func (d *DB) DeleteAsset(ctx context.Context, assetID uint) error {
	return nil
}

func (d *DB) ListAssets(ctx context.Context, userID uint, query domain.QueryAssets) error {
	return nil
}
func (d *DB) FavourAnAsset(ctx context.Context, userID, assetID uint) (uint, error) {
	return 0, nil
}
