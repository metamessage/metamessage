namespace MetaMessage.Ir;

public enum MmNodeType
{
    Unknown = 0,
    Object = 1,
    Array = 2,
    Value = 3,
    Doc = 4
}

public interface IMmTree
{
    Tag Tag { get; }
    string Path { get; set; }
    MmNodeType NodeType { get; }
}

public class MmScalar : IMmTree
{
    public object Data { get; set; }
    public string Text { get; set; }
    public Tag Tag { get; set; }
    public string Path { get; set; } = "";
    public MmNodeType NodeType => MmNodeType.Value;

    public MmScalar(object data, string text, Tag tag)
    {
        Data = data;
        Text = text;
        Tag = tag;
    }
}

public class MmArray : IMmTree
{
    public List<IMmTree> Children { get; set; }
    public Tag Tag { get; set; }
    public string Path { get; set; } = "";
    public MmNodeType NodeType => MmNodeType.Array;

    public MmArray(List<IMmTree> children, Tag tag)
    {
        Children = children;
        Tag = tag;
    }
}

public class MmMap : IMmTree
{
    public List<KeyValuePair<MmScalar, IMmTree>> Entries { get; set; }
    public Tag Tag { get; set; }
    public string Path { get; set; } = "";
    public MmNodeType NodeType => MmNodeType.Object;

    public MmMap(List<KeyValuePair<MmScalar, IMmTree>> entries, Tag tag)
    {
        Entries = entries;
        Tag = tag;
    }
}

public class MmDoc : IMmTree
{
    public List<KeyValuePair<MmScalar, IMmTree>> Fields { get; set; }
    public Tag Tag { get; set; }
    public string Path { get; set; } = "";
    public MmNodeType NodeType => MmNodeType.Doc;

    public MmDoc(List<KeyValuePair<MmScalar, IMmTree>> fields, Tag tag)
    {
        Fields = fields;
        Tag = tag;
    }
}