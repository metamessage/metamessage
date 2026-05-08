cmake_minimum_required(VERSION 3.10)
project(mmcpp CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include_directories(src)

set(JSONC_SOURCES
    src/jsonc/scanner.cpp
    src/jsonc/parser.cpp
)

set(MM_SOURCES
    src/mm/codec.cpp
)

add_library(mmcpp_jsonc STATIC ${JSONC_SOURCES})
add_library(mmcpp_mm STATIC ${MM_SOURCES})

enable_testing()

add_executable(mmcpp_test
    tests/test_jsonc.cpp
    src/jsonc/scanner.cpp
    src/jsonc/parser.cpp
)

add_executable(mmcpp_x64dbg_test
    tests/test_x64dbg.cpp
    src/jsonc/scanner.cpp
    src/jsonc/parser.cpp
    src/mm/codec.cpp
)

target_include_directories(mmcpp_test PRIVATE src)
target_include_directories(mmcpp_x64dbg_test PRIVATE src)

add_test(NAME mmcpp_jsonc_test COMMAND mmcpp_test)
add_test(NAME mmcpp_x64dbg_test COMMAND mmcpp_x64dbg_test)
