/*
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     helper
 * @date        2018-01-25 19:19
 */
package helper

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"unsafe"
	"encoding/binary"
	"math"
)

const N = int(unsafe.Sizeof(0))

func IsBigEndian() bool {
	x := 0x1234
	p := unsafe.Pointer(&x)
	p2 := (*[N]byte)(p)
	if p2[0] == 0 {
		return true
	}
	return false
}

//文件是否存在
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//是否全部为数字
func IsAllNumber(str string) bool {
	ret := strings.Trim(str, "0123456789")
	return len(ret) == 0
}

//字符串直接转为json
func StrToJSON(str string) string {
	var jsons bytes.Buffer
	for _, rn := range str {
		rint := int(rn)
		if rint < 128 {
			jsons.WriteRune(rn)
		} else {
			jsons.WriteString("\\u")
			jsons.WriteString(strconv.FormatInt(int64(rint), 16))
		}
	}
	return jsons.String()
}


//int64转byte
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

//bytes转int64
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

//float64转byte
func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)
	return bytes
}

//bytes转float
func ByteToFloat64(bytes []byte) float64 {
	bits := binary.BigEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

//字符串转为字节切片
func StrToByte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//字节切片转为字符串
func ByteToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
