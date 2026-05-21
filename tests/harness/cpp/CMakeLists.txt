cmake_minimum_required(VERSION 3.10)
project(mm_harness_cpp CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

set(MMCPP_DIR "${CMAKE_CURRENT_SOURCE_DIR}/../../../mm-cpp")

add_executable(mm_harness_cpp harness.cpp)
target_include_directories(mm_harness_cpp PRIVATE ${MMCPP_DIR}/src)