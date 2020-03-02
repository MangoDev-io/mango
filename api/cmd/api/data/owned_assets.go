package data

import "github.com/haardikk21/algorand-asset-manager/api/cmd/api/models"

func (s *DatabaseService) InsertNewAsset(addr, assetID string) error {
	var record models.OwnedAssets
	record.Address = addr

	_, err := s.Model(record).
		Where("address = ?address").
		Returning("*").
		SelectOrInsert()
	if err != nil {
		return err
	}

	record.AssetIds = append(record.AssetIds, assetID)

	err = s.Update(record)
	if err != nil {
		return err
	}

	return nil
}

func (s *DatabaseService) SelectAllAssetsForAddress(addr string) (*models.OwnedAssets, error) {
	var record models.OwnedAssets
	err := s.Model(&record).
		Where("address = ?", addr).
		Select()

	if err != nil {
		return nil, err
	}

	return &record, nil
}
