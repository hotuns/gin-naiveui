package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	svg "github.com/ajstarks/svgo"
)

func GenerateSVG(width, height int) ([]byte, string) {
	rand.Seed(time.Now().UnixNano())

	var svgContent bytes.Buffer
	canvas := svg.New(&svgContent)
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:white")

	code := fmt.Sprintf("%04d", rand.Intn(10000)) // 确保生成4位数的随机数
	canvas.Text(width/2, height/2, code, "text-anchor:middle; dominant-baseline:middle; font-size:30px; fill:black;")

	canvas.End()

	return svgContent.Bytes(), code
}
