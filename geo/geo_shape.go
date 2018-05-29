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

//左下角的坐标：经纬度最小
func (gr *GeoRectangle) LeftBottomPoint() GeoPoint {
	return GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng}
}

//左上角的坐标：经度最小、纬度最大
func (gr *GeoRectangle) LeftUpPoint() GeoPoint {
	return GeoPoint{Lat: gr.MaxLat, Lng: gr.MinLng}
}

//右上角的坐标：经纬度最大
func (gr *GeoRectangle) RightUpPoint() GeoPoint {
	return GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng}
}

//右下角的坐标：经度最大、纬度最小
func (gr *GeoRectangle) RightBottomPoint() GeoPoint {
	return GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng}
}

//左边框线段，从上往下指
func (gr *GeoRectangle) LeftBorder() GeoLine {
	return GeoLine{
		Point1: GeoPoint{Lat: gr.MaxLat, Lng: gr.MinLng},
		Point2: GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng},
	}
}

//右边框线段，从上往下指
func (gr *GeoRectangle) RightBorder() GeoLine {
	return GeoLine{
		Point1: GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng},
		Point2: GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng},
	}
}

//上边框线段，从左往右指
func (gr *GeoRectangle) TopBorder() GeoLine {
	return GeoLine{
		Point1: GeoPoint{Lat: gr.MaxLat, Lng: gr.MinLng},
		Point2: GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng},
	}
}

//下边框线段，从左往右指
func (gr *GeoRectangle) BottomBorder() GeoLine {
	return GeoLine{
		Point1: GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng},
		Point2: GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng},
	}
}

//矩形的左上角、右下角的对象线
func (gr *GeoRectangle) LeftUp2RightBottomLine() GeoLine {
	return GeoLine{
		Point1: GeoPoint{Lat: gr.MaxLat, Lng: gr.MinLng},
		Point2: GeoPoint{Lat: gr.MinLat, Lng: gr.MaxLng},
	}
}

