# goutils
另一个家：https://liuyongshuai.com/article/1518878567244959744
断断续续开发的一些golang的工具类，主要包括：

## color
在终端上显示彩色字体、闪烁效果、下划线效果的小工具（但在MAC上并不能展现出闪烁效果）。

## elem
golang里基本类型数据的转换操作。试图传入一个任意类型的数据，然后提供一些判断、转换成任意类型的数据。详见相关的测试代码。
试图在诸如提取数据库字段、获取请求参数时只返回这样的类型，然后再做自由的类型转换。

## file
逐行迭代文件的小工具，还有一个操作文件的常用工具。

## helper
杂七杂八的函数，包括IP转换、截字符串、全角/半角之间转换、繁体字转简体字等

## http
自己封装的一个发起http请求的库，主要是自用。

## geo
跟地理位置相关的一些操作，如：
* geohash编解码（包括返回字符串格式的及int64格式的）
* 用一个点去构造另一个点（转换一定角度）
* 点是否在多边形内部（射线法）
* 两点间的距离
* 两直线是否相交
* 点是否在直线上
* 用geohash的小格子去切割多边形，返回完全在多边形内部的小格子及部分区域重合多边形的小格子，测试用例里生成了好些html文件，是将所有的小格子及多边形在百度地图上画出来以便直观看结果。

## mysql
自己封装的请求mysql等操作的库，主要是自用。

## slice
封装了对slice类型的常用操作。

# snowflake
这是对SnowFlake算法的一个改进版，在原算法基础上提供了对各域的位数的自定义设置。
如在一些场景下，dateCenterId、workerId并不需要占太多的位数，反而sequence部分需要占较多的位数。此时就可以根据具体业务场景自己设置。
并支持在多线程访问时的安全问题。自己测试时连续生成了大约一亿多的ID，并没有发现重复的。详见相关的测试代码。
