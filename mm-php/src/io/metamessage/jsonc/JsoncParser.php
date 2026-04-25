<?php

namespace io\metamessage\jsonc;

class JsoncParser {
    /** @var JsoncToken[] */
    private array $tokens;
    private int $pos = 0;
    /** @var JsoncToken[] */
    private array $pendingComments = [];

    public function __construct(array $tokens) {
        $this->tokens = $tokens;
    }

    private function peek(): JsoncToken {
        if ($this->pos >= count($this->tokens)) {
            return $this->tokens[count($this->tokens) - 1];
        }
        return $this->tokens[$this->pos];
    }

    private function next(): JsoncToken {
        $t = $this->peek();
        $this->pos++;
        return $t;
    }

    private function consumeComments(): ?JsoncTag {
        if (count($this->pendingComments) === 0) {
            return null;
        }

        $merged = new JsoncTag();

        foreach ($this->pendingComments as $c) {
            $t = $this->tagFromComment($c->literal);
            if ($t !== null) {
                $merged = $this->mergeTag($merged, $t);
            }
        }

        $this->pendingComments = [];
        return $merged;
    }

    public function parse(): JsoncNode {
        $tok = $this->peek();
        if ($tok->type === JsoncTokenType::EOF) {
            return new JsoncValue();
        }

        return match ($tok->type) {
            JsoncTokenType::LBrace => $this->parseObject(),
            JsoncTokenType::LBracket => $this->parseArray(),
            default => $this->parseValue(),
        };
    }

    private function parseObject(): JsoncObject {
        $this->next();
        $tag = $this->consumeComments();
        $obj = new JsoncObject($tag);

        while ($this->peek()->type !== JsoncTokenType::RBrace && $this->peek()->type !== JsoncTokenType::EOF) {
            $tok = $this->peek();
            if ($tok->type === JsoncTokenType::LeadingComment || $tok->type === JsoncTokenType::TrailingComment) {
                $this->next();
                continue;
            }

            $keyToken = $this->next();
            if ($keyToken->type !== JsoncTokenType::String) {
                break;
            }

            $this->next();

            $fieldPath = $obj->path !== "" ? $obj->path . "." . $keyToken->literal : $keyToken->literal;
            $value = $this->parseValue();
            $obj->fields[] = new JsoncField($keyToken->literal, $value);

            if ($this->peek()->type === JsoncTokenType::Comma) {
                $this->next();
            }
        }

        if ($this->peek()->type === JsoncTokenType::RBrace) {
            $this->next();
        }

        return $obj;
    }

    private function parseArray(): JsoncArray {
        $this->next();
        $tag = $this->consumeComments();
        $arr = new JsoncArray($tag);
        $index = 0;

        while ($this->peek()->type !== JsoncTokenType::RBracket && $this->peek()->type !== JsoncTokenType::EOF) {
            $tok = $this->peek();
            if ($tok->type === JsoncTokenType::LeadingComment || $tok->type === JsoncTokenType::TrailingComment) {
                $this->next();
                continue;
            }

            $itemPath = ($arr->path !== "" ? $arr->path : "") . "[" . $index . "]";
            $item = $this->parseValue();
            $arr->items[] = $item;
            $index++;

            if ($this->peek()->type === JsoncTokenType::Comma) {
                $this->next();
            }
        }

        if ($this->peek()->type === JsoncTokenType::RBracket) {
            $this->next();
        }

        return $arr;
    }

    private function parseValue(): JsoncNode {
        $tok = $this->peek();

        if ($tok->type === JsoncTokenType::LBrace) {
            return $this->parseObject();
        }

        if ($tok->type === JsoncTokenType::LBracket) {
            return $this->parseArray();
        }

        $actualToken = $this->next();
        $tag = $this->consumeComments() ?? new JsoncTag();

        return match ($actualToken->type) {
            JsoncTokenType::String => (function() use ($actualToken, $tag): JsoncValue {
                if ($tag->type === JsoncValueType::Unknown) {
                    $tag->type = JsoncValueType::String;
                }
                return new JsoncValue($actualToken->literal, $actualToken->literal, $tag);
            })(),
            JsoncTokenType::Number => (function() use ($actualToken, $tag): JsoncValue {
                if ($tag->type === JsoncValueType::Unknown) {
                    $tag->type = str_contains($actualToken->literal, ".") ? JsoncValueType::Float64 : JsoncValueType::Int;
                }
                $data = str_contains($actualToken->literal, ".")
                    ? floatval($actualToken->literal)
                    : intval($actualToken->literal);
                return new JsoncValue($data, $actualToken->literal, $tag);
            })(),
            JsoncTokenType::True => (function() use ($tag): JsoncValue {
                $tag->type = JsoncValueType::Bool;
                return new JsoncValue(true, "true", $tag);
            })(),
            JsoncTokenType::False => (function() use ($tag): JsoncValue {
                $tag->type = JsoncValueType::Bool;
                return new JsoncValue(false, "false", $tag);
            })(),
            JsoncTokenType::Null => (function() use ($tag): JsoncValue {
                $tag->isNull = true;
                return new JsoncValue(null, "null", $tag);
            })(),
            default => new JsoncValue(null, "", $tag),
        };
    }

