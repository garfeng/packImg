package compress

import (
	"archive/zip"
	"fmt"
	"gocv.io/x/gocv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*
func CompressByWebp(src gocv.Mat, quality float32) ([]byte, error) {

	img, err := src.ToImage()
	if err != nil {
		return nil, err
	}
	return webp.EncodeGray(img, quality)
}

*/

func ReadImage(name string, flag gocv.IMReadFlag) (img gocv.Mat, size int, err error) {
	buff, err := os.ReadFile(name)
	if err != nil {
		return gocv.Mat{}, 0, err
	}
	img, err = gocv.IMDecode(buff, flag)
	if err != nil {
		return gocv.Mat{}, len(buff), err
	}
	return img, len(buff), nil
}

func CompressOne(name string, w io.Writer, quality float32) (float64, error) {
	img, size, err := ReadImage(name, gocv.IMReadUnchanged)

	if err != nil {
		return 0, err
	}
	defer img.Close()

	buff, err := gocv.IMEncodeWithParams(".webp", img, []int{
		gocv.IMWriteWebpQuality, int(quality),
	})
	if err != nil {
		return 0, err
	}
	defer buff.Close()
	size2 := buff.Len()

	_, err = w.Write(buff.GetBytes())
	return float64(size2) / float64(size), err
}

func CompressAGroup(bmpFiles []string, root string, prefix string, zipId int, quality float32, startId, endId int) error {
	dstZipFile := fmt.Sprintf("%s.%d.zip", prefix, zipId)

	os.Remove(dstZipFile)
	f, err := os.Create(dstZipFile)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	_, rootName := filepath.Split(root)
	rootIdx := len(root)

	num := len(bmpFiles)

	if startId >= num {
		startId = num - 1
	}
	if endId > num {
		endId = num
	}

	for i := startId; i < endId; i++ {
		name := bmpFiles[i]

		relativeName := filepath.Join(rootName, name[rootIdx:])
		dstRelativeName := relativeName + ".webp" //ReplaceExtTo(relativeName, ".webp")

		oneWriter, err := w.Create(dstRelativeName)
		if err != nil {
			log.Println("[Err]", name, err)
			continue
		}
		percent, err := CompressOne(name, oneWriter, quality)
		if err != nil {
			log.Println("[Err]", name, err)
			continue
		}
		ss := fmt.Sprintf("[%d/%d] [%0.1f%%]", i+1, num, percent*100)
		log.Println("[OK]", ss, name)
	}

	return nil
}

func ScanBmpAndCompress(root string, quality float32) error {
	bmpFiles := scanSubRoots(root, ".bmp", ".png", ".jpg", ".jpeg", ".tif")

	root = strings.Trim(root, "/")
	root = strings.Trim(root, "\\")
	dstZipFilePrefix := root + ".webp"

	zipFileId := 0

	num := len(bmpFiles)
	batchNumber := 100

	for i := 0; i < num; i += batchNumber {
		zipFileId += 1
		CompressAGroup(bmpFiles, root, dstZipFilePrefix, zipFileId, quality, i, i+batchNumber)
	}

	return nil
}
