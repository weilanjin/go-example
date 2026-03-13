package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestStreamChain(t *testing.T) {
	got := Map(
		Distinct(
			Of(1, 2, 2, 3, 4, 5).
				Filter(func(v int) bool { return v > 1 }).
				Skip(1).
				Limit(3),
		),
		func(v int) int { return v * 10 },
	).CollectToSlice()

	want := []int{20, 30, 40}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected chain result, got=%v want=%v", got, want)
	}
}

func TestMapFlatMap(t *testing.T) {
	mapped := Map(Of(1, 2, 3), func(v int) string {
		return string(rune('a' + v - 1))
	}).CollectToSlice()

	if !reflect.DeepEqual(mapped, []string{"a", "b", "c"}) {
		t.Fatalf("unexpected mapped result: %v", mapped)
	}

	flat := FlatMap(Of(1, 2, 3), func(v int) Stream[int] {
		return Of(v, v*10)
	}).CollectToSlice()

	if !reflect.DeepEqual(flat, []int{1, 10, 2, 20, 3, 30}) {
		t.Fatalf("unexpected flat result: %v", flat)
	}
}

func TestSortedReduceAndCollectToMap(t *testing.T) {
	sorted := Of(3, 1, 2).Sorted(func(a, b int) bool { return a < b }).CollectToSlice()
	if !reflect.DeepEqual(sorted, []int{1, 2, 3}) {
		t.Fatalf("unexpected sorted result: %v", sorted)
	}

	sum := Of(1, 2, 3, 4).Reduce(0, func(acc, v int) int { return acc + v })
	if sum != 10 {
		t.Fatalf("unexpected reduce result: %d", sum)
	}

	m := CollectToMap(Of("aa", "bbb", "c"), func(v string) int { return len(v) }, func(v string) string { return v })
	if len(m) != 3 || m[1] != "c" || m[2] != "aa" || m[3] != "bbb" {
		t.Fatalf("unexpected map result: %v", m)
	}
}

func TestLazyEvaluation(t *testing.T) {
	hit := 0
	s := Map(
		Of(1, 2, 3).Filter(func(v int) bool {
			hit++
			return v%2 == 1
		}),
		func(v int) int { return v * 2 },
	)

	if hit != 0 {
		t.Fatalf("lazy broken before terminal, hit=%d", hit)
	}

	_ = s.CollectToSlice()
	if hit != 3 {
		t.Fatalf("unexpected traversal count, hit=%d", hit)
	}
}

