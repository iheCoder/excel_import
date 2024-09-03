package correct_checker

import (
	util "excel_import/utils"
	"fmt"
	"testing"
)

type ResourceTestModel struct {
	ID           int     `json:"id" gorm:"column:id"`
	Name         string  `json:"name" gorm:"column:name"`
	ResourceType int32   `json:"resource_type" gorm:"column:resource_type"`
	ResourceID   int64   `json:"resource_id" gorm:"column:resource_id"`
	Sort         float64 `json:"sort" gorm:"column:sort"`
}

func (r *ResourceTestModel) TableName() string {
	return "resource"
}

func TestGetTableModels(t *testing.T) {
	db := util.InitDB()
	tx := db.Begin()
	defer tx.Rollback()

	slc := &SimpleLinkChecker{
		TableModel: &ResourceTestModel{},
		RangeWhere: "id <= 3",
	}

	models, err := slc.GetTableModels(tx)
	if err != nil {
		t.Fatal(err)
	}

	if len(models) != 3 {
		t.Fatalf("expected 3 models, got %d", len(models))
	}

	t.Log("done")
}

type DrinkTestModel struct {
	ID    int     `json:"id" gorm:"column:id"`
	Name  string  `json:"name" gorm:"column:name"`
	Price float64 `json:"price" gorm:"column:price"`
	Brand string  `json:"brand" gorm:"column:brand"`
}

func (d *DrinkTestModel) TableName() string {
	return "drink"
}

type ShoesTestModel struct {
	ID    int     `json:"id" gorm:"column:id"`
	Name  string  `json:"name" gorm:"column:name"`
	Price float64 `json:"price" gorm:"column:price"`
	Brand string  `json:"brand" gorm:"column:brand"`
	Size  int     `json:"size" gorm:"column:size"`
}

func (s *ShoesTestModel) TableName() string {
	return "shoes"
}

type BuildingTestModel struct {
	ID   int    `json:"id" gorm:"column:id"`
	Info string `json:"info" gorm:"column:info"`
}

type BuildingInfo struct {
	ConstructOrg string  `json:"construct_org" gorm:"column:construct_org"`
	Cost         float64 `json:"cost" gorm:"column:cost"`
	WorkerCount  int     `json:"worker_count" gorm:"column:worker_count"`
}

func (b *BuildingTestModel) TableName() string {
	return "building"
}

func TestCheckCorrect(t *testing.T) {
	db := util.InitDB()
	tx := db.Begin()
	defer tx.Rollback()

	ltw := func(model any) (any, string) {
		resource := model.(*ResourceTestModel)
		switch resource.ResourceType {
		case 1:
			return &DrinkTestModel{}, fmt.Sprintf("id = %d", resource.ResourceID)
		case 2:
			return &ShoesTestModel{}, fmt.Sprintf("id = %d", resource.ResourceID)
		case 3:
			return &BuildingTestModel{}, fmt.Sprintf("id = %d", resource.ResourceID)
		}

		return nil, ""
	}
	slc := NewSimpleLinkChecker(ltw, &ResourceTestModel{})

	// check correct
	// id in (4,5), it's correct
	slc.RangeWhere = "id in (4,5)"
	if err := slc.CheckCorrect(tx); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// id in (6), it's wrong
	slc.RangeWhere = "id in (6)"
	if err := slc.CheckCorrect(tx); err == nil {
		t.Fatalf("expected error, got nil")
	}

	t.Log("done")
}
