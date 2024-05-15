/*
 * PDF to text: Extract all text for each page of a pdf file.
 *
 * Run as: go run pdf_extract_text.go input.pdf
 */

package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	common "github.com/lxt1045/PdfTextExtract/common"
	. "github.com/lxt1045/PdfTextExtract/core"
	. "github.com/lxt1045/PdfTextExtract/extractor"
	pdf "github.com/lxt1045/PdfTextExtract/model"

	// "github.com/otiai10/gosseract"

	//"runtime"
	"strings"
)

type ContentPair struct {
	s     *PdfObjectStream
	index int
}

func parseText(this *pdf.PdfReader) (string, error) {
	pageList := this.GetPageList()
	parser := this.GetParser()
	mFontsForPages := this.GetFontsForPages()
	// mPageResources := this.GetPageResources()

	contentStreamChan := make(chan ContentPair, 10)

	// client := gosseract.NewClient()
	// client.SetLanguage("chi_sim", "eng")
	// defer client.Close()

	go func() {
		for i := 0; i < len(pageList); i++ {
			if pageObjDict, ok := pageList[i].PdfObject.(*PdfObjectDictionary); ok {
				if contentsArray, ok := pageObjDict.Get("Contents").(*PdfObjectArray); ok {
					for j := 0; j < len(*contentsArray); j++ {
						contentObj, err := parser.Trace((*contentsArray)[j])
						if err != nil {
							common.Log.Debug("Error: trace content to obj failed, err: %s", err)
							continue
						}
						if contentStmObj, ok := contentObj.(*PdfObjectStream); ok {
							produce := true
							for produce {
								select {
								case contentStreamChan <- ContentPair{contentStmObj, i}:
									produce = false
								default:
									time.Sleep(2 * time.Millisecond)
								}
							}
						}
					}
				} else if contentObj, err := parser.Trace(pageObjDict.Get("Contents")); err == nil {
					if contentStmObj, ok := contentObj.(*PdfObjectStream); ok {
						produce := true
						for produce {
							select {
							case contentStreamChan <- ContentPair{contentStmObj, i}:
								produce = false
							default:
								time.Sleep(2 * time.Millisecond)
							}
						}
					}
				}
			}

			// //switch off
			// if false {
			// 	if xObjectObjDict, ok := mPageResources[i].Get("XObject").(*PdfObjectDictionary); ok {
			// 		common.Log.Debug("xobject %s", xObjectObjDict)
			// 		for imgName, _ := range xObjectObjDict.Dict() {
			// 			common.Log.Debug("picture: %s try to ocr", imgName)
			// 			if imageObj, err := parser.Trace(xObjectObjDict.Get(imgName)); err == nil {
			// 				if imageObjStm, ok := imageObj.(*PdfObjectStream); ok {
			// 					client.SetImageFromBytes(imageObjStm.Stream)
			// 					text, _ := client.Text()
			// 					common.Log.Debug("image text: %s", text)
			// 				}
			// 			}
			// 		}
			// 	}
			// }
		}
		close(contentStreamChan)
	}()

	var textBuffer bytes.Buffer
	for {
		if pair, ok := <-contentStreamChan; ok {
			streamData, err := DecodeStream(pair.s)
			if err != nil {
				return "", err
			}

			common.Log.Trace("stream data: %s", streamData)

			e := New(string(streamData), mFontsForPages[pair.index])
			s, _ := e.ExtractText()
			textBuffer.WriteString(s)
			textBuffer.WriteString("\n\n")
		} else {
			break
		}
	}

	return textBuffer.String(), nil
}

// outputPdfText prints out contents of PDF file to stdout.
func ExtractPdfContent(content string) (string, error) {

	f := strings.NewReader(content)

	pdfReader, err := pdf.NewPdfReader(f)

	if err != nil {
		fmt.Printf("parser pdf failed, err: %s\n", err)
		return "", err
	}

	err = pdfReader.ParseFonts()
	if err != nil {
		fmt.Printf("parse fonts err: %s\n", err)
		return "", err
	}

	text, err := parseText(pdfReader)

	return text, err
}

// outputPdfText prints out contents of PDF file to stdout.
func ExtractPdfFile(inputPath string) (string, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return "", err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)

	if err != nil {
		fmt.Printf("parser pdf failed, err: %s\n", err)
		return "", err
	}

	err = pdfReader.ParseFonts()
	if err != nil {
		fmt.Printf("parse fonts err: %s\n", err)
		return "", err
	}

	text, err := parseText(pdfReader)

	return text, err
}
