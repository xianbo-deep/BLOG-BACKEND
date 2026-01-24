package utils

import (
	"Blog-Backend/core"
	"Blog-Backend/dto/common"
	"net"
)

func LookupIP(ipSyr string) (res common.GeoInfo, ok bool) {
	ip := net.ParseIP(ipSyr)
	if ip == nil {
		return common.GeoInfo{}, false
	}
	// 获取记录
	record, err := core.GeoDB.City(ip)
	if err != nil {
		return common.GeoInfo{}, false
	}
	// 获取国家码
	res.CountryCode = record.Country.IsoCode
	res.CountryEn = pick(record.Country.Names, "zh-CN", res.CountryCode)
	res.CountryEn = pick(record.Country.Names, "en", res.CountryCode)

	if len(record.Subdivisions) > 0 {
		sub := record.Subdivisions[0]
		res.RegionCode = sub.IsoCode
		res.RegionZh = pick(sub.Names, "zh-CN", res.RegionCode)
		res.RegionEn = pick(sub.Names, "en", res.RegionCode)
	}

	// 城市
	res.CityZh = pick(record.City.Names, "zh-CN", "")
	res.CityEn = pick(record.City.Names, "en", "")

	// 经纬度
	res.Lat = record.Location.Latitude
	res.Lon = record.Location.Longitude

	return res, true
}

func pick(m map[string]string, lang, fallback string) string {
	if m == nil || len(m) == 0 {
		return fallback
	}
	if v := m[lang]; v != "" {
		return v
	}
	return fallback
}
