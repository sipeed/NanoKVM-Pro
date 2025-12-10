const zh = {
  translation: {
    head: {
      desktop: '远程桌面',
      login: '登录',
      changePassword: '修改密码',
      terminal: '终端',
      wifi: 'Wi-Fi'
    },
    auth: {
      login: '登录',
      placeholderUsername: '请输入用户名',
      placeholderPassword: '请输入密码',
      placeholderPassword2: '请再次输入密码',
      noEmptyUsername: '用户名不能为空',
      noEmptyPassword: '密码不能为空',
      noAccount: '获取用户信息失败，请刷新重试或重置密码',
      invalidUser: '用户名或密码错误',
      error: '未知错误',
      changePassword: '修改密码',
      changePasswordDesc: '为了您的设备安全，请修改密码!',
      differentPassword: '两次密码不一致',
      illegalUsername: '用户名中包含非法字符',
      illegalPassword: '密码中包含非法字符',
      forgetPassword: '忘记密码',
      ok: '确定',
      cancel: '取消',
      loginButtonText: '登录',
      tips: {
        reset1: '网页默认帐号：',
        reset2: 'SSH 默认帐号：',
        reset3: '如需重置密码，请参考 ',
        change1: '请注意，此操作将同时更新以下密码：',
        change2: '网页的登录密码',
        change3: '系统 root 用户的密码（SSH 登录密码）',
        change4: '如需重置密码，请参考 '
      }
    },
    wifi: {
      title: 'Wi-Fi',
      description: '配置 NanoKVM Wi-Fi 信息',
      success: '请前往设备检查 NanoKVM 的网络状态。',
      failed: '操作失败，请重试。',
      invalidMode: '当前模式不支持配置网络。请先前往设备启用 Wi-Fi 配置模式。',
      confirmBtn: '确定',
      finishBtn: '完成'
    },
    notification: {
      browser: {
        title: '浏览器推荐',
        description:
          '您当前使用的浏览器可能无法完全支持所有高级功能与特性，推荐使用 Chrome 以获得最佳体验。'
      },
      h265: {
        title: '风险提示',
        description:
          '本设备的出厂固件已明确配置为不含 H.265/HEVC 视频编码或解码功能。\n\n' +
          '任何由用户自行发起、对固件进行的修改、手动安装第三方 H.265 模块，或独立开发 HEVC 功能的行为，均构成最终用户对所有相关风险的承担。\n\n' +
          '这些风险包括但不限于专利授权费用、法律责任、设备兼容性问题及性能下降。\n\n' +
          '制造商特此明确声明，对任何用户自定义 H.265 功能的运行、安全性或法律合规性，不承担任何责任。',
        notSupported: '当前浏览器不支持 H.265！'
      },
      btnOk: '确认'
    },
    screen: {
      scale: '缩放',
      title: '屏幕',
      video: '视频模式',
      quality: '图像质量',
      qualityLossless: '无损',
      qualityHigh: '高',
      qualityMedium: '中',
      qualityLow: '低',
      resetHdmi: '重置 HDMI',
      noSignal: 'HDMI 无信号',
      inconsistentVideoMode: '已暂停，其他视频模式正在播放中',
      upload: '上传'
    },
    keyboard: {
      title: '键盘',
      paste: '粘贴',
      tips: '仅支持标准键盘的字母和符号',
      placeholder: '请输入内容',
      submit: '确定',
      virtual: '虚拟键盘',
      ctrlaltdel: 'Ctrl+Alt+Del'
    },
    assistant: {
      title: 'Smart Assistant / 智能助手 ',
      tips: {
        intro:
          'Smart Assistant 是NanoKVM的实验性功能，可让你体验使用VLM大模型对电脑进行自动化操作，常被称为"Computer Use".\n在使用该功能前请调整视频模式为MJPEG传输，建议屏幕分辨率在720P或以下.\n目前智能助手依赖于外部VLM模型服务，需要你自行填入相应的模型服务器地址和API key，这意味着你填写的模型服务器会收到你的电脑屏幕截图，并可能产生相关tokens费用。\n你也可以自建运行VLM模型的内网服务器来避免上述问题。',
        content:
          '大模型目前仍有不可预知的幻觉风险和受限的操作能力，有可能做出意外的损坏电脑数据的操作，请你务必知晓这个风险，在智能助手操作时，请务必观察其操作，随时中断可能的风险操作！\n点击以下按钮即会在新窗口中打开智能助手，也表示你认可并接受以上风险。\n出于安全性考虑，该服务目前仅允许一个用户使用，且关闭页面后自动关闭服务，需要重新点击按钮来开启。',
        contribute: '该代码已开源，你可放心查看，也欢迎贡献！',
        warning: '警告'
      },
      install: '安装运行依赖',
      tryNow: '立即体验'
    },
    mouse: {
      title: '鼠标',
      cursor: '光标样式',
      default: '默认光标',
      pointer: '悬浮指针',
      cell: '单元指针',
      text: '文本指针',
      grab: '抓取指针',
      hide: '隐藏指针',
      mode: '鼠标模式',
      absolute: '绝对模式',
      relative: '相对模式',
      speed: '滚轮速度',
      fast: '快',
      slow: '慢',
      requestPointer: '正在使用鼠标相对模式，请点击桌面获取鼠标指针。',
      resetHid: '重置 HID',
      hidOnly: {
        title: 'HID-Only 模式',
        desc: '若使用过程中遇到鼠标键盘无响应，且重置 HID 无效，可能是 NanoKVM 与您的设备存在兼容性问题。建议尝试启用 HID-Only 模式以提升兼容性。',
        tip1: '启用该模式会禁用虚拟U盘和虚拟网卡',
        tip2: '该模式下无法使用镜像挂载功能',
        enable: '启用 HID-Only 模式',
        disable: '关闭 HID-Only 模式'
      }
    },
    image: {
      title: '镜像',
      loading: '加载中',
      empty: '无镜像文件',
      upload: '上传镜像',
      uploaded: '已上传',
      mount: '挂载镜像',
      mounted: '已挂载',
      fetchError: '获取镜像列表失败，请刷新重试',
      uploadError: '上传失败，请重试',
      mountMode: '模式',
      readOnly: '只读',
      mountFailed: '挂载失败',
      mountDesc: '在某些系统中，需要在远程主机中弹出虚拟硬盘后再挂载镜像。',
      unmountFailed: '卸载失败',
      unmountDesc: '在某些系统中，需要在远程主机中手动弹出后再卸载镜像。',
      refresh: '刷新镜像列表',
      download: '下载镜像',
      attention: '注意',
      deleteConfirm: '确定要删除该镜像吗？',
      okBtn: '确定',
      cancelBtn: '取消'
    },
    script: {
      title: '脚本',
      upload: '上传',
      run: '运行',
      runBackground: '后台运行',
      runFailed: '运行失败',
      attention: '注意',
      delDesc: '确定要删除该文件吗？',
      confirm: '确定',
      cancel: '取消',
      delete: '删除',
      close: '关闭'
    },
    terminal: {
      title: '终端',
      nanokvm: 'NanoKVM 终端',
      serial: '串口终端',
      serialPort: '串口',
      serialPortPlaceholder: '请输入串口',
      baudrate: '波特率',
      confirm: '确定'
    },
    wol: {
      title: 'Wake-on-LAN',
      sending: '指令发送中...',
      sent: '指令已发送',
      input: '请输入MAC地址',
      ok: '确定'
    },
    download: {
      title: '下载镜像',
      input: '请输入镜像的下载地址',
      ok: '确定',
      disabled: '/data 是只读分区，无法下载镜像'
    },
    power: {
      title: '电源',
      showConfirm: '显示确认框',
      showConfirmTip: '电源操作需要二次确认',
      reset: '重启',
      power: '电源',
      powerShort: '电源（短按）',
      powerLong: '电源（长按）',
      resetConfirm: '确认执行重启操作吗？',
      powerConfirm: '确认执行电源操作吗？',
      okBtn: '确认',
      cancelBtn: '取消'
    },
    settings: {
      title: '设置',
      about: {
        title: '关于 NanoKVM',
        information: '信息',
        ip: 'IP',
        mdns: 'mDNS',
        arch: '指令架构',
        device_number: '设备编号',
        application: '应用版本',
        applicationTip: 'NanoKVM 网页应用版本',
        image: '镜像版本',
        imageTip: 'NanoKVM 系统镜像版本',
        deviceKey: '设备码',
        community: '社区',
        hostname: '主机名',
        hostnameUpdated: '主机名修改成功，重启后生效',
        pikvm: {
          title: '切换 PiKVM',
          attention: '注意',
          desc1: '是否要将当前系统切换为 PiKVM？',
          desc2: '切换系统后设备会自动重启。',
          confirm: '确定'
        }
      },
      appearance: {
        title: '外观',
        display: '显示',
        language: '语言',
        languageDesc: '选择界面语言',
        menuBar: '菜单栏',
        menuBarDesc: '是否在菜单栏中显示图标',
        webTitle: '标题',
        webTitleDesc: '自定义网站页面标题'
      },
      device: {
        title: '设备',
        oled: {
          title: 'OLED',
          description: '设置 OLED 屏幕自动休眠时间',
          0: '永不',
          15: '15秒',
          30: '30秒',
          60: '1分钟',
          180: '3分钟',
          300: '5分钟',
          600: '10分钟',
          1800: '30分钟',
          3600: '1小时'
        },
        wifi: {
          title: 'Wi-Fi',
          description: '配置 Wi-Fi 信息',
          apMode: 'AP 模式已启用，该模式下仅支持手动连接',
          others: '其他...',
          connect: '连接 Wi-Fi',
          connectDesc1: '请输入网络名称和密码',
          connectDesc2: '请输入密码以连接此网络',
          disconnect: '是否要断开该网络连接？',
          ssid: '名称',
          password: '密码',
          joinBtn: '加入',
          confirmBtn: '确定',
          cancelBtn: '取消'
        },
        led: {
          title: 'LED 灯带',
          description: '设置 LED 灯带',
          horizontal: '水平灯珠数量',
          vertical: '垂直灯珠数量',
          brightness: '灯带亮度',
          restart: '需要重新开关 LED 灯带以应用新配置。'
        },
        passthrough: {
          title: 'HDMI 环出',
          tip: '默认开启，但是采集速率会被限制到4K30等常规显示器帧率，关闭以实现4K40等非标准帧率。该功能可能对部分显示器不适用。',
          description: '启用 HDMI 环出功能'
        },
        capture: {
          title: 'HDMI 采集',
          tip: '默认开启，但是环出会被限制到4K30，本地工作时可以关闭以实现4K60环出。',
          description: '启用 HDMI 采集功能'
        },
        ssh: {
          description: '启用 SSH 远程访问',
          tip: '启用前请务必设置强密码（帐号 - 修改密码）'
        },
        advanced: '高级设置',
        staticIp: {
          title: '静态 IP',
          config: '设置',
          tip: '建议在可物理接触设备的情况下进行设置，以防止配置错误导致设备失联，无法远程访问',
          document: '查看文档',
          enable: '启用静态 IP',
          placeholder: '格式：IP/CIDR [Gateway]\n示例：192.168.1.10/24',
          notEmpty: 'IP 地址不能为空',
          notValid: 'IP 地址内容有误',
          success: '设置完成，请检查 IP 是否可用',
          failed: '设置失败',
          okBtn: '确定',
          cancelBtn: '取消'
        },
        mouseJiggler: {
          title: '鼠标抖动',
          description: '防止远程主机休眠',
          disable: '关闭',
          absolute: '绝对模式',
          relative: '相对模式'
        },
        mdns: {
          description: '启用 mDNS 发现服务'
        },
        hdmi: {
          description: '启用 HDMI/显示器 输出功能'
        },
        hidOnly: 'HID-Only 模式',
        disk: '虚拟U盘',
        diskDesc: '在远程主机中挂载虚拟U盘',
        noSdCard: '请插入SD卡以启用此功能',
        network: '虚拟网卡',
        networkDesc: '在远程主机中挂载虚拟网卡',
        datetime: {
          timezone: '时区',
          datetime: '时间',
          format: '时间格式',
          synchronize: '同步时间',
          lastSynchronization: '上次同步：',
          notSynchronized: '系统时间尚未同步',
          syncNow: '立即同步'
        },
        okBtn: '是',
        cancelBtn: '否'
      },
      tailscale: {
        title: 'Tailscale',
        restart: '确定要重启 Tailscale 吗？',
        stop: '确定要停止 Tailscale 吗？',
        stopDesc: '该操作会退出登录，并停止开机自动启动。',
        loading: '加载中...',
        notInstall: '未检测到 Tailscale，请先安装',
        install: '安装',
        installing: '安装中',
        failed: '安装失败',
        retry: '请刷新后重试，或尝试手动安装',
        download: '下载',
        package: '安装包',
        unzip: '并解压',
        upTailscale: '将 tailscale 文件上传到 /usr/bin/ 目录',
        upTailscaled: '将 tailscaled 文件上传到 /usr/sbin/ 目录',
        refresh: '刷新页面',
        notRunning: 'Tailscale 尚未运行，请先执行启动操作',
        run: '启动',
        notLogin: '该设备尚未绑定，请点击登录并将这台设备绑定到您的账号。',
        urlPeriod: '该链接10分钟内有效',
        login: '登录',
        loginSuccess: '登录完成',
        enable: '启用 Tailscale',
        deviceName: '设备名称',
        deviceIP: '设备地址',
        account: '账号',
        logout: '退出',
        logoutDesc: '确定要退出吗？',
        uninstall: '卸载 Tailscale',
        uninstallDesc: '确定要卸载 Tailscale 吗？',
        okBtn: '确认',
        cancelBtn: '取消'
      },
      update: {
        title: '检查更新',
        queryFailed: '获取版本号失败',
        updateFailed: '更新失败，请重试',
        isLatest: '已经是最新版本。',
        available: '有新的可用版本，确定要更新吗？',
        updating: '更新中，请稍候...',
        confirm: '确定',
        cancel: '取消',
        preview: '预览更新',
        previewDesc: '率先体验即将推出的新功能和优化',
        previewTip: '预览版更新可能包含一些不稳定因素或未完善的功能！',
        download: '下载',
        install: '安装',
        changelog: '更新日志',
        restart: '重启服务'
      },
      account: {
        title: '帐号',
        webAccount: '网页帐号',
        password: '密码',
        updateBtn: '修改',
        logoutBtn: '退出',
        logoutDesc: '确定要退出吗？'
      },
      kvmadmin: {
        title: 'NanoKVM Admin',
        description: '批量管理 NanoKVM 设备',
        tip: '该功能使用 mDNS 协议，请确保设备中 mDNS 功能已开启',
        install: '安装',
        start: '运行',
        visit: '访问',
        uninstall: '卸载',
        attention: '注意',
        confirmUninstall: '确定要卸载 NanoKVM Admin 吗？',
        clearData: '该操作会删除所有已保存的数据！',
        installFailed: '安装失败，请重试',
        startFailed: '运行失败，请重试',
        okBtn: '确定',
        cancelBtn: '取消'
      }
    },
    error: {
      title: '我们遇到了问题',
      refresh: '刷新'
    },
    fullscreen: {
      toggle: '切换全屏'
    },
    menu: {
      collapse: '收起',
      expand: '展开'
    }
  }
};

export default zh;
