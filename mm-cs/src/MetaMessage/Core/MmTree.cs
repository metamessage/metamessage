namespace MetaMessage.Core;

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
    MmTag Tag { get; }
    string Path { get; set; }
    MmNodeType NodeType { get; }
}

public class MmScalar : IMmTree
{
    public object Data { get; set; }
    public string Text { get; set; }
    public MmTag Tag { get; set; }
    public string Path { get; set; } = "";
    public MmNodeType NodeType => MmNodeType.Value;

    public MmScalar(object data, string text, MmTag tag)
    {
        Data = data;
        Text = text;
        Tag = tag;
    }
}

public class MmArray : IMmTree
{
    public List<IMmTree> Children { get; set; }
    public MmTag Tag { get; set; }
    public string Path { get; set; } = "";
    public MmNodeType NodeType => MmNodeType.Array;

    public MmArray(List<IMmTree> children, MmTag tag)
    {
        Children = children;
        Tag = tag;
    }
}

public class MmMap : IMmTree
{
    public List<KeyValuePair<MmScalar, IMmTree>> Entries { get; set; }
    public MmTag Tag { get; set; }
    public string Path { get; set; } = "";
    public MmNodeType NodeType => MmNodeType.Object;

    public MmMap(List<KeyValuePair<MmScalar, IMmTree>> entries, MmTag tag)
    {
        Entries = entries;
        Tag = tag;
    }
}

public class MmDoc : IMmTree
{
    public List<KeyValuePair<MmScalar, IMmTree>> Fields { get; set; }
    public MmTag Tag { get; set; }
    public string Path { get; set; } = "";
    public MmNodeType NodeType => MmNodeType.Doc;

    public MmDoc(List<KeyValuePair<MmScalar, IMmTree>> fields, MmTag tag)
    {
        Fields = fields;
        Tag = tag;
    }
}