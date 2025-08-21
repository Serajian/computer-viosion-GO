package main

import (
	"fmt"
	"log/slog"

	"gocv.io/x/gocv"
)

func main() {
	fmt.Println("Hello World")
	camID := 0
	webcam, err := gocv.OpenVideoCapture(camID)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer func() {
		err = webcam.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}()
	if !webcam.IsOpened() {
		slog.Error("Webcam is not opened")
		return
	}
	//window := gocv.NewWindow("My Camera")
	//defer func() {
	//	err = window.Close()
	//	if err != nil {
	//		slog.Error(err.Error())
	//	}
	//}()
	frame := gocv.NewMat()
	defer func() {
		err = frame.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	//for {
	//	webcam.Read(&img)
	//	_ = window.IMShow(img)
	//	if window.WaitKey(1) == 27 {
	//		break
	//	}
	//}

	ok := webcam.Read(&frame)
	if !ok || frame.Empty() {
		slog.Error("Unable to read frame from cam")
		return
	}
	fmt.Printf("OK. got one frame: %dx%d\n", frame.Cols(), frame.Rows())
}
