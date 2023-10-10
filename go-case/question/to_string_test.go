package question

type ConfigOne struct {
	Daemon string
}

func (c *ConfigOne) String() string {
	// 占位符 %v 会调用 xx.String 方法，导致递归调用
	// return fmt.Sprintf("%v", c)
	return ""
}