package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

//-----------------コマンドラインオプション-----------------------------------

var (
	XrfOpt  = flag.String("Xrf", "default", "help message for \"s\" option")
	YrfOpt  = flag.String("Yrf", "default", "help message for \"s\" option")
	wfpOpt  = flag.String("wfp", "default", "help message for \"s\" option")
	distOpt = flag.String("d", "default", "help message for \"s\" option")
)

//-------------------------- エラー処理-----------------------------------
func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

//----------------スライスにおける実数・文字列変換-----------------------------------

func strToFloat64_1d(i []string) []float64 {
	// 1次元スライスに対して、文字列型から実数型に変換
	// 引数：文字列型の1次元スライス
	// 返り値：実数型の1次元スライス

	f := make([]float64, len(i))
	for n := range i {
		f[n], _ = strconv.ParseFloat(i[n], 64)
	}
	return f
}

func float64ToStr(x float64) string {
	// 実数型から文字列型に変換
	// 引数：実数型の1変数
	// 返り値：文字列型の変数

	f := strconv.FormatFloat(x, 'f', 8, 64)

	return f
}

//-------------------------- ファイル 入出力関係-----------------------------------

func csvReader_1d(readFilename string) []string {
	// csv読み込み
	// 引数：ファイル名(ファイルパス含む)
	// 返り値：csvファイルから読み込んだ1次元スライス

	// 読み込み
	fr, err := os.Open(readFilename)
	failOnError(err)

	defer fr.Close()

	r := csv.NewReader(fr)

	rows, err := r.ReadAll() // rowsのshapeは(len(rows),1)と2次元となってしまっている
	failOnError(err)

	// 1次元へflatten
	str_data := make([]string, len(rows))
	for i := 0; i < len(rows); i++ {
		str_data[i] = rows[i][0]
	}

	return str_data
}

func csvWriter(str string, writeFilename string) {
	// csv 書き込み
	// 引数：csvファイルに書き込む変数．ファイル名(ファイルパス含む)
	// 返り値：なし

	strCSV := str

	// 以下txtファイル出力の要領でcsvファイルを出力

	// ファイルを書き込み用にオープン (mode=0666)
	file, err := os.Create(writeFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// テキストを書き込む
	_, err = file.WriteString(strCSV)
	if err != nil {
		panic(err)
	}
}

//------------------------- 要素同士の距離------------------------

func dist(a float64, b float64, selectDist string) float64 {
	// 距離計算方法の選択
	// 引数：実数a, 実数b, 距離計算方法
	// 返り値：距離
	if selectDist == "se" {
		return squareError(a, b)
	} else {
		return squareError(a, b)
	}
}

func squareError(a float64, b float64) float64 {
	// 2乗誤差
	// 引数：実数a, 実数b
	// 返り値：2乗誤差
	return math.Pow(a-b, 2)
}

//--------------要素が3のスライスの中から最小値を探索------------------

func selectMin(List []float64) float64 {
	// 動的計画法を行う上ですでに計算された3値より最小の値を選択
	// 引数：3値が要素となるスライス
	// 返り値：最小値
	min := 99999.9
	for _, v := range List {
		if min > v {
			min = v
		}
	}
	return min
}

func dtw(x []float64, y []float64, selectDist string) ([][]float64, [][]float64) {
	// DTWを動的計画法で計算する
	// 引数：距離を計算する時系列x,y 距離計算方法
	// 返り値：動的計画法で計算された2重スライス, 全組み合わせの距離が計算された2重スライス

	// 要素が空の2重スライスを生成
	dtwArray := make([][]float64, len(x))
	for i := range dtwArray {
		dtwArray[i] = make([]float64, len(y))
	}
	distArray := make([][]float64, len(x))
	for i := range distArray {
		distArray[i] = make([]float64, len(y))
	}

	dtwArray[0][0] = dist(x[0], y[0], selectDist)
	distArray[0][0] = dist(x[0], y[0], selectDist)

	for i := 1; i < len(x); i++ {
		dtwArray[i][0] = dist(x[i], y[0], selectDist) + dtwArray[i-1][0]
		distArray[i][0] = dist(x[i], y[0], selectDist) + distArray[i-1][0]
	}

	for j := 1; j < len(y); j++ {
		dtwArray[0][j] = dist(x[0], y[j], selectDist) + dtwArray[0][j-1]
		distArray[0][j] = dist(x[0], y[j], selectDist) + distArray[0][j-1]
	}

	for i := 1; i < len(x); i++ {
		for j := 1; j < len(y); j++ {
			candMin := []float64{dtwArray[i-1][j], dtwArray[i][j-1], dtwArray[i-1][j-1]}
			dtwArray[i][j] = dist(x[i], y[j], selectDist) + selectMin(candMin)
			distArray[i][j] = dist(x[i], y[j], selectDist)
		}
	}

	return dtwArray, distArray
}

//------------------------------ Main -----------------------------------

func main() {

	// コマンドライン引数の確認
	flag.Parse()
	XreadFilename := *XrfOpt // 読み込みxファイル名
	YreadFilename := *YrfOpt // 読み込みyファイル名
	writeFilePath := *wfpOpt // 書き出しファイル名
	selectDist := *distOpt
	if XreadFilename == "default" || YreadFilename == "default" || selectDist == "default" {
		fmt.Print("コマンドライン引数エラー")
		return
	}

	// 配列Xのcsvファイル読み込み
	Xstr_data := csvReader_1d(XreadFilename)
	Ystr_data := csvReader_1d(YreadFilename)

	// 各要素を文字列型から実数値型へ変換
	x := strToFloat64_1d(Xstr_data)
	y := strToFloat64_1d(Ystr_data)

	// DTWを計算
	dtwArray, _ := dtw(x, y, selectDist)

	if writeFilePath == "default" {
		// DTW距離を標準出力する
		fmt.Print(dtwArray[len(x)-1][len(y)-1])
	} else {
		// 実数値型から文字列型へ変換
		dtwDistance := float64ToStr(dtwArray[len(x)-1][len(y)-1])
		// DTW距離をファイル書き出し
		csvWriter(dtwDistance, writeFilePath)
	}
}
