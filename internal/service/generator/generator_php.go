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

	if err := serv.generateEnv(); err != nil {
		return err
	}
	if err := serv.generateLogs(); err != nil {
		return err
	}
	if err := serv.generateVar(); err != nil {
		return err
	}
	if err := serv.generateHtaccess(); err != nil {
		return err
	}
	if err := serv.generateComposerJson(); err != nil {
		return err
	}
	if err := serv.generateApp(); err != nil {
		return err
	}
	if err := serv.generateSrc(); err != nil {
		return err
	}
	if err := serv.generatePublic(); err != nil {
		return err
	}
	if err := serv.generateTemplates(); err != nil {
		return err
	}
	
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

// .htaccess生成
func (serv *sourceGeneratorPhp) generateHtaccess() error {
	source := "_originalcopy_/php/.htaccess"
	destination := serv.path + ".htaccess"

	err := CopyFile(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// composer.json生成
func (serv *sourceGeneratorPhp) generateComposerJson() error {
	source := "_originalcopy_/php/composer.json"
	destination := serv.path + "composer.json"

	err := CopyFile(source, destination)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// app生成
func (serv *sourceGeneratorPhp) generateApp() error {
	source := "_originalcopy_/php/app"
	destination := serv.path + "app/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateAppFiles(destination)
}

// app内のファイル生成
func (serv *sourceGeneratorPhp) generateAppFiles(path string) error {
	if err := serv.generateAppFileDependencies(path); err != nil {
		return err
	}
	if err := serv.generateAppFileMiddleware(path); err != nil {
		return err
	}
	if err := serv.generateAppFileSettings(path); err != nil {
		return err
	}
	if err := serv.generateAppFileRepositories(path); err != nil {
		return err
	}
	if err := serv.generateAppFileRoutes(path); err != nil {
		return err
	}
	return nil
}

// app/dependencies.php生成
func (serv *sourceGeneratorPhp) generateAppFileDependencies(path string) error {
	// _originalcopy_からコピー
	return nil
}

// app/middleware.php生成
func (serv *sourceGeneratorPhp) generateAppFileMiddleware(path string) error {
	// _originalcopy_からコピー
	return nil
}

// app/settings.php生成
func (serv *sourceGeneratorPhp) generateAppFileSettings(path string) error {
	// _originalcopy_からコピー
	return nil
}

// app/repositories.php生成
func (serv *sourceGeneratorPhp) generateAppFileRepositories(path string) error {
	code := "<?php\n\ndeclare(strict_types=1);\n\n"

	for _, table := range *serv.tables {
		code += fmt.Sprintf(
			"use App\\Application\\Models\\Daos\\%sDao;\n",
			SnakeToPascal(table.TableName),
		)
	}
	for _, table := range *serv.tables {
		code += fmt.Sprintf(
			"use App\\Application\\Models\\DaoImpls\\%sDaoImpl;\n",
			SnakeToPascal(table.TableName),
		)
	}
	code += "use DI\\ContainerBuilder;\n\n"
	code += "return function (ContainerBuilder $containerBuilder) {\n\n"
	code += "\t$containerBuilder->addDefinitions([\n"

	for _, table := range *serv.tables {
		tnp := SnakeToPascal(table.TableName)
		code += fmt.Sprintf("\t\t%sDao::class => \\DI\\autowire(%sDaoImpl::class),\n", tnp, tnp)
	}
	code += "\t]);\n};"

	err := WriteFile(fmt.Sprintf("%srepositories.php", path), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// app/routes.php生成
func (serv *sourceGeneratorPhp) generateAppFileRoutes(path string) error {
	code := "<?php\n\ndeclare(strict_types=1);\n\n"
	code += "use App\\Application\\Controllers\\IndexController;\n"

	for _, table := range *serv.tables {
		code += fmt.Sprintf(
			"use App\\Application\\Controllers\\%sController;\n",
			SnakeToPascal(table.TableName),
		)
	}
	code += "use Psr\\Http\\Message\\ResponseInterface as Response;\n"
	code += "use Psr\\Http\\Message\\ServerRequestInterface as Request;\n"
	code += "use Slim\\App;\n"
	code += "use Slim\\Interfaces\\RouteCollectorProxyInterface as Group;\n\n"
	code += "return function (App $app) {\n"
	code += "\t$app->options('/{routes:.*}', function (Request $request, Response $response) {\n"
	code += "\t\t// CORS Pre-Flight OPTIONS Request Handler\n"
	code += "\t\treturn $response;\n\t});\n\n"
	code += "\t$app->group('/mastertables', function (Group $group) {\n"
	code += "\t\t$group->get('', IndexController::class. ':indexPage');\n"
	code += "\t\t$group->get('/', IndexController::class. ':indexPage');\n"

	for _, table := range *serv.tables {
		tn := table.TableName
		tnp := SnakeToPascal(tn)
		tnc := SnakeToCamel(tn)
		code += fmt.Sprintf("\n\t\t$group->get('/%s', %sController::class. ':%sPage');\n", tn, tnp, tnc)
        code += fmt.Sprintf("\t\t$group->get('/api/%s', %sController::class. ':get%s');\n", tn, tnp, tnp)
        code += fmt.Sprintf("\t\t$group->post('/api/%s', %sController::class. ':post%s');\n", tn, tnp, tnp)
        code += fmt.Sprintf("\t\t$group->put('/api/%s', %sController::class. ':put%s');\n", tn, tnp, tnp)
        code += fmt.Sprintf("\t\t$group->delete('/api/%s', %sController::class. ':delete%s');\n", tn, tnp, tnp)
	}
	code += "\t});\n};"

	err := WriteFile(fmt.Sprintf("%sroutes.php", path), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// src生成
func (serv *sourceGeneratorPhp) generateSrc() error {
	source := "_originalcopy_/php/src"
	destination := serv.path + "src"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateApplication()
}

// Application生成
func (serv *sourceGeneratorPhp) generateApplication() error {
	if err := serv.generateHandlers(); err != nil {
		return err
	}
	if err := serv.generateMiddleware(); err != nil {
		return err
	}
	if err := serv.generateResponseEmitter(); err != nil {
		return err
	}
	if err := serv.generateSettings(); err != nil {
		return err
	}
	if err := serv.generateControllers(); err != nil {
		return err
	}
	if err := serv.generateServices(); err != nil {
		return err
	}
	if err := serv.generateModels(); err != nil {
		return err
	}

	return nil
}

// Handlers生成
func (serv *sourceGeneratorPhp) generateHandlers() error {
	//path := serv.path + "src/Application/Handlers/"
	// _originalcopy_からコピー
	return nil
}

// Middleware生成
func (serv *sourceGeneratorPhp) generateMiddleware() error {
	//path := serv.path + "src/Application/Middleware/"
	// _originalcopy_からコピー
	return nil
}

// ResponseEmitter生成
func (serv *sourceGeneratorPhp) generateResponseEmitter() error {
	//path := serv.path + "src/Application/Emitter/"
	// _originalcopy_からコピー
	return nil
}

// Settings生成
func (serv *sourceGeneratorPhp) generateSettings() error {
	//path := serv.path + "src/Application/Settings/"
	// _originalcopy_からコピー
	return nil
}

// Controllers生成
func (serv *sourceGeneratorPhp) generateControllers() error {
	path := serv.path + "src/Application/Controllers/"
	return serv.generateControllersFiles(path)
}

// Controllers内のファイル生成
func (serv *sourceGeneratorPhp) generateControllersFiles(path string) error {
	if err := serv.generateControllersFileBase(path); err != nil {
		return err
	}
	if err := serv.generateControllersFileIndex(path); err != nil {
		return err
	}
	for _, table := range *serv.tables {
		if err := serv.generateControllersFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

// Controllers/BaseController.php生成
func (serv *sourceGeneratorPhp) generateControllersFileBase(path string) error {
	// _originalcopy_からコピー
	return nil
}

// Controllers/IndexController.php生成
func (serv *sourceGeneratorPhp) generateControllersFileIndex(path string) error {
	// _originalcopy_からコピー
	return nil
}

const PHP_CONTROLLER_FORMAT =
`<?php

declare(strict_types=1);

namespace App\Application\Controllers;

use App\Application\Controllers\BaseController;
use App\Application\Services\%sService ;
use Psr\Log\LoggerInterface;
use Psr\Container\ContainerInterface;
use Psr\Http\Message\ResponseInterface as Response;
use Slim\Views\Twig;

class %sController extends BaseController
{

    private Twig $twig;
    protected %sService $%sService;

    public function __construct(ContainerInterface $container, LoggerInterface $logger, Twig $twig, %sService $%sService)
    {
        parent::__construct($container, $logger);
        $this->twig = $twig;
        $this->%sService = $%sService;
    }

    public function %sPage($request, $response, $args): Response
    {
        $response = $this->twig->render($response, '%s.html', []);
        return $response;
    }

    public function get%s($request, $response, $args): Response
    {
        $results = $this->%sService->getAll();
        $response->getBody()->write(json_encode($results));
        return $response->withHeader('Content-Type', 'application/json');
    }

    public function post%s($request, $response, $args): Response
    {
        $data = $request->getParsedBody();
        try {
            $result = $this->%sService->create($data);
            $response->getBody()->write(json_encode($result));

        } catch (\InvalidArgumentException $e) {
            $this->logger->debug($e->getMessage());
            return $response
            ->withHeader('Content-Type', 'application/json')
            ->withStatus(400);
            
        } catch (\Exception $e) {
            $this->logger->error($e->getMessage());
            return $response
            ->withHeader('Content-Type', 'application/json')
            ->withStatus(500);
        }
        return $response->withHeader('Content-Type', 'application/json');
    }

    public function put%s($request, $response, $args): Response
    {
        $data = $request->getParsedBody();
        try {
            $result = $this->%sService->update($data);
            $response->getBody()->write(json_encode($result));

        } catch (\InvalidArgumentException $e) {
            $this->logger->debug($e->getMessage());
            return $response
            ->withHeader('Content-Type', 'application/json')
            ->withStatus(400);

        } catch (\Exception $e) {
            $this->logger->error($e->getMessage());
            return $response
            ->withHeader('Content-Type', 'application/json')
            ->withStatus(500);
        }
        return $response->withHeader('Content-Type', 'application/json');
    }

    public function delete%s($request, $response, $args): Response
    {
        $data = $request->getParsedBody();
        try {
            $this->%sService->delete($data);

        } catch (\InvalidArgumentException $e) {
            $this->logger->debug($e->getMessage());
            return $response
            ->withHeader('Content-Type', 'application/json')
            ->withStatus(400);

        } catch (\Exception $e) {
            $this->logger->error($e->getMessage());
            return $response
            ->withHeader('Content-Type', 'application/json')
            ->withStatus(500);
        }
        return $response->withHeader('Content-Type', 'application/json');
    }

}
`

// Controllers/*Controller.php生成
func (serv *sourceGeneratorPhp) generateControllersFile(table *dto.Table, path string) error {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)

	code := fmt.Sprintf(
		PHP_CONTROLLER_FORMAT,
		tnp, tnp, tnp, tnc, tnp, tnc, tnc, tnc, tnc, tn, tnp, tnc, tnp, tnc, tnp, tnc, tnp, tnc,
	)

	err := WriteFile(fmt.Sprintf("%s%sController.php", path, tnp), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// Services生成
func (serv *sourceGeneratorPhp) generateServices() error {
	path := serv.path + "src/Application/Services/"
	return serv.generateServicesFiles(path)
}

// Services内のファイル生成
func (serv *sourceGeneratorPhp) generateServicesFiles(path string) error {
	if err := serv.generateServicesFileBase(path); err != nil {
		return err
	}

	for _, table := range *serv.tables {
		if err := serv.generateServicesFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

// Services/BaseService.php生成
func (serv *sourceGeneratorPhp) generateServicesFileBase(path string) error {
	// _originalcopy_からコピー
	return nil
}

const PHP_SERVICE_FORMAT = 
`<?php

declare(strict_types=1);

namespace App\Application\Services;

use App\Application\Models\Daos\%sDao;
use App\Application\Models\Entities\%s;
use Psr\Log\LoggerInterface;
use Psr\Container\ContainerInterface;
use Psr\Http\Message\ResponseInterface as Response;
use Slim\Views\Twig;

class %sService extends BaseService
{

    protected %sDao $%sDao;

    public function __construct(ContainerInterface $container, LoggerInterface $logger, %sDao $%sDao)
    {
        parent::__construct($container, $logger);
        $this->%sDao = $%sDao;
    }

    public function getAll(): array
    {
        return $this->%sDao->findAll();
    }

%s

%s

%s

}
`

// Services/*Service.php生成
func (serv *sourceGeneratorPhp) generateServicesFile(table *dto.Table, path string) error {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)

	code := fmt.Sprintf(
		PHP_SERVICE_FORMAT,
		tnp, tnp, tnp, tnp, tnc, tnp, tnc, tnc, tnc, tnc,
		serv.generateServicesFileCodeCreate(table),
		serv.generateServicesFileCodeUpdate(table),
		serv.generateServicesFileCodeDelete(table),
	)

	err := WriteFile(fmt.Sprintf("%s%sService.php", path, tnp), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// ServicesのCreateメソッドプログラム生成
func (serv *sourceGeneratorPhp) generateServicesFileCodeCreate(table *dto.Table) string {
	tnc := SnakeToCamel(table.TableName)
	tnp := SnakeToPascal(table.TableName)

	code := fmt.Sprintf("\tpublic function create($data): %s\n\t{\n", tnp)
	code += fmt.Sprintf("\t\t$%s = new %s();\n", tnc, tnp)

	for _, col := range table.Columns {
		if col.IsInsAble {
			cn := col.ColumnName
			cnp := SnakeToPascal(cn)
			code += fmt.Sprintf("\t\t$%s->set%s($data['%s']);\n", tnc, cnp, cn)
 		}
	}
	code += fmt.Sprintf("\n\t\treturn $this->%sDao->create($%s);\n\t}", tnc, tnc)
	return code
}

// ServicesのUpdateメソッドプログラム生成
func (serv *sourceGeneratorPhp) generateServicesFileCodeUpdate(table *dto.Table) string {
	tnc := SnakeToCamel(table.TableName)
	tnp := SnakeToPascal(table.TableName)

	code := fmt.Sprintf("\tpublic function update($data): %s\n\t{\n", tnp)
	code += fmt.Sprintf("\t\t$%s = new %s();\n", tnc, tnp)

	for _, col := range table.Columns {
		if col.IsPrimaryKey || col.IsUpdAble {
			cn := col.ColumnName
			cnp := SnakeToPascal(cn)
			code += fmt.Sprintf("\t\t$%s->set%s($data['%s']);\n", tnc, cnp, cn)
 		}
	}

	code += fmt.Sprintf("\n\t\treturn $this->%sDao->update($%s);\n\t}", tnc, tnc)
	return code
}

// ServicesのDeleteメソッドプログラム生成
func (serv *sourceGeneratorPhp) generateServicesFileCodeDelete(table *dto.Table) string {
	tnc := SnakeToCamel(table.TableName)
	tnp := SnakeToPascal(table.TableName)

	code := fmt.Sprintf("\tpublic function delete($data): %s\n\t{\n", tnp)
	code += fmt.Sprintf("\t\t$%s = new %s();\n", tnc, tnp)

	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			cn := col.ColumnName
			cnp := SnakeToPascal(cn)
			code += fmt.Sprintf("\t\t$%s->set%s($data['%s']);\n", tnc, cnp, cn)
 		}
	}

	code += fmt.Sprintf("\n\t\treturn $this->%sDao->delete($%s);\n\t}", tnc, tnc)
	return code
}

// Models生成
func (serv *sourceGeneratorPhp) generateModels() error {
	path := serv.path + "src/Application/Models/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	if err := serv.generateEntities(); err != nil {
		return err
	}
	if err := serv.generateDaos(); err != nil {
		return err
	}
	if err := serv.generateDaoImpls(); err != nil {
		return err
	}

	return nil
}

// Entities生成
func (serv *sourceGeneratorPhp) generateEntities() error {
	path := serv.path + "src/Application/Models/Entities/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateEntitiesFiles(path)
}

// Entities内のファイル生成
func (serv *sourceGeneratorPhp) generateEntitiesFiles(path string) error {
	for _, table := range *serv.tables {
		if err := serv.generateEntitiesFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

// CSVフォーマットのカラム型からentityフィールド用の型取得
func (serv *sourceGeneratorPhp) getEntityFieldType(col *dto.Column) string {
	isNotNull := col.IsNotNull
	isPrimaryKey := col.IsPrimaryKey
	colType := col.ColumnType

	if colType == "s" || colType == "t" {
		if isNotNull || isPrimaryKey {
			return "string"
		}
		return "?string"
	}
	if colType == "i" {
		if isNotNull || isPrimaryKey {
			return "int"
		}
		return "?int"
	}
	if colType == "f" {
		if isNotNull || isPrimaryKey {
			return "float"
		}
		return "?float"
	}
	return ""
}

// entityのセッタープログラム生成
func (serv *sourceGeneratorPhp) generateEntitiesFileCodeSetter(col *dto.Column) string {
	colType := serv.getEntityFieldType(col)
	cnp := SnakeToPascal(col.ColumnName)
	cnc := SnakeToCamel(col.ColumnName)

	code := fmt.Sprintf("\tpublic function set%s($%s)\n\t{\n", cnp, cnc)

	switch colType {
	case "string", "?string":
		code += fmt.Sprintf("\t\t$this->%s = $%s;\n", cnp, cnc)

	case "int":
		code += fmt.Sprintf("\t\tif ($%s === null || $%s === \"\") {\n", cnc, cnc) +
			fmt.Sprintf("\t\t\tthrow new \\InvalidArgumentException(\"error: set%s\");\n", cnp) +
			fmt.Sprintf("\t\t} else if (filter_var($%s, FILTER_VALIDATE_INT) !== false) {\n", cnc) +
			fmt.Sprintf("\t\t\t$this->%s = (int) $%s;\n", cnc, cnc) +
			fmt.Sprintf("\t\t} else {\n") + 
			fmt.Sprintf("\t\t\tthrow new \\InvalidArgumentException(\"error: set%s\");\n\t\t}\n", cnp)

	case "?int":
		code += fmt.Sprintf("\t\tif ($%s === null || $%s === \"\") {\n", cnc, cnc) +
			fmt.Sprintf("\t\t\t$this->%s = null;\n", cnc) +
			fmt.Sprintf("\t\t} else if (filter_var($%s, FILTER_VALIDATE_INT) !== false) {\n", cnc) +
			fmt.Sprintf("\t\t\t$this->%s = (int) $%s;\n", cnc, cnc) +
			fmt.Sprintf("\t\t} else {\n") + 
			fmt.Sprintf("\t\t\tthrow new \\InvalidArgumentException(\"error: set%s\");\n\t\t}\n", cnp)

	case "float":
		code += fmt.Sprintf("\t\tif ($%s === null || $%s === \"\") {\n", cnc, cnc) +
			fmt.Sprintf("\t\t\tthrow new \\InvalidArgumentException(\"error: set%s\");\n", cnp) +
			fmt.Sprintf("\t\t} else if (is_numeric($%s)) {\n", cnc) +
			fmt.Sprintf("\t\t\t$this->%s = (float) $%s;\n", cnc, cnc) +
			fmt.Sprintf("\t\t} else {\n") + 
			fmt.Sprintf("\t\t\tthrow new \\InvalidArgumentException(\"error: set%s\");\n\t\t}\n", cnp)
	case "?float":
		code += fmt.Sprintf("\t\tif ($%s === null || $%s === \"\") {\n", cnc, cnc) +
			fmt.Sprintf("\t\t\t$this->%s = null;\n", cnc) +
			fmt.Sprintf("\t\t} else if (is_numeric($%s)) {\n", cnc) +
			fmt.Sprintf("\t\t\t$this->%s = (float) $%s;\n", cnc, cnc) +
			fmt.Sprintf("\t\t} else {\n") + 
			fmt.Sprintf("\t\t\tthrow new \\InvalidArgumentException(\"error: set%s\");\n\t\t}\n", cnp)
	}

	code += "\t}\n\n"
	return code
}

// Entities/*.php生成
func (serv *sourceGeneratorPhp) generateEntitiesFile(table *dto.Table, path string) error {
	code := "<?php\n\ndeclare(strict_types=1);\n\nnamespace App\\Application\\Models\\Entities;\n\n" +
		fmt.Sprintf(
			"use JsonSerializable;\n\nclass %s implements JsonSerializable\n{\n", 
			SnakeToPascal(table.TableName),
		)

	//フィールド
	for _, col := range table.Columns {
		code += fmt.Sprintf("\tprivate %s $%s;\n\n", serv.getEntityFieldType(&col), SnakeToCamel(col.ColumnName))
	}

	//ゲッター
	for _, col := range table.Columns {
		code += fmt.Sprintf(
			"\tpublic function get%s(): %s\n\t{\n\t\treturn $this->%s;\n\t}\n\n",
			SnakeToPascal(col.ColumnName), serv.getEntityFieldType(&col), SnakeToCamel(col.ColumnName),
		)
	}

	//セッター
	for _, col := range table.Columns {
		code += serv.generateEntitiesFileCodeSetter(&col)
	}

	//jsonSerialize
	code += "\t//json_encode()でエンコードされるときに呼ばれる\n" +
		"\t#[\\ReturnTypeWillChange]\n" +
		"\tpublic function jsonSerialize(): array\n\t{\n\t\treturn [\n"
	for _, col := range table.Columns {
		code += fmt.Sprintf("\t\t\t'%s' => $this->%s,\n", col.ColumnName, SnakeToCamel(col.ColumnName))
	}

	code += "\t\t];\n\t}\n}"
	err := WriteFile(fmt.Sprintf("%s%s.php", path, SnakeToPascal(table.TableName)), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// Daos生成
func (serv *sourceGeneratorPhp) generateDaos() error {
	path := serv.path + "src/Application/Models/Daos/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateDaosFiles(path)
}

// Daos内のファイル生成
func (serv *sourceGeneratorPhp) generateDaosFiles(path string) error {
	for _, table := range *serv.tables {
		if err := serv.generateDaosFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

const PHP_DAO_FORMAT =
`<?php

declare(strict_types=1);

namespace App\Application\Models\Daos;

use App\Application\Models\Entities\%s;

interface %sDao
{

    public function findAll(): array;

    public function create(%s $%s): %s;

    public function update(%s $%s): %s;

    public function delete(%s $%s);
}
`
// Daos/*Dao.php生成
func (serv *sourceGeneratorPhp) generateDaosFile(table *dto.Table, path string) error {
	tnc := SnakeToCamel(table.TableName)
	tnp := SnakeToPascal(table.TableName)

	code := fmt.Sprintf(
		PHP_DAO_FORMAT,
		tnp, tnp, tnp, tnc, tnp, tnp, tnc, tnp, tnp, tnc,
	)
	err := WriteFile(fmt.Sprintf("%s%sDao.php", path, SnakeToPascal(table.TableName)), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// DaoImpls生成
func (serv *sourceGeneratorPhp) generateDaoImpls() error {
	path := serv.path + "src/Application/Models/DaoImpls/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateDaoImplsFiles(path)
}

// DaoImpls内のファイル生成
func (serv *sourceGeneratorPhp) generateDaoImplsFiles(path string) error {
	for _, table := range *serv.tables {
		if err := serv.generateDaoImplsFile(&table, path); err != nil {
			return err
		}
	}
	return nil
}

const PHP_DAOIMPL_FORMAT = 
`<?php

declare(strict_types=1);

namespace App\Application\Models\DaoImpls;

use App\Application\Models\Daos\%sDao;
use App\Application\Models\Entities\%s;

use \PDOException;
use \PDO;
use Psr\Log\LoggerInterface;

class %sDaoImpl implements %sDao
{

    private LoggerInterface $logger;
    private PDO $db;

    public function __construct(LoggerInterface $logger, PDO $db){
        $this->logger = $logger;
        $this->db = $db;
    }

%s

%s

%s

%s

%s
}
`
// DaoImplsの*DaoImpl.php生成
func (serv *sourceGeneratorPhp) generateDaoImplsFile(table *dto.Table, path string) error {
	tnp := SnakeToPascal(table.TableName)

	codeCreate := ""
	codeUpdate := ""

	if serv.rdbms == constant.MYSQL {
		// RETURNING が使えない場合
		codeCreate = serv.generateDaoImplsFileCodeCreate_MySQL(table)
		codeUpdate = serv.generateDaoImplsFileCodeUpdate_MySQL(table)
	} else {
		codeCreate = serv.generateDaoImplsFileCodeCreate(table)
		codeUpdate = serv.generateDaoImplsFileCodeUpdate(table)
	}

	code := fmt.Sprintf(
		PHP_DAOIMPL_FORMAT,
		tnp, tnp, tnp, tnp,
		serv.generateDaoImplsFileCodeFindAll(table),
		serv.generateDaoImplsFileCodeFindOne(table),
		codeCreate,
		codeUpdate,
		serv.generateDaoImplsFileCodeDelete(table),
	)

	err := WriteFile(fmt.Sprintf("%s%s.DaoImpl.php", path, tnp), code)
	if err != nil {
		logger.LogError(err.Error())
	}
	return err
}

// SELECT ORDER BYで指定するpkのカンマ区切り文字列生成
func (serv *sourceGeneratorPhp) concatPrimaryKeyWithCommas(cols []dto.Column) string {
	var ls []string
	for _, col := range cols {
		if col.IsPrimaryKey {
			ls = append(ls, col.ColumnName)
		}
	}
	return strings.Join(ls, ", ")
}

// DaoImplのfindAllメソッド生成
func (serv *sourceGeneratorPhp) generateDaoImplsFileCodeFindAll(table *dto.Table) string {
	tn := table.TableName
	tnp := SnakeToPascal(tn)

	code := "\tpublic function findAll(): array\n\t{\n"
	code += "\t\t$query = \n\t\t\t\"SELECT\n"

	for i, col := range table.Columns {
		if i == 0 {
			code += fmt.Sprintf("\t\t\t\t%s\n", col.ColumnName)
		} else {
			code += fmt.Sprintf("\t\t\t\t,%s\n", col.ColumnName)
		}
	}
	code += fmt.Sprintf(
		"\t\t\tFROM %s\n\t\t\tORDER BY %s ASC\";\n\n", 
		tn, serv.concatPrimaryKeyWithCommas(table.Columns),
	)
	code += "\t\ttry {\n\t\t\t$stmt = $this->db->prepare($query);\n" +
		"\t\t\t$stmt->execute();\n\t\t} catch (PDOException $e) {\n" +
		"\t\t\t$this->logger->error($e->getMessage());\n\t\t}\n\n"

	code += "\t\t$result = $stmt->fetchAll(PDO::FETCH_ASSOC);\n\t\t$ret = [];\n"
	code += "\t\tforeach ($result as $row) {\n"
	code += fmt.Sprintf("\t\t\t$x = new %s();\n", tnp)

	for _, col := range table.Columns {
		code += fmt.Sprintf("\t\t\t$x->set%s($row['%s']);\n", SnakeToPascal(col.ColumnName), col.ColumnName)
	}
	code += "\n\t\t\t$ret[] = $x;\n\t\t}\n\n\t\treturn $ret;\n\t}"

	return code
}

// DaoImplのfindOneメソッド生成
func (serv *sourceGeneratorPhp) generateDaoImplsFileCodeFindOne(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)

	code := fmt.Sprintf("\tpublic function findOne(%s, $%s): %s\n\t{\n", tnp, tnc, tnp)
	code += "\t\t$query = \n\t\t\t\"SELECT"

	for i, col := range table.Columns {
		if i == 0 {
			code += fmt.Sprintf("\n\t\t\t\t%s", col.ColumnName)
		} else {
			code += fmt.Sprintf("\n\t\t\t\t,%s", col.ColumnName)
		}
	}
	code += fmt.Sprintf("\n\t\t\tFROM %s\n\t\t\tWHERE ", tn)

	isFirst := true
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			if isFirst {
				code += fmt.Sprintf("%s = :%s", col.ColumnName, SnakeToCamel(col.ColumnName))
				isFirst = false
			} else {
				code += fmt.Sprintf("\n\t\t\t  AND %s = :%s", col.ColumnName, SnakeToCamel(col.ColumnName))
			}
		}
	}
	code += "\";\n\n"
	code += "\t\ttry {\n\t\t\t$stmt = $this->db->prepare($query);\n"

	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			code += fmt.Sprintf(
				"\t\t\t$stmt->bindValue(':%s', $%s->get%s());\n", 
				col.ColumnName, tnc, SnakeToPascal(col.ColumnName,
			))
		}
	}
	code += "\t\t\t$stmt->execute();\n\t\t} catch (PDOException $e) {\n" +
		"\t\t\t$this->logger->error($e->getMessage());\n\t\t}\n\n"

	code += "\t\t$result = $stmt->fetch(PDO::FETCH_ASSOC);\n\n"
	code += fmt.Sprintf("\t\t$ret = new %s();\n", tnp)

	for _, col := range table.Columns {
		code += fmt.Sprintf("\t\t$ret->set%s($result['%s']);\n", SnakeToPascal(col.ColumnName), col.ColumnName)
	}
	code += "\n\t\treturn $ret;\n\t}"

	return code
}

// DaoImplのcreateメソッド生成
func (serv *sourceGeneratorPhp) generateDaoImplsFileCodeCreate(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)

	code := fmt.Sprintf("\tpublic function create(%s, $%s): %s\n\t{\n", tnp, tnc, tnp)
	code += fmt.Sprintf("\t\t$query = \n\t\t\t\"INSERT INTO %s (\n", tn)

	isFirst := true
	for _, col := range table.Columns {
		if col.IsInsAble {
			if isFirst {
				code += fmt.Sprintf("\t\t\t\t%s\n", col.ColumnName)
				isFirst = false
			} else {
				code += fmt.Sprintf("\t\t\t\t,%s\n", col.ColumnName)
			}
		}
	}
	code += "\t\t\t) VALUES (\n"

	isFirst = true
	for _, col := range table.Columns {
		if col.IsInsAble {
			if isFirst {
				code += fmt.Sprintf("\t\t\t\t:%s\n", col.ColumnName)
				isFirst = false
			} else {
				code += fmt.Sprintf("\t\t\t\t,:%s\n", col.ColumnName)
			}
		}
	}
	code += "\t\t\t) RETURNING"

	for i, col := range table.Columns {
		if i == 0 {
			code += fmt.Sprintf("\n\t\t\t\t%s", col.ColumnName)
		} else {
			code += fmt.Sprintf("\n\t\t\t\t,%s", col.ColumnName)
		}
	}
	code += "\";\n\n"
	code += "\t\ttry {\n\t\t\t$stmt = $this->db->prepare($query);\n"

	for _, col := range table.Columns {
		if col.IsInsAble {
			code += fmt.Sprintf(
				"\t\t\t$stmt->bindValue(':%s', $%s->get%s());\n", 
				col.ColumnName, tnc, SnakeToPascal(col.ColumnName,
			))
		}
	}
	code += "\t\t\t$stmt->execute();\n\t\t} catch (PDOException $e) {\n" +
		"\t\t\t$this->logger->error($e->getMessage());\n\t\t}\n\n"

	code += "\t\t$result = $stmt->fetch(PDO::FETCH_ASSOC);\n\n"
	code += fmt.Sprintf("\t\t$ret = new %s();\n", tnp)

	for _, col := range table.Columns {
		code += fmt.Sprintf("\t\t$ret->set%s($result['%s']);\n", SnakeToPascal(col.ColumnName), col.ColumnName)
	}
	code += "\n\t\treturn $ret;\n\t}"

	return code
}

// DaoImplのupdateメソッド生成
func (serv *sourceGeneratorPhp) generateDaoImplsFileCodeUpdate(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)

	code := fmt.Sprintf("\tpublic function update(%s, $%s): %s\n\t{\n", tnp, tnc, tnp)
	code += fmt.Sprintf("\t\t$query = \n\t\t\t\"UPDATE %s SET\n", tn)

	isFirst := true
	for _, col := range table.Columns {
		if col.IsUpdAble {
			if isFirst {
				code += fmt.Sprintf("\t\t\t\t%s = :%s\n", col.ColumnName, SnakeToCamel(col.ColumnName))
				isFirst = false
			} else {
				code += fmt.Sprintf("\t\t\t\t,%s = :%s\n", col.ColumnName, SnakeToCamel(col.ColumnName))
			}
		}
	}
	code += "\t\t\tWHERE "

	isFirst = true
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			if isFirst {
				code += fmt.Sprintf("%s = :%s", col.ColumnName, SnakeToCamel(col.ColumnName))
				isFirst = false
			} else {
				code += fmt.Sprintf("\n\t\t\t  AND %s = :%s", col.ColumnName, SnakeToCamel(col.ColumnName))
			}
		}
	}
	code += "\n\t\t\tRETURNING"

	for i, col := range table.Columns {
		if i == 0 {
			code += fmt.Sprintf("\n\t\t\t\t%s", col.ColumnName)
		} else {
			code += fmt.Sprintf("\n\t\t\t\t,%s", col.ColumnName)
		}
	}
	code += "\";\n\n"
	code += "\t\ttry {\n\t\t\t$stmt = $this->db->prepare($query);\n"

	for _, col := range table.Columns {
		if col.IsUpdAble || col.IsPrimaryKey {
			code += fmt.Sprintf(
				"\t\t\t$stmt->bindValue(':%s', $%s->get%s());\n", 
				col.ColumnName, tnc, SnakeToPascal(col.ColumnName,
			))
		}
	}
	code += "\t\t\t$stmt->execute();\n\t\t} catch (PDOException $e) {\n" +
		"\t\t\t$this->logger->error($e->getMessage());\n\t\t}\n\n"

	code += "\n\t\treturn $ret;\n\t}"

	return code
}

// DaoImplのcreateメソッド生成(MYSQLの8.0.21以前用)
// RETURNINGを使わない
func (serv *sourceGeneratorPhp) generateDaoImplsFileCodeCreate_MySQL(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)

	code := fmt.Sprintf("\tpublic function create(%s, $%s): %s\n\t{\n", tnp, tnc, tnp)
	code += fmt.Sprintf("\t\t$query = \n\t\t\t\"INSERT INTO %s (\n", tn)

	isFirst := true
	for _, col := range table.Columns {
		if col.IsInsAble {
			if isFirst {
				code += fmt.Sprintf("\t\t\t\t%s\n", col.ColumnName)
				isFirst = false
			} else {
				code += fmt.Sprintf("\t\t\t\t,%s\n", col.ColumnName)
			}
		}
	}
	code += "\t\t\t) VALUES (\n"

	isFirst = true
	for _, col := range table.Columns {
		if col.IsInsAble {
			if isFirst {
				code += fmt.Sprintf("\t\t\t\t:%s\n", col.ColumnName)
				isFirst = false
			} else {
				code += fmt.Sprintf("\t\t\t\t,:%s\n", col.ColumnName)
			}
		}
	}
	code += "\t\t\t)\";\n\n"
	code += "\t\ttry {\n\t\t\t$stmt = $this->db->prepare($query);\n"

	for _, col := range table.Columns {
		if col.IsInsAble {
			code += fmt.Sprintf(
				"\t\t\t$stmt->bindValue(':%s', $%s->get%s());\n", 
				col.ColumnName, tnc, SnakeToPascal(col.ColumnName,
			))
		}
	}
	code += "\t\t\t$stmt->execute();\n\t\t} catch (PDOException $e) {\n" +
		"\t\t\t$this->logger->error($e->getMessage());\n\t\t}\n\n"

	code += fmt.Sprintf("\t\treturn $this->findOne($%s);\n\t}", tnc)

	return code
}

// DaoImplのupdateメソッド生成(MYSQLの8.0.21以前用)
// RETURNINGを使わない
func (serv *sourceGeneratorPhp) generateDaoImplsFileCodeUpdate_MySQL(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)

	code := fmt.Sprintf("\tpublic function update(%s, $%s): %s\n\t{\n", tnp, tnc, tnp)
	code += fmt.Sprintf("\t\t$query = \n\t\t\t\"UPDATE %s SET(\n", tn)

	isFirst := true
	for _, col := range table.Columns {
		if col.IsUpdAble {
			if isFirst {
				code += fmt.Sprintf("\t\t\t\t%s = :%s\n", col.ColumnName, SnakeToCamel(col.ColumnName))
				isFirst = false
			} else {
				code += fmt.Sprintf("\t\t\t\t,%s = :%s\n", col.ColumnName, SnakeToCamel(col.ColumnName))
			}
		}
	}
	code += "\t\t\tWHERE "

	isFirst = true
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			if isFirst {
				code += fmt.Sprintf("%s = :%s", col.ColumnName, SnakeToCamel(col.ColumnName))
				isFirst = false
			} else {
				code += fmt.Sprintf("\n\t\t\t  AND %s = :%s", col.ColumnName, SnakeToCamel(col.ColumnName))
			}
		}
	}
	code += "\";\n\n"
	code += "\t\ttry {\n\t\t\t$stmt = $this->db->prepare($query);\n"

	for _, col := range table.Columns {
		if col.IsUpdAble || col.IsPrimaryKey {
			code += fmt.Sprintf(
				"\t\t\t$stmt->bindValue(':%s', $%s->get%s());\n", 
				col.ColumnName, tnc, SnakeToPascal(col.ColumnName,
			))
		}
	}
	code += "\t\t\t$stmt->execute();\n\t\t} catch (PDOException $e) {\n" +
		"\t\t\t$this->logger->error($e->getMessage());\n\t\t}\n\n"

	code += fmt.Sprintf("\t\treturn $this->findOne($%s);\n\t}", tnc)

	return code
}

// DaoImplのdeleteメソッド生成
func (serv *sourceGeneratorPhp) generateDaoImplsFileCodeDelete(table *dto.Table) string {
	tn := table.TableName
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)

	code := fmt.Sprintf("\tpublic function delete(%s, $%s)\n\t{\n", tnp, tnc)
	code += fmt.Sprintf("\t\t\t$query = \n\t\t\t\"DELETE FROM %s\n\t\t\tWHERE ", tn)

	isFirst := true
	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			if isFirst {
				code += fmt.Sprintf("%s = :%s", col.ColumnName, SnakeToCamel(col.ColumnName))
				isFirst = false
			} else {
				code += fmt.Sprintf("\n\t\t\t  AND %s = :%s", col.ColumnName, SnakeToCamel(col.ColumnName))
			}
		}
	}
	code += "\";\n\n"
	code += "\t\ttry {\n\t\t\t$stmt = $this->db->prepare($query);\n"

	for _, col := range table.Columns {
		if col.IsPrimaryKey {
			code += fmt.Sprintf(
				"\t\t\t$stmt->bindValue(':%s', $%s->get%s());\n", 
				col.ColumnName, tnc, SnakeToPascal(col.ColumnName,
			))
		}
	}

	code += "\t\t\t$stmt->execute();\n\t\t} catch (PDOException $e) {\n" +
		"\t\t\t$this->logger->error($e->getMessage());\n\t\t}\n\n"

	code += "\t\treturn;\n\t}"

	return code
}

// public生成
func (serv *sourceGeneratorPhp) generatePublic() error {
	source := "_originalcopy_/php/public"
	destination := serv.path + "public/"

	err := CopyDir(source, destination)
	if err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateStatic()
}

// static生成
func (serv *sourceGeneratorPhp) generateStatic() error {
	if err := serv.generateStaticFileHtaccess(); err != nil {
		return err
	}
	if err := serv.generateStaticFileIndex(); err != nil {
		return err
	}
	if err := serv.generateCss(); err != nil {
		return err
	}
	if err := serv.generateJs(); err != nil {
		return err
	}

	return nil
}

// static/.htaccess生成
func (serv *sourceGeneratorPhp) generateStaticFileHtaccess() error {
	// _originalcopy_からコピー
	return nil
}

// static/index.php生成
func (serv *sourceGeneratorPhp) generateStaticFileIndex() error {
	// _originalcopy_からコピー
	return nil
}

// css生成
func (serv *sourceGeneratorPhp) generateCss() error {
	//path := serv.path + "public/static/css/"
	// _originalcopy_からコピー
	return nil
}

// js生成
func (serv *sourceGeneratorPhp) generateJs() error {
	path := serv.path + "public/static/js/"

	if err := os.MkdirAll(path, 0777); err != nil {
		logger.LogError(err.Error())
		return err
	}

	return serv.generateJsFiles(path)
}

// js/*.js生成
func (serv *sourceGeneratorPhp) generateJsFiles(path string) error {
	for _, table := range *serv.tables {
		code := GenerateJsCode(&table)
		if err := WriteFile(fmt.Sprintf("%s%s.js", path, table.TableName), code); err != nil {
			logger.LogError(err.Error())
			return err
		}
	}
	return nil
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

// templates/base.html生成
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

// templates/index.html生成
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

// templates/*.html生成
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