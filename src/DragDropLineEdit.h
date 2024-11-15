#ifndef READFILEEDIT
#define READFILEEDIT

#include <QLineEdit>
#include <QDragEnterEvent>
#include <QDropEvent>
#include <QFile>
#include <QTextStream>
#include <QMimeData>
#include <QMessageBox>

class DragDropLineEdit : public QLineEdit {
    Q_OBJECT

public:
    explicit DragDropLineEdit(QWidget *parent = nullptr) : QLineEdit(parent) {
        // 启用拖拽功能
        setAcceptDrops(true);
    }

protected:
    // 重写 dragEnterEvent，处理拖拽进入事件
    void dragEnterEvent(QDragEnterEvent *event) override {
        // 只有当拖拽的是文件时，才允许接受
        if (event->mimeData()->hasUrls()) {
            event->acceptProposedAction();
        } else {
            event->ignore();
        }
    }

    // 重写 dropEvent，处理文件被放下时的事件
    void dropEvent(QDropEvent *event) override {
        const QMimeData *mimeData = event->mimeData();
        if (mimeData->hasUrls()) {
            // 获取拖放的文件路径
            QList<QUrl> urls = mimeData->urls();
            if (!urls.isEmpty()) {
                QString filePath = urls.first().toLocalFile();
                readFileContent(filePath);
            }
        }
    }

private:
    // 读取文件内容并将其设置到 QLineEdit 中
    void readFileContent(const QString &filePath) {
        QFile file(filePath);
        if (file.open(QIODevice::ReadOnly | QIODevice::Text)) {
            QTextStream in(&file);
            QString fileContent = in.readAll();
            file.close();

            // 将文件内容设置到 QLineEdit 中
            setText(fileContent);
        } else {
            // 读取文件失败，显示错误消息
            QMessageBox::warning(this, "Error", "Failed to open the file.");
        }
    }
};
#endif
