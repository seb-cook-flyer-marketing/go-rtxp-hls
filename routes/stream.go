package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/seb-cook-flyer-marketing/go-rtxp-hls/config"
	"github.com/seb-cook-flyer-marketing/go-rtxp-hls/lib/ffmpeg"
	"github.com/seb-cook-flyer-marketing/go-rtxp-hls/lib/types"
)

func RegisterStreamRoutes(router *gin.Engine) {
	router.POST("/stream/convert", AuthenticationMiddleware(handleStreamConvert))
	router.POST("/stream/stop", AuthenticationMiddleware(handleStreamStop))
}

func handleStreamConvert(c *gin.Context) {
	var data types.StreamConvertRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	sourceURL := data.URL
	if strings.Contains(data.URL, "?") {
		sourceURL = strings.Split(data.URL, "?")[0]
	}

	streamID := types.GetStreamID(sourceURL)
	outputPath, err := ffmpeg.ConvertStream(data.URL, streamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": config.Config.URL + "/static/" + outputPath})
}

func handleStreamStop(c *gin.Context) {
	var data types.StreamConvertRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	sourceURL := data.URL
	if strings.Contains(data.URL, "?") {
		sourceURL = strings.Split(data.URL, "?")[0]
	}

	streamID := types.GetStreamID(sourceURL)
	result, err := ffmpeg.StopStream(streamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": result})
}
