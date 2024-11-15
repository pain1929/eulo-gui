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
