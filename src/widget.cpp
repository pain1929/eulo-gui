#include "widget.h"

#include <filesystem>

#include "./ui_widget.h"
#include <QMessageBox>
#include <boost/process/windows.hpp>
#include "MsgQue.h"
Widget::Widget(QWidget *parent)
    : QWidget(parent)
    , ui(new Ui::Widget)
{

    SettingDlg dlg;
    config_ = dlg.config;

    ui->setupUi(this);
    this->setWindowIcon(QIcon(":/image/ico.png"));
    // 设置样式
    // 创建 QPixmap，加载图片
    QPixmap pixmap(":/image/title.png");

    // 设置 QLabel 显示图片
    ui->title->setPixmap(pixmap);

    // 可选：设置图片的缩放方式
    ui->title->setScaledContents(true);  // 图片会自动缩放以适应标签大小
    this->setStyleSheet(R"(

QPushButton {
            background-color: #6b6b6b;
            color: #ffffff;
            font: bold 16px "Minecraft";
            border: 2px solid #3c3c3c;
            border-radius: 5px;
            padding: 5px;
        }
        QPushButton:hover {
            background-color: #828282;
            border: 2px solid #5c5c5c;
        }
        QPushButton:pressed {
            background-color: #505050;
            border: 2px solid #2c2c2c;
        }



QMessageBox {
      background-color: white;
      color: black;
}
QMessageBox QLabel {
     color: black;
}



 QLineEdit {
            background-color: rgba(50, 50, 50, 180);
            color: #ffffff;
            font: 14px "Minecraft";
            border: 2px solid #3c3c3c;
            border-radius: 5px;
            padding: 5px;
        }
        QLineEdit:focus {
            border: 2px solid #5c5c5c;
        }
        QLineEdit::placeholder {
            color: #aaaaaa;
            font-style: italic;
        }

/* 标签样式 */
QLabel {
    color: white; /* 标签文字颜色 */
    font-family: 'Minecraft', sans-serif;
    font-size: 16px;
}

/* 下拉框样式 */
QComboBox {
    background-color: #2c2c2c;
    border: 2px solid #555555;
    color: white;
    font-family: 'Minecraft', sans-serif;
    font-size: 14px;
    padding: 5px;
    border-radius: 4px;
}

/* 下拉框项目样式 */
QComboBox QAbstractItemView {
    background-color: #2c2c2c;
    border: 1px solid #555555;
    color: white;
    font-family: 'Minecraft', sans-serif;
    font-size: 14px;
}

/* 滚动条样式 */
QScrollBar:vertical, QScrollBar:horizontal {
    border: 2px solid #3a3a3a;
    background-color: #2c2c2c;
    width: 10px;
    height: 10px;
}

QScrollBar::handle:vertical, QScrollBar::handle:horizontal {
    background-color: #5a5a5a;
    border-radius: 5px;
}

QScrollBar::handle:vertical:hover, QScrollBar::handle:horizontal:hover {
    background-color: #7a7a7a;
}
    )");

    this->setWindowTitle("Eulo-GUI");
    load();
    ui->code->setText(token.code.c_str());
    ui->pwd->setText(token.pwd.c_str());
    ui->token->setText(token.token.c_str());
    ui->code->setPlaceholderText("请输入租赁服号");
    ui->pwd->setPlaceholderText("请输入租赁服密码 没有则空白");
    ui->token->setPlaceholderText("输入你的 FastBuilder Token 或者拖入Token文件");

}

Widget::~Widget()
{
    delete ui;
    dead = true;
    if (this->minecraft) { this->minecraft->terminate();}
    if (this->eulo) { this->eulo->terminate();}
    if (this->mcClientThread) {this->mcClientThread->terminate();}
}

void Widget::paintEvent(QPaintEvent* event)
{
    // 创建 QPainter 对象
    QPainter painter(this);

    // 加载背景图片
    QPixmap background(":/image/bp.png");

    // 绘制背景图片，指定其位置和大小
    painter.drawPixmap(0, 0, width(), height(), background);

    // 设置文字的字体和颜色
    QFont font("Arial", 10);
    painter.setFont(font);
    painter.setPen(Qt::white); // 白色文字，方便在背景上显示

    // 要绘制的文字
    QString text = "赞颂者GUI v" + QString(APP_VERSION);

    // 计算右下角的位置
    int margin = 10; // 距离窗口边界的边距
    int x = 5;
    int y = height() - margin;

    // 绘制文字在右下角
    painter.drawText(x, y, text);
}

