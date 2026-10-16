package io.github.metamessage

import io.github.metamessage.ir.ValueType
import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test

/**
 * 空容器 / is_null 绑定回归用例（对齐 Go zzrepro_test.go）：
 * 编码侧为空容器时会注入 example 哨兵项，绑定时必须跳过，
 * 产出空的非 null 集合；is_null 绑定成 null。
 */
class EmptyContainerBindTest {

    @MM
    class StrBox(
            @MM(type = ValueType.VEC, childType = ValueType.STR, allowEmpty = true)
            var names: List<String> = listOf()
    )

    @MM
    class IntBox(
            @MM(type = ValueType.VEC, childType = ValueType.I32, allowEmpty = true)
            var nums: List<Int> = listOf()
    )

    @MM
    class MapBox(
            @MM(type = ValueType.MAP, childType = ValueType.STR, allowEmpty = true)
            var m: Map<String, String> = mapOf()
    )

    @Test
    fun emptyStringListBindsToEmptyList() {
        val wire = MetaMessage.encodeFromValue(StrBox(listOf()))
        val out = MetaMessage.decodeToValue(wire, StrBox::class.java)
        assertNotNull(out.names, "空容器应绑定成非 null 空列表")
        assertEquals(0, out.names.size, "example 哨兵项不得进入绑定结果")
    }

    @Test
    fun emptyIntListBindsToEmptyList() {
        val wire = MetaMessage.encodeFromValue(IntBox(listOf()))
        val out = MetaMessage.decodeToValue(wire, IntBox::class.java)
        assertNotNull(out.nums)
        assertEquals(0, out.nums.size)
    }

    @Test
    fun realElementsStillBind() {
        val wire = MetaMessage.encodeFromValue(IntBox(listOf(5)))
        val out = MetaMessage.decodeToValue(wire, IntBox::class.java)
        assertEquals(listOf(5), out.nums)
    }

    @Test
    fun emptyMapBindsToEmptyMap() {
        val wire = MetaMessage.encodeFromValue(MapBox(mapOf()))
        val out = MetaMessage.decodeToValue(wire, MapBox::class.java)
        assertNotNull(out.m)
        assertEquals(0, out.m.size)
    }
}
