using System.Net;
using System.Net.Sockets;
using System.Numerics;
using System.Text.RegularExpressions;
using MetaMessage.Ir;
using ValueType = MetaMessage.Ir.ValueType;

namespace MetaMessage.Core;

public class ValidationResult
{
    public bool IsValid { get; set; }
    public List<string> Errors { get; set; }

    public ValidationResult()
    {
        IsValid = true;
        Errors = new List<string>();
    }

    public void AddError(string error)
    {
        IsValid = false;
        Errors.Add(error);
    }
}

public class MmValidator
{
    private static readonly Regex EmailRegex = new Regex(@"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$");
    private static readonly Regex DecimalRegex = new Regex(@"^-?\d+\.\d+$");
    private static readonly Regex UuidRegex = new Regex(@"^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$");

    private bool ShouldSkipValidate(Tag tag)
    {
        return tag.Example;
    }

    public ValidationResult Validate(dynamic value, Tag tag)
    {
        var result = new ValidationResult();

        if (value == null)
        {
            if (!tag.Nullable)
            {
                result.AddError("value is required");
            }
            return result;
        }

        if (ShouldSkipValidate(tag))
        {
            return result;
        }

        switch (tag.Type)
        {
            case ValueType.Bool:
                ValidateBool(value, tag, result);
                break;
            case ValueType.I:
                ValidateI(value, tag, result);
                break;
            case ValueType.I8:
                ValidateI8(value, tag, result);
                break;
            case ValueType.I16:
                ValidateI16(value, tag, result);
                break;
            case ValueType.I32:
                ValidateI32(value, tag, result);
                break;
            case ValueType.I64:
                ValidateI64(value, tag, result);
                break;
            case ValueType.U:
                ValidateU(value, tag, result);
                break;
            case ValueType.U8:
                ValidateU8(value, tag, result);
                break;
            case ValueType.U16:
                ValidateU16(value, tag, result);
                break;
            case ValueType.U32:
                ValidateU32(value, tag, result);
                break;
            case ValueType.U64:
                ValidateU64(value, tag, result);
                break;
            case ValueType.Bigint:
                ValidateBigint(value, tag, result);
                break;
            case ValueType.F32:
                ValidateF32(value, tag, result);
                break;
            case ValueType.F64:
                ValidateF64(value, tag, result);
                break;
            case ValueType.Decimal:
                ValidateDecimal(value, tag, result);
                break;
            case ValueType.Str:
                ValidateStr(value, tag, result);
                break;
            case ValueType.Email:
                ValidateEmail(value, tag, result);
                break;
            case ValueType.Url:
                ValidateUrl(value, tag, result);
                break;
            case ValueType.Bytes:
                ValidateBytes(value, tag, result);
                break;
            case ValueType.Uuid:
                ValidateUuid(value, tag, result);
                break;
            case ValueType.Ip:
                ValidateIp(value, tag, result);
                break;
            case ValueType.Media:
                ValidateMedia(value, tag, result);
                break;
            case ValueType.Datetime:
            case ValueType.Date:
            case ValueType.Time:
                ValidateDatetime(value, tag, result);
                break;
            case ValueType.Enums:
                ValidateEnum(value, tag, result);
                break;
            case ValueType.Arr:
            case ValueType.Vec:
                ValidateArr(value, tag, result);
                break;
            case ValueType.Obj:
                ValidateObj(value, tag, result);
                break;
            case ValueType.Map:
                ValidateMap(value, tag, result);
                break;
        }

        return result;
    }

    public void ValidateVec(dynamic value, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (tag.Location != 0)
        {
            result.AddError($"type slice not support location UTC{tag.Location}");
            return;
        }

        if (value is System.Collections.IEnumerable enumerable)
        {
            int count = 0;
            foreach (var _ in enumerable)
            {
                count++;
            }

            if (count == 0)
            {
                if (tag.AllowEmpty)
                {
                    return;
                }
                result.AddError("type slice not allow empty");
                return;
            }

            if (tag.Size > 0 && count != tag.Size)
            {
                result.AddError($"size mismatch, want={tag.Size}, got={count}");
                return;
            }

            if (tag.ChildUnique)
            {
                var seen = new Dictionary<object, bool>();
                int index = 0;
                foreach (var item in enumerable)
                {
                    if (item != null)
                    {
                        if (seen.ContainsKey(item))
                        {
                            result.AddError($"vec duplicate value found: {item}, index: {index}");
                            return;
                        }
                        seen[item] = true;
                    }
                    index++;
                }
            }
        }
        else
        {
            result.AddError("value must be an array");
        }
    }

