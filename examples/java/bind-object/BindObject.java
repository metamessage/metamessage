package io.github.metamessage.examples;

import io.github.metamessage.jsonc.Jsonc;
import io.github.metamessage.mm.MM;
import io.github.metamessage.mm.ValueType;
import java.util.List;

@MM
class User {
    public String name;
    public int age;
    public boolean active;
    @MM(childType = ValueType.INT)
    public List<Integer> scores;
}

public class BindObject {
    public static void main(String[] args) {
        // JSONC 字符串
        String jsonc = """
            {
                // mm: type=str; desc=姓名
                "name": "Alice",
                // mm: type=i; desc=年龄
                "age": 25,
                // mm: type=bool; desc=是否激活
                "active": true,
                // mm: type=array; child_type=i; desc=分数
                "scores": [95, 87, 92]
            }
            """;

        System.out.println("Input JSONC:");
        System.out.println(jsonc);

        // 从 JSONC 绑定到 User 对象
        User user = Jsonc.bindFromString(jsonc, User.class);

        System.out.println("\nBound to object:");
        System.out.println("Name: " + user.name);
        System.out.println("Age: " + user.age);
        System.out.println("Active: " + user.active);
        System.out.println("Scores: " + user.scores);
    }
}
