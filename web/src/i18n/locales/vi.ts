const vi = {
  translation: {
    head: {
      desktop: 'Remote Desktop',
      login: 'Đăng nhập',
      changePassword: 'Đổi mật khẩu',
      terminal: 'Terminal',
      wifi: 'Wi-Fi'
    },
    auth: {
      login: 'Đăng nhập',
      placeholderUsername: 'Vui lòng nhập tên người dùng',
      placeholderPassword: 'vui lòng nhập mật khẩu',
      placeholderPassword2: 'vui lòng nhập lại mật khẩu',
      noEmptyUsername: 'tên người dùng không được để trống',
      noEmptyPassword: 'mật khẩu không được để trống',
      noAccount:
        'Không thể lấy thông tin người dùng, vui lòng làm mới trang web hoặc đặt lại mật khẩu',
      invalidUser: 'tên người dùng hoặc mật khẩu không hợp lệ',
      error: 'lỗi không mong đợi',
      changePassword: 'Đổi mật khẩu',
      changePasswordDesc: 'Để bảo mật thiết bị của bạn, vui lòng thay đổi mật khẩu đăng nhập web.',
      differentPassword: 'mật khẩu không khớp',
      illegalUsername: 'tên người dùng chứa ký tự không hợp lệ',
      illegalPassword: 'mật khẩu chứa ký tự không hợp lệ',
      forgetPassword: 'Quên mật khẩu',
      ok: 'OK',
      cancel: 'Hủy',
      loginButtonText: 'Đăng nhập'
    },
    wifi: {
      title: 'Wi-Fi',
      description: 'Configure Wi-Fi for NanoKVM',
      success: 'Please check the network status of NanoKVM and visit the new IP address.',
      failed: 'Operation failed, please try again.',
      confirmBtn: 'Ok',
      finishBtn: 'Finished'
    },
    screen: {
      video: 'Chế độ video',
      quality: 'Chất lượng',
      qualityLossless: 'Không mất dữ liệu',
      qualityHigh: 'Cao',
      qualityMedium: 'Trung bình',
      qualityLow: 'Thấp',
      resetHdmi: 'Reset HDMI'
    },
    keyboard: {
      paste: 'Dán',
      tips: 'Chỉ hỗ trợ các chữ cái và ký hiệu bàn phím tiêu chuẩn',
      placeholder: 'Vui lòng nhập',
      submit: 'Gửi',
      virtual: 'Bàn phím',
      ctrlaltdel: 'Ctrl+Alt+Del'
    },
    mouse: {
      default: 'Con trỏ mặc định',
      pointer: 'Con trỏ trỏ',
      cell: 'Con trỏ ô',
      text: 'Con trỏ văn bản',
      grab: 'Con trỏ nắm',
      hide: 'Ẩn con trỏ',
      mode: 'Chế độ chuột',
      absolute: 'Chế độ tuyệt đối',
      relative: 'Chế độ tương đối',
      requestPointer:
        'Đang sử dụng chế độ tương đối. Vui lòng nhấp vào màn hình để lấy con trỏ chuột.',
      resetHid: 'Đặt lại HID'
    },
    image: {
      title: 'Hình ảnh',
      loading: 'Đang tải...',
      empty: 'Không tìm thấy',
      mountFailed: 'Mount thất bại',
      mountDesc: 'Trong một số hệ thống, cần phải eject đĩa ảo trên máy remote trước khi mount.',
      unmountFailed: 'Gỡ bỏ thất bại',
      unmountDesc:
        'Trong một số hệ thống, bạn cần phải tháo gỡ thủ công trên máy chủ từ xa trước khi gỡ bỏ ảnh.'
    },
    script: {
      title: 'Script',
      upload: 'Tải lên',
      run: 'Chạy',
      runBackground: 'Chạy nền',
      runFailed: 'Chạy thất bại',
      attention: 'Chú ý',
      delDesc: 'Bạn có chắc chắn muốn xóa tệp này không?',
      confirm: 'Có',
      cancel: 'Không',
      delete: 'Xóa',
      close: 'Đóng'
    },
    terminal: {
      title: 'Terminal',
      nanokvm: 'Terminal NanoKVM',
      serial: 'Terminal Cổng Nối Tiếp',
      serialPort: 'Cổng Nối Tiếp',
      serialPortPlaceholder: 'Vui lòng nhập cổng nối tiếp',
      baudrate: 'Tốc độ Baud',
      confirm: 'OK'
    },
    wol: {
      title: 'Wake-on-LAN',
      sending: 'Đang gửi lệnh...',
      sent: 'Đã gửi lệnh',
      input: 'Vui lòng nhập địa chỉ MAC',
      ok: 'OK'
    },
    power: {
      title: 'Nguồn',
      reset: 'Đặt lại',
      power: 'Nguồn',
      powerShort: 'Nguồn (nhấp ngắn)',
      powerLong: 'Nguồn (nhấp dài)'
    },
    settings: {
      title: 'Settings',
      about: {
        title: 'Giới thiệu về NanoKVM',
        information: 'Thông tin',
        ip: 'IP',
        mdns: 'mDNS',
        application: 'Phiên bản Ứng dụng',
        applicationTip: 'NanoKVM web application version',
        image: 'Phiên bản Hình ảnh',
        imageTip: 'NanoKVM system image version',
        deviceKey: 'Khóa Thiết bị',
        community: 'Cộng đồng'
      },
      appearance: {
        title: 'Appearance',
        display: 'Display',
        language: 'Language',
        menuBar: 'Menu Bar',
        menuBarDesc: 'Display icons in the menu bar'
      },
      device: {
        title: 'Device',
        oled: {
          title: 'OLED',
          description: 'OLED screen automatically sleep',
          0: 'Never',
          15: '15 sec',
          30: '30 sec',
          60: '1 min',
          180: '3 min',
          300: '5 min',
          600: '10 min',
          1800: '30 min',
          3600: '1 hour'
        },
        wifi: {
          title: 'Wi-Fi',
          description: 'Configure Wi-Fi'
        },
        disk: 'Virtual Disk',
        diskDesc: 'Mount virtual U-disk on the remote host',
        network: 'Virtual Network',
        networkDesc: 'Mount virtual network card on the remote host'
      },
      tailscale: {
        title: 'Tailscale',
        restart: 'Are you sure to restart Tailscale?',
        stop: 'Are you sure to stop Tailscale?',
        stopDesc: 'Log out Tailscale and disable its automatic startup on boot.',
        loading: 'Đang tải...',
        notInstall: 'Không tìm thấy Tailscale! Vui lòng cài đặt.',
        install: 'Cài đặt',
        installing: 'Đang cài đặt',
        failed: 'Cài đặt thất bại',
        retry: 'Vui lòng làm mới và thử lại. Hoặc thử cài đặt thủ công',
        download: 'Tải xuống',
        package: 'gói cài đặt',
        unzip: 'và giải nén nó',
        upTailscale: 'Tải tailscale lên thư mục /usr/bin/ của NanoKVM',
        upTailscaled: 'Tải tailscaled lên thư mục /usr/sbin/ của NanoKVM',
        refresh: 'Làm mới trang hiện tại',
        notLogin:
          'Thiết bị chưa được liên kết. Vui lòng đăng nhập và liên kết thiết bị này với tài khoản của bạn.',
        urlPeriod: 'URL này có hiệu lực trong 10 phút',
        login: 'Đăng nhập',
        loginSuccess: 'Đăng nhập thành công',
        enable: 'Kích hoạt Tailscale',
        deviceName: 'Tên Thiết bị',
        deviceIP: 'IP Thiết bị',
        account: 'Tài khoản',
        logout: 'Đăng xuất',
        okBtn: 'Yes',
        cancelBtn: 'No'
      },
      update: {
        title: 'Kiểm tra cập nhật',
        queryFailed: 'Lấy phiên bản thất bại',
        updateFailed: 'Cập nhật thất bại. Vui lòng thử lại.',
        isLatest: 'Bạn đã có phiên bản mới nhất.',
        available: 'Có bản cập nhật mới. Bạn có chắc chắn muốn cập nhật không?',
        updating: 'Bắt đầu cập nhật. Vui lòng chờ...',
        confirm: 'Xác nhận',
        cancel: 'Hủy'
      },
      account: {
        title: 'Account',
        webAccount: 'Web Account Name',
        password: 'Password',
        updateBtn: 'Update',
        logoutBtn: 'Logout'
      }
    }
  }
};

export default vi;
