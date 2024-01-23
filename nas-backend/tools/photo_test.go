package tools

import "testing"

func TestPhotoCropping(t *testing.T) {
	PhotoCropping("IMG_5652.JPEG", 690, 967, 971, 1315)
}
