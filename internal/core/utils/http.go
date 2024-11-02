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

type QueryMap map[string]string

func NewQueryMapper(queries map[string]string) *QueryMap {
	var queryBuilder QueryMap
	queryBuilder = queries
	return &queryBuilder
}

func (qb QueryMap) PaginationQuery() (map[string]int, error) {
	paginationQueries := map[string]int{}
	limit, ok := qb["limit"]
	if !ok {
		paginationQueries["limit"] = LIMIT
	} else {
		limitUint, err := strconv.ParseUint(limit, 10, 32)
		if err != nil {
			return paginationQueries, err
		}
		paginationQueries["limit"] = int(limitUint)
	}

	page, ok := qb["page"]
	if !ok {
		paginationQueries["page"] = PAGE
	} else {
		pageUint, err := strconv.ParseUint(page, 10, 32)
		if err != nil {
			return paginationQueries, err
		}
		if pageUint > 0 {
			pageUint--
		}
		paginationQueries["page"] = int(pageUint)
	}
	return paginationQueries, nil
}

func ExtractID(context *gin.Context) int {
	bearer := context.Request.Header.Get("Authorization")
	token := strings.Split(bearer, " ")[1]
	payload, _ := GetIdentity(token)
	return payload.UserID
}
