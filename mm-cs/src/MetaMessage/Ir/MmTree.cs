namespace MetaMessage.Ir;

public enum MmNodeType
{
    Unknown = 0,
    Object = 1,
    Array = 2,
    Value = 3,
    Doc = 4
}

public interface INode
{
    Tag Tag { get; set; }
    string Path { get; set; }
    MmNodeType NodeType { get; }
}

public class NodeScalar : INode
{
    public object? Data { get; set; }
    public string Text { get; set; }
    public Tag TTag { get; set; }
    public string Path { get; set; } = "";
    public MmNodeType NodeType => MmNodeType.Value;

    public NodeScalar(object? data, string text, Tag tag)
    {
        Data = data;
        Text = text;
        TTag = tag;
    }
}

public class MmArray : INode
{
    public List<INode> Children { get; set; }
    public Tag Tag { get; set; }
    public string Path { get; set; } = "";
    public MmNodeType NodeType => MmNodeType.Array;

    public MmArray(List<INode> children, Tag tag)
    {
        Children = children;
        Tag = tag;
    }
}

public class MmMap : INode
{
    public List<KeyValuePair<NodeScalar, INode>> Entries { get; set; }
    public Tag Tag { get; set; }
    public string Path { get; set; } = "";
    public MmNodeType NodeType => MmNodeType.Object;

    public MmMap(List<KeyValuePair<NodeScalar, INode>> entries, Tag tag)
    {
        Entries = entries;
        Tag = tag;
    }
}

public class MmDoc : INode
{
    public List<KeyValuePair<NodeScalar, INode>> Fields { get; set; }
    public Tag Tag { get; set; }
    public string Path { get; set; } = "";
    public MmNodeType NodeType => MmNodeType.Doc;

    public MmDoc(List<KeyValuePair<NodeScalar, INode>> fields, Tag tag)
    {
        Fields = fields;
        Tag = tag;
    }
}