package compress

import (
	"testing"
)

func TestScanBmpAndCompress(t *testing.T) {
	ScanBmpAndCompress("../images", 90)
}
