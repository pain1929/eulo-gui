#ifndef SETTINGDLG_H
#define SETTINGDLG_H

#include <QDialog>

struct GameConfig
{
    std::u16string gamePath;

    size_t toData (uint8_t * data)
    {
        auto gamePathLen = gamePath.length() * 2; //单位字节
        memcpy(data , &gamePathLen , 4);
        memcpy(data + 4 , gamePath.data() ,gamePathLen );

        return 4 + gamePathLen;
    }

    void load (uint8_t * data)
    {
        int32_t len{};
        memcpy(&len , data , 4);
        gamePath.resize(len / 2);
        memcpy(gamePath.data() , data +4 , len);
    }


};


namespace Ui {
class SettingDlg;
}

class SettingDlg : public QDialog
{
    Q_OBJECT

public:
    GameConfig config;
    explicit SettingDlg(QWidget *parent = nullptr);

    ~SettingDlg();

public slots:
    void on_save_clicked(bool);

private:
    Ui::SettingDlg *ui;
};

#endif // SETTINGDLG_H
