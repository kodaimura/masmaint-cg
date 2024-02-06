package service

import (
	"errors"
	"time"
	"os/exec"

	"masmaint-cg/internal/core/utils"
	"masmaint-cg/internal/core/logger"
	"masmaint-cg/internal/shared/dto"
	"masmaint-cg/internal/shared/constant"
	"masmaint-cg/internal/service/generator"
)

type SourceGenerator interface {
	GenerateSource() error
}


type generateService struct {}

func NewGenerateService() *generateService {
	return &generateService{}
}


func (serv *generateService) Generate(tables *[]dto.Table, lang, rdbms string) (string, error) {
	path := serv.generateSourcePath()

	if err := serv.generateSource(tables, lang, rdbms, path); err != nil {
		logger.LogError(err.Error())
		return "", err
	}

	pathZip := path + ".zip"
	if err := exec.Command("zip", "-rm", pathZip, path).Run(); err != nil {
		logger.LogError(err.Error())
		return "", err
	}

	return pathZip, nil
}


func (serv *generateService) generateSourcePath() string {
	return "./output/masmaint-" + time.Now().Format("2006-01-02-15-04-05") + 
		"-" + utils.RandomString(10)
}


func (serv *generateService) generateSource(tables *[]dto.Table, lang, rdbms, path string) error {
	var sg SourceGenerator

	if !(rdbms == constant.POSTGRESQL ||
		rdbms == constant.MYSQL ||
		rdbms == constant.SQLITE_3350) {
		return errors.New("未対応RDBMS");
	}

	if lang == constant.GOLANG {
		sg = generator.NewSourceGeneratorGolang(tables, rdbms, path + "/masmaint/")
	} else if lang == constant.PHP {
		sg = generator.NewSourceGeneratorPhp(tables, rdbms, path + "/masmaint/")
	} else {
		return errors.New("未対応言語");
	}

	return sg.GenerateSource();
}