const en = {
  translation: {
    head: {
      desktop: 'Escritorio remoto',
      login: 'Inicio de sesión',
      changePassword: 'Cambiar contraseña',
      terminal: 'Consola',
      wifi: 'Wi-Fi'
    },
    auth: {
      login: 'Iniciar sesión',
      placeholderUsername: 'Introduce tu usuario',
      placeholderPassword: 'Introduce tu contraseña',
      placeholderPassword2: 'Introduce tu contraseña de nuevo',
      noEmptyUsername: 'El usuario no puede estar vacío',
      noEmptyPassword: 'La contraseña no puede estar vacía',
      noAccount:
        'No se ha encontrado la cuenta. Por favor, recarga la página o recupera tu contraseña.',
      invalidUser: 'Usuario o contraseña incorrectos',
      error: 'Error inesperado',
      changePassword: 'Cambiar contraseña',
      differentPassword: 'Las contraseñas no coinciden',
      changePasswordDesc:
        'Para la seguridad de su dispositivo, por favor modifique la contraseña de inicio de sesión en la web.',
      illegalUsername: 'El usuario contiene caracteres no permitidos',
      illegalPassword: 'La contraseña contiene caracteres no permitidos',
      forgetPassword: 'Contraseña olvidada',
      ok: 'Aceptar',
      cancel: 'Cancelar',
      loginButtonText: 'Iniciar sesión'
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
      video: 'Modo de vídeo',
      quality: 'Calidad',
      qualityLossless: 'Sin pérdida',
      qualityHigh: 'Alto',
      qualityMedium: 'Medio',
      qualityLow: 'Bajo',
      resetHdmi: 'Reset HDMI'
    },
    keyboard: {
      paste: 'Pegar',
      tips: 'Solo están soportadas las letras y símbolos estándar del teclado',
      placeholder: 'Por favor, introduce el texto',
      submit: 'Enviar',
      virtual: 'Teclado virtual',
      ctrlaltdel: 'Ctrl+Alt+Del'
    },
    mouse: {
      default: 'Cursor por defecto',
      pointer: 'Cursor de puntero',
      cell: 'Cursor de celda',
      text: 'Cursor de texto',
      grab: 'Cursor de agarre',
      hide: 'Ocultar cursor',
      mode: 'Modo de ratón',
      absolute: 'Modo absoluto',
      relative: 'Modo relativo',
      requestPointer:
        'Usando modo relativo. Por favor, haz clic en el escritorio para obtener el cursor del ratón.',
      resetHid: 'Restablecer HID'
    },
    image: {
      title: 'Imágenes',
      loading: 'Cargando...',
      empty: 'No se ha encontrado nada',
      mountFailed: 'Fallo al montar',
      mountDesc:
        'En algunos sistemas, es necesario expulsar el disco virtual en el host remoto antes de montar una imagen.',
      unmountFailed: 'Desmontaje fallido',
      unmountDesc:
        'En algunos sistemas, es necesario expulsar manualmente en el host remoto antes de desmontar la imagen.'
    },
    script: {
      title: 'Script',
      upload: 'Subir',
      run: 'Ejecutar',
      runBackground: 'Ejecutar en segundo plano',
      runFailed: 'Ejecución fallida',
      attention: 'Atención',
      delDesc: '¿Estás seguro que deseas eliminar este archivo?',
      confirm: 'Sí',
      cancel: 'No',
      delete: 'Eliminar',
      close: 'Cerrar'
    },
    terminal: {
      title: 'Consola',
      nanokvm: 'Consola del NanoKVM',
      serial: 'Consola del Puerto Serie',
      serialPort: 'Puerto Serie',
      serialPortPlaceholder: 'Por favor, introduce el puerto serie',
      baudrate: 'Tasa de baudios',
      confirm: 'Confirmar'
    },
    wol: {
      title: 'Wake-on-LAN',
      sending: 'Enviando comando...',
      sent: 'Comando enviado',
      input: 'Por favor, introduce la dirección MAC',
      ok: 'Vale'
    },
    power: {
      title: 'Encender / Apagar',
      reset: 'Reiniciar',
      power: 'Encender / Apagar',
      powerShort: 'Encender / Apagar (pulsación corta)',
      powerLong: 'Encender / Apagar (pulsación larga)'
    },
    settings: {
      title: 'Settings',
      about: {
        title: 'Sobre NanoKVM',
        information: 'Información',
        ip: 'IP',
        mdns: 'mDNS',
        application: 'Versión de la aplicación',
        applicationTip: 'NanoKVM web application version',
        image: 'Versión de la imagen',
        imageTip: 'NanoKVM system image version',
        deviceKey: 'Clave del dispositivo',
        community: 'Comunidad',
        pikvm: {
          title: 'Cambiar a PiKVM',
          attention: 'Atención',
          desc1: '¿Estás seguro de que quieres cambiar el sistema actual a PiKVM?',
          desc2: 'El dispositivo se reiniciará automáticamente después de cambiar de sistema.',
          confirm: 'Confirmar'
        }
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
        led: {
          title: 'Tira LED',
          description: 'Configurar tira LED'
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
        loading: 'Cargando...',
        notInstall: '¡Tailscale no encontrado! Por favor, instálalo.',
        install: 'Instalar',
        installing: 'Instalando',
        failed: 'La instalación falló',
        retry:
          'Por favor, actualiza la página e inténtalo de nuevo. O intenta instalarlo manualmente',
        download: 'Descargar el',
        package: 'paquete de instalación',
        unzip: 'y descomprimirlo',
        upTailscale: 'Sube tailscale al directorio /usr/bin/ del NanoKVM',
        upTailscaled: 'Sube tailscaled al directorio /usr/sbin/ del NanoKVM',
        refresh: 'Actualizar la página actual',
        notLogin:
          'El dispositivo aún no ha sido vinculado. Por favor, inicia sesión y vincula este dispositivo a tu cuenta.',
        urlPeriod: 'Esta URL es válida por 10 minutos',
        login: 'Iniciar sesión',
        loginSuccess: 'Inicio de sesión exitoso',
        enable: 'Habilitar Tailscale',
        deviceName: 'Nombre del dispositivo',
        deviceIP: 'IP del dispositivo',
        account: 'Cuenta',
        logout: 'Cerrar sesión',
        okBtn: 'Yes',
        cancelBtn: 'No'
      },
      update: {
        title: 'Buscar actualizaciones',
        queryFailed: 'Error al obtener la versión',
        updateFailed: 'La actualización falló. Por favor, inténtalo de nuevo.',
        isLatest: 'Ya tienes la última versión.',
        available: 'Hay una actualización disponible. ¿Estás seguro de que quieres actualizar?',
        updating: 'Actualización iniciada. Por favor, espera...',
        confirm: 'Confirmar',
        cancel: 'Cancelar'
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

export default en;
