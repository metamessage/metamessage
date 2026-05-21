// MetaMessage Kotlin test harness - parse JSONC file and re-print to JSONC.
import io.github.metamessage.jsonc.parseFromJsonc
import io.github.metamessage.jsonc.toJsonc
import java.io.File
import kotlin.system.exitProcess

fun main(args: Array<String>) {
    if (args.isEmpty()) {
        System.err.println("usage: harness <file.jsonc>")
        exitProcess(1)
    }

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