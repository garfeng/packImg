package compress

import (
	"gocv.io/x/gocv"
	"image"
)

func Region(src gocv.Mat, roi image.Rectangle) gocv.Mat {
	return CloneRegion(src, roi)
}

func CloneRegion(src gocv.Mat, roi image.Rectangle) gocv.Mat {
	if roi.Dx() <= 0 || roi.Dy() <= 0 {
		return gocv.NewMat()
	}

	srcROI := image.Rect(0, 0, src.Cols(), src.Rows())

	if roi.In(srcROI) {
		region := src.Region(roi)
		defer region.Close()
		return region.Clone()
	}

	dst := gocv.NewMatWithSize(roi.Dy(), roi.Dx(), src.Type())
	dst.SetTo(gocv.Scalar{})
	realROIOfSrc := srcROI.Intersect(roi)

	// roi 完全超出了原始图像区域，直接返回黑图
	if realROIOfSrc.Dx() <= 0 || realROIOfSrc.Dy() <= 0 {
		return dst
	}

	realROIOfDst := realROIOfSrc.Sub(roi.Min)

	srcRegion := src.Region(realROIOfSrc)
	defer srcRegion.Close()

	dstRegion := dst.Region(realROIOfDst)
	defer dstRegion.Close()

	srcRegion.CopyTo(&dstRegion)
	return dst
}
