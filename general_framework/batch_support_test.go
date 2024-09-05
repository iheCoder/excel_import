package general_framework

import (
	util "excel_import/utils"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type ProductTestModel struct {
	ID         int       `gorm:"column:id"`
	Name       string    `gorm:"column:name"`
	Price      float64   `gorm:"column:price"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func (ProductTestModel) TableName() string {
	return "product"
}

func TestBatchSupportAddModelBatchInsert(t *testing.T) {
	count := 102

	// init db
	db := util.InitDB()
	tx := db.Begin()
	defer tx.Rollback()

	// query product count
	var productCount int64
	err := tx.Model(&ProductTestModel{}).Count(&productCount).Error
	if err != nil {
		t.Fatal(err)
	}

	// generate products
	products := generateProductsByCount(count)

	// batch insert
	batchSupport := newBatchSupportFeature(10)
	for _, product := range products {
		err = batchSupport.AddModel(tx, product)
		if err != nil {
			t.Fatal(err)
		}
	}

	// post handle
	err = batchSupport.PostHandle(tx)
	if err != nil {
		t.Fatal(err)
	}

	// check product count
	var newProductCount int64
	err = tx.Model(&ProductTestModel{}).Count(&newProductCount).Error
	if err != nil {
		t.Fatal(err)
	}

	if newProductCount != productCount+int64(count) {
		t.Fatalf("new product count %d != old product count %d + %d", newProductCount, productCount, count)
	}

	t.Log("batch insert success")
}

func generateProductsByCount(count int) []*ProductTestModel {
	products := make([]*ProductTestModel, 0, count)
	for i := 0; i < count; i++ {
		products = append(products, &ProductTestModel{
			Name:       "product" + strconv.Itoa(i),
			Price:      float64(rand.Intn(10000)+1) / 100,
			CreateTime: time.Now(),
		})
	}
	return products
}
