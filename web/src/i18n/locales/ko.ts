const ko = {
  translation: {
    head: {
      desktop: '원격 데스크톱',
      login: '로그인',
      changePassword: '비밀번호 변경',
      terminal: '터미널',
      wifi: 'Wi-Fi'
    },
    auth: {
      login: '로그인',
      placeholderUsername: '사용자 이름을 입력하세요.',
      placeholderPassword: '비밀번호를 입력하세요.',
      placeholderPassword2: '비밀번호를 다시 입력하세요.',
      noEmptyUsername: '사용자 이름은 비어있을 수 없습니다.',
      noEmptyPassword: '비밀번호는 비어있을 수 없습니다.',
      noAccount:
        '사용자 정보를 불러오는 데 실패했습니다. 페이지를 새로고침하거나 비밀번호를 초기화하세요.',
      invalidUser: '사용자 이름이나 비밀번호가 틀렸습니다.',
      error: '알 수 없는 오류',
      changePassword: '비밀번호 변경',
      changePasswordDesc: '보안을 위해 웹 로그인 비밀번호를 변경하세요.',
      differentPassword: '비밀번호가 서로 일치하지 않습니다.',
      illegalUsername: '사용자 이름에 사용할 수 없는 문자가 있습니다.',
      illegalPassword: '비밀번호에 사용할 수 없는 문자가 있습니다.',
      forgetPassword: '비밀번호 분실',
      ok: '확인',
      cancel: '취소',
      loginButtonText: '로그인'
    },
    wifi: {
      title: 'Wi-Fi',
      description: 'NanoKVM Wi-Fi 설정',
      success: 'NanoKVM의 네트워크 상태를 확인하고 새 IP 주소로 접속하세요.',
      failed: '작업에 실패했습니다. 다시 시도하세요.',
      confirmBtn: '확인',
      finishBtn: '완료'
    },
    screen: {
      video: '비디오 모드',
      quality: '품질',
      qualityLossless: '무손실',
      qualityHigh: '높음',
      qualityMedium: '중간',
      qualityLow: '낮음',
      resetHdmi: 'HDMI 초기화'
    },
    keyboard: {
      paste: '붙여넣기',
      tips: '표준 키보드 문자 및 기호만 지원됩니다',
      placeholder: '입력하세요',
      submit: '전송',
      virtual: '키보드',
      ctrlaltdel: 'Ctrl+Alt+Del'
    },
    mouse: {
      default: '기본 커서',
      pointer: '포인터 커서',
      cell: '셀 커서',
      text: '텍스트 커서',
      grab: '잡기 커서',
      hide: '커서 숨기기',
      mode: '마우스 모드',
      absolute: '절대값 모드',
      relative: '상대값 모드',
      requestPointer: '상대값 모드를 사용 중입니다. 커서를 찾으려면 데스크톱을 클릭하세요.',
      resetHid: 'HID 초기화'
    },
    image: {
      title: '이미지',
      loading: '불러오는 중...',
      empty: '아무것도 없습니다.',
      mountFailed: '이미지 마운트 실패',
      mountDesc:
        '일부 시스템에서는 이미지를 마운트하기 전에 원격 호스트에서 가상 디스크를 제거해야 합니다.',
      unmountFailed: '마운트 해제 실패',
      unmountDesc:
        '일부 시스템에서는 이미지를 마운트 해제하기 전에 원격 호스트에서 수동으로 꺼내야 합니다.',
      refresh: '이미지 목록 새로고침'
    },
    script: {
      title: '스크립트',
      upload: '업로드',
      run: '실행',
      runBackground: '백그라운드에서 실행',
      runFailed: '실행 실패',
      attention: '주의',
      delDesc: '이 파일을 정말로 삭제합니까?',
      confirm: '네',
      cancel: '아니오',
      delete: '삭제',
      close: '닫기'
    },
    terminal: {
      title: '터미널',
      nanokvm: 'NanoKVM 터미널',
      serial: '시리얼 포트 터미널',
      serialPort: '시리얼 포트',
      serialPortPlaceholder: '시리얼 포트를 입력하세요',
      baudrate: 'Baud rate',
      confirm: '확인'
    },
    wol: {
      title: 'Wake-on-LAN',
      sending: '패킷 전송 중...',
      sent: '패킷 전송 완료',
      input: 'MAC주소를 입력하세요.',
      ok: '확인'
    },
    download: {
      title: '이미지 다운로드',
      input: '원격 이미지 URL을 입력하세요.',
      ok: '확인',
      disabled: '/data 파티션이 읽기 전용(RO) 상태이므로 이미지를 다운로드할 수 없습니다.'
    },
    power: {
      title: '전원',
      reset: '리셋',
      power: '전원',
      powerShort: '전원 (짧게 누르기)',
      powerLong: '전원 (길게 누르기)'
    },
    settings: {
      title: '설정',
      about: {
        title: 'NanoKVM 정보',
        information: '정보',
        ip: 'IP',
        mdns: 'mDNS',
        application: '펌웨어 버전',
        applicationTip: 'NanoKVM 웹 애플리케이션 버전',
        image: '이미지 버전',
        imageTip: 'NanoKVM 시스템 이미지 버전',
        deviceKey: '장치 키',
        community: '커뮤니티',
        pikvm: {
          title: 'PiKVM으로 전환',
          attention: '주의',
          desc1: '현재 시스템을 PiKVM로 전환하시겠습니까?',
          desc2: '시스템 전환 후 장치가 자동으로 재부팅됩니다.',
          confirm: '확인'
        }
      },
      appearance: {
        title: '디자인',
        display: '표시',
        language: '언어',
        menuBar: '메뉴 바',
        menuBarDesc: '메뉴 바에 아이콘을 표시'
      },
      device: {
        title: '장치',
        oled: {
          title: 'OLED',
          description: 'OLED 화면 자동 절전',
          0: '사용 안 함',
          15: '15초',
          30: '30초',
          60: '1분',
          180: '3분',
          300: '5분',
          600: '10분',
          1800: '30분',
          3600: '1시간'
        },
        wifi: {
          title: 'Wi-Fi',
          description: 'Wi-Fi 설정'
        },
        ssh: {
          description: 'SSH 원격 접속 활성화',
          tip: '활성화하기 전에 강력한 비밀번호를 설정하세요. (계정 - 비밀번호 변경)'
        },
        led: {
          title: 'LED 스트립',
          description: 'LED 스트립 설정'
        },
        disk: '가상 디스크',
        diskDesc: '원격 호스트에서 가상 U-디스크를 마운트합니다.',
        network: '가상 네트워크',
        networkDesc: '원격 호스트에서 가상 네트워크 카드를 마운트합니다.'
      },
      tailscale: {
        title: 'Tailscale',
        restart: '정말로 Tailscale을 다시 시작하시겠습니까?',
        stop: '정말로 Tailscale을 중지하시겠습니까?',
        stopDesc: 'Tailscale에서 로그아웃하고 자동 시작을 비활성화합니다.',
        loading: '불러오는 중...',
        notInstall: 'Tailscale이 없습니다. 설치해주세요.',
        install: '설치',
        installing: '설치중',
        failed: '설치 실패',
        retry: '새로고침하고 다시 시도하거나, 수동으로 설치하세요',
        download: '다운로드 중 :',
        package: '패키지 설치',
        unzip: '압축 해제',
        upTailscale: 'tailscale을 NanoKVM 의 다음 경로에 업로드 했습니다. : /usr/bin/',
        upTailscaled: 'tailscaled을 NanoKVM 의 다음 경로에 업로드 했습니다. :  /usr/sbin/',
        refresh: '현재 페이지 새로고침',
        notLogin: '이 기기는 현재 연동 되지 않았습니다. 로그인해서 계정에 이 장치를 연동하세요.',
        urlPeriod: '이 주소는 10분간 유효합니다.',
        login: '로그인',
        loginSuccess: '로그인 성공',
        enable: 'Tailscale 활성화',
        deviceName: '장치 이름',
        deviceIP: '장치 IP',
        account: '계정',
        logout: '로그아웃',
        uninstall: 'Tailscale 제거',
        okBtn: '네',
        cancelBtn: '아니오'
      },
      update: {
        title: '업데이트 확인',
        queryFailed: '버전 확인 실패',
        updateFailed: '업데이트 실패, 재시도하세요.',
        isLatest: '이미 최신 버전입니다.',
        available: '업데이트가 가능합니다. 정말로 업데이트 할까요?',
        updating: '업데이트 시작. 잠시 기다려주세요...',
        confirm: '확인',
        cancel: '취소'
      },
      account: {
        title: '계정',
        webAccount: 'Web 계정',
        password: '비밀번호',
        updateBtn: '업데이트',
        logoutBtn: '로그아웃'
      }
    },
    error: {
      title: '문제가 발생했습니다.',
      refresh: '새로고침'
    }
  }
};

export default ko;
