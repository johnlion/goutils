/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     geo
 * @date        2018-04-19 14:31
 */
package geo

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	//地球半径
	EARTH_RADIUS = 6378137
)

//格式化距离，最终都要输出以米为单位
// 输入"5.5km"=>"5500"
// 输入"5000m"=>"5000"
func FormatDistance(distance string) (ret float64) {
	ret = 0
	distance = strings.ToLower(distance)
	reg := regexp.MustCompile(`^[\d][\d\.]+km$`)
	if reg.MatchString(distance) { //以km为单位
		distance = strings.Replace(distance, "km", "", -1)
		dist, err := strconv.ParseFloat(distance, 64)
		if err != nil {
			return
		}
		ret = dist * 1000
		return
	}
	reg = regexp.MustCompile(`^[\d][\d\.]+m$`)
	if reg.MatchString(distance) { //以m为单位
		distance = strings.Replace(distance, "m", "", -1)
		dist, err := strconv.ParseFloat(distance, 64)
		if err != nil {
			return
		}
		ret = dist
		return
	}
	return
}

//计算两个经纬度间的中间位置
func MidPoint(point1, point2 GeoPoint) GeoPoint {
	if point2.IsEqual(point1) {
		return point2
	}
	lat1Arc := point1.Lat * math.Pi / 180.0
	lat2Arc := point2.Lat * math.Pi / 180.0
	lng1Arc := point1.Lng * math.Pi / 180.0
	diffLng := (point2.Lng - point1.Lng) * math.Pi / 180.0

	bx := math.Cos(lat2Arc) * math.Cos(diffLng)
	by := math.Cos(lat2Arc) * math.Sin(diffLng)

	lat3Rad := math.Atan2(math.Sin(lat1Arc)+math.Sin(lat2Arc), math.Sqrt(math.Pow(math.Cos(lat1Arc)+bx, 2)+math.Pow(by, 2)))
	lng3Rad := lng1Arc + math.Atan2(by, math.Cos(lat1Arc)+bx)

	lat3 := lat3Rad * 180.0 / math.Pi
	lng3 := lng3Rad * 180.0 / math.Pi

	return GeoPoint{Lat: lat3, Lng: lng3}
}

//在指定距离、角度上，返回另一个经纬度坐标
//lat、lng：源经纬度
//dist：距离，单位米
//angle：角度，如"45"
func PointAtDistAndAngle(point GeoPoint, dist, angle float64) GeoPoint {
	if dist <= 0 {
		return point
	}
	dr := dist / EARTH_RADIUS
	angle = angle * (math.Pi / 180.0)
	lat1 := point.Lat * (math.Pi / 180.0)
	lng1 := point.Lng * (math.Pi / 180.0)

	lat2 := math.Asin(math.Sin(lat1)*math.Cos(dr) + math.Cos(lat1)*math.Sin(dr)*math.Cos(angle))
	lng2 := lng1 + math.Atan2(math.Sin(angle)*math.Sin(dr)*math.Cos(lat1), math.Cos(dr)-(math.Sin(lat1)*math.Sin(lat2)))
	lng2 = math.Mod(lng2+3*math.Pi, 2*math.Pi) - math.Pi

	lat2 = lat2 * (180.0 / math.Pi)
	lng2 = lng2 * (180.0 / math.Pi)
	return GeoPoint{Lat: lat2, Lng: lng2}
}

//计算地球上的曲线距离，返回值为米
func EarthDistance(point1, point2 GeoPoint) float64 {
	if point1.IsEqual(point2) {
		return 0
	}
	rad := math.Pi / 180.0
	lat1 := point1.Lat * rad
	lng1 := point1.Lng * rad
	lat2 := point2.Lat * rad
	lng2 := point2.Lng * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * float64(EARTH_RADIUS)
}

/**
随机生成指定数量的多边形，必须指定基本的矩形区域、顶点的最大个数及最小个数
所有的这些多边形的边，非相邻的不能有交点！
仅于校验之用，不保证性能！
*/
func GenPolygons(baseRect GeoRectangle, polygonNum, pointMinNum, pointMaxNum int) (ret []GeoPolygon) {
	width := baseRect.Width()
	height := baseRect.Height()
	if width <= 0 || height <= 0 {
		return
	}
	if polygonNum <= 0 {
		return
	}
	if pointMaxNum < pointMinNum {
		return
	}
	if pointMinNum < 3 {
		return
	}
	fmt.Fprintf(os.Stdout, "start GenPolygons：baseRect[%d x %d] polygonNum[%d] pointMinNum[%d] pointMaxNum[%d]\n", int(width), int(height), polygonNum, pointMinNum, pointMaxNum)
	diffNum := pointMaxNum - pointMinNum + 1
	var randFloat float64
	diffLat := baseRect.MaxLat - baseRect.MinLat
	diffLng := baseRect.MaxLng - baseRect.MinLng
	for i := 0; i < polygonNum; i++ {
		rand.Seed(time.Now().UnixNano() + int64(i))
		//顶点的数量
		vertexNum := rand.Intn(diffNum) + pointMinNum
		var points []GeoPoint
		for vn := 0; vn < vertexNum; vn++ {
			//确保一个多边形的边都不相交
			for {
				rand.Seed(time.Now().UnixNano() + int64(vn))
				randFloat = rand.Float64()
				lat := baseRect.MinLat + randFloat*diffLat
				rand.Seed(time.Now().UnixNano() + int64(vn+1))
				randFloat = rand.Float64()
				lng := baseRect.MinLng + randFloat*diffLng
				point := MakeGeoPoint(lat, lng)
				points = append(points, point)
				polygon := MakeGeoPolygon(points)
				if polygon.IsBorderInterect() {
					points = points[0 : len(points)-1]
				} else {
					break
				}
			}
		}
		ret = append(ret, MakeGeoPolygon(points))
	}
	return
}
