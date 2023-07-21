package generator

import (
	//"os"
	//"fmt"
	//"strings"

	"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/shared/dto"
	//"masmaint-cg/internal/shared/constant"
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

	if err := serv.generateEnv(); err != nil {
		return err
	}
	if err := serv.generateLogs(); err != nil {
		return err
	}
	if err := serv.generateVar(); err != nil {
		return err
	}
	if err := serv.generateSettingFiles(); err != nil {
		return err
	}
	/*
	if err := serv.generateTemplates(); err != nil {
		return err
	}
	if err := serv.generatePublic(); err != nil {
		return err
	}
	if err := serv.generateApp(); err != nil {
		return err
	}
	if err := serv.generateSrc(); err != nil {
		return err
	}
	*/
	return nil	
}

// env生成
func (serv *sourceGeneratorPhp) generateEnv() error {
	source := "_originalcopy_/php/env"
	destination := serv.path + "env/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// logs生成
func (serv *sourceGeneratorPhp) generateLogs() error {
	source := "_originalcopy_/php/logs"
	destination := serv.path + "logs/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// var生成
func (serv *sourceGeneratorPhp) generateVar() error {
	source := "_originalcopy_/php/var"
	destination := serv.path + "var/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// .htaccess/composer.json生成
func (serv *sourceGeneratorPhp) generateSettingFiles() error {
	source := "_originalcopy_/php/.htaccess"
	destination := serv.path + ".htaccess"

	err := CopyFile(source, destination)
	if err != nil {
		logger.LogError(err.Error())
		return err
	}

	source = "_originalcopy_/php/composer.json"
	destination = serv.path + "composer.json"

	err = CopyFile(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}