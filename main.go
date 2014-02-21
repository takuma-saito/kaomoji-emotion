
package main

import (
	"fmt"
	"math"
	"time"
	"strings"
	"./server"
	"core/std"
	"core/bit"
)

// MAX: 50,  DELTA: 1
const MAX_LOOP = 10
const MERGIN = 0.0
const SEPARETOR = '$'
var basetime time.Time

type FaceVector map[rune]float64
type Name []rune

// 分類集合の情報
type Class struct {	
	name map[uint64]int  // クラス集合（分類集合）の名前
	id []Name            // クラスID
}

// 学習データ
type LearningItems struct {
	histogram map[rune]int    // 文字の頻度分布
	faces []Name              // 顔文字集合
	faceVectors []FaceVector  // 顔文字の特徴ベクトル
	weights []FaceVector       // クラス学習用の重み関数
	answers []int             // 正解クラス
	class Class
}

func MakeHash(name []rune) uint64 {
	hash := uint64(1)
	for _, k := range name {hash = hash << 6 + uint64(k)}
	return hash
}

func ReadLines(filename string) []string {
	var lines []string
	std.ReadFile(filename, func(line string) {
		lines = append(lines, line)
	})
	return lines
}

func GetFaces(lines []string) (faces []Name, answers []int, class Class) {
	var ans int
	className := make(map[uint64]int)
	classID := []Name{}
	id := 0
	for linenum, line := range lines {
		words := strings.Split(line, string(SEPARETOR))
		if len(words) != 2 {
			panic(fmt.Sprintf("Informal Learning data: %d %v\n", linenum + 1, words))}
		faces = append(faces, Name(words[0]))
		h := MakeHash([]rune(words[1]))
		if c, ok := className[h]; ok {ans = c} else {
			className[h] = id
			classID = append(classID, Name(words[1]))
			ans = id
			id += 1
		}
		answers = append(answers, ans)
	}
	class = Class{name:className, id:classID}
	return
}

// Debug function

func ShowFaces(faces []Name) {
	for _, face := range faces {
		fmt.Println(string(face))
	}
}

func ShowVector(vector FaceVector, face Name) {
	fmt.Println(string(face))
	for _, char := range face {
		fmt.Printf("%s -> %3.1f\n", string(char), vector[char])
	}
	fmt.Println()
}

func ShowVectors(vectors []FaceVector, faces []Name) {
	for i, vector := range vectors {
		ShowVector(vector, faces[i])
	}
	fmt.Println()
}

func ShowWeight(weight FaceVector) {
	for char, w := range weight {
		fmt.Printf("%-5s %3.1f\n", string(char), w)
	}
	fmt.Printf("\n\n")
}

func ShowWeights(weights []FaceVector) {
	for _, weight := range weights {
		ShowWeight(weight)
	}
	fmt.Println()
}

func ShowHistogram(histogram map[rune]int) {
	for char, count := range histogram {
		fmt.Printf("%s : %d\n", string(char), count)
	}
	fmt.Println()
}

func  ShowClass(class Class) {
	for id, className := range class.id {
		fmt.Printf("%d : %s\n", id, string(className))
	}
	fmt.Println()
}

func (items *LearningItems) Show() {
	ShowHistogram(items.histogram)
	ShowFaces(items.faces)
	ShowVectors(items.faceVectors, items.faces)
	ShowWeights(items.weights)
	ShowClass(items.class)
}

func Add(x, y FaceVector) FaceVector {
	vector := make(FaceVector)
	for c, v := range x {vector[c] = v}
	for c, v := range y {vector[c] += v}
	return vector
}

func Sub(x, y FaceVector) FaceVector {
	vector := make(FaceVector)
	for c, v := range x {vector[c] = v}
	for c, v := range y {vector[c] -= v}
	return vector
}

