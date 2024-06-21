# gi模块不需要引入，需要下载gtk链接库
import gi
gi.require_version('Gtk', '3.0')

from gi.repository import Gtk

from gui.DBAscreen import DAB

# 这是Python应用程序的入口点
if __name__ == '__main__':
    # 创建DAB类的实例，这个实例是应用程序的主窗口
    win = DAB()

    # 当窗口被关闭时，Gtk.main_quit将被调用，以便干净地退出GTK主事件循环
    win.connect("destroy", Gtk.main_quit)

    # 显示窗口及其所有子组件
    win.show_all()

    Gtk.main()
