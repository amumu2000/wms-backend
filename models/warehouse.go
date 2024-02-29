package models

type Warehouse struct {
	ID            int64  `gorm:"column:id"`
	WarehouseName string `gorm:"column:warehouse_name"`
	Location      string `gorm:"column:location"`
	ManagerID     int64  `gorm:"column:manager_id"`
	Comment       string `gorm:"column:comment"`
}

func GetWarehouseByID(id int64) *Warehouse {
	var count int64
	warehouse := Warehouse{}
	db.Table("warehouse").Where("id = ?", id).Count(&count).First(&warehouse)

	if count == 0 {
		return nil
	} else {
		return &warehouse
	}
}

func FindWarehouses(id int64, warehouseName, location string, managerID int64) []Warehouse {
	query := db.Table("warehouse")

	if id >= 0 {
		query = query.Where("id = ?", id)
	}

	if warehouseName != "" {
		query = query.Where("warehouse_name LIKE ?", "%"+warehouseName+"%")
	}

	if location != "" {
		query = query.Where("location LIKE ?", "%"+location+"%")
	}

	if managerID >= 0 {
		query = query.Where("manager_id = ?", managerID)
	}

	var warehouses []Warehouse
	query.Find(&warehouses)

	return warehouses
}

func AddWarehouse(warehouseName, location string, managerID int64, comment string) int64 {
	warehouse := Warehouse{
		WarehouseName: warehouseName,
		Location:      location,
		ManagerID:     managerID,
		Comment:       comment,
	}

	db.Table("warehouse").Create(&warehouse)

	return warehouse.ID
}

func DeleteWarehouse(id []int64) {
	db.Table("warehouse").Delete(&Warehouse{}, id)
}

func EditWarehouse(warehouse Warehouse) {
	db.Table("warehouse").Save(&warehouse)
}
