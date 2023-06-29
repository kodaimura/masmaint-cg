package generator

import (
	//"os"

	//"masmaint-cg/internal/core/logger"
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
	if err := serv.generateSourceConfig(); err != nil {
		return err
	}
	if err := serv.generateSourceCore(); err != nil {
		return err
	}
	/*
	if err := serv.generateSourceController(); err != nil {
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
	source := "_originalcopy_/golang/cmd"
	destination := serv.path + "cmd/"
	
	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err

}

func (serv *SourceGeneratorGolang) generateSourceConfig() error {
	source := "_originalcopy_/golang/config"
	destination := serv.path + "config/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

func (serv *SourceGeneratorGolang) generateSourceCore() error {
	source := "_originalcopy_/golang/core"
	destination := serv.path + "core/"
	
	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}