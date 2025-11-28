# cua_engine.py

import os, subprocess, time, base64
import cua_func

IMG_KEEP_N = 3
MAX_ROUNDS = 20
default_api_type = "DashScope"
default_base_url = ""   #"https://dashscope.aliyuncs.com/compatible-mode/v1"
default_model_name = "qwen3-vl-plus-2025-09-23"
default_initial_prompt="""
现在你是一个运行在IP-KVM上的智能助手, IP-KVM可以获取电脑的桌面截图和执行键鼠操作,使用你的能力帮助用户自动化实现目标任务。
你的回复内容会按以下格式解析为键鼠操作指令,请严格按照指定格式输出,如果完成任务,输出"Done"表示结束。
注意每次回复的操作指令只能有2~3行,,如果需要多步操作,可以分多轮回复,每轮回复后会等待你操作完成并截图反馈给你,然后你再根据最新的截图继续下一步操作。
1. 鼠标操作: mouse x y btn  
    x,y 鼠标移动到的归一化坐标,如果为0,0则不动; 
    btn 鼠标按键操作, 0不操作,1左键单击,2右键单击,3左键双击,4右键双击,5向下滚动一次,6向上滚动一次
2. 拖拽操作: drag x0 y0 x1 y1
    x0,y0 鼠标左键按下时的起始归一化位置
    x1,y1 鼠标左键拖拽松开的终点归一化位置
3. 键盘操作: keyboard key_desc
    key_desc is str like this: a,A,F11,CTRL+C,CTRL+SHIFT+A,ALT+F4,GUI+L
    special key str define like this:
    ENTER,ESCAPE,BACKSPACE,TAB,SPACE,CAPS_LOCK,
    RIGHT_ARROW,LEFT_ARROW,DOWN_ARROW,UP_ARROW,
    LEFT_CONTROL,LEFT_SHIFT,LEFT_ALT,LEFT_GUI,RIGHT_CONTROL,
    RIGHT_SHIFT,RIGHT_ALT,RIGHT_GUI,
4. 输入操作: typestr content
    content 待输入的内容,注意只能为ascii可见字符组成的字符串
5. 删除操作: backspace n
    n 退格键删除的字符个数
5. 延时操作: sleep second
    second 延时的秒数,可以为小数,如0.5表示延时0.5秒

注意点： 
1. 在移动光标之前，一定要通过截图检查目标元素的坐标，如果点击失败，请根据截图里的上次光标位置调整新光标坐标来点击
2. 在执行打开页面,打开程序,打开文件,保存文件，拖动文件等需要等待的操作后,必须等待截图反馈，确认上一步操作完成，一定要注意！！！
3. 操作的时候注意当前焦点窗口，如果有弹窗或者其它覆盖在上面的窗口可能影响你的操作，可以先关闭它们
4. 操作指令前，加上两行#开头文本，第一行简短描述当前页面重点内容，第二行描述你下面的指令要做什么

现在我下达的操作目标为："""

class VLM_dashscope:
    def __init__(self, api_key, base_url=default_base_url, MODEL_NAME=default_model_name, IMG_KEEP_N=3):
        self.API_KEY = api_key
        self.BASE_URL = base_url
        self.MODEL_NAME = MODEL_NAME 
        self.IMG_KEEP_N = IMG_KEEP_N
        self.reset()
    def reset(self):
        self.messages = [
            {
            "role": "system",
            "content": [{"text": "You are a helpful assistant."}]}
        ]
    def step(self, img_base64, prompt):
        from dashscope import MultiModalConversation  
        user_msg = {
            "role": "user",
            "content": [
                {"image": f"data:image/jpeg;base64,{img_base64}"},
                {"text": prompt}
            ]
        }
        self.messages.append(user_msg)
        # clear old image to save Token
        if len(self.messages) - 1 - self.IMG_KEEP_N * 2 >= 0:
            del self.messages[len(self.messages) - 1 - self.IMG_KEEP_N * 2]["content"][0]
        # call VLM
        try:
            response = MultiModalConversation.call(
                api_key=self.API_KEY,
                model=self.MODEL_NAME,  
                messages=self.messages)
            if response.output == None:
                return {'status': 'error', 'log': response['message']}
            res = response.output.choices[0].message.content[0]['text']
            self.messages.append(response['output']['choices'][0]['message'])
        except Exception as e:
            return {'status': 'error', 'log': f"LLM API Error: {str(e)}"} 
        return {'status': 'ok', 'res': res} 

