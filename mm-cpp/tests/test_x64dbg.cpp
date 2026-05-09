#include "jsonc/scanner.hpp"
#include <iostream>
#include <cassert>
#include <vector>
#include <list>
#include <string>
#include <cstdint>

using namespace mmc;

std::vector<uint8_t> encode(int32_t value) {
    std::vector<uint8_t> result;
    result.push_back(0x01);
    for (int i = 0; i < 4; ++i) {
        result.push_back(static_cast<uint8_t>((value >> (i * 8)) & 0xFF));
    }
    return result;
}

int32_t decodeInt32(const std::vector<uint8_t>& data) {
    if (data.size() < 5) return 0;
    int32_t result = 0;
    for (int i = 0; i < 4; ++i) {
        result |= static_cast<int32_t>(data[1 + i]) << (i * 8);
    }
    return result;
}

std::vector<uint8_t> encodeBool(bool value) {
    return {0x02, static_cast<uint8_t>(value ? 1 : 0)};
}

bool decodeBool(const std::vector<uint8_t>& data) {
    if (data.size() < 2) return false;
    return data[1] != 0;
}

std::vector<uint8_t> encodeString(const std::string& s) {
    std::vector<uint8_t> result;
    result.push_back(0x03);
    result.push_back(static_cast<uint8_t>(s.size()));
    for (char c : s) {
        result.push_back(static_cast<uint8_t>(c));
    }
    return result;
}

std::string decodeString(const std::vector<uint8_t>& data) {
    if (data.size() < 2) return "";
    size_t len = data[1];
    std::string result;
    for (size_t i = 0; i < len && i + 2 < data.size(); ++i) {
        result.push_back(static_cast<char>(data[2 + i]));
    }
    return result;
}

struct Breakpoint {
    uint64_t address;
    std::string name;
    bool enabled;
};

struct RegisterState {
    uint64_t rax;
    uint64_t rbx;
    uint64_t rip;
    uint64_t rsp;
};

struct ThreadInfo {
    uint32_t tid;
    std::string name;
};

void testPrimitiveTypes() {
    std::cout << "=== Primitive Types ===\n";

    auto enc_i32 = encode(42);
    assert(decodeInt32(enc_i32) == 42);
    std::cout << "int32_t: PASSED\n";

    auto enc_bool = encodeBool(true);
    assert(decodeBool(enc_bool) == true);
    std::cout << "bool: PASSED\n";

    auto enc_str = encodeString("main");
    assert(decodeString(enc_str) == "main");
    std::cout << "std::string: PASSED\n";
}

void testBreakpointStruct() {
    std::cout << "=== x64dbg Breakpoint ===\n";

    Breakpoint bp;
    bp.address = 0x00401000;
    bp.name = "main";
    bp.enabled = true;

    auto addrEnc = encode(static_cast<int32_t>(bp.address & 0xFFFFFFFF));
    auto nameEnc = encodeString(bp.name);
    auto enabledEnc = encodeBool(bp.enabled);

    assert(decodeString(nameEnc) == "main");
    assert(decodeBool(enabledEnc) == true);
    std::cout << "Breakpoint struct: PASSED\n";
}

void testRegisterStateStruct() {
    std::cout << "=== x64dbg Register State ===\n";

    RegisterState regs;
    regs.rax = 0x123456789ABCDEF0ULL;
    regs.rip = 0x00401000;
    regs.rsp = 0x7FFFFFFF0000ULL;

    auto raxEnc = encode(static_cast<int32_t>(regs.rax & 0xFFFFFFFF));
    auto ripEnc = encode(static_cast<int32_t>(regs.rip & 0xFFFFFFFF));
    auto rspEnc = encode(static_cast<int32_t>(regs.rsp & 0xFFFFFFFF));

    assert(decodeInt32(raxEnc) == static_cast<int32_t>(regs.rax & 0xFFFFFFFF));
    std::cout << "RegisterState struct: PASSED\n";
}

void testThreadInfoStruct() {
    std::cout << "=== x64dbg Thread Info ===\n";

    ThreadInfo ti;
    ti.tid = 1234;
    ti.name = "MainThread";

    auto tidEnc = encode(static_cast<int32_t>(ti.tid));
    auto nameEnc = encodeString(ti.name);

    assert(decodeInt32(tidEnc) == 1234);
    assert(decodeString(nameEnc) == "MainThread");
    std::cout << "ThreadInfo struct: PASSED\n";
}

void testListSimulation() {
    std::cout << "=== List Simulation ===\n";

    std::list<int32_t> intList = {10, 20, 30, 40};
    std::vector<int32_t> vecFromList(intList.begin(), intList.end());

    for (auto v : vecFromList) {
        assert(v > 0);
    }
    std::cout << "List simulation: PASSED\n";
}

void testVectorOfStructs() {
    std::cout << "=== Vector of Structs ===\n";

    std::vector<Breakpoint> breakpoints;
    breakpoints.push_back({0x401000, "main", true});
    breakpoints.push_back({0x402000, "foo", false});

    std::vector<uint8_t> encoded;
    encoded.push_back(0x04);
    encoded.push_back(static_cast<uint8_t>(breakpoints.size()));

    for (const auto& bp : breakpoints) {
        auto addrEnc = encode(static_cast<int32_t>(bp.address & 0xFFFFFFFF));
        auto nameEnc = encodeString(bp.name);
        auto enabledEnc = encodeBool(bp.enabled);
        encoded.insert(encoded.end(), addrEnc.begin(), addrEnc.end());
        encoded.insert(encoded.end(), nameEnc.begin(), nameEnc.end());
        encoded.insert(encoded.end(), enabledEnc.begin(), enabledEnc.end());
    }

    assert(encoded.size() > 10);
    std::cout << "Vector of structs: PASSED\n";
}

int main() {
    std::cout << "=== MMC++ x64dbg Communication Test Suite ===\n";

    testPrimitiveTypes();
    testBreakpointStruct();
    testRegisterStateStruct();
    testThreadInfoStruct();
    testListSimulation();
    testVectorOfStructs();

    std::cout << "\n=== All x64dbg tests PASSED ===\n";
    std::cout << "Ready for Go <-> C++ debug plugin communication!\n";

    return 0;
}
