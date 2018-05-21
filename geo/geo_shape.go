/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     geo
 * @date        2018-04-27 20:25
 */
package geo

import (
	"fmt"
	"math"
)

//构造点
func MakeGeoPoint(lat, lng float64) GeoPoint {
	return GeoPoint{Lat: lat, Lng: lng}
}

//一个点
type GeoPoint struct {
	Lat float64 `json:"lat"` //纬度
	Lng float64 `json:"lng"` //经度
}

//返回字符串表示的形式
func (gp *GeoPoint) FormatStr() string {
	return fmt.Sprintf("%v,%v", gp.Lat, gp.Lng)
}

//返回数组表示的形式
func (gp *GeoPoint) FormatArray() [2]float64 {
	return [2]float64{gp.Lat, gp.Lng}
}

//根据指定的距离、角度构造另一个点
func (gp *GeoPoint) PointAtDistAndAngle(distance, angle float64) GeoPoint {
	return PointAtDistAndAngle(*gp, distance, angle)
}

//跟另一个点是否相等
func (gp *GeoPoint) IsEqual(p GeoPoint) bool {
	return gp.Lat == p.Lat && gp.Lng == p.Lng
}

//判断一个点的经纬度是否合法
func (gp *GeoPoint) Check() bool {
	if gp.Lng > MAX_LONGITUDE || gp.Lng < MIN_LONGITUDE || gp.Lat > MAX_LATITUDE || gp.Lat < MIN_LATITUDE {
		return false
	}
	return true
}

//构造直线
func MakeGeoLine(p1 GeoPoint, p2 GeoPoint) GeoLine {
	return GeoLine{Point1: p1, Point2: p2}
}

//一条直接
type GeoLine struct {
	Point1 GeoPoint `json:"point1"` //起点
	Point2 GeoPoint `json:"point2"` //终点
}

//直线的长度
func (gl *GeoLine) Length() float64 {
	return EarthDistance(gl.Point1, gl.Point2)
}

//获取直线的最小外包矩形，如果是条平行线或竖线的话，可能会有问题
func (gl *GeoLine) GetBoundsRect() GeoRectangle {
	return GeoRectangle{
		MaxLat: math.Max(gl.Point2.Lat, gl.Point1.Lat),
		MaxLng: math.Max(gl.Point2.Lng, gl.Point1.Lng),
		MinLat: math.Min(gl.Point2.Lat, gl.Point1.Lat),
		MinLng: math.Min(gl.Point2.Lng, gl.Point1.Lng),
	}
}

//是否包含某个点，基本思路：
//点为Q，线段为P1P2，判断点Q在线段上的依据是：(Q-P1)×(P2-P1)=0
//且Q在以P1P2为对角定点的矩形内
func (gl *GeoLine) IsContainPoint(p GeoPoint) bool {
	rect := gl.GetBoundsRect()
	if !rect.IsPointInRect(p) {
		return false
	}
	if p.IsEqual(gl.Point2) || p.IsEqual(gl.Point1) {
		return true
	}
	p1 := vectorDifference(gl.Point1, gl.Point2)
	p2 := vectorDifference(p, gl.Point1)
	cross := vectorCrossProduct(p1, p2)
	//浮点类型计算时候与0比较时候的容差
	var floatDiff = 2e-10
	if math.Abs(0-cross) < floatDiff {
		return true
	}
	return false
}

//临时方法
func vectorDifference(p1 GeoPoint, p2 GeoPoint) GeoPoint {
	return GeoPoint{Lat: p1.Lat - p2.Lat, Lng: p1.Lng - p2.Lng}
}

func vectorCrossProduct(p1 GeoPoint, p2 GeoPoint) float64 {
	return p1.Lat*p2.Lng - p1.Lng*p2.Lat
}

//与另一条直线是否相交【当一条直接的两个顶点相同时判断有点问题】
func (gl *GeoLine) IsIntersectWithLine(line GeoLine) bool {
	_, b := gl.GetIntersectPoint(line)
	return b
}

