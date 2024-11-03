package utils

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	LIMIT = 10
	PAGE  = 0
)

type JSON map[string]any

type QueryMap struct {
	Queries map[string]string
	Pq      paginationParams
}

func (qb *QueryMap) SetPaginate() error {
	params := paginationParams{}
	paginationQueries := map[string]int{}

	limit, ok := qb.Queries["limit"]
	if !ok {
		paginationQueries["limit"] = LIMIT
	} else {
		limitUint, err := strconv.ParseUint(limit, 10, 32)
		if err != nil {
			return err
		}
		params.Limit = uint8(limitUint)
	}

	page, ok := qb.Queries["page"]
	if !ok {
		paginationQueries["page"] = PAGE
	} else {
		pageUint, err := strconv.ParseUint(page, 10, 32)
		if err != nil {
			return err
		}
		params.page = uint8(pageUint)
	}

	params.calculateOffset()
	qb.Pq = params
	return nil
}

func (qb *QueryMap) Set(k, v string) {
	qb.Queries[k] = v
}

type paginationParams struct {
	page   uint8
	Limit  uint8
	Offset uint8
}

func (p *paginationParams) calculateOffset() {
	p.Offset = p.Limit * (p.page - 1)
}

func NewQueryMap(queries map[string]string) *QueryMap {
	var queryBuilder QueryMap
	queryBuilder.Queries = queries
	return &queryBuilder
}
func ExtractID(context *gin.Context) int {
	bearer := context.Request.Header.Get("Authorization")
	token := strings.Split(bearer, " ")[1]
	payload, _ := GetIdentity(token)
	return payload.UserID
}
