<?php

declare(strict_types=1);

namespace App\Application\Controllers;

use App\Application\Controllers\BaseController;
use Psr\Log\LoggerInterface;
use Psr\Container\ContainerInterface;
use Psr\Http\Message\ResponseInterface as Response;
use Slim\Views\Twig;

class IndexController extends BaseController
{

    private Twig $twig;

    public function __construct(ContainerInterface $container, LoggerInterface $logger, Twig $twig)
    {
        parent::__construct($container, $logger);
        $this->twig = $twig;
    }

    public function indexPage($request, $response, $args): Response
    {
        $response = $this->twig->render($response, 'index.html', []);
        return $response;
    }

}