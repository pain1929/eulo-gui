cmake_minimum_required(VERSION 3.5)

project(Eulo-GUI VERSION 2.1.1 LANGUAGES CXX)

if(MSVC)
    # 设置 MSVC 编译器使用 UTF-8 编码
    add_compile_options("/utf-8")
endif()

set(CMAKE_INCLUDE_CURRENT_DIR ON)

set(CMAKE_AUTOUIC ON)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

if (${CMAKE_BUILD_TYPE} STREQUAL "Debug")
    set(CMAKE_PREFIX_PATH $ENV{BOOST_PATH_DEBUG})
else ()
    set(CMAKE_PREFIX_PATH $ENV{BOOST_PATH_RELEASE})
endif ()


find_package(Qt5 COMPONENTS Widgets REQUIRED)
find_package(Boost COMPONENTS system filesystem thread REQUIRED)
add_executable(${PROJECT_NAME} WIN32
        main.cpp
        widget.cpp
        widget.h
        widget.ui
        res.qrc
        app.rc
        TcpServer.hpp
        MsgQue.h
        message/Msg.h
        DragDropLineEdit.h
        settingdlg.cpp
        settingdlg.h
        settingdlg.ui
        FileReader.h
)

target_link_libraries(${PROJECT_NAME} PRIVATE Qt5::Widgets Boost::system Boost::filesystem Boost::thread)
