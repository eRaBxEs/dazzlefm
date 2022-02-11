package util

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/sessions"

	"bitbucket.org/mayowa/helpers/convert"
	"github.com/labstack/echo"

	"bitbucket.org/mayowa/helpers/dbtools"
	"golang.org/x/crypto/bcrypt"
)

// Pager ...
func Pager(page, limit int64) (offset, lmt int64) {
	if page <= 0 {
		page = 1
	}

	lmt = limit
	if limit <= 0 {
		lmt = 15
	}

	offset = (page - 1) * lmt
	return
}

// IntParam converts c.Param(p) or c.QueryParam(p) to int
func IntParam(c echo.Context, p string) int {
	val := c.Param(p)
	if len(val) > 0 {
		return convert.Atoi(val)
	}

	return convert.Atoi(c.QueryParam(p))

}

// Int64Param converts c.Param() to int64
func Int64Param(c echo.Context, p string) int64 {
	val := c.Param(p)
	if len(val) > 0 {
		return convert.Atoi64(val)
	}

	return convert.Atoi64(c.QueryParam(p))
}

// RecToPageCount converts as record count to page count based on the provided limit
func RecToPageCount(recCount, limit uint64) (pageCount uint64) {
	pageCount = recCount / limit
	if pageCount%limit != 0 {
		pageCount++
	}

	return
}

// ToDBUID ...
func ToDBUID(v string) dbtools.DBUID {
	retv, err := dbtools.NewDBUIDStr(v)
	if err != nil {
		return dbtools.NewDBUIDInt(0)
	}

	return retv
}

// GetFileName extract the filename from a path and sanitize it
func GetFileName(file string) string {
	fName := filepath.Base(filepath.Clean(file))

	fName = strings.Replace(fName, " ", "", -1)

	return fName
}

// HashPassword returns a bcrypt hash of the input string
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares a bcrypt hash with a string
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ParseISOTime converts YYYY-MM-DD HH:MM:SS to time.Time
func ParseISOTime(val string, end bool) time.Time {
	retv, err := time.Parse("2006-01-02 15:04:05", val)
	if err != nil {
		retv, err = time.Parse("2006-01-02", val)
		if err != nil {
			retv, err = time.Parse("20060102", val)
			if err != nil {
				return time.Time{}
			}
		}
	}

	if end {
		return time.Date(retv.Year(), retv.Month(), retv.Day(), 23, 59, 59, 0, time.Local)
	}

	return retv
}

// GetSessionValue gets an item from the session. performs check for existence
func GetSessionValue(s *sessions.Session, key string) (val interface{}, exists bool) {
	if s == nil {
		return nil, false
	}

	if _, exists = s.Values[key]; !exists {
		return nil, false
	}

	return s.Values[key], true
}

// GetSessionStrValue string typed version of GetSessionValue
func GetSessionStrValue(s *sessions.Session, key string) (val string, exists bool) {
	retv, exists := GetSessionValue(s, key)
	if !exists {
		return
	}

	return retv.(string), true
}

// GetSessionIntValue in64 typed version of GetSessionValue
func GetSessionIntValue(s *sessions.Session, key string) (val int, exists bool) {
	retv, exists := GetSessionValue(s, key)
	if !exists {
		return
	}

	return retv.(int), true
}

// GetSessionInt64Value in64 typed version of GetSessionValue
func GetSessionInt64Value(s *sessions.Session, key string) (val int64, exists bool) {
	retv, exists := GetSessionValue(s, key)
	if !exists {
		return
	}

	return retv.(int64), true
}

// ImageData ...
type ImageData struct {
	Data string `json:"data"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// SaveImageData save a data uri to a file and return a url where the file can be reached
func SaveImageData(env *Environment, imgData json.RawMessage) (image *ImageData, err error) {
	saver := &AssetSave{}
	saver.Init(env)

	asset := &Asset{}
	image = &ImageData{}

	// convert images in json
	if len(imgData) > 0 {
		err = json.Unmarshal(imgData, image)
		if err != nil {
			return
		}

		if !strings.HasPrefix(image.Data, "data:") {
			// not a data uri
			return
		}

		// extract data uri and save to disk
		asset, err = saver.Save(&image.Data)
		if err != nil {
			return
		}

		image.Data = asset.FileURL
	}

	return
}

// SaveAudioData save a data uri to a file and return a url where the file can be reached
func SaveAudioData(env *Environment, imgData json.RawMessage) (image *ImageData, err error) {
	saver := &AssetSave{}
	saver.Init(env)

	asset := &Asset{}
	image = &ImageData{}

	// convert images in json
	if len(imgData) > 0 {
		err = json.Unmarshal(imgData, image)
		if err != nil {
			return
		}

		if !strings.HasPrefix(image.Data, "data:audio") {
			// not a data uri
			return
		}

		// extract data uri and save to disk
		asset, err = saver.Save(&image.Data)
		if err != nil {
			return
		}

		image.Data = asset.FileURL
	}

	return
}
