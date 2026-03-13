package main

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestStreamChain(t *testing.T) {
	got := Map(
		Distinct(
			From(1, 2, 2, 3, 4, 5).
				Filter(func(v int) bool { return v > 1 }).
				Skip(1).
				Limit(3),
		),
		func(v int) int { return v * 10 },
	).Slice()

	want := []int{20, 30, 40}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected chain result, got=%v want=%v", got, want)
	}
}

func TestMapFlatMap(t *testing.T) {
	mapped := Map(From(1, 2, 3), func(v int) string {
		return string(rune('a' + v - 1))
	}).Slice()

	if !reflect.DeepEqual(mapped, []string{"a", "b", "c"}) {
		t.Fatalf("unexpected mapped result: %v", mapped)
	}

	flat := FlatMap(From(1, 2, 3), func(v int) Stream[int] {
		return From(v, v*10)
	}).Slice()

	if !reflect.DeepEqual(flat, []int{1, 10, 2, 20, 3, 30}) {
		t.Fatalf("unexpected flat result: %v", flat)
	}
}

func TestSortedReduceAndCollectToMap(t *testing.T) {
	sorted := From(3, 1, 2).Sorted(func(a, b int) bool { return a < b }).Slice()
	if !reflect.DeepEqual(sorted, []int{1, 2, 3}) {
		t.Fatalf("unexpected sorted result: %v", sorted)
	}

	sum := From(1, 2, 3, 4).Reduce(0, func(acc, v int) int { return acc + v })
	if sum != 10 {
		t.Fatalf("unexpected reduce result: %d", sum)
	}

	m := CollectToMap(From("aa", "bbb", "c"), func(v string) int { return len(v) }, func(v string) string { return v })
	if len(m) != 3 || m[1] != "c" || m[2] != "aa" || m[3] != "bbb" {
		t.Fatalf("unexpected map result: %v", m)
	}
}

func TestLazyEvaluation(t *testing.T) {
	hit := 0
	s := Map(
		From(1, 2, 3).Filter(func(v int) bool {
			hit++
			return v%2 == 1
		}),
		func(v int) int { return v * 2 },
	)

	if hit != 0 {
		t.Fatalf("lazy broken before terminal, hit=%d", hit)
	}

	_ = s.Slice()
	if hit != 3 {
		t.Fatalf("unexpected traversal count, hit=%d", hit)
	}
}

func TestShortCircuit(t *testing.T) {
	countAny := 0
	any := From(1, 2, 3, 4, 5).Any(func(v int) bool {
		countAny++
		return v == 3
	})
	if !any {
		t.Fatalf("anyMatch should be true")
	}
	if countAny != 3 {
		t.Fatalf("anyMatch should short-circuit at 3, count=%d", countAny)
	}

	countAll := 0
	all := From(2, 4, 5, 6).Every(func(v int) bool {
		countAll++
		return v%2 == 0
	})
	if all {
		t.Fatalf("allMatch should be false")
	}
	if countAll != 3 {
		t.Fatalf("allMatch should short-circuit at first odd, count=%d", countAll)
	}

	first, ok := From(9, 8, 7).Head()
	if !ok || first != 9 {
		t.Fatalf("findFirst unexpected, first=%d ok=%v", first, ok)
	}
}

func TestLimitSkipCount(t *testing.T) {
	got := From(1, 2, 3, 4, 5).Skip(2).Limit(2).Slice()
	if !reflect.DeepEqual(got, []int{3, 4}) {
		t.Fatalf("unexpected skip/limit result: %v", got)
	}

	if c := From(1, 2, 3).Len(); c != 3 {
		t.Fatalf("unexpected count result: %d", c)
	}

	if c := From(1, 2, 3).Filter(func(v int) bool { return v > 1 }).Len(); c != 2 {
		t.Fatalf("unexpected filtered count result: %d", c)
	}

	if c := From(1, 2, 3).Limit(0).Len(); c != 0 {
		t.Fatalf("unexpected limit zero count result: %d", c)
	}
}

func TestDistinctStability(t *testing.T) {
	got := Distinct(From(3, 1, 3, 2, 1, 2)).Slice()
	want := []int{3, 1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected distinct result, got=%v want=%v", got, want)
	}

	res := Distinct(From(5, 4, 4, 3, 2, 2, 1)).Sorted(func(a, b int) bool { return a < b }).Slice()
	if !sort.IntsAreSorted(res) {
		t.Fatalf("sorted distinct result should be ordered: %v", res)
	}
}

func TestFactoriesAndBuilder(t *testing.T) {
	if c := Empty[int]().Len(); c != 0 {
		t.Fatalf("empty count should be 0, got=%d", c)
	}

	if c := FromPointer[int](nil).Len(); c != 0 {
		t.Fatalf("ofNullable nil count should be 0, got=%d", c)
	}

	v := 42
	if got := FromPointer(&v).Slice(); !reflect.DeepEqual(got, []int{42}) {
		t.Fatalf("ofNullable value mismatch: %v", got)
	}

	b := NewBuilder[int]().Add(1).Add(2).Add(3)
	if got := b.Build().Slice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("builder mismatch: %v", got)
	}

	joined := Concat(From(1, 2), From(3, 4)).Slice()
	if !reflect.DeepEqual(joined, []int{1, 2, 3, 4}) {
		t.Fatalf("concat mismatch: %v", joined)
	}
}

