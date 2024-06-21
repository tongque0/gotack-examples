
import time
import json
import math
import socket
import simplejson
import threading
import os
from gui.util import save_record_to_file
class Robot(threading.Thread):
    def __init__(self, dab):
        super().__init__()
        self.daemon = True  # 设置线程为守护线程
        self.dab = dab

    def run(self):
        while self.dab.who == 1:
            self.dab.thinking = True
            self.dab.queue_draw()
            if self.dab.first == 0:
                s0, s1 = self.dab.human, self.dab.robot
            else:
                s1, s0 = self.dab.human, self.dab.robot
            now = 0
            if self.dab.who != self.dab.first:
                now = 1
            h, v = 0, 0
            for move in self.dab.record:
                x, y = move[1], move[2]
                if move[0] == 0:
                    v |= (1 << (y * 6 + x))
                else:
                    h |= (1 << (x * 6 + y))

            algorithm = "quctann"
            print(self.dab.record)
            timeout = int(1 + self.dab.timeout) * 1000
            try:
                s = socket.create_connection(("localhost", 12345))
                arg = {
                    "id": int(time.time()),
                    "method": "Server.MakeMove",
                    "params": [{
                        "Algorithm": algorithm,
                        "Board": {"H": h, "V": v, "S": [s0, s1], "Now": now, "Turn": self.dab.turn},
                        "Timeout": timeout
                    }]
                }
                data = simplejson.dumps(arg).encode()
                s.sendall(data)
                response = s.recv(1024).decode()
                s.close()
                res = simplejson.loads(response)
                print(res)
                ms = (res["result"]["H"], res["result"]["V"])
                moves = []
                for i in range(2):
                    for n in range(30):
                        if ((1 << n) & ms[i]) != 0:
                            moves.append(self.dab.num2move(((1 << n) | (i << 31)), 1, 1))
                while len(moves) > 1:
                    for m in moves:
                        if not self.dab.change(m):
                            self.dab.move(m)
                            moves.remove(m)
                            break
                if moves:
                    self.dab.move(moves[0])
                    moves.remove(moves[0])
            except Exception as e:
                print("Robot failed:", e)
            finally:
                self.dab.thinking = False
                self.dab.queue_draw()
                save_record_to_file(self.dab.record,s0,s1,"record.txt")  # 保存棋谱记录到文件
