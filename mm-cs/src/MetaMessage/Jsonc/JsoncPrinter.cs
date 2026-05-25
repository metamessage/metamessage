using MetaMessage.Ir;
using ValueType = MetaMessage.Ir.ValueType;

namespace MetaMessage.Jsonc;

public class JsoncPrinter
{
    private readonly bool _prettyPrint;
    private const string IndentUnit = "\t";

    public JsoncPrinter(bool prettyPrint = true)
    {
        _prettyPrint = prettyPrint;
    }

    public string Print(IJsoncNode node)
    {
        var sb = new System.Text.StringBuilder();
        WriteLeadingComments(sb, node.Tag, 0);
        PrintNode(node, sb, 0);
        return sb.ToString();
    }

    private void PrintNode(IJsoncNode node, System.Text.StringBuilder sb, int indent)
    {
        if (node is JsoncDoc doc)
        {
            WriteObjectJSONC(sb, doc.Fields.Select(f => new KeyValuePair<string, IJsoncNode>(f.Key, f.Value)).ToList(), doc.Tag, indent);
        }
        else if (node is JsoncObject obj)
        {
            WriteObjectJSONC(sb, obj.Fields.ToList(), obj.Tag, indent);
        }
        else if (node is JsoncArray array)
        {
            WriteArrayJSONC(sb, array.Elements, array.Tag, indent);
        }
        else if (node is JsoncValue value)
        {
            WriteValueJSONC(sb, value, indent);
        }
    }

    private void WriteLeadingComments(System.Text.StringBuilder b, Tag? tag, int indent)
    {
        if (tag == null) return;
        var tagStr = tag.ToString();
        if (string.IsNullOrEmpty(tagStr)) return;
        b.Append('\n');
        WriteIndent(b, indent);
        b.Append(tagStr);
        b.Append('\n');
    }

    private void WriteIndent(System.Text.StringBuilder b, int indent)
    {
        for (int i = 0; i < indent; i++)
        {
            b.Append(IndentUnit);
        }
    }

    private void WriteValueJSONC(System.Text.StringBuilder b, JsoncValue v, int indent)
    {
        if (v.Tag != null)
        {
            switch (v.Tag.Type)
            {
                case ValueType.Str:
                case ValueType.Bytes:
                case ValueType.Datetime:
                case ValueType.Date:
                case ValueType.Time:
                case ValueType.Uuid:
                case ValueType.Ip:
                case ValueType.Url:
                case ValueType.Email:
                case ValueType.Enums:
                    b.Append('"');
                    b.Append(EscapeString(v.Value?.ToString() ?? ""));
                    b.Append('"');
                    return;

                case ValueType.I:
                case ValueType.I8:
                case ValueType.I16:
                case ValueType.I32:
                case ValueType.I64:
                case ValueType.U:
                case ValueType.U8:
                case ValueType.U16:
                case ValueType.U32:
                case ValueType.U64:
                case ValueType.Bigint:
                case ValueType.Decimal:
                case ValueType.Bool:
                    if (v.Value is bool bVal)
                    {
                        b.Append(bVal ? "true" : "false");
                    }
                    else
                    {
                        b.Append(v.Value?.ToString() ?? "null");
                    }
                    return;

                case ValueType.F32:
                case ValueType.F64:
                    b.Append(v.Value?.ToString() ?? "null");
                    return;

                default:
                    b.Append(v.Value?.ToString() ?? "null");
                    return;
            }
        }

        switch (v.TokenType)
        {
            case JsoncTokenType.String:
                b.Append('"');
                b.Append(EscapeString(v.Value?.ToString() ?? ""));
                b.Append('"');
                break;
            case JsoncTokenType.Number:
                b.Append(v.Value?.ToString() ?? "null");
                break;
            case JsoncTokenType.True:
                b.Append("true");
                break;
            case JsoncTokenType.False:
                b.Append("false");
                break;
            case JsoncTokenType.Null:
                b.Append("null");
                break;
            default:
                b.Append(v.Value?.ToString() ?? "null");
                break;
        }
    }

    private void WriteObjectJSONC(System.Text.StringBuilder b, List<KeyValuePair<string, IJsoncNode>> fields, Tag? tag, int indent)
    {
        WriteLeadingComments(b, tag, indent);
        b.Append('{');
        if (!_prettyPrint)
        {
            foreach (var field in fields)
            {
                b.Append('"');
                b.Append(field.Key);
                b.Append("\":");
                PrintNode(field.Value, b, indent);
                b.Append(',');
            }
            b.Append('}');
            return;
        }

        b.Append('\n');
        foreach (var field in fields)
        {
            WriteLeadingComments(b, field.Value.Tag, indent + 1);
            WriteIndent(b, indent + 1);
            b.Append('"');
            b.Append(field.Key);
            b.Append("\": ");
            PrintNode(field.Value, b, indent + 1);
            b.Append(",\n");
        }
        WriteIndent(b, indent);
        b.Append('}');
    }

    private void WriteArrayJSONC(System.Text.StringBuilder b, List<IJsoncNode> items, Tag? tag, int indent)
    {
        WriteLeadingComments(b, tag, indent);
        b.Append('[');
        if (!_prettyPrint)
        {
            foreach (var item in items)
            {
                PrintNode(item, b, indent);
                b.Append(',');
            }
            b.Append(']');
            return;
        }

        b.Append('\n');
        foreach (var item in items)
        {
            WriteLeadingComments(b, item.Tag, indent + 1);
            WriteIndent(b, indent + 1);
            PrintNode(item, b, indent + 1);
            b.Append(",\n");
        }
        WriteIndent(b, indent);
        b.Append(']');
    }

    private static string EscapeString(string s)
    {
        var sb = new System.Text.StringBuilder();
        foreach (char c in s)
        {
            switch (c)
            {
                case '"':
                    sb.Append("\\\"");
                    break;
                case '\\':
                    sb.Append("\\\\");
                    break;
                case '\n':
                    sb.Append("\\n");
                    break;
                case '\r':
                    sb.Append("\\r");
                    break;
                case '\t':
                    sb.Append("\\t");
                    break;
                default:
                    sb.Append(c);
                    break;
            }
        }
        return sb.ToString();
    }
}