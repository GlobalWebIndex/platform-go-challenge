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
	case domain.ChartAssetType:
		chs := []Chart{}
		err := gormQuery.Limit(query.Limit).Find(&chs).Error
		if err != nil {
			return nil, err
		}
		assets = listRowsToAssets(chs)
	case domain.AudienceAssetType:
		aus := []Audience{}
		err := gormQuery.Limit(query.Limit).Find(&aus).Error
		if err != nil {
			return nil, err
		}
		assets = listRowsToAssets(aus)
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
		Limit:   query.Limit,
		Assets:  assets,
	}
	return &dl, nil
}

func (d *DB) FavouriteAsset(ctx context.Context, userID, assetID uint, at domain.AssetType, isFavourite bool) (uint, error) {
	var nid uint = 0
	var err error
	switch at {
	case domain.InsightAssetType:
		if isFavourite {
			count := int64(0)
			d.db.Model(FavouriteInsight{}).Where("user_id = ? AND insight_id = ? ", userID, assetID).Count(&count)
			if count == 0 {
				in := &FavouriteInsight{UserID: userID, InsightID: assetID}
				err = d.db.Create(in).Error
				nid = in.ID
			} else {
				err = errors.New("record exists")
			}
		} else {
			in := &FavouriteInsight{}
			err = d.db.Where("user_id = ? AND insight_id = ? ", userID, assetID).Unscoped().Delete(in).Error
		}
	case domain.ChartAssetType:
		if isFavourite {
			count := int64(0)
			d.db.Model(FavouriteChart{}).Where("user_id = ? AND chart_id = ? ", userID, assetID).Count(&count)
			if count == 0 {
				ch := &FavouriteChart{UserID: userID, ChartID: assetID}
				err = d.db.Create(ch).Error
				nid = ch.ID
			} else {
				err = errors.New("record exists")
			}
		} else {
			in := &FavouriteChart{}
			err = d.db.Where("user_id = ? AND chart_id = ? ", userID, assetID).Unscoped().Delete(in).Error
		}
	case domain.AudienceAssetType:
		if isFavourite {
			count := int64(0)
			d.db.Model(FavouriteAudience{}).Where("user_id = ? AND audience_id = ? ", userID, assetID).Count(&count)
			if count == 0 {
				au := &FavouriteAudience{UserID: userID, AudienceID: assetID}
				err = d.db.Create(au).Error
				nid = au.ID
			} else {
				err = errors.New("record exists")
			}
		} else {
			in := &FavouriteAudience{}
			err = d.db.Where("user_id = ? AND audience_id = ? ", userID, assetID).Unscoped().Delete(in).Error
		}
	}
	if err != nil {
		return 0, err
	}

	return nid, nil
}

func (d *DB) ListFavouriteAssets(ctx context.Context, userID uint, query domain.QueryAssets) (*domain.ListedAssets, error) {
	var aus []Audience
	gormQuery := d.db.Model(FavouriteAudience{}).Select("audiences.*").
		Joins("JOIN users ON favourite_audiences.user_id = users.id AND favourite_audiences.user_id =? ", userID)
	if query.IsDesc {
		gormQuery = gormQuery.Joins("JOIN audiences ON favourite_audiences.audience_id = audiences.id AND audiences.id < ?", query.LastID).Order("audiences.id desc")
	} else {
		gormQuery = gormQuery.Joins("JOIN audiences ON favourite_audiences.audience_id = audiences.id AND audiences.id > ?", query.LastID)
	}

	gormQuery.Limit(query.Limit).Find(&aus)
	assets := listRowsToAssets(aus)
	var firstID uint = 0
	var lastID uint = 0
	if len(assets) > 0 {
		firstID = uint(assets[0].ID)
		lastID = uint(assets[len(assets)-1].ID)
	}
	la := domain.ListedAssets{
		FirstID: firstID,
		LastID:  lastID,
		Limit:   query.Limit,
		Type:    query.Type,
		Assets:  assets,
	}

	return &la, nil
}
