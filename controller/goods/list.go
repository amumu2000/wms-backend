package goods

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type listRequest struct {
	GoodsID           *int64  `json:"goods_id"`
	Name              *string `json:"name"`
	Category          *string `json:"category"`
	PriceGreater      *int    `json:"price_greater"`
	PriceGreaterEqual *bool   `json:"price_greater_equal"`
	PriceLess         *int    `json:"price_less"`
	PriceLessEqual    *bool   `json:"price_less_equal"`
}

func List(c *gin.Context) {
	if !utils.CheckToken(c, 1) {
		return
	}

	req := listRequest{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	goodsID, name, category := int64(-1), "", ""
	priceGreater, priceGreaterEqual, priceLess, priceLessEqual := -1, false, -1, false

	if req.GoodsID != nil {
		goodsID = *req.GoodsID
	}

	if req.Name != nil {
		name = *req.Name
	}

	if req.Category != nil {
		category = *req.Category
	}

	if req.PriceGreater != nil {
		priceGreater = *req.PriceGreater
	}

	if req.PriceGreaterEqual != nil {
		priceGreaterEqual = *req.PriceGreaterEqual
	}

	if req.PriceLess != nil {
		priceLess = *req.PriceLess
	}

	if req.PriceLessEqual != nil {
		priceLessEqual = *req.PriceGreaterEqual
	}

	goods := models.FindGoods(goodsID, name, category, priceGreater, priceLess, priceGreaterEqual, priceLessEqual)

	data := make([]gin.H, 0)
	for _, good := range goods {
		inventories := models.GetInventories(good.ID, -1)
		inventoryData := make([]gin.H, 0)

		totalCount := 0

		for _, inventory := range inventories {
			inventoryData = append(inventoryData, gin.H{
				"warehouse_id": inventory.WarehouseID,
				"count":        inventory.Count,
			})

			totalCount += inventory.Count
		}

		data = append(data, gin.H{
			"id":          good.ID,
			"good_name":   good.GoodName,
			"category":    good.Category,
			"price":       good.Price,
			"comment":     good.Comment,
			"inventory":   inventoryData,
			"total_count": totalCount,
		})
	}

	utils.Status200(c, gin.H{
		"goods": data,
	}, "")
}
