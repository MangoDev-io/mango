package models

// OwnedAssets is the mapping between an Algorand account and assets owned by them
type OwnedAssets struct {
	CreatorAddress  string `pg:",pk" json:"creatorAddr"`
	ManagerAddress  string `json:"managerAddr"`
	ReserveAddress  string `json:"reserveAddr`
	FreezeAddress   string `json:"freezeAddr"`
	ClawbackAddress string `json:"clawbackAddr"`
	AssetId         string `json:"assetId"`
}
