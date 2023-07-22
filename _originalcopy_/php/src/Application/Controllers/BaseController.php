<?php

declare(strict_types=1);

namespace App\Application\Controllers;

use Psr\Log\LoggerInterface;
use Psr\Container\ContainerInterface;

class BaseController
{
    protected ContainerInterface $container;
    protected LoggerInterface $logger;

    public function __construct(ContainerInterface $container, LoggerInterface $logger)
    {
        $this->container = $container;
        $this->logger = $logger;
    }
}