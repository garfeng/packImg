package main

import (
	"fmt"
	"github.com/garfeng/packImg/compress"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		return
	}

	root := os.Args[1]

	fmt.Println("请输入压缩质量，（1~100）：")
	quality := 80
	fmt.Scanf("%d\n", &quality)
	if quality < 1 || quality > 100 {
		quality = 80
	}

	err := compress.ScanBmpAndCompress(root, float32(quality))
	if err != nil {
		fmt.Println(err)
		pause()
		return
	}

	fmt.Println("完成，窗口将在5秒后关闭")
	<-time.After(time.Second * 5)
}

func pause() {
	fmt.Println("请手动关闭窗口")
	a := ""
	fmt.Scanf("%s", &a)
}