void Widget::on_login_clicked(bool) {

    auto code = ui->code->text();
    auto pwd = ui->pwd->text();
    auto token = ui->token->text();
    this->token.token = token.toStdString();
    this->token.code = code.toStdString();
    this->token.pwd = pwd.toStdString();
    save();
    if (this->eulo) { eulo->terminate();}
    if (this->minecraft) {minecraft->terminate();}
    if (this->mcClientThread) {mcClientThread->terminate();}

    on_btnTitle("启动中...");

    try
    {
        this->eulo = std::make_shared<boost::process::child>("Eulogist.exe" ,
            code.toStdString() ,
            pwd.toStdString() ,
            token.toStdString() ,boost::process::windows::hide);
    }
    catch (const std::exception &e)
    {
        QMessageBox::critical(this , "提示" , "无法找到启动 Eulogist 进程");
    }

    on_btnUsed(false);

    if (this->mcClientThread) {this->mcClientThread->wait();}
    this->mcClientThread = std::make_shared<WorkThread>(this);
    connect(this->mcClientThread.get() , &WorkThread::titleMsg , this , &Widget::on_titleMsg);
    connect(this->mcClientThread.get() , &WorkThread::btnUsed , this , &Widget::on_btnUsed);
    connect(this->mcClientThread.get() , &WorkThread::btnTitle , this , &Widget::on_btnTitle);
    connect(this->mcClientThread.get() , &WorkThread::min , this , &Widget::on_min);
    connect(this->mcClientThread.get() , &WorkThread::normal , this , &Widget::on_normal);
    this->mcClientThread->start();

}

void Widget::on_titleMsg(QString msg)
{
    QMessageBox::information(this , "提示" ,  msg);
}

void Widget::on_btnUsed(bool used)
{
    ui->game->setEnabled(used);
}

void Widget::on_btnTitle(QString msg)
{
    ui->game->setText(msg);
}

void Widget::on_setting_clicked(bool)
{
    SettingDlg dlg(this);
    dlg.loadMCPEPath();
    dlg.exec();
    this->config_ = dlg.config;
}

void Widget::on_game_clicked(bool) {

    //结束进行
    if (this->eulo) {this->eulo->terminate();}
    if (this->minecraft) {this->minecraft->terminate();}
    //检查客户端是否存在
    this->gamePath = std::make_shared<std::filesystem::path>(config_.gamePath + u"/Minecraft.Windows.exe");
    this->configPath = std::make_shared<std::filesystem::path>( std::filesystem::current_path().generic_u16string() +  u"/netease.cppconfig");

    if (!exists(*gamePath))
    {
        QMessageBox::critical(this , "错误" , "路径不正确");
        QMessageBox::information(this , "提示" , "请配置启动路径！");
        on_setting_clicked(true);
        return;
    }

    //先启动代理
    on_login_clicked(true);

}

void Widget::on_min()
{
    this->showMinimized();
}

void Widget::on_normal()
{
    this->showNormal();
}

void WorkThread::run()
{
    while(!widget_->dead && widget_->eulo->running() || MsgRegister::obj().getMsg())
    {
        auto msg = MsgRegister::obj().getMsg();
        if (!msg)
        {
            std::this_thread::sleep_for(std::chrono::microseconds(200));
            continue;
        }

        if(!msg->started)
        {
            std::string msgStr;
            msgStr.resize(msg->errorMsgLen);
            memcpy(msgStr.data() , msg->errorMsg.data() , msg->errorMsgLen);
            emit titleMsg(QString(msgStr.c_str()));
            emit btnUsed(true);
            emit btnTitle("启动游戏");
            MsgRegister::obj().setMsg(nullptr);
            return;
        }

        emit titleMsg ("启动成功！");
        emit btnTitle("服务已启动");
        emit btnUsed(true);
        MsgRegister::obj().setMsg(nullptr);

        // 开始启动客户端
        try
        {
            widget_->minecraft = std::make_shared<boost::process::child>(widget_->gamePath->generic_wstring() ,
                L"config=" + widget_->configPath->generic_wstring()
            ,boost::process::windows::hide);
        }
        catch (const std::exception &e)
        {
            emit titleMsg("无法找到启动 Minecraft.Windows 进程 请配置正确路径");
        }
    }
    emit btnUsed(true);
    emit btnTitle("启动游戏");
    emit this->normal();


}


