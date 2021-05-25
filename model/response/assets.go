package response

import "github.com/timoth-y/chainmetric-core/models"

// AssetResponseItem defines a single item structure for AssetsResponse.
type AssetResponseItem struct {
	models.Asset
	Requirements *models.Requirements `json:"requirements,omitempty" metadata:",optional"`
}

// AssetsResponse defines a structure of the assets query request.
type AssetsResponse struct {
	Items    []AssetResponseItem `json:"items"`
	ScrollID string              `json:"scroll_id"`
}
