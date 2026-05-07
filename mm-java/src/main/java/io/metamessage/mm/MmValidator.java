package io.github.metamessage.mm;

import java.math.BigInteger;
import java.net.InetAddress;
import java.net.URI;
import java.net.URISyntaxException;
import java.nio.charset.StandardCharsets;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.ZoneOffset;
import java.time.format.DateTimeParseException;
import java.util.Base64;
import java.util.HashSet;
import java.util.Set;
import java.util.regex.Pattern;
import java.util.regex.PatternSyntaxException;

public class MmValidator {

    private static final Pattern EMAIL_REGEX = Pattern.compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$");
    private static final Pattern DECIMAL_REGEX = Pattern.compile("^-?\\d+\\.\\d+$");
    private static final Pattern UUID_REGEX = Pattern.compile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$");

    public static ValidationResult validateArray(Object value, MmTag tag) {
        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type array not support location UTC" + tag.locationHours);
        }

        if (value == null) {
            return ValidationResult.error("array value cannot be null");
        }

        if (!(value instanceof Object[] arr)) {
            return ValidationResult.error("value is not an array");
        }

        int length = arr.length;

        if (length == 0) {
            if (tag.allowEmpty) {
                return ValidationResult.success();
            }
            return ValidationResult.error("type array not allow empty");
        }

        if (tag.size > 0) {
            if (length > tag.size) {
                return ValidationResult.error("type array over size");
            }
        }

        if (tag.childUnique) {
            Set<Object> seen = new HashSet<>();
            for (int i = 0; i < length; i++) {
                Object item = arr[i];
                if (seen.contains(item)) {
                    return ValidationResult.error("array duplicate value found: " + item + ", index: " + i);
                }
                seen.add(item);
            }
        }

