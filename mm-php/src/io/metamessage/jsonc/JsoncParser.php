<?php

namespace io\metamessage\jsonc;

use io\metamessage\ir\Tag;
use io\metamessage\ir\ValueType;
use io\metamessage\ir\Node;
use io\metamessage\ir\Value;
use io\metamessage\ir\Object_;
use io\metamessage\ir\Array_;
use io\metamessage\ir\Field;
use io\metamessage\core\CamelToSnake;

class JsoncParser
{
    private const MAX_DEPTH = 32;
    private const TRUE_STR = 'true';
    private const FALSE_STR = 'false';

    private const SIMPLE_STRINGS = [
        'code' => true,
        'message' => true,
        'data' => true,
        'success' => true,
        'error' => true,
        'unknown' => true,
        'page' => true,
        'limit' => true,
        'offset' => true,
        'total' => true,
        'id' => true,
        'name' => true,
        'description' => true,
        'type' => true,
        'version' => true,
        'status' => true,
        'url' => true,
        'create_time' => true,
        'update_time' => true,
        'delete_time' => true,
        'account' => true,
        'token' => true,
        'expire_time' => true,
        'key' => true,
        'value' => true,
    ];

    /** @var JsoncToken[] */
    private array $toks;
    private int $pos = 0;
    /** @var JsoncToken[] */
    private array $pending = [];
    private int $depth = 0;

    public function __construct(array $tokens)
    {
        $this->toks = $tokens;
    }

    private function peek(): JsoncToken
    {
        if ($this->pos >= count($this->toks)) {
            return new JsoncToken(JsoncTokenType::EOF);
        }
        return $this->toks[$this->pos];
    }

    private function next(): JsoncToken
    {
        $t = $this->peek();
        $this->pos++;
        return $t;
    }

    private function consumeCommentsFor(int $anchorLine): ?Tag
    {
        if (count($this->pending) === 0) {
            return null;
        }

        $last = $this->pending[count($this->pending) - 1];
        if ($anchorLine - $last->line > 1) {
            $this->pending = [];
            return null;
        }

        $out = null;
        foreach ($this->pending as $ct) {
            $parsed = self::parseCommentsToTag($ct->literal);
            if ($parsed === null) {
                continue;
            }
            $out = Tag::mergeTag($out, $parsed);
        }

        return $out;
    }

    public function parse(): ?Node
    {
        $val = null;
        while (true) {
            $tok = $this->peek();
            if ($tok->type === JsoncTokenType::EOF) {
                return $val;
            }

            if ($tok->type === JsoncTokenType::LeadingComment) {
                if (count($this->pending) > 0) {
                    $last = $this->pending[count($this->pending) - 1];
                    if ($tok->line - $last->line > 1) {
                        $this->pending = [];
                    }
                }
                $this->pending[] = $tok;
                $this->next();
                continue;
            }

            if ($tok->type === JsoncTokenType::TrailingComment) {
                if ($val !== null) {
                    $parsed = self::parseCommentsToTag($tok->literal);
                    if ($parsed !== null) {
                        self::mergeNodeTag($val, $parsed);
                    }
                }
                $this->next();
                continue;
            }

            $val = $this->parseValue('');
        }
    }

