package cache

type Store interface {
	Add(obj any) error
	Update(obj any) error
	Delete(obj any) error
	List() []any
	ListKeys() []string
	Get(obj any) (item any, exists bool, err error)
	GetByKey(key string) (item any, exists bool, err error)
	Replace([]any, string) error
	Resync() error
}

type KeyFunc func(obj any) (string, error)