        return ValidationResult.success();
    }

    public static ValidationResult validateStruct(MmTag tag) {
        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type struct not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success();
    }

    public static ValidationResult validateMap(MmTag tag) {
        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type map not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success();
    }

    public static ValidationResult validateString(String value, MmTag tag) {
        if (value == null) {
            return ValidationResult.error("string value cannot be null");
        }

        if (value.isEmpty()) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, value);
            }
            return ValidationResult.error("type string not allow empty value \"\"" + value);
        }

        if (!tag.pattern.isEmpty()) {
            try {
                Pattern pattern = Pattern.compile(tag.pattern);
                if (!pattern.matcher(value).matches()) {
                    return ValidationResult.error("value \"" + value + "\" does not match pattern " + tag.pattern);
                }
            } catch (PatternSyntaxException e) {
                return ValidationResult.error("pattern \"" + tag.pattern + "\" compile err: " + e.getMessage());
            }
        }

        int length = value.length();

        if (!tag.min.isEmpty()) {
            try {
                int mini = Integer.parseInt(tag.min);
                if (length < mini) {
                    return ValidationResult.error("string length " + length + " is less than the minimum limit " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse min as int: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                int maxi = Integer.parseInt(tag.max);
                if (length > maxi) {
                    return ValidationResult.error("string length " + length + " exceeds the maximum limit " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse max as int: " + e.getMessage());
            }
        }

        if (tag.size != 0) {
            if (length != tag.size) {
                return ValidationResult.error("string length " + length + " != size " + tag.size);
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type string not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, value);
    }

    public static ValidationResult validateBytes(byte[] value, MmTag tag) {
        if (value == null) {
            return ValidationResult.error("bytes value cannot be null");
        }

        int length = value.length;

        if (length == 0) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, "");
            }
            return ValidationResult.error("type []byte not allow empty value []byte{}");
        }

        if (!tag.min.isEmpty()) {
            try {
                int mini = Integer.parseInt(tag.min);
                if (length < mini) {
                    return ValidationResult.error("[]byte length " + length + " is less than the minimum limit " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse min as int: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                int maxi = Integer.parseInt(tag.max);
                if (length > maxi) {
                    return ValidationResult.error("[]byte length " + length + " exceeds the maximum limit " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse max as int: " + e.getMessage());
            }
        }

        if (tag.size != 0) {
            if (length != tag.size) {
                return ValidationResult.error("[]byte length " + length + " != size " + tag.size);
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type []byte not support location UTC" + tag.locationHours);
        }

        String base64 = Base64.getEncoder().encodeToString(value);
        return ValidationResult.success(value, base64);
    }

    public static ValidationResult validateBool(Boolean value, MmTag tag) {
        if (value == null) {
            return ValidationResult.error("bool value cannot be null");
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.allowEmpty) {
            return ValidationResult.error("type bool not support allow empty");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type bool not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, value.toString());
    }

    public static ValidationResult validateInt8(byte value, MmTag tag) {
        if (value == 0) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, "0");
            }
            return ValidationResult.error("type int8 not allow empty value " + value);
        }

        long val = (long) value;

        if (!tag.min.isEmpty()) {
            try {
                long mini = Long.parseLong(tag.min);
                if (mini < -128 || mini > 127) {
                    return ValidationResult.error("tag.min " + mini + " is out of int8 range [-128, 127]");
                }
                if (val < mini) {
                    return ValidationResult.error("value " + value + " is less than the minimum limit " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.min as int8: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                long maxi = Long.parseLong(tag.max);
                if (maxi < -128 || maxi > 127) {
                    return ValidationResult.error("tag.max " + maxi + " is out of int8 range [-128, 127]");
                }
                if (val > maxi) {
                    return ValidationResult.error("value " + value + " exceeds the maximum limit " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.max as int8: " + e.getMessage());
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type int8 not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, Byte.toString(value));
    }

    public static ValidationResult validateInt16(short value, MmTag tag) {
        if (value == 0) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, "0");
            }
            return ValidationResult.error("type int16 not allow empty value " + value);
        }

        long val = (long) value;

        if (!tag.min.isEmpty()) {
            try {
                long mini = Long.parseLong(tag.min);
                if (mini < -32768 || mini > 32767) {
                    return ValidationResult.error("tag.min " + mini + " is out of int16 range [-32768, 32767]");
                }
                if (val < mini) {
                    return ValidationResult.error("value " + value + " is less than the minimum limit " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.min as int16: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                long maxi = Long.parseLong(tag.max);
                if (maxi < -32768 || maxi > 32767) {
                    return ValidationResult.error("tag.max " + maxi + " is out of int16 range [-32768, 32767]");
                }
                if (val > maxi) {
                    return ValidationResult.error("value " + value + " exceeds the maximum limit " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.max as int16: " + e.getMessage());
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type int16 not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, Short.toString(value));
    }

    public static ValidationResult validateInt32(int value, MmTag tag) {
        if (value == 0) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, "0");
            }
            return ValidationResult.error("type int32 not allow empty value " + value);
        }

        long val = (long) value;

        if (!tag.min.isEmpty()) {
            try {
                long mini = Long.parseLong(tag.min);
                if (mini < -2147483648L || mini > 2147483647L) {
                    return ValidationResult.error("tag.min " + mini + " is out of int32 range [-2147483648, 2147483647]");
                }
                if (val < mini) {
                    return ValidationResult.error("value " + value + " is less than the minimum limit " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.min as int32: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                long maxi = Long.parseLong(tag.max);
                if (maxi < -2147483648L || maxi > 2147483647L) {
                    return ValidationResult.error("tag.max " + maxi + " is out of int32 range [-2147483648, 2147483647]");
                }
                if (val > maxi) {
                    return ValidationResult.error("value " + value + " exceeds the maximum limit " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.max as int32: " + e.getMessage());
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type int32 not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, Integer.toString(value));
    }

    public static ValidationResult validateInt64(long value, MmTag tag) {
        if (value == 0) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, "0");
            }
            return ValidationResult.error("type int64 not allow empty value " + value);
        }

        if (!tag.min.isEmpty()) {
            try {
                long mini = Long.parseLong(tag.min);
                if (mini < -9223372036854775808L || mini > 9223372036854775807L) {
                    return ValidationResult.error("tag.min " + mini + " is out of int64 range [-9223372036854775808, 9223372036854775807]");
                }
                if (value < mini) {
                    return ValidationResult.error("value " + value + " is less than the minimum limit " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.min as int64: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                long maxi = Long.parseLong(tag.max);
                if (maxi < -9223372036854775808L || maxi > 9223372036854775807L) {
                    return ValidationResult.error("tag.max " + maxi + " is out of int64 range [-9223372036854775808, 9223372036854775807]");
                }
                if (value > maxi) {
                    return ValidationResult.error("value " + value + " exceeds the maximum limit " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.max as int64: " + e.getMessage());
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type int64 not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, Long.toString(value));
    }

    public static ValidationResult validateUint(int value, MmTag tag) {
        if (value == 0) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, "0");
            }
            return ValidationResult.error("type uint not allow empty value " + value);
        }

        if (!tag.min.isEmpty()) {
            try {
                long mini = Long.parseLong(tag.min);
                if (mini < 0 || mini > 4294967295L) {
                    return ValidationResult.error("tag.min " + mini + " is out of uint range [0, 4294967295]");
                }
                if (value < mini) {
                    return ValidationResult.error("value " + value + " is less than the minimum limit " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.min as uint: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                long maxi = Long.parseLong(tag.max);
                if (maxi < 0 || maxi > 4294967295L) {
                    return ValidationResult.error("tag.max " + maxi + " is out of uint range [0, 4294967295]");
                }
                if (value > maxi) {
                    return ValidationResult.error("value " + value + " exceeds the maximum limit " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.max as uint: " + e.getMessage());
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type uint not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, Integer.toString(value));
    }

    public static ValidationResult validateUint8(short value, MmTag tag) {
        if (value == 0) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, "0");
            }
            return ValidationResult.error("type uint8 not allow empty value " + value);
        }

        if (!tag.min.isEmpty()) {
            try {
                long mini = Long.parseLong(tag.min);
                if (mini < 0 || mini > 255) {
                    return ValidationResult.error("tag.min " + mini + " is out of uint8 range [0, 255]");
                }
                if (value < mini) {
                    return ValidationResult.error("value " + value + " is less than the minimum limit " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.min as uint8: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                long maxi = Long.parseLong(tag.max);
                if (maxi < 0 || maxi > 255) {
                    return ValidationResult.error("tag.max " + maxi + " is out of uint8 range [0, 255]");
                }
                if (value > maxi) {
                    return ValidationResult.error("value " + value + " exceeds the maximum limit " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.max as uint8: " + e.getMessage());
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type uint8 not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, Short.toString(value));
    }

    public static ValidationResult validateUint16(int value, MmTag tag) {
        if (value == 0) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, "0");
            }
            return ValidationResult.error("type uint16 not allow empty value " + value);
        }

        if (!tag.min.isEmpty()) {
            try {
                long mini = Long.parseLong(tag.min);
                if (mini < 0 || mini > 65535) {
                    return ValidationResult.error("tag.min " + mini + " is out of uint16 range [0, 65535]");
                }
                if (value < mini) {
                    return ValidationResult.error("value " + value + " is less than the minimum limit " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.min as uint16: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                long maxi = Long.parseLong(tag.max);
                if (maxi < 0 || maxi > 65535) {
                    return ValidationResult.error("tag.max " + maxi + " is out of uint16 range [0, 65535]");
                }
                if (value > maxi) {
                    return ValidationResult.error("value " + value + " exceeds the maximum limit " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.max as uint16: " + e.getMessage());
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type uint16 not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, Integer.toString(value));
    }

    public static ValidationResult validateUint32(long value, MmTag tag) {
        if (value == 0) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, "0");
            }
            return ValidationResult.error("type uint32 not allow empty value " + value);
        }

        if (!tag.min.isEmpty()) {
            try {
                long mini = Long.parseLong(tag.min);
                if (mini < 0 || mini > 4294967295L) {
                    return ValidationResult.error("tag.min " + mini + " is out of uint32 range [0, 4294967295]");
                }
                if (value < mini) {
                    return ValidationResult.error("value " + value + " is less than the minimum limit " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.min as uint32: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                long maxi = Long.parseLong(tag.max);
                if (maxi < 0 || maxi > 4294967295L) {
                    return ValidationResult.error("tag.max " + maxi + " is out of uint32 range [0, 4294967295]");
                }
                if (value > maxi) {
                    return ValidationResult.error("value " + value + " exceeds the maximum limit " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.max as uint32: " + e.getMessage());
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type uint32 not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, Long.toString(value));
    }

    public static ValidationResult validateUint64(BigInteger value, MmTag tag) {
        if (value.signum() == 0) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, "0");
            }
            return ValidationResult.error("type uint64 not allow empty value 0");
        }

        if (!tag.min.isEmpty()) {
            try {
                BigInteger mini = new BigInteger(tag.min);
                if (mini.signum() < 0 || mini.compareTo(new BigInteger("18446744073709551615")) > 0) {
                    return ValidationResult.error("tag.min " + mini + " is out of uint64 range [0, 18446744073709551615]");
                }
                if (value.compareTo(mini) < 0) {
                    return ValidationResult.error("value " + value + " is less than the minimum limit " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.min as uint64: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                BigInteger maxi = new BigInteger(tag.max);
                if (maxi.signum() < 0 || maxi.compareTo(new BigInteger("18446744073709551615")) > 0) {
                    return ValidationResult.error("tag.max " + maxi + " is out of uint64 range [0, 18446744073709551615]");
                }
                if (value.compareTo(maxi) > 0) {
                    return ValidationResult.error("value " + value + " exceeds the maximum limit " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.max as uint64: " + e.getMessage());
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type uint64 not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, value.toString());
    }

    public static ValidationResult validateFloat32(float value, MmTag tag) {
        if (value == 0.0f) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, "0.0");
            }
            return ValidationResult.error("type float32 not allow empty value 0.0");
        }

        double val = (double) value;

        if (!tag.min.isEmpty()) {
            try {
                double mini = Double.parseDouble(tag.min);
                if (mini < -3.402823466e+38 || mini > 3.402823466e+38) {
                    return ValidationResult.error("tag.min " + mini + " is out of float32 range [-3.402823466e+38, 3.402823466e+38]");
                }
                if (val < mini) {
                    return ValidationResult.error("value " + value + " is less than the minimum limit " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.min as float32: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                double maxi = Double.parseDouble(tag.max);
                if (maxi < -3.402823466e+38 || maxi > 3.402823466e+38) {
                    return ValidationResult.error("tag.max " + maxi + " is out of float32 range [-3.402823466e+38, 3.402823466e+38]");
                }
                if (val > maxi) {
                    return ValidationResult.error("value " + value + " exceeds the maximum limit " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.max as float32: " + e.getMessage());
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type float32 not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, Float.toString(value));
    }

    public static ValidationResult validateFloat64(double value, MmTag tag) {
        if (value == 0.0) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, "0.0");
            }
            return ValidationResult.error("type float64 not allow empty value 0.0");
        }

        if (!tag.min.isEmpty()) {
            try {
                double mini = Double.parseDouble(tag.min);
                if (mini < -1.7976931348623157e+308 || mini > 1.7976931348623157e+308) {
                    return ValidationResult.error("tag.min " + mini + " is out of float64 range [-1.7976931348623157e+308, 1.7976931348623157e+308]");
                }
                if (value < mini) {
                    return ValidationResult.error("value " + value + " is less than the minimum limit " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.min as float64: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                double maxi = Double.parseDouble(tag.max);
                if (maxi < -1.7976931348623157e+308 || maxi > 1.7976931348623157e+308) {
                    return ValidationResult.error("tag.max " + maxi + " is out of float64 range [-1.7976931348623157e+308, 1.7976931348623157e+308]");
                }
                if (value > maxi) {
                    return ValidationResult.error("value " + value + " exceeds the maximum limit " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse tag.max as float64: " + e.getMessage());
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type float64 not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, Double.toString(value));
    }

    public static ValidationResult validateBigInteger(BigInteger value, MmTag tag) {
        if (value.signum() == 0) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, "0");
            }
            return ValidationResult.error("type big.Int not allow empty value 0");
        }

        if (!tag.min.isEmpty()) {
            try {
                BigInteger mini = new BigInteger(tag.min);
                if (value.compareTo(mini) < 0) {
                    return ValidationResult.error("big.Int " + value + " < min " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse min as big.Int: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                BigInteger maxi = new BigInteger(tag.max);
                if (value.compareTo(maxi) > 0) {
                    return ValidationResult.error("big.Int " + value + " > max " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse max as big.Int: " + e.getMessage());
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type big.Int not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, value.toString());
    }

    public static ValidationResult validateDateTime(LocalDateTime value, MmTag tag) {
        if (value == null) {
            return ValidationResult.error("datetime value cannot be null");
        }

        ZoneOffset offset = ZoneOffset.UTC;
        if (tag.locationHours != 0) {
            offset = ZoneOffset.ofHours(tag.locationHours);
        }

        LocalDateTime truncated = value.truncatedTo(java.time.temporal.ChronoUnit.SECONDS);
        String format = truncated.format(java.time.format.DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"));

        if (truncated.toEpochSecond(offset) == 0) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, format);
            }
            return ValidationResult.error("type datetime not allow empty " + format);
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        return ValidationResult.success(value, format);
    }

    public static ValidationResult validateDate(LocalDate value, MmTag tag) {
        if (value == null) {
            return ValidationResult.error("date value cannot be null");
        }

        String format = value.format(java.time.format.DateTimeFormatter.ofPattern("yyyy-MM-dd"));

        if (value.equals(LocalDate.ofEpochDay(0))) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, format);
            }
            return ValidationResult.error("type date not allow empty " + format);
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        return ValidationResult.success(value, format);
    }

    public static ValidationResult validateTime(LocalTime value, MmTag tag) {
        if (value == null) {
            return ValidationResult.error("time value cannot be null");
        }

        String format = value.format(java.time.format.DateTimeFormatter.ofPattern("HH:mm:ss"));

        if (value.equals(LocalTime.MIDNIGHT)) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, format);
            }
            return ValidationResult.error("type time not allow empty " + format);
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        return ValidationResult.success(value, format);
    }

    public static ValidationResult validateUUID(String value, MmTag tag) {
        if (value == null) {
            return ValidationResult.error("uuid value cannot be null");
        }

        if (value.isEmpty()) {
            if (tag.allowEmpty) {
                return ValidationResult.success(new byte[16], "");
            }
            return ValidationResult.error("type uuid not allow empty value \"\"");
        }

        if (!UUID_REGEX.matcher(value).matches()) {
            return ValidationResult.error("value '" + value + "' does not match UUID pattern");
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type uuid not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, value);
    }

    public static ValidationResult validateDecimal(String value, MmTag tag) {
        if (value == null) {
            return ValidationResult.error("decimal value cannot be null");
        }

        if (value.isEmpty()) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, value);
            }
            return ValidationResult.error("type decimal not allow empty value \"\"");
        }

        if (!DECIMAL_REGEX.matcher(value).matches()) {
            return ValidationResult.error("invalid decimal \"" + value + "\", must be like \"0.0\"");
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type decimal not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, value);
    }

    public static ValidationResult validateIP(String value, MmTag tag) {
        if (value == null) {
            return ValidationResult.error("ip value cannot be null");
        }

        if (value.isEmpty()) {
            if (tag.allowEmpty) {
                return ValidationResult.success(null, "");
            }
            return ValidationResult.error("type ip not allow empty value \"\"");
        }

        try {
            InetAddress addr = InetAddress.getByName(value);
            if (tag.version == 4) {
                if (addr.getAddress().length != 4) {
                    return ValidationResult.error("invalid ipv4: " + value);
                }
            } else if (tag.version == 6) {
                if (addr.getAddress().length != 16) {
                    return ValidationResult.error("invalid ipv6: " + value);
                }
            }
        } catch (Exception e) {
            return ValidationResult.error("invalid ip: " + value);
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type ip not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, value);
    }

    public static ValidationResult validateURL(String value, MmTag tag) {
        if (value == null) {
            return ValidationResult.error("url value cannot be null");
        }

        if (value.isEmpty()) {
            if (tag.allowEmpty) {
                return ValidationResult.success(null, "");
            }
            return ValidationResult.error("type url not allow empty value \"\"");
        }

        try {
            URI uri = new URI(value);
            if (!uri.getScheme().equals("http") && !uri.getScheme().equals("https")) {
                return ValidationResult.error("invalid url: " + value);
            }
            if (uri.getHost() == null || uri.getHost().isEmpty()) {
                return ValidationResult.error("invalid url: " + value);
            }
        } catch (URISyntaxException e) {
            return ValidationResult.error("invalid url: " + value);
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type url not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, value);
    }

    public static ValidationResult validateEmail(String value, MmTag tag) {
        if (value == null) {
            return ValidationResult.error("email value cannot be null");
        }

        if (value.isEmpty()) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, value);
            }
            return ValidationResult.error("type email not allow empty value \"\"");
        }

        if (!EMAIL_REGEX.matcher(value).matches()) {
            return ValidationResult.error("value '" + value + "' does not match email pattern");
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type email not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(value, value);
    }

    public static ValidationResult validateEnum(String value, MmTag tag) {
        if (value == null) {
            return ValidationResult.error("enum value cannot be null");
        }

        if (value.isEmpty()) {
            if (tag.allowEmpty) {
                return ValidationResult.success(-1, "");
            }
            return ValidationResult.error("type enum not allow empty value \"\"");
        }

        String[] enums = tag.enumValues.split("\\|");
        int idx = -1;
        for (int i = 0; i < enums.length; i++) {
            if (enums[i].trim().equals(value)) {
                idx = i;
                break;
            }
        }

        if (idx == -1) {
            return ValidationResult.error("value '" + value + "' not found in enum: " + tag.enumValues);
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type enum not support location UTC" + tag.locationHours);
        }

        return ValidationResult.success(idx, value);
    }

    public static ValidationResult validateImage(byte[] value, MmTag tag) {
        if (value == null) {
            return ValidationResult.error("image value cannot be null");
        }

        int length = value.length;

        if (length == 0) {
            if (tag.allowEmpty) {
                return ValidationResult.success(value, "");
            }
            return ValidationResult.error("type image not allow empty value []byte{}");
        }

        if (!tag.min.isEmpty()) {
            try {
                int mini = Integer.parseInt(tag.min);
                if (length < mini) {
                    return ValidationResult.error("[]byte length " + length + " < min " + mini);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse min as int: " + e.getMessage());
            }
        }

        if (!tag.max.isEmpty()) {
            try {
                int maxi = Integer.parseInt(tag.max);
                if (length > maxi) {
                    return ValidationResult.error("[]byte length " + length + " > max " + maxi);
                }
            } catch (NumberFormatException e) {
                return ValidationResult.error("failed to parse max as int: " + e.getMessage());
            }
        }

        if (tag.size != 0) {
            if (length != tag.size) {
                return ValidationResult.error("[]byte length " + length + " != size " + tag.size);
            }
        }

        if (tag.desc.length() > 65535) {
            return ValidationResult.error("desc length exceeds 65535 bytes");
        }

        if (tag.locationHours != 0) {
            return ValidationResult.error("type image not support location UTC" + tag.locationHours);
        }

        String base64 = Base64.getEncoder().encodeToString(value);
        return ValidationResult.success(value, base64);
    }

    public static class ValidationResult {
        private final boolean success;
        private final String error;
        private final Object data;
        private final String text;

        private ValidationResult(boolean success, String error, Object data, String text) {
            this.success = success;
            this.error = error;
            this.data = data;
            this.text = text;
        }

        public static ValidationResult success() {
            return new ValidationResult(true, null, null, null);
        }

        public static ValidationResult success(Object data, String text) {
            return new ValidationResult(true, null, data, text);
        }

        public static ValidationResult error(String error) {
            return new ValidationResult(false, error, null, null);
        }

        public boolean isSuccess() {
            return success;
        }

        public String getError() {
            return error;
        }

        public Object getData() {
            return data;
        }

        public String getText() {
            return text;
        }
    }
}
