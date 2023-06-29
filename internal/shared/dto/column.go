package dto


type Column struct {
	ColumnName string
	ColumnType string
	IsPrimaryKey bool
	IsNotNull bool
	IsReadOnly bool
}