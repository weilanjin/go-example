package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events

func main() {
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*") // 允许跨域访问

		memT := time.NewTicker(time.Second)
		defer memT.Stop()

		cpuT := time.NewTicker(time.Second)
		defer cpuT.Stop()

		cltDone := r.Context().Done() // 确保在请求结束时停止定时器
		rc := http.NewResponseController(w)
		for {
			select {
			case <-cltDone:
				// 客户端断开连接，停止处理
				log.Println("Client disconnected")
				return
			case <-memT.C:
				// 模拟内存使用情况
				m, err := mem.VirtualMemory()
				if err != nil {
					log.Printf("Error getting memory info: %s", err.Error())
				}
				if _, err := fmt.Fprintf(w, "event: mem\ndata: total: %s, used: %s, perc:%.2f%%\n\n",
					formatBytes(m.Total), formatBytes(m.Used), m.UsedPercent); err != nil {
					log.Printf("Error writing memory data: %s", err.Error())
				}
				rc.Flush() // 确保数据被发送到客户端
			case <-cpuT.C:
				// 模拟CPU使用情况
				c, err := cpu.Times(false)
				if err != nil {
					log.Printf("Error getting memory info: %s", err.Error())
				}
				if _, err := fmt.Fprintf(w, "event: cpu\ndata: User: %s Sys: %s, Idle:%s\n\n",
					formatCPUSeconds(c[0].User), formatCPUSeconds(c[0].System), formatCPUSeconds(c[0].Idle)); err != nil {
					log.Printf("Error writing memory data: %s", err.Error())
				}
				rc.Flush() // 确保数据被发送到客户端
			}
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}

func formatCPUSeconds(seconds float64) string {
	if seconds < 60 {
		return fmt.Sprintf("%.2f s", seconds)
	} else if seconds < 3600 {
		return fmt.Sprintf("%.2f m", seconds/60)
	} else {
		return fmt.Sprintf("%.2f h", seconds/3600)
	}
}

func formatBytes(size uint64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	index := 0
	for size >= 1024 && index < len(units)-1 {
		size /= 1024
		index++
	}
	return fmt.Sprintf("%d %s", size, units[index])
}
