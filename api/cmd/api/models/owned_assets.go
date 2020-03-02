package models

// OwnedAssets is the mapping between an Algorand account and assets owned by them
type OwnedAssets struct {
	Address  string   `pg:",pk" json:"address"`
	AssetIds []string `json:"assetIds"`
}
