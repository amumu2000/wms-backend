package models

type Good struct {
	ID       int64  `gorm:"column:id"`
	GoodName string `gorm:"column:good_name"`
	Category string `gorm:"column:category"`
	Comment  string `gorm:"column:comment"`
	Price    int    `gorm:"column:price"`
}

func GetGoodByID(id int64) *Good {
	var count int64
	good := Good{}
	db.Table("goods").Where("id = ?", id).Count(&count).First(&good)

	if count == 0 {
		return nil
	} else {
		return &good
	}
}

func FindGoods(goodID int64, goodName, category string, priceGreater, priceLess int,
	priceGreaterEqual, priceLessEqual bool) []Good {
	query := db.Table("goods")

	if goodID >= 0 {
		query = query.Where("id = ?", goodID)
	}

	if goodName != "" {
		query = query.Where("good_name LIKE ?", "%"+goodName+"%")
	}

	if category != "" {
		query = query.Where("category LIKE ?", "%"+category+"%")
	}

	if priceGreater >= 0 {
		op := ""

		if priceGreaterEqual {
			op = ">="
		} else {
			op = ">"
		}

		query = query.Where("price "+op+" ?", priceGreater)
	}

	if priceLess >= 0 {
		op := ""

		if priceLessEqual {
			op = "<="
		} else {
			op = "<"
		}

		query = query.Where("price "+op+" ?", priceLess)
	}

	var goods []Good
	query.Find(&goods)

	return goods
}

func AddGood(goodName, category, comment string, price int) int64 {
	good := Good{
		GoodName: goodName,
		Category: category,
		Comment:  comment,
		Price:    price,
	}

	db.Table("goods").Create(&good)

	return good.ID
}

func DeleteGood(id []int64) {
	db.Table("goods").Delete(&Good{}, id)
}

func EditGood(good Good) {
	db.Table("goods").Save(&good)
}
