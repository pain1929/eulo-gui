#ifndef FILEREADER
#define FILEREADER
#include <qglobal.h>
#include <string>
#include <vector>
#include <fstream>
class FileReader {


public:
    static std::vector<uint8_t> read(const std::string &fileName) {

        std::ifstream in(fileName , std::ios::binary | std::ios::ate);
        if (!in.is_open()) {return {};}

        // 获取文件大小
        std::streamsize size = in.tellg();
        in.seekg(0, std::ios::beg);

        // 创建 uint8_t 数组来保存文件内容
        std::vector<uint8_t> buffer(size);
        in.read(reinterpret_cast<char*>(buffer.data()), size);
        return buffer;

    }

    static void write(const std::string &fileName , uint8_t * data , size_t size) {

        std::ofstream out(fileName , std::ios::binary | std::ios::out);

        if (!out.is_open()) {return;}


        out.write(reinterpret_cast<const char*>(data) , size);

    }

};


#endif