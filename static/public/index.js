async function submit() {
    const secretElem = document.getElementById('secretInput');
    const value = secretElem.value.trim();
    if (value === '') {
        return;
    }
    const uses = document.getElementById('uses').value;
    const expiration = document.getElementById('expiration').value;

    const key = crypto.getRandomValues(new Uint8Array(32));
    const encryptedValue = await encrypt(value, key);

    try {
        const id = await uploadValue({encryptedValue, uses, expiration});
        window.location = `/link#${new URLSearchParams({
            id: btoa(id),
            key: btoa(String.fromCharCode(...key))
        }).toString()}`;
    } catch (e) {
        alert('Failed to upload secret: ' + e.message);
    }
}

async function encrypt(text, key) {
    const encoder = new TextEncoder();
    const data = encoder.encode(text);
    const iv = crypto.getRandomValues(new Uint8Array(12));
    const algorithm = {name: 'AES-GCM', iv: iv};
    const cryptoKey = await crypto.subtle.importKey('raw', key, algorithm, false, ['encrypt']);
    const encrypted = await crypto.subtle.encrypt(algorithm, cryptoKey, data);
    return btoa(String.fromCharCode(...iv) + String.fromCharCode(...new Uint8Array(encrypted)));
}

async function uploadValue(data) {
    const response = await fetch('/secrets', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            encrypted_secret: data.encryptedValue,
            uses: Number.parseInt(data.uses, 10),
            expiration: data.expiration,
        })
    });
    if (!response.ok) {
        throw new Error('Failed to upload secret: ' + await response.text());
    }
    return await response.text();
}