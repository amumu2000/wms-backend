package models

type Inventory struct {
	ID          int64 `gorm:"column:id"`
	GoodsID     int64 `gorm:"column:goods_id"`
	Count       int   `gorm:"column:count"`
	WarehouseID int64 `gorm:"column:warehouse_id"`
}

func GetInventoryByID(id int64) *Inventory {
	var count int64
	inventory := Inventory{}
	db.Table("inventory").Where("id = ?", id).Count(&count).First(&inventory)

	if count == 0 {
		return nil
	} else {
		return &inventory
	}
}

func GetInventories(goodsID, whID int64) []Inventory {
	query := db.Table("inventory")

	if goodsID >= 0 {
		query = query.Where("goods_id = ?", goodsID)
	}

	if whID >= 0 {
		query = query.Where("warehouse_id = ?", whID)
	}

	var inventories []Inventory
	query.Find(&inventories)

	return inventories
}

func FindInventories(goodsID, whID int64, countGreater, countLess int, countGreaterEqual, countLessEqual bool) []Inventory {
	query := db.Table("inventory")

	if goodsID >= 0 {
		query = query.Where("goods_id = ?", goodsID)
	}

	if whID >= 0 {
		query = query.Where("warehouse_id = ?", whID)
	}

	if countGreater >= 0 {
		op := ""

		if countGreaterEqual {
			op = ">="
		} else {
			op = ">"
		}

		query = query.Where("count "+op+" ?", countGreater)
	}

	if countLess >= 0 {
		op := ""

		if countLessEqual {
			op = "<="
		} else {
			op = "<"
		}

		query = query.Where("count "+op+" ?", countLess)
	}

	var inventories []Inventory
	query.Find(&inventories)

	return inventories
}

func AddInventory(goodsID, whID int64, count int) int64 {
	inventory := Inventory{
		GoodsID:     goodsID,
		WarehouseID: whID,
		Count:       count,
	}

	db.Table("inventory").Create(&inventory)

	return inventory.ID
}

func DeleteInventory(id []int64) {
	db.Table("inventory").Delete(&Inventory{}, id)
}

func EditInventory(inventory Inventory) {
	db.Table("inventory").Save(&inventory)
}
