package classifiers

import (
	"log/slog"

	"gocv.io/x/gocv"
)

func Face() gocv.CascadeClassifier {
	classifier := gocv.NewCascadeClassifier()
	if ok := classifier.Load("./assets/haarcascade_frontalface_alt2.xml"); !ok {
		slog.Error("failed to load Haar cascade")
	}
	return classifier
}

func Eye() gocv.CascadeClassifier {
	eyeClf := gocv.NewCascadeClassifier()
	if ok := eyeClf.Load("./assets/haarcascade_eye_tree_eyeglasses.xml"); !ok {
		slog.Error("failed to load Eye cascade")
	}
	return eyeClf
}
