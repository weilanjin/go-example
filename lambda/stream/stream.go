package stream

type Predicate func(any) bool

type Function func(any) any

type Consumer func(any)

type BinaryOperator func(any, any) any

type Supplier func() any

type BiConsumer func(any, any)

type Stream struct {
	values []any
}

func Of(values ...any) Stream {
	return Stream{values: values}
}

func (s Stream) Filter(fn Predicate) Stream {
	var filtered []any
	for _, v := range s.values {
		if fn(v) {
			filtered = append(filtered, v)
		}
	}
	return Stream{values: filtered}
}

func (s Stream) Map(fn Function) Stream {
	var mapped []any
	for _, v := range s.values {
		mapped = append(mapped, fn(v))
	}
	return Stream{values: mapped}
}

func (s Stream) Reduce(fn BinaryOperator, initial any) any {
	var result any = initial
	for _, v := range s.values {
		result = fn(result, v)
	}
	return result
}

func (s Stream) ForEach(fn Consumer) {
	for _, v := range s.values {
		fn(v)
	}
}

func (s Stream) Collect(supplier Supplier, consumer BiConsumer) any {
	var result any = supplier()
	for _, v := range s.values {
		consumer(result, v)
	}
	return result
}
