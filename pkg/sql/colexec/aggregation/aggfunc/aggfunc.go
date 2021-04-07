package aggfunc

import (
	"matrixbase/pkg/container/types"
	"matrixbase/pkg/sql/colexec/aggregation"
	"matrixbase/pkg/sql/colexec/aggregation/avg"
	"matrixbase/pkg/sql/colexec/aggregation/count"
	"matrixbase/pkg/sql/colexec/aggregation/max"
	"matrixbase/pkg/sql/colexec/aggregation/min"
	"matrixbase/pkg/sql/colexec/aggregation/starcount"
	"matrixbase/pkg/sql/colexec/aggregation/sum"
)

func NewCount(typ types.Type) aggregation.Aggregation {
	return count.New(aggregation.ReturnType(aggregation.Count, typ))
}

func NewStarCount(typ types.Type) aggregation.Aggregation {
	return starcount.New(aggregation.ReturnType(aggregation.StarCount, typ))
}

func NewMax(typ types.Type) aggregation.Aggregation {
	switch typ.Oid {
	case types.T_int8:
		return max.NewInt8(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_int16:
		return max.NewInt16(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_int32:
		return max.NewInt32(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_int64:
		return max.NewInt64(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_uint8:
		return max.NewUint8(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_uint16:
		return max.NewUint16(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_uint32:
		return max.NewUint32(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_uint64:
		return max.NewUint64(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_float32:
		return max.NewFloat32(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_float64:
		return max.NewFloat64(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_varchar:
		return max.NewStr(aggregation.ReturnType(aggregation.Max, typ))
	}
	return nil
}

func NewMin(typ types.Type) aggregation.Aggregation {
	switch typ.Oid {
	case types.T_int8:
		return min.NewInt8(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_int16:
		return min.NewInt16(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_int32:
		return min.NewInt32(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_int64:
		return min.NewInt64(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_uint8:
		return min.NewUint8(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_uint16:
		return min.NewUint16(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_uint32:
		return min.NewUint32(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_uint64:
		return min.NewUint64(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_float32:
		return min.NewFloat32(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_float64:
		return min.NewFloat64(aggregation.ReturnType(aggregation.Max, typ))
	case types.T_varchar:
		return min.NewStr(aggregation.ReturnType(aggregation.Max, typ))
	}
	return nil
}

func NewSum(typ types.Type) aggregation.Aggregation {
	switch typ.Oid {
	case types.T_float32, types.T_float64:
		return sum.NewFloat(aggregation.ReturnType(aggregation.Sum, typ))
	case types.T_int8, types.T_int16, types.T_int32, types.T_int64:
		return sum.NewInt(aggregation.ReturnType(aggregation.Sum, typ))
	case types.T_uint8, types.T_uint16, types.T_uint32, types.T_uint64:
		return sum.NewUint(aggregation.ReturnType(aggregation.Sum, typ))
	}
	return nil
}

func NewAvg(typ types.Type) aggregation.Aggregation {
	switch typ.Oid {
	case types.T_tuple:
		return avg.NewSumCount(aggregation.ReturnType(aggregation.Avg, typ))
	case types.T_float32, types.T_float64:
		return avg.NewFloat(aggregation.ReturnType(aggregation.Avg, typ))
	case types.T_int8, types.T_int16, types.T_int32, types.T_int64:
		return avg.NewInt(aggregation.ReturnType(aggregation.Avg, typ))
	case types.T_uint8, types.T_uint16, types.T_uint32, types.T_uint64:
		return avg.NewUint(aggregation.ReturnType(aggregation.Avg, typ))
	}
	return nil
}

func NewSumCount(typ types.Type) aggregation.Aggregation {
	switch typ.Oid {
	case types.T_float32, types.T_float64:
		return sum.NewFloatSumCount(aggregation.ReturnType(aggregation.SumCount, typ))
	case types.T_int8, types.T_int16, types.T_int32, types.T_int64:
		return sum.NewIntSumCount(aggregation.ReturnType(aggregation.SumCount, typ))
	case types.T_uint8, types.T_uint16, types.T_uint32, types.T_uint64:
		return sum.NewUintSumCount(aggregation.ReturnType(aggregation.SumCount, typ))
	}
	return nil
}