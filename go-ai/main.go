package main

import (
	"fmt"
	"image/color"
	"log"

	"gocv.io/x/gocv"
)

// https://github.com/opencv/opencv/blob/4.x/data/haarcascades/haarcascade_frontalface_default.xml

func main() {
	// 打开摄像头
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		log.Fatalf("Error opening webcam: %v", err)
	}
	defer webcam.Close()

	// 加载人脸分类器
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load("haarcascade_frontalface_default.xml") {
		log.Fatalf("Error reading cascade file: haarcascade_frontalface_default.xml")
	}

	// 打开窗口以显示视频
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	// 创建一个图像矩阵以保存帧
	img := gocv.NewMat()
	defer img.Close()

	fmt.Printf("Press ESC to stop\n")

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed\n")
			return
		}
		if img.Empty() {
			continue
		}

		// 转换图像为灰度
		gray := gocv.NewMat()
		defer gray.Close()
		gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

		// 探测人脸
		rects := classifier.DetectMultiScale(gray)
		for _, r := range rects {
			// 在原图上画矩形
			gocv.Rectangle(&img, r, color.RGBA{0, 255, 0, 0}, 3)
		}

		// 显示图像
		window.IMShow(img)
		if window.WaitKey(1) == 27 {
			break
		}
	}
}