func TestTakeDropPeekAndNoneMatch(t *testing.T) {
	peeked := 0
	res := From(1, 2, 3, 4, 5).
		Peek(func(int) { peeked++ }).
		TakeWhile(func(v int) bool { return v < 4 }).
		DropWhile(func(v int) bool { return v < 2 }).
		Slice()

	if !reflect.DeepEqual(res, []int{2, 3}) {
		t.Fatalf("take/drop mismatch: %v", res)
	}
	if peeked != 4 {
		t.Fatalf("peek should observe until takeWhile stop, peeked=%d", peeked)
	}

	if !From(2, 4, 6).None(func(v int) bool { return v%2 == 1 }) {
		t.Fatalf("noneMatch should be true")
	}
}

func TestGenerateIterateAndSortedByKey(t *testing.T) {
	gen := Generate(func() int { return 7 }).Limit(3).Slice()
	if !reflect.DeepEqual(gen, []int{7, 7, 7}) {
		t.Fatalf("generate mismatch: %v", gen)
	}

	it := Iterate(1, func(v int) int { return v + 1 }).Limit(5).Slice()
	if !reflect.DeepEqual(it, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("iterate mismatch: %v", it)
	}

	itw := IterateWhile(1, func(v int) bool { return v <= 3 }, func(v int) int { return v + 1 }).Slice()
	if !reflect.DeepEqual(itw, []int{1, 2, 3}) {
		t.Fatalf("iterateWhile mismatch: %v", itw)
	}

	type item struct {
		name string
		key  int
	}
	sorted := SortedByKey(From(item{"b", 2}, item{"a", 1}), func(v item) int { return v.key }).Slice()
	if sorted[0].name != "a" {
		t.Fatalf("sortedByKey mismatch: %v", sorted)
	}
}

func TestOptionalAndMinMax(t *testing.T) {
	opt := From(3, 1, 2).ReduceOptional(func(a, b int) int { return a + b })
	if v := opt.OrElse(0); v != 6 {
		t.Fatalf("reduceOptional mismatch: %d", v)
	}

	minV := MinOrdered(From(3, 1, 2)).OrElse(0)
	maxV := MaxOrdered(From(3, 1, 2)).OrElse(0)
	if minV != 1 || maxV != 3 {
		t.Fatalf("min/max mismatch min=%d max=%d", minV, maxV)
	}

	collected := Collect(From(1, 2, 3), func() []int { return make([]int, 0) }, func(dst *[]int, v int) {
		*dst = append(*dst, v)
	}, nil)
	if !reflect.DeepEqual(collected, []int{1, 2, 3}) {
		t.Fatalf("collect mismatch: %v", collected)
	}

	reduced := ReduceWithCombiner(From(1, 2, 3), 0, func(acc, v int) int { return acc + v }, nil)
	if reduced != 6 {
		t.Fatalf("reduceWithCombiner mismatch: %d", reduced)
	}
}

func TestStateAndClose(t *testing.T) {
	s := From(1, 2, 3).Parallel()
	if !s.IsParallel() {
		t.Fatalf("parallel state mismatch")
	}
	if s.Sequential().IsParallel() {
		t.Fatalf("sequential should clear parallel state")
	}

	steps := make([]int, 0, 2)
	From(1).
		OnClose(func() { steps = append(steps, 1) }).
		OnClose(func() { steps = append(steps, 2) }).
		Close()
	if !reflect.DeepEqual(steps, []int{2, 1}) {
		t.Fatalf("close order mismatch: %v", steps)
	}

	any, ok := From(8, 9).Head()
	if !ok || any != 8 {
		t.Fatalf("findAny mismatch: any=%d ok=%v", any, ok)
	}
}

func TestParallelMapOrdered(t *testing.T) {
	got := MapPar(From(1, 2, 3, 4), 2, func(v int) int {
		time.Sleep(5 * time.Millisecond)
		return v * 10
	}).Slice()

	want := []int{10, 20, 30, 40}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parallel ordered map mismatch: got=%v want=%v", got, want)
	}
}

func TestParallelMapUnordered(t *testing.T) {
	got := MapParUnordered(From(1, 2, 3, 4, 5), 3, func(v int) int {
		time.Sleep(time.Duration(6-v) * time.Millisecond)
		return v * 100
	}).Slice()

	if len(got) != 5 {
		t.Fatalf("parallel unordered map size mismatch: %v", got)
	}

	counts := make(map[int]int)
	for i := range got {
		counts[got[i]]++
	}
	for _, expect := range []int{100, 200, 300, 400, 500} {
		if counts[expect] != 1 {
			t.Fatalf("parallel unordered map missing value=%d got=%v", expect, got)
		}
	}
}

func TestParallelForEach(t *testing.T) {
	var sum int64
	var mu sync.Mutex
	seen := make(map[int]struct{})

	From(1, 2, 3, 4, 5, 6).EachPar(3, func(v int) {
		atomic.AddInt64(&sum, int64(v))
		mu.Lock()
		seen[v] = struct{}{}
		mu.Unlock()
	})

	if sum != 21 {
		t.Fatalf("parallel forEach sum mismatch: %d", sum)
	}
	if len(seen) != 6 {
		t.Fatalf("parallel forEach seen mismatch: %v", seen)
	}
}

