<?php

namespace io\metamessage\jsonc;

enum JsoncTokenType {
    case EOF;
    case LBrace;
    case RBrace;
    case LBracket;
    case RBracket;
    case Colon;
    case Comma;
    case String;
    case Number;
    case True;
    case False;
    case Null;
    case LeadingComment;
    case TrailingComment;
}

class JsoncToken {
    public function __construct(
        public readonly JsoncTokenType $type,
        public readonly string $literal = "",
        public readonly int $line = 0,
        public readonly int $column = 0
    ) {}
}
