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
    email_regex = re.compile(r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$')
    decimal_regex = re.compile(r'^-?\d+\.\d+$')
    uuid_regex = re.compile(r'^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$')

    @staticmethod
    def validate_array(value: List[Any], tag: Tag) -> ValidationResult:
        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type array not support location UTC{tag.location}")

        length = len(value)

        if length == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type array not allow empty")
            return ValidationResult(True, data=value, text=str(value))

        if tag.size > 0 and length > tag.size:
            return ValidationResult(False, "type array over size")

        if tag.child_unique:
            seen = set()
            for i, item in enumerate(value):
                key = item if not isinstance(item, (dict, list)) else str(item)
                if key in seen:
                    return ValidationResult(False, f"array duplicate value found: {item}, index: {i}")
                seen.add(key)

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_struct(tag: Tag) -> ValidationResult:
        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type struct not support location UTC{tag.location}")

        return ValidationResult(True)

    @staticmethod
    def validate_string(value: str, tag: Tag) -> ValidationResult:
        if value == "":
            if not tag.allow_empty:
                return ValidationResult(False, f"type string not allow empty value \"{value}\"")
            return ValidationResult(True, data=value, text=value)

        if tag.pattern:
            try:
                regex = re.compile(tag.pattern)
                if not regex.match(value):
                    return ValidationResult(False, f"value \"{value}\" does not match pattern {tag.pattern}")
            except re.error as e:
                return ValidationResult(False, f"pattern \"{tag.pattern}\" compile err: {e}")

        length = len(value)

        if tag.min:
            try:
                mini = int(tag.min)
                if length < mini:
                    return ValidationResult(False, f"string length {length} is less than the minimum limit {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as int: {tag.min}")

        if tag.max:
            try:
                maxi = int(tag.max)
                if length > maxi:
                    return ValidationResult(False, f"string length {length} exceeds the maximum limit {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as int: {tag.max}")

        if tag.size > 0 and length != tag.size:
            return ValidationResult(False, f"string length {length} != size {tag.size}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type string not support location UTC{tag.location}")

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
                    return ValidationResult(False, f"[]byte length {length} is less than the minimum limit {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as int: {tag.min}")

        if tag.max:
            try:
                maxi = int(tag.max)
                if length > maxi:
                    return ValidationResult(False, f"[]byte length {length} exceeds the maximum limit {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as int: {tag.max}")

        if tag.size > 0 and length != tag.size:
            return ValidationResult(False, f"[]byte length {length} != size {tag.size}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type []byte not support location UTC{tag.location}")

        text = base64.b64encode(value).decode('utf-8')

        return ValidationResult(True, data=value, text=text)

    @staticmethod
    def validate_bool(value: bool, tag: Tag) -> ValidationResult:
        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.allow_empty:
            return ValidationResult(False, "type bool not support allow empty")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type bool not support location UTC{tag.location}")

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_int(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, f"type int not allow empty value {value}")
            return ValidationResult(True, data=value, text="0")

        if tag.min:
            try:
                mini = int(tag.min)
                if value < mini:
                    return ValidationResult(False, f"value {value} is less than the minimum limit {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as int: {tag.min}")

        if tag.max:
            try:
                maxi = int(tag.max)
                if value > maxi:
                    return ValidationResult(False, f"value {value} exceeds the maximum limit {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as int: {tag.max}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type int not support location UTC{tag.location}")

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_int8(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, f"type int8 not allow empty value {value}")
            return ValidationResult(True, data=value, text="0")

        if value < -128 or value > 127:
            return ValidationResult(False, f"value {value} out of int8 range [-128, 127]")

        if tag.min:
            try:
                mini = int(tag.min)
                if mini < -128 or mini > 127:
                    return ValidationResult(False, f"tag.min {mini} is out of int8 range [-128, 127]")
                if value < mini:
                    return ValidationResult(False, f"value {value} is less than the minimum limit {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as int8: {tag.min}")

        if tag.max:
            try:
                maxi = int(tag.max)
                if maxi < -128 or maxi > 127:
                    return ValidationResult(False, f"tag.max {maxi} is out of int8 range [-128, 127]")
                if value > maxi:
                    return ValidationResult(False, f"value {value} exceeds the maximum limit {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as int8: {tag.max}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type int8 not support location UTC{tag.location}")

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_int16(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, f"type int16 not allow empty value {value}")
            return ValidationResult(True, data=value, text="0")

        if value < -32768 or value > 32767:
            return ValidationResult(False, f"value {value} out of int16 range [-32768, 32767]")

        if tag.min:
            try:
                mini = int(tag.min)
                if mini < -32768 or mini > 32767:
                    return ValidationResult(False, f"tag.min {mini} is out of int16 range [-32768, 32767]")
                if value < mini:
                    return ValidationResult(False, f"value {value} is less than the minimum limit {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as int16: {tag.min}")

        if tag.max:
            try:
                maxi = int(tag.max)
                if maxi < -32768 or maxi > 32767:
                    return ValidationResult(False, f"tag.max {maxi} is out of int16 range [-32768, 32767]")
                if value > maxi:
                    return ValidationResult(False, f"value {value} exceeds the maximum limit {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as int16: {tag.max}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type int16 not support location UTC{tag.location}")

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_int32(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, f"type int32 not allow empty value {value}")
            return ValidationResult(True, data=value, text="0")

        if value < -2147483648 or value > 2147483647:
            return ValidationResult(False, f"value {value} out of int32 range [-2147483648, 2147483647]")

        if tag.min:
            try:
                mini = int(tag.min)
                if mini < -2147483648 or mini > 2147483647:
                    return ValidationResult(False, f"tag.min {mini} is out of int32 range [-2147483648, 2147483647]")
                if value < mini:
                    return ValidationResult(False, f"value {value} is less than the minimum limit {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as int32: {tag.min}")

        if tag.max:
            try:
                maxi = int(tag.max)
                if maxi < -2147483648 or maxi > 2147483647:
                    return ValidationResult(False, f"tag.max {maxi} is out of int32 range [-2147483648, 2147483647]")
                if value > maxi:
                    return ValidationResult(False, f"value {value} exceeds the maximum limit {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as int32: {tag.max}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type int32 not support location UTC{tag.location}")

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_int64(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, f"type int64 not allow empty value {value}")
            return ValidationResult(True, data=value, text="0")

        if value < -9223372036854775808 or value > 9223372036854775807:
            return ValidationResult(False, f"value {value} out of int64 range [-9223372036854775808, 9223372036854775807]")

        if tag.min:
            try:
                mini = int(tag.min)
                if mini < -9223372036854775808 or mini > 9223372036854775807:
                    return ValidationResult(False, f"tag.min {mini} is out of int64 range [-9223372036854775808, 9223372036854775807]")
                if value < mini:
                    return ValidationResult(False, f"value {value} is less than the minimum limit {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as int64: {tag.min}")

        if tag.max:
            try:
                maxi = int(tag.max)
                if maxi < -9223372036854775808 or maxi > 9223372036854775807:
                    return ValidationResult(False, f"tag.max {maxi} is out of int64 range [-9223372036854775808, 9223372036854775807]")
                if value > maxi:
                    return ValidationResult(False, f"value {value} exceeds the maximum limit {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as int64: {tag.max}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type int64 not support location UTC{tag.location}")

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_uint(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, f"type uint not allow empty value {value}")
            return ValidationResult(True, data=value, text="0")

        if value < 0 or value > 4294967295:
            return ValidationResult(False, f"value {value} out of uint range [0, 4294967295]")

        if tag.min:
            try:
                mini = int(tag.min)
                if mini < 0 or mini > 4294967295:
                    return ValidationResult(False, f"tag.min {mini} is out of uint range [0, 4294967295]")
                if value < mini:
                    return ValidationResult(False, f"value {value} is less than the minimum limit {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as uint: {tag.min}")

        if tag.max:
            try:
                maxi = int(tag.max)
                if maxi < 0 or maxi > 4294967295:
                    return ValidationResult(False, f"tag.max {maxi} is out of uint range [0, 4294967295]")
                if value > maxi:
                    return ValidationResult(False, f"value {value} exceeds the maximum limit {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as uint: {tag.max}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type uint not support location UTC{tag.location}")

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_uint8(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, f"type uint8 not allow empty value {value}")
            return ValidationResult(True, data=value, text="0")

        if value < 0 or value > 255:
            return ValidationResult(False, f"value {value} out of uint8 range [0, 255]")

        if tag.min:
            try:
                mini = int(tag.min)
                if mini < 0 or mini > 255:
                    return ValidationResult(False, f"tag.min {mini} is out of uint8 range [0, 255]")
                if value < mini:
                    return ValidationResult(False, f"value {value} is less than the minimum limit {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as uint8: {tag.min}")

        if tag.max:
            try:
                maxi = int(tag.max)
                if maxi < 0 or maxi > 255:
                    return ValidationResult(False, f"tag.max {maxi} is out of uint8 range [0, 255]")
                if value > maxi:
                    return ValidationResult(False, f"value {value} exceeds the maximum limit {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as uint8: {tag.max}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type uint8 not support location UTC{tag.location}")

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_uint16(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, f"type uint16 not allow empty value {value}")
            return ValidationResult(True, data=value, text="0")

        if value < 0 or value > 65535:
            return ValidationResult(False, f"value {value} out of uint16 range [0, 65535]")

        if tag.min:
            try:
                mini = int(tag.min)
                if mini < 0 or mini > 65535:
                    return ValidationResult(False, f"tag.min {mini} is out of uint16 range [0, 65535]")
                if value < mini:
                    return ValidationResult(False, f"value {value} is less than the minimum limit {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as uint16: {tag.min}")

        if tag.max:
            try:
                maxi = int(tag.max)
                if maxi < 0 or maxi > 65535:
                    return ValidationResult(False, f"tag.max {maxi} is out of uint16 range [0, 65535]")
                if value > maxi:
                    return ValidationResult(False, f"value {value} exceeds the maximum limit {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as uint16: {tag.max}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type uint16 not support location UTC{tag.location}")

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_uint32(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, f"type uint32 not allow empty value {value}")
            return ValidationResult(True, data=value, text="0")

        if value < 0 or value > 4294967295:
            return ValidationResult(False, f"value {value} out of uint32 range [0, 4294967295]")

        if tag.min:
            try:
                mini = int(tag.min)
                if mini < 0 or mini > 4294967295:
                    return ValidationResult(False, f"tag.min {mini} is out of uint32 range [0, 4294967295]")
                if value < mini:
                    return ValidationResult(False, f"value {value} is less than the minimum limit {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as uint32: {tag.min}")

        if tag.max:
            try:
                maxi = int(tag.max)
                if maxi < 0 or maxi > 4294967295:
                    return ValidationResult(False, f"tag.max {maxi} is out of uint32 range [0, 4294967295]")
                if value > maxi:
                    return ValidationResult(False, f"value {value} exceeds the maximum limit {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as uint32: {tag.max}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type uint32 not support location UTC{tag.location}")

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_uint64(value: int, tag: Tag) -> ValidationResult:
        if value == 0:
            if not tag.allow_empty:
                return ValidationResult(False, f"type uint64 not allow empty value {value}")
            return ValidationResult(True, data=value, text="0")

        if value < 0 or value > 18446744073709551615:
            return ValidationResult(False, f"value {value} out of uint64 range [0, 18446744073709551615]")

        if tag.min:
            try:
                mini = int(tag.min)
                if mini < 0 or mini > 18446744073709551615:
                    return ValidationResult(False, f"tag.min {mini} is out of uint64 range [0, 18446744073709551615]")
                if value < mini:
                    return ValidationResult(False, f"value {value} is less than the minimum limit {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as uint64: {tag.min}")

        if tag.max:
            try:
                maxi = int(tag.max)
                if maxi < 0 or maxi > 18446744073709551615:
                    return ValidationResult(False, f"tag.max {maxi} is out of uint64 range [0, 18446744073709551615]")
                if value > maxi:
                    return ValidationResult(False, f"value {value} exceeds the maximum limit {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as uint64: {tag.max}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type uint64 not support location UTC{tag.location}")

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_float32(value: float, tag: Tag) -> ValidationResult:
        if value == 0.0:
            if not tag.allow_empty:
                return ValidationResult(False, "type float32 not allow empty value 0.0")
            return ValidationResult(True, data=value, text="0.0")

        if tag.min:
            try:
                mini = float(tag.min)
                if mini < -3.4e38 or mini > 3.4e38:
                    return ValidationResult(False, f"tag.min {mini} is out of float32 range [-3.4e38, 3.4e38]")
                if value < mini:
                    return ValidationResult(False, f"value {value} is less than the minimum limit {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as float32: {tag.min}")

        if tag.max:
            try:
                maxi = float(tag.max)
                if maxi < -3.4e38 or maxi > 3.4e38:
                    return ValidationResult(False, f"tag.max {maxi} is out of float32 range [-3.4e38, 3.4e38]")
                if value > maxi:
                    return ValidationResult(False, f"value {value} exceeds the maximum limit {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as float32: {tag.max}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type float32 not support location UTC{tag.location}")

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_float64(value: float, tag: Tag) -> ValidationResult:
        if value == 0.0:
            if not tag.allow_empty:
                return ValidationResult(False, "type float64 not allow empty value 0.0")
            return ValidationResult(True, data=value, text="0.0")

        if tag.min:
            try:
                mini = float(tag.min)
                if mini < -1.7976931348623157e308 or mini > 1.7976931348623157e308:
                    return ValidationResult(False, f"tag.min {mini} is out of float64 range [-1.7976931348623157e308, 1.7976931348623157e308]")
                if value < mini:
                    return ValidationResult(False, f"value {value} is less than the minimum limit {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as float64: {tag.min}")

        if tag.max:
            try:
                maxi = float(tag.max)
                if maxi < -1.7976931348623157e308 or maxi > 1.7976931348623157e308:
                    return ValidationResult(False, f"tag.max {maxi} is out of float64 range [-1.7976931348623157e308, 1.7976931348623157e308]")
                if value > maxi:
                    return ValidationResult(False, f"value {value} exceeds the maximum limit {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as float64: {tag.max}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type float64 not support location UTC{tag.location}")

        return ValidationResult(True, data=value, text=str(value))

    @staticmethod
    def validate_big_int(value: str, tag: Tag) -> ValidationResult:
        try:
            from decimal import Decimal
            val = Decimal(value)
        except Exception:
            return ValidationResult(False, f"invalid big.Int value: {value}")

        if val == 0:
            if not tag.allow_empty:
                return ValidationResult(False, "type big.Int not allow empty value 0")
            return ValidationResult(True, data=value, text="0")

        if tag.min:
            try:
                from decimal import Decimal
                mini = Decimal(tag.min)
                if val < mini:
                    return ValidationResult(False, f"big.Int {value} < min {tag.min}")
            except Exception:
                return ValidationResult(False, f"invalid min {tag.min} for big.Int")

        if tag.max:
            try:
                from decimal import Decimal
                maxi = Decimal(tag.max)
                if val > maxi:
                    return ValidationResult(False, f"big.Int {value} > max {tag.max}")
            except Exception:
                return ValidationResult(False, f"invalid max {tag.max} for big.Int")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type big.Int not support location UTC{tag.location}")

        return ValidationResult(True, data=value, text=str(val))

    @staticmethod
    def validate_datetime(value: Any, tag: Tag) -> ValidationResult:
        if not hasattr(value, 'timestamp'):
            return ValidationResult(False, f"expected datetime object, got {type(value).__name__}")

        if value.timestamp() == 0:
            if not tag.allow_empty:
                return ValidationResult(False, f"type datetime not allow empty {value}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        format_str = value.strftime("%Y-%m-%d %H:%M:%S")

        return ValidationResult(True, data=value, text=format_str)

    @staticmethod
    def validate_uuid(value: str, tag: Tag) -> ValidationResult:
        if value == "":
            if not tag.allow_empty:
                return ValidationResult(False, "type uuid not allow empty value \"\"")
            return ValidationResult(True, data=b'\x00' * 16, text=value)

        if not MmValidator.uuid_regex.match(value):
            return ValidationResult(False, f"value '{value}' does not match UUID pattern")

        try:
            uuid_bytes = bytes.fromhex(value.replace('-', ''))
        except ValueError:
            return ValidationResult(False, f"invalid uuid: {value}")

        if tag.version > 0:
            version = int(value[14], 16)
            if tag.version != version:
                return ValidationResult(False, "invalid uuid version")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type uuid not support location UTC{tag.location}")

        return ValidationResult(True, data=uuid_bytes, text=value)

    @staticmethod
    def validate_email(value: str, tag: Tag) -> ValidationResult:
        if value == "":
            if not tag.allow_empty:
                return ValidationResult(False, "type email not allow empty value \"\"")
            return ValidationResult(True, data=value, text=value)

        if not MmValidator.email_regex.match(value):
            return ValidationResult(False, f"value '{value}' does not match email pattern")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type email not support location UTC{tag.location}")

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
            return ValidationResult(False, f"value '{value}' not found in enum: {enums}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type enum not support location UTC{tag.location}")

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
                    return ValidationResult(False, f"[]byte length {length} < min {mini}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.min as int: {tag.min}")

        if tag.max:
            try:
                maxi = int(tag.max)
                if length > maxi:
                    return ValidationResult(False, f"[]byte length {length} > max {maxi}")
            except ValueError:
                return ValidationResult(False, f"failed to parse tag.max as int: {tag.max}")

        if tag.size > 0 and length != tag.size:
            return ValidationResult(False, f"[]byte length {length} != size {tag.size}")

        if len(tag.desc) > 65535:
            return ValidationResult(False, "desc length exceeds 65535 bytes")

        if tag.location is not None and tag.location != 0:
            return ValidationResult(False, f"type image not support location UTC{tag.location}")

        text = base64.b64encode(value).decode('utf-8')

        return ValidationResult(True, data=value, text=text)

    @staticmethod
    def validate(value: Any, tag: Tag) -> ValidationResult:
        if tag.type == ValueType.Array:
            if isinstance(value, list):
                return MmValidator.validate_array(value, tag)
            else:
                return ValidationResult(False, f"expected array, got {type(value).__name__}")
        elif tag.type == ValueType.Struct:
            return MmValidator.validate_struct(tag)
        elif tag.type == ValueType.String:
            if isinstance(value, str):
                return MmValidator.validate_string(value, tag)
            else:
                return ValidationResult(False, f"expected string, got {type(value).__name__}")
        elif tag.type == ValueType.Bytes:
            if isinstance(value, bytes):
                return MmValidator.validate_bytes(value, tag)
            else:
                return ValidationResult(False, f"expected bytes, got {type(value).__name__}")
        elif tag.type == ValueType.Bool:
            if isinstance(value, bool):
                return MmValidator.validate_bool(value, tag)
            else:
                return ValidationResult(False, f"expected bool, got {type(value).__name__}")
        elif tag.type == ValueType.Int:
            if isinstance(value, int):
                return MmValidator.validate_int(value, tag)
            else:
                return ValidationResult(False, f"expected int, got {type(value).__name__}")
        elif tag.type == ValueType.Int8:
            if isinstance(value, int):
                return MmValidator.validate_int8(value, tag)
            else:
                return ValidationResult(False, f"expected int8, got {type(value).__name__}")
        elif tag.type == ValueType.Int16:
            if isinstance(value, int):
                return MmValidator.validate_int16(value, tag)
            else:
                return ValidationResult(False, f"expected int16, got {type(value).__name__}")
        elif tag.type == ValueType.Int32:
            if isinstance(value, int):
                return MmValidator.validate_int32(value, tag)
            else:
                return ValidationResult(False, f"expected int32, got {type(value).__name__}")
        elif tag.type == ValueType.Int64:
            if isinstance(value, int):
                return MmValidator.validate_int64(value, tag)
            else:
                return ValidationResult(False, f"expected int64, got {type(value).__name__}")
        elif tag.type == ValueType.Uint:
            if isinstance(value, int):
                return MmValidator.validate_uint(value, tag)
            else:
                return ValidationResult(False, f"expected uint, got {type(value).__name__}")
        elif tag.type == ValueType.Uint8:
            if isinstance(value, int):
                return MmValidator.validate_uint8(value, tag)
            else:
                return ValidationResult(False, f"expected uint8, got {type(value).__name__}")
        elif tag.type == ValueType.Uint16:
            if isinstance(value, int):
                return MmValidator.validate_uint16(value, tag)
            else:
                return ValidationResult(False, f"expected uint16, got {type(value).__name__}")
        elif tag.type == ValueType.Uint32:
            if isinstance(value, int):
                return MmValidator.validate_uint32(value, tag)
            else:
                return ValidationResult(False, f"expected uint32, got {type(value).__name__}")
        elif tag.type == ValueType.Uint64:
            if isinstance(value, int):
                return MmValidator.validate_uint64(value, tag)
            else:
                return ValidationResult(False, f"expected uint64, got {type(value).__name__}")
        elif tag.type == ValueType.Float32:
            if isinstance(value, float):
                return MmValidator.validate_float32(value, tag)
            else:
                return ValidationResult(False, f"expected float32, got {type(value).__name__}")
        elif tag.type == ValueType.Float64:
            if isinstance(value, float):
                return MmValidator.validate_float64(value, tag)
            else:
                return ValidationResult(False, f"expected float64, got {type(value).__name__}")
        elif tag.type == ValueType.BigInt:
            if isinstance(value, str):
                return MmValidator.validate_big_int(value, tag)
            else:
                return ValidationResult(False, f"expected string for big.Int, got {type(value).__name__}")
        elif tag.type in (ValueType.DateTime, ValueType.Date, ValueType.Time):
            return MmValidator.validate_datetime(value, tag)
        elif tag.type == ValueType.UUID:
            if isinstance(value, str):
                return MmValidator.validate_uuid(value, tag)
            else:
                return ValidationResult(False, f"expected string, got {type(value).__name__}")
        elif tag.type == ValueType.Email:
            if isinstance(value, str):
                return MmValidator.validate_email(value, tag)
            else:
                return ValidationResult(False, f"expected string, got {type(value).__name__}")
        elif tag.type == ValueType.Enum:
            if isinstance(value, str):
                return MmValidator.validate_enum(value, tag)
            else:
                return ValidationResult(False, f"expected string, got {type(value).__name__}")
        elif tag.type == ValueType.Image:
            if isinstance(value, bytes):
                return MmValidator.validate_image(value, tag)
            else:
                return ValidationResult(False, f"expected bytes, got {type(value).__name__}")
        else:
            return ValidationResult(True, data=value, text=str(value))

validator = MmValidator