//求两直线交点的坐标
//参考：https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect
//返回值：交点、是否平行、是否相交
func (gl *GeoLine) GetIntersectPoint(line GeoLine) (GeoPoint, bool) {
	if gl.Length() == 0 {
		if line.IsContainPoint(gl.Point1) || line.IsContainPoint(gl.Point2) {
			return gl.Point1, true
		} else {
			return GeoPoint{}, false
		}
	}
	if line.Length() == 0 {
		if gl.IsContainPoint(line.Point1) || gl.IsContainPoint(line.Point2) {
			return line.Point1, true
		} else {
			return GeoPoint{}, false
		}
	}
	p := gl.Point1
	//一线段的向量差
	r := vectorDifference(gl.Point2, gl.Point1)
	q := line.Point1
	//另一线段的向量差
	s := vectorDifference(line.Point2, line.Point1)
	//两线段向量差的叉乘
	rCrossS := vectorCrossProduct(r, s)
	qMinusP := vectorDifference(q, p)
	if rCrossS == 0 {
		//If r × s = 0 and (q − p) × r = 0, then the two lines are collinear.
		//同一条直线，随便返回一个点即可
		if vectorCrossProduct(qMinusP, r) == 0 {
			return line.Point1, true
		} else {
			//If r × s = 0 and (q − p) × r ≠ 0,
			//then the two lines are parallel and non-intersecting.
			//两条为平行线，没有交点
			return GeoPoint{}, false
		}
	}
	//If r × s ≠ 0 and 0 ≤ t ≤ 1 and 0 ≤ u ≤ 1,
	//the two line segments meet at the point p + t r = q + u s
	t := vectorCrossProduct(qMinusP, s) / rCrossS
	u := vectorCrossProduct(qMinusP, r) / rCrossS
	//有相交点
	if t >= 0 && t <= 1 && u >= 0 && u <= 1 {
		p1 := GeoPoint{Lat: gl.Point1.Lat + t*r.Lat, Lng: gl.Point1.Lng + t*r.Lng}
		p2 := GeoPoint{Lat: line.Point1.Lat + u*s.Lat, Lng: line.Point1.Lng + u*s.Lng}
		interPoint := p1
		//如果在计算的时候有点小小的误差，这里直接取中间得了，理论上这两个点应该相等
		if !p1.IsEqual(p2) {
			interPoint = MidPoint(p1, p2)
		}
		return interPoint, true
	}
	return GeoPoint{}, false
}

//一个圆点
type GeoCircle struct {
	Center GeoPoint `json:"center"` //圆心
	Radius float64  `json:"radius"` //半径（米）
}

//一个点是否在圆内
func (gc *GeoCircle) InCircle(point GeoPoint) bool {
	dist := EarthDistance(gc.Center, point)
	return dist <= gc.Radius
}

//一个矩形
type GeoRectangle struct {
	MinLat float64 `json:"min_lat"` //最小纬度
	MinLng float64 `json:"min_lng"` //最小经度
	MaxLat float64 `json:"max_lat"` //最大纬度
	MaxLng float64 `json:"max_lng"` //最大经度
}

//经度方向的跨度
func (gr *GeoRectangle) LngSpan() float64 {
	return gr.MaxLng - gr.MinLng
}

//纬度方向的跨度
func (gr *GeoRectangle) LatSpan() float64 {
	return gr.MaxLat - gr.MinLat
}

//判断给定的经纬度是否在小格子内
func (gr *GeoRectangle) IsPointInRect(point GeoPoint) bool {
	if point.Lat <= gr.MaxLat && point.Lat >= gr.MinLat && point.Lng >= gr.MinLng && point.Lng <= gr.MaxLng {
		return true
	}
	return false
}

//获取中心点坐标
func (gr *GeoRectangle) MidPoint() GeoPoint {
	point := MidPoint(GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng}, GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng})
	return point
}

//矩形X方向的边长，即纬度线方向，保持纬度相同即可
func (gr *GeoRectangle) Width() float64 {
	return EarthDistance(
		GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng},
		GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng},
	)
}

//矩形Y方向的边长，即经度线方向，保持经度相同即可
func (gr *GeoRectangle) Height() float64 {
	return EarthDistance(
		GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng},
		GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng},
	)
}

//矩形的所有的边
func (gr *GeoRectangle) GetRectVertex() (ret []GeoPoint) {
	ret = append(ret,
		GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng},
		GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng},
		GeoPoint{Lat: gr.MaxLat, Lng: gr.MinLng},
		GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng},
	)
	return
}

//矩形的所有的边
func (gr *GeoRectangle) GetRectBorders() (ret []GeoLine) {
	p := gr.GetRectVertex()
	ret = append(ret,
		GeoLine{Point1: p[0], Point2: p[1]},
		GeoLine{Point1: p[1], Point2: p[2]},
		GeoLine{Point1: p[2], Point2: p[3]},
		GeoLine{Point1: p[3], Point2: p[0]},
	)
	return
}

//构造多边形
func MakeGeoPloygon(points []GeoPoint) GeoPloygon {
	return GeoPloygon{Points: points}
}

