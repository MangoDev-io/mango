package models

// OwnedAssets is the mapping between an Algorand account and assets owned by them
type OwnedAssets struct {
	ID       string   `pg:"default_gen_random_uuid()"`
	Address  string   `json:"address"`
	AssetIDs []string `json:"assetIds"`
}
