/**
 * @author      Liu Yongshuai<liuyongshuai@didichuxing.com>
 * @package     es
 * @date        2018-05-23 15:37
 */
package geo

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"runtime"
	"testing"
	"time"
)

func TestGeoPolygon_SplitGeoHashRect(t *testing.T) {
	pc, _, _, _ := runtime.Caller(0)
	f := runtime.FuncForPC(pc)
	fmt.Printf("\n\n\n------------start %s------------\n", f.Name())
	var polygon GeoPolygon
	polygon = getPolygon1()
	splitGeoHashRect(polygon, "polygon1", 13)
	polygon = getPolygon2()
	splitGeoHashRect(polygon, "polygon2", 13)
	polygon = getPolygon3()
	splitGeoHashRect(polygon, "polygon3", 13)
	polygon = getPolygon4()
	splitGeoHashRect(polygon, "polygon4", 13)
	polygon = getPolygon5()
	splitGeoHashRect(polygon, "polygon5", 13)
	polygon = getPolygon6()
	splitGeoHashRect(polygon, "polygon6", 13)
	polygon = getPolygon7()
	splitGeoHashRect(polygon, "polygon7", 13)
	polygon = getPolygon8()
	splitGeoHashRect(polygon, "polygon8", 13)
	polygon = getPolygon9()
	splitGeoHashRect(polygon, "polygon9", 19)
	polygon = getPolygon10()
	splitGeoHashRect(polygon, "polygon10", 15)
	fmt.Printf("------------end %s------------\n", f.Name())
}

//凸多边形
func getPolygon1() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.385297, Lat: 39.993252},
		{Lng: 116.325505, Lat: 39.974235},
		{Lng: 116.290435, Lat: 39.931314},
		{Lng: 116.346777, Lat: 39.879508},
		{Lng: 116.436464, Lat: 39.911836},
		{Lng: 116.451987, Lat: 39.93751},
		{Lng: 116.449687, Lat: 39.971138},
		{Lng: 116.415767, Lat: 39.994579},
	})
	return polygon
}
func getPolygon2() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.399669, Lat: 40.004307},
		{Lng: 116.360575, Lat: 39.952114},
		{Lng: 116.281812, Lat: 39.954326},
		{Lng: 116.3623, Lat: 39.916706},
		{Lng: 116.309983, Lat: 39.863559},
		{Lng: 116.401969, Lat: 39.892352},
		{Lng: 116.503729, Lat: 39.861344},
		{Lng: 116.469234, Lat: 39.929101},
		{Lng: 116.529025, Lat: 39.978215},
		{Lng: 116.440488, Lat: 39.956981},
	})
	return polygon
}

