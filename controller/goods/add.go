package goods

import (
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type addReq struct {
	GoodName *string `json:"good_name"`
	Category *string `json:"category"`
	Comment  *string `json:"comment"`
	Price    *int    `json:"price"`
}

func Add(c *gin.Context) {
	ok, session := utils.CheckTokenWithSession(c, 1)
	if !ok {
		return
	}

	req := addReq{}

	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		utils.BadRequest(c)
		return
	}

	goodName := *req.GoodName
	category := *req.Category
	comment := *req.Comment
	price := *req.Price

	goodsID := models.AddGood(goodName, category, comment, price)

	models.AddLog("添加商品", fmt.Sprintf("已添加商品%s", goodName), session.UserID)

	utils.Status200(c, gin.H{
		"goods_id": goodsID,
	}, "")
}
