package dataquery

import (
	"errors"
	"net/url"
	"strconv"
)

type DataQuery struct {
	Limit   int
	Page    int
	Offset  int
	Search  string
	OrderBy string
	Sort    string
	ShowAll bool
}

func (f *DataQuery) CalculateOffset() int {
	f.Offset = (f.Page - 1) * f.Limit
	return f.Offset
}

func ParseFromURLQuery(uv url.Values) (DataQuery, error) {
	var dq DataQuery
	var err error
	if pg := uv.Get("page"); pg != "" {
		dq.Page, err = strconv.Atoi(pg)
		if err != nil {
			return DataQuery{}, errors.New("query 'page' is expected to be integer")
		}
	}

	if dq.Page == 0 {
		dq.Page = 1
	}

	if lm := uv.Get("limit"); lm != "" {
		dq.Limit, err = strconv.Atoi(lm)
		if err != nil {
			return DataQuery{}, errors.New("query 'limit' is expected to be integer")
		}
	}
	if dq.Limit == 0 {
		dq.Limit = 10
	}

	dq.Search = uv.Get("search")
	dq.OrderBy = uv.Get("order_by")
	dq.Sort = uv.Get("sort")

	return dq, nil
}
