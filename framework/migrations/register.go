package migrations

import (
	"fmt"

	"gorm.io/gorm"
)

var (
	MigrationsTable = &Migrations{}
)

type Migrations struct {
	Models []AbstractTableModel
}

func (m *Migrations) AddMigrationModel(models ...AbstractTableModel) {
	if len(models) > 0 {
		MigrationsTable.Models = append(MigrationsTable.Models, models...)
	}
}

func (m *Migrations) Migrate(db *gorm.DB) []error {
	errorList := make([]error, 0)
	for _, models := range m.Models {
		fmt.Println("    [!] Migrate: ", models.TableName())
		err := db.AutoMigrate(models)
		if err != nil {
			errMessage := fmt.Errorf("%s: %w", models.TableName(), err)
			errorList = append(errorList, errMessage)
		}
	}
	return errorList
}
