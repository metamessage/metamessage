"""
空容器 / nil 绑定的回归用例（对齐 Go zzrepro_test.go）：

1. 空 List[str] 编码后 wire 带 example 哨兵 ""，绑定回来必须是长度 0 的列表；
2. 空 List[int]（child_type=i32）不再因 example 哨兵类型不匹配而报错；
3. 空 Dict 绑定回长度为 0 的 dict；
4. 编码侧带 is_null（真实值为 None）绑定回 None。
"""
from typing import Dict, List, Optional

from metamessage import (
    Field,
    NewTag,
    NodeObject,
    NodeScalar,
    ValueType,
    decode_to_value,
    encode_from_value,
    mm,
    node_to_value,
)


@mm()
class ZzStrBox:
    names: List[str] = mm(allow_empty=True)


@mm()
class ZzIntBox:
    nums: List[int] = mm(allow_empty=True, child_type=ValueType.I32)


@mm()
class ZzInner:
    a: int = mm()
    b: str = mm()


@mm()
class ZzNullBox:
    ptr: Optional[ZzInner] = mm(allow_empty=True)
    m: Dict[str, ZzInner] = mm(allow_empty=True)


@mm()
class ZzNullScalar:
    num: Optional[int] = mm()
    s: Optional[str] = mm()


def test_empty_str_list_binds_to_empty_list():
    wire = encode_from_value(ZzStrBox(names=[]))
    out = decode_to_value(wire, ZzStrBox)
    assert out.names == []


def test_empty_int_list_binds_to_empty_list():
    wire = encode_from_value(ZzIntBox(nums=[]))
    out = decode_to_value(wire, ZzIntBox)
    assert out.nums == []


def test_real_elements_still_bind():
    wire = encode_from_value(ZzIntBox(nums=[5]))
    out = decode_to_value(wire, ZzIntBox)
    assert out.nums == [5]


def test_empty_map_binds_to_empty_dict():
    wire = encode_from_value(ZzNullBox(ptr=ZzInner(a=1, b="x"), m={}))
    out = decode_to_value(wire, ZzNullBox)
    assert out.m == {}


def test_is_null_binds_to_none():
    wire = encode_from_value(ZzNullScalar(num=None, s=None))
    out = decode_to_value(wire, ZzNullScalar)
    assert out.num is None
    assert out.s is None


def test_is_null_obj_node_binds_to_none():
    # Go 侧对 nil 指针产出带 is_null 的 obj 空壳节点，绑定必须是 None
    tag = NewTag()
    tag.type = ValueType.Obj
    tag.is_null = True
    obj = NodeObject(tag=tag, fields=[Field(key="a", value=NodeScalar(data=0, text="0", tag=NewTag()))])
    assert node_to_value(obj, dict) is None


def test_empty_container_without_target_type():
    wire = encode_from_value(ZzStrBox(names=[]))
    assert decode_to_value(wire) == {"names": []}