//爆炸形状
func getPolygon3() GeoPolygon {
	p := MakeGeoPoint(39.923664, 116.403424)
	var points []GeoPoint
	stp := 20
	num := 360 / stp
	for i := 0; i < num; i++ {
		dist := 4000
		if i%2 == 0 {
			dist = 8000
		}
		angle := float64(i * stp)
		p0 := PointAtDistAndAngle(p, float64(dist), angle)
		points = append(points, p0)
	}
	return MakeGeoPolygon(points)
}
func getPolygon4() GeoPolygon {
	//{39.869384765625 116.279296875 39.9957275390625 116.455078125}
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lat: 39.869384765625, Lng: 116.279296875},
		{Lat: 39.9957275390625, Lng: 116.279296875},
		{Lat: 39.9957275390625, Lng: 116.455078125},
		{Lat: 39.869384765625, Lng: 116.455078125},
	})
	return polygon

}
func getPolygon5() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.315013, Lat: 39.969147},
		{Lng: 116.340453, Lat: 39.993584},
		{Lng: 116.364456, Lat: 39.96771},
		{Lng: 116.37193, Lat: 39.967488},
		{Lng: 116.38429, Lat: 39.994358},
		{Lng: 116.411024, Lat: 39.964724},
		{Lng: 116.423241, Lat: 39.994247},
		{Lng: 116.463916, Lat: 39.966935},
		{Lng: 116.423816, Lat: 39.940387},
		{Lng: 116.441926, Lat: 39.91405},
		{Lng: 116.399382, Lat: 39.89202},
		{Lng: 116.350514, Lat: 39.912832},
		{Lng: 116.300209, Lat: 39.939281},
		{Lng: 116.342753, Lat: 39.955654},
	})
	return polygon
}
func getPolygon6() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.314438, Lat: 39.968926},
		{Lng: 116.331829, Lat: 39.991594},
		{Lng: 116.350227, Lat: 39.963949},
		{Lng: 116.369343, Lat: 39.992921},
		{Lng: 116.386159, Lat: 39.96406},
		{Lng: 116.42281, Lat: 39.9948},
		{Lng: 116.463485, Lat: 39.966493},
		{Lng: 116.496543, Lat: 39.95134},
		{Lng: 116.442069, Lat: 39.929543},
		{Lng: 116.485332, Lat: 39.913718},
		{Lng: 116.448681, Lat: 39.878068},
		{Lng: 116.425541, Lat: 39.91416},
		{Lng: 116.414617, Lat: 39.846499},
		{Lng: 116.390327, Lat: 39.896338},
		{Lng: 116.360144, Lat: 39.846499},
		{Lng: 116.33456, Lat: 39.886705},
		{Lng: 116.28713, Lat: 39.854808},
		{Lng: 116.319756, Lat: 39.90298},
		{Lng: 116.281668, Lat: 39.931093},
		{Lng: 116.331254, Lat: 39.952778},
		{Lng: 116.277356, Lat: 39.976446},
	})
	return polygon
}
func getPolygon7() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.30797, Lat: 39.991926},
		{Lng: 116.345627, Lat: 40.006739},
		{Lng: 116.385297, Lat: 39.993695},
		{Lng: 116.426978, Lat: 40.020664},
		{Lng: 116.451124, Lat: 39.993252},
		{Lng: 116.498267, Lat: 39.959857},
		{Lng: 116.467797, Lat: 39.952114},
		{Lng: 116.439051, Lat: 39.975562},
		{Lng: 116.32838, Lat: 39.972907},
		{Lng: 116.315444, Lat: 39.968041},
		{Lng: 116.319181, Lat: 39.873749},
		{Lng: 116.353964, Lat: 39.854919},
		{Lng: 116.462335, Lat: 39.864888},
		{Lng: 116.485907, Lat: 39.846721},
		{Lng: 116.46291, Lat: 39.801281},
		{Lng: 116.408006, Lat: 39.838078},
		{Lng: 116.349652, Lat: 39.78487},
		{Lng: 116.299634, Lat: 39.836527},
		{Lng: 116.234956, Lat: 39.911836},
		{Lng: 116.302509, Lat: 39.939281},
	})
	return polygon
}
func getPolygon8() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.328667, Lat: 39.972907},
		{Lng: 116.362012, Lat: 39.949238},
		{Lng: 116.441063, Lat: 39.947246},
		{Lng: 116.457161, Lat: 39.970475},
		{Lng: 116.465785, Lat: 39.874413},
		{Lng: 116.436177, Lat: 39.910508},
		{Lng: 116.364887, Lat: 39.906301},
		{Lng: 116.322056, Lat: 39.873749},
		{Lng: 116.34534, Lat: 39.930871},
	})
	return polygon
}
func getPolygon9() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.403685, Lat: 39.909262},
		{Lng: 116.40461, Lat: 39.909255},
		{Lng: 116.40461, Lat: 39.908543},
		{Lng: 116.403676, Lat: 39.908543},
	})
	return polygon
}

func getPolygon10() GeoPolygon {
	polygon := MakeGeoPolygon([]GeoPoint{
		{Lng: 116.363126, Lat: 39.913468},
		{Lng: 116.363162, Lat: 39.912777},
		{Lng: 116.442465, Lat: 39.914188},
		{Lng: 116.442213, Lat: 39.915046},
	})
	return polygon
}

