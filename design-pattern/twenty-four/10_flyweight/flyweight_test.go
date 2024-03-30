package flyweight

import "testing"

func TestFlyweight(t *testing.T) {
	iv := NewImageViewer("image1.png")
	iv.Display()
	iv2 := NewImageViewer("image1.png")
	if iv.ImageFlyweight != iv2.ImageFlyweight {
		t.Fatal()
	}
}
