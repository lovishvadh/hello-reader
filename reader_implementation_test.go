package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReaderWithFile(t *testing.T) {
	mockTxtFilePath := filepath.Join("mockdata", "mockdata.txt")
	reader, err := NewReader(mockTxtFilePath)
	require.NoError(t, err)

	p1 := make([]byte, 100)
	n1, err1 := reader.Read(p1)
	require.NoError(t, err1)

	assert.Contains(t, string(p1[:n1]), "Hello there if you are seeing this data your function works! ;)")

	mockJsonFilePath := filepath.Join("mockdata", "mockdata.json")
	reader, err = NewReader(mockJsonFilePath)
	require.NoError(t, err)

	p2 := make([]byte, 100)
	n2, err2 := reader.Read(p2)
	require.NoError(t, err2)

	assert.Contains(t, string(p2[:n2]), "Lovish")

}

func TestNewReaderWithUrl(t *testing.T) {
	mockJsonURLPath := "https://jsonplaceholder.typicode.com/todos"
	reader, err := NewReader(mockJsonURLPath)
	require.NoError(t, err)

	p1 := make([]byte, 100)
	n1, err1 := reader.Read(p1)
	require.NoError(t, err1)

	assert.Contains(t, string(p1[:n1]), "delectus aut autem")
}

func TestNewReaderWithInvalidUrl(t *testing.T) {
	mockJsonURLPath := "https://jsonplaceholder.typiiiicode.com/todos"
	_, err := NewReader(mockJsonURLPath)
	require.Error(t, err)
}

func TestNewReaderWithInvalidFile(t *testing.T) {
	mockTxtFilePath := filepath.Join("mockdata", "sample.txt")
	_, err := NewReader(mockTxtFilePath)
	require.Error(t, err)
}

func TestNewReaderWithEmptyFile(t *testing.T) {
	mockTxtFilePath := filepath.Join("mockdata", "empty.txt")
	_, err := NewReader(mockTxtFilePath)
	require.NoError(t, err)
}

func TestNewReaderWithEmptyUrl(t *testing.T) {
	mockJsonURLPath := ""
	_, err := NewReader(mockJsonURLPath)
	require.Error(t, err)
}

func TestNewReaderWithEmptyFilePath(t *testing.T) {
	mockJsonURLPath := "https://httpbin.org/gzip"
	reader, err := NewReader(mockJsonURLPath)
	require.NoError(t, err)

	p1 := make([]byte, 100)
	n1, err1 := reader.Read(p1)
	require.NoError(t, err1)

	assert.Contains(t, string(p1[:n1]), "gzipped")
}
