package prototype

// 原型模式 使对象能复制自身，并且暴露到接口中，使客户端面向接口编程时，不需要知道接口实际对的情况下生成新的对象。
type Cloneable interface {
	Clone() Cloneable
}

type PrototypeManager struct {
	prototypes map[string]Cloneable
}

func NewPrototypeManager() *PrototypeManager {
	return &PrototypeManager{
		prototypes: make(map[string]Cloneable),
	}
}

func (p *PrototypeManager) Get(name string) (Cloneable, bool) {
	c, ok := p.prototypes[name]
	if !ok {
		return nil, ok
	}
	return c.Clone(), true
}

func (p *PrototypeManager) Set(name string, prototype Cloneable) {
	p.prototypes[name] = prototype
}
