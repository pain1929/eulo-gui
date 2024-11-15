#include "widget.h"

#include <QApplication>
#include <TcpServer.hpp>
#include <boost/asio/thread_pool.hpp>
static boost::asio::io_context g_ctx;
static boost::asio::thread_pool pool{1};

int main(int argc, char *argv[])
{

    TcpServer server(g_ctx , 1930);


    boost::asio::post(pool , [] ()
    {
       g_ctx.run();
    });

    QApplication a(argc, argv);
    Widget w;
    w.show();
    return QApplication::exec();
}