func TestCollectorsFamily(t *testing.T) {
	count := CollectWith(From(1, 2, 3, 4), CountingCollector[int]())
	if count != 4 {
		t.Fatalf("counting collector mismatch: %d", count)
	}

	joined := CollectWith(From("a", "b", "c"), JoiningCollector("-", "[", "]"))
	if joined != "[a-b-c]" {
		t.Fatalf("joining collector mismatch: %s", joined)
	}

	set := CollectWith(From(1, 2, 2, 3), ToSetCollector[int]())
	if len(set) != 3 {
		t.Fatalf("set collector mismatch: %v", set)
	}

	grouped := CollectWith(From("go", "java", "c"), GroupingByCollector(func(v string) int { return len(v) }))
	if !reflect.DeepEqual(grouped[2], []string{"go"}) || !reflect.DeepEqual(grouped[1], []string{"c"}) {
		t.Fatalf("groupingBy collector mismatch: %v", grouped)
	}

	parts := CollectWith(From(1, 2, 3, 4), PartitioningByCollector(func(v int) bool { return v%2 == 0 }))
	if !reflect.DeepEqual(parts[true], []int{2, 4}) || !reflect.DeepEqual(parts[false], []int{1, 3}) {
		t.Fatalf("partitioningBy collector mismatch: %v", parts)
	}
}

func TestCollectorsComposed(t *testing.T) {
	mapped := CollectWith(
		From(1, 2, 3),
		MappingCollector(func(v int) string { return string(rune('a' + v - 1)) }, ToSliceCollector[string]()),
	)
	if !reflect.DeepEqual(mapped, []string{"a", "b", "c"}) {
		t.Fatalf("mapping collector mismatch: %v", mapped)
	}

	filtered := CollectWith(
		From(1, 2, 3, 4),
		FilteringCollector(func(v int) bool { return v%2 == 0 }, ToSliceCollector[int]()),
	)
	if !reflect.DeepEqual(filtered, []int{2, 4}) {
		t.Fatalf("filtering collector mismatch: %v", filtered)
	}

	flatMapped := CollectWith(
		From("a,b", "c"),
		FlatMappingCollector(func(v string) []string {
			if v == "a,b" {
				return []string{"a", "b"}
			}
			return []string{v}
		}, ToSliceCollector[string]()),
	)
	if !reflect.DeepEqual(flatMapped, []string{"a", "b", "c"}) {
		t.Fatalf("flatMapping collector mismatch: %v", flatMapped)
	}

	toMap := CollectWith(
		From("a", "aa", "b"),
		ToMapCollector(
			func(v string) byte { return v[0] },
			func(v string) int { return len(v) },
			func(existing, incoming int) int { return existing + incoming },
		),
	)
	if toMap['a'] != 3 || toMap['b'] != 1 {
		t.Fatalf("toMap collector mismatch: %v", toMap)
	}

	groupMap := CollectWith(
		From("go", "java", "js"),
		GroupingByMappingCollector(func(v string) int { return len(v) }, func(v string) byte { return v[0] }),
	)
	if !reflect.DeepEqual(groupMap[2], []byte{'g', 'j'}) {
		t.Fatalf("groupingByMapping collector mismatch: %v", groupMap)
	}
}

func TestCollectorsStatistics(t *testing.T) {
	sumInt := CollectWith(From(1, 2, 3), SummingIntCollector(func(v int) int { return v }))
	if sumInt != 6 {
		t.Fatalf("summingInt mismatch: %d", sumInt)
	}

	sumI64 := CollectWith(From(int64(1), int64(2), int64(3)), SummingInt64Collector(func(v int64) int64 { return v }))
	if sumI64 != 6 {
		t.Fatalf("summingInt64 mismatch: %d", sumI64)
	}

	sumF64 := CollectWith(From(1.5, 2.5), SummingFloat64Collector(func(v float64) float64 { return v }))
	if sumF64 != 4.0 {
		t.Fatalf("summingFloat64 mismatch: %f", sumF64)
	}

	avgInt := CollectWith(From(1, 2, 3, 4), AveragingIntCollector(func(v int) int { return v }))
	if avgInt != 2.5 {
		t.Fatalf("averagingInt mismatch: %f", avgInt)
	}

	avgI64 := CollectWith(From(int64(2), int64(4)), AveragingInt64Collector(func(v int64) int64 { return v }))
	if avgI64 != 3.0 {
		t.Fatalf("averagingInt64 mismatch: %f", avgI64)
	}

	avgF64 := CollectWith(From(1.0, 2.0, 5.0), AveragingFloat64Collector(func(v float64) float64 { return v }))
	if avgF64 != 8.0/3.0 {
		t.Fatalf("averagingFloat64 mismatch: %f", avgF64)
	}

	intStats := CollectWith(From(3, 1, 2), SummarizingIntCollector(func(v int) int { return v }))
	if intStats.Count != 3 || intStats.Sum != 6 || intStats.Min != 1 || intStats.Max != 3 || intStats.Average() != 2.0 {
		t.Fatalf("summarizingInt mismatch: %+v", intStats)
	}

	longStats := CollectWith(From(int64(7), int64(9)), SummarizingInt64Collector(func(v int64) int64 { return v }))
	if longStats.Count != 2 || longStats.Sum != 16 || longStats.Min != 7 || longStats.Max != 9 || longStats.Average() != 8.0 {
		t.Fatalf("summarizingInt64 mismatch: %+v", longStats)
	}

	doubleStats := CollectWith(From(1.5, 2.5, 0.5), SummarizingFloat64Collector(func(v float64) float64 { return v }))
	if doubleStats.Count != 3 || doubleStats.Sum != 4.5 || doubleStats.Min != 0.5 || doubleStats.Max != 2.5 || doubleStats.Average() != 1.5 {
		t.Fatalf("summarizingFloat64 mismatch: %+v", doubleStats)
	}
}

