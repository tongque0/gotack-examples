
import json
import os


def move_to_str(move):
    """将棋谱中的单步记录转换为指定格式的字符串"""
    direction = 'h' if move[0] == 0 else 'v'
    column = chr(ord('a') + move[2])  # 将列数转换为字母
    if direction == 'h':
        row = str(6 - move[1])  # 将行数转换为对应的棋盘行数
    else:
        row = str(6 - move[1]-1)
    player = 'r' if move[3] == 0 else 'b'
    return f"{player}({column}{row},{direction})"

def convert_record_to_json(record, s0, s1):
    """将记录和分数转换为指定格式的JSON对象"""
    game_moves = [{"piece": move_to_str(move)} for move in record]
    return {
        "R": f"B{s0}",
        "B": f"B{s1}",
        "winner": "R" if s0 > s1 else "B",
        "RScore": s0,
        "BScore": s1,
        "Date": "2023-11-3",  # 可以使用当前日期
        "Event": "tack",
        "game": game_moves
    }

def save_record_to_file(record, s0, s1, filename, folder='record'):
    """将记录转换为JSON并保存到指定文件夹中，使用GBK编码"""
    try:
        # 确保目标文件夹存在
        if not os.path.exists(folder):
            os.makedirs(folder)
        # 构建完整的文件路径
        file_path = os.path.join(folder, filename)
        print(file_path)
        # 将记录转换为JSON格式
        json_data = convert_record_to_json(record, s0, s1)

        # 保存文件
        with open(file_path, 'w', encoding='gbk') as file:
            json.dump(json_data, file, indent=4, ensure_ascii=False)

        print(f"Record saved to {file_path}")
    except Exception as e:
        print('Save record failed:', e)
