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

const (
	//浮点类型计算时候与0比较时候的容差
	FLOAT_DIFF = 2e-10
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

//直线的长度
func (gl *GeoLine) FormatStr() string {
	return fmt.Sprintf("%s-%s", gl.Point1.FormatStr(), gl.Point2.FormatStr())
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
	p1 := VectorDifference(gl.Point1, gl.Point2)
	p2 := VectorDifference(p, gl.Point1)
	cross := VectorCrossProduct(p1, p2)
	if math.Abs(0-cross) < FLOAT_DIFF {
		return true
	}
	return false
}

//两点的向量差
func VectorDifference(p1 GeoPoint, p2 GeoPoint) GeoPoint {
	return GeoPoint{Lat: p1.Lat - p2.Lat, Lng: p1.Lng - p2.Lng}
}

//两向量叉乘
func VectorCrossProduct(p1 GeoPoint, p2 GeoPoint) float64 {
	return p1.Lat*p2.Lng - p1.Lng*p2.Lat
}

//与另一条直线是否相交、平行
func (gl *GeoLine) IsIntersectWithLine(line GeoLine) (isIntersect bool, isParallel bool) {
	_, isParallel, isIntersect = gl.GetIntersectPoint(line)
	return
}

//求两直线交点的坐标
//参考：https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect
//返回值：交点、是否平行、是否相交
func (gl *GeoLine) GetIntersectPoint(line GeoLine) (interPoint GeoPoint, isParallel bool, isIntersect bool) {
	//其中一条直线是点的情况
	if gl.Length() == 0 {
		if line.IsContainPoint(gl.Point1) {
			interPoint = gl.Point1
			isParallel = false
			isIntersect = true
			return
		} else {
			isParallel = false
			isIntersect = false
			return
		}
	}
	if line.Length() == 0 {
		if gl.IsContainPoint(line.Point1) {
			interPoint = line.Point1
			isParallel = false
			isIntersect = true
			return
		} else {
			isParallel = false
			isIntersect = false
			return
		}
	}
	p := gl.Point1
	//一线段的向量差
	r := VectorDifference(gl.Point2, gl.Point1)
	q := line.Point1
	//另一线段的向量差
	s := VectorDifference(line.Point2, line.Point1)
	//两线段向量差的叉乘
	rCrossS := VectorCrossProduct(r, s)
	qMinusP := VectorDifference(q, p)
	if rCrossS == 0 {
		//If r × s = 0 and (q − p) × r = 0, then the two lines are collinear.
		//同一条直线，随便返回一个点即可
		if VectorCrossProduct(qMinusP, r) == 0 {
			isParallel = true
			isIntersect = true
			interPoint = line.Point1
			return
		} else {
			//If r × s = 0 and (q − p) × r ≠ 0,
			//then the two lines are parallel and non-intersecting.
			//两条为平行线，没有交点
			isParallel = true
			isIntersect = false
			return
		}
	}
	//If r × s ≠ 0 and 0 ≤ t ≤ 1 and 0 ≤ u ≤ 1,
	//the two line segments meet at the point p + t r = q + u s
	t := VectorCrossProduct(qMinusP, s) / rCrossS
	u := VectorCrossProduct(qMinusP, r) / rCrossS
	//有相交点
	if t >= 0 && t <= 1 && u >= 0 && u <= 1 {
		p1 := GeoPoint{Lat: gl.Point1.Lat + t*r.Lat, Lng: gl.Point1.Lng + t*r.Lng}
		p2 := GeoPoint{Lat: line.Point1.Lat + u*s.Lat, Lng: line.Point1.Lng + u*s.Lng}
		interPoint = p1
		//如果在计算的时候有点小小的误差，这里直接取中间得了，理论上这两个点应该相等
		if !p1.IsEqual(p2) {
			interPoint = MidPoint(p1, p2)
		}
		isParallel = false
		isIntersect = true
		return
	}
	isParallel = false
	isIntersect = false
	return
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

//判断给定的经纬度是否在小格子内，包括边界
func (gr *GeoRectangle) IsPointInRect(point GeoPoint) bool {
	if point.Lat <= gr.MaxLat && point.Lat >= gr.MinLat && point.Lng >= gr.MinLng && point.Lng <= gr.MaxLng {
		return true
	}
	return false
}

//判断给定的经纬度是否完全在小格子内，不包括边界
func (gr *GeoRectangle) IsPointAllInRect(point GeoPoint) bool {
	if point.Lat < gr.MaxLat && point.Lat > gr.MinLat && point.Lng > gr.MinLng && point.Lng < gr.MaxLng {
		return true
	}
	return false
}

//获取中心点坐标
func (gr *GeoRectangle) MidPoint() GeoPoint {
	point := MidPoint(GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng}, GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng})
	return point
}

