package utils

import (
	"Blog-Backend/core"
	"net"
)

func LookupIP(ipSyr string) (country, province, city string) {
	ip := net.ParseIP(ipSyr)
	if ip == nil {
		return
	}
	// 获取记录
	record, err := core.GeoDB.City(ip)
	if err != nil {
		return
	}
	// 获取国家码
	country = record.Country.IsoCode

	if len(record.Subdivisions) > 0 {
		// 获取省份
		province = record.Subdivisions[0].Names["zh-CN"]
	}
	// 获取城市
	city = record.City.Names["zh-CN"]

	return
}
