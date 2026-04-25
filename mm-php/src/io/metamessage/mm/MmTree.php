<?php

namespace io\metamessage\mm;

class MmTree {
    public MmTag $tag;

    public function __construct(MmTag $tag) {
        $this->tag = $tag;
    }
}

namespace io\metamessage\mm\MmTree;

use io\metamessage\mm\MmTag;

class MmScalar extends \io\metamessage\mm\MmTree {
    public $data;
    public string $text;

    public function __construct($data, string $text, MmTag $tag) {
        parent::__construct($tag);
        $this->data = $data;
        $this->text = $text;
    }
}

class MmObject extends \io\metamessage\mm\MmTree {
    public array $fields;

    public function __construct(MmTag $tag, array $fields) {
        parent::__construct($tag);
        $this->fields = $fields;
    }
}

class MmArray extends \io\metamessage\mm\MmTree {
    public array $items;

    public function __construct(MmTag $tag, array $items) {
        parent::__construct($tag);
        $this->items = $items;
    }
}