//矩形X方向的边长，即纬度线方向，保持纬度相同即可，单位米
func (gr *GeoRectangle) Width() float64 {
	return EarthDistance(
		GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng},
		GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng},
	)
}

//矩形Y方向的边长，即经度线方向，保持经度相同即可，单位米
func (gr *GeoRectangle) Height() float64 {
	return EarthDistance(
		GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng},
		GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng},
	)
}

//矩形的所有的点
func (gr *GeoRectangle) GetRectVertex() (ret []GeoPoint) {
	ret = append(ret,
		GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng},
		GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng},
		GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng},
		GeoPoint{Lat: gr.MaxLat, Lng: gr.MinLng},
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
func MakeGeoPolygon(points []GeoPoint) GeoPolygon {
	pl := len(points)
	if pl >= 3 {
		if !points[0].IsEqual(points[pl-1]) {
			points = append(points, points[0])
		}
	}
	return GeoPolygon{Points: points}
}

//一个多边形
type GeoPolygon struct {
	Points  []GeoPoint   `json:"points"`  //一堆顶点，必须是首尾相连有序的
	Borders []GeoLine    `json:"borders"` //所有的边
	Rect    GeoRectangle `json:"rect"`    //最小外包矩形
}

//获取所有的顶点
func (gp *GeoPolygon) GetPoints() []GeoPoint {
	return gp.Points
}

