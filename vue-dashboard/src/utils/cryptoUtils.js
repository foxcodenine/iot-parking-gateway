import CryptoJS from "crypto-js";

function decryptEnv(encryptedText) {
    const secretKey = import.meta.env.VITE_APP_SECRET_KEY;
    const decodedData = CryptoJS.enc.Base64.parse(encryptedText).toString(CryptoJS.enc.Hex);
    const iv = CryptoJS.enc.Hex.parse(decodedData.slice(0, 32));
    const ciphertext = CryptoJS.enc.Hex.parse(decodedData.slice(32));
    const key = CryptoJS.enc.Utf8.parse(secretKey);
    const decrypted = CryptoJS.AES.decrypt({ ciphertext: ciphertext }, key, {
        iv: iv,
        mode: CryptoJS.mode.CFB,
        padding: CryptoJS.pad.NoPadding,
    });

    return CryptoJS.enc.Utf8.stringify(decrypted);
}

export {
    decryptEnv
};