//矩形的左下角、右上角的对象线
func (gr *GeoRectangle) LeftBottom2RightUpLine() GeoLine {
	return GeoLine{
		Point1: GeoPoint{Lat: gr.MinLat, Lng: gr.MinLng},
		Point2: GeoPoint{Lat: gr.MaxLat, Lng: gr.MaxLng},
	}
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

//是否有边相交（相怜的不算）
func (gp *GeoPolygon) IsBorderInterect() bool {
	if !gp.Check() {
		return false
	}
	borders := gp.GetPolygonBorders()
	for _, line1 := range borders {
		for _, line2 := range borders {
			if line1.Point1.IsEqual(line2.Point1) || line1.Point1.IsEqual(line2.Point2) {
				continue
			}
			if line1.Point2.IsEqual(line2.Point1) || line1.Point1.IsEqual(line2.Point2) {
				continue
			}
			isInter, _ := line1.IsIntersectWithLine(line2)
			if isInter {
				return true
			}
		}
	}
	return false
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
				if p3.IsEqual(p2) {
					p3 = points[(i+2)%PNum]
				}
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

/**
用类似射线法的思想去将多边形切成多个小格子
见测试用例：geopolygon_test.go，它会生成一堆html文件，用百度地图画出来所有的格子
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
func (gp *GeoPolygon) SplitGeoHashRect(precision int) (inRect, interRect []string) {
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
	lineInterCache := map[string]map[GeoLine]GeoPoint{}

	//小格子的四个边的延长线与多边形交点情况
	var topInters, bottomInters, leftInters, rightInters map[GeoLine]GeoPoint
	var ok bool

	//从左到右、从上到下遍历所有的小格子
	for verInterator := 0; verInterator < verNum; verInterator++ {
		//最左边的小格子的经度都一样，纬度逐次减小
		baseLat := basePoint.Lat - float64(verInterator)*diffLat
		baseLng := basePoint.Lng
		for horiInterator := 0; horiInterator < horiNum; horiInterator++ {
			//当前小格子的geo值及小格子矩形
			geo, tRect := GeoHashEncode(baseLat, baseLng+float64(horiInterator)*diffLng, stp)

			//小格子上边框的延长线、及与多边形每个边的交点情况
			topLine := GeoLine{
				Point1: GeoPoint{Lat: tRect.MaxLat, Lng: geoRect.MinLng - 1},
				Point2: GeoPoint{Lat: tRect.MaxLat, Lng: geoRect.MaxLng + 1},
			}
			if topInters, ok = lineInterCache[topLine.FormatStr()]; !ok {
				topInters = gp.interPointsWithHorizontalLine(topLine)
				lineInterCache[topLine.FormatStr()] = topInters
			}

			//小格子下边框延长线、及与多边形各边交点情况
			bottomLine := GeoLine{
				Point1: GeoPoint{Lat: tRect.MinLat, Lng: geoRect.MinLng - 1},
				Point2: GeoPoint{Lat: tRect.MinLat, Lng: geoRect.MaxLng + 1},
			}
			if bottomInters, ok = lineInterCache[bottomLine.FormatStr()]; !ok {
				bottomInters = gp.interPointsWithHorizontalLine(bottomLine)
				lineInterCache[bottomLine.FormatStr()] = bottomInters
			}

			//小格子左边框延长线、及与多边形各边交点情况
			leftLine := GeoLine{
				Point1: GeoPoint{Lat: geoRect.MaxLat + 1, Lng: tRect.MinLng},
				Point2: GeoPoint{Lat: geoRect.MinLat - 1, Lng: tRect.MinLng},
			}
			if leftInters, ok = lineInterCache[leftLine.FormatStr()]; !ok {
				leftInters = gp.interPointsWithVertialLine(leftLine)
				lineInterCache[leftLine.FormatStr()] = leftInters
			}

			//小格子右边框延长线、及与多边形各边交点情况
			rightLine := GeoLine{
				Point1: GeoPoint{Lat: geoRect.MaxLat + 1, Lng: tRect.MaxLng},
				Point2: GeoPoint{Lat: geoRect.MinLat - 1, Lng: tRect.MaxLng},
			}
			if rightInters, ok = lineInterCache[rightLine.FormatStr()]; !ok {
				rightInters = gp.interPointsWithVertialLine(rightLine)
				lineInterCache[rightLine.FormatStr()] = rightInters
			}

			//TODO 既然得到了上下左右四线的交点情况，有没有可能将这一排或这一列的小格子都一起判断？

			isContinue := false

			//底边框跟多边形的交点不符合在多边形内部的情况
			for border, interPoint := range bottomInters {
				//如果交点位于小格子底边框的非顶点上，即只位于底边框线段中间某个位置上
				//且多边形的相应边的另一顶点在下边框的上方，此时必定是半包围的小格子
				if interPoint.Lng > tRect.MinLng && interPoint.Lng < tRect.MaxLng {
					if border.Point1.Lat > interPoint.Lat || border.Point2.Lat > interPoint.Lat {
						interRect = append(interRect, geo)
						isContinue = true
						break
					}
				}
				//考虑到特殊情况，即此边正好和小格子对角线部分重合
				//如果交点在小格子的某个角上，同时又跟对角线上对应的点相交
				topInterPoint, ok := topInters[border]
				if !ok {
					continue
				}
				if interPoint.Lng == tRect.MinLng && topInterPoint.Lng == tRect.MaxLng ||
					interPoint.Lng == tRect.MaxLng && topInterPoint.Lng == tRect.MinLng {
					interRect = append(interRect, geo)
					isContinue = true
					break
				}
			}
			if isContinue {
				continue
			}

			//左边垂线的交点在边框上
			for border, interPoint := range leftInters {
				//如果交点位于小格子左边框的非顶点上，即只位于左边框线段中间某个位置上
				//且多边形的相应边的另一顶点在左边框的右方，此时必定是半包围的小格子
				if interPoint.Lat < tRect.MaxLat && interPoint.Lat > tRect.MinLat {
					if border.Point1.Lng > interPoint.Lng || border.Point2.Lng > interPoint.Lng {
						interRect = append(interRect, geo)
						isContinue = true
						break
					}
				}
				//考虑到特殊情况，即此边正好和小格子对角线部分重合
				//如果交点在小格子的某个角上，同时又跟对角线上对应的点相交
				rightInterPoint, ok := rightInters[border]
				if !ok {
					continue
				}
				if interPoint.Lat == tRect.MaxLat && rightInterPoint.Lat == tRect.MinLat ||
					interPoint.Lat == tRect.MinLat && rightInterPoint.Lat == tRect.MaxLat {
					interRect = append(interRect, geo)
					isContinue = true
					break
				}
			}
			if isContinue {
				continue
			}

			//右边框的交点在边框上
			for border, interPoint := range rightInters {
				//如果交点位于小格子右边框的非顶点上，即只位于右边框线段中间某个位置上
				//且多边形的相应边的另一顶点在右边框的左方，此时必定是半包围的小格子
				if interPoint.Lat < tRect.MaxLat && interPoint.Lat > tRect.MinLat {
					if border.Point1.Lng < interPoint.Lng || border.Point2.Lng < interPoint.Lng {
						interRect = append(interRect, geo)
						isContinue = true
						break
					}
				}
			}
			if isContinue {
				continue
			}

			//对于上下边框，判断小格子左右两边的交点情况
			leftNum := 0
			rightNum := 0
			for border, interPoint := range topInters {
				//上边框向左的射线跟多边形的交点情况
				if interPoint.Lng <= tRect.MinLng {
					leftNum++
					continue
				}
				//上边框向右的射线跟多边形的交点情况
				if interPoint.Lng >= tRect.MaxLng {
					rightNum++
					continue
				}
				//如果交点位于小格子上边框的非顶点上，即只位于上边框线段中间某个位置上
				//且多边形的相应边的另一顶点在上边框的下方，此时必定是半包围的小格子
				if border.Point1.Lat < interPoint.Lat || border.Point2.Lat < interPoint.Lat {
					interRect = append(interRect, geo)
					isContinue = true
					break
				}
			}
			if isContinue {
				continue
			}
			if leftNum%2 == 1 && rightNum%2 == 1 {
				inRect = append(inRect, geo)
				continue
			}
		}
	}

	return
}

//一条横线和多边形的交点，与横线部分重合的不算、交点在多边形顶点的时位于直接上方的不算
func (gp *GeoPolygon) interPointsWithHorizontalLine(line GeoLine) (ret map[GeoLine]GeoPoint) {
	ret = map[GeoLine]GeoPoint{}
	maxLng := math.Max(line.Point1.Lng, line.Point2.Lng)
	minLng := math.Min(line.Point1.Lng, line.Point2.Lng)
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
		//平行线不要
		if border.Point1.Lat == border.Point2.Lat {
			continue
		}
		//如果交点在其顶点上，并且另一点的纬度大于横线的不要，否则就算有交点
		if border.Point1.Lat == lineLat && border.Point1.Lng >= minLng && border.Point1.Lng <= maxLng {
			if border.Point2.Lat <= lineLat {
				border.Point1.Lat = lineLat
				ret[border] = border.Point1
			}
			continue
		}
		if border.Point2.Lat == lineLat && border.Point2.Lng >= minLng && border.Point2.Lng <= maxLng {
			if border.Point1.Lat <= lineLat {
				border.Point1.Lat = lineLat
				ret[border] = border.Point2
			}
			continue
		}
		//普通的相交
		p, isParallel, isInter := border.GetIntersectPoint(line)
		if isInter && !isParallel {
			p.Lat = lineLat
			ret[border] = p
		}
	}
	return
}

//一条垂线和多边形的交点
func (gp *GeoPolygon) interPointsWithVertialLine(line GeoLine) (ret map[GeoLine]GeoPoint) {
	ret = map[GeoLine]GeoPoint{}
	lineLng := line.Point2.Lng
	borders := gp.GetPolygonBorders()
	for _, border := range borders {
		if border.Point2.Lng > lineLng && border.Point1.Lng > lineLng {
			continue
		}
		if border.Point2.Lng < lineLng && border.Point1.Lng < lineLng {
			continue
		}
		//垂线不要
		if border.Point1.Lng == border.Point2.Lng {
			continue
		}
		//普通的相交
		p, isParallel, isInter := border.GetIntersectPoint(line)
		if isInter && !isParallel {
			p.Lng = lineLng
			ret[border] = p
		}
	}
	return
}
