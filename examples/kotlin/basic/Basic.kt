package io.metamessage.examples

import io.metamessage.mm.MetaMessage
import io.metamessage.mm.MM

@MM
class Person(var name: String = "Ed", var age: Int = 30)

fun main() {
    // 创建 Person 对象
    val person = Person()
    println("Original: Name=${person.name}, Age=${person.age}")

    // 编码到 Wire 格式
    val wire = MetaMessage.encode(person)
    println("Encoded: ${bytesToHex(wire)}")

    // 从 Wire 解码
    val decoded = MetaMessage.decode(wire, Person::class.java)
    println("Decoded: Name=${decoded.name}, Age=${decoded.age}")
}

fun bytesToHex(bytes: ByteArray): String {
    val sb = StringBuilder()
    for (b in bytes) {
        sb.append(String.format("%02x", b))
    }
    return sb.toString()
}
