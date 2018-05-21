/**
 * @author      Liu Yongshuai<liuyongshuai@hotmail.com>
 * @package     geo
 * @date        2018-04-19 14:34
 */
package geo

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
)

var Testpoints = []GeoPoint{
	{Lat: 31.473712, Lng: 120.270125},
	{Lat: 31.552752, Lng: 120.2372409},
	{Lat: 31.575529, Lng: 120.3026},
	{Lat: 31.676580, Lng: 120.357945},
	{Lat: 31.675095, Lng: 120.309268963407},
	{Lat: 31.580200, Lng: 120.35321},
	{Lat: 31.547996, Lng: 120.365778},
	{Lat: 31.549302, Lng: 120.367247421931},
	{Lat: 31.534714, Lng: 120.295785},
	{Lat: 31.541530, Lng: 120.3526},
	{Lat: 31.546304, Lng: 120.383614},
	{Lat: 31.479810, Lng: 120.449},
	{Lat: 31.542248, Lng: 120.378810983558},
}

func TestGeoHashEncode(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	fmt.Printf("\n\n\n------------start %s------------\n", f.Name())

	for _, pInfo := range Testpoints {
		fmt.Fprintf(os.Stdout, "Lat:%f\tLng:%f\n", pInfo.Lat, pInfo.Lng)
		for i := 4; i <= 10; i++ {
			geo, square := GeoHashEncode(pInfo.Lat, pInfo.Lng, i)
			xdist := square.Width()
			ydist := square.Height()
			//附近9个格子
			gs := GetNeighborsGeoCodes(pInfo.Lat, pInfo.Lng, i)
			dt := fmt.Sprintf("%fm x %fm", xdist, ydist)
			fmt.Fprintf(os.Stdout, "\tprecision:%d\tgeoHash:%-10s\tdist:%-14s\tNeighbors:%s\n", i, geo, dt, strings.Join(gs, ","))
		}
	}
	fmt.Printf("------------end %s------------\n", f.Name())
}

func TestGeoHashDecode(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	fmt.Printf("\n\n\n------------start %s------------\n", f.Name())

	//一坨的地理位置信息
	points := []GeoPoint{
		{Lat: 31.675095455, Lng: 120.309268963407},
		{Lat: 31.54224843, Lng: 120.378810983558},
	}
	for _, pInfo := range points {
		fmt.Fprintf(os.Stdout, "Lat:%f\tLng:%f\n", pInfo.Lat, pInfo.Lng)
		for i := 4; i <= 10; i++ {
			geo, deGeo := GeoHashEncode(pInfo.Lat, pInfo.Lng, i)
			p := deGeo.MidPoint()
			midLat, midLng := p.Lat, p.Lng
			fmt.Fprintf(os.Stdout, "\tprecision:%d\tgeoHash:%-10s\t\n", i, geo)
			fmt.Fprintf(os.Stdout, "\t\tmaxLat:%v\tminLat:%v\n", deGeo.MaxLat, deGeo.MinLat)
			fmt.Fprintf(os.Stdout, "\t\tminLng:%v\tmaxLng:%v\n", deGeo.MinLng, deGeo.MaxLng)
			fmt.Fprintf(os.Stdout, "\t\tmidLat:%v\tmidLng:%v\n", midLat, midLng)
			fmt.Println("\tdeGeo")
			deGeo = GeoHashDecode(geo)
			fmt.Fprintf(os.Stdout, "\t\tmaxLat:%v\tminLat:%v\n", deGeo.MaxLat, deGeo.MinLat)
			fmt.Fprintf(os.Stdout, "\t\tminLng:%v\tmaxLng:%v\n", deGeo.MinLng, deGeo.MaxLng)
			fmt.Fprintf(os.Stdout, "\t\tmidLat:%v\tmidLng:%v\n", midLat, midLng)

		}
	}
	fmt.Printf("------------end %s------------\n", f.Name())
}

