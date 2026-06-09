namespace MetaMessage;

/// <summary>
/// Convenience API for MetaMessage encoding/decoding.
/// Provides quick access to common operations without needing to import the Core namespace.
/// </summary>
public static class Mm
{
    /// <summary>
    /// Convert a value directly to MetaMessage binary format.
    /// </summary>
    /// <param name="value">The value to encode.</param>
    /// <param name="tagString">Optional mm tag string (e.g. "name=foo; type=i").</param>
    /// <returns>Byte array of the encoded data.</returns>
    public static byte[] EncodeFromValue(object value, string tagString = "")
    {
        return Core.MetaMessage.FromValue(value, tagString);
    }

    /// <summary>
    /// Convert a JSONC string to MetaMessage binary format.
    /// </summary>
    public static byte[] EncodeFromJsonc(string jsonc)
    {
        return Core.MetaMessage.FromJSONC(jsonc);
    }

    /// <summary>
    /// Decode MetaMessage binary format to a plain value (Dictionary, List, or scalar).
    /// </summary>
    public static object? DecodeToValue(byte[] data)
    {
        return Core.MetaMessage.DecodeToValue(data);
    }

    /// <summary>
    /// Decode MetaMessage binary format to a JSONC string.
    /// </summary>
    public static string DecodeToJsonc(byte[] data)
    {
        return Core.MetaMessage.DecodeToJsonc(data);
    }

    /// <summary>
    /// Convert a value to a JSONC string.
    /// </summary>
    /// <param name="value">The value to convert.</param>
    /// <param name="tagString">Optional mm tag string.</param>
    /// <returns>The JSONC string.</returns>
    public static string ValueToJsonc(object value, string tagString = "")
    {
        if (string.IsNullOrEmpty(tagString))
        {
            return Core.MetaMessage.ValueToJsonc(value);
        }
        return Core.MetaMessage.ValueToJsonc(value, tagString);
    }

    /// <summary>
    /// Convert a JSONC string to a plain value.
    /// </summary>
    public static object? JsoncToValue(string jsonc)
    {
        return Core.MetaMessage.JsoncToValue(jsonc);
    }
}