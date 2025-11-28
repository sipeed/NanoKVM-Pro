import os,sys,subprocess,time
from datetime import datetime
import requests, base64, io
from hid_keyboard import HIDKeyboardHelper
from PIL import Image


HID_KB_DEVICE    = '/dev/hidg0'
HID_MOUSE_DEVICE = '/dev/hidg1'
HID_TP_DEVICE    = '/dev/hidg2'
if not os.path.exists(HID_MOUSE_DEVICE):
    print(f"警告：设备文件 {HID_MOUSE_DEVICE} 不存在。请确认 HID Gadget 已正确配置。")
hid_helper = HIDKeyboardHelper()


def send_tp_hid(x,y, w,h):
    print("--> move to %d,%d"%(int(x*w),int(y*h)))
    x = int(x*0x7fff)
    y = int(y*0x7fff)
    pos_str = "00%02x%02x%02x%02x00"%(x%256, x//256, y%256, y//256)  #以下处理为了防止出现0x0a特殊字符
    command =  f'echo "{pos_str}" | xxd -r -p > {HID_TP_DEVICE}' 
    print(command)
    subprocess.run(command, shell=True, check=True)

# [buttons, x_movement, y_movement, wheel_movement]
def send_mouse_hid_full(btn, dx, dy, wheel):
    pos_str = "%02x%02x%02x%02x"%(btn, dx&0xff, dy&0xff, wheel)  #以下处理为了防止出现0x0a特殊字符
    command = f'echo "{pos_str}" | xxd -r -p > {HID_MOUSE_DEVICE}' 
    #print(command)
    subprocess.run(command, shell=True, check=True)

def send_mouse_hid(btn, wheel):
    pos_str = "%02x0000%02x"%(btn, wheel)  #以下处理为了防止出现0x0a特殊字符
    command = f'echo "{pos_str}" | xxd -r -p > {HID_MOUSE_DEVICE}' 
    print(command)
    subprocess.run(command, shell=True, check=True)


def send_keyboard_hid(hid):
    hid_str = "".join(["%02x" % b for b in hid])
    command =  f'echo "{hid_str}" | xxd -r -p > {HID_KB_DEVICE}' 
    print(command)
    subprocess.run(command, shell=True, check=True)

def encode_image(path):
    with open(path, "rb") as f:
        return base64.b64encode(f.read()).decode("utf-8")

def mouse(x, y, btn, w,h):
    # move mouse
    if x!=0 or y!=0:
        send_tp_hid(x, y, w,h)
        time.sleep(0.1)
    # btn action
    if btn == 1 or btn == 2:
        val = btn  #0x01 left, 0x02 right
        send_mouse_hid(btn, 0)
        time.sleep(0.05)
        send_mouse_hid(0, 0)
        time.sleep(0.5)
    elif btn == 3 or btn ==4:
        val = btn-2  #left/right double
        send_mouse_hid(val, 0)
        time.sleep(0.05)
        send_mouse_hid(0, 0)
        time.sleep(0.1) 
        send_mouse_hid(val, 0)
        time.sleep(0.05)
        send_mouse_hid(0, 0)
    elif btn == 5 :
        #scroll down 
        send_mouse_hid(0, 0xfe)
        send_mouse_hid(0, 0xfe)
        send_mouse_hid(0, 0xfe)
        send_mouse_hid(0, 0xfe)
        time.sleep(0.1)
        send_mouse_hid(0, 0)
    elif btn == 6:
        #scroll up 
        send_mouse_hid(0, 0x02)
        send_mouse_hid(0, 0x02)
        send_mouse_hid(0, 0x02)
        send_mouse_hid(0, 0x02)
        time.sleep(0.1)
        send_mouse_hid(0, 0)
    else:
        return -1
    return 0

def drag(x0, y0, x1, y1, w, h):
    # move mouse
    send_tp_hid(x0, y0, w,h)
    time.sleep(0.1)
    # press left btn
    send_mouse_hid(1, 0)
    time.sleep(0.1)
    # move mouse
    dx_rest = int((x1-x0)*w)
    dy_rest = int((y1-y0)*h)
    while dx_rest!=0 or dy_rest!=0:
        dd = 16
        dx = max(-dd, min(dd, dx_rest))
        dy = max(-dd, min(dd, dy_rest))
        x0 = x0+dx
        y0 = y0+dy
        #print("move %d,%d to %d,%d, target %d,%d"%(dx,dy,x0,y0,x1,y1))
        send_mouse_hid_full(1, dx, dy, 0)
        time.sleep(0.01)
        dx_rest -= dx
        dy_rest -= dy
    # release left btn
    send_mouse_hid(0, 0)
    return 0


'''
key_desc is str like this
'a','A','F11','CTRL+C','CTRL+SHIFT+A','ALT+F4','GUI+L'

special key str define like this:
'F1','F2','F3','F4','F5','F6','F7','F8','F9','F10','F11','F12',
'ENTER','ESCAPE','BACKSPACE','TAB','SPACE','CAPS_LOCK',
'RIGHT_ARROW','LEFT_ARROW','DOWN_ARROW','UP_ARROW',
'LEFT_CONTROL','LEFT_SHIFT','LEFT_ALT','LEFT_GUI','RIGHT_CONTROL',
'RIGHT_SHIFT','RIGHT_ALT','RIGHT_GUI',
'''
def keyboard(key_desc):
    sequence = hid_helper.generate_hid_sequence(key_desc)
    if isinstance(sequence, str):  #error str
        return -1, sequence
    else:
        for i, report in enumerate(sequence):
            send_keyboard_hid(report)
            time.sleep(0.02)
        return 0, "OK"


def typestr(content):
    lower_chars = "`1234567890-=qwertyuiop[]\\asdfghjkl;'zxcvbnm,./"
    upper_chars = "~!@#$%^&*()_+QWERTYUIOP{}|ASDFGHJKL:\"ZXCVBNM<>?"
    desc_chars =  "`1234567890-=QWERTYUIOP[]\\ASDFGHJKL;'ZXCVBNM,./"
    print("Type in: ", end="")
    err = 0
    info = "OK"
    for c in content:
        desc_char = ''
        if c in lower_chars:
            desc_char = desc_chars[lower_chars.index(c)]
        elif c in upper_chars:
            desc_char = desc_chars[upper_chars.index(c)]
            desc_char = "SHIFT+" + desc_char
        elif c == ' ': #space
            desc_char = 'SPACE'
        elif c == 0x08: # Backspace
            desc_char = 'BACKSPACE'
        elif c == 0x09: # Tab
            desc_char = 'TAB'
        elif c == 0x0d: # Enter
            desc_char = 'ENTER'
        if desc_char != '':
            print("%s "%desc_char, end="")
            err,info = keyboard(desc_char)
            if err!=0:
                break
        else:
            err = -1
            info = "Unsupported char '%s'"%c
            break
    return err, info

def backspace(n):
    for i in range(n):
        err,info = keyboard('BACKSPACE')
        time.sleep(0.02)
    return 0, "OK"

def sleep(sec):
    time.sleep(sec)
    return 0

if 0:
    from maix import image, nn, app
    model = "/root/models/pp_ocr.mud"
    ocr = nn.PP_OCR(model)

def capture_screen(prefix):
    t0 =time.time()
    url = "https://localhost/api/stream/mjpeg"
    try:
        resp = requests.get(url, stream=True, timeout=10, verify=False)
        resp.raise_for_status()
        content_type = resp.headers.get('Content-Type', '')
        if 'multipart/x-mixed-replace' not in content_type:
            raise ValueError(f"Unexpected content type: {content_type}")
        # boundary 提取
        if 'boundary=' in content_type:
            boundary = '--' + content_type.split('boundary=')[-1].strip('";')
        else:
            raise ValueError("Cannot find boundary in Content-Type header")
        boundary_bytes = boundary.encode()
        buffer = bytearray()
        frame_count = 0
        max_frames_to_try = 3


        for chunk in resp.iter_content(chunk_size=16384):  # 增大 chunk
            buffer.extend(chunk)
            
            while True:
                boundary_pos = buffer.find(boundary_bytes)
                if boundary_pos == -1:
                    break

                frame_data = buffer[:boundary_pos]
                del buffer[:boundary_pos + len(boundary_bytes)]

                jpeg_start = frame_data.find(b'\xff\xd8')
                jpeg_end = frame_data.find(b'\xff\xd9', jpeg_start)
                if jpeg_start != -1 and jpeg_end != -1:
                    jpeg_data = frame_data[jpeg_start:jpeg_end+2]

                    try:
                        img = Image.open(io.BytesIO(jpeg_data))
                        filename = None
                        if prefix:
                            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
                            filename = prefix + f"{timestamp}.jpg"
                            img.save(filename)
                        buffered = io.BytesIO()
                        img.save(buffered, format="JPEG")  
                        img_base64 = base64.b64encode(buffered.getvalue()).decode('utf-8')
                        t1 =time.time();print("take %.3fs to snapshoot"%(t1-t0))  # ~0.4s
                        return img_base64, filename

                    except Exception as e_img:
                        print(f"Error processing image: {e_img}")
                        frame_count += 1
                        if frame_count >= max_frames_to_try:
                            raise Exception("Failed processing multiple frames")
        raise Exception("No valid JPEG frame found")

    except Exception as e:
        print(f"Error fetching frame: {e}")
        return None, None


def capture_screen_pikvm(prefix):
    if prefix == None:
        prefix = "/tmp/cua_"
    now = datetime.now()
    timestamp = now.strftime("%Y%m%d_%H%M%S")
    filename = prefix + f"{timestamp}.jpg"
    cmd = ["curl","-k","-u", "admin:admin",
        "https://localhost/api/streamer/snapshot?save=1&preview_quality=95",
        "--output", filename]
    subprocess.run(cmd, check=True)
    img = encode_image(filename)
    return img, filename

