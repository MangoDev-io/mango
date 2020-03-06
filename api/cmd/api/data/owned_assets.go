package data

import "github.com/haardikk21/algorand-asset-manager/api/cmd/api/models"

// InsertNewAsset inserts a new asset into the database
func (s *DatabaseService) InsertNewAsset(creatorAddr, managerAddr, reserveAddr, freezeAddr, clawbackAddr, assetID string) error {
	var record models.OwnedAssets
	record.CreatorAddress = creatorAddr
	record.ManagerAddress = managerAddr
	record.ReserveAddress = reserveAddr
	record.FreezeAddress = freezeAddr
	record.ClawbackAddress = clawbackAddr
	record.AssetId = assetID

	_, err := s.Model(&record).
		Insert()

	if err != nil {
		return err
	}

	return nil
}

// UpdateAssetAddresses updates the mutable addresses linked to an asset
func (s *DatabaseService) UpdateAssetAddresses(creatorAddr, managerAddr, reserveAddr, freezeAddr, clawbackAddr, assetID string) error {
	var record models.OwnedAssets
	record.CreatorAddress = creatorAddr
	record.ManagerAddress = managerAddr
	record.ReserveAddress = reserveAddr
	record.FreezeAddress = freezeAddr
	record.ClawbackAddress = clawbackAddr
	record.AssetId = assetID

	_, err := s.Model(&record).
		Where("asset_id = ?asset_id").
		Update()

	if err != nil {
		return err
	}

	return nil
}

// SelectAllAssetsForAddress selects all assets where address is one of the addresses linked to it
func (s *DatabaseService) SelectAllAssetsForAddress(addr string) ([]*models.OwnedAssets, error) {
	var record []*models.OwnedAssets

	err := s.Model(&record).
		Where("creator_address = ?", addr).
		WhereOr("manager_address = ?", addr).
		WhereOr("reserve_address = ?", addr).
		WhereOr("freeze_address = ?", addr).
		WhereOr("clawback_address = ?", addr).
		Select()

	if err != nil {
		return nil, err
	}

	return record, nil
}
