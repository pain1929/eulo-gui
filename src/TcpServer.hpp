#ifndef TCPSERVER_HPP
#define TCPSERVER_HPP
#include <qglobal.h>
#include <boost/asio.hpp>

#include "message/Msg.h"
#include "MsgQue.h"



using boost::asio::ip::tcp;
namespace asio = boost::asio;

// 会话类：管理单个客户端的连接
class TcpSession : public std::enable_shared_from_this<TcpSession> {
public:
    TcpSession(tcp::socket socket) : socket_(std::move(socket)) {}

    void start() {
        read();
    }

private:
    void read() {
        auto self(shared_from_this());

        auto msgHdr = new MessageHdr;

        boost::asio::async_read(socket_ , asio::buffer(msgHdr , sizeof(MessageHdr)) , [this , self , msgHdr](boost::system::error_code ec , size_t len)
        {
            if (ec)
            {
                close();
                return;
            }

            if (msgHdr->zero != 0xAB || !msgHdr->msgLen) //消息类型不匹配
            {
                close();
                return;
            }
            auto buf = new uint8_t[msgHdr->msgLen];
            boost::asio::read(socket_ , asio::buffer(buf , msgHdr->msgLen));

            // buf 前8个字节是 消息类型 见 message/Msg.h 中 MsgType 类型
            auto msgType = *reinterpret_cast<MsgType*>(buf);
            switch (msgType)
            {
            case MsgType::EuloMsg :
                {
                    auto euloMsg = std::make_shared<EuloMsgType>();
                    euloMsg->load(buf + sizeof (MsgType));
                    MsgRegister::obj().setMsg(euloMsg);
                    break;
                }
                default:
                    break;
            }
            delete [] buf;
            delete msgHdr;
            read();
        });

    }

    void write(std::size_t length) {

    }

    void close()
    {
        socket_.close();
    }

    tcp::socket socket_;

};


// 服务器类：管理所有客户端的连接
class TcpServer : public std::enable_shared_from_this<TcpServer>{
public:
    TcpServer(asio::io_context& io_context, short port)
        : acceptor_(io_context, tcp::endpoint(tcp::v4(), port)) {
        accept();
    }

private:
    void accept() {
        acceptor_.async_accept(
            [this](boost::system::error_code ec, tcp::socket socket) {
                if (!ec) {
                    std::make_shared<TcpSession>(std::move(socket))->start();
                    accept();  // 接受下一个连接
                }

            });
    }

    tcp::acceptor acceptor_;
};





#endif