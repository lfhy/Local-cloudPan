package static_test

import (
	"fmt"
	"local-cloud-api/static"
	"testing"
)

func TestGetFile(t *testing.T) {
	dirs, err := static.Fs.ReadDir("dist")
	if err != nil {
		t.Failed()
		return
	}
	for _, dir := range dirs {
		fmt.Printf("dir: %v\n", dir)
	}
}
