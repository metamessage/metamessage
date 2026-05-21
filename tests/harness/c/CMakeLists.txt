cmake_minimum_required(VERSION 3.10)
project(mm_harness_c C)

set(CMAKE_C_STANDARD 99)
set(CMAKE_C_STANDARD_REQUIRED ON)

set(MMC_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../../../mm-c")

set(JSONC_SOURCES
    ${MMC_DIR}/src/jsonc/mc_token.c
    ${MMC_DIR}/src/jsonc/mc_scanner.c
    ${MMC_DIR}/src/jsonc/mc_ast.c
    ${MMC_DIR}/src/jsonc/mc_parser.c
)

set(IR_SOURCES
    ${MMC_DIR}/src/ir/mc_ast.c
    ${MMC_DIR}/src/ir/mc_tag.c
)

set(CORE_SOURCES
    ${MMC_DIR}/src/core/mc_encoder.c
    ${MMC_DIR}/src/core/mc_decoder.c
)

set(MM_SOURCES
    ${MMC_DIR}/src/jsonc/mc_printer.c
    ${MMC_DIR}/src/jsonc/mc_ir_parser.c
    ${MMC_DIR}/src/mm.c
)

set(ALL_SOURCES
    ${JSONC_SOURCES}
    ${IR_SOURCES}
    ${CORE_SOURCES}
    ${MM_SOURCES}
    harness.c
)

add_executable(mm_harness_c ${ALL_SOURCES})
target_include_directories(mm_harness_c PRIVATE ${MMC_DIR}/src)