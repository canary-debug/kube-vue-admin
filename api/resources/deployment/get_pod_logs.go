package deployment

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
)

// LogHandler 处理日志请求
func GetPodLogs(c *gin.Context) {
	namespace := c.Param("namespace")
	podName := c.Param("pod")
	container := c.Query("container") // 选填，多容器 Pod 需指定
	// 参数校验
	if namespace == "" || podName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace and pod are required"})
		return
	}

	// 参数：是否实时流
	follow, _ := strconv.ParseBool(c.DefaultQuery("follow", "false"))
	// 参数：展示最后多少行
	tailLinesStr := c.Query("tail")
	var tailLines *int64
	if tailLinesStr != "" {
		t, _ := strconv.ParseInt(tailLinesStr, 10, 64)
		tailLines = &t
	}

	// 1. 构造 PodLogOptions
	opts := &corev1.PodLogOptions{
		Container: container,
		Follow:    follow,    // 关键：true 表示实时流，false 表示一次性快照
		TailLines: tailLines, // 关键：控制行号
		//Timestamps: true,      // 可选：显示时间戳
	}

	// 发起日志请求
	req := resources.Clientset.CoreV1().Pods(namespace).GetLogs(podName, opts)
	stream, err := req.Stream(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open stream: " + err.Error()})
		return
	}
	defer stream.Close()

	// 根据是否实时，选择输出方式
	if follow {
		// 实时流：使用 SSE (Server-Sent Events) 协议
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		reader := bufio.NewReader(stream)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return
			}
			// 发送到 Gin 响应流
			fmt.Fprintf(c.Writer, "data: %s\n\n", line)
			c.Writer.Flush() // 强制刷新缓冲区，实现“实时”
		}
	} else {
		// 非实时：直接将整个 Stream 拷贝到响应体
		c.Writer.Header().Set("Content-Type", "text/plain")
		_, _ = io.Copy(c.Writer, stream)
	}
}
