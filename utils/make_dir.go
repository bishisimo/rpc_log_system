/*
@author '彼时思默'
@time 2020/5/9 上午9:46
@describe:
*/
package utils

import "os"

func MakeDir(dir string) {
	if _, err := os.Stat(dir); !os.IsExist(err) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}
}
