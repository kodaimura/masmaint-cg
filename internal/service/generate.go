package service

import (
	"errors"
	"time"
	"os/exec"

	"masmaint-cg/pkg/utils"
	"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/shared/dto"
	"masmaint-cg/internal/shared/constant"
	"masmaint-cg/internal/service/generator"
)

type SourceGenerator interface {
	Generate() error
}


type GenerateService struct {}

func NewGenerateService() *GenerateService {
	return &GenerateService{}
}


func (serv *GenerateService) Generate(tables *[]dto.Table, lang, rdbms string) (string, error) {
	path := serv.generateSourcePath()

	if err := serv.generateSource(tables, lang, rdbms, path); err != nil {
		logger.LogError(err.Error())
		return "", err
	}

	pathZip := path + ".zip"
	if err := exec.Command("zip", "-r", pathZip, path).Run(); err != nil {
		logger.LogError(err.Error())
		return "", err
	}

	return pathZip, nil
}


func (serv *GenerateService) generateSourcePath() string {
	return "./tmp/masmaint-" + time.Now().Format("2006-01-02-15-04-05") + 
		"-" + utils.RandomString(10) + "/masmaint"
}


func (serv *GenerateService) generateSource(tables *[]dto.Table, lang, rdbms, path string) error {
	var sg SourceGenerator

	if lang == constant.GOLANG {
		sg = generator.NewSourceGeneratorGolang(tables, rdbms, path)
	} else {
		return errors.New("未対応言語");
	}

	return sg.Generate();
}