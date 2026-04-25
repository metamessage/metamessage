package io.metamessage.examples;

import io.metamessage.jsonc.Jsonc;
import io.metamessage.mm.MetaMessage;
import io.metamessage.jsonc.JcNode;

public class JsoncToWire {
    public static void main(String[] args) throws Exception {
        // JSONC 字符串
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

        System.out.println("Input JSONC:");
        System.out.println(jsonc);

        // 解析 JSONC 到 JcNode
        JcNode node = Jsonc.parseFromString(jsonc);
        System.out.println("\nParsed JcNode:");
        System.out.println(Jsonc.toString(node));

        // 编码到 Wire 格式
        byte[] wire = MetaMessage.encode(node);
        System.out.println("\nEncoded Wire:");
        System.out.println(bytesToHex(wire));
    }

    private static String bytesToHex(byte[] bytes) {
        StringBuilder sb = new StringBuilder();
        for (byte b : bytes) {
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }
}
