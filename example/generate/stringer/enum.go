package enum

// go install golang.org/x/tools/cmd/stringer
// 执行命令
// stringer -type=HeroType -trimprefix=HeroType -linecomment
// 	-trimprefix：删除名称的前缀
//  -linecomment：设置一个完全不同的枚举值名称，只需要在枚举值后面的加上注释

//go:generate stringer -type=HeroType -trimprefix=HeroType -linecomment
type HeroType int

const (
	HeroTypeTank     HeroType = iota + 1
	HeroTypeAssassin          // aa
	HeroTypeMage
)
