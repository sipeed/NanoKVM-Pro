HID_KEYCODES = {
    # 基础按键
    'A': 0x04, 'B': 0x05, 'C': 0x06, 'D': 0x07, 'E': 0x08,
    'F': 0x09, 'G': 0x0A, 'H': 0x0B, 'I': 0x0C, 'J': 0x0D,
    'K': 0x0E, 'L': 0x0F, 'M': 0x10, 'N': 0x11, 'O': 0x12,
    'P': 0x13, 'Q': 0x14, 'R': 0x15, 'S': 0x16, 'T': 0x17,
    'U': 0x18, 'V': 0x19, 'W': 0x1A, 'X': 0x1B, 'Y': 0x1C, 'Z': 0x1D,
    
    # 数字键
    '1': 0x1E, '2': 0x1F, '3': 0x20, '4': 0x21, '5': 0x22,
    '6': 0x23, '7': 0x24, '8': 0x25, '9': 0x26, '0': 0x27,

    # 符号键
    '`':0x35, '-': 0x2d, '=': 0x2e, '[': 0x2f, ']': 0x30,
    '\\': 0x31, ';': 0x33, '\'': 0x34, ',': 0x36, '.': 0x37, '/': 0x38,
    
    # 特殊键
    'ENTER': 0x28, 'ESCAPE': 0x29, 'BACKSPACE': 0x2A, 'TAB': 0x2B,
    'SPACE': 0x2C, 'CAPS_LOCK': 0x39, 'DELETE': 0x4C, 'INSERT': 0x49,
    
    # 功能键 F1-F12
    'F1': 0x3A, 'F2': 0x3B, 'F3': 0x3C, 'F4': 0x3D, 'F5': 0x3E, 'F6': 0x3F,
    'F7': 0x40, 'F8': 0x41, 'F9': 0x42, 'F10': 0x43, 'F11': 0x44, 'F12': 0x45,
    
    # 方向键
    'RIGHT_ARROW': 0x4F, 'LEFT_ARROW': 0x50, 'DOWN_ARROW': 0x51, 'UP_ARROW': 0x52,
    
    # 修饰键
    'LEFT_CONTROL': 0xE0, 'LEFT_SHIFT': 0xE1, 'LEFT_ALT': 0xE2, 'LEFT_GUI': 0xE3,
    'RIGHT_CONTROL': 0xE4, 'RIGHT_SHIFT': 0xE5, 'RIGHT_ALT': 0xE6, 'RIGHT_GUI': 0xE7,
}

# 消费者控制键码（媒体键）
CONSUMER_CONTROL_CODES = {
    # 媒体控制
    'PLAY_PAUSE': 0xCD,
    'SCAN_NEXT_TRACK': 0xB5,
    'SCAN_PREVIOUS_TRACK': 0xB6,
    'STOP': 0xB7,
    'FAST_FORWARD': 0xB3,
    'REWIND': 0xB4,
    'RECORD': 0xB2,
    
    # 音量控制
    'VOLUME_UP': 0xE9,
    'VOLUME_DOWN': 0xEA,
    'MUTE': 0xE2,
    
    # 电源
    'POWER': 0x30,
    'SLEEP': 0x32,
    'WAKE': 0x33,
    
    # 浏览器/应用
    'HOME': 0x223,
    'BACK': 0x224,
    'FORWARD': 0x225,
    'REFRESH': 0x227,
    'BOOKMARKS': 0x22A,
    
    # 系统
    'CALCULATOR': 0x192,
    'EMAIL': 0x18A,
    'BROWSER': 0x196,
    'FILE_EXPLORER': 0x194,
}

