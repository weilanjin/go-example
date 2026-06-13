package pagination

import (
	"context"
	"fmt"
)

// 基于游标（Cursor）的滚动查询（Seek Pagination）
// 每次查询都记住上一批最后一条数据的特定标识（游标），下一次查询从这个标识往后找。

/*
	-- 每次查询利用 B+ 树索引直接定位起点，速度极快且稳定
	SELECT * FROM operate_log WHERE id > {last_max_id} ORDER BY id ASC LIMIT 1000;
*/

const DefaultBatchSize = 1000

// SeekPage 通用游标查询探针处理器
func SeekPage[T any, C comparable](
	ctx context.Context,
	cursor C, limit int,
	fetchFunc func(ctx context.Context, cursor C, fetchLimit int) ([]T, error), // 剥离查询逻辑
	getCursorFunc func(item T) C, // 剥离游标提取逻辑
) ([]T, *C, error) {
	if fetchFunc == nil || getCursorFunc == nil {
		return make([]T, 0), new(C), nil // 参数校验失败，返回空结果
	}
	// 1. 安全防护：限制最大单次查询量
	if limit <= 0 || limit > DefaultBatchSize {
		limit = DefaultBatchSize
	}

	// 2. 核心探针：向数据库发起请求，但强制将 limit 设为 limit + 1 [cite: 8]
	fetchLimit := limit + 1

	// 3. 调用外部业务实际的数据库查询逻辑
	rows, err := fetchFunc(ctx, cursor, fetchLimit)
	if err != nil {
		return nil, nil, fmt.Errorf("查询数据失败: %w", err)
	}
	if len(rows) == 0 {
		return make([]T, 0), new(C), nil // 无数据直接返回
	}

	// 4. “limit + 1” 判断逻辑
	var nextCursor *C
	if len(rows) > limit {
		// 结果集长度 > limit，说明数据库里至少还有更多数据 [cite: 9]
		// 截取前 limit 条数据返回给前端，设置 has_more = true [cite: 9]
		rows = rows[:limit]

		// 提取截取后最后一条数据的游标
		nextCursor = new(getCursorFunc(rows[len(rows)-1]))
	}
	return rows, nextCursor, nil
}

// SeekPageWithFilter 通用游标查询探针处理器，支持额外的查询过滤条件
func SeekPageWithFilter[T any, C comparable, F any](
	ctx context.Context,
	cursor C,
	limit int,
	filter F, // 查询过滤条件
	fetchFunc func(ctx context.Context, cursor C, filter F, fetchLimit int) ([]T, error), // 实际查询逻辑
	getCursorFunc func(item T) C, // 游标提取逻辑
) ([]T, *C, error) {
	if fetchFunc == nil || getCursorFunc == nil {
		return make([]T, 0), new(C), nil // 参数校验失败，返回空结果
	}

	// 1. 安全防护：限制最大单次查询量
	if limit <= 0 || limit > DefaultBatchSize {
		limit = DefaultBatchSize
	}

	// 2. 核心探针：向数据库发起请求，但强制将 limit 设为 limit + 1 [cite: 8]
	fetchLimit := limit + 1

	// 3. 调用外部业务实际的数据库查询逻辑
	rows, err := fetchFunc(ctx, cursor, filter, fetchLimit)
	if err != nil {
		return nil, nil, fmt.Errorf("查询数据失败: %w", err)
	}
	if len(rows) == 0 {
		return make([]T, 0), new(C), nil // 无数据直接返回
	}

	// 4. “limit + 1” 判断逻辑
	var nextCursor *C
	if len(rows) > limit {
		// 结果集长度 > limit，说明数据库里至少还有更多数据 [cite: 9]
		// 截取前 limit 条数据返回给前端，设置 has_more = true [cite: 9]
		rows = rows[:limit]

		// 提取截取后最后一条数据的游标
		nextCursor = new(getCursorFunc(rows[len(rows)-1]))
	}
	return rows, nextCursor, nil
}
