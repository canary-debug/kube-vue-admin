package pod

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodInfo Pod信息结构体
type PodInfo struct {
	Name         string            `json:"name"`
	Status       string            `json:"status"`
	RestartCount int32             `json:"restart_count"`
	Ports        []int32           `json:"ports"`
	NodeName     string            `json:"node_name"`
	PodIP        string            `json:"pod_ip"`
	CreatedAt    metav1.Time       `json:"created_at"`
	Labels       map[string]string `json:"labels"`
	Namespace    string            `json:"namespace,omitempty"`
}

// ConvertPodToPodInfo 将Pod对象转换为PodInfo结构体
func ConvertPodToPodInfo(pod *v1.Pod) *PodInfo {
	// 计算总重启次数
	totalRestartCount := int32(0)
	for _, status := range pod.Status.InitContainerStatuses {
		totalRestartCount += status.RestartCount
	}
	for _, status := range pod.Status.ContainerStatuses {
		totalRestartCount += status.RestartCount
	}

	// 获取容器端口
	var ports []int32
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			ports = append(ports, port.ContainerPort)
		}
	}

	return &PodInfo{
		Name:         pod.Name,
		Status:       GetPodStatus(pod),
		RestartCount: totalRestartCount,
		Ports:        ports,
		NodeName:     pod.Spec.NodeName,
		PodIP:        pod.Status.PodIP,
		CreatedAt:    pod.CreationTimestamp,
		Labels:       pod.Labels,
		Namespace:    pod.Namespace,
	}
}

// GetPodStatus 获取Pod的简化状态
func GetPodStatus(pod *v1.Pod) string {
	switch pod.Status.Phase {
	case v1.PodPending:
		return "Pending"
	case v1.PodRunning:
		// 检查是否有容器处于等待或终止状态
		for _, status := range pod.Status.ContainerStatuses {
			if status.State.Waiting != nil {
				return "Waiting"
			}
			if status.State.Terminated != nil {
				return "Terminated"
			}
		}
		return "Running"
	case v1.PodSucceeded:
		return "Succeeded"
	case v1.PodFailed:
		return "Failed"
	case v1.PodUnknown:
		return "Unknown"
	default:
		return string(pod.Status.Phase)
	}
}
