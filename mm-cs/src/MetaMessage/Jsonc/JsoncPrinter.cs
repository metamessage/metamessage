namespace MetaMessage.Jsonc;

public class JsoncPrinter
{
    private readonly bool _prettyPrint;
    private readonly int _indentSize;
    private int _currentIndent;

    public JsoncPrinter(bool prettyPrint = true, int indentSize = 2)
    {
        _prettyPrint = prettyPrint;
        _indentSize = indentSize;
    }

    public string Print(IJsoncNode node)
    {
        var sb = new System.Text.StringBuilder();
        PrintNode(node, sb);
        return sb.ToString();
    }

    private void PrintNode(IJsoncNode node, System.Text.StringBuilder sb)
    {
        if (node is JsoncObject obj)
        {
            PrintObject(obj, sb);
        }
        else if (node is JsoncArray array)
        {
            PrintArray(array, sb);
        }
        else if (node is JsoncValue value)
        {
            PrintValue(value, sb);
        }
    }

    private void PrintObject(JsoncObject obj, System.Text.StringBuilder sb)
    {
        if (obj.LeadingComment != null)
        {
            PrintComment(obj.LeadingComment, sb);
        }

        sb.Append('{');

        var fields = obj.Fields;
        if (fields.Count == 0)
        {
            sb.Append('}');
            return;
        }

        if (_prettyPrint)
        {
            sb.AppendLine();
            _currentIndent += _indentSize;
        }

        int index = 0;
        foreach (var kvp in fields)
        {
            if (_prettyPrint)
            {
                AppendIndent(sb);
            }

            if (kvp.Value.LeadingComment != null)
            {
                PrintComment(kvp.Value.LeadingComment, sb);
                if (_prettyPrint)
                {
                    AppendIndent(sb);
                }
            }

            sb.Append('"');
            sb.Append(kvp.Key);
            sb.Append('"');
            sb.Append(':');

            if (_prettyPrint)
            {
                sb.Append(' ');
            }

            PrintNode(kvp.Value, sb);

            index++;
            if (index < fields.Count)
            {
                sb.Append(',');
            }

            if (_prettyPrint)
            {
                sb.AppendLine();
            }
        }

        if (_prettyPrint)
        {
            sb.AppendLine();
            _currentIndent -= _indentSize;
            AppendIndent(sb);
        }

        sb.Append('}');
    }

    private void PrintArray(JsoncArray array, System.Text.StringBuilder sb)
    {
        if (array.LeadingComment != null)
        {
            PrintComment(array.LeadingComment, sb);
        }

        sb.Append('[');

        if (array.Elements.Count == 0)
        {
            sb.Append(']');
            return;
        }

        if (_prettyPrint)
        {
            sb.AppendLine();
            _currentIndent += _indentSize;
        }

        for (int i = 0; i < array.Elements.Count; i++)
        {
            if (_prettyPrint)
            {
                AppendIndent(sb);
            }

            var element = array.Elements[i];
            if (element.LeadingComment != null)
            {
                PrintComment(element.LeadingComment, sb);
                if (_prettyPrint)
                {
                    AppendIndent(sb);
                }
            }

            PrintNode(element, sb);

            if (i < array.Elements.Count - 1)
            {
                sb.Append(',');
            }

            if (_prettyPrint)
            {
                sb.AppendLine();
            }
        }

        if (_prettyPrint)
        {
            sb.AppendLine();
            _currentIndent -= _indentSize;
            AppendIndent(sb);
        }

        sb.Append(']');
    }

    private void PrintValue(JsoncValue value, System.Text.StringBuilder sb)
    {
        switch (value.TokenType)
        {
            case JsoncTokenType.String:
                sb.Append('"');
                sb.Append(EscapeString(value.GetString() ?? ""));
                sb.Append('"');
                break;
            case JsoncTokenType.Number:
                sb.Append(value.Value?.ToString() ?? "null");
                break;
            case JsoncTokenType.True:
                sb.Append("true");
                break;
            case JsoncTokenType.False:
                sb.Append("false");
                break;
            case JsoncTokenType.Null:
                sb.Append("null");
                break;
            default:
                var type = value.Tag?.Type ?? ValueType.Unknown;
                bool needsQuotes = type.NeedsQuotes();
                var text = value.Value?.ToString() ?? "";
                if (needsQuotes)
                {
                    sb.Append('"');
                    sb.Append(EscapeString(text));
                    sb.Append('"');
                }
                else
                {
                    sb.Append(text);
                }
                break;
        }
    }

    private void PrintComment(JsoncComment comment, System.Text.StringBuilder sb)
    {
        if (_prettyPrint)
        {
            AppendIndent(sb);
        }

        if (comment.IsBlock)
        {
            sb.Append("/*");
            var commentText = comment.Text.TrimStart('/', '*').TrimEnd('*', '/');
            if (commentText.StartsWith("mm:"))
            {
                sb.Append(commentText);
            }
            else
            {
                sb.Append(commentText.TrimStart('/', '*').TrimEnd('*', '/'));
            }
            sb.AppendLine("*/");
        }
        else
        {
            var commentText = comment.Text;
            if (commentText.StartsWith("//"))
            {
                sb.AppendLine(commentText);
            }
            else if (commentText.StartsWith("mm:"))
            {
                sb.AppendLine("// " + commentText);
            }
            else
            {
                sb.AppendLine(commentText);
            }
        }
    }

    private void AppendIndent(System.Text.StringBuilder sb)
    {
        for (int i = 0; i < _currentIndent; i++)
        {
            sb.Append(' ');
        }
    }

    private string EscapeString(string s)
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