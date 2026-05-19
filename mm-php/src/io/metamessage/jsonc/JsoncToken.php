<?php

namespace io\metamessage\jsonc;

class JsoncTokenType
{
    const EOF = 0;

    const LBrace   = 1;
    const RBrace   = 2;
    const LBracket = 3;
    const RBracket = 4;
    const Colon    = 5;
    const Comma    = 6;

    const String = 7;
    const Number = 8;
    const True   = 9;
    const False  = 10;
    const Null   = 11;

    const LeadingComment  = 12;
    const TrailingComment = 13;
}

class LiteralType
{
    const LiteralUnknown = 0;
    const LiteralTrue    = 1;
    const LiteralFalse   = 2;
    const LiteralNull    = 3;
    const LiteralNumber  = 4;
    const LiteralString  = 5;
    const LiteralObject  = 6;
    const LiteralArray   = 7;
}

class JsoncToken
{
    public int $type;
    public string $literal;
    public int $line;
    public int $column;

    public function __construct(int $type = 0, string $literal = "", int $line = 1, int $column = 1)
    {
        $this->type    = $type;
        $this->literal = $literal;
        $this->line    = $line;
        $this->column  = $column;
    }
}