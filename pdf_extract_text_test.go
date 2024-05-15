/*
 * PDF to text: Extract all text for each page of a pdf file.
 *
 * Run as: go run pdf_extract_text.go input.pdf
 */

package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"

	common "github.com/lxt1045/PdfTextExtract/common"
	"github.com/lxt1045/PdfTextExtract/extractor"
	// "github.com/otiai10/gosseract"
)

func TestExtract(t *testing.T) {
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	// inputPath := os.Args[1]
	inputPath := "D:/book/天涯/kk预测.pdf"

	m := new(runtime.MemStats)
	runtime.GC()
	runtime.ReadMemStats(m)

	text, err := extractor.ExtractPdfFile(inputPath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(text)

	fmt.Printf("before load, heap memory: %d, head in use: %d\n", m.HeapAlloc, m.HeapInuse)
	runtime.GC()
	runtime.ReadMemStats(m)
	fmt.Printf("after load, heap memory: %d, head in use: %d\n", m.HeapAlloc, m.HeapInuse)

}

func TestExtractByContent(t *testing.T) {
	common.SetLogger(common.NewConsoleLogger(common.LogLevelDebug))

	// inputPath := os.Args[1]
	inputPath := "D:/book/天涯/kk预测.pdf"

	m := new(runtime.MemStats)
	runtime.GC()
	runtime.ReadMemStats(m)

	fi, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := io.ReadAll(fi)
	if err != nil {
		panic(err)
	}
	content := string(fd)
	text, err := extractor.ExtractPdfContent(content)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(text)

	fmt.Printf("before load, heap memory: %d, head in use: %d\n", m.HeapAlloc, m.HeapInuse)
	runtime.GC()
	runtime.ReadMemStats(m)
	fmt.Printf("after load, heap memory: %d, head in use: %d\n", m.HeapAlloc, m.HeapInuse)

}
