package gin_mock_case

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	//nolint:golint,unused
	maxUploadLength = 1024 * 1024 * 1024 * 2
	//nolint:golint,unused
	defaultUploadReadSize = 1024 * 12
	//nolint:golint,unused
	defaultUploadBoundarySize = 1024 * 8
	//nolint:golint,unused
	defaultUploadBufferSize = 1024 * 4

	BOUNDARY                 = "; boundary="
	UploadContentDisposition = "Content-Disposition: "
	UploadNAME               = "name=\""
	UploadFILENAME           = "filename=\""
	ContentType              = "Content-Type: "
	ContentLength            = "Content-Length: "
)

var (
	//nolint:golint,unused
	uploadBoundaryIndex = []byte("\r\n\r\n")
	//nolint:golint,unused
	uploadBoundaryLine = []byte("\r\n")
)

// pathExists
//
//	path exists
//
//nolint:golint,unused
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// pathExistsFast
//
//	path exists fast
//
//nolint:golint,unused
func pathExistsFast(path string) bool {
	exists, _ := pathExists(path)
	return exists
}

//nolint:golint,unused
func mkdir(folderPath string) (string, error) {
	if pathExistsFast(folderPath) {
		return folderPath, nil
	}
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	return folderPath, nil
}

//nolint:golint,unused
func saveFileHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		jsonErr(c, err)
		return
	}
	uploadFileName := c.Request.FormValue("upload_name")
	if len(uploadFileName) == 0 {
		jsonErr(c, nil, "not found form upload_name")
		return
	}
	fileName := c.Request.FormValue("file_name")
	if len(fileName) == 0 {
		jsonErr(c, nil, "not found form file_name")
		return
	}

	files := form.File[uploadFileName]
	if len(files) == 0 {
		jsonErr(c, nil, "not found file")
		return
	}
	for _, fileItem := range files {
		fileNameInt := time.Now().Unix()
		fileNameStr := strconv.FormatInt(fileNameInt, 10)
		fileSaveName := fileNameStr + "-" + fileItem.Filename + ".tmp"

		uploadPath, errSavePath := mkdir("upload")
		if errSavePath != nil {
			jsonErr(c, errSavePath)
			return
		}
		fileSavePath := filepath.Join(uploadPath, fileSaveName)
		errSave := c.SaveUploadedFile(fileItem, fileSavePath)
		if errSave != nil {
			jsonErr(c, errSave)
			return
		}
		jsonSuccess(c, &FileHeader{
			Name:          fileItem.Filename,
			ContentType:   fileItem.Header.Get("Content-Type"),
			ContentLength: fileItem.Size,
		})
	}
}

//nolint:golint,unused
func saveMaxFileHandler(c *gin.Context) {
	contentLength := c.Request.ContentLength
	if contentLength <= 0 || contentLength > maxUploadLength {
		jsonErr(c, nil, "ContentLength greater max or less than 0")
		return
	}
	contentTypeList, hasContentTypeKey := c.Request.Header["Content-Type"]
	if !hasContentTypeKey {
		jsonErr(c, nil, "not found head Content-Type")
		return
	}
	if len(contentTypeList) == 0 {
		jsonErr(c, nil, "head Content-Type is empty")
		return
	}
	contentType := contentTypeList[0]

	locContentType := strings.Index(contentType, BOUNDARY)
	if locContentType == -1 {
		jsonErr(c, nil, "Content-Type error, no boundary")
		return
	}
	boundary := []byte(contentType[(locContentType + len(BOUNDARY)):])

	readData := make([]byte, defaultUploadReadSize)
	var readTotal = 0
	for {
		fileHeader, fileData, err := parseFromHead(readData, readTotal, append(boundary, uploadBoundaryLine...), c.Request.Body)
		if err != nil {
			jsonErr(c, err)
			return
		}
		f, errFileCreate := os.Create(fileHeader.FileName)
		if errFileCreate != nil {
			jsonErr(c, errFileCreate)
			return
		}
		_, errWrite := f.Write(fileData)
		if errWrite != nil {
			jsonErr(c, errWrite)
			return
		}

		//fileData = nil

		tempData, reachEnd, err := readToBoundary(boundary, c.Request.Body, f)
		errClose := f.Close()
		if errClose != nil {
			jsonErr(c, errClose)
			return
		}
		if err != nil {
			jsonErr(c, err)
			return
		}
		if reachEnd {
			break
		} else {
			copy(readData[0:], tempData)
			readTotal = len(tempData)
			continue
		}
	}
}

//nolint:golint,unused
type FileHeader struct {
	ContentDisposition string
	Name               string
	FileName           string
	ContentType        string
	ContentLength      int64
}

