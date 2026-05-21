using System;
using System.Collections.Generic;
using System.Numerics;
using System.Text.RegularExpressions;

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
    public ValidationResult Validate(dynamic value, MmTag tag)
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

        switch (tag.Type)
        {
            case ValueType.BOOL:
                ValidateBool(value, tag, result);
                break;
            case ValueType.INT:
                ValidateI(value, tag, result);
                break;
            case ValueType.INT8:
                ValidateI8(value, tag, result);
                break;
            case ValueType.INT16:
                ValidateI16(value, tag, result);
                break;
            case ValueType.INT32:
                ValidateI32(value, tag, result);
                break;
            case ValueType.INT64:
                ValidateI64(value, tag, result);
                break;
            case ValueType.UINT:
                ValidateU(value, tag, result);
                break;
            case ValueType.UINT8:
                ValidateU8(value, tag, result);
                break;
            case ValueType.UINT16:
                ValidateU16(value, tag, result);
                break;
            case ValueType.UINT32:
                ValidateU32(value, tag, result);
                break;
            case ValueType.UINT64:
                ValidateU64(value, tag, result);
                break;
            case ValueType.BIGINT:
                ValidateBiginteger(value, tag, result);
                break;
            case ValueType.FLOAT32:
                ValidateF32(value, tag, result);
                break;
            case ValueType.FLOAT64:
                ValidateF64(value, tag, result);
                break;
            case ValueType.DECIMAL:
                ValidateDecimal(value, tag, result);
                break;
            case ValueType.STRING:
                ValidateStr(value, tag, result);
                break;
            case ValueType.EMAIL:
                ValidateEmail(value, tag, result);
                break;
            case ValueType.URL:
                ValidateUrl(value, tag, result);
                break;
            case ValueType.BYTES:
                ValidateBytes(value, tag, result);
                break;
            case ValueType.UUID:
                ValidateUuid(value, tag, result);
                break;
            case ValueType.IP:
                ValidateIp(value, tag, result);
                break;
            case ValueType.IMAGE:
                ValidateImage(value, tag, result);
                break;
            case ValueType.VIDEO:
                ValidateVideo(value, tag, result);
                break;
            case ValueType.DATETIME:
            case ValueType.DATE:
            case ValueType.TIME:
                ValidateDatetime(value, tag, result);
                break;
            case ValueType.ENUM:
                ValidateEnum(value, tag, result);
                break;
            case ValueType.ARRAY:
            case ValueType.VEC:
                ValidateArr(value, tag, result);
                break;
            case ValueType.OBJ:
                ValidateObj(value, tag, result);
                break;
        }

        return result;
    }

    private void ValidateBool(dynamic value, MmTag tag, ValidationResult result)
    {
        if (!(value is bool))
        {
            result.AddError("value must be a boolean");
            return;
        }

        if (tag.AllowEmpty)
        {
            result.AddError("type bool not support allow empty");
            return;
        }
    }

    private void ValidateI(dynamic value, MmTag tag, ValidationResult result)
    {
        long longValue;
        if (!TryGetInt64Value(value, out longValue))
        {
            result.AddError("value must be a number");
            return;
        }

        if (longValue == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type int not allow empty value 0");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (long.TryParse(tag.Min, out long minVal))
            {
                if (longValue < minVal)
                {
                    result.AddError($"value {longValue} is less than the minimum limit {minVal}");
                }
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (long.TryParse(tag.Max, out long maxVal))
            {
                if (longValue > maxVal)
                {
                    result.AddError($"value {longValue} exceeds the maximum limit {maxVal}");
                }
            }
        }
    }

    private void ValidateI8(dynamic value, MmTag tag, ValidationResult result)
    {
        long longValue;
        if (!TryGetInt64Value(value, out longValue))
        {
            result.AddError("value must be a number");
            return;
        }

        if (longValue == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type int8 not allow empty value 0");
            return;
        }

        if (longValue < sbyte.MinValue || longValue > sbyte.MaxValue)
        {
            result.AddError($"value {longValue} out of range for int8 ({sbyte.MinValue} to {sbyte.MaxValue})");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (long.TryParse(tag.Min, out long minVal))
            {
                if (minVal < sbyte.MinValue || minVal > sbyte.MaxValue)
                {
                    result.AddError($"tag.min {minVal} is out of int8 range [{sbyte.MinValue}, {sbyte.MaxValue}]");
                }
                else if (longValue < minVal)
                {
                    result.AddError($"value {longValue} is less than the minimum limit {minVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as int8: {tag.Min}");
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (long.TryParse(tag.Max, out long maxVal))
            {
                if (maxVal < sbyte.MinValue || maxVal > sbyte.MaxValue)
                {
                    result.AddError($"tag.max {maxVal} is out of int8 range [{sbyte.MinValue}, {sbyte.MaxValue}]");
                }
                else if (longValue > maxVal)
                {
                    result.AddError($"value {longValue} exceeds the maximum limit {maxVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as int8: {tag.Max}");
            }
        }
    }

    private void ValidateI16(dynamic value, MmTag tag, ValidationResult result)
    {
        long longValue;
        if (!TryGetInt64Value(value, out longValue))
        {
            result.AddError("value must be a number");
            return;
        }

        if (longValue == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type int16 not allow empty value 0");
            return;
        }

        if (longValue < short.MinValue || longValue > short.MaxValue)
        {
            result.AddError($"value {longValue} out of range for int16 ({short.MinValue} to {short.MaxValue})");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (long.TryParse(tag.Min, out long minVal))
            {
                if (minVal < short.MinValue || minVal > short.MaxValue)
                {
                    result.AddError($"tag.min {minVal} is out of int16 range [{short.MinValue}, {short.MaxValue}]");
                }
                else if (longValue < minVal)
                {
                    result.AddError($"value {longValue} is less than the minimum limit {minVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as int16: {tag.Min}");
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (long.TryParse(tag.Max, out long maxVal))
            {
                if (maxVal < short.MinValue || maxVal > short.MaxValue)
                {
                    result.AddError($"tag.max {maxVal} is out of int16 range [{short.MinValue}, {short.MaxValue}]");
                }
                else if (longValue > maxVal)
                {
                    result.AddError($"value {longValue} exceeds the maximum limit {maxVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as int16: {tag.Max}");
            }
        }
    }

    private void ValidateI32(dynamic value, MmTag tag, ValidationResult result)
    {
        long longValue;
        if (!TryGetInt64Value(value, out longValue))
        {
            result.AddError("value must be a number");
            return;
        }

        if (longValue == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type int32 not allow empty value 0");
            return;
        }

        if (longValue < int.MinValue || longValue > int.MaxValue)
        {
            result.AddError($"value {longValue} out of range for int32 ({int.MinValue} to {int.MaxValue})");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (long.TryParse(tag.Min, out long minVal))
            {
                if (minVal < int.MinValue || minVal > int.MaxValue)
                {
                    result.AddError($"tag.min {minVal} is out of int32 range [{int.MinValue}, {int.MaxValue}]");
                }
                else if (longValue < minVal)
                {
                    result.AddError($"value {longValue} is less than the minimum limit {minVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as int32: {tag.Min}");
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (long.TryParse(tag.Max, out long maxVal))
            {
                if (maxVal < int.MinValue || maxVal > int.MaxValue)
                {
                    result.AddError($"tag.max {maxVal} is out of int32 range [{int.MinValue}, {int.MaxValue}]");
                }
                else if (longValue > maxVal)
                {
                    result.AddError($"value {longValue} exceeds the maximum limit {maxVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as int32: {tag.Max}");
            }
        }
    }

    private void ValidateI64(dynamic value, MmTag tag, ValidationResult result)
    {
        long longValue;
        if (!TryGetInt64Value(value, out longValue))
        {
            result.AddError("value must be a number");
            return;
        }

        if (longValue == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type int64 not allow empty value 0");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (long.TryParse(tag.Min, out long minVal))
            {
                if (minVal < long.MinValue || minVal > long.MaxValue)
                {
                    result.AddError($"tag.min {minVal} is out of int64 range [{long.MinValue}, {long.MaxValue}]");
                }
                else if (longValue < minVal)
                {
                    result.AddError($"value {longValue} is less than the minimum limit {minVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as int64: {tag.Min}");
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (long.TryParse(tag.Max, out long maxVal))
            {
                if (maxVal < long.MinValue || maxVal > long.MaxValue)
                {
                    result.AddError($"tag.max {maxVal} is out of int64 range [{long.MinValue}, {long.MaxValue}]");
                }
                else if (longValue > maxVal)
                {
                    result.AddError($"value {longValue} exceeds the maximum limit {maxVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as int64: {tag.Max}");
            }
        }
    }

    private void ValidateU(dynamic value, MmTag tag, ValidationResult result)
    {
        ulong ulongValue;
        if (!TryGetUInt64Value(value, out ulongValue))
        {
            result.AddError("value must be a number");
            return;
        }

        if (ulongValue == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type uint not allow empty value 0");
            return;
        }

        if (ulongValue < 0 || ulongValue > uint.MaxValue)
        {
            result.AddError($"value {ulongValue} out of range for uint (0 to {uint.MaxValue})");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (ulong.TryParse(tag.Min, out ulong minVal))
            {
                if (minVal < 0 || minVal > uint.MaxValue)
                {
                    result.AddError($"tag.min {minVal} is out of uint range [0, {uint.MaxValue}]");
                }
                else if (ulongValue < minVal)
                {
                    result.AddError($"value {ulongValue} is less than the minimum limit {minVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as uint: {tag.Min}");
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (ulong.TryParse(tag.Max, out ulong maxVal))
            {
                if (maxVal < 0 || maxVal > uint.MaxValue)
                {
                    result.AddError($"tag.max {maxVal} is out of uint range [0, {uint.MaxValue}]");
                }
                else if (ulongValue > maxVal)
                {
                    result.AddError($"value {ulongValue} exceeds the maximum limit {maxVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as uint: {tag.Max}");
            }
        }
    }

    private void ValidateU8(dynamic value, MmTag tag, ValidationResult result)
    {
        ulong ulongValue;
        if (!TryGetUInt64Value(value, out ulongValue))
        {
            result.AddError("value must be a number");
            return;
        }

        if (ulongValue == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type uint8 not allow empty value 0");
            return;
        }

        if (ulongValue > byte.MaxValue)
        {
            result.AddError($"value {ulongValue} out of range for uint8 (0 to {byte.MaxValue})");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (ulong.TryParse(tag.Min, out ulong minVal))
            {
                if (ulongValue < minVal)
                {
                    result.AddError($"value {ulongValue} is less than the minimum limit {minVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as uint8: {tag.Min}");
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (ulong.TryParse(tag.Max, out ulong maxVal))
            {
                if (maxVal > byte.MaxValue)
                {
                    result.AddError($"tag.max {maxVal} is out of uint8 range [0, {byte.MaxValue}]");
                }
                else if (ulongValue > maxVal)
                {
                    result.AddError($"value {ulongValue} exceeds the maximum limit {maxVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as uint8: {tag.Max}");
            }
        }
    }

    private void ValidateU16(dynamic value, MmTag tag, ValidationResult result)
    {
        ulong ulongValue;
        if (!TryGetUInt64Value(value, out ulongValue))
        {
            result.AddError("value must be a number");
            return;
        }

        if (ulongValue == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type uint16 not allow empty value 0");
            return;
        }

        if (ulongValue < ushort.MinValue || ulongValue > ushort.MaxValue)
        {
            result.AddError($"value {ulongValue} out of range for uint16 ({ushort.MinValue} to {ushort.MaxValue})");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (ulong.TryParse(tag.Min, out ulong minVal))
            {
                if (minVal < ushort.MinValue || minVal > ushort.MaxValue)
                {
                    result.AddError($"tag.min {minVal} is out of uint16 range [{ushort.MinValue}, {ushort.MaxValue}]");
                }
                else if (ulongValue < minVal)
                {
                    result.AddError($"value {ulongValue} is less than the minimum limit {minVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as uint16: {tag.Min}");
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (ulong.TryParse(tag.Max, out ulong maxVal))
            {
                if (maxVal < ushort.MinValue || maxVal > ushort.MaxValue)
                {
                    result.AddError($"tag.max {maxVal} is out of uint16 range [{ushort.MinValue}, {ushort.MaxValue}]");
                }
                else if (ulongValue > maxVal)
                {
                    result.AddError($"value {ulongValue} exceeds the maximum limit {maxVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as uint16: {tag.Max}");
            }
        }
    }

    private void ValidateU32(dynamic value, MmTag tag, ValidationResult result)
    {
        ulong ulongValue;
        if (!TryGetUInt64Value(value, out ulongValue))
        {
            result.AddError("value must be a number");
            return;
        }

        if (ulongValue == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type uint32 not allow empty value 0");
            return;
        }

        if (ulongValue < uint.MinValue || ulongValue > uint.MaxValue)
        {
            result.AddError($"value {ulongValue} out of range for uint32 ({uint.MinValue} to {uint.MaxValue})");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (ulong.TryParse(tag.Min, out ulong minVal))
            {
                if (minVal < uint.MinValue || minVal > uint.MaxValue)
                {
                    result.AddError($"tag.min {minVal} is out of uint32 range [{uint.MinValue}, {uint.MaxValue}]");
                }
                else if (ulongValue < minVal)
                {
                    result.AddError($"value {ulongValue} is less than the minimum limit {minVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as uint32: {tag.Min}");
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (ulong.TryParse(tag.Max, out ulong maxVal))
            {
                if (maxVal < uint.MinValue || maxVal > uint.MaxValue)
                {
                    result.AddError($"tag.max {maxVal} is out of uint32 range [{uint.MinValue}, {uint.MaxValue}]");
                }
                else if (ulongValue > maxVal)
                {
                    result.AddError($"value {ulongValue} exceeds the maximum limit {maxVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as uint32: {tag.Max}");
            }
        }
    }

    private void ValidateU64(dynamic value, MmTag tag, ValidationResult result)
    {
        ulong ulongValue;
        if (!TryGetUInt64Value(value, out ulongValue))
        {
            result.AddError("value must be a number");
            return;
        }

        if (ulongValue == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type uint64 not allow empty value 0");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (ulong.TryParse(tag.Min, out ulong minVal))
            {
                if (minVal < ulong.MinValue || minVal > ulong.MaxValue)
                {
                    result.AddError($"tag.min {minVal} is out of uint64 range [{ulong.MinValue}, {ulong.MaxValue}]");
                }
                else if (ulongValue < minVal)
                {
                    result.AddError($"value {ulongValue} is less than the minimum limit {minVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as uint64: {tag.Min}");
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (ulong.TryParse(tag.Max, out ulong maxVal))
            {
                if (maxVal < ulong.MinValue || maxVal > ulong.MaxValue)
                {
                    result.AddError($"tag.max {maxVal} is out of uint64 range [{ulong.MinValue}, {ulong.MaxValue}]");
                }
                else if (ulongValue > maxVal)
                {
                    result.AddError($"value {ulongValue} exceeds the maximum limit {maxVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as uint64: {tag.Max}");
            }
        }
    }

    private void ValidateBiginteger(dynamic value, MmTag tag, ValidationResult result)
    {
        BigInteger bigIntValue;
        if (value is BigInteger)
        {
            bigIntValue = value;
        }
        else if (value is long)
        {
            bigIntValue = value;
        }
        else if (value is int)
        {
            bigIntValue = value;
        }
        else if (value is short)
        {
            bigIntValue = value;
        }
        else if (value is byte)
        {
            bigIntValue = value;
        }
        else if (value is sbyte)
        {
            bigIntValue = value;
        }
        else if (value is ulong)
        {
            bigIntValue = value;
        }
        else if (value is uint)
        {
            bigIntValue = value;
        }
        else if (value is ushort)
        {
            bigIntValue = value;
        }
        else if (value is string strVal && BigInteger.TryParse(strVal, out bigIntValue))
        {
        }
        else
        {
            result.AddError("value must be a BigInteger");
            return;
        }

        if (bigIntValue == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type big.Int not allow empty value 0");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (BigInteger.TryParse(tag.Min, out BigInteger minVal))
            {
                if (bigIntValue < minVal)
                {
                    result.AddError($"big.Int value {bigIntValue} is less than the minimum limit {minVal}");
                }
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (BigInteger.TryParse(tag.Max, out BigInteger maxVal))
            {
                if (bigIntValue > maxVal)
                {
                    result.AddError($"big.Int value {bigIntValue} exceeds the maximum limit {maxVal}");
                }
            }
        }
    }

    private void ValidateF32(dynamic value, MmTag tag, ValidationResult result)
    {
        double doubleValue;
        if (value is double)
        {
            doubleValue = value;
        }
        else if (value is float)
        {
            doubleValue = value;
        }
        else if (value is decimal)
        {
            doubleValue = (double)value;
        }
        else if (value is int || value is long || value is short || value is byte || value is sbyte)
        {
            doubleValue = Convert.ToDouble(value);
        }
        else
        {
            result.AddError("value must be a float");
            return;
        }

        if (doubleValue == 0.0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type float32 not allow empty value 0.0");
            return;
        }

        if (doubleValue < float.MinValue || doubleValue > float.MaxValue)
        {
            result.AddError($"value {doubleValue} out of range for float32 ({float.MinValue} to {float.MaxValue})");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (double.TryParse(tag.Min, out double minVal))
            {
                if (minVal < float.MinValue || minVal > float.MaxValue)
                {
                    result.AddError($"tag.min {minVal} is out of float32 range [{float.MinValue}, {float.MaxValue}]");
                }
                else if (doubleValue < minVal)
                {
                    result.AddError($"value {doubleValue} is less than the minimum limit {minVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as float32: {tag.Min}");
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (double.TryParse(tag.Max, out double maxVal))
            {
                if (maxVal < float.MinValue || maxVal > float.MaxValue)
                {
                    result.AddError($"tag.max {maxVal} is out of float32 range [{float.MinValue}, {float.MaxValue}]");
                }
                else if (doubleValue > maxVal)
                {
                    result.AddError($"value {doubleValue} exceeds the maximum limit {maxVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as float32: {tag.Max}");
            }
        }
    }

    private void ValidateF64(dynamic value, MmTag tag, ValidationResult result)
    {
        double doubleValue;
        if (value is double)
        {
            doubleValue = value;
        }
        else if (value is float)
        {
            doubleValue = value;
        }
        else if (value is decimal)
        {
            doubleValue = (double)value;
        }
        else if (value is int || value is long || value is short || value is byte || value is sbyte)
        {
            doubleValue = Convert.ToDouble(value);
        }
        else
        {
            result.AddError("value must be a float");
            return;
        }

        if (doubleValue == 0.0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type float64 not allow empty value 0.0");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (double.TryParse(tag.Min, out double minVal))
            {
                if (minVal < double.MinValue || minVal > double.MaxValue)
                {
                    result.AddError($"tag.min {minVal} is out of float64 range [{double.MinValue}, {double.MaxValue}]");
                }
                else if (doubleValue < minVal)
                {
                    result.AddError($"value {doubleValue} is less than the minimum limit {minVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.min as float64: {tag.Min}");
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (double.TryParse(tag.Max, out double maxVal))
            {
                if (maxVal < double.MinValue || maxVal > double.MaxValue)
                {
                    result.AddError($"tag.max {maxVal} is out of float64 range [{double.MinValue}, {double.MaxValue}]");
                }
                else if (doubleValue > maxVal)
                {
                    result.AddError($"value {doubleValue} exceeds the maximum limit {maxVal}");
                }
            }
            else
            {
                result.AddError($"failed to parse tag.max as float64: {tag.Max}");
            }
        }
    }

    private void ValidateDecimal(dynamic value, MmTag tag, ValidationResult result)
    {
        if (value is decimal decimalVal)
        {
            if (decimalVal == 0)
            {
                if (tag.AllowEmpty)
                {
                    return;
                }
                result.AddError("type decimal not allow empty value 0");
                return;
            }

            if (!string.IsNullOrEmpty(tag.Min))
            {
                if (decimal.TryParse(tag.Min, out decimal minVal))
                {
                    if (minVal < decimal.MinValue || minVal > decimal.MaxValue)
                    {
                        result.AddError($"tag.min {minVal} is out of decimal range [{decimal.MinValue}, {decimal.MaxValue}]");
                    }
                    else if (decimalVal < minVal)
                    {
                        result.AddError($"value {decimalVal} is less than the minimum limit {minVal}");
                    }
                }
                else
                {
                    result.AddError($"failed to parse tag.min as decimal: {tag.Min}");
                }
            }

            if (!string.IsNullOrEmpty(tag.Max))
            {
                if (decimal.TryParse(tag.Max, out decimal maxVal))
                {
                    if (maxVal < decimal.MinValue || maxVal > decimal.MaxValue)
                    {
                        result.AddError($"tag.max {maxVal} is out of decimal range [{decimal.MinValue}, {decimal.MaxValue}]");
                    }
                    else if (decimalVal > maxVal)
                    {
                        result.AddError($"value {decimalVal} exceeds the maximum limit {maxVal}");
                    }
                }
                else
                {
                    result.AddError($"failed to parse tag.max as decimal: {tag.Max}");
                }
            }
        }
        else if (value is double doubleVal)
        {
            ValidateF64(doubleVal, tag, result);
        }
        else if (value is float floatVal)
        {
            ValidateF32(floatVal, tag, result);
        }
        else
        {
            result.AddError("value must be a decimal");
        }
    }

    private void ValidateStr(dynamic value, MmTag tag, ValidationResult result)
    {
        if (!(value is string))
        {
            result.AddError("value must be a string");
            return;
        }

        string strValue = value;

        if (string.IsNullOrEmpty(strValue) && !tag.AllowEmpty)
        {
            result.AddError("value is empty");
            return;
        }

        // Pattern validation
        if (!string.IsNullOrEmpty(tag.Pattern) && !string.IsNullOrEmpty(strValue))
        {
            try
            {
                if (!Regex.IsMatch(strValue, tag.Pattern))
                {
                    result.AddError($"value \"{strValue}\" does not match pattern {tag.Pattern}");
                }
            }
            catch (Exception)
            {
                result.AddError($"invalid pattern {tag.Pattern}");
            }
        }
    }

    private void ValidateEmail(dynamic value, MmTag tag, ValidationResult result)
    {
        if (!(value is string))
        {
            result.AddError("value must be a string");
            return;
        }

        string email = value;

        if (string.IsNullOrEmpty(email) && !tag.AllowEmpty)
        {
            result.AddError("value is empty");
            return;
        }

        if (!string.IsNullOrEmpty(email))
        {
            string emailRegex = @"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$";
            if (!Regex.IsMatch(email, emailRegex))
            {
                result.AddError("value is not a valid email");
            }
        }
    }

    private void ValidateUrl(dynamic value, MmTag tag, ValidationResult result)
    {
        if (!(value is string))
        {
            result.AddError("value must be a string");
            return;
        }

        string url = value;

        if (string.IsNullOrEmpty(url) && !tag.AllowEmpty)
        {
            result.AddError("value is empty");
            return;
        }

        if (!string.IsNullOrEmpty(url))
        {
            if (!Uri.IsWellFormedUriString(url, UriKind.Absolute))
            {
                result.AddError("value is not a valid url");
            }
        }
    }

    private void ValidateBytes(dynamic value, MmTag tag, ValidationResult result)
    {
        if (!(value is byte[]))
        {
            result.AddError("value must be a byte array");
            return;
        }

        byte[] bytes = value;

        if (bytes.Length == 0 && !tag.AllowEmpty)
        {
            result.AddError("value is empty");
        }
    }

    private void ValidateUuid(dynamic value, MmTag tag, ValidationResult result)
    {
        if (!(value is string))
        {
            result.AddError("value must be a string");
            return;
        }

        string uuid = value;

        if (string.IsNullOrEmpty(uuid) && !tag.AllowEmpty)
        {
            result.AddError("value is empty");
            return;
        }

        if (!string.IsNullOrEmpty(uuid))
        {
            if (!Guid.TryParse(uuid, out var guid))
            {
                result.AddError("value is not a valid uuid");
                return;
            }

            // Version check
            if (tag.Version != 0)
            {
                var bytes = guid.ToByteArray();
                int uuidVersion = (bytes[7] >> 4) & 0x0F;
                if (tag.Version != uuidVersion)
                {
                    result.AddError($"invalid uuid version: expected {tag.Version}, got {uuidVersion}");
                }
            }
        }
    }

    private void ValidateIp(dynamic value, MmTag tag, ValidationResult result)
    {
        if (!(value is string))
        {
            result.AddError("value must be a string");
            return;
        }

        string ip = value;

        if (string.IsNullOrEmpty(ip))
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type ip not allow empty value");
            return;
        }

        if (!System.Net.IPAddress.TryParse(ip, out var parsedIp))
        {
            result.AddError("value is not a valid IP address");
            return;
        }

        // IP version check using tag.Version (4 or 6)
        if (tag.Version != 0)
        {
            if (tag.Version == 4 && parsedIp.AddressFamily != System.Net.Sockets.AddressFamily.InterNetwork)
            {
                result.AddError("expected IPv4 address");
            }
            else if (tag.Version == 6 && parsedIp.AddressFamily != System.Net.Sockets.AddressFamily.InterNetworkV6)
            {
                result.AddError("expected IPv6 address");
            }
        }
    }

    private void ValidateImage(dynamic value, MmTag tag, ValidationResult result)
    {
        byte[] bytes = value is byte[] b ? b : null;
        if (bytes == null)
        {
            result.AddError("value must be a byte array");
            return;
        }

        if (bytes.Length == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type image not allow empty value");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (int.TryParse(tag.Min, out int minVal) && bytes.Length < minVal)
            {
                result.AddError($"image byte length {bytes.Length} < min {minVal}");
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (int.TryParse(tag.Max, out int maxVal) && bytes.Length > maxVal)
            {
                result.AddError($"image byte length {bytes.Length} > max {maxVal}");
            }
        }

        if (tag.Size != 0 && bytes.Length != tag.Size)
        {
            result.AddError($"image byte length {bytes.Length} != size {tag.Size}");
        }
    }

    private void ValidateVideo(dynamic value, MmTag tag, ValidationResult result)
    {
        byte[] bytes = value is byte[] b ? b : null;
        if (bytes == null)
        {
            result.AddError("value must be a byte array");
            return;
        }

        if (bytes.Length == 0)
        {
            if (tag.AllowEmpty)
            {
                return;
            }
            result.AddError("type video not allow empty value");
            return;
        }

        if (!string.IsNullOrEmpty(tag.Min))
        {
            if (int.TryParse(tag.Min, out int minVal) && bytes.Length < minVal)
            {
                result.AddError($"video byte length {bytes.Length} < min {minVal}");
            }
        }

        if (!string.IsNullOrEmpty(tag.Max))
        {
            if (int.TryParse(tag.Max, out int maxVal) && bytes.Length > maxVal)
            {
                result.AddError($"video byte length {bytes.Length} > max {maxVal}");
            }
        }

        if (tag.Size != 0 && bytes.Length != tag.Size)
        {
            result.AddError($"video byte length {bytes.Length} != size {tag.Size}");
        }
    }

    private void ValidateDatetime(dynamic value, MmTag tag, ValidationResult result)
    {
        if (!(value is DateTime))
        {
            result.AddError("value must be a datetime");
            return;
        }
    }

    private void ValidateEnum(dynamic value, MmTag tag, ValidationResult result)
    {
        if (!(value is string))
        {
            result.AddError("value must be a string");
            return;
        }

        string enumValue = value;

        if (string.IsNullOrEmpty(enumValue) && !tag.AllowEmpty)
        {
            result.AddError("value is empty");
            return;
        }

        if (!string.IsNullOrEmpty(enumValue) && !string.IsNullOrEmpty(tag.Enum))
        {
            var enumValues = tag.Enum.Split('|');
            if (!enumValues.Contains(enumValue))
            {
                result.AddError("value is not in enum");
            }
        }
    }

    private void ValidateArr(dynamic value, MmTag tag, ValidationResult result)
    {
        if (!(value is System.Collections.IEnumerable))
        {
            result.AddError("value must be an array");
            return;
        }

        System.Collections.IEnumerable array = value;
        int count = 0;
        foreach (var item in array)
        {
            count++;
            var childTag = new MmTag();
            childTag.InheritFromArrayParent(tag);
            var itemResult = Validate(item, childTag);
            if (!itemResult.IsValid)
            {
                foreach (var error in itemResult.Errors)
                {
                    result.AddError($"item {count - 1}: {error}");
                }
            }
        }

        if (count == 0 && !tag.AllowEmpty)
        {
            result.AddError("value is empty");
        }
    }

    private void ValidateObj(dynamic value, MmTag tag, ValidationResult result)
    {
        if (value == null)
        {
            result.AddError("value must be an object");
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

    public static ValidationResult Validate(dynamic value, MmTag tag)
    {
        return _validator.Validate(value, tag);
    }
}