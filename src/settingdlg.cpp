#include "settingdlg.h"
#include "ui_settingdlg.h"
#include <filesystem>
#include <iostream>
#include <QMessageBox>
#include "FileReader.h"

// 递归寻找游戏
static std::shared_ptr<std::u16string> findBEClient(const std::u16string &path) {

    for (const auto & e : std::filesystem::directory_iterator(path) )
    {
        if (!std::filesystem::is_directory(e) && e.path().filename().generic_u16string() == u"Minecraft.Windows.exe")
        {
            return std::make_shared<std::u16string>(e.path().parent_path().generic_u16string());
        } else if (std::filesystem::is_directory(e)){
            auto ret =  findBEClient(e.path().generic_u16string());
            if (ret ) {return ret;}
        }
    }

    return nullptr;

}



SettingDlg::SettingDlg(QWidget *parent) :
    QDialog(parent),
    ui(new Ui::SettingDlg)
{
    ui->setupUi(this);
    setWindowTitle("配置项");
    ui->gamePath->setPlaceholderText("设置游戏启动路径");
    auto buf = FileReader::read("config.data");

    // 读到了配置文件
    if (!buf.empty())
    {
        config.load(buf.data());
        ui->gamePath->setText(QString::fromStdU16String(config.gamePath));
    }


}

void SettingDlg::on_save_clicked(bool)
{
    auto str = ui->gamePath->text().toStdU16String();
    if (std::filesystem::exists(str) && std::filesystem::is_directory(str))
    {
        config.gamePath = str;
        auto buf = new uint8_t[1024];
        auto size = config.toData(buf);
        FileReader::write("config.data" ,buf , size );
        delete [] buf;
        QMessageBox::information(this , "提示" , "保存成功");
        this->reject();
    }else
    {
        QMessageBox::critical(this ,"错误" , "目录不存在");
        return;
    }


}

void SettingDlg::loadMCPEPath() {
    // 没有读到配置文件 寻找默认路径
    static const std::u16string defPath = u"MCLDownload/";
    static const std::list<std::u16string> prefix = {u"c:/" , u"d:/" , u"e:/" , u"f:/" , u"g:/" ,u"h:/" , u"i:/"};

    //如果路径为空自动寻找默认路径
    if (config.gamePath.empty()) {
        // 创建消息框
        QMessageBox msgBox;
        msgBox.setText("是否自动寻找游戏路径？这可能会花费几秒钟。");
        msgBox.setStandardButtons(QMessageBox::Yes | QMessageBox::No | QMessageBox::Cancel);
        msgBox.setDefaultButton(QMessageBox::Yes);

        // 显示消息框并判断按钮
        if (msgBox.exec() != QMessageBox::Yes){
            return;
        }

        for (const auto & p : prefix) {

            if (std::filesystem::exists(config.gamePath)){
                break;
            }


            auto path = p + defPath;

            if (!std::filesystem::exists(path) || !std::filesystem::is_directory(path)){
                continue;
            }

            for (const auto & e : std::filesystem::directory_iterator(path)) {

                //查看目录是否包含MinecraftBE字符串如果没有则略过
                if (e.path().generic_u16string().find(u"MinecraftBE") == std::string::npos ) {
                    continue;
                }

                //在包含MinecraftBE字符串的目录下寻找 游戏客户端
                auto ret = findBEClient(e.path().generic_u16string());

                //成功找到后保存配置
                if (ret){
                    config.gamePath = *ret;
                    ui->gamePath->setText(QString::fromStdU16String(config.gamePath));
                    break;
                }

            }

        }

    }

    //再次检查
    if(config.gamePath.empty()) {
        QMessageBox::information(this , "提示" , "未找到游戏路径 请检查是否安装游戏并曾经游玩过基岩版游戏");
    }

}

SettingDlg::~SettingDlg()
{
    delete ui;
}
