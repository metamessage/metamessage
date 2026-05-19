<?php

namespace io\metamessage\core;

class MmDecodeException extends \Exception
{
    public function __construct(string $message)
    {
        parent::__construct($message);
    }
}