//多边形的所有的边
func (gp *GeoPolygon) GetPolygonBorders() (ret []GeoLine) {
	if !gp.Check() {
		return
	}
	if len(gp.Borders) > 0 {
		return gp.Borders
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
	gp.Borders = ret
	return
}

//添加点
func (gp *GeoPolygon) AddPoint(p GeoPoint) {
	gp.Points = append(gp.Points, p)
}

//将多边形处理成字符串的切片格式
func (gp *GeoPolygon) FormatStringArray() (ret []string) {
	for _, p := range gp.Points {
		ret = append(ret, p.FormatStr())
	}
	return
}

//判断是否是合法的多边形
func (gp *GeoPolygon) Check() bool {
	if len(gp.Points) < 3 {
		return false
	}
	//如果多边形边长超过100km，玩去，没法处理了！
	rect := gp.GetBoundsRect()
	width := rect.Width()
	height := rect.Height()
	var maxDist float64 = 100000
	if width >= maxDist || height >= maxDist {
		return false
	}
	return true
}

//获取最小外包矩形
func (gp *GeoPolygon) GetBoundsRect() GeoRectangle {
	if gp.Rect.Width() > 0 || gp.Rect.Height() > 0 {
		return gp.Rect
	}
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
	rect := GeoRectangle{MaxLat: maxLat, MaxLng: maxLng, MinLat: minLat, MinLng: minLng}
	gp.Rect = rect
	return rect
}

//判断点是否在多边形内部，此处使用最简单的射线法判断
//边数较多时性能不高，只适合在写入时小批量判断
//计算射线与多边形各边的交点，如果是偶数，则点在多边形外，否则在多边形内。
//还会考虑一些特殊情况，如点在多边形顶点上，点在多边形边上等特殊情况。
//参考：http://api.map.baidu.com/library/GeoUtils/1.2/src/GeoUtils.js
func (gp *GeoPolygon) IsPointInPolygon(p GeoPoint) bool {
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
		//射线有可能有交点
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
					if math.Abs(p.Lng-xInters) < FLOAT_DIFF {
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

//按geohash方式，将多边形切割成一个一个小格子，这些小格子至少有一个点是位于多边形内的
//主要应对ES支持自定义配送范围时，有的商家的配送范围巨特么的变态的情况。
//直接用多边形查会把ES查死，这里将配送范围化成一个个小格子，写到term索引里
//基本思路：取多边形的最小外包矩形，先切这个矩形成一个个小的geohash小格子，
//再判断每个格子与多边形是否存在包含或相交的情况，效率较低！还要找更好的办法
//返回值：所有跟多边形有交集的小格子、完全被包围的小格子、部分在多边形内部的小格子
func (gp *GeoPolygon) SplitGeoHashRect(precision int) (inRect, interRect []string) {
	if !gp.Check() {
		return
	}

	boundsRect := gp.GetBoundsRect()
	width := boundsRect.Width()
	height := boundsRect.Height()
	polygonBroders := gp.GetPolygonBorders()

	//先从左下角开始：经度、纬度均最小！然后将此小格子逐步向右、向上推进
	leftUpBaseGeoHash, geoHashBaseRect := GeoHashEncode(boundsRect.MinLat, boundsRect.MinLng, precision)

	//将当前格子添加进去
	tmpGeoHashList := map[string]bool{}
	tmpGeoHashList[leftUpBaseGeoHash] = true

	//计算外包矩形左下角到该geohash格子的右框与矩形交点的距离，即第一个小格子跟最小外包矩形重合区域的宽高
	xBaseLen := EarthDistance(
		GeoPoint{Lat: boundsRect.MinLat, Lng: boundsRect.MinLng},
		GeoPoint{Lat: boundsRect.MinLat, Lng: geoHashBaseRect.MaxLng},
	)
	yBaseLen := EarthDistance(
		GeoPoint{Lat: boundsRect.MinLat, Lng: boundsRect.MinLng},
		GeoPoint{Lat: geoHashBaseRect.MaxLat, Lng: boundsRect.MinLng},
	)

	//小格子的宽高
	geoHashRectWidth := geoHashBaseRect.Width()
	geoHashRectHeight := geoHashBaseRect.Height()

	//修正最小外包矩形的宽高
	width += geoHashRectWidth
	height += geoHashRectHeight

	//设置初始值
	xLen := xBaseLen
	var rect1 = geoHashBaseRect
	var xrect *GeoRectangle
	var ghash string

	//然后向右扩展
	for xLen <= width {
		midPoint := rect1.MidPoint()
		ghash, xrect = GeoHashEncode(midPoint.Lat, midPoint.Lng+rect1.LngSpan(), precision)
		tmpGeoHashList[ghash] = true
		xLen += xrect.Width()
		yLen := yBaseLen
		//开始向上推进
		for yLen <= height {
			midPoint := rect1.MidPoint()
			ghash, rect1 = GeoHashEncode(midPoint.Lat+rect1.LatSpan(), midPoint.Lng, precision)
			tmpGeoHashList[ghash] = true
			yLen += geoHashRectHeight
		}
		//再将小格子向右推进一位
		rect1 = xrect
	}

	//再逐个对小格子判断：要么小格子的边跟多边形的边有相交、要么小格子有顶点在多边形内
	for ghash := range tmpGeoHashList {
		//小格子对应的矩形区域
		grect := GeoHashDecode(ghash)
		//小格子的所有的边
		borders := grect.GetRectBorders()
		//小格子所有的顶点
		points := grect.GetRectVertex()
		isHit := false
		//遍历小格子所有的边，跟多边形的任一边相交就算成功
		for _, gb := range borders {
			if isHit {
				break
			}
			//遍历多边形的边
			for _, pb := range polygonBroders {
				//如果两边相交的话
				if inter, _ := gb.IsIntersectWithLine(pb); inter {
					isHit = true
					break
				}
			}
		}
		//跟任一边相交
		if isHit {
			interRect = append(interRect, ghash)
			continue
		}
		//遍历小格子所有的顶点
		inPointNum := 0
		for _, p := range points {
			if gp.IsPointInPolygon(p) {
				inPointNum++
			}
		}
		//至少有一个点在多边形的内部
		if inPointNum > 0 {
			//所有的顶点均在多边形内部
			if inPointNum == 4 {
				inRect = append(inRect, ghash)
			} else {
				interRect = append(interRect, ghash)
			}
			continue
		}
	}

	return
}

/**
用类似射线法的思想去将多边形切成多个小格子
先找最小的被外包的geohash的矩形，这个可能要比最小外包矩形大一些，按指定精度填满各个小格子，这个矩形一定包含整数个geohash格子
跟多边形有交点的格子肯定是只有部分区域和多边形重合的格子，无须判断直接将其拎出来即可。
从左至右逐个遍历剩下的小格子，对每个小格子做如下的判断：
	以小格子的顶边向两边延伸，求与此直接相交的边数，如果两方向的相交的边均为奇数个则为在多边形内部，否则在外部。
这里判断交点的个数的规则：
	在此直线上侧相交的边都不算，只要下侧的
	如果一条边跟此直线相交，则被切成两个边，只保留下侧的边，看看与多少条边相交
	左右两个方向上都为奇数个的说明此小格子完全在多边形内部，否则在外部
http://willdemaine.ghost.io/filling-geofences-with-geohashes/
*/
func (gp *GeoPolygon) RaySplitGeoHashRect(precision int) (inRect, interRect []string) {
	if !gp.Check() {
		return
	}

	//切格子用的geohash精度
	stp := precision

	//先提取最小外包矩形，并取取四个点的geohash格子
	minRect := gp.GetBoundsRect()

	//一个临时小格子，计算经纬度差用的
	_, tmpRect := GeoHashEncode(minRect.MaxLat, minRect.MinLng, stp)
	tmpMidPoint := tmpRect.MidPoint()
	diffLat := tmpRect.MaxLat - tmpRect.MinLat
	diffLng := tmpRect.MaxLng - tmpRect.MinLng

	//如果外包矩形太小，取其中心点的
	if minRect.Width() <= tmpRect.Width() && minRect.Height() <= tmpRect.Height() {
		mid := minRect.MidPoint()
		geo, _ := GeoHashEncode(mid.Lat, mid.Lng, stp)
		interRect = append(interRect, geo)
		return
	}

	//左上角的格子，基准点
	_, leftUpRect := GeoHashEncode(tmpMidPoint.Lat, tmpMidPoint.Lng, stp)
	//左下角的小格子，经纬、纬度全最小
	_, leftBottom := GeoHashEncode(minRect.MinLat, minRect.MinLng, stp)
	//右上角的小格子，经纬、纬度全最大
	_, rightUp := GeoHashEncode(minRect.MaxLat, minRect.MaxLng, stp)

	//最小的geohash整数倍的外包矩形，可能要比最小外包矩形大一些
	geoRect := GeoRectangle{
		MaxLat: rightUp.MaxLat,
		MaxLng: rightUp.MaxLng,
		MinLat: leftBottom.MinLat,
		MinLng: leftBottom.MinLng,
	}

	//垂直、水平方向的格子数，可能出现xx.9999或者xxx.0001这样的情况
	tmpVNum := geoRect.Height() / rightUp.Height()
	tmpHNum := geoRect.Width() / rightUp.Width()
	verNum := int(tmpVNum)
	horiNum := int(tmpHNum)
	if math.Abs(tmpVNum-float64(verNum)) > 0.1 || verNum <= 0 {
		verNum++
	}
	if math.Abs(tmpHNum-float64(horiNum)) > 0.1 || horiNum <= 0 {
		horiNum++
	}

	//基准点，从这个点开始往右、往下推进
	basePoint := leftUpRect.MidPoint()

	//线段与多边形的交点情况，缓存交点用的，保证每个切线只计算一次
	lineInterMap := map[string][]intersectLineAndPoint{}

	//小格子的四个边的延长线与多边形交点情况
	var topInters, bottomInters, leftInters, rightInters []intersectLineAndPoint
	var ok bool

	//从左到右、从上到下遍历所有的小格子
	for vi := 0; vi < verNum; vi++ {
		//最左边的小格子的经度都一样，纬度逐次减小
		baseLat := basePoint.Lat - float64(vi)*diffLat
		baseLng := basePoint.Lng
		for hi := 0; hi < horiNum; hi++ {
			//当前小格子的geo值及小格子矩形
			geo, tRect := GeoHashEncode(baseLat, baseLng+float64(hi)*diffLng, stp)

			//小格子顶边及与多边形交点情况
			topLine := GeoLine{
				Point1: GeoPoint{Lat: tRect.MaxLat, Lng: geoRect.MinLng - 1},
				Point2: GeoPoint{Lat: tRect.MaxLat, Lng: geoRect.MaxLng + 1},
			}
			if topInters, ok = lineInterMap[topLine.FormatStr()]; !ok {
				topInters = gp.interPointsWithTopLine(topLine)
				lineInterMap[topLine.FormatStr()] = topInters
			}

			//小格子底边及与多边形交点情况
			bottomLine := GeoLine{
				Point1: GeoPoint{Lat: tRect.MinLat, Lng: geoRect.MinLng - 1},
				Point2: GeoPoint{Lat: tRect.MinLat, Lng: geoRect.MaxLng + 1},
			}
			if bottomInters, ok = lineInterMap[bottomLine.FormatStr()]; !ok {
				bottomInters = gp.interPointsWithBottomLine(bottomLine)
				lineInterMap[bottomLine.FormatStr()] = bottomInters
			}

			//小格子左边框及与多边形交点情况
			leftLine := GeoLine{
				Point1: GeoPoint{Lat: geoRect.MaxLat + 1, Lng: tRect.MinLng},
				Point2: GeoPoint{Lat: geoRect.MinLat - 1, Lng: tRect.MinLng},
			}
			if leftInters, ok = lineInterMap[leftLine.FormatStr()]; !ok {
				leftInters = gp.interPointsWithLeftLine(leftLine)
				lineInterMap[leftLine.FormatStr()] = leftInters
			}

			//小格子右边框及与多边形交点情况
			rightLine := GeoLine{
				Point1: GeoPoint{Lat: geoRect.MaxLat + 1, Lng: tRect.MaxLng},
				Point2: GeoPoint{Lat: geoRect.MinLat - 1, Lng: tRect.MaxLng},
			}
			if rightInters, ok = lineInterMap[rightLine.FormatStr()]; !ok {
				rightInters = gp.interPointsWithRightLine(rightLine)
				lineInterMap[rightLine.FormatStr()] = rightInters
			}

			isContinue := false

			//底边框跟多边形的交点不符合在多边形内部的情况
			for _, inter := range bottomInters {
				if inter.point.Lng > tRect.MinLng && inter.point.Lng < tRect.MaxLng {
					interRect = append(interRect, geo)
					isContinue = true
					break
				}
			}
			if isContinue {
				continue
			}

			//左边垂线的交点在边框上
			for _, inter := range leftInters {
				if inter.point.Lat < tRect.MaxLat && inter.point.Lat > tRect.MinLat {
					interRect = append(interRect, geo)
					isContinue = true
					break
				}
			}
			if isContinue {
				continue
			}

			//右边框的交点在边框上
			for _, inter := range rightInters {
				if inter.point.Lat < tRect.MaxLat && inter.point.Lat > tRect.MinLat {
					interRect = append(interRect, geo)
					isContinue = true
					break
				}
			}
			if isContinue {
				continue
			}

			//对于上下边框，判断小格子左右两边的交点情况
			leftNum := 0
			rightNum := 0
			internalInterNum := 0
			for _, inter := range topInters {
				if inter.point.Lng <= tRect.MinLng {
					leftNum++
					continue
				}
				if inter.point.Lng >= tRect.MaxLng {
					rightNum++
					continue
				}
				if inter.point.Lng > tRect.MinLng && inter.point.Lng < tRect.MaxLng {
					interRect = append(interRect, geo)
					isContinue = true
					break
				}
				internalInterNum++
			}
			if isContinue {
				continue
			}
			if leftNum <= 0 || rightNum <= 0 || leftNum%2 == 0 || rightNum%2 == 0 {
				if internalInterNum > 0 {
					interRect = append(interRect, geo)
					continue
				}
			} else {
				inRect = append(inRect, geo)
				continue
			}
		}
	}

	return
}

//一条横线跟多边形的哪些边相交，以及交点
type intersectLineAndPoint struct {
	border GeoLine  //相交的边
	point  GeoPoint //跟此边的交点
}

//一条横线和多边形的交点，与横线部分重合的不算、交点在多边形顶点的时位于直接上方的不算
func (gp *GeoPolygon) interPointsWithTopLine(line GeoLine) (ret []intersectLineAndPoint) {
	maxLng := math.Max(line.Point1.Lng, line.Point2.Lng)
	minLng := math.Max(line.Point1.Lng, line.Point2.Lng)
	lineLat := line.Point2.Lat
	borders := gp.GetPolygonBorders()
	for _, border := range borders {
		//由于line是一条直线，纬度相等，只不过经度有变化
		if border.Point2.Lat > lineLat && border.Point1.Lat > lineLat {
			continue
		}
		if border.Point2.Lat < lineLat && border.Point1.Lat < lineLat {
			continue
		}
		//如果交点在其顶点上，并且另一点的纬度大于横线的不要，否则就算有交点
		if border.Point1.Lat == lineLat && border.Point1.Lng >= minLng && border.Point1.Lng <= maxLng {
			if border.Point2.Lat <= line.Point2.Lat {
				ret = append(ret, intersectLineAndPoint{
					border: border,
					point:  border.Point1,
				})
			}
			continue
		}
		if border.Point2.Lat == lineLat && border.Point2.Lng >= minLng && border.Point2.Lng <= maxLng {
			if border.Point1.Lat <= line.Point2.Lat {
				ret = append(ret, intersectLineAndPoint{
					border: border,
					point:  border.Point2,
				})
			}
			continue
		}
		//普通的相交
		p, isParallel, isInter := border.GetIntersectPoint(line)
		if isInter && !isParallel {
			ret = append(ret, intersectLineAndPoint{
				border: border,
				point:  p,
			})
		}
	}
	return
}

//geohash的小格子的底边延长线和多边形各边的交点，只要此线上方部分
func (gp *GeoPolygon) interPointsWithBottomLine(line GeoLine) (ret []intersectLineAndPoint) {
	borders := gp.GetPolygonBorders()
	for _, border := range borders {
		//由于line是一条直线，纬度相等，只不过经度有变化
		if border.Point2.Lat > line.Point2.Lat && border.Point1.Lat > line.Point2.Lat {
			continue
		}
		if border.Point2.Lat < line.Point2.Lat && border.Point1.Lat < line.Point2.Lat {
			continue
		}
		//普通的相交
		p, isParallel, isInter := border.GetIntersectPoint(line)
		if isInter && !isParallel {
			ret = append(ret, intersectLineAndPoint{
				border: border,
				point:  p,
			})
		}
	}
	return
}

//一条垂线和多边形的交点
func (gp *GeoPolygon) interPointsWithLeftLine(line GeoLine) (ret []intersectLineAndPoint) {
	borders := gp.GetPolygonBorders()
	for _, border := range borders {
		if border.Point2.Lng > line.Point2.Lng && border.Point1.Lng > line.Point2.Lng {
			continue
		}
		if border.Point2.Lng < line.Point2.Lng && border.Point1.Lng < line.Point2.Lng {
			continue
		}
		//普通的相交
		p, isParallel, isInter := border.GetIntersectPoint(line)
		if isInter && !isParallel {
			ret = append(ret, intersectLineAndPoint{
				border: border,
				point:  p,
			})
		}
	}
	return
}

//一条垂线和多边形的交点
func (gp *GeoPolygon) interPointsWithRightLine(line GeoLine) (ret []intersectLineAndPoint) {
	borders := gp.GetPolygonBorders()
	for _, border := range borders {
		if border.Point2.Lng > line.Point2.Lng && border.Point1.Lng > line.Point2.Lng {
			continue
		}
		if border.Point2.Lng < line.Point2.Lng && border.Point1.Lng < line.Point2.Lng {
			continue
		}
		//普通的相交
		p, isParallel, isInter := border.GetIntersectPoint(line)
		if isInter && !isParallel {
			ret = append(ret, intersectLineAndPoint{
				border: border,
				point:  p,
			})
		}
	}
	return
}
