const en = {
  translation: {
    head: {
      desktop: 'Remote Desktop',
      login: 'Login',
      changePassword: 'Change Password',
      terminal: 'Terminal',
      wifi: 'Wi-Fi'
    },
    auth: {
      login: 'Login',
      placeholderUsername: 'Username',
      placeholderPassword: 'Password',
      placeholderPassword2: 'Please enter password again',
      noEmptyUsername: 'Username required',
      noEmptyPassword: 'Password required',
      noAccount: 'Failed to get user information, please refresh web page or reset password',
      invalidUser: 'Invalid username or password',
      error: 'Unexpected error',
      changePassword: 'Change Password',
      changePasswordDesc: 'For the security of your device, please change the password!',
      differentPassword: 'Passwords do not match',
      illegalUsername: 'Username contains illegal characters',
      illegalPassword: 'Password contains illegal characters',
      forgetPassword: 'Forgot Password',
      ok: 'Ok',
      cancel: 'Cancel',
      loginButtonText: 'Login',
      tips: {
        reset1: 'Web default account:',
        reset2: 'SSH default account:',
        reset3: 'To reset the passwords, please refer to ',
        change1: 'Please note that this action will change the following passwords:',
        change2: 'Web login password',
        change3: 'System root password (SSH login password)',
        change4: 'To reset the passwords, please refer to '
      }
    },
    assistant: {
      title: 'Smart Assistant',
      tips: {
        intro:
          'Smart Assistant is an experimental feature of NanoKVM that allows you to experience using a Vision-Language Model (VLM) for automated computer operations, often called "Computer Use."\nBefore using this feature, please set the video mode to MJPEG transfer and your screen resolution to 720p or lower.\nThe Smart Assistant currently relies on an external VLM model service, which requires you to provide the corresponding model server address and API key. \nThis means the model server you specify will receive screenshots of your computer screen, and you may incur related token fees. \nYou can also set up your own VLM model server on your local network to avoid these issues.',
        content:
          "Large models still have an unpredictable risk of hallucinations and limited operational capabilities, which may lead to unintended actions that could damage your computer's data. You must be aware of this risk. While the Smart Assistant is operating, you must monitor its actions and be ready to interrupt any potentially risky operations at any time!\nClicking the button below will open the Smart Assistant in a new window, and by doing so, you acknowledge and accept the risks mentioned above.\nFor security reasons, this service is currently limited to a single user. It automatically shutdown service when the page is closed, requiring you to click the button again to start a new session.",
        contribute:
          'The code has been open-sourced, so feel free to review it. Contributions are also welcome!',
        warning: 'Warning'
      },
      install: 'Install Dependencies',
      tryNow: 'Try It Now'
    },
    wifi: {
      title: 'Wi-Fi',
      description: 'Configure Wi-Fi for NanoKVM',
      success: 'Please go to the device to check the network status of NanoKVM.',
      failed: 'Operation failed, please try again.',
      invalidMode:
        'The current mode does not support network setup. Please go to your device and enable Wi-Fi configuration mode.',
      confirmBtn: 'Ok',
      finishBtn: 'Finished'
    },
    notification: {
      browser: {
        title: 'Browser Recommendation',
        description:
          'For the best user experience, we recommend using the Chrome browser. While we strive for compatibility across all major browsers, some advanced features may not work as intended outside of Chrome.'
      },
      h265: {
        title: 'Warning',
        description:
          "The device's factory-installed firmware is strictly provisioned without H.265/HEVC video encoding or decoding functionality.\n\n" +
          'Any user-initiated modification of the firmware, manual installation of third-party H.265 modules, or independent development of HEVC features constitutes an assumption of all related risks by the end-user. These risks include, but are not limited to, patent licensing fees, legal liabilities, device compatibility issues, and performance degradation.\n\n' +
          'The manufacturer expressly disclaims all responsibility for the operation, security, or legal compliance of any user-customized H.265 features.',
        notSupported: 'The current browser does not support H.265!'
      },
      btnOk: 'OK'
    },
    screen: {
      title: 'Screen',
      video: 'Video Mode',
      quality: 'Quality',
      qualityAuto: 'Automatic',
      qualityLossless: 'Lossless',
      qualityHigh: 'High',
      qualityMedium: 'Medium',
      qualityLow: 'Low',
      scale: 'Scale',
      resetHdmi: 'Reset HDMI',
      noSignal: 'HDMI no signal',
      inconsistentVideoMode: 'Play paused. Another video mode is playing.',
      upload: 'Upload',
      advanced: 'Advanced'
    },
    volume: {
      title: 'Volume'
    },
    microphone: {
      title: 'Microphone',
      notMounted: 'Microphone Not Mounted',
      tip1: 'The virtual microphone is not mounted on the remote host. Would you like to mount and enable it now?',
      tip2: 'You can also manage this later in "Settings > Device > Virtual Mic".',
      confirm: 'Enable',
      cancel: 'Cancel'
    },
    keyboard: {
      title: 'Keyboard',
      paste: 'Paste',
      tips: 'Only standard keyboard letters and symbols are supported',
      placeholder: 'Please input',
      submit: 'Submit',
      virtual: 'Keyboard',
      shortcut: {
        title: 'Shortcuts',
        custom: 'Custom',
        capture: 'Click here to capture shortcut',
        clear: 'Clear',
        save: 'Save',
        captureTips:
          'Capturing system-level keys (such as the Windows key) requires full-screen permission.',
        enterFullScreen: 'Toggle full-screen mode.'
      },
      leaderKey: {
        title: 'Leader Key',
        desc: 'Bypass browser restrictions and send system shortcuts directly to the remote host.',
        howToUse: 'How to Use',
        simultaneous: {
          title: 'Simultaneous Mode',
          desc1: 'Press and hold the Leader Key, then press the shortcut.',
          desc2: 'Intuitive, but may conflict with system shortcuts.'
        },
        sequential: {
          title: 'Sequential Mode',
          desc1:
            'Press the Leader Key → press the shortcut in sequence → press the Leader Key again.',
          desc2: 'Requires more steps, but completely avoids system conflicts.'
        },
        enable: 'Enable Leader Key',
        placeholder: 'Please press the Leader Key',
        shiftRight: 'Right Shift',
        ctrlRight: 'Right Ctrl',
        metaRight: 'Right Win',
        submit: 'Submit',
        recorder: {
          rec: 'REC',
          activate: 'Activate keys',
          input: 'Please press the shortcut...'
        }
      }
    },
    mouse: {
      title: 'Mouse',
      cursor: 'Cursor style',
      default: 'Default cursor',
      pointer: 'Pointer cursor',
      cell: 'Cell cursor',
      text: 'Text cursor',
      grab: 'Grab cursor',
      hide: 'Hide cursor',
      mode: 'Mouse mode',
      absolute: 'Absolute mode',
      relative: 'Relative mode',
      direction: 'Wheel direction',
      scrollUp: 'Scroll up',
      scrollDown: 'Scroll down',
      speed: 'Wheel speed',
      fast: 'Fast',
      slow: 'Slow',
      requestPointer: 'Using relative mode. Please click desktop to get mouse pointer.',
      resetHid: 'Reset HID',
      hidOnly: {
        title: 'HID-Only mode',
        desc: "If your mouse and keyboard stop responding and resetting HID doesn't help, it could be a compatibility issue between the NanoKVM and the device. Try to enable HID-Only mode for better compatibility.",
        tip1: 'Enabling HID-Only mode will unmount the virtual U-disk and virtual network',
        tip2: 'In HID-Only mode, image mounting is disabled',
        enable: 'Enable HID-Only mode',
        disable: 'Disable HID-Only mode'
      }
    },
    image: {
      title: 'Images',
      loading: 'Loading...',
      empty: 'Nothing Found',
      upload: 'Upload Image',
      uploaded: 'Image Uploaded',
      mount: 'Mount Image',
      mounted: 'Mounted',
      fetchError: 'Failed to fetch image list, please refresh and try again',
      uploadError: 'Upload failed, please try again',
      mountMode: 'Mount mode',
      readOnly: 'Read only',
      mountFailed: 'Mount failed',
      mountDesc:
        'On some systems, you need to eject the virtual hard drive from the remote host before mounting the image.',
      unmountFailed: 'Unmount failed',
      unmountDesc:
        'On some systems, you need to manually eject from the remote host before unmounting the image.',
      refresh: 'Refresh the image list',
      download: 'Download Image',
      attention: 'Attention',
      deleteConfirm: 'Are you sure you want to delete this image?',
      okBtn: 'Yes',
      cancelBtn: 'No'
    },
    script: {
      title: 'Scripts',
      upload: 'Upload',
      run: 'Run',
      runBackground: 'Run Background',
      runFailed: 'Run failed',
      attention: 'Attention',
      delDesc: 'Are you sure you want to delete this file?',
      confirm: 'Yes',
      cancel: 'No',
      delete: 'Delete',
      close: 'Close'
    },
    terminal: {
      title: 'Terminal',
      nanokvm: 'NanoKVM Terminal',
      serial: 'Serial Port Terminal',
      serialPort: 'Serial Port',
      serialPortPlaceholder: 'Please enter the serial port',
      baudrate: 'Baud rate',
      parity: 'Parity',
      parityNone: 'None',
      parityEven: 'Even',
      parityOdd: 'Odd',
      flowControl: 'Flow control',
      flowControlNone: 'None',
      flowControlSoft: 'Soft',
      flowControlHard: 'Hard',
      dataBits: 'Data bits',
      stopBits: 'Stop bits',
      confirm: 'Ok'
    },
    wol: {
      title: 'Wake-on-LAN',
      sending: 'Sending command...',
      sent: 'Command sent',
      input: 'Please enter the MAC',
      ok: 'Ok'
    },
    download: {
      title: 'Image Downloader',
      input: 'Please enter a remote image URL',
      ok: 'Ok',
      disabled: '/data partition is RO, so we cannot download the image'
    },
    power: {
      title: 'Power',
      showConfirm: 'Confirmation',
      showConfirmTip: 'Power operations require an extra confirmation',
      reset: 'Reset',
      power: 'Power',
      powerShort: 'Power (short click)',
      powerLong: 'Power (long click)',
      resetConfirm: 'Proceed reset operation?',
      powerConfirm: 'Proceed power operation?',
      okBtn: 'Yes',
      cancelBtn: 'No'
    },
    settings: {
      title: 'Settings',
      about: {
        title: 'About NanoKVM',
        information: 'Information',
        ip: 'IP',
        mdns: 'mDNS',
        application: 'Application Version',
        applicationTip: 'NanoKVM web application version',
        arch: 'Arch',
        device_number: 'Device Number',
        image: 'Image Version',
        imageTip: 'NanoKVM system image version',
        deviceKey: 'Device Key',
        community: 'Community',
        hostname: 'Hostname',
        hostnameUpdated: 'Hostname updated. Reboot to apply.',
        pikvm: {
          title: 'Switch to PiKVM',
          attention: 'Attention',
          desc1: 'Are you sure you want to switch the current system to PiKVM?',
          desc2: 'The device will automatically reboot after switching systems.',
          confirm: 'Confirm'
        },
        ipType: {
          Wired: 'Wired',
          Wireless: 'Wireless',
          Other: 'Other'
        }
      },
      appearance: {
        title: 'Appearance',
        display: 'Display',
        language: 'Language',
        languageDesc: 'Select the language for the interface',
        webTitle: 'Web Title',
        webTitleDesc: 'Customize the web page title',
        menuBar: {
          title: 'Menu Bar',
          mode: 'Display Mode',
          modeDesc: 'Display menu bar on the screen',
          modeOff: 'Off',
          modeAuto: 'Auto hide',
          modeAlways: 'Always visible',
          icons: 'Submenu Icons',
          iconsDesc: 'Display submenu icons in the menu bar'
        }
      },
      screen: {
        title: 'Screen',
        videoMode: {
          title: 'Video Mode',
          description: 'Select the encoding and transmission method'
        },
        edid: {
          description: 'Select display parameters',
          default: 'Default',
          custom: 'Custom',
          document: 'Document'
        },
        rateControlMode: {
          title: 'Rate Control Mode',
          description: 'Select the method to control the bitrate',
          cbr: 'CBR (Constant Bitrate): Prioritizes connection stability. Minimizes stuttering, but image quality may drop during high-motion scenes.',
          vbr: 'VBR (Variable Bitrate): Prioritizes image clarity. High-motion scenes may cause bandwidth spikes, potentially leading to latency or stutter.'
        },
        bitrate: {
          title: 'Bitrate',
          description: 'Balance image clarity with bandwidth usage',
          lower:
            'Lower values reduce latency on poor networks but compromise image quality, potentially causing blur or artifacts.',
          higher:
            'Higher values improve image clarity but consume more bandwidth. May increase latency if the connection is slow.'
        },
        gop: {
          description: 'Set the interval between full keyframes',
          lower:
            'Lower values send full frames more frequently. This speeds up recovery from visual artifacts but consumes more bandwidth and reduces image detail.',
          higher:
            'Higher values optimize for incremental updates. This improves image detail and saves bandwidth, but recovery from visual artifacts will take longer.'
        },
        quality: {
          title: 'Quality',
          description: 'Set the display image quality'
        },
        fps: {
          auto: 'Auto',
          description: 'Set the display frame rate',
          tip1: "Setting the value to 0 enables Auto mode, synchronizing the display with the remote host's frame rate.",
          tip2: 'Lowering the frame rate significantly reduces bandwidth usage. Manually decrease it during network instability or when high frame rates are not required.'
        },
        scale: {
          title: 'Scale',
          description: 'Set the display scaling ratio'
        }
      },
      device: {
        title: 'Device',
        oled: {
          title: 'OLED',
          description: 'Turn off OLED screen after',
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
          description: 'Configure Wi-Fi',
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
          description: 'Enable HDMI passthrough',
          tip: 'Enabling HDMI capture limits passthrough to 4K30. Disable HDMI capture to achieve 4K60 passthrough.'
        },
        capture: {
          title: 'HDMI',
          description: 'Enable HDMI output/capture',
          tip: 'Enabling HDMI passthrough limits the output frame rate to standard 4K30. Disable passthrough to support non-standard rates like 4K40. (Monitor compatibility may vary).'
        },
        ssh: {
          description: 'Enable SSH remote access',
          tip: 'Set a strong password before enabling (Account - Change Password)'
        },
        advanced: 'Advanced Settings',
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
        lowPower: {
          title: 'Low Power Mode',
          tips: 'Tips',
          tip1: 'When enabled, HDMI output will be automatically turned off during idle periods to reduce power consumption and heat.',
          tip2: 'However, waking up HDMI requires re-initialization. You may experience a HDMI No Signal state for approximately 10 seconds.',
          okBtn: 'Ok',
          cancelBtn: 'Cancel'
        },
        mouseJiggler: {
          title: 'Mouse Jiggler',
          description: 'Prevent the remote host from sleeping',
          disable: 'Disable',
          absolute: 'Absolute Mode',
          relative: 'Relative Mode'
        },
        mdns: {
          description: 'Enable mDNS discovery service'
        },
        hidOnly: 'HID-Only Mode',
        hidOnlyDesc: 'HID-Only Mode is enabled and virtual devices are currently unavailable',
        imageMounted: 'Image Mounted',
        imageMountedDesc: 'Operating the virtual devices will unmount the current image',
        disk: 'Virtual Disk',
        diskDesc: 'Mount virtual U-disk on the remote host',
        emmc: {
          warning: 'Warning',
          tip1: "Mounting the eMMC for the first time will wipe all data in the NanoKVM's /data directory. Please ensure you have a backup before proceeding.",
          tip2: 'Do you want to continue?'
        },
        network: 'Virtual Network',
        networkDesc: 'Mount virtual network card on the remote host',
        mic: 'Virtual Mic',
        micDesc: 'Mount virtual microphone on the remote host',
        micTip:
          'Due to USB endpoint limitations, enabling the Virtual Mic will disable the USB Secondary Screen.\n' +
          'And it is recommended to disable this feature during BIOS operations to ensure optimal compatibility.',
        reboot: 'Reboot',
        rebootDesc: 'Are you sure you want to reboot NanoKVM?',
        datetime: {
          timezone: 'Time zone',
          datetime: 'Datetime',
          format: 'Time format',
          synchronize: 'Synchronize time',
          lastSynchronization: 'Last synchronization: ',
          notSynchronized: 'The system time is not synchronized yet',
          syncNow: 'Sync now'
        },
        okBtn: 'Yes',
        cancelBtn: 'No'
      },
      tailscale: {
        title: 'Tailscale',
        restart: 'Restart Tailscale?',
        stop: 'Stop Tailscale?',
        stopDesc: 'Log out Tailscale and disable automatic startup on boot.',
        loading: 'Loading...',
        notInstall: 'Tailscale not found! Please install.',
        install: 'Install',
        installing: 'Installing',
        failed: 'Install failed',
        retry: 'Please refresh and try again. Or try to install manually',
        download: 'Download the',
        package: 'installation package',
        unzip: 'and unzip it',
        upTailscale: 'Upload tailscale to NanoKVM directory /usr/bin/',
        upTailscaled: 'Upload tailscaled to NanoKVM directory /usr/sbin/',
        refresh: 'Refresh current page',
        notRunning: 'Tailscale is not running. Please start it to continue.',
        run: 'Start',
        notLogin:
          'The device has not been bound yet. Please login and bind this device to your account.',
        urlPeriod: 'This url is valid for 10 minutes',
        login: 'Login',
        loginSuccess: 'Login Success',
        enable: 'Enable Tailscale',
        deviceName: 'Device Name',
        deviceIP: 'Device IP',
        account: 'Account',
        logout: 'Logout',
        logoutDesc: 'Are you sure you want to logout?',
        uninstall: 'Uninstall Tailscale',
        uninstallDesc: 'Are you sure you want to uninstall Tailscale?',
        okBtn: 'Yes',
        cancelBtn: 'No'
      },
      update: {
        title: 'Check for Updates',
        queryFailed: 'Get version failed',
        updateFailed: 'Update failed. Please retry.',
        isLatest: 'You already have the latest version.',
        available: 'An update is available. Are you sure you want to update now?',
        updating: 'Update started. Please wait...',
        confirm: 'Confirm',
        upgrade_tip_nano:
          'A new version of nanokvm is available. Are you sure you want to update now?',
        upgrade_tip_firmware:
          'A new firmware version is available. Are you sure you want to update now?',
        cancel: 'Cancel',
        preview: 'Preview Updates',
        previewDesc: 'Get early access to new features and improvements',
        previewTip:
          'Please be aware that preview releases may contain bugs or incomplete functionality!',
        download: 'Download',
        install: 'Install',
        changelog: 'Changelog',
        restart: 'Restart'
      },
      account: {
        title: 'Account',
        webAccount: 'Web Account Name',
        password: 'Password',
        updateBtn: 'Change',
        logoutBtn: 'Logout',
        logoutDesc: 'Are you sure you want to logout?',
        okBtn: 'Yes',
        cancelBtn: 'No'
      },
      kvmadmin: {
        title: 'NanoKVM Admin',
        description: 'Batch management of NanoKVM devices',
        tip: 'This feature uses the mDNS protocol. Please make sure the mDNS is enabled on your device.',
        install: 'Install',
        start: 'Run',
        visit: 'Open',
        uninstall: 'Uninstall',
        attention: 'Attention',
        confirmUninstall: 'Are you sure to uninstall NanoKVM Admin?',
        clearData: 'This operation will delete all saved data!',
        installFailed: 'Install failed. Please try again.',
        startFailed: 'Service startup failed. Please try again.',
        okBtn: 'Yes',
        cancelBtn: 'No'
      }
    },
    error: {
      title: "We've ran into an issue",
      refresh: 'Refresh'
    },
    fullscreen: {
      toggle: 'Toggle Fullscreen'
    },
    menu: {
      collapse: 'Collapse Menu',
      expand: 'Expand Menu'
    }
  }
};

export default en;
