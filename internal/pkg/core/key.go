package core

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func Range(start, end int) []int {
	result := make([]int, 0)
	if start <= end {
		for start <= end {
			result = append(result, start)
			start++
		}
	} else {
		for start >= end {
			result = append(result, start)
			start--
		}
	}
	return result
}

// IndexList 返回符合条件的索引列表
// key: 关键字
//	1.选择前6个：'1,2,3,4,5,6' 或 '1-3,4-6' 或 '1-6' 或 '-6'
//	2.选择第6个及后面的所有：'6-'
//	3.选择第6个：'6'
//	4.选择所有：'all' 或 '-'
//注意：超出部分会被忽略，'all' 只能单独使用
// max: 最大索引
func IndexList(key string, max int) []int {
	if max == 0 {
		return []int{}
	}
	if key == "all" {
		return Range(1, max)
	}
	result := make([]int, 0)
	for _, item := range strings.Split(key, ",") {
		item = strings.Trim(item, " ")
		re1 := "^[1-9][0-9]*$"
		re2 := "(^[0-9]*)-([0-9]*$)"
		if re, _ := regexp.Compile(re1); re.MatchString(item) {
			i, _ := strconv.Atoi(item)
			if i > 0 && i <= max {
				result = append(result, i)
			}
			continue
		}
		if re, _ := regexp.Compile(re2); re.MatchString(item) {
			start := 1
			end := max
			s := re.FindStringSubmatch(item)
			if s[1] != "" {
				start, _ = strconv.Atoi(s[1])
			}
			if s[2] != "" {
				end, _ = strconv.Atoi(s[2])
			}
			if start > end {
				start, end = end, start
			}
			if start > max || end < 1 {
				continue
			}
			if start < 1 {
				start = 1
			}
			if end > max {
				end = max
			}
			result = append(result, Range(start, end)...)
			continue
		}
	}
	result = RemoveRepByMap(result)
	sort.Ints(result)
	return result
}

// RemoveRepByMap 通过map主键唯一的特性过滤重复元素
func RemoveRepByMap(slc []int) []int {
	var result []int
	tempMap := map[int]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

func Reverse(slc []int) []int {
	for i := 0; i < len(slc)/2; i++ {
		j := len(slc) - i - 1
		slc[i], slc[j] = slc[j], slc[i]
	}
	return slc
}
