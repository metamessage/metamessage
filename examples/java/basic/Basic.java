package io.github.metamessage.examples;

import io.github.metamessage.mm.MetaMessage;
import io.github.metamessage.mm.MM;

@MM
class Person {
    public String name = "Ed";
    public int age = 30;
}

public class Basic {
    public static void main(String[] args) throws Exception {
        // 创建 Person 对象
        Person person = new Person();
        System.out.println("Original: Name=" + person.name + ", Age=" + person.age);

        // 编码到 Wire 格式
        byte[] wire = MetaMessage.encode(person);
        System.out.println("Encoded: " + bytesToHex(wire));

        // 从 Wire 解码
        Person decoded = MetaMessage.decode(wire, Person.class);
        System.out.println("Decoded: Name=" + decoded.name + ", Age=" + decoded.age);
    }

    private static String bytesToHex(byte[] bytes) {
        StringBuilder sb = new StringBuilder();
        for (byte b : bytes) {
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }
}
