package io.metamessage.jsonc;

import io.metamessage.mm.MmTag;
import java.util.List;

public sealed interface JcNode permits JcNode.JcValue, JcNode.JcObject, JcNode.JcArray {

    MmTag tag();

    String path();

    record JcValue(Object data, String text, MmTag tag, String path) implements JcNode {
        @Override
        public MmTag tag() {
            return tag;
        }
    }

    record JcObject(MmTag tag, String path, List<JcField> fields) implements JcNode {
        @Override
        public MmTag tag() {
            return tag;
        }
    }

    record JcArray(MmTag tag, String path, List<JcNode> items) implements JcNode {
        @Override
        public MmTag tag() {
            return tag;
        }
    }
}
