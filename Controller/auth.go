package controller

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"net/http"
	"slip/utils"

	"github.com/gin-gonic/gin"
	"slip/config"
)


const (
	expectedString = "slip"
)

func Login(c *gin.Context) {
	// 修改为从查询参数获取 KeyID 和 EncryptedString
	encryptedString := c.Query("encrypted_string")
	clientID := c.Query("client_id")

	if encryptedString == "" || clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少参数"})
		return
	}

	if config.AppConfig.Keys["client_id"] != clientID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的密钥 ID"})
		return
	}

	// 解密字符串
	decrypted, err := decrypt(encryptedString, config.AppConfig.Keys["secret_key"])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "解密失败"})
		return
	}

	// 验证解密后的字符串
	if decrypted != expectedString {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的凭证"})
		return
	}

	// 生成 JWT
	token, err := utils.GenerateToken(expectedString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成令牌"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func decrypt(encryptedString, key string) (string, error) {
	// 解码 Base64 字符串
	combined, err := base64.StdEncoding.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	// 提取 IV 和密文
	if len(combined) < 12 { // GCM 推荐使用 12 字节的 IV
		return "", err
	}
	iv := combined[:12]
	ciphertext := combined[12:]

	// 创建 AES-GCM 加密器
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 解密
	plaintext, err := aesGCM.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
