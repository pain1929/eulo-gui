#include "widget.h"

#include <QApplication>
#include <TcpServer.hpp>
#include <boost/asio/thread_pool.hpp>
static boost::asio::io_context g_ctx;
static boost::asio::thread_pool pool{1};

void createJobObject() {
    HANDLE job = CreateJobObject(NULL, NULL);
    if (job == NULL) {
        return; // Handle error
    }

    JOBOBJECT_EXTENDED_LIMIT_INFORMATION jobInfo = {0};
    jobInfo.BasicLimitInformation.LimitFlags = JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE;

    if (!SetInformationJobObject(job, JobObjectExtendedLimitInformation, &jobInfo, sizeof(jobInfo))) {
        CloseHandle(job); // Handle error
        return;
    }

    // Assign the current process and its child processes to the job
    if (!AssignProcessToJobObject(job, GetCurrentProcess())) {
        CloseHandle(job); // Handle error
    }
}


int main(int argc, char *argv[])
{
    createJobObject();

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
