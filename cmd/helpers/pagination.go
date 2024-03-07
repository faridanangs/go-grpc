package helpers

import (
	"go_grpc_yt/pb/pagination"
	"math"

	"gorm.io/gorm"
)

func Pagination(sql *gorm.DB, pagination *pagination.Pagination, page int64) (int64, int64) {
	var (
		total  int64
		limit  int64 = 3
		offset int64
	)

	sql.Count(&total)
	if page == 1 {
		offset = 0
	} else {
		offset = (page - 1) * limit
	}
	pagination.Total = uint64(total)
	pagination.CurentPage = uint32(page)
	pagination.PerPage = uint32(limit)
	pagination.LastPage = uint32(math.Ceil(float64(total) / float64(limit)))

	return offset, limit
}
