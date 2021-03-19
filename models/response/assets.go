package response

import "github.com/timoth-y/iot-blockchain-contracts/models"

type AssetResponseItem struct {
	models.Asset
	Requirements *models.Requirements `json:"requirements,omitempty" metadata:",optional"`
}

type AssetsResponse struct {
	Items    []AssetResponseItem `json:"items"`
	ScrollID string              `json:"scroll_id"`
}
