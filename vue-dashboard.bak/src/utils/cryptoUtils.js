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

async function decryptAES(encryptedText) {
    const key = import.meta.env.VITE_AES_SECRET_KEY;
    const decoder = new TextDecoder();
    const encoder = new TextEncoder();

    // Convert the Base64 encrypted text to a Uint8Array
    const encryptedBytes = Uint8Array.from(atob(encryptedText), c => c.charCodeAt(0));

    // Extract the nonce (IV) and the ciphertext
    const ivSize = 12; // 12-byte nonce for AES-GCM
    const iv = encryptedBytes.slice(0, ivSize);
    const ciphertext = encryptedBytes.slice(ivSize);

    // Encode the key into a CryptoKey
    const cryptoKey = await crypto.subtle.importKey(
        "raw",
        encoder.encode(key),
        { name: "AES-GCM" },
        false,
        ["decrypt"]
    );

    // Decrypt the ciphertext
    const decryptedBytes = await crypto.subtle.decrypt(
        { name: "AES-GCM", iv: iv },
        cryptoKey,
        ciphertext
    );

    // Convert decrypted bytes back to a UTF-8 string
    return decoder.decode(decryptedBytes);
}


export {
    encryptString,
    decryptString,
    decryptAES,
};
