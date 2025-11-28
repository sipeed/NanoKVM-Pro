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
      noAccount:
        'Не удалось получить информацию о пользователе, пожалуйста, обновите веб-страницу или сбросьте пароль',
      invalidUser: 'Неверное имя пользователя или пароль',
      error: 'Непредвиденная ошибка',
      changePassword: 'Изменить пароль',
      changePasswordDesc: 'Для безопасности вашего устройства, пожалуйста поменяйте пароль!',
      differentPassword: 'Пароли не совпадают',
      illegalUsername: 'Имя пользователя содержит недопустимые символы',
      illegalPassword: 'Пароль содержит недопустимые символы',
      forgetPassword: 'Забыли пароль?',
      ok: 'ОК',
      cancel: 'Отмена',
      loginButtonText: 'Войти'
    },
    wifi: {
      title: 'Wi-Fi',
      description: 'Настройка Wi-Fi для NanoKVM',
      success: 'Пожалуйста, проверьте сетевой статус NanoKVM и откройте в браузере новый IP адрес.',
      failed: 'Операция не удалась, пожалуйста попробуйте снова.',
      confirmBtn: 'OK',
      finishBtn: 'Готово'
    },
    screen: {
      title: 'Экран',
      video: 'Видеорежим',
      quality: 'Качество',
      qualityLossless: 'Без потерь',
      qualityHigh: 'Высокий',
      qualityMedium: 'Средний',
      qualityLow: 'Низкий',
      resetHdmi: 'Перезагрузить HDMI подсистему'
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
      mountFailed: 'Монтирование образа не удалось',
      mountDesc:
        'В некоторых системах необходимо отсоединить виртуальный диск на удаленном хосте перед монтированием образа.',
      unmountFailed: 'Отключение не удалось',
      unmountDesc:
        'В некоторых системах необходимо вручную извлечь образ на удаленном хосте перед отключением.',
      refresh: 'Обновить список образов'
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
      nanokvm: 'Терминал NanoKVM',
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
      ok: 'ОК'
    },
    download: {
      title: 'Скачать образ',
      input: 'Введите адрес удаленного образа',
      ok: 'OK',
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
      okBtn: 'Да',
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
        image: 'Версия прошивки',
        imageTip: 'Версия системного образа NanoKVM',
        deviceKey: 'Ключ устройства',
        community: 'Сообщество',
        hostname: 'Имя хоста',
        hostnameUpdated:
          'Имя хоста изменено. Перезагрузите NanoKVM, чтобы настройки вступили в силу.',
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
        menuBar: 'Панель меню',
        menuBarDesc: 'Иконки для отображения в панели меню',
        webTitle: 'Заголовок страницы',
        webTitleDesc: 'Изменить заголовок страницы'
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
          description: 'Настройка Wi-Fi'
        },
        ssh: {
          description: 'Включить удаленный доступ по SSH',
          tip: 'Задайте сильный пароль перед включением (Аккаунт - Пароль - Обновить)'
        },
        advanced: 'Дополнительные настройки',
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
        hidOnly: 'Режим только HID',
        disk: 'Виртуальный диск',
        diskDesc: 'Смонтировать виртуальный USB диск на удаленном хосте',
        network: 'Виртуальная сеть',
        networkDesc: 'Смонтировать виртуальное сетевое устройство на удаленном хосте',
        reboot: 'Перезагрузка',
        rebootDesc: 'Вы уверены, что хотите перезагрузить NanoKVM?',
        okBtn: 'Да',
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
        okBtn: 'Да',
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
        cancel: 'Отмена',
        preview: 'Ранние обновления',
        previewDesc: 'Получайте ранний доступ к новым функциям и улучшениям',
        previewTip:
          'Обратите внимание: в ранних версиях могут быть ошибки или незавершённый функционал!'
      },
      account: {
        title: 'Аккаунт',
        webAccount: 'Имя веб-аккаунта',
        password: 'Пароль',
        updateBtn: 'Обновить',
        logoutBtn: 'Выйти',
        logoutDesc: 'Вы действительно хотите выйти?',
        okBtn: 'Да',
        cancelBtn: 'Нет'
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
