package exec

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/canary-debug/kube-vue-admin/api/resources"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// 1. 定义 WebSocket 消息格式，用于处理窗口缩放
type TerminalMessage struct {
	Operation string `json:"operation"` // "stdin" 或 "resize"
	Data      string `json:"data"`      // 输入内容
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

// 2. 实现 remotecommand.TerminalSizeQueue 接口
type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
}

func (t *TerminalSession) Read(p []byte) (int, error) {
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		return 0, err
	}

	var msg TerminalMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		return copy(p, message), nil
	}

	switch msg.Operation {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	default:
		return 0, nil
	}
}

func (t *TerminalSession) Write(p []byte) (int, error) {
	// 添加控制序列过滤，移除可能导致显示问题的序列
	output := string(p)

	// 移除常见的控制序列干扰
	output = strings.ReplaceAll(output, "[?2004l", "") // 禁用bracketed paste mode
	output = strings.ReplaceAll(output, "[?2004h", "") // 启用bracketed paste mode

	// 如果输出以这些序列开头，可能是初始化序列，可以过滤掉
	output = strings.TrimPrefix(output, "[?2004l")

	err := t.wsConn.WriteMessage(websocket.TextMessage, []byte(output))
	return len(output), err
}

func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	size := <-t.sizeChan
	return &size
}

// 3. Gin 接口编写
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Terminal(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer ws.Close()

	session := &TerminalSession{
		wsConn:   ws,
		sizeChan: make(chan remotecommand.TerminalSize),
	}

	// 获取参数
	namespace := c.DefaultQuery("namespace", "default")
	podName := c.Query("pod")
	containerName := c.Query("container")

	// 构造 Exec 请求
	req := resources.Clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   []string{"env", "TERM=xterm-256color", "/bin/sh", "-c", "exec ${SHELL:-/bin/sh}"}, // 修改命令
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(resources.RestConfig, "POST", req.URL())
	if err != nil {
		log.Printf("Create executor error: %v", err)
		return
	}

	// 开始流式传输
	err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             session,
		Stdout:            session,
		Stderr:            session,
		TerminalSizeQueue: session,
		Tty:               true,
	})
	if err != nil {
		log.Printf("Stream error: %v", err)
	}
}