func TestDiffLatLng(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	fmt.Printf("\n\n\n------------start %s------------\n", f.Name())

	distance := 54392.02
	fmt.Fprintf(os.Stdout, "distance=%v\n", distance)
	var angle float64 = 33
	for index, point := range Testpoints {
		ang := float64(index) * angle
		p := PointAtDistAndAngle(point, distance, ang)
		dist := EarthDistance(p, point)
		fmt.Fprintf(os.Stdout, "oriLat:%f\toriLng:%f\tdiff:%v\tdist:%f\tangle=%v\n", point.Lat, point.Lng, p, dist, ang)
	}
	fmt.Printf("------------end %s------------\n", f.Name())
}
func TestGeoLine_GetIntersectPoint(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	fmt.Printf("\n\n\n------------start %s------------\n", f.Name())
	p1 := GeoPoint{Lat: 39.998006, Lng: 116.322415}
	p2 := GeoPoint{Lat: 39.851969, Lng: 116.434991}
	p3 := GeoPoint{Lat: 40.004818, Lng: 116.474552}
	p4 := GeoPoint{Lat: 39.850391, Lng: 116.326583}
	line1 := MakeGeoLine(p1, p2)
	line2 := MakeGeoLine(p3, p4)
	p, b := line1.GetIntersectPoint(line2)
	fmt.Println(p, b)
	p1 = GeoPoint{Lat: 39.998006, Lng: 116.322415}
	p2 = GeoPoint{Lat: 39.851969, Lng: 116.434991}
	p3 = GeoPoint{Lat: 39.908006, Lng: 116.322415}
	p4 = GeoPoint{Lat: 39.908006, Lng: 116.922415}
	line1 = MakeGeoLine(p1, p2)
	line2 = MakeGeoLine(p3, p4)
	p, b = line1.GetIntersectPoint(line2)
	fmt.Println(p, b)
	fmt.Println(line1.IsContainPoint(p))
	fmt.Printf("------------end %s------------\n", f.Name())
}
func TestMidPoint(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	fmt.Printf("\n\n\n------------start %s------------\n", f.Name())
	point := GeoPoint{Lat: 39.43373712, Lng: 120.378810983558}
	for _, p := range Testpoints {
		midPoint := MidPoint(point, p)
		dist1 := EarthDistance(point, p)
		dist2 := EarthDistance(point, midPoint)
		dist3 := EarthDistance(p, midPoint)
		fmt.Fprintf(os.Stdout, "allDist:%f\tdist1:%f\tdist2:%f\n", dist1, dist2, dist3)
	}
	fmt.Printf("------------end %s------------\n", f.Name())
}

func TestFormatDistance(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	fmt.Printf("\n\n\n------------start %s------------\n", f.Name())
	str := []string{
		"5000.00m",
		"44.1km",
		"43535345",
		"km222",
		"44km.01",
		"44.44KM",
	}
	for _, s := range str {
		ret := FormatDistance(s)
		fmt.Fprintf(os.Stdout, "str:%s\tdist:%v\n", s, ret)
	}
	fmt.Printf("------------end %s------------\n", f.Name())
}

func TestGeoGridBox_InGridBox(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	fmt.Printf("\n\n\n------------start %s------------\n", f.Name())
	box := GeoRectangle{
		MaxLng: 120.378810983558,
		MinLng: 120.309268963407,
		MaxLat: 31.547996,
		MinLat: 31.047996,
	}
	fmt.Println(box.Width())
	fmt.Println(box.Height())
	fmt.Println(box.IsPointInRect(GeoPoint{Lat: 31.547996, Lng: 120.378810}))
	fmt.Println(box.IsPointInRect(GeoPoint{Lat: 31.947996, Lng: 120.378810}))
	fmt.Println(box.IsPointInRect(GeoPoint{Lat: 31.547996, Lng: 120.978810}))
	fmt.Println(box.IsPointInRect(GeoPoint{Lat: 31.947996, Lng: 120.978810}))
	fmt.Println(box.IsPointInRect(box.MidPoint()))
	fmt.Printf("------------end %s------------\n", f.Name())
}