class VLM_openai:
    def __init__(self, api_key, base_url=default_base_url, MODEL_NAME=default_model_name, IMG_KEEP_N=3):
        print("loading openai, it is slow...")
        print(time.time())
        from openai import OpenAI
        print(time.time())
        self.API_KEY = api_key
        self.BASE_URL = base_url
        self.MODEL_NAME = MODEL_NAME 
        self.IMG_KEEP_N = IMG_KEEP_N
        self.client = OpenAI(
            api_key=self.API_KEY,
            base_url=self.BASE_URL,
        )  
        self.reset()
    def reset(self):
        self.messages = [
            {
                "role": "system",
                "content": [{"type": "text", "text": "You are a helpful assistant."}],
            }
        ]
    def step(self, img_base64, prompt):
        user_msg = {
            "role": "user",
            "content": [
                {"type": "image_url", "image_url": {"url":f"data:image/jpeg;base64,{img_base64}"}},
                {"type": "text", "text": prompt}
            ]
        }
        self.messages.append(user_msg)
        # clear old image to save Token
        if len(self.messages) - 1 - self.IMG_KEEP_N * 2 >= 0:
            del self.messages[len(self.messages) - 1 - self.IMG_KEEP_N * 2]["content"][0]
        # call VLM
        try:
            completion = self.client.chat.completions.create(
                model=self.MODEL_NAME,
                messages=self.messages
            )
            res = completion.choices[0].message.content
            self.messages.append(completion.choices[0].message.model_dump())
        except Exception as e:
            return {'status': 'error', 'log': f"LLM API Error: {str(e)}"}
        return {'status': 'ok', 'res': res} 


