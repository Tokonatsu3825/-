package contentssecurity

import (
	"fmt"
	"encoding/csv"
	"os"
	"strconv"
	"gonum.org/v1/gonum/mat"
	"errors"

	"time"
	"math/rand"
	
	conn "github.com/uecconsecexp/secexp2022/se_go/connector"
)

// ========== kadai1用 ==========

// main.goで使用したい機能はパッケージとして切り出しておくと便利です。
// lib.goに記述する必要はなく、別なファイルにしてもかまいません。
// ただしファイルの行頭はpackage contentssecurityとする必要があります。

// main.goで使用したい関数はパブリックとするため、大文字で始めます。
func Hello() string {
	return "Hello, world!"
}

func ReadData(csvPath string) ([][]float64, error) {

	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	str_table, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	table := make([][]float64, len(str_table)-1)
	
	for i, row := range str_table[1:] {
		table[i] = make([]float64, len(row)-1)
		for j, cell := range row[1:] {
			table[i][j], err = strconv.ParseFloat(cell, 64)
			if err != nil {
				return nil, err
			}
		}	
	}	
		
	return table, nil
}

func Calc_matrix(a [][]float64, b [][]float64)([][]float64, error) {
	rowsA := len(a)
	colsA := len(a[0])
	colsB := len(b[0])

	result := make([][]float64, rowsA)
	for i := range result {
		result[i] = make([]float64, colsB)
	}

	for i := 0; i < rowsA; i++ {
		for j := 0; j < colsB; j++ {
			for k := 0; k < colsA; k++ {
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}

	return result, nil
}

func Sum_matrix(a [][]float64, b [][]float64) ([][]float64, error) {

	rowsA, colsA := len(a), len(a[0])
	rowsB, colsB := len(b), len(b[0])

	if rowsA != rowsB || colsA != colsB {
		return nil, errors.New("行列のサイズの不一致")
	}

	result := make([][]float64, rowsA)
	for i := range result {
		result[i] = make([]float64, colsA)
	}

	for i := 0; i < rowsA; i++ {
		for j := 0; j < colsA; j++ {
			result[i][j] = a[i][j] + b[i][j]
		}
	}

	return result, nil
}

func Hantei_matrix(seiseki [][]float64, saitei [][]float64) ([][]float64, error) {

	rowsA, colsA := len(seiseki), len(seiseki[0])
	rowsB, colsB := len(saitei), len(saitei[0])

	if rowsA != rowsB || colsA != colsB {
		return nil, errors.New("行列のサイズの不一致")
	}

	result := make([][]float64, rowsA)
	for i := range result {
		result[i] = make([]float64, colsA)
	}

	for i := 0; i < rowsA; i++ {
		for j := 0; j < colsA; j++ {
			if seiseki[i][j] < saitei[i][j] {
				result[i][j] = 0
			} else {
				result[i][j] = 1
			}
		}
	}

	return result, nil
}


// 便利なモジュールやその他必要な情報は README.md にまとめています。

// ========== kadai2用 ==========

// 実際に通信を行う機能もパッケージ側に書いてしまいましょう。
// YobikouSide、ChugakuSideの中身を書き換えてください。

func YobikouSide() {
	yobikou, err := conn.NewYobikouServer()
	if err != nil {
		panic(err)
	}
	defer yobikou.Close() 

	//2
	// Mの受信
	matrixM, err := yobikou.ReceiveTable()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Received M: %v\n", matrixM)


	//3
	// 二次元配列を行列に変換
	matA := mat.NewDense(len(matrixM), len(matrixM[0]), nil)
	for i := 0; i < len(matrixM); i++ {
		for j := 0; j < len(matrixM[0]); j++ {
			matA.Set(i, j, matrixM[i][j])
		}
	}
	// 逆行列の計算
	var matAInv mat.Dense
	err = matAInv.Inverse(matA)
	if err != nil {
		panic(err)
	}

	//4
	topMat := mat.NewDense(3, 6, nil)
	for i := 0; i < 3; i++ {
		for j := 0; j < 6; j++ {
			topMat.Set(i, j, matAInv.At(i, j))
		}
	}
	bottomMat := mat.NewDense(3, 6, nil)
	for i := 0; i < 3; i++ {
		for j := 0; j < 6; j++ {
			bottomMat.Set(i, j, matAInv.At(i+3, j))
		}
	}
	// topMat を [][]float64 に変換
	topMatData := topMat.RawMatrix().Data
	topMatRows, topMatCols := topMat.Dims()
	topMatSlice := make([][]float64, topMatRows)
	for i := 0; i < topMatRows; i++ {
		topMatSlice[i] = topMatData[i*topMatCols : (i+1)*topMatCols]
	}
	// bottomMat を [][]float64 に変換
	bottomMatData := bottomMat.RawMatrix().Data
	bottomMatRows, bottomMatCols := bottomMat.Dims()
	bottomMatSlice := make([][]float64, bottomMatRows)
	for i := 0; i < bottomMatRows; i++ {
		bottomMatSlice[i] = bottomMatData[i*bottomMatCols : (i+1)*bottomMatCols]
	}

	//5
	// matrixBDash の計算
	omomi, err := ReadData("./omomi.txt")
	if err != nil {
		panic(err)
	}
	matrixBDash, err := Calc_matrix(bottomMatSlice, omomi)

	//8
	//A'を受け取る
	matrixADash, err := yobikou.ReceiveTable()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received A': %v\n", matrixADash)

	//9
	//B'を送信
	err = yobikou.SendTable(matrixBDash)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Send B': %v\n", matrixBDash)


	//11
	//matrixBDashDashの計算
	matrixX, err :=  Calc_matrix(matrixADash, topMatSlice)
	if err != nil {
		panic(err)
	}
	matrixBDashDash, err := Calc_matrix(matrixX, omomi)
	if err != nil {
		panic(err)
	}

	//12
	//A''を受け取る
	matrixADashDash, err := yobikou.ReceiveTable()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received A'': %v\n", matrixADashDash)
	fmt.Printf(" B'': %v\n", matrixBDashDash)

	//13 
	tekisei, err := Sum_matrix(matrixADashDash, matrixBDashDash)
	if err != nil {
		panic(err)
	}

	fmt.Printf("13 done\n")

	//14
	saiteiten, err := ReadData("./saiteiten.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("saiteiten: %v\n", saiteiten)

	saiteiten4 := make([][]float64, 4)
	for i := range saiteiten4 {
		saiteiten4[i] = make([]float64, 4)
		for j := 0; j < 4; j++ {
			saiteiten4[i][j] = saiteiten[0][j] 
		}
	}

	fmt.Printf("tekisei: %v\n", tekisei)
	fmt.Printf("saiteiten4: %v\n", saiteiten4)


	gouhi, err := Hantei_matrix(tekisei, saiteiten4)
	if err != nil {
		panic(err)
	}

	fmt.Printf("14 done\n")

	//15
	err = yobikou.SendTable(gouhi)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Send gouhi: %v\n", gouhi)
}



/////////////////////////
////////////////////////

// CSVを二次元配列に保存する関数
func ReadCSV(filename string) ([][]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var lines [][]float64
	isFirstRow := true
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}
		if isFirstRow {
			isFirstRow = false
			continue
		}
		var line []float64
		for i, item := range record {
			if i == 0 {
				continue
			}
			f, err := strconv.ParseFloat(item, 64)
			if err != nil {
				return nil, err
			}
			line = append(line, f)
		}
		lines = append(lines, line)
	}
	return lines, nil
}


// 行列を左右半分にする関数
func Splitmatrix(matrix [][]float64) ([][]float64, [][]float64) {
    // MleftとMrightを適切なサイズで初期化します。
    Mleft := make([][]float64, 6)
    Mright := make([][]float64, 6)
    for i := 0; i < 6; i++ {
        Mleft[i] = make([]float64, 3)
        Mright[i] = make([]float64, 3)
    }

    for i := 0; i < 6; i++ {
        for j := 0; j < 6; j++ {
            if j < 3 {
                Mleft[i][j] = matrix[i][j]
            } else {
                Mright[i][j-3] = matrix[i][j]
            }
        }
    }
    return Mleft, Mright
}


// 正則であるかを判定
func IsRegular(matrix [][]float64) bool {
	// 2次元スライスを*mat.Dense型に変換します。
	r, c := len(matrix), len(matrix[0])
	data := make([]float64, r*c)
	for i := range matrix {
		for j := range matrix[i] {
			data[i*c+j] = matrix[i][j]
		}
	}
	m := mat.NewDense(r, c, data)

	// 行列式を計算します。
	determinant := mat.Det(m)
	if determinant != 0 {
		return true
	}
	return false
}

func Gouhi_henkan(a [][]float64)([][]string, error){
	rows := 5
	cols := 5



	result := make([][]string, rows)
	for i := range result {
		result[i] = make([]string, cols)
	}

	result[0][0] = "　　 "
	result[1][0] = "生徒1"
	result[2][0] = "生徒2"
	result[3][0] = "生徒3"
	result[4][0] = "生徒4"

	result[0][1] = " 高校A"
	result[0][2] = " 高校B"
	result[0][3] = " 高校C"
	result[0][4] = " 高校D"

    for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if a[i][j] == 1 {
				result[i+1][j+1] = "　合格"
			} else {
				result[i+1][j+1] = "不合格"
			}

		}
	}

	return result, nil
}





