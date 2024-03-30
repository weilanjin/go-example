package flyweight

// 享元模式 从对象剥离出不发生改变且多个实例需要的重复数据，
// 独立出一个享元，使多个对象共享，从而节省内存以及减少数量。
import (
	"fmt"
	"log"
)

type ImageFlyweightFactory struct {
	maps map[string]*ImageFlyweight
}

var imageFactory *ImageFlyweightFactory

func GetImageFlyweightFactory() *ImageFlyweightFactory {
	if imageFactory == nil {
		imageFactory = &ImageFlyweightFactory{
			maps: make(map[string]*ImageFlyweight),
		}
	}
	return imageFactory
}

func (f *ImageFlyweightFactory) Get(filename string) *ImageFlyweight {
	img := f.maps[filename]
	if img == nil {
		img = NewImageFlyweight(filename)
		f.maps[filename] = img
	}
	return img
}

type ImageFlyweight struct {
	data string
}

func NewImageFlyweight(filename string) *ImageFlyweight {
	s := fmt.Sprintf("image data %s", filename)
	return &ImageFlyweight{
		data: s,
	}
}

func (i *ImageFlyweight) Data() string {
	return i.data
}

type ImageViewer struct {
	*ImageFlyweight
}

func NewImageViewer(filename string) *ImageViewer {
	img := GetImageFlyweightFactory().Get(filename)
	return &ImageViewer{
		ImageFlyweight: img,
	}
}

func (i *ImageViewer) Display() {
	log.Printf("Display: %s", i.Data())
}