//nolint:golint,unused
func parseFromHead(readData []byte, readTotal int, boundary []byte, stream io.ReadCloser) (FileHeader, []byte, error) {
	buf := make([]byte, defaultUploadBufferSize)
	foundBoundary := false
	boundaryLoc := -1
	var fileHeader FileHeader
	for {
		readLen, err := stream.Read(buf)
		if err != nil {
			if err != io.EOF {
				return fileHeader, nil, err
			}
			break
		}
		if readTotal+readLen > cap(readData) {
			return fileHeader, nil, fmt.Errorf("not found boundary")
		}
		copy(readData[readTotal:], buf[:readLen])
		readTotal += readLen
		if !foundBoundary {
			boundaryLoc = bytes.Index(readData[:readTotal], boundary)
			if -1 == boundaryLoc {
				continue
			}
			foundBoundary = true
		}
		startLoc := boundaryLoc + len(boundary)
		fileHeadLoc := bytes.Index(readData[startLoc:readTotal], uploadBoundaryIndex)
		if -1 == fileHeadLoc {
			continue
		}
		fileHeadLoc += startLoc
		ret := false
		fileHeader, ret = parseUploadFileHeader(readData[startLoc:fileHeadLoc])
		if !ret {
			return fileHeader, nil, fmt.Errorf("parseFileHeader fail:%s", string(readData[startLoc:fileHeadLoc]))
		}
		return fileHeader, readData[fileHeadLoc+4 : readTotal], nil
	}
	return fileHeader, nil, fmt.Errorf("reach to sream EOF")
}

//nolint:golint,unused
func parseUploadFileHeader(h []byte) (FileHeader, bool) {
	arr := bytes.Split(h, []byte("\r\n"))
	var outHeader FileHeader
	outHeader.ContentLength = -1
	for _, item := range arr {
		if bytes.HasPrefix(item, []byte(UploadContentDisposition)) {
			l := len(UploadContentDisposition)
			arr1 := bytes.Split(item[l:], []byte("; "))
			outHeader.ContentDisposition = string(arr1[0])
			if bytes.HasPrefix(arr1[1], []byte(UploadNAME)) {
				outHeader.Name = string(arr1[1][len(UploadNAME) : len(arr1[1])-1])
			}
			l = len(arr1[2])
			if bytes.HasPrefix(arr1[2], []byte(UploadFILENAME)) && arr1[2][l-1] == 0x22 {
				outHeader.FileName = string(arr1[2][len(UploadFILENAME) : l-1])
			}
		} else if bytes.HasPrefix(item, []byte(ContentType)) {
			l := len(ContentType)
			outHeader.ContentType = string(item[l:])
		} else if bytes.HasPrefix(item, []byte(ContentLength)) {
			l := len(ContentLength)
			s := string(item[l:])
			contentLength, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				//log.Printf("content length error:%s", string(item))
				return outHeader, false
			} else {
				outHeader.ContentLength = contentLength
			}
		}
		//else {
		// log.Printf("unknown:%s\n", string(item))
		//}
	}
	if len(outHeader.FileName) == 0 {
		return outHeader, false
	}
	return outHeader, true
}

//nolint:golint,unused
func readToBoundary(boundary []byte, stream io.ReadCloser, target io.WriteCloser) ([]byte, bool, error) {
	readData := make([]byte, defaultUploadBoundarySize)
	readDataLen := 0
	buf := make([]byte, defaultUploadBufferSize)
	bLen := len(boundary)
	reachEnd := false
	for !reachEnd {
		readLen, err := stream.Read(buf)
		if err != nil {
			if err != io.EOF && readLen <= 0 {
				return nil, true, err
			}
			reachEnd = true
		}
		// The following sentence is stupid and worth optimizing.
		copy(readData[readDataLen:], buf[:readLen]) // Append to another buffer, just for search convenience.
		readDataLen += readLen
		if readDataLen < bLen+4 {
			continue
		}
		loc := bytes.Index(readData[:readDataLen], boundary)
		if loc >= 0 {
			//找到了结束位置
			_, errWrite := target.Write(readData[:loc-4])
			if errWrite != nil {
				return nil, reachEnd, errWrite
			}
			return readData[loc:readDataLen], reachEnd, nil
		}

		_, errWrite := target.Write(readData[:readDataLen-bLen-4])
		if errWrite != nil {
			return nil, reachEnd, errWrite
		}
		copy(readData[0:], readData[readDataLen-bLen-4:])
		readDataLen = bLen + 4
	}
	_, errWrite := target.Write(readData[:readDataLen])
	if errWrite != nil {
		return nil, reachEnd, errWrite
	}
	return nil, reachEnd, nil
}
