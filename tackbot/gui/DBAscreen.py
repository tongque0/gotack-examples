import gi
import time
import json
import math
import threading
import simplejson
gi.require_version('Gtk', '3.0')
from gi.repository import Gtk, Gdk, GLib
from gui.Robot import Robot

class DAB(Gtk.Window):
    def __init__(self):
        # 调用父类构造函数初始化窗口
        super(DAB, self).__init__()

        # 设置窗口的基本属性
        self.myname, self.version = 'TackDot', '1.0.0'  # 应用程序名称和版本
        self.set_title(self.myname)  # 设置窗口标题
        self.set_resizable(False)  # 禁用窗口大小调整
        self.set_position(Gtk.WindowPosition.CENTER)  # 将窗口居中显示
        self.GridSize = 80  # 设置绘图区域的网格大小
        self.thinking = False  # 标志AI是否正在思考
        self.timeout = 1  # AI思考时间,单位s
        self.algorithm = 'quctann'

        # 创建快捷键加速器组
        agr = Gtk.AccelGroup()
        self.add_accel_group(agr)

        # 创建菜单栏
        self.mb = Gtk.MenuBar()

        # 文件菜单设置
        self.fileitem = Gtk.MenuItem.new_with_label('对局')
        self.filemenu = Gtk.Menu()
        self.fileitem.set_submenu(self.filemenu)
        # 新游戏菜单项
        self.newitem = Gtk.MenuItem.new_with_label('新游戏')
        key, mod = Gtk.accelerator_parse('<Control>N')
        self.newitem.add_accelerator('activate', agr, key, mod, Gtk.AccelFlags.VISIBLE)
        self.newitem.connect('activate', self.on_new_game)
        self.filemenu.append(self.newitem)
        # 人类先手切换菜单项
        self.humanitem = Gtk.CheckMenuItem.new_with_label('人类先手')
        self.humanitem.set_sensitive(False)
        self.humanitem.set_active(True)
        self.humanitem.connect('activate', self.on_human_first)
        self.filemenu.append(self.humanitem)
        # 机器人先手切换菜单项
        self.robotitem = Gtk.CheckMenuItem.new_with_label('机器人先手')
        self.robotitem.set_active(False)
        self.robotitem.connect('activate', self.on_robot_first)
        self.filemenu.append(self.robotitem)
        # 添加分隔符
        self.filemenu.append(Gtk.SeparatorMenuItem())
        # 撤销菜单项
        self.undoitem = Gtk.MenuItem.new_with_label('悔棋')
        self.undoitem.set_sensitive(False)
        key, mod = Gtk.accelerator_parse('<Control>X')
        self.undoitem.add_accelerator('activate', agr, key, mod, Gtk.AccelFlags.VISIBLE)
        self.undoitem.connect('activate', self.on_undo)
        self.filemenu.append(self.undoitem)
        # 重做菜单项
        self.redoitem = Gtk.MenuItem.new_with_label('撤销悔棋')
        self.redoitem.set_sensitive(False)
        key, mod = Gtk.accelerator_parse('<Control>Z')
        self.redoitem.add_accelerator('activate', agr, key, mod, Gtk.AccelFlags.VISIBLE)
        self.redoitem.connect('activate', self.on_redo)
        self.filemenu.append(self.redoitem)
        # 退出菜单项
        self.filemenu.append(Gtk.SeparatorMenuItem())
        self.exititem = Gtk.MenuItem.new_with_label('退出')
        key, mod = Gtk.accelerator_parse('<Control>Q')
        self.exititem.add_accelerator('activate', agr, key, mod, Gtk.AccelFlags.VISIBLE)
        self.exititem.connect('activate', Gtk.main_quit)
        self.filemenu.append(self.exititem)
        # 上膛
        self.mb.append(self.fileitem)

                # 时间设置菜单
        time_menu = Gtk.Menu()
        time_menu_item = Gtk.MenuItem.new_with_label("时间设置")
        time_menu_item.set_submenu(time_menu)
        group = None
        for seconds in [1, 5, 10, 20, 30]:
            menu_item = Gtk.RadioMenuItem.new_with_label(group, f"{seconds} 秒")
            menu_item.connect("toggled", self.on_timeout_change, seconds)
            if group is None:
                group = menu_item.get_group()
            menu_item.set_active(self.timeout == seconds)
            time_menu.append(menu_item)
        self.mb.append(time_menu_item)

        # 算法选择菜单
        algo_menu = Gtk.Menu()
        algo_menu_item = Gtk.MenuItem.new_with_label("算法选择")
        algo_menu_item.set_submenu(algo_menu)
        algo_group = None
        for algo in ["Alpha-Beta", "MCTS", "UCT", "quctann"]:
            algo_item = Gtk.RadioMenuItem.new_with_label(algo_group, algo)
            algo_item.connect("toggled", self.on_algorithm_change, algo)
            if algo_group is None:
                algo_group = algo_item.get_group()
            algo_item.set_active(self.algorithm == algo)
            algo_menu.append(algo_item)
        self.mb.append(algo_menu_item)

        # 创建绘图区域
        self.darea = Gtk.DrawingArea()
        # 根据新的GridSize调整绘图区域的尺寸
        self.darea.set_size_request(self.GridSize * 7, self.GridSize * 7)

        # 创建状态栏
        self.statusbar = Gtk.Statusbar()
        context_id = self.statusbar.get_context_id("info")
        self.statusbar.push(context_id, ' 0m00s 移动次数: 0')
        align = Gtk.Alignment(xalign=0.9, yalign=0, xscale=0, yscale=1)
        self.scorelabel = Gtk.Label(label='人类: 0, 机器人: 0')
        align.add(self.scorelabel)
        self.statusbar.pack_end(align, False, False, 0)

        # 创建主布局容器
        self.vbox = Gtk.Box(orientation=Gtk.Orientation.VERTICAL, spacing=2)
        self.vbox.pack_start(self.mb, False, False, 0)
        self.vbox.pack_start(self.darea, False, False, 0)
        self.vbox.pack_start(self.statusbar, False, False, 0)
        self.add(self.vbox)

        # 连接窗口销毁事件
        self.connect('destroy', Gtk.main_quit)
        self.darea.connect('draw', self.draw)
        self.darea.set_events(Gdk.EventMask.ALL_EVENTS_MASK)
        self.darea.connect('button-press-event', self.on_darea_button_press)
        self.darea.connect('button-release-event', self.on_darea_button_release)
        self.darea.connect('motion-notify-event', self.on_darea_motion_notify)
        self.darea.connect('leave-notify-event', self.on_darea_leave_notify)
        self.init_board()

        # 显示所有组件
        self.show_all()

    # 绘制窗口内容
    def draw(self, widget, cr):
        self.undoitem.set_sensitive(self.moves > 0)
        self.redoitem.set_sensitive(self.moves < len(self.record))
        self.scorelabel.set_text('Timeout: %ds, 人类得分: %d, AI得分: %d' % (self.timeout,  self.human, self.robot))
        self.draw_board(cr)


    # 初始化棋盘
    def init_board(self, first=0):
        self.degree = [[4] * 5 for i in range(5)]
        self.belong = [[-1] * 5 for i in range(5)]
        self.hexist = [[0] * 5 for i in range(6)]
        self.vexist = [[0] * 6 for i in range(5)]
        self.human, self.robot = 0, 0
        self.first, self.who = first, first
        self.moves, self.record = 0, []
        self.turn = 0
        self.cursor = (-1, -1, -1)
        self.begtime = time.time()
        self.update_time_elapse()
        self.queue_draw()
        if self.who != 0:
            Robot(self).start()

    # 更新经过的时间
    def update_time_elapse(self):
        self.statusbar.remove_all(1)
        if self.thinking:
            tips = 'thinking...'
        else:
            tips = f'Turn(s): {self.turn}'
        self.statusbar.push(
            1, f'{int(time.time() - self.begtime) // 60}m{int(time.time() - self.begtime) % 60:02d}s {tips}')
        if self.human + self.robot < 25:
            GLib.timeout_add(1000, self.update_time_elapse)
        self.queue_draw()

    def on_new_game(self, widget):
        if self.thinking:
            self.thinking=False
        self.init_board(self.first)

    def on_human_first(self, widget):
        if self.thinking:
            self.thinking=False
        if self.first == 1:
            self.robotitem.set_active(False)
            self.robotitem.set_sensitive(True)
            self.humanitem.set_sensitive(False)
            self.first = 0
            self.init_board(0)

    def on_robot_first(self, widget):
        if self.thinking:
            self.thinking=False
        if self.first == 0:
            self.humanitem.set_active(False)
            self.humanitem.set_sensitive(True)
            self.robotitem.set_sensitive(False)
            self.first = 1
            self.init_board(1)

    def on_undo(self, widget):
        if self.thinking:
            return
        for i in range(self.moves)[::-1]:
            if self.record[i][3] == 0:
                self.moves = i
                break
        self.degree = [[4] * 5 for i in range(5)]
        self.belong = [[-1] * 5 for i in range(5)]
        self.hexist = [[0] * 5 for i in range(6)]
        self.vexist = [[0] * 6 for i in range(5)]
        self.human, self.robot = 0, 0
        self.who = self.first
        self.turn = 0
        for move in self.record[:self.moves]:
            self.move(move, False)
        self.queue_draw()

    def on_redo(self, widget):
        if self.thinking:
            return
        cnt = 0
        for move in self.record[self.moves:]:
            if move[3] == 0:
                cnt += 1
            if cnt > 1:
                break
            self.move(move, False)
            self.moves += 1
        self.queue_draw()

    def on_darea_button_press(self, widget, event):
        if self.thinking:
            return

    def on_darea_button_release(self, widget, event):
        if self.thinking:
            return
        if self.who == 0:  # human's turn
            move = self.xy2move(event.x, event.y, self.who)
            if move[0] < 0:
                return
            self.move(move)
            self.cursor = (-1, -1, -1)
            self.queue_draw()
            if self.who != 0:
                Robot(self).start()

    def on_darea_motion_notify(self, widget, event):
        cursor = self.cursor
        self.cursor = self.xy2move(event.x, event.y, self.who)[:3]
        if self.cursor[0] == 0 and self.hexist[self.cursor[1]][self.cursor[2]] != 0:
            self.cursor = (-1, -1, -1)
        if self.cursor[0] == 1 and self.vexist[self.cursor[1]][self.cursor[2]] != 0:
            self.cursor = (-1, -1, -1)
        if self.cursor != cursor:
            self.queue_draw()
        if self.human + self.robot == 25 and self.who == -1:
            self.who = -2
            if self.human > self.robot:
                msg_text = 'You win.\nGood job!'
            else:
                msg_text = 'I win!\nHahaha'
            self.show_message_dialog(msg_text)

    def on_darea_leave_notify(self, widget, event):
        self.cursor = (-1, -1, -1)
        self.queue_draw()

    def on_ai_difficulty_change(self, widget):
       print("AI difficulty change requested.")

    def show_message_dialog(self, text):
        dialog = Gtk.MessageDialog(
            transient_for=self,
            flags=0,
            message_type=Gtk.MessageType.INFO,
            buttons=Gtk.ButtonsType.OK,
            text=text
        )
        dialog.run()
        dialog.destroy()

    def draw_board(self, cr):
        orgx, orgy = self.GridSize, self.GridSize
        w, h = self.darea.get_allocated_width(), self.darea.get_allocated_height()
        size = (min(w, h) - (orgx + orgy)) / 5

        # 设置背景
        if not self.thinking:
            cr.set_source_rgb(240 / 255.0, 248 / 255.0, 255 / 255.0)  # 淡蓝色
        else:
            cr.set_source_rgb(229 / 255.0, 252 / 255.0, 163 / 255.0)  # 淡绿色


        cr.paint()
        # 绘制棋盘格
        for i in range(5):
            for j in range(5):
                if self.belong[i][j] == 0:  # 人类占领
                    cr.rectangle(orgx + size * j, orgy + size * i, size, size)
                    cr.set_source_rgb(250 / 255.0, 199 / 255.0, 199 / 255.0)
                elif self.belong[i][j] == 1:  # 机器人占领
                    cr.rectangle(orgx + size * j, orgy + size * i, size, size)
                    cr.set_source_rgb(0 , 255, 255)
                cr.fill()

        # 绘制线条
        cr.set_line_width(2)
        for i in range(6):
            for j in range(5):
                if self.hexist[i][j] == 1:
                    cr.set_source_rgb(1, 0, 0)  # 红色
                    cr.move_to(orgx + size * j, orgy + size * i)
                    cr.line_to(orgx + size * (j + 1), orgy + size * i)
                elif self.hexist[i][j] == 2:
                    cr.set_source_rgb(0, 0, 1)  # 蓝色
                    cr.move_to(orgx + size * j, orgy + size * i)
                    cr.line_to(orgx + size * (j + 1), orgy + size * i)
                cr.stroke()

        for i in range(5):
            for j in range(6):
                if self.vexist[i][j] == 1:
                    cr.set_source_rgb(1, 0, 0)  # 红色
                    cr.move_to(orgx + size * j, orgy + size * i)
                    cr.line_to(orgx + size * j, orgy + size * (i + 1))
                elif self.vexist[i][j] == 2:
                    cr.set_source_rgb(0, 0, 1)  # 蓝色
                    cr.move_to(orgx + size * j, orgy + size * i)
                    cr.line_to(orgx + size * j, orgy + size * (i + 1))
                cr.stroke()

         # 绘制游戏光标
        if self.cursor != (-1, -1, -1):
            cr.set_line_width(4)
            cr.set_source_rgb(1, 215 / 255.0, 0)  # 金黄色
        if self.cursor[0] == 0:  # 水平
            cr.move_to(orgx + size * self.cursor[2], orgy + size * self.cursor[1])
            cr.line_to(orgx + size * (self.cursor[2] + 1), orgy + size * self.cursor[1])
        else:  # 垂直
            cr.move_to(orgx + size * self.cursor[2], orgy + size * self.cursor[1])
            cr.line_to(orgx + size * self.cursor[2], orgy + size * (self.cursor[1] + 1))
        cr.stroke()


        # 绘制点
        cr.set_line_width(6)
        cr.set_source_rgb(0, 0, 0)
        for i in range(6):
            for j in range(6):
                cr.arc(orgx + size * i, orgy + size * j, 1, 0, 2 * math.pi)
                cr.stroke()


         # 读取 self.record 的最后一个元素，并且确保当前有动作并且未撤销至该动作前
        if self.record and self.moves == len(self.record):
            last_move = self.record[-1]
            direction, row, col, player = last_move

            # 设置颜色为黑色
            cr.set_source_rgb(0, 0, 0)  # 黑色

            cr.set_line_width(4)  # 设置线宽为4，使线条更明显

            # 根据最后一个记录的方向和位置绘制线条
            if direction == 0:  # 水平线条
                cr.move_to(orgx + size * col, orgy + size * row)
                cr.line_to(orgx + size * (col + 1), orgy + size * row)
            else:  # 垂直线条
                cr.move_to(orgx + size * col, orgy + size * row)
                cr.line_to(orgx + size * col, orgy + size * (row + 1))

            cr.stroke()  # 描绘线条

    def xy2move(self, x, y, who):
        x, y = int(x), int(y)
        orgx, orgy = self.GridSize, self.GridSize
        w, h = self.darea.get_allocated_width(), self.darea.get_allocated_height()
        size = int((min(w, h) - (orgx + orgy)) / 5)

        # 计算点击位置相对于棋盘原点的坐标
        adjusted_x = x - orgx
        adjusted_y = y - orgy

        # 计算对应的棋盘格索引
        grid_x = adjusted_x // size
        grid_y = adjusted_y // size

        # 检查点击是否在棋盘有效区域外
        if (adjusted_x < 0 or adjusted_y < 0 or
            grid_x > 5 or grid_y > 5 or
            # (adjusted_x < size and adjusted_y < size) or
            (adjusted_x < size and grid_y > 5) or
            (grid_x > 5 and adjusted_y < size) or
            (grid_x > 5 and grid_y > 5) or
            (x>size*6 and y>size*6) ):
            return (-1, -1, -1, who)

        if x < size:
            return (1, y // size - 1, 0, who)
        if x >= size * (5 + 1):
            return (1, y // size - 1, 5, who)
        if y < size:
            return (0, 0, x // size - 1, who)
        if y >= size * (5 + 1):
            return (0, 5, x // size - 1, who)

        zero = 1e-5
        x1, y1 = float(x // size * size), float(y // size * size)
        x2, y2 = float(x1 + size), y1
        x3, y3 = x1, y1 + size
        x4, y4 = x2, y3
        p1x, p1y = x4 - x1, y4 - y1
        p2x, p2y = x - x1, y - y1
        p3x, p3y = x2 - x3, y2 - y3
        p4x, p4y = x - x3, y - y3
        c1 = (p1x * p2y) - (p1y * p2x)
        c2 = (p3x * p4y) - (p3y * p4x)

        if (c1 >= -zero and c1 <= zero) or (c1 >= -zero and c1 <= zero):
            return (-1, -1, -1, who)
        if c1 < 0 and c2 < 0:
            return (0, y // size - 1, x // size - 1, who)
        if c1 < 0 and c2 > 0:
            return (1, y // size - 1, x // size, who)
        if c1 > 0 and c2 > 0:
            return (0, y // size, x // size - 1, who)
        if c1 > 0 and c2 < 0:
            return (1, y // size - 1, x // size - 1, who)
        return (-1, -1, -1, who)

    def change(self, move):
        x, y = move[1], move[2]
        if move[0] == 0:  # horizon
            if x > 0 and self.degree[x-1][y] == 1:
                return False
            if x < 5 and self.degree[x][y] == 1:
                return False
        else:  # vertical
            if y > 0 and self.degree[x][y-1] == 1:
                return False
            if y < 5 and self.degree[x][y] == 1:
                return False
        return True

    def num2move(self, value, who, step=-1):
        ty = 0 if (value & (1 << 31)) != 0 else 1
        x, y = -1, -1
        for i in range(5)[::step]:
            for j in range(6)[::step]:
                if (value & 1) == 1:
                    x, y = (j, i) if ty == 0 else (i, j)
                    break
                value >>= 1
            if x != -1:
                break
        return (ty, x, y, who)

    def move(self, move, change_record=True):
        if len(move) != 4:
            move = self.num2move(move[0], move[1])
        flag, x, y, who = False, move[1], move[2], move[3]
        if move[0] == 0:  # horizon
            if self.hexist[x][y]:
                return
            self.hexist[x][y] = who + 1
            if x > 0:
                self.degree[x-1][y] -= 1
                if self.degree[x-1][y] == 0:
                    self.belong[x-1][y] = who
                    if who == 0: self.human += 1
                    else: self.robot += 1
                    flag = True
            if x < 5:
                self.degree[x][y] -= 1
                if self.degree[x][y] == 0:
                    self.belong[x][y] = who
                    if who == 0: self.human += 1
                    else: self.robot += 1
                    flag = True
        else:  # vertical
            if self.vexist[x][y]:
                return
            self.vexist[x][y] = who + 1
            if y > 0:
                self.degree[x][y-1] -= 1
                if self.degree[x][y-1] == 0:
                    self.belong[x][y-1] = who
                    if who == 0: self.human += 1
                    else: self.robot += 1
                    flag = True
            if y < 5:
                self.degree[x][y] -= 1
                if self.degree[x][y] == 0:
                    self.belong[x][y] = who
                    if who == 0: self.human += 1
                    else: self.robot += 1
                    flag = True
        if change_record:
            del self.record[self.moves:]
            self.record.append(move)
            self.moves += 1
        if self.human + self.robot == 25:
            self.who = -1
            self.queue_draw()
            return
        if not flag:
            self.who ^= 1
            self.turn += 1

    def on_timeout_change(self, menu_item, seconds):
        if menu_item.get_active():
            self.timeout = seconds
            print(f"AI思考时间已设置为 {self.timeout} 秒")

    def on_algorithm_change(self, menu_item, algo):
        if menu_item.get_active():
            self.algorithm = algo
            print(f"AI算法已设置为 {self.algorithm}")