// 中学校　プログラム
func Generaterandommatrix() [][]float64 {
	rand.Seed(time.Now().UnixNano()) // 乱数のシードを設定
	matrix := make([][]float64, 6)   // 6x6の行列を初期化
   
	for i := range matrix {
	 matrix[i] = make([]float64, 6)
	 for j := range matrix[i] {
	  matrix[i][j] = rand.Float64() * 100 // 0から100までの乱数を生成
	 }
	}
   
	return matrix
}

func ChugakuSide(addr string) {
	chugaku, err := conn.NewChugakuClient(addr)
	if err != nil {
		panic(err)
	}
	defer chugaku.Close() // 消さないでください。

	//1
	// 乱数行列を生成
	randommatrix := Generaterandommatrix()

	for !IsRegular(randommatrix){ // 正則であるかを確認
		randommatrix = Generaterandommatrix()
	}

	// //テスト用乱数行列
	// randommatrix := [][]float64{
	// 	{1.0, 2.0, 3.0, 4.0, 5.0, 6.0},
	// 	{0.0, 1.0, 0.0, 0.0, 0.0, 0.0},
	// 	{0.0, 0.0, 1.0, 0.0, 0.0, 0.0},
	// 	{0.0, 0.0, 0.0, 1.0, 0.0, 0.0},
	// 	{0.0, 0.0, 0.0, 0.0, 1.0, 0.0},
	// 	{0.0, 0.0, 0.0, 0.0, 0.0, 1.0},
	// }

	// 2
	// 乱数行列を送る
	err = chugaku.SendTable(randommatrix)
	if err != nil {
		panic(err)
	}

	// 6
	// Mleft, Mrightを求める
	Mleft, Mright := Splitmatrix(randommatrix)
	
	fmt.Printf("Mleft: %v\n", Mleft)
	fmt.Printf("Mright: %v\n", Mright)

	// 成績データ読み込み
	seiseki_data, err := ReadCSV("seiseki.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("成績: %v\n", seiseki_data)
	
	// 7
	// A'を求める
	A_prime, err := Calc_matrix(seiseki_data, Mleft)
	fmt.Printf("A': %v\n", A_prime)

	// 8
	// A'を送る
	err = chugaku.SendTable(A_prime)
	if err != nil {
		panic(err)
	}

	// ↑ここまで動作確認済み

	
	// B' を受信
	b_prime, err := chugaku.ReceiveTable()
	if err != nil {
		panic(err)
	}
	fmt.Printf("中学Received2: %v\n", b_prime)

	// 10
	// A''を求める
	A_dprime, err := Calc_matrix(seiseki_data, Mright)
	A_dprime, err = Calc_matrix(A_dprime, b_prime)

	// 12
	// A''を送る
	err = chugaku.SendTable(A_dprime)
	if err != nil {
		panic(err)
	}

	// 合否行列を受信
	gouhi_data, err := chugaku.ReceiveTable()
	if err != nil {
		panic(err)
	}
	gouhi, err := Gouhi_henkan(gouhi_data)

	// 16
    // 合否表を出力
	for i := 0; i < 5; i++{
		for j := 0; j < 5; j++{
			fmt.Print(gouhi[i][j])
			fmt.Print(" ")

		}
		fmt.Print("\n")
	}	

}
