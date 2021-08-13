package main

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	strings "strings"
	"sync"
)

var mutex sync.Mutex

var m = make(map[string]*ShortURLObj)

var newID uint64
var maxID = uint64(math.Pow(62, 8))

var base62 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

type URLObj interface {
	visitShortURL()
	visitLongURL()
	equalID(id uint64) bool
	getID() uint64
	getShortURL() string
	getVisitShortURLCnt() uint64
	getVisitLongURLCnt() uint64
}

type ShortURLObj struct {
	ID               uint64
	ShortURL         string
	visitShortURLCnt uint64
	visitLongURLCnt  uint64
}

func (obj *ShortURLObj) visitShortURL() {
	obj.visitShortURLCnt++
}

func (obj *ShortURLObj) visitLongURL() {
	obj.visitLongURLCnt++
}

func (obj *ShortURLObj) equalID(id uint64) bool {
	return obj.ID == id
}

func (obj *ShortURLObj) getID() uint64 {
	return obj.ID
}

func (obj *ShortURLObj) getVisitShortURLCnt() uint64 {
	return obj.visitShortURLCnt
}

func (obj *ShortURLObj) getVisitLongURLCnt() uint64 {
	return obj.visitLongURLCnt
}

func (obj *ShortURLObj) getShortURL() string {
	return obj.ShortURL
}

func main() {
	r := gin.Default()
	r.GET("/long-url/:url/short-url", getShortURL)
	r.GET("/short-url/:url/long-url", getLongURL)
	r.Run()
}

func getShortURL(c *gin.Context) {

	longURL := c.Param("url")

	obj, err := getShortURLObjByLongURL(longURL)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":              obj.getID(),
		"shortUrl":        obj.getShortURL(),
		"visitLongURLCnt": obj.getVisitLongURLCnt(),
	})
}

func getLongURL(c *gin.Context) {
	longURL := c.Param("url")

	url, visitCnt, err := getLongURLByShortURL(longURL)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"longUrl":          url,
		"visitShortURLCnt": visitCnt,
	})
}

func getShortURLObjByLongURL(longURL string) (URLObj, error) {

	mutex.Lock()
	defer mutex.Unlock()
	obj, ok := m[longURL]

	if !ok {
		if newID >= maxID {
			return nil, ErrIDBiggerThanMaxID
		}

		newID++

		shortURL := encode(newID)

		shortURLObj := &ShortURLObj{
			newID,
			shortURL,
			0,
			1,
		}

		m[longURL] = shortURLObj
		return shortURLObj, nil
	}

	obj.visitLongURL()

	return obj, nil

}

func getLongURLByShortURL(shortURL string) (longURL string, visitCnt uint64, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	id := decode(shortURL)
	for key, value := range m {
		if value.equalID(id) {
			value.visitShortURL()
			return key, value.visitShortURLCnt, nil
		}
	}
	return "", 0, ErrNotFoundShortURL
}

func encode(id uint64) string {

	var b bytes.Buffer

	i := id % 62
	id /= 62

	b.WriteString(string(base62[i]))

	for id > 0 {
		i := id % 62
		id /= 62
		b.WriteString(string(base62[i]))
	}
	return b.String()
}

func decode(url string) uint64 {
	result := uint64(0)
	power := uint64(1)

	for _, char := range url {
		digit := uint64(strings.Index(base62, string(char)))
		result += digit * power
		power *= 62
	}
	return result
}

var ErrIDBiggerThanMaxID = errors.New("can not generate short url because id is bigger than maxID")
var ErrNotFoundShortURL = errors.New("cat not found short url")
