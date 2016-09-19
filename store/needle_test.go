package store

import "testing"

func TestGetPadding(t *testing.T) {
	fixSize := totalHeaderLen + totalFooterLen
	dataSize := 12132
	paddingSize := len(getNeedlePadding(uint32(dataSize)))
	total := dataSize + fixSize + paddingSize
	if total%align != 0 {
		t.Errorf("total is %d", total)
	}
}
