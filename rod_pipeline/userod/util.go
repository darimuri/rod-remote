package userod

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ctessum/geom"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func GetImageAreaCentroid(img, area *rod.Element) (*proto.Point, error) {
	quads, errShape := img.Shape()
	if errShape != nil {
		return nil, errShape
	}

	coords, errCoords := parseCoords(area)
	if errCoords != nil {
		return nil, errCoords
	}

	points := make([]geom.Point, 0)
	for i := 0; i < len(coords)-1; i += 2 {
		points = append(points, geom.Point{X: float64(coords[i]), Y: float64(coords[i+1])})
	}

	polygon := geom.Polygon{points}
	c := polygon.Centroid()

	x := quads.Quads[0][0] + c.X
	y := quads.Quads[0][1] + c.Y

	center := proto.NewPoint(x, y)

	return &center, nil
}

func parseCoords(el *rod.Element) ([]int, error) {
	attr, err := el.Attribute("coords")
	if err != nil {
		return nil, err
	} else if attr == nil {
		return nil, errors.New("nil coords attribute value")
	}

	vals := strings.Split(*attr, ",")

	if len(vals) < 8 {
		return nil, errors.New("coords attribute value does not have enough element(8)")
	}

	if len(vals)%2 != 0 {
		return nil, errors.New(fmt.Sprintf("number of coords(%d) is not even", len(vals)))
	}

	coords := make([]int, len(vals))
	for idx, v := range vals {
		n, errConv := strconv.Atoi(v)
		if errConv != nil {
			return nil, errConv
		}

		coords[idx] = n
	}

	return coords, nil
}
