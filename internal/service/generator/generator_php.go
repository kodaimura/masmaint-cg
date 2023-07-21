package generator

import (
	"os"
	"fmt"
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
	if err := serv.generateTemplates(); err != nil {
		return err
	}
	/*
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

// templates生成
func (serv *sourceGeneratorPhp) generateTemplates() error {
	path := serv.path + "templates/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateTemplatesFiles(path)
}

// templates内のファイル生成
func (serv *sourceGeneratorPhp) generateTemplatesFiles(path string) error {
	if err := serv.generateTemplatesFileBase(path); err != nil {
		return err
	}
	if err := serv.generateTemplatesFileIndex(path); err != nil {
		return err
	}
	for _, table := range *serv.tables {
		if err := serv.generateTemplatesFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

// templates内のbase.html生成
func (serv *sourceGeneratorPhp) generateTemplatesFileBase(path string) error {
	code := fmt.Sprintf(
		"%s\n\t{%s block content %s}{%s endblock %s}\n%s",
		GenerateHtmlCodeHeader(serv.tables),
		"%", "%", "%", "%",
		GenerateHtmlCodeFooter(),
	)

	err := WriteFile(path + "base.html", code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// templates内のindex.html生成
func (serv *sourceGeneratorPhp) generateTemplatesFileIndex(path string) error {
	code := fmt.Sprintf(
		"{%s extends 'base.html' %s}\n\n{%s block content %s}\n{%s endblock %s}",
		"%", "%", "%", "%", "%", "%",
	)

	err := WriteFile(path + "index.html", code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// templates内の*.html生成
func (serv *sourceGeneratorPhp) generateTemplatesFile(table *dto.Table, path string) error {
	code := fmt.Sprintf(
		"{%s extends 'base.html' %s}\n\n{%s block content %s}%s{%s endblock %s}",
		"%", "%", "%", "%", GenerateHtmlCodeMain(table), "%", "%",
	)

	err := WriteFile(fmt.Sprintf("%s%s.html", path, table.TableName), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}