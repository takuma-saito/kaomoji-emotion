
package main
import ("fmt"; "core/std"; "core/bit")

type FaceVector map[rune]int

var charSet map[rune]int
var charHistogram map[rune]int
var chars []rune
var	IdToFace [][]rune
var faceVectors []FaceVector

// 正解集合

// g(c)   = tf(c) * idf(c)
// tf(c)  = 文字 C が入力テキスト中に出現する回数
// idf(c) = log2(N / histogram(c))
// histogram(c)   = 訓練データ中に文字 C が出現する頻度
// N      = 学習データの総数
func MakeFaceVector(face []rune, histogram map[rune]int, N int) FaceVector {
	vector := make(FaceVector)
	tf := make(map[rune]int)
	for _, char := range face {
		tf[char] += 1
	}
	for _, char := range face {
		vector[char] = tf[char] * bit.Log2(N / histogram[char])
	}
	return vector
}

func ShowVector(vector FaceVector, face []rune) {
	fmt.Println(string(face))
	for _, char := range face {
		fmt.Printf("%s -> %d\n", string(char), vector[char])
	}
	fmt.Println()
}

func ShowHistogram(histogram map[rune]int) {
	for char, count := range histogram {
		fmt.Printf("%s : %d\n", string(char), count)
	}
	fmt.Println()
}

func EstimateWeight(faceVectors []FaceVector, ans [][]rune)

func main() {
	chars = make([]rune, 0)
	IdToFace = make([][]rune, 0)
	charSet = make(map[rune]int)
	charHistogram = make(map[rune]int)
	std.ReadFile("test/kaomoji-200.txt", func(face string) {
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
	})
	// ShowHistogram(charHistogram)
	for _, face := range IdToFace {
		vector := MakeFaceVector(face, charHistogram, len(IdToFace))
		faceVectors = append(faceVectors, vector)
		// ShowVector(vector, face)
	}
	weight := EstimateWeight(faceVectors, ans)
}


