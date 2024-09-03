package correct_checker

import (
	util "excel_import/utils"
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
