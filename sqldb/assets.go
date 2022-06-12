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
	newAsset := &domain.Asset{}
	switch at {
	case domain.InsightAssetType:
		in := &Insight{}
		err := d.db.First(in, assetID).Error
		if err != nil {
			return nil, err
		}
		newAsset.ID = in.ID
		newAsset.Data = in.ToDomain()
	case domain.ChartAssetType:
		ch := &Chart{}
		err := d.db.First(ch, assetID).Error
		if err != nil {
			return nil, err
		}
		newAsset.ID = ch.ID
		newAsset.Data = ch.ToDomain()
	case domain.AudienceAssetType:
		au := &Audience{}
		err := d.db.First(au, assetID).Error
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

func (d *DB) DeleteAsset(ctx context.Context, at domain.AssetType, assetID uint) error {
	switch at {
	case domain.InsightAssetType:
		err := d.db.Delete(&Insight{}, assetID).Error
		if err != nil {
			return err
		}

	case domain.ChartAssetType:
		err := d.db.Delete(&Chart{}, assetID).Error
		if err != nil {
			return err
		}

	case domain.AudienceAssetType:
		err := d.db.Delete(&Audience{}, assetID).Error
		if err != nil {
			return err
		}

	default:
		return errors.New("this asset type does not exist")
	}
	return nil
}

func (d *DB) ListAssets(ctx context.Context, query domain.QueryAssets) (*domain.ListedAssets, error) {
	gormQuery := d.db
	if query.IsDesc {
		gormQuery = gormQuery.Where("id < ?", query.LastID).Order("id desc")
	} else {
		gormQuery = gormQuery.Where("id > ?", query.LastID)
	}
	assets := []domain.Asset{}
	switch query.Type {
	case domain.InsightAssetType:
		ins := []Insight{}
		err := gormQuery.Limit(query.Limit).Find(&ins).Error
		if err != nil {
			return nil, err
		}
		assets = listRowsToAssets(ins)
	}
	var firstID uint = 0
	var lastID uint = 0
	if len(assets) > 0 {
		firstID = uint(assets[0].ID)
		lastID = uint(assets[len(assets)-1].ID)
	}

	dl := domain.ListedAssets{
		FirstID: firstID,
		LastID:  lastID,
		Assets:  assets,
	}
	return &dl, nil
}
func (d *DB) FavouriteAsset(ctx context.Context, userID, assetID uint, isFavourite bool) (uint, error) {
	return 0, nil
}
