// MetaMessage Kotlin test harness - parse JSONC file and re-print to JSONC.
import io.github.metamessage.MetaMessage
import io.github.metamessage.jsonc.parseFromJsonc
import io.github.metamessage.jsonc.toJsonc
import java.io.File
import kotlin.system.exitProcess

fun main(args: Array<String>) {
    if (args.isEmpty()) {
        System.err.println("usage: harness [--encode|--decode] <file.jsonc>")
        exitProcess(1)
    }

    if (args[0] == "--encode") {
        if (args.size < 2) {
            System.err.println("usage: harness --encode <file.jsonc>")
            exitProcess(1)
        }
        val input: String
        try {
            input = File(args[1]).readText()
        } catch (e: Exception) {
            System.err.println("read error: ${e.message}")
            exitProcess(1)
        }
        try {
            val wire = MetaMessage.encodeFromJsonc(input)
            val hex = wire.joinToString("") { "%02x".format(it) }
            print(hex)
        } catch (e: Exception) {
            System.err.println("encode error: ${e.message}")
            exitProcess(1)
        }
        return
    }

    if (args[0] == "--decode") {
        val hexStr = System.`in`.bufferedReader().readText().trim()
        try {
            val wire = hexStr.chunked(2).map { it.toInt(16).toByte() }.toByteArray()
            val output = MetaMessage.decodeToJsonc(wire)
            print(output)
        } catch (e: Exception) {
            System.err.println("decode error: ${e.message}")
            exitProcess(1)
        }
        return
    }

    // Existing behavior
    val input: String
    try {
        input = File(args[0]).readText()
    } catch (e: Exception) {
        System.err.println("read error: ${e.message}")
        exitProcess(1)
    }

    val node = try {
        parseFromJsonc(input)
    } catch (e: Exception) {
        System.err.println("parse error: ${e.message}")
        exitProcess(1)
    }

    val output = toJsonc(node)
    print(output)
}