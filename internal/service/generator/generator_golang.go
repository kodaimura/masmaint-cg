package generator

import (
	"os"

	"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/shared/dto"
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

func (serv *SourceGeneratorGolang) GenerateSource() error {
	if err := serv.generateSourceCmd(); err != nil {
		return err
	}
	/*
	if err := serv.generateSourceConfig(); err != nil {
		return err
	}
	if err := serv.generateSourceController(); err != nil {
		return err
	}
	if err := serv.generateSourceCore(); err != nil {
		return err
	}
	if err := serv.generateSourceDto(); err != nil {
		return err
	}
	if err := serv.generateSourceModel(); err != nil {
		return err
	}
	if err := serv.generateSourceService(); err != nil {
		return err
	}
	if err := serv.generateSourceWeb(); err != nil {
		return err
	}
	*/
	return nil	
}

func (serv *SourceGeneratorGolang) generateSourceCmd() error {
	path := serv.path + "cmd/masmaint/"
	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateSourceCmdFile(path)
}

func (serv *SourceGeneratorGolang) generateSourceCmdFile(path string) error {
	code := "package main\n\n" + 
	"import (\n\t\"masmaint/core/server\"\n)\n\n" +
	"func main() {\n\tserver.Run()\n}"

	return WriteFile(path + "main.go", code)
}