package models

// AssetCreate is the structure received from web for creating a new asset
type AssetCreate struct {
	CreatorAddr   string `json:"creatorAddr"`
	AssetName     string `json:"assetName"`
	UnitName      string `json:"unitName"`
	Total         uint64 `json:"total"`
	Decimals      uint32 `json:"decimals"`
	DefaultFrozen bool   `json:"defaultFrozen"`
	URL           string `json:"url"`
	MetaDataHash  string `json:"metadataHash"`
	ManagerAddr   string `json:"managerAddr"`
	ReserveAddr   string `json:"reserveAddr"`
	FreezeAddr    string `json:"freezeAddr"`
	ClawbackAddr  string `json:"clawbackAddr"`
}

// AssetDetails contains details of an asset
type AssetDetails struct {
	AssetID string `json:"assetId"`
	AssetCreate
}
