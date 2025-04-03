package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

const (
	HookTypePostFinish = "post-finish"
	HookTypePreFinish  = "pre-finish"
	HookTypePreCreate  = "pre-create"
)

func main() {
	r := gin.Default()

	//tusd web hook
	r.POST("/api/v1/upload/handle-tusd-hooks", HandleTusdHooks)
	err := r.Run(":8005")
	if err != nil {
		log.Fatalf("impossible to start server: %s", err)
	}
}

func HandleTusdHooks(c *gin.Context) {
	var dto handleHookRequest

	//pre-create hook
	if dto.Type == HookTypePreCreate {
		id := "120342141823123"
		c.JSON(http.StatusOK, gin.H{
			"ChangeFileInfo": map[string]interface{}{
				"ID": id,
				"MetaData": map[string]string{
					"video_id": id,
				},
			},
		})
		return
	}

	//pre-finish hook
	if dto.Type == HookTypePreFinish {
		videoIDStr := dto.Event.UploadEvent.MetaData.VideoID
		videoID, _ := strconv.ParseInt(videoIDStr, 10, 64)
		res := GetSuccessResponse(clientResponse{dto.Event.UploadEvent.MetaData.VideoID})
		r, _ := json.Marshal(res)

		fmt.Println("video id", videoID)
		c.JSON(http.StatusOK, gin.H{
			"HTTPResponse": map[string]interface{}{
				"StatusCode": 200,
				"Body":       string(r),
				"Headers": map[string]interface{}{
					"Content-Type": "application/json",
				},
			},
		})
		return
	}

	if dto.Type == HookTypePostFinish {
		// Get information about the current, incoming HTTP request.
		path := dto.Event.HttpRequestEvent.URI
		videoID := dto.Event.UploadEvent.MetaData.VideoID
		fmt.Println("file path", path)
		fmt.Println("video id received in post finish hook", videoID)
		c.JSON(200, gin.H{})
		return
	}
	c.JSON(200, gin.H{})
	return
}

type clientResponse struct {
	ID string `json:"id"`
}

type handleHookRequest struct {
	Type  string `binding:"required" json:"Type"`
	Event event  `binding:"required" json:"Event"`
}

type event struct {
	HttpRequestEvent httpRequestEvent `json:"HTTPRequest"`
	UploadEvent      uploadEvent      `json:"Upload"`
}
type httpRequestEvent struct {
	URI string `json:"URI"`
}

type uploadEvent struct {
	ID       string   `json:"ID"`
	MetaData metaData `json:"MetaData"`
}

type metaData struct {
	VideoID string `json:"video_id"`
}
type StandardResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors,omitempty"`
}

func GetSuccessResponse(data interface{}) StandardResponse {
	return StandardResponse{
		Message: "success",
		Data:    data,
	}
}
