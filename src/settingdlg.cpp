#include "settingdlg.h"
#include "ui_settingdlg.h"
#include <filesystem>
#include <QMessageBox>
#include "FileReader.h"
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

    }else {
        // 没有读到配置文件 寻找默认路径
        static const std::u16string defPath = u"MCLDownload/MinecraftBENetease/windowsmc";
        static const std::list<std::u16string> prefix = {u"c:/" , u"d:/" , u"e:/" , u"f:/" , u"g:/" ,u"h:/" , u"i:/"};

        //如果路径为空自动寻找默认路径
        if (config.gamePath.empty()) {

            for (const auto & p : prefix) {
                auto path = p + defPath + u"/Minecraft.Windows.exe";
                if (std::filesystem::exists(path)) {
                    config.gamePath = p + defPath;
                    break;
                }
            }

        }
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

SettingDlg::~SettingDlg()
{
    delete ui;
}
