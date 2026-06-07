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

    public string Print(INode node)
    {
        var sb = new System.Text.StringBuilder();
        WriteLeadingComments(sb, node.Tag, 0);
        PrintNode(node, sb, 0);
        return sb.ToString();
    }

    private void PrintNode(INode node, System.Text.StringBuilder sb, int indent)
    {
        if (node is NodeScalar scalar)
        {
            WriteScalarJSONC(sb, scalar, indent);
        }
        else if (node is NodeArray array)
        {
            WriteArrayJSONC(sb, array.Children, array.Tag, indent);
        }
        else if (node is NodeObject map)
        {
            WriteMapJSONC(sb, map.Entries, map.Tag, indent);
        }
    }

    private void WriteLeadingComments(System.Text.StringBuilder b, Tag? tag, int indent)
    {
        if (tag == null) return;
        var tagStr = tag.ToString();
        if (string.IsNullOrEmpty(tagStr)) return;
        b.Append('\n');
        WriteIndent(b, indent);
        b.Append("// mm: ");
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

    private void WriteScalarJSONC(System.Text.StringBuilder b, NodeScalar v, int indent)
    {
        if (v.Tag != null)
        {
            if (v.Tag.IsNull)
            {
                WriteNullDefault(b, v.Tag);
                return;
            }
            switch (v.Tag.Type)
            {
                case ValueType.Str:
                case ValueType.Bytes:
                case ValueType.Media:
                case ValueType.Image:
                case ValueType.Video:
                case ValueType.Datetime:
                case ValueType.Date:
                case ValueType.Time:
                case ValueType.Uuid:
                case ValueType.Enums:
                case ValueType.Ip:
                case ValueType.Url:
                case ValueType.Email:
                    b.Append('"');
                    b.Append(EscapeString(v.Text));
                    b.Append('"');
                    return;

                case ValueType.Bigint:
                case ValueType.Decimal:
                    b.Append(v.Text);
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
                    {
                        if (v.Data is bool boolVal)
                        {
                            b.Append(boolVal ? "true" : "false");
                        }
                        else
                        {
                            b.Append(v.Text);
                        }
                    }
                    return;

                case ValueType.Bool:
                    {
                        if (v.Data is bool boolVal)
                        {
                            b.Append(boolVal ? "true" : "false");
                        }
                        else
                        {
                            b.Append(v.Text);
                        }
                    }
                    return;

                case ValueType.F32:
                case ValueType.F64:
                    b.Append(v.Text);
                    return;

                default:
                    b.Append(v.Text);
                    return;
            }
        }

        // No tag — infer from data type
        if (v.Data == null)
        {
            b.Append("null");
        }
        else if (v.Data is string)
        {
            b.Append('"');
            b.Append(EscapeString(v.Text));
            b.Append('"');
        }
        else if (v.Data is bool bVal)
        {
            b.Append(bVal ? "true" : "false");
        }
        else
        {
            b.Append(v.Text);
        }
    }

    private void WriteNullDefault(System.Text.StringBuilder b, Tag tag)
    {
        switch (tag.Type)
        {
            case ValueType.Bool:
                b.Append("false");
                break;
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
            case ValueType.F32:
            case ValueType.F64:
            case ValueType.Decimal:
                b.Append('0');
                break;
            case ValueType.Str:
            case ValueType.Email:
            case ValueType.Url:
            case ValueType.Enums:
            case ValueType.Datetime:
            case ValueType.Date:
            case ValueType.Time:
            case ValueType.Uuid:
            case ValueType.Bytes:
            default:
                b.Append("\"\"");
                break;
        }
    }

    private void WriteMapJSONC(System.Text.StringBuilder b, List<KeyValuePair<NodeScalar, INode>> entries, Tag? tag, int indent)
    {
        b.Append('{');
        if (!_prettyPrint)
        {
            foreach (var entry in entries)
            {
                b.Append('"');
                b.Append(entry.Key.Text);
                b.Append("\":");
                PrintNode(entry.Value, b, indent);
                b.Append(',');
            }
            b.Append('}');
            return;
        }

        b.Append('\n');
        foreach (var entry in entries)
        {
            WriteLeadingComments(b, entry.Value.Tag, indent + 1);
            WriteIndent(b, indent + 1);
            b.Append('"');
            b.Append(entry.Key.Text);
            b.Append("\": ");
            PrintNode(entry.Value, b, indent + 1);
            b.Append(",\n");
        }
        WriteIndent(b, indent);
        b.Append('}');
    }

    private void WriteArrayJSONC(System.Text.StringBuilder b, List<INode> items, Tag? tag, int indent)
    {
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