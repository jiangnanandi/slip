// 客户端使用的JS加密函数
async function encrypt(plainText, key) {
    // 将密钥转换为 Uint8Array
    const keyBytes = new TextEncoder().encode(key);
    
    // 生成随机的 IV
    const iv = window.crypto.getRandomValues(new Uint8Array(12)); // GCM 推荐使用 12 字节的 IV
    
    // 导入密钥
    const cryptoKey = await window.crypto.subtle.importKey(
        "raw",
        keyBytes,
        { name: "AES-GCM", length: 128 },
        false,
        ["encrypt"]
    );

    // 将明文转换为 Uint8Array
    const plainTextBytes = new TextEncoder().encode(plainText);

    // 使用 AES GCM 模式加密
    const encrypted = await window.crypto.subtle.encrypt(
        {
            name: "AES-GCM",
            iv: iv,
            additionalData: new Uint8Array(), // 可选的附加数据
            tagLength: 128 // 认证标签长度
        },
        cryptoKey,
        plainTextBytes
    );

    // 将 IV 和密文合并，并进行 Base64 编码
    const combined = new Uint8Array(iv.length + encrypted.byteLength);
    combined.set(iv);
    combined.set(new Uint8Array(encrypted), iv.length);

    // 转换为 Base64 字符串
    return btoa(String.fromCharCode(...combined));
}

// 示例用法 16字节的字符串
const key = "your_secret_key"; // 这个key其实是注册在服务器上的，不要泄露
const plainText = "slip";

encrypt(plainText, key).then(encryptedString => {
    console.log(encryptedString);
}).catch(err => {
    console.error("Encryption error:", err);
});