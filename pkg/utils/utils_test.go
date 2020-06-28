package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExportCsv(t *testing.T) {

	finalPath, err := ExportCsv("../export/test.csv", []string{"col1", "col2"}, [][]string{{"Line1", "Hello Readers of"}, {"Line2", "golangcode.com"}})
	assert.Nil(t, err, "should not be error")
	assert.NotEmpty(t, finalPath, "should not empty")
}
