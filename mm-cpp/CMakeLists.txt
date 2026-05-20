cmake_minimum_required(VERSION 3.10)
project(mmcpp CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include_directories(src)

set(JSONC_SOURCES
    src/jsonc/scanner.cpp
)

set(IR_SOURCES
    src/ir/value_type.hpp
    src/ir/tag.hpp
    src/ir/ast.hpp
)

set(CORE_SOURCES
    src/core/constants.hpp
    src/core/encoder.hpp
    src/core/decoder.hpp
)

set(MM_SOURCES
    src/mm/macro.hpp
    src/mm/mm.hpp
)

enable_testing()

add_executable(mmcpp_test
    tests/test_jsonc.cpp
    src/jsonc/scanner.cpp
)

add_executable(mmcpp_comprehensive_test
    tests/test_comprehensive.cpp
    src/jsonc/scanner.cpp
)

target_include_directories(mmcpp_test PRIVATE src)
target_include_directories(mmcpp_comprehensive_test PRIVATE src)

add_test(NAME mmcpp_jsonc_test COMMAND mmcpp_test)
add_test(NAME mmcpp_comprehensive_test COMMAND mmcpp_comprehensive_test)