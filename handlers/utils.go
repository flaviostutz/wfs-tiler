package handlers

import "math"

//BBox bounding box
type BBox struct {
	west  float64
	north float64
	east  float64
	south float64
}

func tile2BBox(x int, y int, z int) BBox {
	north, west := tile2latlon(x, y, z)
	south, east := tile2latlon(x+1, y+1, z)
	return BBox{west, north, east, south}
}

func tile2latlon(x int, y int, z int) (lat float64, lon float64) {
	n := math.Pi - (2.0*math.Pi*float64(y))/math.Pow(2.0, float64(z))
	latRad := math.Atan(math.Sinh(n))
	lat1 := latRad * (180 / math.Pi)
	lon1 := float64(x)/math.Pow(2.0, float64(z))*360.0 - 180.0
	return lat1, lon1
}