    private function tagFromComment(string $comment): ?JsoncTag {
        $trimmed = trim($comment);
        if (!str_starts_with($trimmed, "mm:")) {
            return null;
        }
        $tagStr = trim(substr($trimmed, 3));
        if ($tagStr === "") {
            return null;
        }
        return $this->parseTag($tagStr);
    }

    private function parseTag(string $tagStr): JsoncTag {
        $tag = new JsoncTag();
        $parts = array_map('trim', explode(';', $tagStr));

        foreach ($parts as $part) {
            if ($part === "") {
                continue;
            }

            $kv = explode('=', $part, 2);
            $key = strtolower($kv[0]);
            $value = count($kv) > 1 ? trim($kv[1]) : "";

            switch ($key) {
                case "is_null":
                    $tag->isNull = true;
                    $tag->nullable = true;
                    break;
                case "example":
                    break;
                case "desc":
                    $tag->desc = $value;
                    break;
                case "type":
                    $tag->type = match ($value) {
                        "str" => JsoncValueType::String,
                        "i" => JsoncValueType::Int,
                        "i8" => JsoncValueType::Int8,
                        "i16" => JsoncValueType::Int16,
                        "i32" => JsoncValueType::Int32,
                        "i64" => JsoncValueType::Int64,
                        "u" => JsoncValueType::Uint,
                        "u8" => JsoncValueType::Uint8,
                        "u16" => JsoncValueType::Uint16,
                        "u32" => JsoncValueType::Uint32,
                        "u64" => JsoncValueType::Uint64,
                        "f32" => JsoncValueType::Float32,
                        "f64" => JsoncValueType::Float64,
                        "bool" => JsoncValueType::Bool,
                        "bytes" => JsoncValueType::Bytes,
                        "bi" => JsoncValueType::BigInt,
                        "datetime" => JsoncValueType::DateTime,
                        "date" => JsoncValueType::Date,
                        "time" => JsoncValueType::Time,
                        "uuid" => JsoncValueType::UUID,
                        "decimal" => JsoncValueType::Decimal,
                        "ip" => JsoncValueType::IP,
                        "url" => JsoncValueType::URL,
                        "email" => JsoncValueType::Email,
                        "enum" => JsoncValueType::Enum,
                        "arr" => JsoncValueType::Array,
                        "struct" => JsoncValueType::Struct,
                        default => JsoncValueType::Unknown,
                    };
                    break;
                case "nullable":
                    $tag->nullable = true;
                    break;
                case "raw":
                    $tag->raw = true;
                    break;
                case "allow_empty":
                    $tag->allowEmpty = true;
                    break;
                case "unique":
                    $tag->unique = true;
                    break;
                case "default":
                    $tag->defaultValue = $value;
                    break;
                case "min":
                    $tag->min = $value;
                    break;
                case "max":
                    $tag->max = $value;
                    break;
                case "size":
                    $tag->size = intval($value);
                    break;
                case "enum":
                    $tag->type = JsoncValueType::Enum;
                    $tag->enum = $value;
                    break;
                case "pattern":
                    $tag->pattern = $value;
                    break;
                case "location":
                    $tag->location = $value;
                    break;
                case "version":
                    $tag->version = intval($value);
                    break;
                case "mime":
                    $tag->mime = $value;
                    break;
            }
        }

        return $tag;
    }

    private function mergeTag(JsoncTag $a, JsoncTag $b): JsoncTag {
        if ($a->type !== JsoncValueType::Unknown) {
            $b->type = $a->type;
        }
        if ($a->desc !== "") {
            $b->desc = $a->desc;
        }
        if ($a->nullable) {
            $b->nullable = true;
        }
        if ($a->isNull) {
            $b->isNull = true;
        }
        if ($a->defaultValue !== "") {
            $b->defaultValue = $a->defaultValue;
        }
        if ($a->min !== "") {
            $b->min = $a->min;
        }
        if ($a->max !== "") {
            $b->max = $a->max;
        }
        if ($a->size !== 0) {
            $b->size = $a->size;
        }
        if ($a->enum !== "") {
            $b->enum = $a->enum;
        }
        if ($a->pattern !== "") {
            $b->pattern = $a->pattern;
        }
        if ($a->location !== "") {
            $b->location = $a->location;
        }
        if ($a->version !== 0) {
            $b->version = $a->version;
        }
        if ($a->mime !== "") {
            $b->mime = $a->mime;
        }
        return $b;
    }
}

function parseJsonc(string $source): JsoncNode {
    $scanner = new JsoncScanner($source);
    $tokens = [];
    while (true) {
        $tok = $scanner->nextToken();
        $tokens[] = $tok;
        if ($tok->type === JsoncTokenType::EOF) {
            break;
        }
    }
    $parser = new JsoncParser($tokens);
    return $parser->parse();
}
