package dto


type Column struct {
	ColumnName string
	ColumnNameJp string
	ColumnType string
	IsPrimaryKey bool
	IsNotNull bool
	IsReadOnly bool
}