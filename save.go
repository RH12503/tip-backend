package main

import (
	"bufio"
	"encoding/binary"
	"github.com/RH12503/Triangula/geom"
	"github.com/RH12503/Triangula/image"
	"github.com/RH12503/Triangula/normgeom"
	"github.com/RH12503/Triangula/render"
	"github.com/RH12503/Triangula/triangulation"
	"math"
	"os"
)

func SaveFile(filepath string, points normgeom.NormPointGroup, image image.Data) error {
	w, h := image.Size()

	triangles := triangulation.Triangulate(points, w, h)
	triangleData := render.TrianglesOnImage(triangles, image)

	file, err := os.Create(filepath)
	defer file.Close()

	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)

	writer.Write(uint16ToBytes(uint16(w)))
	writer.Write(uint16ToBytes(uint16(h)))

	writer.Write(uint16ToBytes(uint16(len(points))))

	pointsMap := make(map[geom.Point]int)

	for i, p := range points {
		point := geom.Point{
			X: multAndRound(p.X, w),
			Y: multAndRound(p.Y, h),
		}
		writer.Write(uint16ToBytes(uint16(point.X)))
		writer.Write(uint16ToBytes(uint16(point.Y)))

		pointsMap[point] = i
	}

	for _, d := range triangleData {
		tri := d.Triangle.Points
		col := d.Color

		for _, p := range tri {
			point := geom.Point{
				X: multAndRound(p.X, w),
				Y: multAndRound(p.Y, h),
			}

			writer.Write(uint16ToBytes(uint16(pointsMap[point])))
		}

		writer.Write([]byte{uint8(multAndRound(col.R, 255))})
		writer.Write([]byte{uint8(multAndRound(col.G, 255))})
		writer.Write([]byte{uint8(multAndRound(col.B, 255))})
	}

	writer.Flush()

	return nil
}

func uint16ToBytes(num uint16) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, num)
	return bytes
}

func multAndRound(v float64, m int) int {
	return int(math.Round(v * float64(m)))
}
