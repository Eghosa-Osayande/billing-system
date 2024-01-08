package util

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type OrderBy string

type PagedResult[k any] struct {
	Items  []k     `json:"items"`
	Cursor *string `json:"cursor"`
}

func NewPagedResult[k any](items []k, cursor *string) *PagedResult[k] {

	return &PagedResult[k]{
		Items:  items,
		Cursor: cursor,
	}
}

func ListToPagedResult[T any](
	items []T,
	props func(item T) (time.Time, int64),
) PagedResult[T] {

	var nextCursor *string = nil
	if len(items) > 0 {
		createdAt, id := props(items[len(items)-1])
		cursor := EncodeCursor(createdAt, id)
		nextCursor = &cursor
	}
	return PagedResult[T]{
		Items:  items,
		Cursor: nextCursor,
	}
}

func DecodeCursor(encodedCursor string) (res time.Time, uuid int64, err error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return
	}

	arrStr := strings.Split(string(byt), ",")
	if len(arrStr) != 2 {
		err = errors.New("cursor is invalid")
		return
	}

	res, err = time.Parse(time.RFC3339Nano, arrStr[0])
	if err != nil {
		return
	}
	
	uuid,err = strconv.ParseInt(arrStr[1],10,64)
	if err != nil {
		return
	}
	return
}

func EncodeCursor(t time.Time, uuid int64) string {
	key := fmt.Sprintf("%s,%v", t.Format(time.RFC3339Nano), uuid)
	return base64.StdEncoding.EncodeToString([]byte(key))
}