class CUA_Engine:
    def __init__(self, api_type, api_key, base_url=default_base_url, model_name=default_model_name, img_keep_n=3, max_rounds=20, initial_prompt=default_initial_prompt):
        self.API_TYPE = api_type
        self.API_KEY = api_key
        self.BASE_URL = base_url
        self.MODEL_NAME = model_name
        self.IMG_KEEP_N = img_keep_n
        self.MAX_ROUNDS = max_rounds
        self.initial_prompt = initial_prompt

        # Accelerate startup, delay loading
        if self.API_TYPE == "DashScope":
            self.VLM = VLM_dashscope(api_key, self.BASE_URL, self.MODEL_NAME, self.IMG_KEEP_N)
        elif self.API_TYPE == "OpenAI":
            self.VLM = VLM_openai(api_key, self.BASE_URL, self.MODEL_NAME, self.IMG_KEEP_N)
        elif self.API_TYPE == "google-genai":
            # TODO: google-genai is the slowest loading API
            # self.VLM = VLM_genai(api_key, self.BASE_URL, self.MODEL_NAME, self.IMG_KEEP_N)
            raise "google-genai is not implement yet"
        else:
            raise "API_TYPE error!"

        # 初始化屏幕尺寸
        try:
            with open("/proc/lt6911_info/width", 'r') as f:
                self.SCREEN_W = int(f.read().strip())
            with open("/proc/lt6911_info/height", 'r') as f:
                self.SCREEN_H = int(f.read().strip())
        except (ValueError, IndexError, FileNotFoundError) as e:
            print(f"Error reading screen size: {e}")
            self.SCREEN_W = 1280  # 设置默认值
            self.SCREEN_H = 720

        # 初始化任务状态
        self.reset()

    def reset(self):
        self.VLM.reset()
        self.done_flag = False
        self.current_round = 0
        self.last_prompt = self.initial_prompt
        self.is_paused = False  # 标记是否被用户暂停
        self.user_input_queue = []  # 存放用户在暂停期间输入的附加信息

    def encode_image(self, path):
        """将图片文件编码为base64"""
        with open(path, "rb") as f:
            return base64.b64encode(f.read()).decode("utf-8")

    def get_current_screenshot(self):
        """获取当前屏幕截图的base64编码和路径"""
        img_base64, img_name = cua_func.capture_screen(None) #"/tmp/cua_")
        return img_base64, img_name

    def _execute_command(self, cmd_line):
        """内部方法：执行单条命令"""
        parts = cmd_line.split()
        cmd = parts[0]
        try:
            if cmd == "mouse":  # mouse x y btn
                x = float(parts[1])
                y = float(parts[2])
                btn = int(parts[3])
                if x>1 or y >1:
                    x = x/1000.0; y = y/1000.0
                cua_func.mouse(x, y, btn, self.SCREEN_W, self.SCREEN_H)
            elif cmd == "drag":
                x0 = float(parts[1])
                y0 = float(parts[2])
                x1 = float(parts[3])
                y1 = float(parts[4])
                if x0>1 or y0 >1 or x1>1 or y1 >1:
                    x0 = x0/1000.0; y0 = y0/1000.0
                    x1 = x1/1000.0; y1 = y1/1000.0
                cua_func.drag(x0, y0, x1, y1, self.SCREEN_W, self.SCREEN_H)
            elif cmd == "keyboard":
                desc_char = parts[1]
                err, info = cua_func.keyboard(desc_char)
                if err != 0:
                    return False, info
            elif cmd == "typestr":
                content = " ".join(parts[1:])
                err, info = cua_func.typestr(content)
                if err != 0:
                    return False, info
            elif cmd == "backspace":
                n = int(parts[1])
                cua_func.backspace(n)
            elif cmd == "sleep":
                sec = float(parts[1])
                time.sleep(sec)  
            else:
                return False, f"Unknown command: {cmd_line}"
        except (ValueError, IndexError) as e:
            return False, f"Command parameters error: {cmd_line} - {e}"
        return True, "Success"

    def step(self):
        """
        执行一步操作。
        返回一个字典，包含：
        - 'status': 'running', 'paused', 'waiting_for_user', 'done', 'error'
        - 'screenshot_base64': 当前截图的base64
        - 'llm_thoughts': 大模型返回的#描述
        - 'executed_command': 待执行或已执行的命令
        - 'log': 本次执行的详细日志信息
        """
        if self.done_flag:
            return {'status': 'done', 'log': "Task is already done."}
        if self.current_round >= self.MAX_ROUNDS:
            self.done_flag = True
            return {'status': 'error', 'log': "Warning: Too many rounds, task force stopped."}
        # 检查是否被用户暂停
        if self.is_paused:
            if self.user_input_queue:
                # 如果有用户输入，将其加入到prompt中
                user_input = self.user_input_queue.pop(0)
                self.last_prompt += f"\n\n[User Additional Input]: {user_input}"
                self.is_paused = False  # 恢复执行
            else:
                return {'status': 'paused', 'log': "Waiting for user input..."}

        # capture screen snap
        time.sleep(0.5)  # sleep 0.5s wait last operation done
        img_base64, img_path = self.get_current_screenshot()
        if img_base64 == None:
            return {'status': 'error', 'log': f"Capture screen failed!"}
        # call VLM
        result = self.VLM.step(img_base64, self.last_prompt)
        if result['status'] == 'error':
            return {'status': 'error', 'log': result['log']}
        res = result['res']

        # parse VLM response
        lines = res.split('\n')
        llm_thoughts = []
        commands_to_execute = []
        notify_input = None
        log_info = f"=============== Round {self.current_round}:\nUser: {self.last_prompt}\nLLM: {res}"
        self.current_round += 1

        # parse line by line
        status = 'none'
        for line in lines:
            print("@@@@@" +line)
            stripped_line = line.strip()
            if stripped_line == "":
                continue
            elif stripped_line.startswith('#'):
                if stripped_line !='# Done': 
                    llm_thoughts.append(stripped_line[1:].strip())  # remove #, as comment
                    continue
                else: # some model stupid 
                    return {
                        'status': "done",
                        'screenshot_base64': img_base64,
                        'llm_thoughts': llm_thoughts,
                        'log': log_info + "\nTask Completed!"
                    }
            elif stripped_line.startswith('$'):
                self.is_paused = True
                notify_input = stripped_line[1:].strip()    #remove $
                status = 'waiting_for_user'
                break
            commands_to_execute.append(stripped_line)
        # check status
        if status == 'waiting_for_user':
            self.last_prompt = ""
            return {
                'status': 'waiting_for_user',
                'screenshot_base64': img_base64,
                'llm_thoughts': llm_thoughts,
                'notify_input':notify_input,
                'log': log_info + "\nwaiting for user input."
            }
        elif commands_to_execute == []:
            self.last_prompt = "No valid command received from LLM."
            return {
                'status': 'running',
                'screenshot_base64': img_base64,
                'llm_thoughts': llm_thoughts,
                'log': log_info + "\nNo valid command, will retry."
            }
        elif commands_to_execute[0] == "Done":
            return {
                'status': "done",
                'screenshot_base64': img_base64,
                'llm_thoughts': llm_thoughts,
                'log': log_info + "\nTask Completed!"
            }
        elif (commands_to_execute[-1] == "Done" and len(commands_to_execute) > 1) or (len(commands_to_execute) > 3):
            # model in hallucination, notice him to check again
            self.last_prompt = "Please check the screenshot and command again, a maximum of 2 lines of commands can be entered each time."
            return {
                'status': 'running', 
                'screenshot_base64': img_base64,
                'llm_thoughts': llm_thoughts,
                'executed_command': "\n".join(commands_to_execute),
                'log': log_info + f"\nNot Executed due to Hallucination detected",
                # 新增一个字段，专门用于向前端传递错误信息
                'error_message': f"VLM Hallucination Detected, auto notify it to try again."
            }
            # dummy done, remove it, wait next turn to confirm

        # excute commands
        status = 'running' 
        success = True
        exec_logs = []
        for cmd in commands_to_execute:
            success, exec_log = self._execute_command(cmd)
            if not success:
                self.last_prompt = exec_log  # 将错误信息作为下一轮的prompt
                return {
                    'status': 'running',
                    'screenshot_base64': img_base64,
                    'llm_thoughts': llm_thoughts,
                    'executed_command': "\n".join(commands_to_execute),
                    'log': log_info + f"\nExecuted: {"\n".join(commands_to_execute)} -> ❌ ERROR: {exec_log}",
                    # 新增一个字段，专门用于向前端传递错误信息
                    'error_message': f"Command excute error: {exec_log}"
                }
            exec_logs.append(exec_log)

        self.last_prompt = ""  # 清空prompt，等待下一轮截图反馈
        return {
            'status': status,
            'screenshot_base64': img_base64,
            'llm_thoughts': llm_thoughts,
            'executed_command': "\n".join(commands_to_execute),
            'log': log_info + f"\nExecuted: {"\n".join(commands_to_execute)} -> {"\n".join(exec_logs)}"
        }

    def pause(self):
        """暂停当前任务"""
        self.is_paused = True

    def resume_with_input(self, user_input):
        """恢复任务，并传入用户输入"""
        self.user_input_queue.append(user_input)
        #self.is_paused = False  #前面有输入就自动退出暂停的

    def get_status(self):
        """获取当前任务状态"""
        return {
            'done': self.done_flag,
            'current_round': self.current_round,
            'is_paused': self.is_paused,
            'max_rounds': self.MAX_ROUNDS
        }