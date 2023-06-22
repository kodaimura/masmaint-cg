package dto


type Column struct {
	ColumnName string
	ColumnType int
	IsPrimaryKey bool
	IsNotNullF bool
	IsUnique bool
}