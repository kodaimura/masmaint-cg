<?php

declare(strict_types=1);

namespace App\Application\Services;

use Psr\Log\LoggerInterface;
use Psr\Container\ContainerInterface;

class BaseService
{
    protected ContainerInterface $container;
    protected LoggerInterface $logger;

    public function __construct(ContainerInterface $container, LoggerInterface $logger)
    {
        $this->container = $container;
        $this->logger = $logger;
    }
}