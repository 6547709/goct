package adapter

import (
	"context"
	"fmt"

	apiInfo "github.com/smartxworks/cloudtower-go-sdk/v2/client/api_info"
)

// About 调 SDK GetAPIVersion；CloudTower 端点 `/get-version` 返回裸字符串。
// 因此 TowerInfo.Build 在当前 API 下始终为空。
func (c *defaultClient) About(ctx context.Context) (TowerInfo, error) {
	p := apiInfo.NewGetAPIVersionParams().WithContext(ctx)
	res, err := c.api.APIInfo.GetAPIVersion(p)
	if err != nil {
		return TowerInfo{}, fmt.Errorf("get api version: %w", err)
	}
	if res == nil {
		return TowerInfo{}, fmt.Errorf("get api version: nil response")
	}
	return TowerInfo{Version: res.Payload}, nil
}
