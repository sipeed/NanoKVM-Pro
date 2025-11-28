const de = {
  translation: {
    head: {
      desktop: 'Entfernter Desktop',
      login: 'Anmelden',
      changePassword: 'Passwort ändern',
      terminal: 'Terminal',
      wifi: 'Wi-Fi'
    },
    auth: {
      login: 'Anmelden',
      placeholderUsername: 'Benutzername',
      placeholderPassword: 'Passwort',
      placeholderPassword2: 'Bitte Passwort erneut eingeben',
      noEmptyUsername: 'Benutzername benötigt',
      noEmptyPassword: 'Passwort benötigt',
      noAccount:
        'Abfragen der Benutzerdaten fehlgeschlagen, bitte die Seite neu laden oder Passwort zurücksetzen',
      invalidUser: 'Falscher Benutzername oder falsches Passwort',
      error: 'Unerwarteter Fehler',
      changePassword: 'Passwort ändern',
      changePasswordDesc: 'Für die Sicherheit Ihres Geräts ändern Sie bitte das Passwort!',
      differentPassword: 'Passwörter stimmen nicht überein',
      illegalUsername: 'Benutzername enthält ungültige Zeichen',
      illegalPassword: 'Passwort enthält ungültige Zeichen',
      forgetPassword: 'Passwort vergessen',
      ok: 'Ok',
      cancel: 'Abbrechen',
      loginButtonText: 'Anmelden'
    },
    wifi: {
      title: 'Wi-Fi',
      description: 'Wi-Fi Konfiguration für NanoKVM',
      success:
        'Bitte überprüfen Sie den Netzwerk-Status des NanoKVM und greifen Sie über die neue IP Adresse darauf zu.',
      failed: 'Aktion fehlgeschlagen, bitte erneut versuchen.',
      confirmBtn: 'Ok',
      finishBtn: 'Fertig'
    },
    screen: {
      title: 'Bildschirm',
      video: 'Video Modus',
      quality: 'Qualität',
      qualityLossless: 'Verlustfrei',
      qualityHigh: 'Hoch',
      qualityMedium: 'Mittel',
      qualityLow: 'Niedrig',
      resetHdmi: 'HDMI zurücksetzen'
    },
    keyboard: {
      title: 'Tastatur',
      paste: 'Einfügen',
      tips: 'Nur Standard-Tastaturbuchstaben und Symbole werden unterstützt',
      placeholder: 'Bitte eingeben',
      submit: 'Senden',
      virtual: 'Tastatur',
      ctrlaltdel: 'Ctrl+Alt+Del'
    },
    mouse: {
      title: 'Maus',
      cursor: 'Cursor',
      default: 'Standard Cursor',
      pointer: 'Zeiger Cursor',
      cell: 'Zellen Cursor',
      text: 'Text Cursor',
      grab: 'Greif Cursor',
      hide: 'Versteckter Cursor',
      mode: 'Maus Modus',
      absolute: 'Absoluter Modus',
      relative: 'Relativer Modus',
      requestPointer:
        'Relativer Modus aktiv. Klicken Sie auf den Desktop um den Mauszeiger zu sehen.',
      resetHid: 'HID zurücksetzen',
      hidOnly: {
        title: 'HID-Only Modus',
        desc: 'Wenn Ihre Maus und Tastatur nicht mehr reagieren und das Zurücksetzen der HID-Verbindung nicht hilft, könnte es sich um ein Kompatibilitätsproblem zwischen dem NanoKVM und dem Gerät handeln. Versuchen Sie, den HID-Only Modus zu aktivieren, um die Kompatibilität zu verbessern.',
        tip1: 'Die Aktivierung des HID-Only Modus entfernt das virtuelle U-Laufwerk und das virtuelle Netzwerk.',
        tip2: 'Im HID-Only Modus ist das Einbinden von Systemabbilder deaktiviert.',
        enable: 'HID-Only Modus aktivieren',
        disable: 'HID-Only Modus deaktivieren'
      }
    },
    image: {
      title: 'Bilder',
      loading: 'Lädt...',
      empty: 'Nichts gefunden',
      mountFailed: 'Einbinden fehlgeschlagen',
      mountDesc:
        'In einigen Systemen ist es notwendig, die virtuelle Festplatte auf dem entfernten Host auszuwerfen, bevor das Image eingebunden werden kann.',
      unmountFailed: 'Aushängen fehlgeschlagen',
      unmountDesc:
        'In einigen Systemen muss das Abbild auf dem Remote-Host manuell ausgeworfen werden, bevor es ausgehängt wird.',
      refresh: 'Bilder aktualisieren'
    },
    script: {
      title: 'Skripte',
      upload: 'Hochladen',
      run: 'Ausführen',
      runBackground: 'Im Hintergrund ausführen',
      runFailed: '',
      attention: 'Achtung',
      delDesc: 'Möchten Sie diese Datei wirklich löschen?',
      confirm: 'Ja',
      cancel: 'Nein',
      delete: 'Löschen',
      close: 'Schliessen'
    },
    terminal: {
      title: 'Terminal',
      nanokvm: 'NanoKVM Terminal',
      serial: 'Serieller Anschluss Terminal',
      serialPort: 'Serieller Anschluss',
      serialPortPlaceholder: 'Bitte seriellen Anschluss angeben',
      baudrate: 'Baudrate',
      parity: 'Parität',
      parityNone: 'Keine',
      parityEven: 'Gerade',
      parityOdd: 'Ungerade',
      flowControl: 'Fluss-Kontrolle',
      flowControlNone: 'Keine',
      flowControlSoft: 'Wenig',
      flowControlHard: 'Viel',
      dataBits: 'Daten bits',
      stopBits: 'Stopp bits',
      confirm: 'Ok'
    },
    wol: {
      title: 'Wake-on-LAN',
      sending: 'Sende Befehl...',
      sent: 'Befehl gesendet',
      input: 'Bitte MAC Adresse eingeben',
      ok: 'Ok'
    },
    download: {
      title: 'Systemabbild Downloader',
      input: 'Bitte geben Sie die URL für das Remote-Systemabbild ein',
      ok: 'Ok',
      disabled:
        '/data Partition ist nur-lesbar, daher kann das Systemabbild nicht heruntergeladen werden'
    },
    power: {
      title: 'Power',
      showConfirm: 'Bestätigung',
      showConfirmTip: 'Diese Aktionen benötigen eine zusätzliche Bestätigung',
      reset: 'Zurücksetzen',
      power: 'Power',
      powerShort: 'Power (Kurzer Klick)',
      powerLong: 'Power (Langer Klick)',
      resetConfirm: 'Reset-Aktion durchführen?',
      powerConfirm: 'Power-Aktion durchführen?',
      okBtn: 'Ja',
      cancelBtn: 'Nein'
    },
    settings: {
      title: 'Einstellungen',
      about: {
        title: 'Über NanoKVM',
        information: 'Informationen',
        ip: 'IP',
        mdns: 'mDNS',
        application: 'Applikations-Version',
        applicationTip: 'NanoKVM Web Applikations-Version',
        image: 'Systemabbild-Version',
        imageTip: 'NanoKVM Systemabbild-Version',
        deviceKey: 'Geräte Schlüssel',
        community: 'Community',
        hostname: 'Hostname',
        hostnameUpdated: 'Hostname aktualisiert. Neustarten um zu übernehmen.',
        pikvm: {
          title: 'Zu PiKVM wechseln',
          attention: 'Achtung',
          desc1: 'Sind Sie sicher, dass Sie das aktuelle System auf PiKVM umstellen möchten?',
          desc2: 'Das Gerät wird nach dem Wechsel des Systems automatisch neu starten.',
          confirm: 'Bestätigen'
        },
        ipType: {
          Wired: 'Kabel',
          Wireless: 'Drahtlos',
          Other: 'Andere'
        }
      },
      appearance: {
        title: 'Erscheinungsbild',
        display: 'Bildschirm',
        language: 'Sprache',
        menuBar: 'Menu-Leiste',
        menuBarDesc: 'Symbole in der Menu-Leiste anzeigen',
        webTitle: 'Web Titel',
        webTitleDesc: 'Passen Sie den Web-Seite Titel an'
      },
      device: {
        title: 'Gerät',
        oled: {
          title: 'OLED',
          description: 'Schalte OLED Bildschirm aus nach',
          0: 'Nie',
          15: '15 Sek',
          30: '30 Sek',
          60: '1 Min',
          180: '3 Min',
          300: '5 Min',
          600: '10 Min',
          1800: '30 Min',
          3600: '1 Stunde'
        },
        wifi: {
          title: 'Wi-Fi',
          description: 'Wi-Fi konfigurieren'
        },
        ssh: {
          description: 'Aktiviere entfernten SSH-Zugang',
          tip: 'Setzten Sie ein starkes Passwort vor dem aktivieren (Konto - Passwort ändern)'
        },
        advanced: 'Erweiterte Einstellungen',
        mouseJiggler: {
          title: 'Mausaktivitäts-Simulator',
          description: 'Verhindert, dass der remote Host in den Energiesparmodus wechselt',
          disable: 'Deaktivieren',
          absolute: 'Absoluter Modus',
          relative: 'Relativer Modus'
        },
        mdns: {
          description: 'mDNS-Erkennungsdienst aktivieren'
        },
        hidOnly: 'HID-Only Mode',
        disk: 'Virtuelle Festplatte',
        diskDesc: 'Binde das virtuelle U-Laufwerk an den entfernten Host',
        network: 'Virtuelles Netzwerk',
        networkDesc: 'Binde die virtuelle Netzwerkkarte an den entfernten Host',
        reboot: 'Neustarten',
        rebootDesc: 'Sind Sie sicher dass Sie NanoKVM neustarten möchten?',
        okBtn: 'Ja',
        cancelBtn: 'Nein'
      },
      tailscale: {
        title: 'Tailscale',
        restart: 'Tailscale neu starten?',
        stop: 'Tailscale stoppen?',
        stopDesc: 'Von Tailscale abmelden und automatischen Start beim Booten deaktivieren.',
        loading: 'Lädt...',
        notInstall: 'Tailscale nicht gefunden! Bitte installieren.',
        install: 'Installieren',
        installing: 'Installiere',
        failed: 'Installation fehlgeschlagen',
        retry: 'Bitte Seite neu laden und erneut versuchen oder manuelle Installation versuchen.',
        download: 'Laden Sie das',
        package: 'Installations-Paket herunter',
        unzip: 'und entpacken Sie es',
        upTailscale: 'Tailscale nach /usr/bin/ auf NanoKVM hochladen',
        upTailscaled: 'Tailscaled nach /usr/bin/ auf NanoKVM hochladen',
        refresh: 'Aktuelle Seite neu laden',
        notLogin:
          'Das Gerät konnte noch nicht gefunden werden. Bitte melden Sie sich an und verknüpfen Sie dieses Gerät mit Ihrem Konto.',
        urlPeriod: 'Diese URL ist für 10 Minuten gültig',
        login: 'Anmelden',
        loginSuccess: 'Anmeldung erfolgreich',
        enable: 'Tailscale einschalten',
        deviceName: 'Geräte Name',
        deviceIP: 'Geräte IP',
        account: 'Konto',
        logout: 'Abmelden',
        logoutDesc: 'Möchten Sie sich wirklich abmelden?',
        uninstall: 'Tailscale deinstallieren',
        okBtn: 'Ja',
        cancelBtn: 'Nein'
      },
      update: {
        title: 'Nach Aktualisierungen suchen',
        queryFailed: 'Version konnte nicht abgefragt werden',
        updateFailed: 'Aktualisierung fehlgeschlagen. Bitte versuchen Sie es erneut.',
        isLatest: 'Sie haben bereits die aktuellste Version.',
        available: 'Eine Aktualisierung ist verfügbar. Möchten sie diese jetzt durchführen?',
        updating: 'Aktualisierung gestartet. Bitte warten...',
        confirm: 'Bestätigen',
        cancel: 'Abbrechen',
        preview: 'Vorab-Versionen',
        previewDesc: 'Erhalten Sie vorab Zugriff auf neue Funktionen und Verbesserungen',
        previewTip:
          'Bitte beachten Sie, dass Vorab-Versionen womöglich noch Fehler oder unvollständige Funktionen enthalten!'
      },
      account: {
        title: 'Konto',
        webAccount: 'Web Konto Name',
        password: 'Passwort',
        updateBtn: 'Ändern',
        logoutBtn: 'Abmelden',
        logoutDesc: 'Möchten Sie sich wirklich abmelden?',
        okBtn: 'Ja',
        cancelBtn: 'Nein'
      }
    },
    error: {
      title: 'Wir sind auf ein Problem gestossen',
      refresh: 'Neuladen'
    },
    fullscreen: {
      toggle: 'Vollbild ein/aus'
    },
    menu: {
      collapse: 'Menu einblenden',
      expand: 'Menu verbergen'
    }
  }
};

export default de;