func TestCollectorsAndThenAndReducing(t *testing.T) {
	collector := CollectingAndThen(ToSliceCollector[int](), func(items []int) int {
		return len(items)
	})
	lenRes := CollectWith(From(1, 2, 3), collector)
	if lenRes != 3 {
		t.Fatalf("collectingAndThen mismatch: %d", lenRes)
	}

	reduced := CollectWith(
		From("a", "bb", "ccc"),
		ReducingCollector(0, func(v string) int { return len(v) }, func(a, b int) int { return a + b }),
	)
	if reduced != 6 {
		t.Fatalf("reducing collector mismatch: %d", reduced)
	}
}

func TestCollectorsMinMaxGroupingDownstreamAndTeeing(t *testing.T) {
	minOpt := CollectWith(From(3, 1, 2), MinByCollector(func(a, b int) bool { return a < b }))
	if v := minOpt.OrElse(0); v != 1 {
		t.Fatalf("minBy collector mismatch: %d", v)
	}

	maxOpt := CollectWith(From(3, 1, 2), MaxByCollector(func(a, b int) bool { return a < b }))
	if v := maxOpt.OrElse(0); v != 3 {
		t.Fatalf("maxBy collector mismatch: %d", v)
	}

	groupedSum := CollectWith(
		From(1, 2, 3, 4, 5),
		GroupingByDownstreamCollector(
			func(v int) bool { return v%2 == 0 },
			SummingIntCollector(func(v int) int { return v }),
		),
	)
	if groupedSum[true] != 6 || groupedSum[false] != 9 {
		t.Fatalf("groupingByDownstream collector mismatch: %v", groupedSum)
	}

	avgLen := CollectWith(
		From("a", "bb", "ccc"),
		TeeingCollector(
			CountingCollector[string](),
			SummingIntCollector(func(v string) int { return len(v) }),
			func(count int64, sum int64) float64 {
				if count == 0 {
					return 0
				}
				return float64(sum) / float64(count)
			},
		),
	)
	if avgLen != 2.0 {
		t.Fatalf("teeing collector mismatch: %f", avgLen)
	}
}

func TestSparkLikeSetAndZipOperators(t *testing.T) {
	if got := Union(From(1, 2), From(3, 4)).Slice(); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
		t.Fatalf("union mismatch: %v", got)
	}

	inter := Intersection(From(1, 2, 3, 3), From(3, 4)).Slice()
	if !reflect.DeepEqual(inter, []int{3, 3}) {
		t.Fatalf("intersection mismatch: %v", inter)
	}

	sub := Subtract(From(1, 2, 3), From(2)).Slice()
	if !reflect.DeepEqual(sub, []int{1, 3}) {
		t.Fatalf("subtract mismatch: %v", sub)
	}

	car := Cartesian(From(1, 2), From("a", "b")).Slice()
	if len(car) != 4 {
		t.Fatalf("cartesian mismatch: %v", car)
	}

	zipped := Zip(From(1, 2, 3), From("a", "b")).Slice()
	if !reflect.DeepEqual(zipped, []Tuple2[int, string]{{1, "a"}, {2, "b"}}) {
		t.Fatalf("zip mismatch: %v", zipped)
	}

	idx := ZipWithIndex(From("x", "y")).Slice()
	if !reflect.DeepEqual(idx, []Tuple2[string, int]{{"x", 0}, {"y", 1}}) {
		t.Fatalf("zipWithIndex mismatch: %v", idx)
	}
}

func TestSparkLikeSortSampleAndTake(t *testing.T) {
	sorted := SortBy(From("bbb", "a", "cc"), func(v string) int { return len(v) }, true).Slice()
	if !reflect.DeepEqual(sorted, []string{"a", "cc", "bbb"}) {
		t.Fatalf("sortBy mismatch: %v", sorted)
	}

	sampled := Sample(From(1, 2, 3, 4, 5), false, 1.0, 1).Slice()
	if !reflect.DeepEqual(sampled, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("sample mismatch: %v", sampled)
	}

	glom := Glom(From(1, 2, 3, 4, 5), 2).Slice()
	if len(glom) != 3 || !reflect.DeepEqual(glom[0], []int{1, 2}) || !reflect.DeepEqual(glom[2], []int{5}) {
		t.Fatalf("glom mismatch: %v", glom)
	}

	take := From(5, 4, 3).Take(2)
	if !reflect.DeepEqual(take, []int{5, 4}) {
		t.Fatalf("take mismatch: %v", take)
	}

	first, ok := From(9, 8).First()
	if !ok || first != 9 {
		t.Fatalf("first mismatch: %d %v", first, ok)
	}

	if !reflect.DeepEqual(TakeOrdered(From(3, 1, 2), 2), []int{1, 2}) {
		t.Fatalf("takeOrdered mismatch")
	}
	if !reflect.DeepEqual(Top(From(3, 1, 2), 2), []int{3, 2}) {
		t.Fatalf("top mismatch")
	}
}

