
package main
import (
	"fmt"
	"math"
	"strings"
	"strconv"
	"core/std"
)

const MAX_LOOP = 300
const DELTA = 0.01
const SEPARETOR = "$"

type FaceVector map[rune]float64

var charSet map[rune]int
var charHistogram map[rune]int
var chars []rune
var	IdToFace [][]rune

// 日本語に変換
var toJa = map[int]string{-1:"悲しい", 1:"楽しい"}

type LearningItem struct {
	input FaceVector
	answer int
}
var learningItems []LearningItem

// g(c)   = tf(c) * idf(c)
// tf(c)  = 文字 C が入力テキスト中に出現する回数
// idf(c) = log2(N / histogram(c))
// histogram(c)   = 訓練データ中に文字 C が出現する頻度
// N      = 学習データの総数
func MakeFaceVector(face []rune) FaceVector {
	vector := make(FaceVector)
	tf := make(map[rune]int)
	for _, char := range face {
		tf[char] += 1
	}
	for _, char := range face {
		histogram, ok := charHistogram[char];
		if !ok {continue}
		vector[char] = float64(tf[char]) *
			math.Log2(float64(len(IdToFace)) / float64(histogram))
	}
	return vector
}

func ShowVector(vector FaceVector, face []rune) {
	fmt.Println(string(face))
	for _, char := range face {
		fmt.Printf("%s -> %lf\n", string(char), vector[char])
	}
	fmt.Println()
}

func ShowWeight(weight FaceVector) {
	for char, w := range weight {
		fmt.Printf("%s -> %lf\n", string(char), w)
	}
	fmt.Println()
}

func ShowHistogram(histogram map[rune]int) {
	for char, count := range histogram {
		fmt.Printf("%s : %d\n", string(char), count)
	}
	fmt.Println()
}

func Add(x, y FaceVector) FaceVector {
	vector := make(FaceVector)
	for c, v := range x {vector[c] += v}
	for c, v := range y {vector[c] += v}
	return vector
}

// 内積
func InProduct(x, y FaceVector) float64 {
	sum := float64(0)
	for xC, xV := range x {
		for yC, yV := range y {
			if xC == yC {sum += xV * yV}
		}
	}
	return sum
}

func ScalarTimes(a float64, x FaceVector) FaceVector {
	for c, v := range x {
		x[c] = v * a
	}
	return x
}

func Sign(x float64) int {
	if x < 0 {return -1} else {return 1}
}

func EstimateWeight(learningItems []LearningItem) FaceVector {
	weight := make(FaceVector)
	for char, _ := range charHistogram {weight[char] = 0}
	for i := 0; i < MAX_LOOP; i++ {
		for j := 0; j < len(learningItems); j++ {
			item := learningItems[j]
			if item.answer != Sign(InProduct(item.input, weight)) {
				weight = ScalarTimes(float64(item.answer) * DELTA, Add(weight, item.input))
			}
			// ShowVector(ScalarTimes(item.answer, Add(weight, item.input)), IdToFace[j])
		}
	}
	return weight
}

func Init() {
	chars = make([]rune, 0)
	IdToFace = make([][]rune, 0)
	charSet = make(map[rune]int)
	charHistogram = make(map[rune]int)
	learningItems = make([]LearningItem, 0)
}

func Predict(face string, weight FaceVector) int {
	return Sign(InProduct(weight, MakeFaceVector([]rune(face))))
}

func main() {
	Init()
	ans := make([]int, 0)
	linenum := 0
	std.ReadFile("test/fun-sad-face.txt", func(line string) {
		linenum++
		words := strings.Split(line, SEPARETOR)
		face := words[0]
		if len(words) != 2 {
			panic(fmt.Sprintf("Informal Learning data: %d %v\n", linenum, words))}
		exists := make(map[rune]bool)
		for _, char := range []rune(face) {
			if _, ok := charHistogram[char]; ok && !exists[char] {
				charHistogram[char] += 1
				exists[char] = true
			}
			if _, ok := charSet[char]; !ok {
				chars = append(chars, char)
				charSet[char] = len(chars) - 1
				charHistogram[char] = 1
			}
		}
		IdToFace = append(IdToFace, []rune(face))
		x, err := strconv.Atoi(words[1])
		if err != nil {panic(fmt.Sprintf("can't convert string to int: %s\n", words[1]))}
		ans = append(ans, x)
	})
	for i, face := range IdToFace {
		vector := MakeFaceVector(face)
		learningItems = append(learningItems,
			LearningItem{input: vector, answer: ans[i]})
	  // ShowVector(vector, face)
	}
	weight := EstimateWeight(learningItems)
	
	// 学習の判定結果を見る
	std.ReadFile("test/kaomoji-250.txt", func(face string) {
		fmt.Printf("%s  ---> %s\n", string(face), toJa[Predict(face, weight)])
	})
}



