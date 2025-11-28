const id = {
  translation: {
    head: {
      desktop: 'Remote Desktop',
      login: 'Masuk',
      changePassword: 'Ubah Sandi',
      terminal: 'Terminal',
      wifi: 'Wi-Fi'
    },
    auth: {
      login: 'Masuk',
      placeholderUsername: 'Silahkan masukkan username',
      placeholderPassword: 'Silahkan masukkan password',
      placeholderPassword2: 'Silahkan masukkan password again',
      noEmptyUsername: 'nama user tidak boleh kosong',
      noEmptyPassword: 'sandi  tidak boleh kosong',
      noAccount:
        'Gagal mendapatkan informasi user, silahkan segarkan halaman atau atur ulang sandi',
      invalidUser: 'invalid username or password',
      error: 'terjadi kesalahan tak terduga',
      changePassword: 'Ganti Sandi',
      changePasswordDesc: 'Untuk keamanan perangkat Anda, silakan ubah kata sandi masuk web.',
      differentPassword: 'sandi tidak sesuai',
      illegalUsername: 'ada karakter ilegal pada nama user',
      illegalPassword: 'ada karakter ilegal pada sandi',
      forgetPassword: 'Lupa Sandi',
      ok: 'Ok',
      cancel: 'Batalkan',
      loginButtonText: 'Masuk',
      tips: {
        reset1:
          'To reset the passwords, pressing and holding the BOOT button on the NanoKVM for 10 seconds.',
        reset2: 'For detailed steps, please consult this document:',
        reset3: 'Web default account:',
        change1: 'Please note that this action will change the following passwords:',
        change2: 'Web login password',
        change3: 'System root password (SSH login password)',
        change4: 'To reset the passwords, press and hold the BOOT button on the NanoKVM.'
      }
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
      video: 'Mode Video',
      quality: 'Kualitas',
      qualityLossless: 'Tanpa Kehilangan',
      qualityHigh: 'Tinggi',
      qualityMedium: 'Sedang',
      qualityLow: 'Rendah',
      resetHdmi: 'Reset HDMI'
    },
    keyboard: {
      paste: 'Tempel',
      tips: 'Hanya huruf dan simbol keyboard standar yang didukung',
      placeholder: 'Silahkan isi',
      submit: 'Kirimkan',
      virtual: 'Keyboard',
      ctrlaltdel: 'Ctrl+Alt+Del'
    },
    mouse: {
      default: 'Kursor bawaan',
      pointer: 'Kursor penunjuk',
      cell: 'Kursor cell',
      text: 'Kursor teks',
      grab: 'Kursor ambil',
      hide: 'Sembunyikan kursor',
      mode: 'Mode tetikus',
      absolute: 'Mode absolut',
      relative: 'Mode relatif',
      requestPointer:
        'Menggunakan mode relatf. Silakan klik desktop untuk mendapatkan penunjuk tetikus.',
      resetHid: 'Setel ulang HID'
    },
    image: {
      title: 'Gambar',
      loading: 'Memuat...',
      empty: 'Tidak ada yang ditemukan',
      mountFailed: 'Pemasangan Gagal',
      mountDesc:
        'Di beberapa sistem, perlu mengeluarkan disk virtual pada host jarak jauh sebelum memasang gambar.',
      unmountFailed: 'Pelepasan gagal',
      unmountDesc:
        'Di beberapa sistem, Anda perlu mengeluarkan secara manual di host jarak jauh sebelum melepaskan gambar.'
    },
    script: {
      title: 'Script',
      upload: 'Mengunggah',
      run: 'Jalankan',
      runBackground: 'Jalankan di belakang',
      runFailed: 'Gagal menjalankan',
      attention: 'Perhatian',
      delDesc: 'Apa kamu yakin menghapus data ini?',
      confirm: 'Ya',
      cancel: 'Tidak',
      delete: 'Hapus',
      close: 'Tutup'
    },
    terminal: {
      title: 'Terminal',
      nanokvm: 'Terminal NanoKVM',
      serial: 'Terminal Port Serial',
      serialPort: 'Port Serial',
      serialPortPlaceholder: 'Silahkan masukkan port serial',
      baudrate: 'Baud rate',
      confirm: 'Ok'
    },
    wol: {
      title: 'Wake-on-LAN',
      sending: 'Kirim perintah...',
      sent: 'Perintah terkirim',
      input: 'Silahkan masukkan MAC',
      ok: 'Ok'
    },
    power: {
      title: 'Daya',
      reset: 'Mulai Ulang',
      power: 'Daya',
      powerShort: 'Data (tekan sebentar)',
      powerLong: 'Power (tekan lama)'
    },
    settings: {
      title: 'Settings',
      about: {
        title: 'Tentang NanoKVM',
        information: 'Informasi',
        ip: 'IP',
        mdns: 'mDNS',
        application: 'Versi Aplikasi',
        applicationTip: 'NanoKVM web application version',
        image: 'Version Gambar',
        imageTip: 'NanoKVM system image version',
        deviceKey: 'Kunci Perangkat',
        community: 'Komunitas'
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
        loading: 'Memuat...',
        notInstall: 'Tailscale tidak ditemukan! Silahkan pasang.',
        install: 'Memasang',
        installing: 'Memasangkan',
        failed: 'Gagal memasangkan',
        retry: 'Harap segarkan dan coba lagi. Atau coba instal secara manual',
        download: 'Mengunduh',
        package: 'paket instalasi',
        unzip: 'dan unzip itu',
        upTailscale: 'Unggah tailscale ke direktori NanoKVM /usr/bin/',
        upTailscaled: 'Unggah tailscaled ke direktori NanoKVM /usr/sbin/',
        refresh: 'Segarkan halaman ini',
        notLogin:
          'Perangkat belum ditautkan. Silakan masuk dan tautkan perangkat ini ke akun Anda.',
        urlPeriod: 'Url ini berlaku selama 10 menit',
        login: 'Masuk',
        loginSuccess: 'Berhasil masuk',
        enable: 'Aktifkan Tailscale',
        deviceName: 'Nama Perangkat',
        deviceIP: 'IP Perangkat',
        account: 'Akun',
        logout: 'Keluar',
        okBtn: 'Yes',
        cancelBtn: 'No'
      },
      update: {
        title: 'Periksa pembaruan',
        queryFailed: 'Gagal mendapatkan versi',
        updateFailed: 'Gagal memperbarui, tolong coba lagi.',
        isLatest: 'Kamu sudah menggunakan versi terbaru.',
        available: 'Ada pembaruan baru. apa kamu mau memperbarui?',
        updating: 'Pembaruan dimulai. Silahkan tunggu...',
        confirm: 'Konfirmasi',
        cancel: 'Batalkan'
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

export default id;