func TestSparkLikePairOperators(t *testing.T) {
	pairs := From(
		Pair[string, int]{Key: "a", Value: 1},
		Pair[string, int]{Key: "a", Value: 2},
		Pair[string, int]{Key: "b", Value: 3},
	)

	mappedVals := MapValues(pairs, func(v int) int { return v * 10 }).Slice()
	if mappedVals[0].Value != 10 || mappedVals[2].Value != 30 {
		t.Fatalf("mapValues mismatch: %v", mappedVals)
	}

	flatVals := FlatMapValues(pairs, func(v int) []int { return []int{v, -v} }).Slice()
	if len(flatVals) != 6 {
		t.Fatalf("flatMapValues mismatch: %v", flatVals)
	}

	grp := GroupByKey(pairs).Slice()
	if len(grp) != 2 {
		t.Fatalf("groupByKey mismatch: %v", grp)
	}

	reduced := ReduceByKey(pairs, func(a, b int) int { return a + b }).Slice()
	if len(reduced) != 2 {
		t.Fatalf("reduceByKey mismatch: %v", reduced)
	}

	folded := FoldByKey(pairs, 10, func(a, b int) int { return a + b }).Slice()
	if len(folded) != 2 {
		t.Fatalf("foldByKey mismatch: %v", folded)
	}

	sortedPair := SortByKey(From(
		Pair[int, string]{Key: 2, Value: "b"},
		Pair[int, string]{Key: 1, Value: "a"},
	), true).Slice()
	if sortedPair[0].Key != 1 {
		t.Fatalf("sortByKey mismatch: %v", sortedPair)
	}

	if cnt := CountByKey(pairs); cnt["a"] != 2 || cnt["b"] != 1 {
		t.Fatalf("countByKey mismatch: %v", cnt)
	}
	if cnt := CountByValue(From(1, 1, 2)); cnt[1] != 2 || cnt[2] != 1 {
		t.Fatalf("countByValue mismatch: %v", cnt)
	}

	ks := Keys(pairs).Slice()
	vs := Values(pairs).Slice()
	if len(ks) != 3 || len(vs) != 3 {
		t.Fatalf("keys/values mismatch: %v %v", ks, vs)
	}
}

func TestSparkLikeJoinFamilyAndCogroup(t *testing.T) {
	left := From(
		Pair[string, int]{Key: "a", Value: 1},
		Pair[string, int]{Key: "b", Value: 2},
	)
	right := From(
		Pair[string, string]{Key: "a", Value: "x"},
		Pair[string, string]{Key: "c", Value: "y"},
	)

	inner := Join(left, right).Slice()
	if len(inner) != 1 || inner[0].Key != "a" || inner[0].Value.First != 1 || inner[0].Value.Second != "x" {
		t.Fatalf("join mismatch: %v", inner)
	}

	loj := LeftOuterJoin(left, right).Slice()
	if len(loj) != 2 {
		t.Fatalf("leftOuterJoin size mismatch: %v", loj)
	}

	roj := RightOuterJoin(left, right).Slice()
	if len(roj) != 2 {
		t.Fatalf("rightOuterJoin size mismatch: %v", roj)
	}

	foj := FullOuterJoin(left, right).Slice()
	if len(foj) != 3 {
		t.Fatalf("fullOuterJoin size mismatch: %v", foj)
	}

	cg := Cogroup(left, right).Slice()
	if len(cg) != 3 {
		t.Fatalf("cogroup mismatch: %v", cg)
	}
}

