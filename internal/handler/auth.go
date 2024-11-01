package controller

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"net/http"
	"slip/internal/pkg/utils"
	"slip/api/defines"

	"github.com/gin-gonic/gin"
	"slip/internal/config"
)

const (
	expectedString = "slip"
)

func Login(c *gin.Context) {
	var auth defines.Auth

	if err := c.ShouldBindQuery(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if config.AppConfig.Keys["client_id"] != auth.ClientID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的客户端 ID"})
		return
	}

	// 解密字符串
	decrypted, err := decrypt(auth.EncryptedString, config.AppConfig.Keys["secret_key"])
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
