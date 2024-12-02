import CryptoJS from "crypto-js";

function encryptString(text) {
    const secretKey = import.meta.env.VITE_AES_SECRET_KEY
    return CryptoJS.AES.encrypt(text, secretKey).toString();
}

// Decryption function
function decryptString(ciphertext) {
    const secretKey = import.meta.env.VITE_AES_SECRET_KEY
    const bytes  = CryptoJS.AES.decrypt(ciphertext, secretKey);
    return bytes.toString(CryptoJS.enc.Utf8);
}

function decryptAES(encryptedText) {
    const key = import.meta.env.VITE_AES_SECRET_KEY
    const encryptedBytes = CryptoJS.AES.decrypt(encryptedText, key);
    const decryptedText = encryptedBytes.toString(CryptoJS.enc.Utf8);
    return decryptedText;
}

export {
    encryptString,
    decryptString,
    decryptAES,
};
