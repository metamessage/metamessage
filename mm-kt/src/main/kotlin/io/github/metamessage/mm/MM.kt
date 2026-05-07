package io.github.metamessage.mm

@Target(AnnotationTarget.FIELD, AnnotationTarget.CLASS)
@Retention(AnnotationRetention.RUNTIME)
annotation class MM(
    val name: String = "",
    val type: ValueType = ValueType.UNKNOWN,
    val isNull: Boolean = false,
    val example: Boolean = false,
    val desc: String = "",
    val raw: Boolean = false,
    val nullable: Boolean = false,
    val allowEmpty: Boolean = false,
    val unique: Boolean = false,
    val defaultValue: String = "",
    val min: String = "",
    val max: String = "",
    val size: Int = 0,
    val enumValues: String = "",
    val pattern: String = "",
    val location: Int = 0,
    val version: Int = 0,
    val mime: String = "",
    val childDesc: String = "",
    val childType: ValueType = ValueType.UNKNOWN,
    val childRaw: Boolean = false,
    val childNullable: Boolean = false,
    val childAllowEmpty: Boolean = false,
    val childUnique: Boolean = false,
    val childDefault: String = "",
    val childMin: String = "",
    val childMax: String = "",
    val childSize: Int = 0,
    val childEnum: String = "",
    val childPattern: String = "",
    val childLocation: Int = 0,
    val childVersion: Int = 0,
    val childMime: String = ""
)
