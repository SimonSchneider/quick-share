async function retrieveSecret() {
    const params = new URLSearchParams(window.location.hash.substring(1));
    const key = atob(params.get('key'));
    const encryptedSecret = atob(document.getElementById('dataContainer').innerText);
    const secret = await decrypt(encryptedSecret, key);
    document.getElementById('secret').innerText = secret;
}

async function decrypt(encrypted, key) {
    const iv = new Uint8Array([...encrypted].slice(0, 12).map(c => c.charCodeAt(0)));
    const data = new Uint8Array([...encrypted].slice(12).map(c => c.charCodeAt(0)));
    const algorithm = {name: 'AES-GCM', iv: iv};
    const cryptoKey = await crypto.subtle.importKey('raw', new Uint8Array([...key].map(c => c.charCodeAt(0))), algorithm, false, ['decrypt']);
    const decrypted = await crypto.subtle.decrypt(algorithm, cryptoKey, data);
    const decoder = new TextDecoder();
    return decoder.decode(decrypted);
}

function copySecret() {
    const secret = document.getElementById('secret').innerText;
    navigator.clipboard.writeText(secret);
}

window.onload = retrieveSecret;
