package models

// OwnedAssets is the mapping between an Algorand account and assets owned by them
type OwnedAssets struct {
	ID       string   `pg:"default_gen_random_uuid()" json:"id"`
	Address  string   `json:"address"`
	AssetIds []string `json:"assetIds"`
}
