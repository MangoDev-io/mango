package data

import (
	"strconv"

	"github.com/haardikk21/algorand-asset-manager/api/cmd/api/models"
)

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
func (s *DatabaseService) UpdateAssetAddresses(managerAddr, reserveAddr, freezeAddr, clawbackAddr, assetID string) error {
	var record models.OwnedAssets
	record.ManagerAddress = managerAddr
	record.ReserveAddress = reserveAddr
	record.FreezeAddress = freezeAddr
	record.ClawbackAddress = clawbackAddr
	record.AssetId = assetID

	_, err := s.Model(&record).
		Where("asset_id = ?asset_id").
		Column("manager_address", "reserve_address", "freeze_address", "clawback_address").
		Update()

	if err != nil {
		return err
	}

	return nil
}

// DeleteAssetByAssetID removes an OwnedAsset listing from the database
func (s *DatabaseService) DeleteAssetByAssetID(assetID uint64) error {
	var record models.OwnedAssets
	record.AssetId = strconv.FormatUint(assetID, 10)

	_, err := s.Model(&record).
		Where("asset_id = ?asset_id").
		Delete()

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
