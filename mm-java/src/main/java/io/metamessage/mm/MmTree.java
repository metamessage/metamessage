package io.metamessage.mm;

import java.util.List;
import java.util.Map;

/** Decoded MM value tree (Go {@code ast.Node} equivalent). */
public sealed interface MmTree permits MmTree.MmScalar, MmTree.MmObject, MmTree.MmArray {

    MmTag tag();

    record MmScalar(Object data, String text, MmTag tag) implements MmTree {}

    record MmObject(MmTag tag, List<Map.Entry<String, MmTree>> fields) implements MmTree {}

    record MmArray(MmTag tag, List<MmTree> items) implements MmTree {}
}
