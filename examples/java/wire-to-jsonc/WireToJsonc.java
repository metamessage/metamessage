package io.github.metamessage.examples;

import io.github.metamessage.jsonc.Jsonc;
import io.github.metamessage.mm.MetaMessage;
import io.github.metamessage.jsonc.JcNode;

public class WireToJsonc {
    public static void main(String[] args) throws Exception {
        // 首先从 JSONC 编码到 Wire
        String jsonc = """
            {
                // mm: type=datetime; desc=创建时间
                "create_time": "2026-01-01 00:00:00",
                // mm: type=str; desc=用户名称
                "user_name": "Alice",
                // mm: type=bool; desc=是否激活
                "is_active": true,
                // mm: type=array; child_type=i
                "scores": [95, 87, 92]
            }
            """;

        JcNode node = Jsonc.parseFromString(jsonc);
        byte[] wire = MetaMessage.encode(node);

        System.out.println("Original JSONC:");
        System.out.println(jsonc);

        System.out.println("\nEncoded Wire:");
        System.out.println(bytesToHex(wire));

        // 从 Wire 解码到 JcNode
        JcNode decodedNode = MetaMessage.decode(wire, JcNode.class);

        // 转换回 JSONC
        String outputJsonc = Jsonc.toString(decodedNode);

        System.out.println("\nDecoded to JSONC:");
        System.out.println(outputJsonc);
    }

    private static String bytesToHex(byte[] bytes) {
        StringBuilder sb = new StringBuilder();
        for (byte b : bytes) {
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }
}
