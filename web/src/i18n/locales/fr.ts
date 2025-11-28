const fr = {
  translation: {
    head: {
      desktop: 'Bureau à distance',
      login: 'Connexion',
      changePassword: 'Changer le mot de passe',
      terminal: 'Terminal',
      wifi: 'Wi-Fi'
    },
    auth: {
      login: 'Connexion',
      placeholderUsername: "Veuillez entrer votre nom d'utilisateur",
      placeholderPassword: 'Veuillez entrer votre mot de passe',
      placeholderPassword2: 'Veuillez entrer votre mot de passe à nouveau',
      noEmptyUsername: "Le nom d'utilisateur ne peut pas être vide",
      noEmptyPassword: 'Le mot de passe ne peut pas être vide',
      noAccount:
        "Impossible de récupérer les informations de l'utilisateur, veuillez rafraîchir la page ou réinitialiser le mot de passe",
      invalidUser: "Nom d'utilisateur ou mot de passe invalide",
      error: 'Erreur inattendue',
      changePassword: 'Changer le mot de passe',
      changePasswordDesc:
        'Pour la sécurité de votre appareil, veuillez modifier le mot de passe de connexion Web.',
      differentPassword: 'Les mots de passe ne correspondent pas',
      illegalUsername: "Le nom d'utilisateur contient des caractères illégaux",
      illegalPassword: 'Le mot de passe contient des caractères illégaux',
      forgetPassword: 'Mot de passe oublié',
      ok: 'Se connecter',
      cancel: 'Annuler',
      loginButtonText: 'Connexion'
    },
    wifi: {
      title: 'Wi-Fi',
      description: 'Configurez le Wi-Fi pour le NanoKVM',
      success:
        'Veuillez vérifier le statut du réseau du NanoKVM et visitez la nouvelle adresse IP.',
      failed: "L'opération a échoué, veuillez réessayer.",
      confirmBtn: 'Ok',
      finishBtn: 'Terminé'
    },
    screen: {
      video: 'Mode vidéo',
      quality: 'Qualité',
      qualityLossless: 'Sans perte',
      qualityHigh: 'Élevé',
      qualityMedium: 'Moyen',
      qualityLow: 'Bas',
      resetHdmi: 'Réinitialiser le HDMI'
    },
    keyboard: {
      paste: 'Coller',
      tips: 'Seuls les caractères et symboles standard du clavier sont pris en charge',
      placeholder: 'Veuillez saisir',
      submit: 'Soumettre',
      virtual: 'Clavier',
      ctrlaltdel: 'Ctrl+Alt+Del'
    },
    mouse: {
      default: 'Curseur par défaut',
      pointer: 'Curseur de la souris',
      cell: 'Curseur de cellule',
      text: 'Curseur de texte',
      grab: 'Curseur de poignée',
      hide: 'Cacher le curseur',
      mode: 'Mode de la souris',
      absolute: 'Mode absolu',
      relative: 'Mode relatif',
      requestPointer:
        'Pour utiliser le mode relatif, cliquez sur le bureau pour capturer le pointeur de la souris.',
      resetHid: 'Réinitialiser le périphérique HID'
    },
    image: {
      title: 'Images',
      loading: 'Chargement',
      empty: 'Vide',
      mountFailed: "Échec du montage de l'image.",
      mountDesc:
        "Dans certains systèmes, il est nécessaire de déséjecter le disque virtuel sur l'hôte distant avant de monter l'image.",
      unmountFailed: 'Démontage échoué',
      unmountDesc:
        "Dans certains systèmes, vous devez éjecter manuellement sur l'hôte distant avant de démonter l'image."
    },
    script: {
      title: 'Script',
      upload: 'Téléverser',
      run: 'Exécuter',
      runBackground: 'Exécuter en arrière-plan',
      runFailed: "Échec de l'exécution",
      attention: 'Attention',
      delDesc: 'Êtes-vous sûr de vouloir supprimer ce fichier ?',
      confirm: 'Oui',
      cancel: 'Non',
      delete: 'Supprimer',
      close: 'Fermer'
    },
    terminal: {
      title: 'Terminal',
      nanokvm: 'Terminal NanoKVM',
      serial: 'Terminal Port Série',
      serialPort: 'Port série',
      serialPortPlaceholder: 'Veuillez entrer le port série',
      baudrate: 'Débit en bauds',
      confirm: 'Ok'
    },
    wol: {
      title: 'Wake-on-LAN',
      sending: 'Envoi de la commande...',
      sent: 'Commande envoyée',
      input: "Veuillez entrer l'adresse MAC",
      ok: 'Ok'
    },
    download: {
      title: 'Télécharger l’image',
      input: 'Veuillez entrer l’URL d’une image distante',
      ok: 'Ok',
      disabled: 'La partition /data est en lecture seule, impossible de télécharger l’image'
    },
    power: {
      title: 'Power',
      reset: 'Réinitialiser',
      power: 'Power',
      powerShort: 'Power (appui court)',
      powerLong: 'Power (appui long)'
    },
    settings: {
      title: 'Paramètres',
      about: {
        title: 'A propos de NanoKVM',
        information: 'Informations',
        ip: 'IP',
        mdns: 'mDNS',
        application: "Version de l'application",
        applicationTip: "Version de l'application Web NanoKVM",
        image: "Version de l'image",
        imageTip: "Version de l'image système NanoKVM",
        deviceKey: "Clé de l'appareil",
        community: 'Communauté'
      },
      appearance: {
        title: 'Apparence',
        display: 'Affichage',
        language: 'Langue',
        menuBar: 'Barre de menus',
        menuBarDesc: 'Afficher les icônes dans la barre de menus'
      },
      device: {
        title: 'Appareil',
        oled: {
          title: 'OLED',
          description: "Écran OLED s'éteint automatiquement",
          0: 'Jamais',
          15: '15 sec',
          30: '30 sec',
          60: '1 min',
          180: '3 min',
          300: '5 min',
          600: '10 min',
          1800: '30 min',
          3600: '1 heure'
        },
        wifi: {
          title: 'Wi-Fi',
          description: 'Configurez le Wi-Fi'
        },
        disk: 'Disque virtuel',
        diskDesc: "Monter le disque virtuel U sur l'hôte distant",
        network: 'Réseau virtuel',
        networkDesc: "Monter la carte réseau virtuelle sur l'hôte distant"
      },
      tailscale: {
        title: 'Tailscale',
        restart: 'Êtes-vous sûr de vouloir redémarrer Tailscale ?',
        stop: 'Êtes-vous sûr de vouloir arrêter Tailscale ?',
        stopDesc: 'Arrêtez Tailscale et désactivez son démarrage automatique.',
        loading: 'Chargement...',
        notInstall: "Tailscale non trouvé ! Veuillez l'installer.",
        install: 'Installer',
        installing: 'Installation',
        failed: 'Installation échouée',
        retry: "Veuillez rafraîchir et réessayer. Ou essayez d'installer manuellement",
        download: 'Télécharger le',
        package: "paquet d'installation",
        unzip: 'et décompressez-le',
        upTailscale: 'Téléverser tailscale dans le répertoire NanoKVM /usr/sbin/',
        upTailscaled: 'Téléverser tailscaled dans le répertoire NanoKVM /usr/sbin/',
        refresh: 'Rafraîchir la page courante',
        notLogin: "L'appareil n'est pas relié. Connectez-vous et liez cet appareil à votre compte.",
        urlPeriod: "L'URL est valide pendant 10 minutes",
        login: 'Connexion',
        loginSuccess: 'Connexion réussie',
        enable: 'Démarrer Tailscale',
        deviceName: "Nom de l'appareil",
        deviceIP: "IP de l'appareil",
        account: 'Compte',
        logout: 'Déconnexion',
        okBtn: 'Oui',
        cancelBtn: 'Non'
      },
      update: {
        title: 'Vérifier les mises à jour',
        queryFailed: 'Impossible de vérifier les mises à jour. Veuillez réessayer.',
        updateFailed: 'Mise à jour échouée. Veuillez réessayer.',
        isLatest: 'Vous avez déjà la dernière version.',
        available: 'Une mise à jour est disponible. Voulez-vous vraiment mettre à jour?',
        updating: 'Mise à jour en cours. Veuillez patienter...',
        confirm: 'Confirmer',
        cancel: 'Annuler'
      },
      account: {
        title: 'Compte',
        webAccount: 'Nom du compte Web',
        password: 'Mot de passe',
        updateBtn: 'Mettre à jour',
        logoutBtn: 'Déconnexion'
      }
    }
  }
};

export default fr;
