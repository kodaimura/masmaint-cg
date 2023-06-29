package generator

import (
	//"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/shared/dto"
	//"masmaint-cg/internal/shared/constant"
)


type SourceGeneratorGolang struct {
	tables *[]dto.Table
	rdbms string
	path string
}

func NewSourceGeneratorGolang(tables *[]dto.Table, rdbms, path string) *SourceGeneratorGolang {
	return &SourceGeneratorGolang{
		tables, rdbms, path,
	}
}

func (serv *SourceGeneratorGolang) Generate() error {
	return nil
}