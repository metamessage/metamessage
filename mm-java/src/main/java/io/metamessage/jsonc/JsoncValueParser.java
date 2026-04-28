package io.metamessage.jsonc;

import io.metamessage.mm.MmTag;
import io.metamessage.mm.MmValidator;
import io.metamessage.mm.ValueType;
import java.math.BigInteger;
import java.net.URI;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.format.DateTimeFormatter;
import java.util.Base64;
import java.util.UUID;

/**
 * Coerces JSON literals to Java values using {@link MmTag} (subset of Go {@code parser} + {@code ast.Validate*}).
 */
final class JsoncValueParser {
    private JsoncValueParser() {}

    static JcNode.JcValue stringLiteral(String text, MmTag tag, String path) {
        ValueType t = tag.type;
        if (t == ValueType.UNKNOWN) {
            tag = tag.copy();
            tag.type = ValueType.STRING;
            t = ValueType.STRING;
        }
        
        // 验证值
        validateValue(text, tag);
        
        return switch (t) {
            case STRING, EMAIL, DECIMAL, ENUM -> {
                if (tag.isNull && !text.isEmpty()) {
                    throw new JsoncException("null string must be empty");
                }
                yield new JcNode.JcValue(tag.isNull ? "" : text, text, tag, path);
            }
            case URL -> {
                if (tag.isNull && !text.isEmpty()) {
                    throw new JsoncException("null url");
                }
                yield new JcNode.JcValue(URI.create(text.isEmpty() ? "http://localhost" : text), text, tag, path);
            }
            case IP -> {
                if (tag.isNull) {
                    yield new JcNode.JcValue("", text, tag, path);
                }
                yield new JcNode.JcValue(text, text, tag, path);
            }
            case DATETIME -> {
                if (tag.isNull) {
                    yield new JcNode.JcValue(LocalDateTime.of(1970, 1, 1, 0, 0, 0), "1970-01-01 00:00:00", tag, path);
                }
                LocalDateTime d = LocalDateTime.parse(text, DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"));
                yield new JcNode.JcValue(d, text, tag, path);
            }
            case DATE -> {
                if (tag.isNull) {
                    yield new JcNode.JcValue(LocalDate.of(1970, 1, 1), "1970-01-01", tag, path);
                }
                LocalDate d = LocalDate.parse(text, DateTimeFormatter.ISO_LOCAL_DATE);
                yield new JcNode.JcValue(d, text, tag, path);
            }
            case TIME -> {
                if (tag.isNull) {
                    yield new JcNode.JcValue(LocalTime.MIDNIGHT, "00:00:00", tag, path);
                }
                LocalTime tm = LocalTime.parse(text, DateTimeFormatter.ISO_LOCAL_TIME);
                yield new JcNode.JcValue(tm, text, tag, path);
            }
            case BYTES, IMAGE -> {
                if (tag.isNull) {
                    yield new JcNode.JcValue(new byte[0], text, tag, path);
                }
                byte[] dec = Base64.getDecoder().decode(text);
                yield new JcNode.JcValue(dec, text, tag, path);
            }
            case UUID -> {
                if (tag.isNull) {
                    yield new JcNode.JcValue(new byte[16], text, tag, path);
                }
                UUID u = UUID.fromString(text);
                yield new JcNode.JcValue(u, text, tag, path);
            }
            case INT -> {
                if (tag.isNull) {
                    yield new JcNode.JcValue(0, "0", tag, path);
                }
                int v = Integer.parseInt(text);
                yield new JcNode.JcValue(v, text, tag, path);
            }
            case BOOL -> {
                if (!"true".equals(text) && !"false".equals(text)) {
                    throw new JsoncException("bool string");
                }
                yield new JcNode.JcValue(Boolean.parseBoolean(text), text, tag, path);
            }
            case BIGINT -> {
                if (tag.isNull) {
                    yield new JcNode.JcValue(BigInteger.ZERO, "0", tag, path);
                }
                BigInteger bi = new BigInteger(text);
                yield new JcNode.JcValue(bi, text, tag, path);
            }
            default -> throw new JsoncException("string literal for type " + t);
        };
    }
    
    private static void validateValue(Object value, MmTag tag) {
        // 如果是 null 值，跳过验证
        if (tag.isNull) {
            return;
        }
        
        MmValidator.ValidationResult result;
        switch (tag.type) {
            case STRING:
                result = MmValidator.validateString((String) value, tag);
                break;
            case EMAIL:
                result = MmValidator.validateEmail((String) value, tag);
                break;
            case URL:
                result = MmValidator.validateURL((String) value, tag);
                break;
            case IP:
                result = MmValidator.validateIP((String) value, tag);
                break;
            case DATETIME:
                try {
                    LocalDateTime datetime = LocalDateTime.parse((String) value, DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"));
                    result = MmValidator.validateDateTime(datetime, tag);
                } catch (Exception e) {
                    throw new JsoncException("invalid datetime format");
                }
                break;
            case DATE:
                try {
                    LocalDate date = LocalDate.parse((String) value, DateTimeFormatter.ISO_LOCAL_DATE);
                    result = MmValidator.validateDate(date, tag);
                } catch (Exception e) {
                    throw new JsoncException("invalid date format");
                }
                break;
            case TIME:
                try {
                    LocalTime time = LocalTime.parse((String) value, DateTimeFormatter.ISO_LOCAL_TIME);
                    result = MmValidator.validateTime(time, tag);
                } catch (Exception e) {
                    throw new JsoncException("invalid time format");
                }
                break;
            case BYTES:
                try {
                    byte[] bytes = Base64.getDecoder().decode((String) value);
                    result = MmValidator.validateBytes(bytes, tag);
                } catch (Exception e) {
                    throw new JsoncException("invalid base64 format");
                }
                break;
            case IMAGE:
                try {
                    byte[] bytes = Base64.getDecoder().decode((String) value);
                    result = MmValidator.validateImage(bytes, tag);
                } catch (Exception e) {
                    throw new JsoncException("invalid base64 format");
                }
                break;
            case UUID:
                result = MmValidator.validateUUID((String) value, tag);
                break;
            case DECIMAL:
                result = MmValidator.validateDecimal((String) value, tag);
                break;
            case ENUM:
                result = MmValidator.validateEnum((String) value, tag);
                break;
            default:
                return;
        }
        
        if (!result.isSuccess()) {
            throw new JsoncException(result.getError());
        }
    }

    static JcNode.JcValue numberLiteral(String text, MmTag tag, String path) {
        MmTag tg = tag;
        if (tg.type == ValueType.UNKNOWN) {
            if (text.contains(".")) {
                tg = tg.copy();
                tg.type = ValueType.FLOAT64;
            } else {
                tg = tg.copy();
                tg.type = text.startsWith("-") ? ValueType.INT : ValueType.INT;
            }
        }
        if (text.contains(".") && tg.type != ValueType.FLOAT32 && tg.type != ValueType.FLOAT64 && tg.type != ValueType.DECIMAL) {
            if (tg.type == ValueType.UNKNOWN) {
                tg = tg.copy();
                tg.type = ValueType.FLOAT64;
            }
        }
        
        // 验证值
        validateNumberValue(text, tg);
        
        if (text.contains("e") || text.contains("E")) {
            double d = Double.parseDouble(text);
            return new JcNode.JcValue(d, text, tg, path);
        }
        if (text.contains(".") && (tg.type == ValueType.FLOAT32 || tg.type == ValueType.FLOAT64 || tg.type == ValueType.DECIMAL)) {
            if (tg.type == ValueType.FLOAT32) {
                float f = Float.parseFloat(text);
                return new JcNode.JcValue(f, text, tg, path);
            }
            double d = Double.parseDouble(text);
            return new JcNode.JcValue(d, text, tg, path);
        }
        if (text.startsWith("-") || (tg.type != ValueType.UINT && tg.type != ValueType.UINT8 && tg.type != ValueType.UINT16 && tg.type != ValueType.UINT32 && tg.type != ValueType.UINT64)) {
            if (text.startsWith("-") && (tg.type == ValueType.BIGINT)) {
                BigInteger bi = new BigInteger(text);
                return new JcNode.JcValue(bi, text, tg, path);
            }
            long v = Long.parseLong(text);
            return longToJcValue(v, text, tg, path);
        }
        long u = Long.parseLong(text);
        return longToJcValue(u, text, tg, path);
    }
    
    private static void validateNumberValue(String text, MmTag tag) {
        // 如果是 null 值，跳过验证
        if (tag.isNull) {
            return;
        }
        
        MmValidator.ValidationResult result;
        if (text.contains(".") || text.contains("e") || text.contains("E")) {
            double value = Double.parseDouble(text);
            result = MmValidator.validateFloat(value, tag);
        } else if (tag.type == ValueType.BIGINT) {
            // 对于 BIGINT 类型，跳过数值范围验证
            // 因为 BigInteger 可以处理任意大小的整数
            return;
        } else {
            try {
                long value = Long.parseLong(text);
                result = MmValidator.validateInt(value, tag);
            } catch (NumberFormatException e) {
                // 如果数值超过 long 范围，跳过验证
                // 因为这可能是一个 BIGINT 类型的值
                return;
            }
        }
        
        if (!result.isSuccess()) {
            throw new JsoncException(result.getError());
        }
    }

    private static JcNode.JcValue longToJcValue(long v, String text, MmTag tag, String path) {
        return switch (tag.type) {
            case INT -> new JcNode.JcValue((int) v, text, tag, path);
            case INT8 -> new JcNode.JcValue((byte) v, text, tag, path);
            case INT16 -> new JcNode.JcValue((short) v, text, tag, path);
            case INT32 -> new JcNode.JcValue((int) v, text, tag, path);
            case INT64 -> new JcNode.JcValue(v, text, tag, path);
            case UINT, UINT8, UINT16, UINT32, UINT64 -> {
                if (v < 0) {
                    throw new JsoncException("unsigned negative");
                }
                yield new JcNode.JcValue(v, text, tag, path);
            }
            default -> new JcNode.JcValue((int) v, text, tag, path);
        };
    }

    static JcNode.JcValue boolLiteral(boolean v, MmTag tag, String path) {
        MmTag t = tag;
        if (t.type == ValueType.UNKNOWN) {
            t = t.copy();
            t.type = ValueType.BOOL;
        }
        if (t.type != ValueType.BOOL) {
            throw new JsoncException("type mismatch for bool");
        }
        
        // 验证值（如果不是 null 值）
        if (!t.isNull) {
            MmValidator.ValidationResult result = MmValidator.validateBool(v, t);
            if (!result.isSuccess()) {
                throw new JsoncException(result.getError());
            }
        }
        
        return new JcNode.JcValue(v, v ? "true" : "false", t, path);
    }
}
