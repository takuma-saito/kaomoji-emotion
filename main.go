
package main
import (
	"fmt"
	"math"
	"strings"
	"strconv"
	"core/std"
	"core/bit"
)

const MAX_LOOP = 500
const DELTA = 1.0
const SEPARETOR = "$"

type FaceVector map[rune]float64
type Face []rune

type LearningItems struct {
	histogram map[rune]int    // 文字の頻度分布
	faces []Face              // 顔文字集合
	faceVectors []FaceVector  // 顔文字の特徴ベクトル
	ans []int
	weight FaceVector
}

// 日本語に変換
var toJa = map[int]string{-1:"悲しい", 1:"楽しい"}

func GetFacesFromFile(filename string) (faces []Face, ans[]int) {
	linenum := 0
	std.ReadFile(filename, func(line string) {
		words := strings.Split(line, SEPARETOR)
		linenum += 1
		if len(words) != 2 {
			panic(fmt.Sprintf("Informal Learning data: %d %v\n", linenum, words))}
		x, err := strconv.Atoi(words[1])
		if err != nil {
			panic(fmt.Sprintf("can't convert string to int: %s\n", words[1]))}
		faces = append(faces, Face(words[0]))
		ans = append(ans, x)
	})
	return
}

func ShowFaces(faces []Face) {
	for _, face := range faces {
		fmt.Println(string(face))
	}
}

func ShowVector(vector FaceVector, face Face) {
	fmt.Println(string(face))
	for _, char := range face {
		fmt.Printf("%s -> %3.1f\n", string(char), vector[char])
	}
	fmt.Println()
}

func ShowVectors(vectors []FaceVector, faces []Face) {
	for i, vector := range vectors {
		ShowVector(vector, faces[i])
	}
	fmt.Println()
}

func ShowWeight(weight FaceVector) {
	for char, w := range weight {
		fmt.Printf("%s -> %3.1f\n", string(char), w)
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

func ScalarTimes(a float64, vector FaceVector) FaceVector {
	res := make(FaceVector)
	for key, value := range vector {
		res[key] = a * value
	}
	return res
}

func Sign(x float64) int {
	if x < 0 {return -1} else {return 1}
}

func Partition(n int) int {
	return (1 + bit.Log2(n)) * 4
}

func (items *LearningItems) Predict(face string) int {
	return Sign(InProduct(items.weight,
		MakeFaceVector([]rune(face), items.histogram, len(items.faces))))
}

func MakeFaceVector(face []rune, histogram map[rune]int, N int) FaceVector {
	vector := make(FaceVector)
	tf := make(map[rune]int)
	for _, char := range face {
		tf[char] += 1
	}
	for _, char := range face {
		count, ok := histogram[char];
		if !ok {continue}
		vector[char] = float64(tf[char]) * math.Log2(float64(N) / float64(count))
	}
	return vector
}

func MakeLItems(faces []Face, ans []int) *LearningItems {
	items := &LearningItems{
		histogram: make(map[rune]int),
		faces: faces,
		faceVectors: []FaceVector{},
		ans: ans,
	}
	
	// histogram の構成
	for _, face := range faces {
		exists := make(map[rune]bool)
		for _, char := range []rune(face) {
			if _, ok := items.histogram[char]; !ok {
				items.histogram[char] = 1
			} else {
				if !exists[char] {
					items.histogram[char] += 1
					exists[char] = true
				}
			}
		}
	}
	
	for _, face := range items.faces {
		items.faceVectors = append(items.faceVectors,
			MakeFaceVector(face, items.histogram, len(items.faces)))
	}
	return items
}

func (items *LearningItems) EstimateWeight() {
	items.weight = make(FaceVector)
	for char, _ := range items.histogram {items.weight[char] = 0}
	for i := 0; i < MAX_LOOP; i++ {
		for j := 0; j < len(items.faces); j++ {
			if items.ans[j] != Sign(InProduct(items.faceVectors[j], items.weight)) {
				items.weight = Add(items.weight,
					ScalarTimes(float64(items.ans[j]) * DELTA, items.faceVectors[j]))
			}
		}
	}
}

// 交差検定
func CrossValidate(faces []Face, ans []int) {
	k := Partition(len(faces)) // 分割数
	for i := 1; i <= len(faces) / k; i++ {
		// data := 
	}
}

func main() {
	faces, ans := GetFacesFromFile("test/fun-sad-face.txt")
	items := MakeLItems(faces, ans)
	items.EstimateWeight()
	// ShowWeight(items.weight)
	// ShowVectors(items.faceVectors, items.faces)
	std.ReadFile("test/kaomoji-250.txt", func(face string) {
		fmt.Printf("%-20s --> %s\n", string(face), toJa[items.Predict(face)])
	})
}

