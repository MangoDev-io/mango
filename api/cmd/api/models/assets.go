package models

// AssetCreate is the structure received from web for creating a new asset
type AssetCreate struct {
	CreatorAddr   string `json:"creatorAddr"`
	AssetName     string `json:"assetName"`
	UnitName      string `json:"unitName"`
	TotalIssuance uint64 `json:"totalIssuance"`
	Decimals      uint32 `json:"decimals"`
	DefaultFrozen bool   `json:"defaultFrozen"`
	URL           string `json:"url"`
	MetaDataHash  string `json:"metadataHash"`
	ManagerAddr   string `json:"managerAddr"`
	ReserveAddr   string `json:"reserveAddr"`
	FreezeAddr    string `json:"freezeAddr"`
	ClawbackAddr  string `json:"clawbackAddr"`
}

// AssetDestroy is the structure passed to the destroy transaction for destroying an asset
type AssetDestroy struct {
	AssetID     uint64 `json:"assetId"`
	ManagerAddr string `json:"managerAddr"`
}

// AssetFreeze is the structure passed to the freeze transaction for freezing an asset
type AssetFreeze struct {
	AssetID       uint64 `json:"assetId"`
	FreezeAddr    string `json:"freezeAddr"`
	TargetAddr    string `json:"targetAddr"`
	FreezeSetting bool   `json:"freezeSetting"`
}