// 内積
func InProduct(x, y FaceVector, chars map[rune]int) float64 {
	sum := float64(0)
	for c, _ := range chars {
		a, ok1 := x[c]
		b, ok2 := y[c]
		if ok1 && ok2 {sum += a * b}
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

func Sign(x float64) float64 {
	if x < MERGIN {return float64(-1)} else {return float64(1)}
}

func Partition(n int) int {
	return n / ((1 + bit.Log2(n)) * 3)
}

func ArgMax(weights []FaceVector, vector FaceVector, chars map[rune]int) int {
	max := math.Inf(-1); res := 0
	for i, weight := range weights {
		result := InProduct(weight, vector, chars)
		if result > max {
			res = i
			max = result
		}
	}
	return res
}

func (items *LearningItems) Predict(face Name) int {
	return ArgMax(items.weights, 
		MakeFaceVector(face, items.histogram, len(items.faces)), items.histogram)
}

func MakeFaceVector(face Name, histogram map[rune]int, N int) FaceVector {
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
	
	// バイアス項
	vector[rune(SEPARETOR)] = float64(1.0)
	
	return vector
}

func MakeLItems(faces []Name, answers []int, class Class) *LearningItems {
	items := &LearningItems{
		histogram: make(map[rune]int),
		faces: faces,
		faceVectors: []FaceVector{},
		class: class,
		answers: answers,		
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

// OK -> NG の時 true  -> 重みを増やす
// NG -> OK の時 false -> 重みを減らす
func TorF(r bool) float64 {
	if r {return float64(1)} else {return float64(-1)}
}

func (items *LearningItems) EstimateWeight() {
	items.weights = make([]FaceVector, len(items.class.id))
	for i, _ := range items.weights {
		items.weights[i] = make(FaceVector)
		for char, _ := range items.histogram {items.weights[i][char] = 0}
	}
	for i := 0; i < MAX_LOOP; i++ {
		for j := 0; j < len(items.faces); j++ {
			predicted := ArgMax(items.weights, items.faceVectors[j], items.histogram)
			if predicted == items.answers[j] {continue}
			// 間違えた場合に学習を行う
			items.weights[items.answers[j]] =
				Add(items.weights[items.answers[j]], items.faceVectors[j])
			items.weights[predicted] =
				Sub(items.weights[predicted], items.faceVectors[j])
		}
	}
	return
}

// メモリ割り当てを行う slice
func Slice(list []string, x, y int) (main []string, rest []string) {
	if x > y {panic(fmt.Sprintf("slice must be y > x: %d > %d\n", x, y))}
	main = append(make([]string, 0), list[x:y]...)
	rest = append(append(make([]string, 0), list[0:x]...), list[y:len(list)]...)
	return
}

// 交差検定
func CrossValidate(lines []string) {
	k := Partition(len(lines)) // 分割数
	trials := 0
	error := 0
	for i := 1; i <= len(lines) / k; i++ {
		main, rest := Slice(lines, (i - 1) * k, i * k)
		itemsL := MakeLItems(GetFaces(rest))
		itemsL.EstimateWeight()
		items := MakeLItems(GetFaces(main))
		for i, face := range items.faces {
			answer := string(items.class.id[items.answers[i]])
			predicted := string(itemsL.class.id[itemsL.Predict(face)])
			if answer != predicted {
				fmt.Printf("error: %-15s %s\n", string(face),
					string(predicted))
				error += 1
			}
			trials += 1
		}
	}
	fmt.Printf("Partition: %d\n", k)
	fmt.Printf("trials: %d\nerror:%d\nsuccess rate: %3.2f%%\n",
		trials, error, 100 * float64(trials - error) / float64(trials))
}

func Play() {
	items := MakeLItems(GetFaces(ReadLines(("test/category-938.txt"))))
	//items.Show()
	items.EstimateWeight()
	std.ReadFile("test/kaomoji-300.txt", func(face string) {
		fmt.Printf("%-15s %s\n",
			string(face),
			string(items.class.id[items.Predict(Name(face))]))
	})
}

func Test() {
	CrossValidate(ReadLines(("test/category-938.txt")))
}

func StartServer(port int) {
	items := MakeLItems(GetFaces(ReadLines(("test/category-938.txt"))))
	items.EstimateWeight()
	server.Start(port, func(face string) string {
		if len(face) == 0 {return face}
		return string(items.class.id[items.Predict(Name(face))])
	})
}

func main() {
	// Test()
	// StartServer(6666)
	Play()
}


