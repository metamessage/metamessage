package io.github.metamessage.jsonc

sealed class JsoncNode {
    abstract var tag: JsoncTag?
    abstract var path: String
}

data class JsoncValue(
    var data: Any? = null,
    var text: String = "",
    override var tag: JsoncTag? = null,
    override var path: String = ""
) : JsoncNode()

data class JsoncObject(
    var fields: MutableList<JsoncField> = mutableListOf(),
    override var tag: JsoncTag? = null,
    override var path: String = ""
) : JsoncNode()

data class JsoncArray(
    var items: MutableList<JsoncNode> = mutableListOf(),
    override var tag: JsoncTag? = null,
    override var path: String = ""
) : JsoncNode()

data class JsoncField(
    val key: String,
    val value: JsoncNode?
)

data class JsoncTag(
    var name: String = "",
    var isNull: Boolean = false,
    var type: JsoncValueType = JsoncValueType.Unknown,
    var desc: String = "",
    var raw: Boolean = false,
    var nullable: Boolean = false,
    var allowEmpty: Boolean = false,
    var unique: Boolean = false,
    var defaultValue: String = "",
    var min: String = "",
    var max: String = "",
    var size: Int = 0,
    var enum: String = "",
    var pattern: String = "",
    var location: String = "",
    var version: Int = 0,
    var mime: String = "",
    var childDesc: String = "",
    var childType: JsoncValueType = JsoncValueType.Unknown,
    var childRaw: Boolean = false,
    var childNullable: Boolean = false,
    var childAllowEmpty: Boolean = false,
    var childUnique: Boolean = false,
    var childDefault: String = "",
    var childMin: String = "",
    var childMax: String = "",
    var childSize: Int = 0,
    var childEnum: String = "",
    var childPattern: String = "",
    var childLocation: String = "",
    var childVersion: Int = 0,
    var childMime: String = "",
    var isInherit: Boolean = false
)

enum class JsoncValueType {
    Unknown,
    String,
    Int,
    Int8,
    Int16,
    Int32,
    Int64,
    Uint,
    Uint8,
    Uint16,
    Uint32,
    Uint64,
    Float32,
    Float64,
    Bool,
    Bytes,
    BigInt,
    DateTime,
    Date,
    Time,
    UUID,
    Decimal,
    IP,
    URL,
    Email,
    Enum,
    Array,
    Struct,
    Slice,
    Map,
    Null,
    Raw
}