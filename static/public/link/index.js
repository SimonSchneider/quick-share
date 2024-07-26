window.onload = () => {
    window.location.hash.substring(1)
    const params = new URLSearchParams(window.location.hash.substring(1));
    const id = atob(params.get('id'));
    const encKey = atob(params.get('key'));
    const link = `${window.location.origin}/secrets/${id}#${new URLSearchParams({key: btoa(encKey)}).toString()}`;
    document.getElementById("link").href = link;
    const qrCodeContainer = document.getElementById('qrcode');
    new QRCode(qrCodeContainer, link);
}

function reset() {
    window.location = '/';
}

function copyLink() {
    const link = document.getElementById("link").href;
    navigator.clipboard.writeText(link);
}

function submit() {
    reset();
}