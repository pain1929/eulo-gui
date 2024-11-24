#ifndef WIDGET_H
#define WIDGET_H

#include <QWidget>
#include <boost/filesystem.hpp>
#include <boost/process.hpp>
#include <boost/asio.hpp>
#include <iostream>
#include <fstream>
#include <QThread>
#include <QPaintEvent>
#include <QPainter>
#include "FileReader.h"
#include "settingdlg.h"

#define APP_VERSION "2.0.0"

QT_BEGIN_NAMESPACE
namespace Ui { class Widget; }
QT_END_NAMESPACE

struct Token
{
    std::string code;
    std::string pwd;
    std::string token;

    size_t toData(uint8_t * data) {
        int codeLen = code.length();
        int pwdLen = pwd.length();
        int tokenLen = token.length();

        memcpy(data , &codeLen , 4);
        memcpy(data + 4 , &pwdLen , sizeof(pwd));
        memcpy(data + 8 , &tokenLen , sizeof(tokenLen));

        auto codePos = data + 12;
        memcpy(codePos , code.data() , codeLen);
        auto pwdPos = codePos + codeLen;
        memcpy(pwdPos , pwd.data() , pwdLen);
        auto tokenPos = pwdPos + pwdLen;
        memcpy(tokenPos , token.data() , token.length());

        return 12 + codeLen + pwdLen + tokenLen;

    }

    void load (const std::string & fileName) {

        auto buffer = FileReader::read("token.data");
        if (buffer.empty()) return;
        int codeLen,pwdLen,tokenLen;
        memcpy(&codeLen , buffer.data() , 4);
        memcpy(&pwdLen , buffer.data() +4 , 4);
        memcpy(&tokenLen , buffer.data() +8  , 4);

        code.resize(codeLen);
        pwd.resize(pwdLen);
        token.resize(tokenLen);
        memcpy(code.data() , buffer.data() + 12 , codeLen);
        memcpy(pwd.data() , buffer.data() + 12 + codeLen , pwdLen);
        memcpy(token.data() , buffer.data() + 12 + codeLen + pwdLen , tokenLen);
    }

};

class WorkThread;
class Widget : public QWidget{
    Q_OBJECT

    std::shared_ptr<boost::process::child> eulo;
    std::shared_ptr<boost::process::child> minecraft;
    std::shared_ptr<WorkThread> mcClientThread;

    std::atomic<bool> dead{false};
    Token token;
    GameConfig config_;

    void save(){
        auto buf = new uint8_t[10240];
        auto size = token.toData(buf);

       FileReader::write("token.data" , buf , size);
       delete [] buf;
    }

    void load() {
        token.load("token.data");
    }

public:
    friend WorkThread;
    Widget(QWidget *parent = nullptr);
    ~Widget();

    void paintEvent(QPaintEvent* event) override;

signals:
    void titleMsg(QString msg);
public slots:
    void on_login_clicked(bool);
    void on_titleMsg(QString msg);
    void on_btnUsed(bool used);
    void on_btnTitle(QString msg);
    void on_setting_clicked(bool);
    void on_game_clicked(bool);
    void on_min();
    void on_normal();

private:
    Ui::Widget *ui;
    std::shared_ptr<std::filesystem::path> gamePath; //!< 游戏路径
    std::shared_ptr<std::filesystem::path>  configPath; //!< 配置路径
};


class WorkThread : public QThread {

    Q_OBJECT
    Widget * widget_;
signals:
    void titleMsg(QString msg);
    void btnUsed(bool used);
    void btnTitle(QString str);
    void min();
    void normal();
public:
    explicit WorkThread(Widget * widget) :widget_(widget) {}


    void run() override;

};

#endif // WIDGET_H
