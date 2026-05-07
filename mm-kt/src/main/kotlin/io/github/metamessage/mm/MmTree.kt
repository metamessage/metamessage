package io.github.metamessage.mm

sealed interface MmTree {
    val tag: MmTag

    data class MmScalar(val data: Any?, val text: String, override val tag: MmTag) : MmTree
    data class MmObject(override val tag: MmTag, val fields: List<Pair<String, MmTree>>) : MmTree
    data class MmArray(override val tag: MmTag, val items: List<MmTree>) : MmTree
}
