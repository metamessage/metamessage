namespace MetaMessage.Mm;

public interface IMmTree
{
    MmTag Tag { get; }
}

public class MmScalar : IMmTree
{
    public object Data { get; set; }
    public string Text { get; set; }
    public MmTag Tag { get; set; }

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

    public MmMap(List<KeyValuePair<MmScalar, IMmTree>> entries, MmTag tag)
    {
        Entries = entries;
        Tag = tag;
    }
}