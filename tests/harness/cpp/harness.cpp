// MetaMessage C++ test harness - parse JSONC file and re-print to JSONC.
#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <cstdlib>

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

int main(int argc, char **argv) {
    if (argc < 2) {
        std::cerr << "usage: harness <file.jsonc>" << std::endl;
        return 1;
    }

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