package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"time"

	"github.com/Serajian/computer-viosion-GO.git/classifiers"
	"gocv.io/x/gocv"
)

const (
	width  = 740
	height = 320
)

func main() {
	fmt.Println("START CV")
	// 1) run webcam
	camID := 0
	webcam, err := gocv.OpenVideoCaptureWithAPI(camID, gocv.VideoCaptureAVFoundation)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	webcam.Set(gocv.VideoCaptureFrameWidth, width)
	webcam.Set(gocv.VideoCaptureFrameHeight, height)
	webcam.Set(gocv.VideoCaptureFPS, 30)
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

	// 2) show picture
	window := gocv.NewWindow("Computer Vision By GO")
	defer func() {
		_ = window.Close()
	}()
	_ = window.ResizeWindow(width, height)

	// Face cascade for first validation
	face := classifiers.Face()

	// Eye cascade for secondary validation (reduces false positives)
	eyeClf := classifiers.Eye()
	// Object detector (full body)
	objDet, err := NewObjectDetector()
	if err != nil {
		slog.Error("load object cascade failed", "err", err)
		return
	}
	objDet.ScaleFactor = 1.03
	objDet.MinNeighbors = 4

	// 3) make matrix
	frame := gocv.NewMat()
	gray := gocv.NewMat()
	edges := gocv.NewMat()
	defer func() {
		_ = frame.Close()
		_ = gray.Close()
		_ = edges.Close()
		_ = face.Close()
		_ = eyeClf.Close()
		objDet.Close()
	}()

	// State
	mode := "normal"
	flip := false
	var fps float64
	frames := 0
	lastTick := time.Now()
	//minSize := getMinSize(&frame)

	for {
		ok := webcam.Read(&frame)
		if !ok || frame.Empty() {
			slog.Error("Unable to read frame from cam")
			continue
		}

		minSize := getMinSize(&frame)
		objDet.MinSize = minSize

		// Process per mode
		display := frame // will be used for drawing/preview

		if flip {
			_ = gocv.Flip(frame, &frame, 1)
		}

		switch mode {
		case "gray":
			_ = gocv.CvtColor(display, &gray, gocv.ColorBGRToGray)
			display = gray
		case "canny":
			_ = gocv.CvtColor(display, &gray, gocv.ColorBGRToGray)
			_ = gocv.Canny(gray, &edges, 50, 150)
			display = edges

		}

		// FPS (simple per-second calc)
		frames++
		elapsed := time.Since(lastTick)
		if elapsed >= time.Second {
			fps = float64(frames) / elapsed.Seconds()
			frames = 0
			lastTick = time.Now()
		}

		guideTXT(frame, fps, mode, flip)

		// detect faces
		_ = gocv.CvtColor(frame, &gray, gocv.ColorBGRToGray)
		_ = gocv.GaussianBlur(gray, &gray, image.Pt(3, 3), 0, 0, gocv.BorderDefault)
		faces := face.DetectMultiScaleWithParams(
			gray,
			1.03,
			6,
			0,
			minSize,
			image.Pt(0, 0),
		)

		// draw boxes
		acceptedFace := 0
		for _, f := range faces {
			// quick geometry filter to reduce false positives
			aspect := float64(f.Dx()) / float64(f.Dy())
			if aspect < 0.7 || aspect > 1.4 { // typical front faces ~0.8–1.3
				continue
			}

			// secondary validation: detect eyes in ROI (in gray)
			roi := gray.Region(f)
			eyes := eyeClf.DetectMultiScaleWithParams(
				roi, 1.05, 3, 0,
				image.Pt(maxInt(12, f.Dx()/10), maxInt(12, f.Dy()/10)), // min eye size
				image.Pt(0, 0),
			)
			_ = roi.Close()
			if len(eyes) == 0 {
				continue
			}

			// draw final box
			_ = gocv.Rectangle(&display, f, getColor("green"), 2)

			// label "human" above the box (avoid going off-screen)
			y := f.Min.Y - 8
			if y < 18 {
				y = f.Min.Y + 18
			}
			putShadowText(&display, "human", image.Pt(f.Min.X, y), 1.4)

			acceptedFace++
		}

		// --- Object detection on the same gray frame ---
		objects := objDet.Detect(gray)
		for _, r := range objects {
			// گزینه: فیلتر خیلی ریزها (اگر دوست داشتی)
			if r.Dx() < 24 || r.Dy() < 24 {
				continue
			}
			_ = gocv.Rectangle(&display, r, getColor("red"), 2)

			yy := r.Min.Y - 8
			if yy < 18 {
				yy = r.Min.Y + 18
			}
			putShadowText(&display, objDet.Label, image.Pt(r.Min.X, yy), 1.2)
		}

		// draw count
		_ = gocv.PutText(&display,
			fmt.Sprintf("Objects: %d", len(objects)),
			image.Pt(10, 100),
			gocv.FontHersheyPlain, 1.5, getColor("red"), 2,
		)

		_ = gocv.PutText(&display,
			fmt.Sprintf("Faces: %d", acceptedFace),
			image.Pt(10, 75),
			gocv.FontHersheyPlain, 1.5, getColor("green"), 2,
		)

		_ = window.IMShow(display)
		key := window.WaitKey(1)

		switch key {
		case 27: // ESC
			slog.Info("Got ESC")
			return
		case 's', 'S':
			name := fmt.Sprintf("snapshot_%d.jpg", time.Now().Unix())
			if ok = gocv.IMWrite(name, display); ok {
				slog.Info("Saved " + name)
			} else {
				slog.Error("Failed to save snapshot")
			}
		case 'f', 'F':
			flip = !flip
		case 'g', 'G':
			mode = "gray"
		case 'e', 'E':
			mode = "canny"
		case 'n', 'N':
			mode = "normal"
		}
	}

}