func TestSparkLikeAggregateCombinePartitionAndLookup(t *testing.T) {
	pairs := From(
		Pair[string, int]{Key: "a", Value: 1},
		Pair[string, int]{Key: "a", Value: 2},
		Pair[string, int]{Key: "b", Value: 3},
	)

	agg := AggregateByKey(pairs, 0, func(a, b int) int { return a + b }, nil).Slice()
	if len(agg) != 2 {
		t.Fatalf("aggregateByKey mismatch: %v", agg)
	}

	combined := CombineByKey(
		pairs,
		func(v int) int { return v },
		func(c, v int) int { return c + v },
		nil,
	).Slice()
	if len(combined) != 2 {
		t.Fatalf("combineByKey mismatch: %v", combined)
	}

	keyed := KeyBy(From("go", "java"), func(v string) int { return len(v) }).Slice()
	if len(keyed) != 2 || keyed[0].Key != 2 {
		t.Fatalf("keyBy mismatch: %v", keyed)
	}

	vals := Lookup(pairs, "a")
	if !reflect.DeepEqual(vals, []int{1, 2}) {
		t.Fatalf("lookup mismatch: %v", vals)
	}

	asMap := CollectAsMap(pairs)
	if asMap["a"] != 2 || asMap["b"] != 3 {
		t.Fatalf("collectAsMap mismatch: %v", asMap)
	}

	mp := MapPartitions(From(1, 2, 3), func(part []int) []int {
		return []int{len(part)}
	}).Slice()
	if !reflect.DeepEqual(mp, []int{3}) {
		t.Fatalf("mapPartitions mismatch: %v", mp)
	}

	mpi := MapPartitionsWithIndex(From(1, 2), func(idx int, part []int) []int {
		return []int{idx, len(part)}
	}).Slice()
	if !reflect.DeepEqual(mpi, []int{0, 2}) {
		t.Fatalf("mapPartitionsWithIndex mismatch: %v", mpi)
	}

	hit := 0
	ForEachPartition(From(1, 2, 3), func(part []int) {
		hit = len(part)
	})
	if hit != 3 {
		t.Fatalf("forEachPartition mismatch: %d", hit)
	}
}

func TestGoPlusContextAndChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	produced := 0
	stream := FromContext(ctx, func(ctx context.Context, emit func(int) bool) {
		for i := 1; i <= 10; i++ {
			produced++
			if i == 4 {
				cancel()
			}
			if !emit(i) {
				return
			}
		}
	})

	got := stream.Slice()
	if !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("fromContext mismatch: %v", got)
	}
	if produced < 3 {
		t.Fatalf("fromContext should produce values before cancellation: %d", produced)
	}

	src := make(chan int, 3)
	src <- 1
	src <- 2
	src <- 3
	close(src)
	fromCh := FromChan((<-chan int)(src)).Slice()
	if !reflect.DeepEqual(fromCh, []int{1, 2, 3}) {
		t.Fatalf("fromChan mismatch: %v", fromCh)
	}

	ch := From(4, 5, 6).ToChan(1)
	gotCh := make([]int, 0, 3)
	for v := range ch {
		gotCh = append(gotCh, v)
	}
	if !reflect.DeepEqual(gotCh, []int{4, 5, 6}) {
		t.Fatalf("toChan mismatch: %v", gotCh)
	}

	a := make(chan int, 2)
	b := make(chan int, 2)
	a <- 10
	a <- 11
	b <- 20
	b <- 21
	close(a)
	close(b)
	merged := make([]int, 0, 4)
	for v := range MergeChan((<-chan int)(a), (<-chan int)(b)) {
		merged = append(merged, v)
	}
	if len(merged) != 4 {
		t.Fatalf("mergeChan size mismatch: %v", merged)
	}
}

func TestGoPlusErrorAndTapError(t *testing.T) {
	_, err := MapE(From(1, 2, 3), func(v int) (int, error) {
		if v == 2 {
			return 0, errors.New("boom")
		}
		return v * 2, nil
	})
	if err == nil {
		t.Fatalf("mapE should fail")
	}

	flat, err := FlatMapE(From(1, 2), func(v int) ([]int, error) {
		return []int{v, v * 10}, nil
	})
	if err != nil {
		t.Fatalf("flatMapE unexpected err: %v", err)
	}
	if got := flat.Slice(); !reflect.DeepEqual(got, []int{1, 10, 2, 20}) {
		t.Fatalf("flatMapE mismatch: %v", got)
	}

	forEachErr := From(1, 2, 3).ForEachE(func(v int) error {
		if v == 3 {
			return errors.New("stop")
		}
		return nil
	})
	if forEachErr == nil {
		t.Fatalf("forEachE should fail")
	}

	collected, err := CollectE(From(1, 2, 3), func() []int { return make([]int, 0) }, func(dst *[]int, v int) error {
		*dst = append(*dst, v)
		return nil
	}, nil)
	if err != nil || !reflect.DeepEqual(collected, []int{1, 2, 3}) {
		t.Fatalf("collectE mismatch: data=%v err=%v", collected, err)
	}

	hit := false
	_ = TapError(errors.New("x"), func(error) { hit = true })
	if !hit {
		t.Fatalf("tapError should invoke handler")
	}
}

func TestGoPlusConcurrencyAPIs(t *testing.T) {
	res := MapPar(From(1, 2, 3), 2, func(v int) int { return v * 3 }).Slice()
	if !reflect.DeepEqual(res, []int{3, 6, 9}) {
		t.Fatalf("parallelMap mismatch: %v", res)
	}

	flat := FlatMapPar(From(1, 2, 3), 2, func(v int) []int { return []int{v, -v} }).Slice()
	if !reflect.DeepEqual(flat, []int{1, -1, 2, -2, 3, -3}) {
		t.Fatalf("parallelFlatMap mismatch: %v", flat)
	}
}

