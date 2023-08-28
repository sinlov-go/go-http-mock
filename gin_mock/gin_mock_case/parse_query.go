package gin_mock_case

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// parseQueryCommonOffsetAndLimit
//
//	parse query offset and limit
//	default offset is 0
//	default limit is 10
//
//nolint:golint,unused
func parseQueryCommonOffsetAndLimit(c *gin.Context) (int, int, error) {
	offsetStr := c.Query("offset")
	offset := 0
	if offsetStr != "" {
		offsetP, err := strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, fmt.Errorf("offset error err: %v", err)
		}
		offset = offsetP
	}

	limitStr := c.Query("limit")
	var limit int
	if limitStr != "" {
		limitP, err := strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, fmt.Errorf("limit error err: %v", err)
		}
		limit = limitP
	} else {
		limit = 10
	}
	return offset, limit, nil
}
