package admin

import (
	"Blog-Backend/dto/response"
	"Blog-Backend/internal/controller/admin"
)

func GetVisitorMap() ([]response.VisitorMapItem, error) {
	res, err := dao.GetVisitorMap()
}
