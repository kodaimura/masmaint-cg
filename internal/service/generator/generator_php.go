package generator

import (
	"os"
	"fmt"
	"strings"

	"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/shared/dto"
	"masmaint-cg/internal/shared/constant"
)


type sourceGeneratorPhp struct {
	tables *[]dto.Table
	rdbms string
	path string
}

func NewSourceGeneratorPhp(tables *[]dto.Table, rdbms, path string) *sourceGeneratorPhp {
	return &sourceGeneratorPhp{
		tables, rdbms, path,
	}
}

// PHPソース生成
func (serv *sourceGeneratorPhp) GenerateSource() error {
	return nil	
}