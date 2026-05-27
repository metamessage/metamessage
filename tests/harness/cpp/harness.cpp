// MetaMessage C++ test harness - parse JSONC file and re-print to JSONC.
#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <vector>
#include <cstdlib>
#include <iomanip>

#include "jsonc/scanner.hpp"
#include "mm/mm.hpp"

static std::string read_file(const char *path) {
    std::ifstream f(path);
    if (!f) {
        std::cerr << "read error: cannot open " << path << std::endl;
        std::exit(1);
    }
    std::ostringstream ss;
    ss << f.rdbuf();
    return ss.str();
}

static std::string bytes_to_hex(const std::vector<uint8_t>& bytes) {
    std::ostringstream ss;
    for (uint8_t b : bytes) {
        ss << std::hex << std::setfill('0') << std::setw(2) << (int)b;
    }
    return ss.str();
}

static std::vector<uint8_t> hex_to_bytes(const std::string& hex) {
    std::vector<uint8_t> bytes;
    for (size_t i = 0; i < hex.length(); i += 2) {
        unsigned int byte;
        std::istringstream(hex.substr(i, 2)) >> std::hex >> byte;
        bytes.push_back(static_cast<uint8_t>(byte));
    }
    return bytes;
}

int main(int argc, char **argv) {
    if (argc < 2) {
        std::cerr << "usage: harness [--encode|--decode] <file.jsonc>" << std::endl;
        return 1;
    }

    std::string arg1 = argv[1];

    if (arg1 == "--encode") {
        if (argc < 3) {
            std::cerr << "usage: harness --encode <file.jsonc>" << std::endl;
            return 1;
        }
        std::string input = read_file(argv[2]);
        std::vector<uint8_t> wire = mmc::mm::fromJSONC(input);
        std::cout << bytes_to_hex(wire);
        return 0;
    }

    if (arg1 == "--decode") {
        std::string hex_str;
        std::string line;
        while (std::getline(std::cin, line)) {
            hex_str += line;
        }
        std::vector<uint8_t> wire = hex_to_bytes(hex_str);
        std::string output = mmc::mm::toJSONC(wire);
        std::cout << output;
        return 0;
    }

    // Existing behavior
    std::string input = read_file(argv[1]);
    auto node = mmc::mm::parseJSONC(input);
    if (!node) {
        std::cerr << "parse error" << std::endl;
        return 1;
    }

    std::string output = mmc::jsonc::toJSONC(node);
    std::cout << output;
    return 0;
}