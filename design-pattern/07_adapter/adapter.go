package adapter

// 适配器模式
// 用于转换一种接口适配另一种接口
type Target interface {
	Request() string
}

type Adaptee interface {
	SpecificRequest() string
}

func NewAdaptee() Adaptee {
	return &adapteeImpl{}
}

// 实现 Adaptee 接口
type adapteeImpl struct{}

func (*adapteeImpl) SpecificRequest() string {
	return "adaptee method"
}

// 实现 Target 接口， 继承 Adaptee
type adapter struct {
	Adaptee
}

func NewAdapter(adaptee Adaptee) Target {
	return &adapter{
		Adaptee: adaptee,
	}
}

func (a *adapter) Request() string {
	return a.SpecificRequest()
}