    public void ValidateArr(dynamic value, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (tag.Location != 0)
        {
            result.AddError($"type array not support location UTC{tag.Location}");
            return;
        }

        if (value is System.Collections.IEnumerable enumerable)
        {
            int count = 0;
            foreach (var _ in enumerable)
            {
                count++;
            }

            if (count == 0)
            {
                if (tag.AllowEmpty)
                {
                    return;
                }
                result.AddError("type array not allow empty");
                return;
            }

            if (tag.Size > 0 && count != tag.Size)
            {
                result.AddError($"size mismatch, want={tag.Size}, got={count}");
                return;
            }

            if (tag.ChildUnique)
            {
                var seen = new Dictionary<object, bool>();
                int index = 0;
                foreach (var item in enumerable)
                {
                    if (item != null)
                    {
                        if (seen.ContainsKey(item))
                        {
                            result.AddError($"array duplicate value found: {item}, index: {index}");
                            return;
                        }
                        seen[item] = true;
                    }
                    index++;
                }
            }
        }
        else
        {
            result.AddError("value must be an array");
        }
    }

    public void ValidateObj(dynamic value, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (tag.Location != 0)
        {
            result.AddError($"type struct not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateMap(dynamic value, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (tag.Location != 0)
        {
            result.AddError($"type map not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateStr(dynamic value, Tag tag, ValidationResult result)
    {
        if (!(value is string strVal))
        {
            result.AddError("value must be a string");
            return;
        }

        ValidateStr(strVal, tag, result);
    }

    private void ValidateStr(string val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (string.IsNullOrEmpty(val))
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Pattern))
        {
            try
            {
                var re = new Regex(tag.Pattern);
                if (!re.IsMatch(val))
                {
                    result.AddError($"value \"{val}\" does not match pattern {tag.Pattern}");
                    return;
                }
            }
            catch (Exception)
            {
                result.AddError($"pattern \"{tag.Pattern}\" compile error");
                return;
            }
        }

        int l = val.Length;

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (int.TryParse(tag.Min, out int minVal))
            {
                if (l < minVal)
                {
                    result.AddError($"string length {l} is less than the minimum limit {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as int: {tag.Min}");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (int.TryParse(tag.Max, out int maxVal))
            {
                if (l > maxVal)
                {
                    result.AddError($"string length {l} exceeds the maximum limit {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as int: {tag.Max}");
                return;
            }
        }

        if (tag.Size != 0 && l != tag.Size)
        {
            result.AddError($"string length {l} != size {tag.Size}");
            return;
        }

        if (tag.Location != 0)
        {
            result.AddError($"type string not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateBytes(dynamic value, Tag tag, ValidationResult result)
    {
        if (!(value is byte[] bytesVal))
        {
            result.AddError("value must be a byte array");
            return;
        }

        ValidateBytes(bytesVal, tag, result);
    }

    private void ValidateBytes(byte[] val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        int l = val.Length;

        if (l == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (int.TryParse(tag.Min, out int minVal))
            {
                if (l < minVal)
                {
                    result.AddError($"[]byte length {l} is less than the minimum limit {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as int: {tag.Min}");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (int.TryParse(tag.Max, out int maxVal))
            {
                if (l > maxVal)
                {
                    result.AddError($"[]byte length {l} exceeds the maximum limit {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as int: {tag.Max}");
                return;
            }
        }

        if (tag.Size != 0 && l != tag.Size)
        {
            result.AddError($"[]byte length {l} != size {tag.Size}");
            return;
        }

        if (tag.Location != 0)
        {
            result.AddError($"type []byte not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateBool(dynamic value, Tag tag, ValidationResult result)
    {
        if (!(value is bool boolVal))
        {
            result.AddError("value must be a boolean");
            return;
        }

        ValidateBool(boolVal, tag, result);
    }

    private void ValidateBool(bool val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (tag.AllowEmpty)
        {
            result.AddError("type bool not support 'allow_empty' tag");
            return;
        }

        if (tag.Location != 0)
        {
            result.AddError($"type bool not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateI(dynamic value, Tag tag, ValidationResult result)
    {
        if (!TryGetInt64Value(value, out long longVal))
        {
            result.AddError("value must be a number");
            return;
        }

        ValidateI((int)longVal, tag, result);
    }

    private void ValidateI(int val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (val == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        long val64 = val;

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (long.TryParse(tag.Min, out long minVal))
            {
                if (val64 < minVal)
                {
                    result.AddError($"value {val64} is less than the minimum limit {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as int: {tag.Min}");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (long.TryParse(tag.Max, out long maxVal))
            {
                if (val64 > maxVal)
                {
                    result.AddError($"value {val64} exceeds the maximum limit {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as int: {tag.Max}");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type int not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateI8(dynamic value, Tag tag, ValidationResult result)
    {
        if (!TryGetInt64Value(value, out long longVal))
        {
            result.AddError("value must be a number");
            return;
        }

        if (longVal < sbyte.MinValue || longVal > sbyte.MaxValue)
        {
            result.AddError($"value {longVal} out of range for int8 ({sbyte.MinValue} to {sbyte.MaxValue})");
            return;
        }

        ValidateI8((sbyte)longVal, tag, result);
    }

    private void ValidateI8(sbyte val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (val == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        long val64 = val;

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (long.TryParse(tag.Min, out long minVal))
            {
                if (val64 < minVal)
                {
                    result.AddError($"value {val64} is less than the minimum limit {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as int8: {tag.Min}");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (long.TryParse(tag.Max, out long maxVal))
            {
                if (val64 > maxVal)
                {
                    result.AddError($"value {val64} exceeds the maximum limit {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as int8: {tag.Max}");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type int8 not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateI16(dynamic value, Tag tag, ValidationResult result)
    {
        if (!TryGetInt64Value(value, out long longVal))
        {
            result.AddError("value must be a number");
            return;
        }

        if (longVal < short.MinValue || longVal > short.MaxValue)
        {
            result.AddError($"value {longVal} out of range for int16 ({short.MinValue} to {short.MaxValue})");
            return;
        }

        ValidateI16((short)longVal, tag, result);
    }

    private void ValidateI16(short val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (val == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        long val64 = val;

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (long.TryParse(tag.Min, out long minVal))
            {
                if (val64 < minVal)
                {
                    result.AddError($"value {val64} is less than the minimum limit {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as int16: {tag.Min}");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (long.TryParse(tag.Max, out long maxVal))
            {
                if (val64 > maxVal)
                {
                    result.AddError($"value {val64} exceeds the maximum limit {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as int16: {tag.Max}");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type int16 not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateI32(dynamic value, Tag tag, ValidationResult result)
    {
        if (!TryGetInt64Value(value, out long longVal))
        {
            result.AddError("value must be a number");
            return;
        }

        if (longVal < int.MinValue || longVal > int.MaxValue)
        {
            result.AddError($"value {longVal} out of range for int32 ({int.MinValue} to {int.MaxValue})");
            return;
        }

        ValidateI32((int)longVal, tag, result);
    }

    private void ValidateI32(int val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (val == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        long val64 = val;

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (long.TryParse(tag.Min, out long minVal))
            {
                if (val64 < minVal)
                {
                    result.AddError($"value {val64} is less than the minimum limit {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as int32: {tag.Min}");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (long.TryParse(tag.Max, out long maxVal))
            {
                if (val64 > maxVal)
                {
                    result.AddError($"value {val64} exceeds the maximum limit {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as int32: {tag.Max}");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type int32 not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateI64(dynamic value, Tag tag, ValidationResult result)
    {
        if (!TryGetInt64Value(value, out long longVal))
        {
            result.AddError("value must be a number");
            return;
        }

        ValidateI64(longVal, tag, result);
    }

    private void ValidateI64(long val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (val == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (long.TryParse(tag.Min, out long minVal))
            {
                if (val < minVal)
                {
                    result.AddError($"value {val} is less than the minimum limit {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as int64: {tag.Min}");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (long.TryParse(tag.Max, out long maxVal))
            {
                if (val > maxVal)
                {
                    result.AddError($"value {val} exceeds the maximum limit {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as int64: {tag.Max}");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type int64 not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateU(dynamic value, Tag tag, ValidationResult result)
    {
        if (!TryGetUInt64Value(value, out ulong ulongVal))
        {
            result.AddError("value must be a number");
            return;
        }

        if (ulongVal > uint.MaxValue)
        {
            result.AddError($"value {ulongVal} out of range for uint (0 to {uint.MaxValue})");
            return;
        }

        ValidateU((uint)ulongVal, tag, result);
    }

    private void ValidateU(uint val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (val == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        ulong val64 = val;

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (ulong.TryParse(tag.Min, out ulong minVal))
            {
                if (val64 < minVal)
                {
                    result.AddError($"value {val} is less than the minimum limit {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as uint: {tag.Min}");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (ulong.TryParse(tag.Max, out ulong maxVal))
            {
                if (val64 > maxVal)
                {
                    result.AddError($"value {val} exceeds the maximum limit {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as uint: {tag.Max}");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type uint not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateU8(dynamic value, Tag tag, ValidationResult result)
    {
        if (!TryGetUInt64Value(value, out ulong ulongVal))
        {
            result.AddError("value must be a number");
            return;
        }

        if (ulongVal > byte.MaxValue)
        {
            result.AddError($"value {ulongVal} out of range for uint8 (0 to {byte.MaxValue})");
            return;
        }

        ValidateU8((byte)ulongVal, tag, result);
    }

    private void ValidateU8(byte val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (val == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        ulong val64 = val;

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (ulong.TryParse(tag.Min, out ulong minVal))
            {
                if (val64 < minVal)
                {
                    result.AddError($"value {val} is less than the minimum limit {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as uint8: {tag.Min}");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (ulong.TryParse(tag.Max, out ulong maxVal))
            {
                if (val64 > maxVal)
                {
                    result.AddError($"value {val} exceeds the maximum limit {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as uint8: {tag.Max}");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type uint8 not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateU16(dynamic value, Tag tag, ValidationResult result)
    {
        if (!TryGetUInt64Value(value, out ulong ulongVal))
        {
            result.AddError("value must be a number");
            return;
        }

        if (ulongVal > ushort.MaxValue)
        {
            result.AddError($"value {ulongVal} out of range for uint16 (0 to {ushort.MaxValue})");
            return;
        }

        ValidateU16((ushort)ulongVal, tag, result);
    }

    private void ValidateU16(ushort val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (val == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        ulong val64 = val;

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (ulong.TryParse(tag.Min, out ulong minVal))
            {
                if (val64 < minVal)
                {
                    result.AddError($"value {val} is less than the minimum limit {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as uint16: {tag.Min}");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (ulong.TryParse(tag.Max, out ulong maxVal))
            {
                if (val64 > maxVal)
                {
                    result.AddError($"value {val} exceeds the maximum limit {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as uint16: {tag.Max}");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type uint16 not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateU32(dynamic value, Tag tag, ValidationResult result)
    {
        if (!TryGetUInt64Value(value, out ulong ulongVal))
        {
            result.AddError("value must be a number");
            return;
        }

        if (ulongVal > uint.MaxValue)
        {
            result.AddError($"value {ulongVal} out of range for uint32 (0 to {uint.MaxValue})");
            return;
        }

        ValidateU32((uint)ulongVal, tag, result);
    }

    private void ValidateU32(uint val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (val == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        ulong val64 = val;

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (ulong.TryParse(tag.Min, out ulong minVal))
            {
                if (val64 < minVal)
                {
                    result.AddError($"value {val} is less than the minimum limit {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as uint32: {tag.Min}");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (ulong.TryParse(tag.Max, out ulong maxVal))
            {
                if (val64 > maxVal)
                {
                    result.AddError($"value {val} exceeds the maximum limit {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as uint32: {tag.Max}");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type uint32 not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateU64(dynamic value, Tag tag, ValidationResult result)
    {
        if (!TryGetUInt64Value(value, out ulong ulongVal))
        {
            result.AddError("value must be a number");
            return;
        }

        ValidateU64(ulongVal, tag, result);
    }

    private void ValidateU64(ulong val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (val == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (ulong.TryParse(tag.Min, out ulong minVal))
            {
                if (val < minVal)
                {
                    result.AddError($"value {val} is less than the minimum limit {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as uint64: {tag.Min}");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (ulong.TryParse(tag.Max, out ulong maxVal))
            {
                if (val > maxVal)
                {
                    result.AddError($"value {val} exceeds the maximum limit {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as uint64: {tag.Max}");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type uint64 not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateF32(dynamic value, Tag tag, ValidationResult result)
    {
        double doubleVal;
        if (value is double d)
        {
            doubleVal = d;
        }
        else if (value is float f)
        {
            doubleVal = f;
        }
        else if (value is int || value is long || value is short || value is sbyte || value is byte)
        {
            doubleVal = Convert.ToDouble(value);
        }
        else
        {
            result.AddError("value must be a float");
            return;
        }

        if (doubleVal < float.MinValue || doubleVal > float.MaxValue)
        {
            result.AddError($"value {doubleVal} out of range for float32 ({float.MinValue} to {float.MaxValue})");
            return;
        }

        ValidateF32((float)doubleVal, tag, result);
    }

    private void ValidateF32(float val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (val == 0.0f)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        double val64 = val;

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (double.TryParse(tag.Min, out double minVal))
            {
                if (val64 < minVal)
                {
                    result.AddError($"{val64} < min {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as float32: {tag.Min}");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (double.TryParse(tag.Max, out double maxVal))
            {
                if (val64 > maxVal)
                {
                    result.AddError($"{val64} > max {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as float32: {tag.Max}");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type float32 not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateF64(dynamic value, Tag tag, ValidationResult result)
    {
        double doubleVal;
        if (value is double d)
        {
            doubleVal = d;
        }
        else if (value is float f)
        {
            doubleVal = f;
        }
        else if (value is int || value is long || value is short || value is sbyte || value is byte)
        {
            doubleVal = Convert.ToDouble(value);
        }
        else
        {
            result.AddError("value must be a float");
            return;
        }

        ValidateF64(doubleVal, tag, result);
    }

    private void ValidateF64(double val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (val == 0.0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (double.TryParse(tag.Min, out double minVal))
            {
                if (val < minVal)
                {
                    result.AddError($"{val} < min {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as float64: {tag.Min}");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (double.TryParse(tag.Max, out double maxVal))
            {
                if (val > maxVal)
                {
                    result.AddError($"{val} > max {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as float64: {tag.Max}");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type float64 not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateBigint(dynamic value, Tag tag, ValidationResult result)
    {
        BigInteger bigIntVal;
        if (value is BigInteger bi)
        {
            bigIntVal = bi;
        }
        else if (value is string strVal && BigInteger.TryParse(strVal, out bigIntVal))
        {
        }
        else if (value is long l)
        {
            bigIntVal = new BigInteger(l);
        }
        else if (value is int i)
        {
            bigIntVal = new BigInteger(i);
        }
        else if (value is short s)
        {
            bigIntVal = new BigInteger(s);
        }
        else if (value is sbyte sb)
        {
            bigIntVal = new BigInteger(sb);
        }
        else if (value is byte b)
        {
            bigIntVal = new BigInteger(b);
        }
        else if (value is ulong ul)
        {
            bigIntVal = new BigInteger(ul);
        }
        else if (value is uint ui)
        {
            bigIntVal = new BigInteger(ui);
        }
        else if (value is ushort us)
        {
            bigIntVal = new BigInteger(us);
        }
        else
        {
            result.AddError("value must be a BigInteger");
            return;
        }

        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (bigIntVal.IsZero)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (BigInteger.TryParse(tag.Min, out BigInteger minVal))
            {
                if (bigIntVal < minVal)
                {
                    result.AddError($"big.Int length {bigIntVal} < min {minVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"invalid min \"{tag.Min}\" for big.Int");
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (BigInteger.TryParse(tag.Max, out BigInteger maxVal))
            {
                if (bigIntVal > maxVal)
                {
                    result.AddError($"big.Int length {bigIntVal} > max {maxVal}");
                    return;
                }
            }
            else
            {
                result.AddError($"invalid max \"{tag.Max}\" for big.Int");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type big.Int not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateDatetime(dynamic value, Tag tag, ValidationResult result)
    {
        if (!(value is DateTime dt))
        {
            result.AddError("value must be a datetime");
            return;
        }

        switch (tag.Type)
        {
            case ValueType.Datetime:
                ValidateDatetime(dt, tag, result);
                break;
            case ValueType.Date:
                ValidateDate(dt, tag, result);
                break;
            case ValueType.Time:
                ValidateTime(dt, tag, result);
                break;
            default:
                ValidateDatetime(dt, tag, result);
                break;
        }
    }

    private void ValidateDatetime(DateTime val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        val = TruncateToSeconds(val);

        if (TimeUtil.EpochSeconds(val) == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }
    }

    private void ValidateDate(DateTime val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        val = TruncateToSeconds(val);

        if (TimeUtil.EpochSeconds(val) == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }
    }

    private void ValidateTime(DateTime val, Tag tag, ValidationResult result)
    {
        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        val = TruncateToSeconds(val);

        if (TimeUtil.EpochSeconds(val) == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }
    }

    private static DateTime TruncateToSeconds(DateTime dt)
    {
        long ticks = dt.Ticks - (dt.Ticks % TimeSpan.TicksPerSecond);
        return new DateTime(ticks, dt.Kind);
    }

    public void ValidateUuid(dynamic value, Tag tag, ValidationResult result)
    {
        if (!(value is string uuidStr))
        {
            result.AddError("value must be a string");
            return;
        }

        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (string.IsNullOrEmpty(uuidStr))
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        if (!UuidRegex.IsMatch(uuidStr))
        {
            result.AddError($"value '{uuidStr}' does not match UUID pattern");
            return;
        }

        if (!Guid.TryParse(uuidStr, out var guid))
        {
            result.AddError($"invalid uuid: {uuidStr}");
            return;
        }

        if (tag.Version != 0)
        {
            var bytes = guid.ToByteArray();
            int uuidVersion = (bytes[7] >> 4) & 0x0F;
            if (tag.Version != uuidVersion)
            {
                result.AddError("invalid uuid version");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type uuid not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateDecimal(dynamic value, Tag tag, ValidationResult result)
    {
        if (!(value is string decimalStr))
        {
            result.AddError("value must be a string");
            return;
        }

        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (string.IsNullOrEmpty(decimalStr))
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        if (!DecimalRegex.IsMatch(decimalStr))
        {
            result.AddError($"invalid decimal \"{decimalStr}\", must be like \"0.0\"");
            return;
        }

        if (tag.Location != 0)
        {
            result.AddError($"type decimal not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateIp(dynamic value, Tag tag, ValidationResult result)
    {
        if (!(value is IPAddress ipAddr))
        {
            if (value is string ipStr)
            {
                if (!IPAddress.TryParse(ipStr, out var parsed))
                {
                    result.AddError("value must be a valid IP address");
                    return;
                }
                ipAddr = parsed!;
            }
            else
            {
                result.AddError("value must be a IP address");
                return;
            }
        }

        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (ipAddr == null || ipAddr.ToString() == "0.0.0.0" || ipAddr.ToString() == "::")
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        if (tag.Version == 4)
        {
            if (ipAddr.AddressFamily != AddressFamily.InterNetwork)
            {
                result.AddError($"invalid ipv4: {ipAddr}");
                return;
            }
        }

        if (tag.Version == 6)
        {
            if (ipAddr.AddressFamily == AddressFamily.InterNetwork)
            {
                result.AddError($"invalid ipv6: {ipAddr}");
                return;
            }
        }

        if (tag.Location != 0)
        {
            result.AddError($"type ip not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateUrl(dynamic value, Tag tag, ValidationResult result)
    {
        if (!(value is Uri uri))
        {
            if (value is string uriStr)
            {
                if (!Uri.TryCreate(uriStr, UriKind.Absolute, out var parsed))
                {
                    result.AddError("value must be a valid URL");
                    return;
                }
                uri = parsed!;
            }
            else
            {
                result.AddError("value must be a URL");
                return;
            }
        }

        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (uri.ToString() == "")
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        if (uri.Scheme != "http" && uri.Scheme != "https")
        {
            result.AddError($"invalid url: {uri}");
            return;
        }

        if (string.IsNullOrEmpty(uri.Host))
        {
            result.AddError($"invalid url: {uri}");
            return;
        }

        if (tag.Location != 0)
        {
            result.AddError($"type url not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateEmail(dynamic value, Tag tag, ValidationResult result)
    {
        if (!(value is string emailStr))
        {
            result.AddError("value must be a string");
            return;
        }

        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (string.IsNullOrEmpty(emailStr))
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        if (!EmailRegex.IsMatch(emailStr))
        {
            result.AddError($"value '{emailStr}' does not match email pattern");
            return;
        }

        if (tag.Location != 0)
        {
            result.AddError($"type email not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateEnum(dynamic value, Tag tag, ValidationResult result)
    {
        if (!(value is string enumStr))
        {
            result.AddError("value must be a string");
            return;
        }

        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        if (string.IsNullOrEmpty(enumStr))
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        var enums = tag.Enums.Split('|');
        int idx = -1;
        for (int i = 0; i < enums.Length; i++)
        {
            if (enums[i].Trim() == enumStr)
            {
                idx = i;
                break;
            }
        }

        if (idx == -1)
        {
            result.AddError($"value '{enumStr}' not found in enums: {string.Join(", ", enums)}");
            return;
        }

        if (tag.Location != 0)
        {
            result.AddError($"type enum not support location UTC{tag.Location}");
            return;
        }
    }

    public void ValidateMedia(dynamic value, Tag tag, ValidationResult result)
    {
        if (!(value is byte[] mediaBytes))
        {
            result.AddError("value must be a byte array");
            return;
        }

        if (tag.Desc.Length > 65535)
        {
            result.AddError("desc length exceeds 65535 bytes");
            return;
        }

        int l = mediaBytes.Length;

        if (l == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("not allow empty (add 'allow_empty' tag if empty is allowed)");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (int.TryParse(tag.Min, out int minVal))
            {
                if (l < minVal)
                {
                    result.AddError($"[]byte length {l} < min {minVal}");
                    return;
                }
            }
            else
            {
                return;
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (int.TryParse(tag.Max, out int maxVal))
            {
                if (l > maxVal)
                {
                    result.AddError($"[]byte length {l} > max {maxVal}");
                    return;
                }
            }
            else
            {
                return;
            }
        }

        if (tag.Size != 0 && l != tag.Size)
        {
            result.AddError($"[]byte length {l} != size {tag.Size}");
            return;
        }

        if (tag.Location != 0)
        {
            result.AddError($"type video not support location UTC{tag.Location}");
            return;
        }
    }

    private bool TryGetInt64Value(dynamic value, out long result)
    {
        result = 0;
        if (value is long)
        {
            result = value;
            return true;
        }
        else if (value is int)
        {
            result = value;
            return true;
        }
        else if (value is short)
        {
            result = value;
            return true;
        }
        else if (value is sbyte)
        {
            result = value;
            return true;
        }
        else if (value is ulong ulVal)
        {
            result = (long)ulVal;
            return true;
        }
        else if (value is uint)
        {
            result = value;
            return true;
        }
        else if (value is ushort)
        {
            result = value;
            return true;
        }
        else if (value is byte)
        {
            result = value;
            return true;
        }
        return false;
    }

    private bool TryGetUInt64Value(dynamic value, out ulong result)
    {
        result = 0;
        if (value is ulong)
        {
            result = value;
            return true;
        }
        else if (value is uint)
        {
            result = value;
            return true;
        }
        else if (value is ushort)
        {
            result = value;
            return true;
        }
        else if (value is byte)
        {
            result = value;
            return true;
        }
        else if (value is long lVal)
        {
            result = (ulong)lVal;
            return true;
        }
        else if (value is int)
        {
            result = (ulong)value;
            return true;
        }
        else if (value is short)
        {
            result = (ulong)value;
            return true;
        }
        else if (value is sbyte)
        {
            result = (ulong)value;
            return true;
        }
        return false;
    }
}

public static class Validator
{
    private static readonly MmValidator _validator = new MmValidator();

    public static ValidationResult Validate(dynamic value, Tag tag)
    {
        return _validator.Validate(value, tag);
    }
}
