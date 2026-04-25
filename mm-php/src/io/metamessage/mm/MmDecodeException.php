<?php

namespace io\metamessage\mm;

class MmDecodeException extends \Exception {
    public function __construct(string $message) {
        parent::__construct($message);
    }
}
