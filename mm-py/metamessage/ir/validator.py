import re
import base64
from typing import Any, Optional, List

from .tag import Tag, ValueType


class ValidationResult:
    def __init__(self, valid: bool, error: Optional[str] = None, data: Any = None, text: Optional[str] = None):
        self.valid = valid
        self.error = error
        self.data = data
        self.text = text


class MmValidator:
    email_regex = re.compile(r'^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$')
    decimal_regex = re.compile(r'^-?\d+\.\d+$')
    uuid_regex = re.compile(r'^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$')

    @staticmethod
    def validate_arr(value: List[Any], tag: Tag) -> ValidationResult:
        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type array not support location UTC%d" % location_offset)

        length = len(value)

        if length == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type array not allow empty")
            return ValidationResult(True, data=value, text=str(value))

        if tag.size > 0 and length > tag.size:
            return ValidationResult(False, "type array over size")

        if tag.child_unique:
            seen = {}
            for i, item in enumerate(value):
                key = item if not isinstance(item, (dict, list)) else str(item)
                if key in seen:
                    return ValidationResult(False, "array duplicate value found: %s, index: %d" % (item, i))
                seen[key] = True

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_obj(tag: Tag) -> ValidationResult:
        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type struct not support location UTC%d" % location_offset)

        return ValidationResult(True)

    @staticmethod
    def validate_map(tag: Tag) -> ValidationResult:
        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type map not support location UTC%d" % location_offset)

        return ValidationResult(True)

    @staticmethod
    def validate_vec(value: List[Any], tag: Tag) -> ValidationResult:
        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type vec not support location UTC%d" % location_offset)

        length = len(value)

        if length == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type vec not allow empty")
            return ValidationResult(True)

        if tag.child_unique:
            seen = {}
            for i, item in enumerate(value):
                key = item if not isinstance(item, (dict, list)) else str(item)
                if key in seen:
                    return ValidationResult(False, "vec duplicate value found: %s, index: %d" % (item, i))
                seen[key] = True

        return ValidationResult(True)

    @staticmethod
    def validate_str(value: str, tag: Tag) -> ValidationResult:
        if value == "":
            if not tag.allow_empty:
                return ValidationResult(False, "type string not allow empty value %s" % repr(value))
            return ValidationResult(True, data=value, text=value)

        if tag.pattern:
            try:
                regex = re.compile(tag.pattern)
                if not regex.match(value):
                    return ValidationResult(False, "value %s does not match pattern %s" % (repr(value), tag.pattern))
            except re.error as e:
                return ValidationResult(False, "pattern %s compile err: %s" % (repr(tag.pattern), e))

        length = len(value)

        if tag.min:
            try:
                mini = int(tag.min)
                if length < mini:
                    return ValidationResult(False, "string length %d is less than the minimum limit %d" % (length, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as int: %s" % tag.min)

        if tag.max:
            try:
                maxi = int(tag.max)
                if length > maxi:
                    return ValidationResult(False, "string length %d exceeds the maximum limit %d" % (length, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as int: %s" % tag.max)

        if tag.size > 0 and length != tag.size:
            return ValidationResult(False, "string length %d != size %d" % (length, tag.size))

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type string not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=value)

    @staticmethod
    def validate_bytes(value: bytes, tag: Tag) -> ValidationResult:
        length = len(value)

        if length == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type []byte not allow empty value []byte{}")
            return ValidationResult(True, data=value, text="")

        if tag.min:
            try:
                mini = int(tag.min)
                if length < mini:
                    return ValidationResult(False, "[]byte length %d is less than the minimum limit %d" % (length, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as int: %s" % tag.min)

        if tag.max:
            try:
                maxi = int(tag.max)
                if length > maxi:
                    return ValidationResult(False, "[]byte length %d exceeds the maximum limit %d" % (length, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as int: %s" % tag.max)

        if tag.size > 0 and length != tag.size:
            return ValidationResult(False, "[]byte length %d != size %d" % (length, tag.size))

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type []byte not support location UTC%d" % location_offset)

        text = base64.b64encode(value).decode('utf-8')

        return ValidationResult(True, data=value, text=text)

    @staticmethod
    def validate_bool(value: bool, tag: Tag) -> ValidationResult:
        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.allow_empty:
            return ValidationResult(False, "type bool not support allow empty")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type bool not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=str(value).lower())

    @staticmethod
    def validate_i(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type int not allow empty value %d" % value)
            return ValidationResult(True, data=value, text="0")

        if tag.min:
            try:
                mini = int(tag.min)
                if value < mini:
                    return ValidationResult(False, "value %d is less than the minimum limit %d" % (value, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as int: %s" % tag.min)

        if tag.max:
            try:
                maxi = int(tag.max)
                if value > maxi:
                    return ValidationResult(False, "value %d exceeds the maximum limit %d" % (value, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as int: %s" % tag.max)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type int not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_i8(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type int8 not allow empty value %d" % value)
            return ValidationResult(True, data=value, text="0")

        if tag.min:
            try:
                mini = int(tag.min)
                if value < mini:
                    return ValidationResult(False, "value %d is less than the minimum limit %d" % (value, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as int8: %s" % tag.min)

        if tag.max:
            try:
                maxi = int(tag.max)
                if value > maxi:
                    return ValidationResult(False, "value %d exceeds the maximum limit %d" % (value, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as int8: %s" % tag.max)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type int8 not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_i16(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type int16 not allow empty value %d" % value)
            return ValidationResult(True, data=value, text="0")

        if tag.min:
            try:
                mini = int(tag.min)
                if value < mini:
                    return ValidationResult(False, "value %d is less than the minimum limit %d" % (value, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as int16: %s" % tag.min)

        if tag.max:
            try:
                maxi = int(tag.max)
                if value > maxi:
                    return ValidationResult(False, "value %d exceeds the maximum limit %d" % (value, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as int16: %s" % tag.max)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type int16 not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_i32(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type int32 not allow empty value %d" % value)
            return ValidationResult(True, data=value, text="0")

        if tag.min:
            try:
                mini = int(tag.min)
                if value < mini:
                    return ValidationResult(False, "value %d is less than the minimum limit %d" % (value, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as int32: %s" % tag.min)

        if tag.max:
            try:
                maxi = int(tag.max)
                if value > maxi:
                    return ValidationResult(False, "value %d exceeds the maximum limit %d" % (value, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as int32: %s" % tag.max)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type int32 not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_i64(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type int64 not allow empty value %d" % value)
            return ValidationResult(True, data=value, text="0")

        if tag.min:
            try:
                mini = int(tag.min)
                if value < mini:
                    return ValidationResult(False, "value %d is less than the minimum limit %d" % (value, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as int64: %s" % tag.min)

        if tag.max:
            try:
                maxi = int(tag.max)
                if value > maxi:
                    return ValidationResult(False, "value %d exceeds the maximum limit %d" % (value, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as int64: %s" % tag.max)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type int64 not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_u(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type uint not allow empty value %d" % value)
            return ValidationResult(True, data=value, text="0")

        if tag.min:
            try:
                mini = int(tag.min)
                if mini < 0:
                    return ValidationResult(False, "failed to parse tag.min as uint: %s" % tag.min)
                if value < mini:
                    return ValidationResult(False, "value %d is less than the minimum limit %d" % (value, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as uint: %s" % tag.min)

        if tag.max:
            try:
                maxi = int(tag.max)
                if maxi < 0:
                    return ValidationResult(False, "failed to parse tag.max as uint: %s" % tag.max)
                if value > maxi:
                    return ValidationResult(False, "value %d exceeds the maximum limit %d" % (value, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as uint: %s" % tag.max)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type uint not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_u8(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type uint8 not allow empty value %d" % value)
            return ValidationResult(True, data=value, text="0")

        if tag.min:
            try:
                mini = int(tag.min)
                if mini < 0:
                    return ValidationResult(False, "failed to parse tag.min as uint8: %s" % tag.min)
                if value < mini:
                    return ValidationResult(False, "value %d is less than the minimum limit %d" % (value, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as uint8: %s" % tag.min)

        if tag.max:
            try:
                maxi = int(tag.max)
                if maxi < 0:
                    return ValidationResult(False, "failed to parse tag.max as uint8: %s" % tag.max)
                if value > maxi:
                    return ValidationResult(False, "value %d exceeds the maximum limit %d" % (value, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as uint8: %s" % tag.max)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type uint8 not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_u16(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type uint16 not allow empty value %d" % value)
            return ValidationResult(True, data=value, text="0")

        if tag.min:
            try:
                mini = int(tag.min)
                if mini < 0:
                    return ValidationResult(False, "failed to parse tag.min as uint16: %s" % tag.min)
                if value < mini:
                    return ValidationResult(False, "value %d is less than the minimum limit %d" % (value, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as uint16: %s" % tag.min)

        if tag.max:
            try:
                maxi = int(tag.max)
                if maxi < 0:
                    return ValidationResult(False, "failed to parse tag.max as uint16: %s" % tag.max)
                if value > maxi:
                    return ValidationResult(False, "value %d exceeds the maximum limit %d" % (value, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as uint16: %s" % tag.max)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type uint16 not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_u32(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type uint32 not allow empty value %d" % value)
            return ValidationResult(True, data=value, text="0")

        if tag.min:
            try:
                mini = int(tag.min)
                if mini < 0:
                    return ValidationResult(False, "failed to parse tag.min as uint32: %s" % tag.min)
                if value < mini:
                    return ValidationResult(False, "value %d is less than the minimum limit %d" % (value, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as uint32: %s" % tag.min)

        if tag.max:
            try:
                maxi = int(tag.max)
                if maxi < 0:
                    return ValidationResult(False, "failed to parse tag.max as uint32: %s" % tag.max)
                if value > maxi:
                    return ValidationResult(False, "value %d exceeds the maximum limit %d" % (value, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as uint32: %s" % tag.max)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type uint32 not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_u64(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type uint64 not allow empty value %d" % value)
            return ValidationResult(True, data=value, text="0")

        if tag.min:
            try:
                mini = int(tag.min)
                if mini < 0:
                    return ValidationResult(False, "failed to parse tag.min as uint64: %s" % tag.min)
                if value < mini:
                    return ValidationResult(False, "value %d is less than the minimum limit %d" % (value, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as uint64: %s" % tag.min)

        if tag.max:
            try:
                maxi = int(tag.max)
                if maxi < 0:
                    return ValidationResult(False, "failed to parse tag.max as uint64: %s" % tag.max)
                if value > maxi:
                    return ValidationResult(False, "value %d exceeds the maximum limit %d" % (value, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as uint64: %s" % tag.max)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type uint64 not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_f32(value: float, tag: Tag) -> ValidationResult:
        if value == 0.0:
            if not tag.allow_empty:
                return ValidationResult(False, "type float32 not allow empty value 0.0")
            return ValidationResult(True, data=value, text="0.0")

        if tag.min:
            try:
                mini = float(tag.min)
                if value < mini:
                    return ValidationResult(False, "%f < min %f" % (value, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as float32: %s" % tag.min)

        if tag.max:
            try:
                maxi = float(tag.max)
                if value > maxi:
                    return ValidationResult(False, "%f > max %f" % (value, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as float32: %s" % tag.max)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type float32 not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_f64(value: float, tag: Tag) -> ValidationResult:
        if value == 0.0:
            if not tag.allow_empty:
                return ValidationResult(False, "type float64 not allow empty value 0.0")
            return ValidationResult(True, data=value, text="0.0")

        if tag.min:
            try:
                mini = float(tag.min)
                if value < mini:
                    return ValidationResult(False, "%f < min %f" % (value, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as float64: %s" % tag.min)

        if tag.max:
            try:
                maxi = float(tag.max)
                if value > maxi:
                    return ValidationResult(False, "%f > max %f" % (value, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as float64: %s" % tag.max)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type float64 not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_bigint(value: str, tag: Tag) -> ValidationResult:
        try:
            from decimal import Decimal
            val = Decimal(value)
        except Exception:
            return ValidationResult(False, "invalid big.Int value: %s" % value)

        if val == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type big.Int not allow empty value 0")
            return ValidationResult(True, data=value, text="0")

        if tag.min:
            try:
                from decimal import Decimal
                mini = Decimal(tag.min)
                if val < mini:
                    return ValidationResult(False, "big.Int %s < min %s" % (value, tag.min))
            except Exception:
                return ValidationResult(False, "invalid min %s for big.Int" % repr(tag.min))

        if tag.max:
            try:
                from decimal import Decimal
                maxi = Decimal(tag.max)
                if val > maxi:
                    return ValidationResult(False, "big.Int %s > max %s" % (value, tag.max))
            except Exception:
                return ValidationResult(False, "invalid max %s for big.Int" % repr(tag.max))

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type big.Int not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=str(val))

    @staticmethod
    def validate_datetime(value: Any, tag: Tag) -> ValidationResult:
        if not hasattr(value, 'timestamp'):
            return ValidationResult(False, "expected datetime object, got %s" % type(value).__name__)

        if value.timestamp() == 0:
            if not tag.allow_empty:
                format_str = value.strftime("%Y-%m-%d %H:%M:%S")
                return ValidationResult(False,
                    "datetime type does not allow empty \"%s\". you can set allow_empty or child_allow_empty to allow it." % format_str)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        format_str = value.strftime("%Y-%m-%d %H:%M:%S")

        return ValidationResult(True, data=value, text=format_str)

    @staticmethod
    def validate_date(value: Any, tag: Tag) -> ValidationResult:
        if not hasattr(value, 'timestamp'):
            return ValidationResult(False, "expected date object, got %s" % type(value).__name__)

        if value.timestamp() == 0:
            if not tag.allow_empty:
                format_str = value.strftime("%Y-%m-%d")
                return ValidationResult(False,
                    "date type does not allow empty \"%s\". you can set allow_empty or child_allow_empty to allow it." % format_str)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        format_str = value.strftime("%Y-%m-%d")

        return ValidationResult(True, data=value, text=format_str)

    @staticmethod
    def validate_time(value: Any, tag: Tag) -> ValidationResult:
        if not hasattr(value, 'timestamp'):
            return ValidationResult(False, "expected time object, got %s" % type(value).__name__)

        if value.timestamp() == 0:
            if not tag.allow_empty:
                format_str = value.strftime("%H:%M:%S")
                return ValidationResult(False,
                    "time type does not allow empty \"%s\". you can set allow_empty or child_allow_empty to allow it." % format_str)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        format_str = value.strftime("%H:%M:%S")

        return ValidationResult(True, data=value, text=format_str)

    @staticmethod
    def validate_uuid(value: str, tag: Tag) -> ValidationResult:
        if value == "":
            if not tag.allow_empty:
                return ValidationResult(False, "type uuid not allow empty value \"\"")
            return ValidationResult(True, data=b'\x00' * 16, text=value)

        if not MmValidator.uuid_regex.match(value):
            return ValidationResult(False, "value '%s' does not match UUID pattern" % value)

        try:
            uuid_bytes = bytes.fromhex(value.replace('-', ''))
        except ValueError:
            return ValidationResult(False, "invalid uuid: %s" % value)

        if tag.version > 0:
            version = int(value[14], 16)
            if tag.version != version:
                return ValidationResult(False, "invalid uuid version")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type uuid not support location UTC%d" % location_offset)

        return ValidationResult(True, data=uuid_bytes, text=value)

    @staticmethod
    def validate_decimal(value: str, tag: Tag) -> ValidationResult:
        if value == "":
            if not tag.allow_empty:
                return ValidationResult(False, "type decimal not allow empty value \"\"")
            return ValidationResult(True, data=value, text=value)

        if not MmValidator.decimal_regex.match(value):
            return ValidationResult(False, "invalid decimal \"%s\", must be like \"0.0\"" % value)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type decimal not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=value)

    @staticmethod
    def validate_ip(value: str, tag: Tag) -> ValidationResult:
        if value == "":
            if not tag.allow_empty:
                return ValidationResult(False, "type ip not allow empty value \"\"")
            return ValidationResult(True, data=value, text="")

        import ipaddress
        try:
            ip = ipaddress.ip_address(value)
        except ValueError:
            return ValidationResult(False, "invalid ip: %s" % value)

        if tag.version == 4:
            if ip.version != 4:
                return ValidationResult(False, "invalid ipv4: %s" % value)

        if tag.version == 6:
            if ip.version != 6:
                return ValidationResult(False, "invalid ipv6: %s" % value)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type ip not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=value)

    @staticmethod
    def validate_url(value: str, tag: Tag) -> ValidationResult:
        if value == "":
            if not tag.allow_empty:
                return ValidationResult(False, "type url not allow empty value \"\"")
            return ValidationResult(True, data=value, text="")

        from urllib.parse import urlparse
        parsed = urlparse(value)

        if parsed.scheme not in ("http", "https"):
            return ValidationResult(False, "invalid url: %s" % value)

        if parsed.hostname is None:
            return ValidationResult(False, "invalid url: %s" % value)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type url not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=value)

    @staticmethod
    def validate_email(value: str, tag: Tag) -> ValidationResult:
        if value == "":
            if not tag.allow_empty:
                return ValidationResult(False, "type email not allow empty value \"\"")
            return ValidationResult(True, data=value, text=value)

        if not MmValidator.email_regex.match(value):
            return ValidationResult(False, "value '%s' does not match email pattern" % value)

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type email not support location UTC%d" % location_offset)

        return ValidationResult(True, data=value, text=value)

    @staticmethod
    def validate_enum(value: str, tag: Tag) -> ValidationResult:
        if value == "":
            if not tag.allow_empty:
                return ValidationResult(False, "type enum not allow empty value \"\"")
            return ValidationResult(True, data=-1, text=value)

        if not tag.enum:
            return ValidationResult(False, "enum not defined")

        enums = tag.enum.split('|')
        idx = -1
        for i, enum_val in enumerate(enums):
            if enum_val.strip() == value:
                idx = i
                break

        if idx == -1:
            return ValidationResult(False, "value '%s' not found in enum: %s" % (value, enums))

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type enum not support location UTC%d" % location_offset)

        return ValidationResult(True, data=idx, text=value)

    @staticmethod
    def validate_image(value: bytes, tag: Tag) -> ValidationResult:
        length = len(value)

        if length == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type image not allow empty value []byte{}")
            return ValidationResult(True, data=value, text="")

        if tag.min:
            try:
                mini = int(tag.min)
                if length < mini:
                    return ValidationResult(False, "[]byte length %d < min %d" % (length, mini))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.min as int: %s" % tag.min)

        if tag.max:
            try:
                maxi = int(tag.max)
                if length > maxi:
                    return ValidationResult(False, "[]byte length %d > max %d" % (length, maxi))
            except ValueError:
                return ValidationResult(False, "failed to parse tag.max as int: %s" % tag.max)

        if tag.size > 0 and length != tag.size:
            return ValidationResult(False, "[]byte length %d != size %d" % (length, tag.size))

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        location_offset = tag.location if tag.location else 0
        try:
            location_offset = int(location_offset)
        except (ValueError, TypeError):
            location_offset = 0
        if location_offset != 0:
            return ValidationResult(False, "type image not support location UTC%d" % location_offset)

        text = base64.b64encode(value).decode('utf-8')

        return ValidationResult(True, data=value, text=text)

    @staticmethod
    def validate(value: Any, tag: Tag) -> ValidationResult:
        if tag.type == ValueType.Array:
            if isinstance(value, list):
                return MmValidator.validate_arr(value, tag)
            else:
                return ValidationResult(False, "expected array, got %s" % type(value).__name__)
        elif tag.type == ValueType.Slice:
            if isinstance(value, list):
                return MmValidator.validate_vec(value, tag)
            else:
                return ValidationResult(False, "expected slice, got %s" % type(value).__name__)
        elif tag.type == ValueType.Object:
            return MmValidator.validate_obj(tag)
        elif tag.type == ValueType.Map:
            return MmValidator.validate_map(tag)
        elif tag.type == ValueType.String:
            if isinstance(value, str):
                return MmValidator.validate_str(value, tag)
            else:
                return ValidationResult(False, "expected string, got %s" % type(value).__name__)
        elif tag.type == ValueType.Bytes:
            if isinstance(value, bytes):
                return MmValidator.validate_bytes(value, tag)
            else:
                return ValidationResult(False, "expected bytes, got %s" % type(value).__name__)
        elif tag.type == ValueType.Bool:
            if isinstance(value, bool):
                return MmValidator.validate_bool(value, tag)
            else:
                return ValidationResult(False, "expected bool, got %s" % type(value).__name__)
        elif tag.type == ValueType.Int:
            if isinstance(value, int):
                return MmValidator.validate_i(value, tag)
            else:
                return ValidationResult(False, "expected int, got %s" % type(value).__name__)
        elif tag.type == ValueType.Int8:
            if isinstance(value, int):
                return MmValidator.validate_i8(value, tag)
            else:
                return ValidationResult(False, "expected int8, got %s" % type(value).__name__)
        elif tag.type == ValueType.Int16:
            if isinstance(value, int):
                return MmValidator.validate_i16(value, tag)
            else:
                return ValidationResult(False, "expected int16, got %s" % type(value).__name__)
        elif tag.type == ValueType.Int32:
            if isinstance(value, int):
                return MmValidator.validate_i32(value, tag)
            else:
                return ValidationResult(False, "expected int32, got %s" % type(value).__name__)
        elif tag.type == ValueType.Int64:
            if isinstance(value, int):
                return MmValidator.validate_i64(value, tag)
            else:
                return ValidationResult(False, "expected int64, got %s" % type(value).__name__)
        elif tag.type == ValueType.Uint:
            if isinstance(value, int):
                return MmValidator.validate_u(value, tag)
            else:
                return ValidationResult(False, "expected uint, got %s" % type(value).__name__)
        elif tag.type == ValueType.Uint8:
            if isinstance(value, int):
                return MmValidator.validate_u8(value, tag)
            else:
                return ValidationResult(False, "expected uint8, got %s" % type(value).__name__)
        elif tag.type == ValueType.Uint16:
            if isinstance(value, int):
                return MmValidator.validate_u16(value, tag)
            else:
                return ValidationResult(False, "expected uint16, got %s" % type(value).__name__)
        elif tag.type == ValueType.Uint32:
            if isinstance(value, int):
                return MmValidator.validate_u32(value, tag)
            else:
                return ValidationResult(False, "expected uint32, got %s" % type(value).__name__)
        elif tag.type == ValueType.Uint64:
            if isinstance(value, int):
                return MmValidator.validate_u64(value, tag)
            else:
                return ValidationResult(False, "expected uint64, got %s" % type(value).__name__)
        elif tag.type == ValueType.Float32:
            if isinstance(value, (int, float)):
                return MmValidator.validate_f32(float(value), tag)
            else:
                return ValidationResult(False, "expected float32, got %s" % type(value).__name__)
        elif tag.type == ValueType.Float64:
            if isinstance(value, (int, float)):
                return MmValidator.validate_f64(float(value), tag)
            else:
                return ValidationResult(False, "expected float64, got %s" % type(value).__name__)
        elif tag.type == ValueType.BigInt:
            if isinstance(value, str):
                return MmValidator.validate_bigint(value, tag)
            else:
                return ValidationResult(False, "expected string for big.Int, got %s" % type(value).__name__)
        elif tag.type == ValueType.DateTime:
            return MmValidator.validate_datetime(value, tag)
        elif tag.type == ValueType.Date:
            return MmValidator.validate_date(value, tag)
        elif tag.type == ValueType.Time:
            return MmValidator.validate_time(value, tag)
        elif tag.type == ValueType.UUID:
            if isinstance(value, str):
                return MmValidator.validate_uuid(value, tag)
            else:
                return ValidationResult(False, "expected string, got %s" % type(value).__name__)
        elif tag.type == ValueType.Decimal:
            if isinstance(value, str):
                return MmValidator.validate_decimal(value, tag)
            else:
                return ValidationResult(False, "expected string, got %s" % type(value).__name__)
        elif tag.type == ValueType.IP:
            if isinstance(value, str):
                return MmValidator.validate_ip(value, tag)
            else:
                return ValidationResult(False, "expected string, got %s" % type(value).__name__)
        elif tag.type == ValueType.URL:
            if isinstance(value, str):
                return MmValidator.validate_url(value, tag)
            else:
                return ValidationResult(False, "expected string, got %s" % type(value).__name__)
        elif tag.type == ValueType.Email:
            if isinstance(value, str):
                return MmValidator.validate_email(value, tag)
            else:
                return ValidationResult(False, "expected string, got %s" % type(value).__name__)
        elif tag.type == ValueType.Enum:
            if isinstance(value, str):
                return MmValidator.validate_enum(value, tag)
            else:
                return ValidationResult(False, "expected string, got %s" % type(value).__name__)
        elif tag.type == ValueType.Image:
            if isinstance(value, bytes):
                return MmValidator.validate_image(value, tag)
            else:
                return ValidationResult(False, "expected bytes, got %s" % type(value).__name__)
        else:
            return ValidationResult(True, data=value, text=str(value))


validator = MmValidator