//点是否在多边形内
func TestGeoPloygon_IsPointInPolygon(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	fmt.Printf("\n\n\n------------start %s------------\n", f.Name())
	ploygon := GeoPloygon{}
	ploygon.AddPoint(GeoPoint{Lat: 39.972907, Lng: 116.322631})
	ploygon.AddPoint(GeoPoint{Lat: 39.937953, Lng: 116.346777})
	ploygon.AddPoint(GeoPoint{Lat: 39.902095, Lng: 116.322056})
	ploygon.AddPoint(GeoPoint{Lat: 39.883051, Lng: 116.39622})
	ploygon.AddPoint(GeoPoint{Lat: 39.96583, Lng: 116.488206})
	ploygon.AddPoint(GeoPoint{Lat: 39.992368, Lng: 116.436464})
	ploygon.AddPoint(GeoPoint{Lat: 39.972907, Lng: 116.322631})
	is1 := ploygon.IsPointInPolygon(GeoPoint{Lat: 39.946804, Lng: 116.383572})
	fmt.Println(is1)
	is2 := ploygon.IsPointInPolygon(GeoPoint{Lat: 40.043204, Lng: 116.394495})
	fmt.Println(is2)
	is3 := ploygon.IsPointInPolygon(GeoPoint{Lat: 39.992368, Lng: 116.436464})
	fmt.Println(is3)
	fmt.Printf("------------end %s------------\n", f.Name())
}

func TestGeoLine_IsIntersect(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	fmt.Printf("\n\n\n------------start %s------------\n", f.Name())
	p1 := GeoPoint{Lat: 39.990599, Lng: 116.446813}
	p2 := GeoPoint{Lat: 39.895009, Lng: 116.262265}
	line1 := GeoLine{Point1: p1, Point2: p2}
	p3 := GeoPoint{Lat: 40.050716, Lng: 116.418642}
	line2 := GeoLine{Point1: p1, Point2: p3}
	fmt.Println(line1.IsIntersectWithLine(line2))
	p4 := GeoPoint{Lat: 39.875078, Lng: 116.581918}
	line3 := GeoLine{Point1: p3, Point2: p4}
	fmt.Println(line3.IsIntersectWithLine(line2))
	fmt.Println(line2.IsIntersectWithLine(line2))
	p5 := GeoPoint{Lat: 39.983855, Lng: 116.27635}
	line4 := GeoLine{Point1: p5, Point2: p3}
	fmt.Println(line4.IsIntersectWithLine(line2))
	fmt.Printf("------------end %s------------\n", f.Name())
}

func TestGeoLine_IsContainPoint(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	fmt.Printf("\n\n\n------------start %s------------\n", f.Name())
	basePoint := GeoPoint{Lat: 39.983855, Lng: 116.27635}
	p1 := PointAtDistAndAngle(basePoint, 3000, 88)
	p2 := PointAtDistAndAngle(basePoint, 9000, 88)
	line := MakeGeoLine(p1, p2)
	fmt.Println(line.IsContainPoint(PointAtDistAndAngle(basePoint, 9000, 88)))
	fmt.Println(line.IsContainPoint(PointAtDistAndAngle(basePoint, 8888, 88)))
	fmt.Println(line.IsContainPoint(PointAtDistAndAngle(basePoint, 10000, 88)))
	fmt.Println(line.IsContainPoint(GeoPoint{Lat: 39.783855, Lng: 116.27635}))
	fmt.Printf("------------end %s------------\n", f.Name())
}

func TestGeoHashDecodeBits(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	fmt.Printf("\n\n\n------------start %s------------\n", f.Name())
	lat := 39.956981
	lng := 116.440488
	var i uint8
	for i = 1; i < 32; i++ {
		geo, _ := GeoHashBitsEncode(lat, lng, i)
		rect := GeoHashBitsDecode(geo, i)
		fmt.Fprintf(os.Stdout, "precision=%d\tdist=%v x %v\n", i, int(rect.Width()), int(rect.Height()))
	}
	fmt.Printf("------------end %s------------\n", f.Name())
}

func TestGeoHashBitsNeighbors(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	fmt.Printf("\n\n\n------------start %s------------\n", f.Name())
	lat := 39.956981
	lng := 116.440488
	var pre uint8 = 15
	nebs := GeoHashBitsNeighbors(lat, lng, pre)
	for _, n := range nebs {
		rect := GeoHashBitsDecode(n, pre)
		fmt.Fprintf(os.Stdout, "dist=%v x %v\n", int(rect.Width()), int(rect.Height()))
		fmt.Println(rect.GetRectVertex())
	}
	fmt.Printf("------------end %s------------\n", f.Name())
}