func TestShortCircuit(t *testing.T) {
	countAny := 0
	any := Of(1, 2, 3, 4, 5).AnyMatch(func(v int) bool {
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
	all := Of(2, 4, 5, 6).AllMatch(func(v int) bool {
		countAll++
		return v%2 == 0
	})
	if all {
		t.Fatalf("allMatch should be false")
	}
	if countAll != 3 {
		t.Fatalf("allMatch should short-circuit at first odd, count=%d", countAll)
	}

	first, ok := Of(9, 8, 7).FindFirst()
	if !ok || first != 9 {
		t.Fatalf("findFirst unexpected, first=%d ok=%v", first, ok)
	}
}

func TestLimitSkipCount(t *testing.T) {
	got := Of(1, 2, 3, 4, 5).Skip(2).Limit(2).CollectToSlice()
	if !reflect.DeepEqual(got, []int{3, 4}) {
		t.Fatalf("unexpected skip/limit result: %v", got)
	}

	if c := Of(1, 2, 3).Count(); c != 3 {
		t.Fatalf("unexpected count result: %d", c)
	}

	if c := Of(1, 2, 3).Filter(func(v int) bool { return v > 1 }).Count(); c != 2 {
		t.Fatalf("unexpected filtered count result: %d", c)
	}

	if c := Of(1, 2, 3).Limit(0).Count(); c != 0 {
		t.Fatalf("unexpected limit zero count result: %d", c)
	}
}

func TestDistinctStability(t *testing.T) {
	got := Distinct(Of(3, 1, 3, 2, 1, 2)).CollectToSlice()
	want := []int{3, 1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected distinct result, got=%v want=%v", got, want)
	}

	res := Distinct(Of(5, 4, 4, 3, 2, 2, 1)).Sorted(func(a, b int) bool { return a < b }).CollectToSlice()
	if !sort.IntsAreSorted(res) {
		t.Fatalf("sorted distinct result should be ordered: %v", res)
	}
}

func TestFactoriesAndBuilder(t *testing.T) {
	if c := Empty[int]().Count(); c != 0 {
		t.Fatalf("empty count should be 0, got=%d", c)
	}

	if c := OfNullable[int](nil).Count(); c != 0 {
		t.Fatalf("ofNullable nil count should be 0, got=%d", c)
	}

	v := 42
	if got := OfNullable(&v).CollectToSlice(); !reflect.DeepEqual(got, []int{42}) {
		t.Fatalf("ofNullable value mismatch: %v", got)
	}

	b := NewBuilder[int]().Add(1).Add(2).Add(3)
	if got := b.Build().CollectToSlice(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("builder mismatch: %v", got)
	}

	joined := Concat(Of(1, 2), Of(3, 4)).CollectToSlice()
	if !reflect.DeepEqual(joined, []int{1, 2, 3, 4}) {
		t.Fatalf("concat mismatch: %v", joined)
	}
}

func TestTakeDropPeekAndNoneMatch(t *testing.T) {
	peeked := 0
	res := Of(1, 2, 3, 4, 5).
		Peek(func(int) { peeked++ }).
		TakeWhile(func(v int) bool { return v < 4 }).
		DropWhile(func(v int) bool { return v < 2 }).
		CollectToSlice()

	if !reflect.DeepEqual(res, []int{2, 3}) {
		t.Fatalf("take/drop mismatch: %v", res)
	}
	if peeked != 4 {
		t.Fatalf("peek should observe until takeWhile stop, peeked=%d", peeked)
	}

	if !Of(2, 4, 6).NoneMatch(func(v int) bool { return v%2 == 1 }) {
		t.Fatalf("noneMatch should be true")
	}
}

func TestGenerateIterateAndSortedByKey(t *testing.T) {
	gen := Generate(func() int { return 7 }).Limit(3).CollectToSlice()
	if !reflect.DeepEqual(gen, []int{7, 7, 7}) {
		t.Fatalf("generate mismatch: %v", gen)
	}

	it := Iterate(1, func(v int) int { return v + 1 }).Limit(5).CollectToSlice()
	if !reflect.DeepEqual(it, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("iterate mismatch: %v", it)
	}

	itw := IterateWhile(1, func(v int) bool { return v <= 3 }, func(v int) int { return v + 1 }).CollectToSlice()
	if !reflect.DeepEqual(itw, []int{1, 2, 3}) {
		t.Fatalf("iterateWhile mismatch: %v", itw)
	}

	type item struct {
		name string
		key  int
	}
	sorted := SortedByKey(Of(item{"b", 2}, item{"a", 1}), func(v item) int { return v.key }).CollectToSlice()
	if sorted[0].name != "a" {
		t.Fatalf("sortedByKey mismatch: %v", sorted)
	}
}

func TestOptionalAndMinMax(t *testing.T) {
	opt := Of(3, 1, 2).ReduceOptional(func(a, b int) int { return a + b })
	if v := opt.OrElse(0); v != 6 {
		t.Fatalf("reduceOptional mismatch: %d", v)
	}

	minV := MinOrdered(Of(3, 1, 2)).OrElse(0)
	maxV := MaxOrdered(Of(3, 1, 2)).OrElse(0)
	if minV != 1 || maxV != 3 {
		t.Fatalf("min/max mismatch min=%d max=%d", minV, maxV)
	}

	collected := Collect(Of(1, 2, 3), func() []int { return make([]int, 0) }, func(dst *[]int, v int) {
		*dst = append(*dst, v)
	}, nil)
	if !reflect.DeepEqual(collected, []int{1, 2, 3}) {
		t.Fatalf("collect mismatch: %v", collected)
	}

	reduced := ReduceWithCombiner(Of(1, 2, 3), 0, func(acc, v int) int { return acc + v }, nil)
	if reduced != 6 {
		t.Fatalf("reduceWithCombiner mismatch: %d", reduced)
	}
}

func TestStateAndClose(t *testing.T) {
	s := Of(1, 2, 3).Parallel()
	if !s.IsParallel() {
		t.Fatalf("parallel state mismatch")
	}
	if s.Sequential().IsParallel() {
		t.Fatalf("sequential should clear parallel state")
	}

	steps := make([]int, 0, 2)
	Of(1).
		OnClose(func() { steps = append(steps, 1) }).
		OnClose(func() { steps = append(steps, 2) }).
		Close()
	if !reflect.DeepEqual(steps, []int{2, 1}) {
		t.Fatalf("close order mismatch: %v", steps)
	}

	any, ok := Of(8, 9).FindAny()
	if !ok || any != 8 {
		t.Fatalf("findAny mismatch: any=%d ok=%v", any, ok)
	}
}

func TestCollectorsFamily(t *testing.T) {
	count := CollectWith(Of(1, 2, 3, 4), CountingCollector[int]())
	if count != 4 {
		t.Fatalf("counting collector mismatch: %d", count)
	}

	joined := CollectWith(Of("a", "b", "c"), JoiningCollector("-", "[", "]"))
	if joined != "[a-b-c]" {
		t.Fatalf("joining collector mismatch: %s", joined)
	}

	set := CollectWith(Of(1, 2, 2, 3), ToSetCollector[int]())
	if len(set) != 3 {
		t.Fatalf("set collector mismatch: %v", set)
	}

	grouped := CollectWith(Of("go", "java", "c"), GroupingByCollector(func(v string) int { return len(v) }))
	if !reflect.DeepEqual(grouped[2], []string{"go"}) || !reflect.DeepEqual(grouped[1], []string{"c"}) {
		t.Fatalf("groupingBy collector mismatch: %v", grouped)
	}

	parts := CollectWith(Of(1, 2, 3, 4), PartitioningByCollector(func(v int) bool { return v%2 == 0 }))
	if !reflect.DeepEqual(parts[true], []int{2, 4}) || !reflect.DeepEqual(parts[false], []int{1, 3}) {
		t.Fatalf("partitioningBy collector mismatch: %v", parts)
	}
}

func TestCollectorsComposed(t *testing.T) {
	mapped := CollectWith(
		Of(1, 2, 3),
		MappingCollector(func(v int) string { return string(rune('a' + v - 1)) }, ToSliceCollector[string]()),
	)
	if !reflect.DeepEqual(mapped, []string{"a", "b", "c"}) {
		t.Fatalf("mapping collector mismatch: %v", mapped)
	}

	filtered := CollectWith(
		Of(1, 2, 3, 4),
		FilteringCollector(func(v int) bool { return v%2 == 0 }, ToSliceCollector[int]()),
	)
	if !reflect.DeepEqual(filtered, []int{2, 4}) {
		t.Fatalf("filtering collector mismatch: %v", filtered)
	}

	flatMapped := CollectWith(
		Of("a,b", "c"),
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
		Of("a", "aa", "b"),
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
		Of("go", "java", "js"),
		GroupingByMappingCollector(func(v string) int { return len(v) }, func(v string) byte { return v[0] }),
	)
	if !reflect.DeepEqual(groupMap[2], []byte{'g', 'j'}) {
		t.Fatalf("groupingByMapping collector mismatch: %v", groupMap)
	}
}

func TestCollectorsStatistics(t *testing.T) {
	sumInt := CollectWith(Of(1, 2, 3), SummingIntCollector(func(v int) int { return v }))
	if sumInt != 6 {
		t.Fatalf("summingInt mismatch: %d", sumInt)
	}

	sumI64 := CollectWith(Of(int64(1), int64(2), int64(3)), SummingInt64Collector(func(v int64) int64 { return v }))
	if sumI64 != 6 {
		t.Fatalf("summingInt64 mismatch: %d", sumI64)
	}

	sumF64 := CollectWith(Of(1.5, 2.5), SummingFloat64Collector(func(v float64) float64 { return v }))
	if sumF64 != 4.0 {
		t.Fatalf("summingFloat64 mismatch: %f", sumF64)
	}

	avgInt := CollectWith(Of(1, 2, 3, 4), AveragingIntCollector(func(v int) int { return v }))
	if avgInt != 2.5 {
		t.Fatalf("averagingInt mismatch: %f", avgInt)
	}

	avgI64 := CollectWith(Of(int64(2), int64(4)), AveragingInt64Collector(func(v int64) int64 { return v }))
	if avgI64 != 3.0 {
		t.Fatalf("averagingInt64 mismatch: %f", avgI64)
	}

	avgF64 := CollectWith(Of(1.0, 2.0, 5.0), AveragingFloat64Collector(func(v float64) float64 { return v }))
	if avgF64 != 8.0/3.0 {
		t.Fatalf("averagingFloat64 mismatch: %f", avgF64)
	}

	intStats := CollectWith(Of(3, 1, 2), SummarizingIntCollector(func(v int) int { return v }))
	if intStats.Count != 3 || intStats.Sum != 6 || intStats.Min != 1 || intStats.Max != 3 || intStats.Average() != 2.0 {
		t.Fatalf("summarizingInt mismatch: %+v", intStats)
	}

	longStats := CollectWith(Of(int64(7), int64(9)), SummarizingInt64Collector(func(v int64) int64 { return v }))
	if longStats.Count != 2 || longStats.Sum != 16 || longStats.Min != 7 || longStats.Max != 9 || longStats.Average() != 8.0 {
		t.Fatalf("summarizingInt64 mismatch: %+v", longStats)
	}

	doubleStats := CollectWith(Of(1.5, 2.5, 0.5), SummarizingFloat64Collector(func(v float64) float64 { return v }))
	if doubleStats.Count != 3 || doubleStats.Sum != 4.5 || doubleStats.Min != 0.5 || doubleStats.Max != 2.5 || doubleStats.Average() != 1.5 {
		t.Fatalf("summarizingFloat64 mismatch: %+v", doubleStats)
	}
}

func TestCollectorsAndThenAndReducing(t *testing.T) {
	collector := CollectingAndThen(ToSliceCollector[int](), func(items []int) int {
		return len(items)
	})
	lenRes := CollectWith(Of(1, 2, 3), collector)
	if lenRes != 3 {
		t.Fatalf("collectingAndThen mismatch: %d", lenRes)
	}

	reduced := CollectWith(
		Of("a", "bb", "ccc"),
		ReducingCollector(0, func(v string) int { return len(v) }, func(a, b int) int { return a + b }),
	)
	if reduced != 6 {
		t.Fatalf("reducing collector mismatch: %d", reduced)
	}
}

func TestCollectorsMinMaxGroupingDownstreamAndTeeing(t *testing.T) {
	minOpt := CollectWith(Of(3, 1, 2), MinByCollector(func(a, b int) bool { return a < b }))
	if v := minOpt.OrElse(0); v != 1 {
		t.Fatalf("minBy collector mismatch: %d", v)
	}

	maxOpt := CollectWith(Of(3, 1, 2), MaxByCollector(func(a, b int) bool { return a < b }))
	if v := maxOpt.OrElse(0); v != 3 {
		t.Fatalf("maxBy collector mismatch: %d", v)
	}

	groupedSum := CollectWith(
		Of(1, 2, 3, 4, 5),
		GroupingByDownstreamCollector(
			func(v int) bool { return v%2 == 0 },
			SummingIntCollector(func(v int) int { return v }),
		),
	)
	if groupedSum[true] != 6 || groupedSum[false] != 9 {
		t.Fatalf("groupingByDownstream collector mismatch: %v", groupedSum)
	}

	avgLen := CollectWith(
		Of("a", "bb", "ccc"),
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
	if got := Union(Of(1, 2), Of(3, 4)).CollectToSlice(); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
		t.Fatalf("union mismatch: %v", got)
	}

	inter := Intersection(Of(1, 2, 3, 3), Of(3, 4)).CollectToSlice()
	if !reflect.DeepEqual(inter, []int{3, 3}) {
		t.Fatalf("intersection mismatch: %v", inter)
	}

	sub := Subtract(Of(1, 2, 3), Of(2)).CollectToSlice()
	if !reflect.DeepEqual(sub, []int{1, 3}) {
		t.Fatalf("subtract mismatch: %v", sub)
	}

	car := Cartesian(Of(1, 2), Of("a", "b")).CollectToSlice()
	if len(car) != 4 {
		t.Fatalf("cartesian mismatch: %v", car)
	}

	zipped := Zip(Of(1, 2, 3), Of("a", "b")).CollectToSlice()
	if !reflect.DeepEqual(zipped, []Tuple2[int, string]{{1, "a"}, {2, "b"}}) {
		t.Fatalf("zip mismatch: %v", zipped)
	}

	idx := ZipWithIndex(Of("x", "y")).CollectToSlice()
	if !reflect.DeepEqual(idx, []Tuple2[string, int]{{"x", 0}, {"y", 1}}) {
		t.Fatalf("zipWithIndex mismatch: %v", idx)
	}
}

func TestSparkLikeSortSampleAndTake(t *testing.T) {
	sorted := SortBy(Of("bbb", "a", "cc"), func(v string) int { return len(v) }, true).CollectToSlice()
	if !reflect.DeepEqual(sorted, []string{"a", "cc", "bbb"}) {
		t.Fatalf("sortBy mismatch: %v", sorted)
	}

	sampled := Sample(Of(1, 2, 3, 4, 5), false, 1.0, 1).CollectToSlice()
	if !reflect.DeepEqual(sampled, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("sample mismatch: %v", sampled)
	}

	glom := Glom(Of(1, 2, 3, 4, 5), 2).CollectToSlice()
	if len(glom) != 3 || !reflect.DeepEqual(glom[0], []int{1, 2}) || !reflect.DeepEqual(glom[2], []int{5}) {
		t.Fatalf("glom mismatch: %v", glom)
	}

	take := Of(5, 4, 3).Take(2)
	if !reflect.DeepEqual(take, []int{5, 4}) {
		t.Fatalf("take mismatch: %v", take)
	}

	first, ok := Of(9, 8).First()
	if !ok || first != 9 {
		t.Fatalf("first mismatch: %d %v", first, ok)
	}

	if !reflect.DeepEqual(TakeOrdered(Of(3, 1, 2), 2), []int{1, 2}) {
		t.Fatalf("takeOrdered mismatch")
	}
	if !reflect.DeepEqual(Top(Of(3, 1, 2), 2), []int{3, 2}) {
		t.Fatalf("top mismatch")
	}
}

func TestSparkLikePairOperators(t *testing.T) {
	pairs := Of(
		Pair[string, int]{Key: "a", Value: 1},
		Pair[string, int]{Key: "a", Value: 2},
		Pair[string, int]{Key: "b", Value: 3},
	)

	mappedVals := MapValues(pairs, func(v int) int { return v * 10 }).CollectToSlice()
	if mappedVals[0].Value != 10 || mappedVals[2].Value != 30 {
		t.Fatalf("mapValues mismatch: %v", mappedVals)
	}

	flatVals := FlatMapValues(pairs, func(v int) []int { return []int{v, -v} }).CollectToSlice()
	if len(flatVals) != 6 {
		t.Fatalf("flatMapValues mismatch: %v", flatVals)
	}

	grp := GroupByKey(pairs).CollectToSlice()
	if len(grp) != 2 {
		t.Fatalf("groupByKey mismatch: %v", grp)
	}

	reduced := ReduceByKey(pairs, func(a, b int) int { return a + b }).CollectToSlice()
	if len(reduced) != 2 {
		t.Fatalf("reduceByKey mismatch: %v", reduced)
	}

	folded := FoldByKey(pairs, 10, func(a, b int) int { return a + b }).CollectToSlice()
	if len(folded) != 2 {
		t.Fatalf("foldByKey mismatch: %v", folded)
	}

	sortedPair := SortByKey(Of(
		Pair[int, string]{Key: 2, Value: "b"},
		Pair[int, string]{Key: 1, Value: "a"},
	), true).CollectToSlice()
	if sortedPair[0].Key != 1 {
		t.Fatalf("sortByKey mismatch: %v", sortedPair)
	}

	if cnt := CountByKey(pairs); cnt["a"] != 2 || cnt["b"] != 1 {
		t.Fatalf("countByKey mismatch: %v", cnt)
	}
	if cnt := CountByValue(Of(1, 1, 2)); cnt[1] != 2 || cnt[2] != 1 {
		t.Fatalf("countByValue mismatch: %v", cnt)
	}

	ks := Keys(pairs).CollectToSlice()
	vs := Values(pairs).CollectToSlice()
	if len(ks) != 3 || len(vs) != 3 {
		t.Fatalf("keys/values mismatch: %v %v", ks, vs)
	}
}

func TestSparkLikeJoinFamilyAndCogroup(t *testing.T) {
	left := Of(
		Pair[string, int]{Key: "a", Value: 1},
		Pair[string, int]{Key: "b", Value: 2},
	)
	right := Of(
		Pair[string, string]{Key: "a", Value: "x"},
		Pair[string, string]{Key: "c", Value: "y"},
	)

	inner := Join(left, right).CollectToSlice()
	if len(inner) != 1 || inner[0].Key != "a" || inner[0].Value.First != 1 || inner[0].Value.Second != "x" {
		t.Fatalf("join mismatch: %v", inner)
	}

	loj := LeftOuterJoin(left, right).CollectToSlice()
	if len(loj) != 2 {
		t.Fatalf("leftOuterJoin size mismatch: %v", loj)
	}

	roj := RightOuterJoin(left, right).CollectToSlice()
	if len(roj) != 2 {
		t.Fatalf("rightOuterJoin size mismatch: %v", roj)
	}

	foj := FullOuterJoin(left, right).CollectToSlice()
	if len(foj) != 3 {
		t.Fatalf("fullOuterJoin size mismatch: %v", foj)
	}

	cg := Cogroup(left, right).CollectToSlice()
	if len(cg) != 3 {
		t.Fatalf("cogroup mismatch: %v", cg)
	}
}

func TestSparkLikeAggregateCombinePartitionAndLookup(t *testing.T) {
	pairs := Of(
		Pair[string, int]{Key: "a", Value: 1},
		Pair[string, int]{Key: "a", Value: 2},
		Pair[string, int]{Key: "b", Value: 3},
	)

	agg := AggregateByKey(pairs, 0, func(a, b int) int { return a + b }, nil).CollectToSlice()
	if len(agg) != 2 {
		t.Fatalf("aggregateByKey mismatch: %v", agg)
	}

	combined := CombineByKey(
		pairs,
		func(v int) int { return v },
		func(c, v int) int { return c + v },
		nil,
	).CollectToSlice()
	if len(combined) != 2 {
		t.Fatalf("combineByKey mismatch: %v", combined)
	}

	keyed := KeyBy(Of("go", "java"), func(v string) int { return len(v) }).CollectToSlice()
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

	mp := MapPartitions(Of(1, 2, 3), func(part []int) []int {
		return []int{len(part)}
	}).CollectToSlice()
	if !reflect.DeepEqual(mp, []int{3}) {
		t.Fatalf("mapPartitions mismatch: %v", mp)
	}

	mpi := MapPartitionsWithIndex(Of(1, 2), func(idx int, part []int) []int {
		return []int{idx, len(part)}
	}).CollectToSlice()
	if !reflect.DeepEqual(mpi, []int{0, 2}) {
		t.Fatalf("mapPartitionsWithIndex mismatch: %v", mpi)
	}

	hit := 0
	ForEachPartition(Of(1, 2, 3), func(part []int) {
		hit = len(part)
	})
	if hit != 3 {
		t.Fatalf("forEachPartition mismatch: %d", hit)
	}
}
