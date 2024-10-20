package main

import (
	"net/http"

	"github.com/BrondoL/captcha/config"
	"github.com/BrondoL/captcha/constant"
	"github.com/BrondoL/captcha/pkg/cache"
	"github.com/BrondoL/captcha/pkg/captcha"
	"github.com/BrondoL/captcha/util"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config.LoadEnv()
	cfg := config.GetEnv()
	redisCache := cache.NewRedis(cfg)

	api := r.Group("/api/v1")

	// Route untuk generate captcha
	api.GET("/generate-captcha", func(c *gin.Context) {
		// Generate ID dan text captcha
		captcha := captcha.NewCaptcha()
		captchaID, captchaText := captcha.GenerateID()

		// Simpan text captcha ke cache
		err := redisCache.Set(c, constant.CacheKeyCaptcha+"-"+captchaID, captchaText, constant.CacheTTLOneMinute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Send response gambar dan captchaID ke client sebagai json
		c.JSON(http.StatusOK, gin.H{
			"id":  captchaID,
			"url": "/captcha/" + captchaID + ".png",
		})
	})

	// Route untuk menampilkan gambar captcha
	api.GET("/captcha/:id", func(c *gin.Context) {
		// Ambil ID dari URL
		captchaID := c.Param("id")
		// split by .png
		captchaID = captchaID[:len(captchaID)-4]

		// Ambil text captcha dari cache
		var captchaText string
		err := redisCache.Get(c, constant.CacheKeyCaptcha+"-"+captchaID, &captchaText)
		if err != nil {
			err := &util.CaptchaNotFound{}
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if captchaText == "" {
			err := &util.CaptchaNotFound{}
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Generate gambar captcha
		captcha := captcha.NewCaptcha()
		img, err := captcha.GenerateImage(captchaText)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Set header response
		c.Header("Content-Type", "image/png")
		c.Header("Content-Length", string(len(img)))

		// Send response gambar captcha ke client
		c.Writer.Write(img)
	})

	// Route untuk verifikasi captcha
	api.POST("/verify-captcha", func(c *gin.Context) {
		// Ambil ID dan text captcha dari request
		var req struct {
			ID   string `json:"id" binding:"required"`
			Text string `json:"text" binding:"required"`
		}

		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Ambil text captcha dari cache
		var storedCaptchaText string
		err = redisCache.Get(c, constant.CacheKeyCaptcha+"-"+req.ID, &storedCaptchaText)
		if err != nil {
			err := &util.CaptchaNotFound{}
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		println("storedCaptchaText: ", storedCaptchaText)
		println("req.Text: ", req.Text)

		// Bandingkan text captcha dari request dengan text captcha dari cache
		if req.Text != storedCaptchaText {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid captcha",
			})
			return
		}

		// Send response validasi captcha ke client
		c.JSON(http.StatusOK, gin.H{
			"message": "Captcha is valid",
		})
	})

	// Jalankan server Gin
	r.Run(":8080")
}
