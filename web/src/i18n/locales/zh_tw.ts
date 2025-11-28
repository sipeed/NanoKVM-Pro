const zh_tw = {
  translation: {
    head: {
      desktop: '遠端桌面',
      login: '登入',
      changePassword: '更改密碼',
      terminal: '終端機',
      wifi: 'Wi-Fi'
    },
    auth: {
      login: '登入',
      placeholderUsername: '使用者名稱',
      placeholderPassword: '密碼',
      placeholderPassword2: '請再次輸入密碼',
      noEmptyUsername: '使用者名稱不能為空',
      noEmptyPassword: '密碼不能為空',
      noAccount: '取得使用者資訊失敗，請重新整理網頁或重設密碼',
      invalidUser: '使用者名稱或密碼錯誤',
      error: '非預期性錯誤',
      changePassword: '更改密碼',
      changePasswordDesc: '為了您的裝置安全，請修改登入密碼。',
      differentPassword: '密碼不一致',
      illegalUsername: '使用者名稱包含非法字元',
      illegalPassword: '密碼包含非法字元',
      forgetPassword: '忘記密碼',
      ok: '確定',
      cancel: '取消',
      loginButtonText: '登入'
    },
    wifi: {
      title: 'Wi-Fi',
      description: '設定 NanoKVM WiFi',
      success: '請檢查 NanoKVM 的網路狀態，並存取新的 IP 位址。',
      failed: '操作失敗，請重試。',
      confirmBtn: '確定',
      finishBtn: '完成'
    },
    screen: {
      video: '編碼格式',
      quality: '品質',
      qualityLossless: '無損',
      qualityHigh: '高',
      qualityMedium: '中',
      qualityLow: '低',
      resetHdmi: '重置 HDMI'
    },
    keyboard: {
      paste: '貼上',
      tips: '僅支援標準鍵盤字母和符號',
      placeholder: '請輸入內容',
      submit: '送出',
      virtual: '虛擬鍵盤',
      ctrlaltdel: 'Ctrl+Alt+Del'
    },
    mouse: {
      default: '預設游標',
      pointer: '懸浮游標',
      cell: '儲存格游標',
      text: '文字游標',
      grab: '抓取游標',
      hide: '隱藏游標',
      mode: '滑鼠模式',
      absolute: '絕對模式',
      relative: '相對模式',
      requestPointer: '使用相對模式。請按一下桌面以取得滑鼠游標。',
      resetHid: '重設 HID'
    },
    image: {
      title: '映像',
      loading: '載入中...',
      empty: '未找到任何內容',
      mountFailed: '掛載失敗',
      mountDesc: '在某些系統中，在掛載映像之前需要中斷遠端主機上的虛擬磁碟。',
      unmountFailed: '卸載失敗',
      unmountDesc: '在某些系統中，需要在遠端主機中手動彈出後再卸載映像。'
    },
    script: {
      title: '腳本',
      upload: '上傳',
      run: '執行',
      runBackground: '背景執行',
      runFailed: '執行失敗',
      attention: '注意',
      delDesc: '確定刪除該檔案？',
      confirm: '確定',
      cancel: '取消',
      delete: '刪除',
      close: '關閉'
    },
    terminal: {
      title: '終端機',
      nanokvm: 'NanoKVM 終端機',
      serial: 'Serial Port 終端機',
      serialPort: 'Serial Port',
      serialPortPlaceholder: '請輸入 Serial Port',
      baudrate: 'Baud rate',
      confirm: '確定'
    },
    wol: {
      title: 'Wake-on-LAN',
      sending: '發送指令中...',
      sent: '指令已發送',
      input: '請輸入 MAC 位址',
      ok: '確定'
    },
    download: {
      title: '下载映像',
      input: '請輸入映像檔的下載URL',
      ok: '確定',
      disabled: '/data 為唯讀目錄，無法下載映像檔'
    },
    power: {
      title: '電源',
      reset: '重新啟動',
      power: '電源',
      powerShort: '電源 (短按)',
      powerLong: '電源 (長按)'
    },
    settings: {
      title: '設定',
      about: {
        title: '關於 NanoKVM',
        information: '資訊',
        ip: 'IP',
        mdns: 'mDNS',
        application: '應用程式版本',
        applicationTip: 'NanoKVM 網頁程式版本',
        image: '映像版本',
        imageTip: 'NanoKVM 系统镜像版本',
        deviceKey: '設備序號',
        community: '社群'
      },
      appearance: {
        title: '外觀',
        display: '顯示',
        language: '語言',
        menuBar: '選單列',
        menuBarDesc: '是否在選單列中顯示圖案'
      },
      device: {
        title: '設備',
        oled: {
          title: 'OLED',
          description: '設定 OLED 螢幕自動睡眠時間',
          0: '永不',
          15: '15秒',
          30: '30秒',
          60: '1分鐘',
          180: '3分鐘',
          300: '5分鐘',
          600: '10分鐘',
          1800: '30分鐘',
          3600: '1小時'
        },

        wifi: {
          title: 'Wi-Fi',
          description: '設定 Wi-Fi'
        },
        ssh: {
          description: '啟用SSH',
          tip: '啟用前請務必設定強密碼（帳號 - 更改密碼）'
        },
        disk: '虛擬隨身碟',
        diskDesc: '在遠端主機上連接虛擬隨身碟',
        network: '虛擬網卡',
        networkDesc: '在遠端主機上新增虛擬網卡'
      },
      tailscale: {
        title: 'Tailscale',
        restart: '確定要重啟 Tailscale 嗎？',
        stop: '確定要停止 Tailscale 嗎？',
        stopDesc: '此操作將會登出帳號，並停止開機自動啟動。',
        loading: '載入中...',
        notInstall: '未找到 Tailscale ！請先安裝。',
        install: '安裝',
        installing: '安裝中',
        failed: '安裝失敗',
        retry: '請重新整理並重試。或嘗試手動安裝',
        download: '下載',
        package: '安裝包',
        unzip: '並解壓縮它',
        upTailscale: '將 Tailscale 上傳到 NanoKVM 目錄 /usr/bin/',
        upTailscaled: '將 Tailscale 上傳到 NanoKVM 目錄 /usr/sbin/',
        refresh: '重新整理頁面',
        notLogin: '設備尚未綁定。請登入並將該裝置綁定到您的帳戶。',
        urlPeriod: '此網址有效期限為 10 分鐘',
        login: '登入',
        loginSuccess: '登入成功',
        enable: '啟用 Tailscale',
        deviceName: '裝置名稱',
        deviceIP: '裝置 IP',
        account: '帳號',
        logout: '登出',
        uninstall: '移除 Tailscale',
        okBtn: '確認',
        cancelBtn: '取消'
      },
      update: {
        title: '檢查更新',
        queryFailed: '取得版本失敗',
        updateFailed: '更新失敗。請重試。',
        isLatest: '您已經擁有最新版本。',
        available: '有可用更新。確定要更新嗎？',
        updating: '正在更新。請稍等...',
        confirm: '確定',
        cancel: '取消'
      },
      account: {
        title: '帳號',
        webAccount: '網頁帳號',
        password: '密碼',
        updateBtn: '修改',
        logoutBtn: '登出'
      }
    },
    error: {
      title: '我們遇到了一些問題',
      refresh: '重新整理'
    },
    fullscreen: {
      toggle: '進入全螢幕模式'
    },
    menu: {
      collapse: '收起選單',
      expand: '展開選單'
    }
  }
};

export default zh_tw;
