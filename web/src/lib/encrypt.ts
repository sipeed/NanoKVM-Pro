import CryptoJS from 'crypto-js';

// This key is only used to prevent the data from being transmitted in plaintext.
const SECRET_KEY = 'nanokvmpro-sipeed-2025';

export function encrypt(data: string) {
  const dataEncrypt = CryptoJS.AES.encrypt(data, SECRET_KEY).toString();
  return encodeURIComponent(dataEncrypt);
}
