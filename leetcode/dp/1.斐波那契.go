package dp

// 动态规划实现
// 1,2,3,4,5,6,7, 8, 9, 10
// 1,1,2,3,5,8,13,21,34,55
func fibonacci(n int) int {
	if n <= 2 {
		return n
	}
	dp := make([]int, n) // 初始dp table
	dp[0] = 1
	dp[1] = 1
	for i := 2; i < n; i++ {
		dp[i] = dp[i-2] + dp[i-1] // 确定状态转移公式
	}
	return dp[n-1]
}