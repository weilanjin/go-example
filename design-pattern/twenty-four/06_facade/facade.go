package facade

import "fmt"

type API interface {
	Test() string
}

type apiImpl struct {
	a AModuleAPI
	b BModuleAPI
}

func NewAPI() API {
	return &apiImpl{
		a: NewAModuleAPI(),
		b: NewBModuleAPI(),
	}
}

func (a apiImpl) Test() string {
	s := a.a.TestA()
	s2 := a.b.TestB()
	return fmt.Sprintf("%s\n%s", s, s2)
}

// A module

type AModuleAPI interface {
	TestA() string
}

type aModuleAPI struct{}

func NewAModuleAPI() AModuleAPI {
	return &aModuleAPI{}
}

func (*aModuleAPI) TestA() string {
	return "A module running"
}

// B module

type BModuleAPI interface {
	TestB() string
}

type bModuleAPI struct{}

func NewBModuleAPI() BModuleAPI {
	return &bModuleAPI{}
}

func (*bModuleAPI) TestB() string {
	return "B module running"
}
