#ifndef MSG_H
#define MSG_H
#include <cstdint>


struct MessageHdr{
    uint16_t zero{0xAB}; //!<头部 规定为 0xAB
    uint16_t msgLen{}; //!<消息长度
};

struct NormalMsg {
    bool started{false}; //!<端口启动成功
    bool pad{}; //!<填充
    uint16_t errorMsgLen{}; //!<错误消息长度
    std::vector<uint8_t> errorMsg{}; //!<错误消息

    void load (uint8_t * data)
    {
        memcpy(this , data , 4);
        auto errmsgBegin = data + 4;
        auto end = errmsgBegin + errorMsgLen;
        errorMsg = std::vector<uint8_t>(errmsgBegin , end);
    }

};
#endif