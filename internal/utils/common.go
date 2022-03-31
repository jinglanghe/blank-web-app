package utils

import "k8s.io/apimachinery/pkg/api/resource"

func RoundDown(sum, num int64) int64 {
	return sum * num / 100
}

func RoundUp(sum, num int64) int64 {
	if (sum*num)%100 > 0 {
		return (sum * num / 100) + 1
	}
	return sum * num / 100
}

func DividedUp(dividend, divisor int64) int64 {
	if dividend*100%divisor > 0 {
		return dividend*100/divisor + 1
	}
	return dividend * 100 / divisor
}

func Max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func QuantityRoundUp(q *resource.Quantity) int64 {
	req, success := q.AsInt64()
	if success != true {
		return q.Value()
	}
	return req
}
