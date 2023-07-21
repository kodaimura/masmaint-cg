<?php

declare(strict_types=1);

namespace App\Application\Handlers;

use Psr\Http\Message\ResponseInterface as Response;
use Slim\Handlers\ErrorHandler as SlimErrorHandler;
use Slim\Exception\HttpNotFoundException;

class HttpErrorHandler extends SlimErrorHandler
{
    /**
     * @inheritdoc
     */
    protected function respond(): Response
    {
        $exception = $this->exception;
        $statusCode = 500;

        if ($exception instanceof HttpException) {
            $statusCode = $exception->getCode();
        }

        $response = $this->responseFactory->createResponse($statusCode);
        $response->getBody()->write($exception->getMessage());

        return $response->withHeader('Content-Type', 'application/json');
    }
}