//一个多边形
type GeoPloygon struct {
	Points []GeoPoint `json:"points"` //一堆顶点，必须是首尾相连有序的
}

//获取所有的顶点
func (gp *GeoPloygon) GetPoints() []GeoPoint {
	return gp.Points
}

//多边形的所有的边
func (gp *GeoPloygon) GetPloygonBorders() (ret []GeoLine) {
	if !gp.Check() {
		return
	}
	points := gp.GetPoints()
	l := len(points)
	p0 := points[0]
	for i := 1; i < l; i++ {
		p := points[i]
		ret = append(ret, GeoLine{Point1: p0, Point2: p})
		p0 = p
	}
	ret = append(ret, GeoLine{Point1: points[l-1], Point2: points[0]})
	return
}

//添加点
func (gp *GeoPloygon) AddPoint(p GeoPoint) {
	gp.Points = append(gp.Points, p)
}

//判断是否是合法的多边形
func (gp *GeoPloygon) Check() bool {
	if len(gp.Points) < 3 {
		return false
	}
	return true
}

//获取最小外包矩形
func (gp *GeoPloygon) GetBoundsRect() GeoRectangle {
	var maxLat = MIN_LATITUDE
	var maxLng = MIN_LONGITUDE
	var minLat = MAX_LATITUDE
	var minLng = MAX_LONGITUDE
	for _, p := range gp.Points {
		maxLat = math.Max(maxLat, p.Lat)
		minLat = math.Min(minLat, p.Lat)
		maxLng = math.Max(maxLng, p.Lng)
		minLng = math.Min(minLng, p.Lng)
	}
	return GeoRectangle{MaxLat: maxLat, MaxLng: maxLng, MinLat: minLat, MinLng: minLng}
}

//判断点是否在多边形内部，此处使用最简单的射线法判断
//边数较多时性能不高，只适合在写入时小批量判断
//计算射线与多边形各边的交点，如果是偶数，则点在多边形外，否则在多边形内。
//还会考虑一些特殊情况，如点在多边形顶点上，点在多边形边上等特殊情况。
//参考：http://api.map.baidu.com/library/GeoUtils/1.2/src/GeoUtils.js
func (gp *GeoPloygon) IsPointInPolygon(p GeoPoint) bool {
	if !p.Check() || !gp.Check() {
		return false
	}
	//判断最小外包矩形
	rect := gp.GetBoundsRect()
	if !rect.IsPointInRect(p) {
		return false
	}

	//交点总数
	var interCount = 0
	//浮点类型计算时候与0比较时候的容差
	var floatDiff = 2e-10
	//相邻的两个顶点
	var p1, p2 GeoPoint
	//顶点个数
	PNum := len(gp.Points)

	//逐个顶点的判断
	p1 = gp.Points[0]
	points := gp.Points
	for i := 1; i < PNum; i++ {
		//正好落在了顶点上
		if p1.IsEqual(p) {
			return true
		}
		//其他顶点
		p2 = points[i%PNum]
		//射线没有交点
		if p.Lat < math.Min(p1.Lat, p2.Lat) || p.Lat > math.Max(p1.Lat, p2.Lat) {
			p1 = p2
			continue
		}
		//射线相交了
		if p.Lat > math.Min(p1.Lat, p2.Lat) && p.Lat < math.Max(p1.Lat, p2.Lat) {
			//东西向有交点
			if p.Lng <= math.Max(p1.Lng, p2.Lng) {
				//此边为一条横线
				if p1.Lat == p2.Lat && p.Lng >= math.Min(p1.Lng, p2.Lng) {
					return true
				}
				//一条竖线
				if p1.Lng == p2.Lng {
					if p1.Lng == p.Lng {
						return true
					} else {
						interCount++
					}
				} else {
					xInters := (p.Lat-p1.Lat)*(p2.Lng-p1.Lng)/(p2.Lat-p1.Lat) + p1.Lng
					if math.Abs(p.Lng-xInters) < floatDiff {
						return true
					}
					if p.Lng < xInters {
						interCount++
					}
				}
			}
		} else {
			if p.Lat == p2.Lat && p.Lng <= p2.Lng {
				p3 := points[(i+1)%PNum]
				if p.Lat >= math.Min(p1.Lat, p3.Lat) && p.Lat <= math.Max(p1.Lat, p3.Lat) {
					interCount++
				} else {
					interCount += 2
				}
			}
		}
		p1 = p2
	}

	if interCount%2 == 0 {
		return false
	} else {
		return true
	}
}