func TestGoPlusTimeAndBackpressure(t *testing.T) {
	window := Window(From(1, 2, 3, 4, 5), 2).Slice()
	if len(window) != 3 || !reflect.DeepEqual(window[0], []int{1, 2}) || !reflect.DeepEqual(window[2], []int{5}) {
		t.Fatalf("window mismatch: %v", window)
	}

	throttled := From(1, 2, 3).Throttle(20 * time.Millisecond).Slice()
	if len(throttled) != 1 {
		t.Fatalf("throttle mismatch: %v", throttled)
	}

	debounced := From(1, 2, 3).Debounce(10 * time.Millisecond).Slice()
	if !reflect.DeepEqual(debounced, []int{3}) {
		t.Fatalf("debounce mismatch: %v", debounced)
	}

	sampled := From(8, 9, 10).SampleEvery(5 * time.Millisecond).Slice()
	if len(sampled) == 0 {
		t.Fatalf("sampleEvery should emit data")
	}

	if got := From(1, 2, 3).Buffer(2).Slice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("buffer mismatch: %v", got)
	}
	if got := From(1, 2, 3).BlockWhenFull(1).Slice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("blockWhenFull mismatch: %v", got)
	}
	dropped := From(1, 2, 3, 4, 5).DropWhenFull(1).Slice()
	if len(dropped) == 0 {
		t.Fatalf("dropWhenFull should keep at least one item")
	}
	latest := From(1, 2, 3, 4).LatestOnly().Slice()
	if len(latest) == 0 || latest[len(latest)-1] != 4 {
		t.Fatalf("latestOnly mismatch: %v", latest)
	}
}

func TestGoPlusIOAndObserve(t *testing.T) {
	lines := FromReaderLines(strings.NewReader("a\nb\nc\n")).Slice()
	if !reflect.DeepEqual(lines, []string{"a", "b", "c"}) {
		t.Fatalf("fromReaderLines mismatch: %v", lines)
	}

	var buf bytes.Buffer
	err := ToWriter(From(1, 2, 3), &buf, func(v int) string { return strings.Repeat("x", v) })
	if err != nil || buf.String() != "xxxxxx" {
		t.Fatalf("toWriter mismatch: out=%q err=%v", buf.String(), err)
	}

	bres := MapBytes(From([]byte("a"), []byte("b")), func(b []byte) []byte {
		return append(b, '!')
	}).Slice()
	if string(bres[0]) != "a!" || string(bres[1]) != "b!" {
		t.Fatalf("mapBytes mismatch: %q %q", string(bres[0]), string(bres[1]))
	}

	ResetMetrics()
	traceHit := false
	SetTraceSink(func(ev TraceEvent) {
		if ev.Span == "demo" && ev.Count > 0 {
			traceHit = true
		}
	})

	_ = From(1, 2, 3).
		Tap(func(int) {}).
		WithMetrics("m1").
		WithTrace("demo").
		Slice()

	if GetMetricCount("m1") != 3 {
		t.Fatalf("metrics mismatch: %d", GetMetricCount("m1"))
	}
	if !traceHit {
		t.Fatalf("trace should be observed")
	}
}

func TestGoPlusFanOut(t *testing.T) {
	in := make(chan int, 3)
	in <- 7
	in <- 8
	in <- 9
	close(in)
	outs := FanOut((<-chan int)(in), 2, 3)
	if len(outs) != 2 {
		t.Fatalf("fanOut size mismatch: %d", len(outs))
	}
	left := make([]int, 0, 3)
	right := make([]int, 0, 3)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for v := range outs[0] {
			left = append(left, v)
		}
	}()
	go func() {
		defer wg.Done()
		for v := range outs[1] {
			right = append(right, v)
		}
	}()
	wg.Wait()
	if !reflect.DeepEqual(left, []int{7, 8, 9}) || !reflect.DeepEqual(right, []int{7, 8, 9}) {
		t.Fatalf("fanOut mismatch: left=%v right=%v", left, right)
	}
}

func TestGoPlusCtxChannelControls(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := Generate(func() int { return 1 }).ToChanCtx(ctx, 0)
	_, ok := <-ch
	if !ok {
		t.Fatalf("toChanCtx should produce data before cancel")
	}
	cancel()
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatalf("toChanCtx should close after cancel")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("toChanCtx cancel timeout")
	}

	in := make(chan int)
	_ = FanOutDrop((<-chan int)(in), 2, 1)
	sendDone := make(chan struct{})
	go func() {
		defer close(sendDone)
		for i := 0; i < 100; i++ {
			in <- i
		}
		close(in)
	}()
	select {
	case <-sendDone:
		// pass
	case <-time.After(300 * time.Millisecond):
		t.Fatalf("fanOutDrop should not block producer with slow consumers")
	}

	mctx, mcancel := context.WithCancel(context.Background())
	mch := make(chan int)
	go func() {
		defer close(mch)
		for i := 0; i < 1000; i++ {
			mch <- i
		}
	}()
	merged := MergeChanCtx(mctx, (<-chan int)(mch))
	mcancel()
	select {
	case _, ok := <-merged:
		if ok {
			for range merged {
			}
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatalf("mergeChanCtx cancel timeout")
	}
}