func guideTXT(display gocv.Mat, fps float64, mode string, flip bool) {
	putShadowText(&display, fmt.Sprintf("FPS: %.1f | Mode: %s | Flip: %v", fps, mode, flip), image.Pt(10, 20), 1.5)
	putShadowText(&display, "S: save  F: flip  G: gray  E: canny  N: normal  ESC: quit", image.Pt(10, 45), 1.5)
}

// putShadowText HUD text (with thin shadow for readability)
func putShadowText(img *gocv.Mat, text string, pt image.Point, scale float64) {
	_ = gocv.PutText(img, text, image.Pt(pt.X+1, pt.Y+1), gocv.FontHersheyPlain, scale, getColor("black"), 2)
	_ = gocv.PutText(img, text, pt, gocv.FontHersheyPlain, scale, getColor("blue"), 2)
}

// getColor UI colors
func getColor(c string) color.RGBA {
	switch c {
	case "black":
		return color.RGBA{0, 0, 0, 0}
	case "white":
		return color.RGBA{255, 255, 255, 0}
	case "red":
		return color.RGBA{255, 0, 0, 255}
	case "blue":
		return color.RGBA{0, 50, 255, 155}
	case "green":
		return color.RGBA{0, 255, 0, 255}
	default:
		return getColor("white")
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ObjectDetector --- Object detector "class" (Cascade-based) ---
type ObjectDetector struct {
	clf          gocv.CascadeClassifier
	Label        string
	ScaleFactor  float64
	MinNeighbors int
	MinSize      image.Point
}

func NewObjectDetector() (*ObjectDetector, error) {
	clf := gocv.NewCascadeClassifier()
	if ok := clf.Load("./assets/haarcascade_upperbody.xml"); !ok {
		return nil, errors.New("failed to load cascade")
	}
	return &ObjectDetector{
		clf:          clf,
		Label:        "object",
		ScaleFactor:  1.05,
		MinNeighbors: 5,
		MinSize:      image.Pt(30, 30),
	}, nil
}

func (od *ObjectDetector) Close() {
	_ = od.clf.Close()
}

func (od *ObjectDetector) Detect(gray gocv.Mat) []image.Rectangle {
	return od.clf.DetectMultiScaleWithParams(
		gray,
		od.ScaleFactor,
		od.MinNeighbors,
		0,
		od.MinSize,
		image.Pt(0, 0),
	)
}

func getMinSize(mat *gocv.Mat) image.Point {
	w, h := mat.Cols(), mat.Rows()
	short := w
	if h < w {
		short = h
	}
	minSize := image.Pt(short/12, short/12)

	return minSize
}