    private function parseValue(string $path): ?Node
    {
        while (true) {
            $tok = $this->next();
            switch ($tok->type) {
                case JsoncTokenType::EOF:
                    return null;

                case JsoncTokenType::LBrace:
                    return $this->parseObject($tok->line, $path);

                case JsoncTokenType::LBracket:
                    return $this->parseArray($tok->line, $path);

                case JsoncTokenType::String:
                    $tag = $this->consumeCommentsFor($tok->line);

                    $text = $tok->literal;

                    if ($tag === null) {
                        $tag = Tag::newTag();
                    }

                    if ($tag->type === ValueType::UNKNOWN) {
                        $tag->type = ValueType::STR;
                    }

                    if (isset(self::SIMPLE_STRINGS[$text])) {
                        $tag->type = ValueType::STR;
                        $result = $this->validateStr($tag, $text);
                        $data = $result[0];
                        $text = $result[1];
                    } else {
                        switch ($tag->type) {
                            case ValueType::STR:
                                if ($tag->isNull) {
                                    if ($text !== '') {
                                        throw new \Exception(sprintf('invalid string: %s, valid: ""', json_encode($text)));
                                    }
                                    $data = '';
                                } else {
                                    $result = $this->validateStr($tag, $text);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::BYTES:
                                if ($tag->isNull) {
                                    if ($text !== '') {
                                        throw new \Exception(sprintf('invalid bytes: %s, valid: ""', json_encode($text)));
                                    }
                                    $data = '';
                                } else {
                                    $decoded = base64_decode($text, true);
                                    if ($decoded === false) {
                                        throw new \Exception(sprintf('invalid base64 bytes %s', json_encode($text)));
                                    }
                                    $result = $this->validateBytes($tag, $decoded);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::DATETIME:
                                $tz = new \DateTimeZone('UTC');
                                if ($tag->locationHours !== 0) {
                                    $sign = $tag->locationHours >= 0 ? '+' : '-';
                                    $tz = new \DateTimeZone(sprintf('%s%02d:00', $sign, abs($tag->locationHours)));
                                }

                                if ($tag->isNull) {
                                    $defaultDt = new \DateTime('1970-01-01 00:00:00', $tz);
                                    $dtStr = $defaultDt->format('Y-m-d H:i:s');
                                    if ($text !== $dtStr) {
                                        throw new \Exception(sprintf('invalid datetime: %s, valid: %s', json_encode($text), json_encode($dtStr)));
                                    }
                                    $data = $defaultDt;
                                } else {
                                    $d = \DateTime::createFromFormat('Y-m-d H:i:s', $text, $tz);
                                    if ($d === false || $d->getLastErrors() !== false && !empty($d->getLastErrors()['warnings'])) {
                                        throw new \Exception(sprintf('invalid datetime %s', json_encode($text)));
                                    }
                                    $result = $this->validateDatetime($tag, $d);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::DATE:
                                $tz = new \DateTimeZone('UTC');
                                if ($tag->locationHours !== 0) {
                                    $sign = $tag->locationHours >= 0 ? '+' : '-';
                                    $tz = new \DateTimeZone(sprintf('%s%02d:00', $sign, abs($tag->locationHours)));
                                }

                                if ($tag->isNull) {
                                    $defaultDt = new \DateTime('1970-01-01 00:00:00', $tz);
                                    $dtStr = $defaultDt->format('Y-m-d');
                                    if ($text !== $dtStr) {
                                        throw new \Exception(sprintf('invalid date: %s, valid: %s', json_encode($text), json_encode($dtStr)));
                                    }
                                    $data = $defaultDt;
                                } else {
                                    $d = \DateTime::createFromFormat('Y-m-d', $text, $tz);
                                    if ($d === false || $d->getLastErrors() !== false && !empty($d->getLastErrors()['warnings'])) {
                                        throw new \Exception(sprintf('invalid date %s', json_encode($text)));
                                    }
                                    $result = $this->validateDate($tag, $d);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::TIME:
                                $tz = new \DateTimeZone('UTC');
                                if ($tag->locationHours !== 0) {
                                    $sign = $tag->locationHours >= 0 ? '+' : '-';
                                    $tz = new \DateTimeZone(sprintf('%s%02d:00', $sign, abs($tag->locationHours)));
                                }

                                if ($tag->isNull) {
                                    $defaultDt = new \DateTime('1970-01-01 00:00:00', $tz);
                                    $dtStr = $defaultDt->format('H:i:s');
                                    if ($text !== $dtStr) {
                                        throw new \Exception(sprintf('invalid time: %s, valid: %s', json_encode($text), json_encode($dtStr)));
                                    }
                                    $data = $defaultDt;
                                } else {
                                    $d = \DateTime::createFromFormat('H:i:s', $text, $tz);
                                    if ($d === false || $d->getLastErrors() !== false && !empty($d->getLastErrors()['warnings'])) {
                                        throw new \Exception(sprintf('invalid time %s', json_encode($text)));
                                    }
                                    $result = $this->validateTime($tag, $d);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::UUID:
                                if ($tag->isNull) {
                                    if ($text !== '') {
                                        throw new \Exception(sprintf('invalid uuid: %s, valid: ""', json_encode($text)));
                                    }
                                    $data = str_repeat("\x00", 16);
                                } else {
                                    $result = $this->validateUUID($tag, $text);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::DECIMAL:
                                if ($tag->isNull) {
                                    if ($text !== '') {
                                        throw new \Exception(sprintf('invalid decimal: %s, valid: ""', json_encode($text)));
                                    }
                                    $data = '';
                                } else {
                                    $result = $this->validateDecimal($tag, $text);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::IP:
                                if ($tag->isNull) {
                                    if ($text !== '') {
                                        throw new \Exception(sprintf('invalid ip: %s, valid: ""', json_encode($text)));
                                    }
                                    $data = '';
                                } else {
                                    $ip = $text;
                                    $result = $this->validateIP($tag, $ip);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::URL:
                                if ($tag->isNull) {
                                    if ($text !== '') {
                                        throw new \Exception(sprintf('invalid url: %s, valid: ""', json_encode($text)));
                                    }
                                    $data = '';
                                } else {
                                    $parsedUrl = parse_url($text);
                                    if ($parsedUrl === false) {
                                        throw new \Exception(sprintf('invalid url %s', json_encode($text)));
                                    }
                                    $result = $this->validateURL($tag, $text);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::EMAIL:
                                if ($tag->isNull) {
                                    if ($text !== '') {
                                        throw new \Exception(sprintf('invalid email: %s, valid: ""', json_encode($text)));
                                    }
                                    $data = '';
                                } else {
                                    $result = $this->validateEmail($tag, $text);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::ENUM:
                                if ($tag->enumValues === '') {
                                    throw new \Exception('enum empty');
                                }

                                if ($tag->isNull) {
                                    if ($text !== '') {
                                        throw new \Exception(sprintf('invalid enum: %s, valid: ""', json_encode($text)));
                                    }
                                    $data = -1;
                                } else {
                                    $result = $this->validateEnum($tag, $text);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::IMAGE:
                                if ($tag->isNull) {
                                    if ($text !== '') {
                                        throw new \Exception(sprintf('invalid image: %s, valid: ""', json_encode($text)));
                                    }
                                    $data = '';
                                } else {
                                    $decoded = base64_decode($text, true);
                                    if ($decoded === false) {
                                        throw new \Exception(sprintf('invalid base64 image %s', json_encode($text)));
                                    }
                                    $result = $this->validateImage($tag, $decoded);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            default:
                                throw new \Exception(sprintf('unsupported type %s for string literal', $tag->type->wireName()));
                        }
                    }

                    $value = new Value();
                    $value->Data = $data;
                    $value->Text = $text;
                    $value->Tag = $tag;
                    $value->Path = $path;
                    return $value;

                case JsoncTokenType::Number:
                    $tag = $this->consumeCommentsFor($tok->line);

                    $text = $tok->literal;

                    if ($tag === null) {
                        $tag = Tag::newTag();
                    }

                    if (str_contains($text, '.')) {
                        if ($tag->type === ValueType::UNKNOWN) {
                            $tag->type = ValueType::F64;
                        }

                        switch ($tag->type) {
                            case ValueType::F32:
                                if ($tag->isNull) {
                                    if ($text !== '0.0') {
                                        throw new \Exception(sprintf('invalid float32: %s, valid: "0.0"', $text));
                                    }
                                    $data = (float) 0.0;
                                } else {
                                    $f64 = (float) $text;
                                    $result = $this->validateF32($tag, $f64);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::F64:
                                if ($tag->isNull) {
                                    if ($text !== '0.0') {
                                        throw new \Exception(sprintf('invalid float64: %s, valid: "0.0"', $text));
                                    }
                                    $data = (float) 0.0;
                                } else {
                                    $f64 = (float) $text;
                                    $result = $this->validateF64($tag, $f64);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            default:
                                throw new \Exception(sprintf('unsupported numeric type %s for float literal', $tag->type->wireName()));
                        }
                    } elseif (str_starts_with($text, '-')) {
                        if ($tag->type === ValueType::UNKNOWN) {
                            $tag->type = ValueType::I;
                        }

                        switch ($tag->type) {
                            case ValueType::I:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid int: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    $result = $this->validateI($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::I8:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid int8: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    if ($uv < -128 || $uv > 127) {
                                        throw new \Exception(sprintf('invalid int8 %s', json_encode($text)));
                                    }
                                    $result = $this->validateI8($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::I16:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid int16: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    if ($uv < -32768 || $uv > 32767) {
                                        throw new \Exception(sprintf('invalid int16 %s', json_encode($text)));
                                    }
                                    $result = $this->validateI16($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::I32:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid int32: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    $result = $this->validateI32($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::I64:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid int64: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    $result = $this->validateI64($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::BIGINT:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid bigint: %s, valid: "0"', $text));
                                    }
                                    $data = '0';
                                } else {
                                    if (!preg_match('/^-?\d+$/', $text)) {
                                        throw new \Exception(sprintf('invalid bigint %s', json_encode($text)));
                                    }
                                    $result = $this->validateBigint($tag, $text);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            default:
                                throw new \Exception(sprintf('unsupported numeric type %s for negative literal', $tag->type->wireName()));
                        }
                    } else {
                        if ($tag->type === ValueType::UNKNOWN) {
                            $tag->type = ValueType::I;
                        }

                        switch ($tag->type) {
                            case ValueType::I:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid int: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    $result = $this->validateI($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::I8:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid int8: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    if ($uv < -128 || $uv > 127) {
                                        throw new \Exception(sprintf('invalid int8 %s', json_encode($text)));
                                    }
                                    $result = $this->validateI8($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::I16:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid int16: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    if ($uv < -32768 || $uv > 32767) {
                                        throw new \Exception(sprintf('invalid int16 %s', json_encode($text)));
                                    }
                                    $result = $this->validateI16($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::I32:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid int32: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    $result = $this->validateI32($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::I64:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid int64: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    $result = $this->validateI64($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::U:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid uint: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    if ($uv < 0) {
                                        throw new \Exception(sprintf('invalid uint %s', json_encode($text)));
                                    }
                                    $result = $this->validateU($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::U8:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid uint8: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    if ($uv < 0 || $uv > 255) {
                                        throw new \Exception(sprintf('invalid uint8 %s', json_encode($text)));
                                    }
                                    $result = $this->validateU8($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::U16:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid uint16: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    if ($uv < 0 || $uv > 65535) {
                                        throw new \Exception(sprintf('invalid uint16 %s', json_encode($text)));
                                    }
                                    $result = $this->validateU16($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::U32:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid uint32: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    if ($uv < 0) {
                                        throw new \Exception(sprintf('invalid uint32 %s', json_encode($text)));
                                    }
                                    $result = $this->validateU32($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::U64:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid uint64: %s, valid: "0"', $text));
                                    }
                                    $data = 0;
                                } else {
                                    $uv = (int) $text;
                                    if ($uv < 0) {
                                        throw new \Exception(sprintf('invalid uint64 %s', json_encode($text)));
                                    }
                                    $result = $this->validateU64($tag, $uv);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            case ValueType::BIGINT:
                                if ($tag->isNull) {
                                    if ($text !== '0') {
                                        throw new \Exception(sprintf('invalid bigint: %s, valid: "0"', $text));
                                    }
                                    $data = '0';
                                } else {
                                    if (!preg_match('/^-?\d+$/', $text)) {
                                        throw new \Exception(sprintf('invalid bigint %s', json_encode($text)));
                                    }
                                    $result = $this->validateBigint($tag, $text);
                                    $data = $result[0];
                                    $text = $result[1];
                                }
                                break;

                            default:
                                throw new \Exception(sprintf('unsupported numeric type %s', $tag->type->wireName()));
                        }
                    }

                    $value = new Value();
                    $value->Data = $data;
                    $value->Text = $text;
                    $value->Tag = $tag;
                    $value->Path = $path;
                    return $value;

                case JsoncTokenType::True:
                    $tag = $this->consumeCommentsFor($tok->line);

                    if ($tag === null) {
                        $tag = Tag::newTag();
                    }
                    if ($tag->type === ValueType::UNKNOWN) {
                        $tag->type = ValueType::BOOL;
                    }

                    switch ($tag->type) {
                        case ValueType::BOOL:
                            if ($tag->isNull) {
                                throw new \Exception('bool must false when bool is null');
                            } else {
                                $this->validateBool($tag, true);
                            }
                            break;

                        default:
                            throw new \Exception(sprintf('unsupported type %s for boolean literal', $tag->type->wireName()));
                    }

                    $value = new Value();
                    $value->Data = true;
                    $value->Text = self::TRUE_STR;
                    $value->Tag = $tag;
                    $value->Path = $path;
                    return $value;

                case JsoncTokenType::False:
                    $tag = $this->consumeCommentsFor($tok->line);

                    if ($tag === null) {
                        $tag = Tag::newTag();
                    }
                    if ($tag->type === ValueType::UNKNOWN) {
                        $tag->type = ValueType::BOOL;
                    }

                    switch ($tag->type) {
                        case ValueType::BOOL:
                            if ($tag->isNull) {
                            } else {
                                $this->validateBool($tag, false);
                            }
                            break;

                        default:
                            throw new \Exception(sprintf('unsupported type %s for boolean literal', $tag->type->wireName()));
                    }

                    $value = new Value();
                    $value->Data = false;
                    $value->Text = self::FALSE_STR;
                    $value->Tag = $tag;
                    $value->Path = $path;
                    return $value;

                case JsoncTokenType::Null:
                    throw new \Exception('null literal is not supported');

                default:
                    throw new \Exception(sprintf('unexpected token %s', $tok->type->name));
            }
        }
    }

    private function parseObject(int $openLine, string $path): Object_
    {
        $this->depth++;
        if ($this->depth > self::MAX_DEPTH) {
            throw new \Exception(sprintf('max depth: %d', self::MAX_DEPTH));
        }

        $tag = $this->consumeCommentsFor($openLine);
        if ($tag === null) {
            $tag = Tag::newTag();
        }
        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::OBJECT;
        }

        if ($tag->name !== '') {
            if ($path === '') {
                $path = $tag->name;
            } else {
                $path = sprintf('%s.%s', $path, $tag->name);
            }
        }

        $obj = new Object_();
        $obj->Tag = $tag;
        $obj->Path = $path;
        $obj->Fields = [];

        $val = null;
        while (true) {
            $tok = $this->peek();
            if ($tok->type === JsoncTokenType::EOF) {
                break;
            }
            if ($tok->type === JsoncTokenType::RBrace) {
                $this->next();
                break;
            }

            if ($tok->type === JsoncTokenType::LeadingComment) {
                if (count($this->pending) > 0) {
                    $last = $this->pending[count($this->pending) - 1];
                    if ($tok->line - $last->line > 1) {
                        $this->pending = [];
                    }
                }
                $this->pending[] = $tok;
                $this->next();
                continue;
            }

            if ($tok->type === JsoncTokenType::TrailingComment) {
                if ($val !== null) {
                    $parsed = self::parseCommentsToTag($tok->literal);
                    if ($parsed !== null) {
                        self::mergeNodeTag($val, $parsed);
                    }
                }
                $this->next();
                continue;
            }

            $key = $this->next();
            if ($key->type !== JsoncTokenType::String) {
                throw new \Exception('expect string key');
            }
            $keyStr = CamelToSnake::convert($key->literal);

            $this->next();

            $pa = sprintf('%s.%s', $path, $keyStr);
            if ($tag->type === ValueType::MAP) {
                $pa = sprintf('%s[%s]', $path, $keyStr);
            }
            $val = $this->parseValue($pa);
            if ($val === null) {
                continue;
            }

            $childTag = $val->getTag();
            if ($childTag !== null && $tag !== null && $childTag->type === ValueType::MAP) {
                $childTag->inherit($tag);
            }

            $field = new Field();
            $field->Key = $keyStr;
            $field->Value = $val;
            $obj->Fields[] = $field;

            if ($this->peek()->type === JsoncTokenType::Comma) {
                $this->next();
            }
        }

        switch ($tag->type) {
            case ValueType::MAP:
                $this->validateMap($tag);
                break;

            case ValueType::OBJECT:
                $this->validateObj($tag);
                break;
        }

        return $obj;
    }

    private function parseArray(int $openLine, string $path): Array_
    {
        $this->depth++;
        if ($this->depth > self::MAX_DEPTH) {
            throw new \Exception(sprintf('max depth: %d', self::MAX_DEPTH));
        }

        $tag = $this->consumeCommentsFor($openLine);
        if ($tag === null) {
            $tag = Tag::newTag();
        }
        if ($tag->type === ValueType::UNKNOWN) {
            if ($tag->size > 0) {
                $tag->type = ValueType::ARR;
            } else {
                $tag->type = ValueType::VEC;
            }
        }

        if ($tag->name !== '') {
            $path = sprintf('%s.%s', $path, $tag->name);
        }

        $arr = new Array_();
        $arr->Tag = $tag;
        $arr->Path = $path;
        $arr->Items = [];

        $item = null;
        $i = 0;
        while (true) {
            $tok = $this->peek();
            if ($tok->type === JsoncTokenType::EOF) {
                break;
            }
            if ($tok->type === JsoncTokenType::RBracket) {
                $this->next();
                break;
            }

            if ($tok->type === JsoncTokenType::LeadingComment) {
                if (count($this->pending) > 0) {
                    $last = $this->pending[count($this->pending) - 1];
                    if ($tok->line - $last->line > 1) {
                        $this->pending = [];
                    }
                }
                $this->pending[] = $tok;
                $this->next();
                continue;
            }

            if ($tok->type === JsoncTokenType::TrailingComment) {
                if ($item !== null) {
                    $parsed = self::parseCommentsToTag($tok->literal);
                    if ($parsed !== null) {
                        self::mergeNodeTag($item, $parsed);
                    }
                }
                $this->next();
                continue;
            }

            $pa = sprintf('%s[%d]', $path, $i);
            $item = $this->parseValue($pa);
            if ($item === null) {
                continue;
            }
            $childTag = $item->getTag();
            if ($childTag !== null && $tag !== null) {
                $childTag->inherit($tag);
            }
            $arr->Items[] = $item;
            $i++;

            if ($this->peek()->type === JsoncTokenType::Comma) {
                $this->next();
            }
        }

        switch ($tag->type) {
            case ValueType::ARR:
                $this->validateArr($tag, $arr->Items);
                break;

            case ValueType::VEC:
                $this->validateVec($tag, $arr->Items);
                break;
        }

        return $arr;
    }

    private static function mergeNodeTag(Node $n, ?Tag $parsed): void
    {
        if ($n === null || $parsed === null) {
            return;
        }
        $existing = $n->getTag();
        $merged = Tag::mergeTag($existing, $parsed);
        if ($n instanceof Value) {
            $n->Tag = $merged;
        } elseif ($n instanceof Object_) {
            $n->Tag = $merged;
        } elseif ($n instanceof Array_) {
            $n->Tag = $merged;
        }
    }

    private static function parseCommentsToTag(string $cs): ?Tag
    {
        $trimmed = trim($cs);
        $lower = strtolower($trimmed);
        if (str_starts_with($lower, 'mm:')) {
            $after = substr($trimmed, 3);
            $parsed = Tag::parseMMTag($after);
            return $parsed;
        }
        return null;
    }

    private function validateStr(Tag $tag, string $val): array
    {
        if ($tag->isNull) {
            return [null, $val];
        }

        if ($val === '') {
            if (!$tag->allowEmpty) {
                throw new \Exception(sprintf('type string not allow empty value "%s"', $val));
            }
            return [$val, $val];
        }

        if ($tag->pattern !== '') {
            $re = @preg_match($tag->pattern, '');
            if ($re === false) {
                throw new \Exception(sprintf('pattern "%s" compile error', $tag->pattern));
            }
            if (!preg_match($tag->pattern, $val)) {
                throw new \Exception(sprintf('value "%s" does not match pattern %s', $val, $tag->pattern));
            }
        }

        $l = mb_strlen($val, 'UTF-8');

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min) {
                throw new \Exception('failed to parse t.Min as int');
            }
            if ($l < $mini) {
                throw new \Exception(sprintf('string length %d is less than the minimum limit %d', $l, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max) {
                throw new \Exception('failed to parse t.Max as int');
            }
            if ($l > $maxi) {
                throw new \Exception(sprintf('string length %d exceeds the maximum limit %d', $l, $maxi));
            }
        }

        if ($tag->size !== 0) {
            if ($l !== $tag->size) {
                throw new \Exception(sprintf('string length %d != size %d', $l, $tag->size));
            }
        }

        return [$val, $val];
    }

    private function validateBytes(Tag $tag, string $val): array
    {
        if ($tag->isNull) {
            return [null, $val];
        }

        $l = strlen($val);

        if ($l === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception('type []byte not allow empty value []byte{}');
            }
            return [$val, base64_encode($val)];
        }

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min) {
                throw new \Exception('failed to parse t.Min as int');
            }
            if ($l < $mini) {
                throw new \Exception(sprintf('[]byte length %d is less than the minimum limit %d', $l, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max) {
                throw new \Exception('failed to parse t.Max as int');
            }
            if ($l > $maxi) {
                throw new \Exception(sprintf('[]byte length %d exceeds the maximum limit %d', $l, $maxi));
            }
        }

        if ($tag->size !== 0) {
            if ($l !== $tag->size) {
                throw new \Exception(sprintf('[]byte length %d != size %d', $l, $tag->size));
            }
        }

        return [$val, base64_encode($val)];
    }

    private function validateBool(Tag $tag, bool $val): array
    {
        if ($tag->isNull) {
            return [null, $val ? self::TRUE_STR : self::FALSE_STR];
        }

        if ($tag->allowEmpty) {
            throw new \Exception('type bool not support allow empty');
        }

        return [$val, $val ? self::TRUE_STR : self::FALSE_STR];
    }

    private function validateI(Tag $tag, int $val): array
    {
        if ($tag->isNull) {
            return [null, (string) $val];
        }

        if ($val === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception(sprintf('type int not allow empty value %d', $val));
            }
            return [$val, (string) $val];
        }

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min) {
                throw new \Exception('failed to parse t.Min as int');
            }
            if ($val < $mini) {
                throw new \Exception(sprintf('value %d is less than the minimum limit %d', $val, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max) {
                throw new \Exception('failed to parse t.Max as int');
            }
            if ($val > $maxi) {
                throw new \Exception(sprintf('value %d exceeds the maximum limit %d', $val, $maxi));
            }
        }

        return [$val, (string) $val];
    }

    private function validateI8(Tag $tag, int $val): array
    {
        if ($tag->isNull) {
            return [null, (string) $val];
        }

        if ($val === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception(sprintf('type int8 not allow empty value %d', $val));
            }
            return [$val, (string) $val];
        }

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min || $mini < -128 || $mini > 127) {
                throw new \Exception('failed to parse t.Min as int8');
            }
            if ($val < $mini) {
                throw new \Exception(sprintf('value %d is less than the minimum limit %d', $val, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max || $maxi < -128 || $maxi > 127) {
                throw new \Exception('failed to parse t.Max as int8');
            }
            if ($val > $maxi) {
                throw new \Exception(sprintf('value %d exceeds the maximum limit %d', $val, $maxi));
            }
        }

        return [$val, (string) $val];
    }

    private function validateI16(Tag $tag, int $val): array
    {
        if ($tag->isNull) {
            return [null, (string) $val];
        }

        if ($val === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception(sprintf('type int16 not allow empty value %d', $val));
            }
            return [$val, (string) $val];
        }

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min || $mini < -32768 || $mini > 32767) {
                throw new \Exception('failed to parse t.Min as int16');
            }
            if ($val < $mini) {
                throw new \Exception(sprintf('value %d is less than the minimum limit %d', $val, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max || $maxi < -32768 || $maxi > 32767) {
                throw new \Exception('failed to parse t.Max as int16');
            }
            if ($val > $maxi) {
                throw new \Exception(sprintf('value %d exceeds the maximum limit %d', $val, $maxi));
            }
        }

        return [$val, (string) $val];
    }

    private function validateI32(Tag $tag, int $val): array
    {
        if ($tag->isNull) {
            return [null, (string) $val];
        }

        if ($val === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception(sprintf('type int32 not allow empty value %d', $val));
            }
            return [$val, (string) $val];
        }

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min) {
                throw new \Exception('failed to parse t.Min as int32');
            }
            if ($val < $mini) {
                throw new \Exception(sprintf('value %d is less than the minimum limit %d', $val, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max) {
                throw new \Exception('failed to parse t.Max as int32');
            }
            if ($val > $maxi) {
                throw new \Exception(sprintf('value %d exceeds the maximum limit %d', $val, $maxi));
            }
        }

        return [$val, (string) $val];
    }

    private function validateI64(Tag $tag, int $val): array
    {
        if ($tag->isNull) {
            return [null, (string) $val];
        }

        if ($val === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception(sprintf('type int64 not allow empty value %d', $val));
            }
            return [$val, (string) $val];
        }

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min) {
                throw new \Exception('failed to parse t.Min as int64');
            }
            if ($val < $mini) {
                throw new \Exception(sprintf('value %d is less than the minimum limit %d', $val, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max) {
                throw new \Exception('failed to parse t.Max as int64');
            }
            if ($val > $maxi) {
                throw new \Exception(sprintf('value %d exceeds the maximum limit %d', $val, $maxi));
            }
        }

        return [$val, (string) $val];
    }

    private function validateU(Tag $tag, int $val): array
    {
        if ($tag->isNull) {
            return [null, (string) $val];
        }

        if ($val === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception(sprintf('type uint not allow empty value %d', $val));
            }
            return [$val, (string) $val];
        }

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min || $mini < 0) {
                throw new \Exception('failed to parse t.Min as uint');
            }
            if ($val < $mini) {
                throw new \Exception(sprintf('value %d is less than the minimum limit %d', $val, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max || $maxi < 0) {
                throw new \Exception('failed to parse t.Max as uint');
            }
            if ($val > $maxi) {
                throw new \Exception(sprintf('value %d exceeds the maximum limit %d', $val, $maxi));
            }
        }

        return [$val, (string) $val];
    }

    private function validateU8(Tag $tag, int $val): array
    {
        if ($tag->isNull) {
            return [null, (string) $val];
        }

        if ($val === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception(sprintf('type uint8 not allow empty value %d', $val));
            }
            return [$val, (string) $val];
        }

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min || $mini < 0 || $mini > 255) {
                throw new \Exception('failed to parse t.Min as uint8');
            }
            if ($val < $mini) {
                throw new \Exception(sprintf('value %d is less than the minimum limit %d', $val, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max || $maxi < 0 || $maxi > 255) {
                throw new \Exception('failed to parse t.Max as uint8');
            }
            if ($val > $maxi) {
                throw new \Exception(sprintf('value %d exceeds the maximum limit %d', $val, $maxi));
            }
        }

        return [$val, (string) $val];
    }

    private function validateU16(Tag $tag, int $val): array
    {
        if ($tag->isNull) {
            return [null, (string) $val];
        }

        if ($val === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception(sprintf('type uint16 not allow empty value %d', $val));
            }
            return [$val, (string) $val];
        }

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min || $mini < 0 || $mini > 65535) {
                throw new \Exception('failed to parse t.Min as uint16');
            }
            if ($val < $mini) {
                throw new \Exception(sprintf('value %d is less than the minimum limit %d', $val, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max || $maxi < 0 || $maxi > 65535) {
                throw new \Exception('failed to parse t.Max as uint16');
            }
            if ($val > $maxi) {
                throw new \Exception(sprintf('value %d exceeds the maximum limit %d', $val, $maxi));
            }
        }

        return [$val, (string) $val];
    }

    private function validateU32(Tag $tag, int $val): array
    {
        if ($tag->isNull) {
            return [null, (string) $val];
        }

        if ($val === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception(sprintf('type uint32 not allow empty value %d', $val));
            }
            return [$val, (string) $val];
        }

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min || $mini < 0) {
                throw new \Exception('failed to parse t.Min as uint32');
            }
            if ($val < $mini) {
                throw new \Exception(sprintf('value %d is less than the minimum limit %d', $val, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max || $maxi < 0) {
                throw new \Exception('failed to parse t.Max as uint32');
            }
            if ($val > $maxi) {
                throw new \Exception(sprintf('value %d exceeds the maximum limit %d', $val, $maxi));
            }
        }

        return [$val, (string) $val];
    }

    private function validateU64(Tag $tag, int $val): array
    {
        if ($tag->isNull) {
            return [null, (string) $val];
        }

        if ($val === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception(sprintf('type uint64 not allow empty value %d', $val));
            }
            return [$val, (string) $val];
        }

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min || $mini < 0) {
                throw new \Exception('failed to parse t.Min as uint64');
            }
            if ($val < $mini) {
                throw new \Exception(sprintf('value %d is less than the minimum limit %d', $val, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max || $maxi < 0) {
                throw new \Exception('failed to parse t.Max as uint64');
            }
            if ($val > $maxi) {
                throw new \Exception(sprintf('value %d exceeds the maximum limit %d', $val, $maxi));
            }
        }

        return [$val, (string) $val];
    }

    private function validateF32(Tag $tag, float $val): array
    {
        if ($tag->isNull) {
            return [null, (string) $val];
        }

        if ($val === 0.0) {
            if (!$tag->allowEmpty) {
                throw new \Exception('type float32 not allow empty value 0.0');
            }
            return [(float) $val, (string) $val];
        }

        if ($tag->min !== '') {
            $mini = (float) $tag->min;
            if ((string) $mini !== $tag->min && $tag->min !== '0' && $tag->min !== '0.0') {
                throw new \Exception('failed to parse t.Min as float32');
            }
            if ($val < $mini) {
                throw new \Exception(sprintf('%f < min %f', $val, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (float) $tag->max;
            if ((string) $maxi !== $tag->max && $tag->max !== '0' && $tag->max !== '0.0') {
                throw new \Exception('failed to parse t.Max as float32');
            }
            if ($val > $maxi) {
                throw new \Exception(sprintf('%f > max %f', $val, $maxi));
            }
        }

        return [(float) $val, (string) $val];
    }

    private function validateF64(Tag $tag, float $val): array
    {
        if ($tag->isNull) {
            return [null, (string) $val];
        }

        if ($val === 0.0) {
            if (!$tag->allowEmpty) {
                throw new \Exception('type float64 not allow empty value 0.0');
            }
            return [$val, (string) $val];
        }

        if ($tag->min !== '') {
            $mini = (float) $tag->min;
            if ((string) $mini !== $tag->min && $tag->min !== '0' && $tag->min !== '0.0') {
                throw new \Exception('failed to parse t.Min as float64');
            }
            if ($val < $mini) {
                throw new \Exception(sprintf('%f < min %f', $val, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (float) $tag->max;
            if ((string) $maxi !== $tag->max && $tag->max !== '0' && $tag->max !== '0.0') {
                throw new \Exception('failed to parse t.Max as float64');
            }
            if ($val > $maxi) {
                throw new \Exception(sprintf('%f > max %f', $val, $maxi));
            }
        }

        return [$val, (string) $val];
    }

    private function validateBigint(Tag $tag, string $val): array
    {
        if ($tag->isNull) {
            return [null, $val];
        }

        if ($val === '0') {
            if (!$tag->allowEmpty) {
                throw new \Exception('type big.Int not allow empty value 0');
            }
            return [$val, $val];
        }

        if ($tag->min !== '') {
            if (!preg_match('/^-?\d+$/', $tag->min)) {
                throw new \Exception(sprintf('invalid min "%s" for big.Int', $tag->min));
            }
            if ($this->bigIntCmp($val, $tag->min) === -1) {
                throw new \Exception(sprintf('big.Int length %s < min %s', $val, $tag->min));
            }
        }

        if ($tag->max !== '') {
            if (!preg_match('/^-?\d+$/', $tag->max)) {
                throw new \Exception(sprintf('invalid max "%s" for big.Int', $tag->max));
            }
            if ($this->bigIntCmp($val, $tag->max) === 1) {
                throw new \Exception(sprintf('big.Int length %s > max %s', $val, $tag->max));
            }
        }

        return [$val, $val];
    }

    private function validateDatetime(Tag $tag, \DateTime $val): array
    {
        if ($tag->isNull) {
            return [null, $val->format('Y-m-d H:i:s')];
        }

        $text = $val->format('Y-m-d H:i:s');

        if ($val->getTimestamp() === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception(sprintf('datetime type does not allow empty "%s". you can set allow_empty or child_allow_empty to allow it.', $text));
            }
        }

        return [$val, $text];
    }

    private function validateDate(Tag $tag, \DateTime $val): array
    {
        if ($tag->isNull) {
            return [null, $val->format('Y-m-d')];
        }

        $text = $val->format('Y-m-d');

        if ($val->getTimestamp() === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception(sprintf('date type does not allow empty "%s". you can set allow_empty or child_allow_empty to allow it.', $text));
            }
        }

        return [$val, $text];
    }

    private function validateTime(Tag $tag, \DateTime $val): array
    {
        if ($tag->isNull) {
            return [null, $val->format('H:i:s')];
        }

        $text = $val->format('H:i:s');

        if ($val->getTimestamp() === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception(sprintf('time type does not allow empty "%s". you can set allow_empty or child_allow_empty to allow it.', $text));
            }
        }

        return [$val, $text];
    }

    private function validateUUID(Tag $tag, string $val): array
    {
        if ($tag->isNull) {
            return [null, $val];
        }

        if ($val === '') {
            if (!$tag->allowEmpty) {
                throw new \Exception('type uuid not allow empty value ""');
            }
            return [str_repeat("\x00", 16), $val];
        }

        $uuidPattern = '/^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/';
        if (!preg_match($uuidPattern, $val)) {
            throw new \Exception(sprintf("value '%s' does not match UUID pattern", $val));
        }

        $hex = str_replace('-', '', $val);
        if (strlen($hex) !== 32) {
            throw new \Exception(sprintf('invalid uuid: %s', $val));
        }

        $uuidBytes = '';
        for ($i = 0; $i < 32; $i += 2) {
            $uuidBytes .= chr(hexdec(substr($hex, $i, 2)));
        }

        return [$uuidBytes, $val];
    }

    private function validateDecimal(Tag $tag, string $val): array
    {
        if ($tag->isNull) {
            return [null, $val];
        }

        if ($val === '') {
            if (!$tag->allowEmpty) {
                throw new \Exception('type decimal not allow empty value ""');
            }
            return [$val, $val];
        }

        $decimalPattern = '/^-?\d+\.\d+$/';
        if (!preg_match($decimalPattern, $val)) {
            throw new \Exception(sprintf('invalid decimal "%s", must be like "0.0"', $val));
        }

        return [$val, $val];
    }

    private function validateIP(Tag $tag, string $val): array
    {
        if ($tag->isNull) {
            return [null, $val];
        }

        if ($val === '' || $val === '<nil>') {
            if (!$tag->allowEmpty) {
                throw new \Exception('type ip not allow empty value ""');
            }
            return [$val, $val];
        }

        if ($tag->version === 4) {
            if (!filter_var($val, FILTER_VALIDATE_IP, FILTER_FLAG_IPV4)) {
                throw new \Exception(sprintf('invalid ipv4: %s', $val));
            }
        }

        if ($tag->version === 6) {
            if (!filter_var($val, FILTER_VALIDATE_IP, FILTER_FLAG_IPV6) || filter_var($val, FILTER_VALIDATE_IP, FILTER_FLAG_IPV4)) {
                throw new \Exception(sprintf('invalid ipv6: %s', $val));
            }
        }

        return [$val, $val];
    }

    private function validateURL(Tag $tag, string $val): array
    {
        if ($tag->isNull) {
            return [null, $val];
        }

        if ($val === '') {
            if (!$tag->allowEmpty) {
                throw new \Exception('type url not allow empty value ""');
            }
            return [$val, $val];
        }

        $parsed = parse_url($val);
        if ($parsed === false) {
            throw new \Exception(sprintf('invalid url: %s', $val));
        }

        $scheme = $parsed['scheme'] ?? '';
        if ($scheme !== 'http' && $scheme !== 'https') {
            throw new \Exception(sprintf('invalid url: %s', $val));
        }

        $host = $parsed['host'] ?? '';
        if ($host === '') {
            throw new \Exception(sprintf('invalid url: %s', $val));
        }

        return [$val, $val];
    }

    private function validateEmail(Tag $tag, string $val): array
    {
        if ($tag->isNull) {
            return [null, $val];
        }

        if ($val === '') {
            if (!$tag->allowEmpty) {
                throw new \Exception('type email not allow empty value ""');
            }
            return [$val, $val];
        }

        $emailPattern = '/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/';
        if (!preg_match($emailPattern, $val)) {
            throw new \Exception(sprintf("value '%s' does not match email pattern", $val));
        }

        return [$val, $val];
    }

    private function validateEnum(Tag $tag, string $val): array
    {
        if ($tag->isNull) {
            return [null, $val];
        }

        if ($val === '') {
            if (!$tag->allowEmpty) {
                throw new \Exception('type enum not allow empty value ""');
            }
            return [$val, $val];
        }

        $enums = explode('|', $tag->enumValues);
        $idx = -1;
        foreach ($enums as $i => $s) {
            if (trim($s) === $val) {
                $idx = $i;
                break;
            }
        }

        if ($idx === -1) {
            $enumList = implode(', ', $enums);
            throw new \Exception(sprintf("value '%s' not found in enum: [%s]", $val, $enumList));
        }

        return [$idx, $val];
    }

    private function validateImage(Tag $tag, string $val): array
    {
        if ($tag->isNull) {
            return [null, $val];
        }

        $l = strlen($val);

        if ($l === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception('type image not allow empty value []byte{}');
            }
            return [$val, base64_encode($val)];
        }

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min) {
                throw new \Exception('failed to parse t.Min as int');
            }
            if ($l < $mini) {
                throw new \Exception(sprintf('[]byte length %d < min %d', $l, $mini));
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max) {
                throw new \Exception('failed to parse t.Max as int');
            }
            if ($l > $maxi) {
                throw new \Exception(sprintf('[]byte length %d > max %d', $l, $maxi));
            }
        }

        if ($tag->size !== 0) {
            if ($l !== $tag->size) {
                throw new \Exception(sprintf('[]byte length %d != size %d', $l, $tag->size));
            }
        }

        return [$val, base64_encode($val)];
    }

    private function validateArr(Tag $tag, array $value): void
    {
        if ($tag->isNull) {
            return;
        }

        $l = count($value);

        if ($l === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception('type array not allow empty');
            }
            return;
        }

        if ($tag->size > 0) {
            if ($l > $tag->size) {
                throw new \Exception('type array over size');
            }
        }

        if ($tag->childUnique && $l > 0) {
            $seen = [];
            foreach ($value as $i => $item) {
                $key = is_object($item) ? spl_object_hash($item) : (is_array($item) ? serialize($item) : $item);
                if (isset($seen[$key])) {
                    throw new \Exception(sprintf('array duplicate value found at index: %d', $i));
                }
                $seen[$key] = true;
            }
        }
    }

    private function validateVec(Tag $tag, array $value): void
    {
        if ($tag->isNull) {
            return;
        }

        $l = count($value);

        if ($l === 0) {
            if (!$tag->allowEmpty) {
                throw new \Exception('type slice not allow empty');
            }
            return;
        }

        if ($tag->childUnique && $l > 0) {
            $seen = [];
            foreach ($value as $i => $item) {
                $key = is_object($item) ? spl_object_hash($item) : (is_array($item) ? serialize($item) : $item);
                if (isset($seen[$key])) {
                    throw new \Exception(sprintf('slice duplicate value found at index: %d', $i));
                }
                $seen[$key] = true;
            }
        }
    }

    private function validateObj(Tag $tag): void
    {
        if ($tag->isNull) {
            return;
        }
    }

    private function validateMap(Tag $tag): void
    {
        if ($tag->isNull) {
            return;
        }
    }

    private function bigIntCmp(string $a, string $b): int
    {
        $negA = ($a[0] === '-');
        $negB = ($b[0] === '-');

        if ($negA && !$negB) {
            return -1;
        }
        if (!$negA && $negB) {
            return 1;
        }

        $aAbs = ltrim($a, '-');
        $bAbs = ltrim($b, '-');

        $lenA = strlen($aAbs);
        $lenB = strlen($bAbs);

        if ($lenA !== $lenB) {
            if ($negA) {
                return $lenA > $lenB ? -1 : 1;
            }
            return $lenA > $lenB ? 1 : -1;
        }

        $cmp = strcmp($aAbs, $bAbs);
        if ($negA) {
            return -$cmp;
        }
        return $cmp <=> 0;
    }
}
