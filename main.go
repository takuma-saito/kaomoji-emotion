
package main
import ("fmt"; "core/std")

var dp[][]int

// 単純な Cost 距離
func Cost(x string, y string) int {	
	var cost int
	dp := make([][]int, len(x) + 1)
	for i, _ := range dp {dp[i] = make([]int, len(y) + 1)}
	for i, _ := range x {dp[i][0] = i}
	for i, _ := range y {dp[0][i] = i}
	
	for i := 1; i <= len(x); i++ {
		for j := 1; j <= len(y); j++ {
			if x[i - 1] == y[j - 1] {cost = 0} else {cost = 1}
			dp[i][j] = std.Min(dp[i][j - 1] + 1, dp[i - 1][j] + 1, dp[i - 1][j - 1] + cost)
		}
	}
	return dp[len(x)][len(y)]
}

func TestCostEnglish() {
	fmt.Println(Cost("sitting", "kitten") == 3)
	fmt.Println(Cost("abcabccdaba", "abracdabra") == 4)
	fmt.Println(Cost(
		"Good Evening World! This is very nice posts!",
		"Hello World! This is amazing posts!") == 19)
}

func main() {
	TestCostEnglish()
}



