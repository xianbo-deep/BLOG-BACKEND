package utils

import (
	"Blog-Backend/core"
	"Blog-Backend/dto/common"
	"net"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
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
	res.CountryZh = pick(record.Country.Names, "zh-CN", res.CountryCode)
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

	if res.CountryCode == "CN" {
		isV6 := (ip.To4() == nil)
		if region, yep := ip2regionSearch(ipSyr, isV6); yep {
			_, province, city := parseIP2Region(region)

			if res.CountryZh == "" {
				res.CountryZh = "中国"
			}

			if province != "" {
				res.RegionZh = province
			}

			if city != "" {
				res.CityZh = city
				res.CityEn = ""
			}
		}
	}

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

func ip2regionSearch(ipstr string, isV6 bool) (string, bool) {
	var s *xdb.Searcher
	if isV6 {
		s = core.IP2R6
	} else {
		s = core.IP2R4
	}
	if s == nil {
		return "", false
	}
	out, err := s.SearchByStr(ipstr)
	if err != nil {
		return "", false
	}
	return out, true
}

// 获取国家、省、市
func parseIP2Region(region string) (country, province, city string) {
	parts := strings.Split(region, "|")
	if len(parts) >= 1 {
		country = clean(parts[0])
	}
	if len(parts) >= 3 {
		province = clean(parts[2])
	}
	if len(parts) >= 4 {
		city = clean(parts[3])
	}
	return
}

// 过滤无意义数据
func clean(s string) string {
	s = strings.TrimSpace(s)
	if s == "" || s == "0" || s == "null" || s == "unknown" {
		return ""
	}
	return s
}
