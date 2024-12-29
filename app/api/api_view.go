package api

import (
	"strconv"
	"strings"
)

type ApiViewReq struct {
	Short          bool   `query:"short"`
	ReplaceImgPath bool   `query:"replaceImgPath"`
	Range          string `header:"range"`
}

func (req *ApiViewReq) GetRange() (start int, end int) {
	rangeInfo := strings.Split(strings.TrimPrefix(req.Range, "bytes="), "-")
	if len(rangeInfo) == 2 {
		start, _ = strconv.Atoi(rangeInfo[0])
		end, _ = strconv.Atoi(rangeInfo[1])
	}
	return
}

type ApiViewRes struct{}
