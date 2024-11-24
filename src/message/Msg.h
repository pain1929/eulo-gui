#ifndef MSG_H
#define MSG_H
#include <cstdint>


// 消息类型
enum struct MsgType : uint64_t {
    EuloMsg = 0, //!< eulogist 进程 返回的启动信息 通知启动成功或失败 启动失败返回错误消息
    SendCmdMsg = 1, //!< GUI进程 发送给 eulogist进程的命令信息 用于发送给租赁服
    MCPacket = 2  //!< GUI进程 和 eulo 进行通信发送基岩协议数据包 并非原始数据 需要由eulo进程进行分析
};


// 消息头部 8字节
struct MessageHdr{
    uint16_t zero{0xAB}; //!<头部 规定为 0xAB
    uint16_t pad{}; //!<填充
    uint32_t msgLen{}; //!<消息长度
};

// 消息类型
struct MessageType {
    MsgType msgType; //!< 消息类型

};

// 所有的tcp通讯必须使用以下结构
// MessageHdr（8字节） + MessageType（8字节） + 消息结构 例如：EuloMsgType

// eulo进程发送给GUI进程通知eulo启动结果
struct EuloMsgType {

    bool started{false}; //!<eulo进程端口启动成功

    bool pad{}; //!<填充

    uint16_t errorMsgLen{}; //!<启动失败发生错误时此字段表示错误消息长度 若启动成功 本字段必须为0

    std::vector<uint8_t> errorMsg{}; //!<错误消息 UTF-8格式字符



    void load (uint8_t * data) {
        memcpy(this , data , 4);
        auto errmsgBegin = data + 4;
        auto end = errmsgBegin + errorMsgLen;
        errorMsg = std::vector<uint8_t>(errmsgBegin , end);
    }

};

// 由GUI进程发送给 eulo进程的 信息 用于发送命令到租赁服
struct EuloSendCmdMsgType {

    uint32_t cmdSize; //!< 命令长度 单位字节
    std::vector<uint8_t> cmd; //!<命令字符串 utf-8

    /**
     * 返回本数据包大小 单位字节
     * @return
     */
    size_t getSize () const {
        return cmdSize + 4;
    }

    void toData(uint8_t * data) {
        memcpy(data , this , 4);
        memcpy( data + 4 , cmd.data() , cmdSize);
    }

};


#endif