class HIDKeyboardHelper:
    """HID键盘操作助手类"""
    
    def __init__(self):
        self.keycodes = HID_KEYCODES
        self.consumer_codes = CONSUMER_CONTROL_CODES
        self.modifier_map = {
            'CTRL': 0x01, 'SHIFT': 0x02, 'ALT': 0x04, 'GUI': 0x08,
            'CONTROL':0x01, 'LEFT_CONTROL':0x01, 'RIGHT_CONTROL': 0x10,
            'LEFT_CTRL': 0x01, 'LEFT_SHIFT': 0x02, 'LEFT_ALT': 0x04, 'LEFT_GUI': 0x08,
            'RIGHT_CTRL': 0x10, 'RIGHT_SHIFT': 0x20, 'RIGHT_ALT': 0x40, 'RIGHT_GUI': 0x80
        }
    
    def get_keycode(self, key_name):
        """获取按键对应的HID键码"""
        key_name = key_name.upper()
        if key_name in self.keycodes:
            return self.keycodes[key_name]
        elif key_name in self.consumer_codes:
            return self.consumer_codes[key_name]
        else:
            raise ValueError(f"未找到按键: {key_name}")
    
    def parse_key_combination(self, key_combo):
        """
        解析按键组合，如 'CTRL+SHIFT+A', 'F11', 'VOLUME_UP'
        返回: (modifier_mask, keycode, is_consumer_key)
        """
        keys = [k.strip().upper() for k in key_combo.split('+')]
        modifier_mask = 0
        main_key = None
        is_consumer_key = False
        
        err = 0
        for key in keys:
            if key in self.modifier_map:
                modifier_mask |= self.modifier_map[key]
            elif key in self.keycodes:
                main_key = self.keycodes[key]
            elif key in self.consumer_codes:
                main_key = self.consumer_codes[key]
                is_consumer_key = True
            else:
                err = 1
                main_key = key
                break 
        
        return modifier_mask, main_key, is_consumer_key, err 
    
    def generate_hid_sequence(self, key_combo):
        """
        生成HID操作序列
        返回: HID报告列表，每个报告是8字节数组
        """
        modifier_mask, keycode, is_consumer_key, err = self.parse_key_combination(key_combo)
        
        if err == 1:
            return "Error: Unknown key %s in combination"%(keycode)
        else:
            if is_consumer_key:
                # 消费者控制键（媒体键）使用不同的报告格式
                return self._generate_consumer_sequence(keycode)
            else:
                return self._generate_keyboard_sequence(modifier_mask, keycode)
    
    def _generate_keyboard_sequence(self, modifier_mask, keycode):
        """生成标准键盘HID序列"""
        sequence = []
        
        if keycode is None:
            # 仅修饰键
            sequence.append([modifier_mask, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00])
            sequence.append([0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00])
        elif modifier_mask == 0:
            # 仅普通键
            sequence.append([0x00, 0x00, keycode, 0x00, 0x00, 0x00, 0x00, 0x00])
            sequence.append([0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00])
        else:
            # 组合键
            sequence.append([modifier_mask, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00])  # 按下修饰键
            sequence.append([modifier_mask, 0x00, keycode, 0x00, 0x00, 0x00, 0x00, 0x00])  # 按下主键
            sequence.append([modifier_mask, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00])  # 松开主键
            sequence.append([0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00])  # 松开所有键
        
        return sequence
    
    def _generate_consumer_sequence(self, keycode):
        """生成消费者控制键（媒体键）HID序列"""
        # 消费者控制键使用2字节报告
        return [
            [keycode & 0xFF, (keycode >> 8) & 0xFF],  # 按下
            [0x00, 0x00]  # 松开
        ]
    
    def print_sequence(self, key_combo):
        """打印按键组合的HID序列"""
        try:
            modifier_mask, keycode, is_consumer_key, err = self.parse_key_combination(key_combo)
            if err == 1:
                print(f"错误: 未知按键 '{keycode}'")
            sequence = self.generate_hid_sequence(key_combo)
            
            print(f"\n按键组合: {key_combo}")
            print(f"解析结果: modifier=0x{modifier_mask:02X}, keycode=0x{keycode or 0:02X}, consumer_key={is_consumer_key}")
            print("HID操作序列:")
            
            for i, report in enumerate(sequence):
                if is_consumer_key:
                    print(f"  {i+1}. 消费者控制报告: {report}")
                else:
                    print(f"  {i+1}. 键盘报告: {report}")
                    
        except ValueError as e:
            print(f"错误: {e}")

# ===== 使用示例 =====
def main():
    hid_helper = HIDKeyboardHelper()
    
    # 测试各种按键组合
    test_keys = [
        'a',              # 普通小写字母
        'A',              # 普通大写字母
        '`',              # 可见字符
        '~',              # 可见字符shift
        'F11',            # 功能键
        'CTRL+C',         # 组合键
        'CTRL+SHIFT+A',   # 多修饰键组合
        'VOLUME_UP',      # 媒体键
        'PLAY_PAUSE',     # 媒体控制
        'ALT+F4',         # 经典组合键
        'GUI+L',          # Windows锁屏
        "SHIFT+`"
    ]
    
    for key in test_keys:
        hid_helper.print_sequence(key)

if __name__ == "__main__":
    main()