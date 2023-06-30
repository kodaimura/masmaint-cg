package dto


type Column struct {
	ColumnName string
	ColumnType string
	IsPrimaryKey bool
	IsNotNull bool
	IsAuto bool
	IsReadOnly bool
}