
package main
import ("fmt"; "core/std")

var dp[][]int

// マルチバイトも含めた Cost 距離
func Cost(x []rune, y []rune) int {	
	var cost int
	dp := make([][]int, len(x) + 1)
	for i := 0; i <= len(x); i++ {dp[i] = make([]int, len(y) + 1)}
	for i := 0; i <= len(x); i++ {dp[i][0] = i}
	for i := 0; i <= len(y); i++ {dp[0][i] = i}
	
	for i := 1; i <= len(x); i++ {
		for j := 1; j <= len(y); j++ {
			if x[i - 1] == y[j - 1] {cost = 0} else {cost = 1}
			dp[i][j] = std.Min(dp[i][j - 1] + 1, dp[i - 1][j] + 1, dp[i - 1][j - 1] + cost)
		}
	}
	return dp[len(x)][len(y)]
}

func TestCostEn() {
	fmt.Println(Cost([]rune("sitting"), []rune("kitten")) == 3)
	fmt.Println(Cost([]rune("abcabccdaba"), []rune("abracdabra")) == 4)
	fmt.Println(Cost(
		[]rune("Good Evening World! This is very nice posts!"),
		[]rune("Hello World! This is amazing posts!")) == 19)
}

func TestCostJa() {
	fmt.Println(Cost([]rune("こんにちは"), []rune("こんばんは")) == 2)
	fmt.Println(Cost([]rune("他の言語同様に名前は重要です"), []rune("他の言語と異なり名前は不要です")) == 5)
	fmt.Println(Cost(
		[]rune("これがどのように行われるか詳細は言語仕様を見ていただきたいのです。"),
		[]rune("この仕組みは、セミコロンのないすっきりしたコードを書くのに役立っています。")) == 33)
}

func TestCostJaAndEn() {
	fmt.Println(Cost([]rune("あra"), []rune("こha")) == 2)
	fmt.Println(Cost([]rune("あraGo"), []rune("こha")) == 4) // -> 3, error
	fmt.Println(Cost(
		[]rune("これがどのように行われるか詳細はGo言語仕様を見ていただきたいのです。"),
		[]rune("この仕組みは、セミコロンのないすっきりしたコードを書くのに役立っています。")) == 33)
	fmt.Println(Cost(
		[]rune("Levenshutein 距離を計算するアルゴリズムの作成"),
		[]rune("距離を計算するアルゴを作る")) == 18) // -> 13, error
}

func Test() {	
	TestCostEn()
	TestCostJa()
	TestCostJaAndEn()
}

func main() {
	std.ReadFile("test")
}