//将多边形及切格子后的画在地图上
func splitGeoHashRect(polygon GeoPolygon, htmlName string, level int) {
	htmlStr := `<html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8" /><title>切格子效果观察</title>
			<script type="text/javascript" src="http://api.map.baidu.com/api?v=1.2"></script><script type="text/javascript" src="http://api.map.baidu.com/library/GeoUtils/1.2/src/GeoUtils_min.js"></script>
			</head><body>
				<div style="width:100%%;height:100%%;border:1px solid gray" id="container_%s"></div>
			</body></html>
			<script type="text/javascript">var map_%s = new BMap.Map("container_%s");window.stdMapCtrl = new BMap.NavigationControl();map_%s.addControl(window.stdMapCtrl);window.scaleCtrl = new BMap.ScaleControl();map_%s.addControl(window.scaleCtrl);window.overviewCtrl = new BMap.OverviewMapControl();map_%s.addControl(window.overviewCtrl);map_%s.addControl(new BMap.CopyrightControl());`
	htmlStr = fmt.Sprintf(htmlStr, htmlName, htmlName, htmlName, htmlName, htmlName, htmlName, htmlName)
	cu, _ := user.Current()
	outHtmlFile := fmt.Sprintf("%s/%s.html", cu.HomeDir, htmlName)
	htmlFP, err1 := os.Create(outHtmlFile) //创建文件
	if err1 != nil {
		panic(err1)
	}
	polygonRect := polygon.GetBoundsRect()
	midPoint := polygonRect.MidPoint()
	htmlStr = fmt.Sprintf("%svar pt_%s = new BMap.Point(%v,%v);", htmlStr, htmlName, midPoint.Lng, midPoint.Lat)
	htmlStr += fmt.Sprintf(`var mkr_%s = new BMap.Marker(pt_%s);var ply_%s;map_%s.centerAndZoom(pt_%s, %d);map_%s.enableContinuousZoom();polygon1_%s();function polygon1_%s() {var pts = [];`, htmlName, htmlName, htmlName, htmlName, htmlName, level, htmlName, htmlName, htmlName)

	polygonPoints := polygon.GetPoints()
	for _, p := range polygonPoints {
		htmlStr = fmt.Sprintf("%spts.push(new BMap.Point(%v,%v));", htmlStr, p.Lng, p.Lat)
	}
	htmlStr += fmt.Sprintf(`ply_%s = new BMap.Polygon(pts);ply_%s.setStrokeColor("red");map_%s.addOverlay(ply_%s);`, htmlName, htmlName, htmlName, htmlName)
	st := time.Now().UnixNano()
	inGrids, pGrids := polygon.RaySplitGeoHashRect(6)
	fmt.Println(htmlName, pT(st, time.Now().UnixNano()))
	for i, grid := range inGrids {
		htmlStr = fmt.Sprintf("%svar pts%d = [];", htmlStr, i)
		rect := GeoHashDecode(grid)
		ps := rect.GetRectVertex()
		for _, p := range ps {
			htmlStr = fmt.Sprintf("%spts%d.push(new BMap.Point(%v,%v));", htmlStr, i, p.Lng, p.Lat)
		}
		htmlStr = fmt.Sprintf("%svar ply_%s_%d=new BMap.Polygon(pts%d);ply_%s_%d.setStrokeWeight('1');map_%s.addOverlay(ply_%s_%d);", htmlStr, htmlName, i, i, htmlName, i, htmlName, htmlName, i)
	}
	for i, grid := range pGrids {
		htmlStr = fmt.Sprintf("%svar pts%d = [];", htmlStr, i)
		rect := GeoHashDecode(grid)
		ps := rect.GetRectVertex()
		for _, p := range ps {
			htmlStr = fmt.Sprintf("%spts%d.push(new BMap.Point(%v,%v));", htmlStr, i, p.Lng, p.Lat)
		}
		htmlStr = fmt.Sprintf("%svar ply_%s_%d=new BMap.Polygon(pts%d);ply_%s_%d.setStrokeWeight('1');ply_%s_%d.setStrokeStyle('dashed');ply_%s_%d.setFillColor('#F0F8FF');map_%s.addOverlay(ply_%s_%d);", htmlStr, htmlName, i, i, htmlName, i, htmlName, i, htmlName, i, htmlName, htmlName, i)
	}
	htmlStr += `}</script>`
	_, err1 = io.WriteString(htmlFP, htmlStr)
	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println("outHtmlFile：", outHtmlFile)
}
