const ru = {
  translation: {
    head: {
      desktop: 'Удаленный рабочий стол',
      login: 'Войти',
      changePassword: 'Изменить пароль',
      terminal: 'Терминал',
      wifi: 'Wi-Fi'
    },
    auth: {
      login: 'Войти',
      placeholderUsername: 'Введите имя пользователя',
      placeholderPassword: 'Введите пароль',
      placeholderPassword2: 'Введите пароль снова',
      noEmptyUsername: 'Имя пользователя не может быть пустым',
      noEmptyPassword: 'Пароль не может быть пустым',
      noAccount: 'Не удалось получить информацию о пользователе, пожалуйста, обновите веб-страницу или сбросьте пароль',
      invalidUser: 'Неверное имя пользователя или пароль',
      error: 'Непредвиденная ошибка',
      changePassword: 'Изменить пароль',
      changePasswordDesc: 'Для безопасности вашего устройства, пожалуйста поменяйте пароль!',
      differentPassword: 'Пароли не совпадают',
      illegalUsername: 'Имя пользователя содержит недопустимые символы',
      illegalPassword: 'Пароль содержит недопустимые символы',
      forgetPassword: 'Забыли пароль?',
      ОК: 'ОК',
      cancel: 'Отмена',
      loginButtonText: 'Войти'
      tips: {
        reset1: 'Web-аккаунт по умолчанию:',
        reset2: 'SSH-аккаунт по умолчанию:',
        reset3: 'Для сброса паролей, пожалуйста, ознакомьтесь с ',
        change1: 'Пожалуйста, имейте в виду — это действие изменит следующие пароли:',
        change2: 'Пароль веб-доступа',
        change3: 'Системный пароль от root (пароль SSH-доступа)',
        change4: 'Для сброса паролей, пожалуйста, ознакомьтесь с '
      }
    },
    assistant: {
      title: 'Умный Помощник',
      tips: {
        intro:
          'Умный Помощник - это экспериментальная функция NanoKVM, позволяющая испытать применение Модель «зрение-язык» (Vision-Language Model — VLM) для автоматизации компьютерных операций, часто называемых "Применением ПК".\n Перед применением этой функции, пожалуйста, выставьте видеорежим в MJPEG и задайте разрешение экрана 720p или ниже.\nНа текущий момент Умный Помощник полагается на внешний VLM-сервис, что требует предоставить соответствующие параметры сервера модели и API-ключ.\nЭто значит, что указанный сервер модели будет получать снимки Вашего компьютерного экрана, и с Вас может взиматься плата за использование токеном.\nВы также можете настроить собственный сервер с VLM-моделью в Вашей локальной сети для избежания этих проблем.',
        content:
          "Большие модели по-прежнему имеют непредсказуемый риск галлюцинаций и ограниченность функциональных способностей, что может привести к нежелаемым действиям, который могут привести к повреждению компьютерных данных. Вы должны осознавать риск. Когда Умный Помощник работает, Вы должны следить за его действиями и быть готовыми прервать любую потенциально рискованную операцию в любой момент!\nКликая на кнопку ниже Вы откроете Умного Помощника в новой окне, тем самым осознавая и принимая ранее описанные риски.\nВ целях безопасности, сервис ограничен одним пользователем. При закрытии страницы сервис автоматически выключится, требуя от Вас повторного нажатия на кнопку для запуска новой сессии.",
        contribute:
          'Исходный код был открыть, поэтому Вы свободно можете его проверить. Вклады тоже приветствуются!',
        warning: 'Внимание'
      },
      install: 'Установить зависимости',
      tryNow: 'Попробовать сейчас'
    },
    wifi: {
      title: 'Wi-Fi',
      description: 'Настройка Wi-Fi для NanoKVM',
      success: 'Пожалуйста, проверьте сетевой статус NanoKVM и откройте в браузере новый IP адрес.',
      failed: 'Операция не удалась, пожалуйста попробуйте снова.',
      invalidMode: 'Текущий режим не поддерживает настройку сети. Пожалуйста, включите режим настройки Wi-Fi на устройстве.',
      confirmBtn: 'ОК',
      finishBtn: 'Готово'
    },
    notification: {
      browser: {
        title: 'Рекомендация к браузеру',
        description:
          'Для лучшего пользвоательского опыта мы рекомендуем исполдьзовать браузер Chrome. Пока мы стремимся к совместимости с основными браузерами, некоторые продвинутые функции за пределами Chrome могут работать не так, как задумано.'
      },
      h265: {
        title: 'Внимание',
        description:
          "В предустановленной на устройстве прошивке отсутствуют функции кодирования или декодирования видео H.265/HEVC..\n\n" +
          'Любая пользовательская модификация прошивки, ручная установка сторонних модулей H.265 или самостоятельная разработка функций HEVC подразумевает принятие пользователем всех связанных с этим рисков. Эти риски включают, помимо прочего, плату за лицензирование патентов, юридическую ответственность, проблемы совместимости и снижение производительности.\n\n' +
          'Производитель отказывается от любой ответственности за работу, безопасность или соответствие законодательству любых настраиваемых пользователем функций H.265.',
        notSupported: 'Данный бреузер не поддерживает H.265!'
      },
      btnОК: 'ОК'
    },
    screen: {
      title: 'Экран',
      video: 'Видеорежим',
      quality: 'Качество',
      qualityAuto: 'Автоматически',
      qualityLossless: 'Без потерь',
      qualityHigh: 'Высокий',
      qualityMedium: 'Средний',
      qualityLow: 'Низкий',
      scale: 'Масштаб',
      resetHdmi: 'Перезагрузить HDMI подсистему',
      noSignal: 'Нет сигнала HDMI',
      inconsistentVideoMode: 'Воспроизведение приостановлено. Идёт воспроизведение другого видеорежима.',
      upload: 'Загрузить',
      advanced: 'Дополнительно'
    },
      audio: {
        title: 'Аудио'
    },
    keyboard: {
      title: 'Клавиатура',
      paste: 'Вставка текста',
      tips: 'Поддерживаются только стандартные буквы и символы клавиатуры',
      placeholder: 'Текст для ввода',
      submit: 'Вставить',
      virtual: 'Клавиатура',
      ctrlaltdel: 'Ctrl+Alt+Del'
    },
    mouse: {
      title: 'Мышь',
      cursor: 'Стиль курсора',
      default: 'Курсор по умолчанию',
      pointer: 'Курсор: указатель',
      cell: 'Курсор: крест',
      text: 'Курсор: выбор текста',
      grab: 'Курсор: захват',
      hide: 'Скрыть курсор',
      mode: 'Режим мыши',
      absolute: 'Абсолютный режим',
      relative: 'Относительный режим',
      speed: 'Скорость колёсика',
      fast: 'Быстро',
      slow: 'Медленно',
      requestPointer: 'Используется относительный режим. Нажмите на рабочий стол для захвата мыши.',
      resetHid: 'Перезагрузить HID-подсистему',
      hidOnly: {
        title: 'Режим только HID',
        desc: 'Если мышь и клавиатура перестали отвечать, и перезагрузка HID-подсистемы не помогает, это может быть связано с несовместимостью между NanoKVM и устройством. Попробуйте включить режим только HID для улучшения совместимости.',
        tip1: 'Включение режима только HID отключит виртуальный USB-диск и виртуальную сеть',
        tip2: 'В режиме только HID отключена возможность монтирования образов',
        enable: 'Включить режим только HID',
        disable: 'Выключить режим только HID'
      }
    },
    image: {
      title: 'Образы',
      loading: 'Загрузка...',
      empty: 'Пусто',
      upload: 'Загрузить образ',
      uploaded: 'Образ загружен',
      mount: 'Смонтировать образ',
      mounted: 'Смонтирован',
      fetchError: 'Ошибка получения списка образов, пожалуйста, обновите страницу и попробуйте ещё раз',
      uploadError: 'Загрузка не удалась, попробуйте ещё раз',
      mountMode: 'Режим монтирования',
      readOnly: 'Только чтение',
      mountFailed: 'Монтирование образа не удалось',
      mountDesc:
        'В некоторых системах необходимо отсоединить виртуальный диск на удаленном хосте перед монтированием образа.',
      unmountFailed: 'Отключение не удалось',
      unmountDesc:
        'В некоторых системах необходимо вручную извлечь образ на удаленном хосте перед отключением.',
      refresh: 'Обновить список образов',
      download: 'Скачать образ',
      attention: 'Внимание',
      deleteConfirm: 'Вы уверены, что хотите удалить этот образ?',
      okBtn: 'Да',
      cancelBtn: 'Нет'
    },
    script: {
      title: 'Скрипты',
      upload: 'Загрузить',
      run: 'Выполнить',
      runBackground: 'Выполнить в фоне',
      runFailed: 'Не удалось выполнить',
      attention: 'Внимание',
      delDesc: 'Вы уверены, что хотите удалить этот файл?',
      confirm: 'Да',
      cancel: 'Нет',
      delete: 'Удалить',
      close: 'Закрыть'
    },
    terminal: {
      title: 'Терминал',
      NanoKVM: 'Терминал NanoKVM',
      serial: 'Терминал COM-порта',
      serialPort: 'COM-порт',
      serialPortPlaceholder: 'Введите COM-порт',
      baudrate: 'Скорость передачи',
      parity: 'Чётность',
      parityNone: 'Нет',
      parityEven: 'Чётная',
      parityOdd: 'Нечётная',
      flowControl: 'Управление потоком',
      flowControlNone: 'Нет',
      flowControlSoft: 'Программное',
      flowControlHard: 'Аппаратное',
      dataBits: 'Бит данных',
      stopBits: 'Стоп-бит',
      confirm: 'ОК'
    },
    wol: {
      title: 'Wake-on-LAN',
      sending: 'Отправка команды...',
      sent: 'Команда отправлена',
      input: 'Введите MAC-адрес',
      ОК: 'ОК'
    },
    download: {
      title: 'Скачать образ',
      input: 'Введите адрес удаленного образа',
      ОК: 'ОК',
      disabled: 'Невозможно скачать образ, раздел /data находится в режиме только для чтения'
    },
    power: {
      title: 'Питание',
      showConfirm: 'Подтверждение',
      showConfirmTip: 'Если включено, операции с питанием потребуют дополнительного подтверждения',
      reset: 'Экстренная перезагрузка',
      power: 'Питание',
      powerShort: 'Питание (короткое нажатие)',
      powerLong: 'Питание (длинное нажатие)',
      resetConfirm: 'Вы уверены, что хотите выполнить экстренную перезагрузку?',
      powerConfirm: 'Вы уверены, что хотите нажать кнопку питания?',
      ОКBtn: 'Да',
      cancelBtn: 'Нет'
    },
    settings: {
      title: 'Настройки',
      about: {
        title: 'О системе NanoKVM',
        information: 'Информация',
        ip: 'IP-адрес',
        mdns: 'Доменное имя mDNS',
        application: 'Версия ПО',
        applicationTip: 'Версия веб-программы NanoKVM',
        arch: 'Arch',
        device_number: 'Device Number',
        image: 'Версия прошивки',
        imageTip: 'Версия системного образа NanoKVM',
        deviceKey: 'Ключ устройства',
        community: 'Сообщество',
        hostname: 'Имя хоста',
        hostnameUpdated: 'Имя хоста изменено. Перезагрузите NanoKVM, чтобы настройки вступили в силу.',
        pikvm: {
          title: 'Переключиться в PiKVM',
          attention: 'Внимание',
          desc1: 'Вы уверены, что хотите переключить текущую систему в PiKVM?',
          desc2: 'Устройство автоматиченски перезагрузиться после переключения системы.',
          confirm: 'Подтвердить'
        },
        ipType: {
          Wired: 'Проводное',
          Wireless: 'Беспроводное',
          Other: 'Другое'
        }
      },
      appearance: {
        title: 'Внешний вид',
        display: 'Отображение',
        language: 'Язык',
        languageDesc: 'Выберите язык интерфейса',
        menuBar: 'Панель меню',
        menuBarDesc: 'Иконки для отображения в панели меню',
        webTitle: 'Заголовок страницы',
        webTitleDesc: 'Изменить заголовок страницы'
      },
      screen: {
        title: 'Экран',
        videoMode: {
          title: 'Видеорежим',
          description: 'Select the encoding and transmission method'
        },
        edid: {
          description: 'Выберите параметры дисплея',
          default: 'По умолчанию',
          custom: 'Свои параметры'
        },
        rateControlMode: {
          title: 'Режим управления скоростью',
          description: 'Выберите режим управления битрейтом',
          cbr: 'CBR (постоянный битрейт): Приоритет стабильности подключнеия. Минимизирует рывки, но качество изображения может упасть во время сцен с большим количеством интенстивных движений.',
          vbr: 'VBR (переменный битрейт): Приоритет чёткости изображения. Сцены с большим количеством интенстивных движений могут привести к пикам потребления пропускной способности, потенциально приводя к заджержкам и рымкав.'
        },
        bitrate: {
          title: 'Битрейт',
          description: 'Баланс чёткости изображения и использования пропускной способности',
          lower:
            'Низкие значения снизят задержку в слабых сетях, но снизит качество изображения, потенциально приводя к размытию или артефактам.',
          higher:
            'Высокие значения улучшат чёткость изображения, но потребляет больше пропускной способности. Может увеличить задержки при слабой сети.'
        },
        gop: {
          description: 'Задать интервал между полными ключевыми кадрами',
          lower:
            'При низкх значениях полный кадр отправляется чаще. Это снижает время устранениня графический артефактов, но потребляет больше пропускной способности и снижает детализацию изображения.',
          higher:
            'Высокие значения опитмизированы для инкрементального обновления. Это улучшает детализацию изображения и снижает потребление пропускной способности, но увеличивает время устранения графических артефактов.'
        },
        quality: {
          title: 'Качество',
          description: 'Установить качество отображаемого изображения'
        },
        scale: {
          title: 'Масштаб',
          description: 'Установить соотношение масштаба изображения'
        }
      },
      device: {
        title: 'Устройство',
        oled: {
          title: 'OLED экран',
          description: 'Автоматическое отключение OLED экрана',
          0: 'Никогда',
          15: '15 сек',
          30: '30 сек',
          60: '1 мин',
          180: '3 мин',
          300: '5 мин',
          600: '10 мин',
          1800: '30 мин',
          3600: '1 час'
        },
        wifi: {
          title: 'Wi-Fi',
          description: 'Настройка Wi-Fi',
          apMode: 'AP mode is enabled, connect to Wi-Fi manually',
          others: 'Others...',
          connect: 'Join Wi-Fi',
          connectDesc1: 'Please enter the network ssid and password',
          connectDesc2: 'Please enter the password to join this network',
          disconnect: 'Are you sure to disconnect the network?',
          ssid: 'Name',
          password: 'Password',
          joinBtn: 'Join',
          confirmBtn: 'Ok',
          cancelBtn: 'Cancel'
        },
        led: {
          title: 'LED Strip',
          description: 'Configure LED strip',
          horizontal: 'Number of horizontal beads',
          vertical: 'Number of vertical beads',
          brightness: 'LED brightness',
          restart: 'Turn the LED strip off and on to apply the new configuration.'
        },
        passthrough: {
          title: 'HDMI Passthrough',
          tip: 'Enabled by default, the capture rate is limited to standard monitor frame rates like 4K30. Disabling it allows for non-standard frame rates, such as 4K40. This feature may not be applicable to some monitors.',
          description: 'Enable HDMI Passthrough'
        },
        capture: {
          title: 'HDMI Capture',
          tip: 'Enabled by default, the passthrough is limited to 4K30. You can disable it when local use to enable 4K60 passthrough.',
          description: 'Enable HDMI Capture'
        },
        ssh: {
          description: 'Включить удаленный доступ по SSH',
          tip: 'Задайте сильный пароль перед включением (Аккаунт - Пароль - Обновить)'
        },
        advanced: 'Дополнительные настройки',
        staticIp: {
          title: 'Static IP',
          config: 'Config',
          tip: 'To prevent connection loss, it is recommended to configure while you have physical access to the device.',
          document: 'Check the document.',
          enable: 'Enable static IP',
          placeholder: 'Format: IP/CIDR [Gateway]\nExample: 192.168.1.10/24',
          notEmpty: 'IP address cannot be empty.',
          notValid: 'IP address format is invalid.',
          failed: 'Configuration failed.',
          success: 'Configuration complete. Please check if the IP address is available.',
          okBtn: 'Ok',
          cancelBtn: 'Cancel'
        },
        mouseJiggler: {
          title: 'Дрожалка для мыши',
          description: 'Включите данную опцию, если хост засыпает',
          disable: 'Отключить',
          absolute: 'Абсолютный режим',
          relative: 'Относительный режим'
        },
        mdns: {
          description: 'Включить сервис обнаружения mDNS'
        },
        hdmi: {
          description: 'Enable HDMI/monitor output'
        },
        hidOnly: 'Режим только HID',
        disk: 'Виртуальный диск',
        diskDesc: 'Смонтировать виртуальный USB диск на удаленном хосте',
        emmc: {
          warning: 'Warning',
          tip1: "Mounting the eMMC for the first time will wipe all data in the NanoKVM's /data directory. Please ensure you have a backup before proceeding.",
          tip2: 'Do you want to continue?'
        },
        network: 'Виртуальная сеть',
        networkDesc: 'Смонтировать виртуальное сетевое устройство на удаленном хосте',
        reboot: 'Перезагрузка',
        rebootDesc: 'Вы уверены, что хотите перезагрузить NanoKVM?',
        datetime: {
          timezone: 'Time zone',
          datetime: 'Datetime',
          format: 'Time format',
          synchronize: 'Synchronize time',
          lastSynchronization: 'Last synchronization: ',
          notSynchronized: 'The system time is not synchronized yet',
          syncNow: 'Sync now'
        },
        ОКBtn: 'Да',
        cancelBtn: 'Нет'
      },
      tailscale: {
        title: 'Tailscale',
        restart: 'Вы уверены что хотите перезагрузить Tailscale?',
        stop: 'Вы уверены что хотите остановить Tailscale?',
        stopDesc: 'Выйти из Tailscale и отключить его автоматический запуск при включении.',
        loading: 'Загрузка...',
        notInstall: 'Tailscale не найден! Пожалуйста, установите.',
        install: 'Установить',
        installing: 'Установка',
        failed: 'Не удалось установить',
        retry: 'Пожалуйста, обновите и попробуйте снова, или попробуйте установить вручную',
        download: 'Скачайте',
        package: 'установочный пакет',
        unzip: 'и разархивируйте его',
        upTailscale: 'Переместите tailscale в каталог /usr/bin/ на NanoKVM',
        upTailscaled: 'Переместите tailscaled в каталог /usr/sbin/ на NanoKVM',
        refresh: 'Обновите текущую страницу',
        notRunning: 'Tailscale is not running. Please start it to continue.',
        run: 'Start',
        notLogin: 'Устройство не привязано. Войдите, чтобы привязать его к аккаунту.',
        urlPeriod: 'Этот адрес действителен в течение 10 минут',
        login: 'Войти',
        loginSuccess: 'Вход выполнен',
        enable: 'Включить Tailscale',
        deviceName: 'Имя устройства',
        deviceIP: 'IP адрес устройства',
        account: 'Аккаунт',
        logout: 'Выход',
        logoutDesc: 'Вы действительно хотите выйти?',
        uninstall: 'Удалить Tailscale',
        uninstallDesc: 'Are you sure you want to uninstall Tailscale?',
        ОКBtn: 'Да',
        cancelBtn: 'Нет'
      },
      update: {
        title: 'Проверить обновления',
        queryFailed: 'Получить версию не удалось',
        updateFailed: 'Обновление не удалось. Пожалуйста, попробуйте еще раз.',
        isLatest: 'У вас уже последняя версия.',
        available: 'Доступно обновление. Вы уверены, что хотите обновить?',
        updating: 'Начато обновление. Пожалуйста, подождите...',
        confirm: 'Подтвердить',
        upgrade_tip_nano: 'A new version of nanokvm is available. Are you sure you want to update now?',
        upgrade_tip_firmware: 'A new firmware version is available. Are you sure you want to update now?',
        cancel: 'Отмена',
        preview: 'Ранние обновления',
        previewDesc: 'Получайте ранний доступ к новым функциям и улучшениям',
        previewTip: 'Обратите внимание: в ранних версиях могут быть ошибки или незавершённый функционал!',
        download: 'Download',
        install: 'Install',
        changelog: 'Changelog',
        restart: 'Restart'
      },
      account: {
        title: 'Аккаунт',
        webAccount: 'Имя веб-аккаунта',
        password: 'Пароль',
        updateBtn: 'Обновить',
        logoutBtn: 'Выйти',
        logoutDesc: 'Вы действительно хотите выйти?',
        ОКBtn: 'Да',
        cancelBtn: 'Нет'
      },
      kvmadmin: {
        title: 'NanoKVM Admin',
        description: 'Batch management of NanoKVM devices',
        tip: 'Эта функция использует протокол mDNS. Пожалуйста, убедитесть, что mDNS включён на устройстве.',
        install: 'Установить',
        start: 'Запустить',
        visit: 'Открыть',
        uninstall: 'Удалить',
        attention: 'Внимание',
        confirmUninstall: 'Вы уверены, что хотите удалить NanoKVM Admin?',
        clearData: 'Эта операция удалить все сохранённые данные!',
        installFailed: 'Установка не удалась. Попробуйте ещё раз.',
        startFailed: 'Запуск сервиса не удался. Попробуйте ещё раз.',
        okBtn: 'Yes',
        cancelBtn: 'No'
      }
    },
    error: {
      title: 'У нас возникла проблема',
      refresh: 'Обновить страницу'
    },
    fullscreen: {
      toggle: 'Полноэкранный режим'
    },
    menu: {
      collapse: 'Свернуть меню',
      expand: 'Развернуть меню'
    }
  }
};

export default ru;
