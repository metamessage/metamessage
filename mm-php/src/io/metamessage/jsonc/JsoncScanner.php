<?php

namespace io\metamessage\jsonc;

class JsoncScanner {
    private string $source;
    private int $length;
    private int $pos = 0;
    private int $line = 1;
    private int $col = 1;
    private bool $newLine = false;

    public function __construct(string $source) {
        $this->source = $source;
        $this->length = strlen($source);
    }

    private function peek(): string {
        if ($this->pos >= $this->length) {
            return "\0";
        }
        return $this->source[$this->pos];
    }

    private function next(): string {
        if ($this->pos >= $this->length) {
            return "\0";
        }
        $ch = $this->source[$this->pos];
        $this->pos++;

        if ($ch === "\n") {
            $this->newLine = true;
            $this->line++;
            $this->col = 1;
        } else {
            $this->col++;
        }

        return $ch;
    }

    private function skipWhitespace(): void {
        while ($this->pos < $this->length && ctype_space($this->source[$this->pos])) {
            $this->next();
        }
    }

    public function nextToken(): JsoncToken {
        $this->skipWhitespace();

        $ch = $this->peek();
        if ($ch === "\0") {
            return new JsoncToken(JsoncTokenType::EOF, "", $this->line, $this->col);
        }

        $startLine = $this->line;
        $startCol = $this->col;

        switch ($ch) {
            case '{':
                $this->next();
                return new JsoncToken(JsoncTokenType::LBrace, "{", $startLine, $startCol);
            case '}':
                $this->next();
                return new JsoncToken(JsoncTokenType::RBrace, "}", $startLine, $startCol);
            case '[':
                $this->next();
                return new JsoncToken(JsoncTokenType::LBracket, "[", $startLine, $startCol);
            case ']':
                $this->next();
                return new JsoncToken(JsoncTokenType::RBracket, "]", $startLine, $startCol);
            case ':':
                $this->next();
                $this->newLine = false;
                return new JsoncToken(JsoncTokenType::Colon, ":", $startLine, $startCol);
            case ',':
                $this->next();
                $this->newLine = false;
                return new JsoncToken(JsoncTokenType::Comma, ",", $startLine, $startCol);
            case '"':
                return $this->scanString();
            case '/':
                return $this->scanComment();
            default:
                return $this->scanLiteral();
        }
    }

    private function scanString(): JsoncToken {
        $startLine = $this->line;
        $startCol = $this->col;
        $this->next();

        $buf = '';
        while (true) {
            $ch = $this->next();
            if ($ch === "\0" || $ch === "\n") {
                break;
            }
            if ($ch === '"') {
                break;
            }
            if ($ch === '\\') {
                $buf .= $ch;
                $buf .= $this->next();
                continue;
            }
            $buf .= $ch;
        }

        return new JsoncToken(JsoncTokenType::String, $buf, $startLine, $startCol);
    }

    private function scanComment(): JsoncToken {
        $startLine = $this->line;
        $startCol = $this->col;
        $this->next();

        if ($this->peek() === '/') {
            $tokenType = !$this->newLine ? JsoncTokenType::TrailingComment : JsoncTokenType::LeadingComment;
            $this->next();
            $buf = '';
            while (true) {
                $ch = $this->peek();
                if ($ch === "\n" || $ch === "\0") {
                    break;
                }
                $buf .= $this->next();
            }
            return new JsoncToken($tokenType, trim($buf), $startLine, $startCol);
        }

        if ($this->peek() === '*') {
            $tokenType = !$this->newLine ? JsoncTokenType::TrailingComment : JsoncTokenType::LeadingComment;
            $this->next();
            $buf = '';
            while (true) {
                if ($this->peek() === "\0") {
                    break;
                }
                if ($this->peek() === '*' && $this->pos + 1 < $this->length && $this->source[$this->pos + 1] === '/') {
                    $this->next();
                    $this->next();
                    break;
                }
                $buf .= $this->next();
            }
            return new JsoncToken($tokenType, trim($buf), $startLine, $startCol);
        }

        return new JsoncToken(JsoncTokenType::EOF, "", $this->line, $this->col);
    }

    private function scanLiteral(): JsoncToken {
        $startLine = $this->line;
        $startCol = $this->col;
        $buf = '';

        while (true) {
            $ch = $this->peek();
            if ($ch === "\0" || str_contains(" \t\r\n,:{}[]", $ch)) {
                break;
            }
            $buf .= $this->next();
        }

        return match ($buf) {
            'true' => new JsoncToken(JsoncTokenType::True, $buf, $startLine, $startCol),
            'false' => new JsoncToken(JsoncTokenType::False, $buf, $startLine, $startCol),
            'null' => new JsoncToken(JsoncTokenType::Null, $buf, $startLine, $startCol),
            default => new JsoncToken(JsoncTokenType::Number, $buf, $startLine, $startCol),
        };
    }
}