func TestBizCommonOperators(t *testing.T) {
	type row struct {
		k string
		v int
	}
	rows := From(row{"a", 1}, row{"a", 2}, row{"b", 3})
	pairs := MapFn(KeyByFn(rows, func(r row) string { return r.k }), func(p Pair[string, row]) Pair[string, int] {
		return Pair[string, int]{Key: p.Key, Value: p.Value.v}
	})

	sum := SumByKey(pairs).Slice()
	if len(sum) != 2 {
		t.Fatalf("sumByKey size mismatch: %v", sum)
	}

	count := CountByKeyStream(pairs).Slice()
	if len(count) != 2 {
		t.Fatalf("countByKeyStream size mismatch: %v", count)
	}

	left := From(
		Pair[string, int]{Key: "a", Value: 1},
		Pair[string, int]{Key: "b", Value: 2},
	)
	right := From(
		Pair[string, string]{Key: "a", Value: "x"},
	)
	if got := InnerJoinByKey(left, right).Slice(); len(got) != 1 || got[0].Key != "a" {
		t.Fatalf("innerJoinByKey mismatch: %v", got)
	}
	if got := LeftJoinByKey(left, right).Slice(); len(got) != 2 {
		t.Fatalf("leftJoinByKey mismatch: %v", got)
	}

	tw := TumblingWindow(From(1, 2, 3, 4, 5), 2).Slice()
	if len(tw) != 3 || !reflect.DeepEqual(tw[0], []int{1, 2}) {
		t.Fatalf("tumblingWindow mismatch: %v", tw)
	}

	sw := SlidingWindow(From(1, 2, 3, 4), 3, 1).Slice()
	if len(sw) != 2 || !reflect.DeepEqual(sw[1], []int{2, 3, 4}) {
		t.Fatalf("slidingWindow mismatch: %v", sw)
	}

	wr := WindowReduce(From(1, 2, 3, 4), 2, 0, func(acc, v int) int { return acc + v }).Slice()
	if !reflect.DeepEqual(wr, []int{3, 7}) {
		t.Fatalf("windowReduce mismatch: %v", wr)
	}

	pr := Process(From(1, 2, 3, 4), func(v int) (int, bool) {
		if v%2 == 0 {
			return v * 10, true
		}
		return 0, false
	}).Slice()
	if !reflect.DeepEqual(pr, []int{20, 40}) {
		t.Fatalf("process mismatch: %v", pr)
	}

	pm := ProcessMany(From(1, 2), func(v int) []int { return []int{v, -v} }).Slice()
	if !reflect.DeepEqual(pm, []int{1, -1, 2, -2}) {
		t.Fatalf("processMany mismatch: %v", pm)
	}
}

func TestGoStyleDirectNames(t *testing.T) {
	v := 7
	if got := FromPointer(&v).Slice(); !reflect.DeepEqual(got, []int{7}) {
		t.Fatalf("fromPointer mismatch: %v", got)
	}

	chained := Chain(From(1, 2), From(3), Empty[int]()).All()
	if !reflect.DeepEqual(chained, []int{1, 2, 3}) {
		t.Fatalf("chain mismatch: %v", chained)
	}

	s := From(1, 2, 3, 4).Where(func(v int) bool { return v%2 == 0 })
	if s.Len() != 2 {
		t.Fatalf("len mismatch: %d", s.Len())
	}
	if !s.Any(func(v int) bool { return v == 2 }) || !s.Every(func(v int) bool { return v%2 == 0 }) {
		t.Fatalf("any/every mismatch")
	}
	if !s.None(func(v int) bool { return v%2 == 1 }) {
		t.Fatalf("none mismatch")
	}

	var total int
	From(1, 2, 3).Each(func(v int) { total += v })
	if total != 6 {
		t.Fatalf("each mismatch: %d", total)
	}

	par := MapPar(From(1, 2, 3), 2, func(v int) int { return v * 2 }).Slice()
	if !reflect.DeepEqual(par, []int{2, 4, 6}) {
		t.Fatalf("mapPar mismatch: %v", par)
	}

	uniq := Unique(From(1, 1, 2, 2, 3)).Slice()
	if !reflect.DeepEqual(uniq, []int{1, 2, 3}) {
		t.Fatalf("unique mismatch: %v", uniq)
	}

	mapped := MapFn(From(1, 2, 3), func(v int) int { return v + 1 }).Slice()
	if !reflect.DeepEqual(mapped, []int{2, 3, 4}) {
		t.Fatalf("mapFn mismatch: %v", mapped)
	}

	fmapped := FlatMapFn(From(1, 2), func(v int) Stream[int] { return From(v, -v) }).Slice()
	if !reflect.DeepEqual(fmapped, []int{1, -1, 2, -2}) {
		t.Fatalf("flatMapFn mismatch: %v", fmapped)
	}

	head, ok := From(9, 8).Head()
	if !ok || head != 9 {
		t.Fatalf("head mismatch: %d %v", head, ok)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ctxVals := FromCtx(ctx, func(ctx context.Context, emit func(int) bool) {
		_ = emit(1)
	}).Slice()
	if len(ctxVals) != 0 {
		t.Fatalf("fromCtx mismatch: %v", ctxVals)
	}

	batch := Batch(From(1, 2, 3), 2).Slice()
	if len(batch) != 2 || !reflect.DeepEqual(batch[0], []int{1, 2}) {
		t.Fatalf("batch mismatch: %v", batch)
	}
}
