
package main
import (
	"fmt"
	"math"
	"strings"
	"core/std"
	"core/bit"
)

// MAX: 50,  DELTA: 1
const MAX_LOOP = 10
const DELTA = 1
const SEPARETOR = "$"

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

func GetFacesFromFile(filename string) (faces []Name, answers []int, class Class) {
	var ans int
	className := make(map[uint64]int)
	classID := []Name{}
	linenum := 0
	id := 0
	std.ReadFile(filename, func(line string) {
		words := strings.Split(line, SEPARETOR)
		linenum += 1
		if len(words) != 2 {
			panic(fmt.Sprintf("Informal Learning data: %d %v\n", linenum, words))}
		faces = append(faces, Name(words[0]))
		h := MakeHash([]rune(words[1]))
		if c, ok := className[h]; ok {ans = c} else {
			className[h] = id
			classID = append(classID, Name(words[1]))
			ans = id
			id += 1
		}
		answers = append(answers, ans)
	})
	class = Class{name:className, id:classID}
	return
}

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

func Sign(x float64) float64 {
	if x < 0 {return float64(-1)} else {return float64(1)}
}

func Partition(n int) int {
	return n / ((1 + bit.Log2(n)) * 3)
}

func ArgMax(weights []FaceVector, vector FaceVector) int {
	max := math.Inf(-1); res := 0
	for i, weight := range weights {
		result := InProduct(weight, vector)
		if result > max {
			res = i
			max = result
		}
	}
	return res
}

func (items *LearningItems) Predict(face Name) int {
	return ArgMax(items.weights,
		MakeFaceVector(face, items.histogram, len(items.faces)))
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

func (items *LearningItems) EstimateWeight(class int) FaceVector {
	weight := make(FaceVector)
	for char, _ := range items.histogram {weight[char] = 0}
	for i := 0; i < MAX_LOOP; i++ {
		for j := 0; j < len(items.faces); j++ {
			predicted := Sign(InProduct(weight, items.faceVectors[j]))
			answer := TorF(items.answers[j] == class)
			// 間違えた場合に学習を行う *ここがおかしい*
			if predicted != answer  {
				weight = Add(weight,
					ScalarTimes(answer * DELTA, items.faceVectors[j]))
			}
		}
	}
	return weight
}

func (items *LearningItems) EstimateWeights() {
	items.weights = make([]FaceVector, len(items.class.id))
	for i, _ := range items.class.id {
		items.weights[i] = items.EstimateWeight(i)
	}
}

type Tuple struct {
	face Name
	answer int
}

// メモリ割り当てを行う slice
func Slice(list []Tuple, x, y int) (main []Tuple, rest []Tuple) {
	before := make([]Tuple, x)
	after := make([]Tuple, len(list) - y)
	main = make([]Tuple, y - x)
	copy(before, list[0:x])
	copy(main, list[x:y])
	copy(after, list[y:len(list)])
	rest = append(before, after...)
	return
}

func TranposeTuple(tuples []Tuple) (faces []Name, answers []int) {
	for _, tuple := range tuples {
		faces = append(faces, tuple.face)
		answers = append(answers, tuple.answer)
	}
	return
}

// 交差検定
func CrossValidate(faces []Name, answers []int, class Class) {
	k := Partition(len(faces)) // 分割数
	trials := 0
	error := 0
	tuples := make([]Tuple, len(faces))
	for i, _ := range tuples {tuples[i] = Tuple{face:faces[i], answer:answers[i]}}
	for i := 1; i <= len(tuples) / k; i++ {
		main, rest := Slice(tuples, (i - 1) * k, i * k)
		faces, answers := TranposeTuple(rest)
		items := MakeLItems(faces,  answers, class) // *bug* class が適切ではない
		items.EstimateWeights()
		for _, tuple := range main {
			if (tuple.answer != items.Predict(tuple.face)) {
				fmt.Printf("error: %-15s %s\n",
					string(tuple.face),
					string(items.class.id[items.Predict(tuple.face)]))
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
	items := MakeLItems(GetFacesFromFile("test/category.txt"))
	items.EstimateWeights()
	//ShowWeights(items.weights)
	std.ReadFile("test/kaomoji-250.txt", func(face string) {
		fmt.Printf("%-15s %s\n",
			string(face),
			string(items.class.id[items.Predict(Name(face))]))
	})
}

func Test() {
	CrossValidate(GetFacesFromFile("test/fun-sad-face.txt"))
}

func main() {
	// Test()
	Play()
}

