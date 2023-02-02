package migrations

type AbstractTableModel interface {
	TableName() string
}
