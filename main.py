#coding:utf-8
#!/usr/bin/env python

import os
import time
import subprocess
import numpy as np


# 一時保存用のディレクトリ
tmpDir = "./tmp/"
os.makedirs(tmpDir, exist_ok=True)


def dtw1d(x,y):
	np.savetxt(tmpDir + "x.csv", x, delimiter=",")
	np.savetxt(tmpDir + "y.csv", y, delimiter=",")
	dtwGo = subprocess.run(["./dtw1d", "-Xrf", tmpDir + "x.csv", "-Yrf", tmpDir + "y.csv", "-wfp", tmpDir + "dtw.csv", "-d", "se"], check=True, stdout=subprocess.PIPE).stdout
	if dtwGo.decode() == "コマンドライン引数エラー":
		print("実行ファイルへの引数エラー")
	else:
		return float(dtwGo.decode())
		# return np.loadtxt(tmpDir+"dtw.csv", delimiter=",")


def dtw2d(X,Y):
	np.savetxt(tmpDir + "X.csv", X, delimiter=",")
	np.savetxt(tmpDir + "Y.csv", Y, delimiter=",")
	dtwGo = subprocess.run(["./dtw2d", "-Xrf", tmpDir + "X.csv", "-Yrf", tmpDir + "Y.csv", "-wfp", tmpDir + "dtw.csv", "-d", "se"], check=True, stdout=subprocess.PIPE).stdout
	if dtwGo.decode() == "コマンドライン引数エラー":
		print("実行ファイルへの引数エラー")
	else:
		return [float(i) for i in dtwGo.decode().replace("[","").replace("]","").split()]
		# return np.loadtxt(tmpDir+"dtw.csv", delimiter=",")


def main():

	NumOfData = 50 # データの数

	#乱数を用いて配列Xを生成
	X = np.random.rand(NumOfData,500) # 500次元ベクトルのデータ
	#乱数を用いて配列Yを生成
	Y= np.random.rand(NumOfData,300) # 300次元ベクトルのデータ

	# 1次元配列ずつ実行ファイルに渡し計算する
	dtw1 = []
	start = time.time()
	for x,y in zip(X,Y):
		dtw1.append(dtw1d(x,y))
	print("time of loop by dtw1d", time.time() - start)
	# print("result of dtw1d", np.array(dtw1))

	# 2次元配列全体を実行ファイル全体に渡し計算する
	start = time.time()
	dtw2 = dtw2d(X,Y)
	print("time of loop by dtw2d", time.time() - start)
	# print("result of dtw2d", np.array(dtw2))


if __name__ == "__main__":
    main()
