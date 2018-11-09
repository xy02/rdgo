package rdgo

import "time"

type Expr map[Operator]Data

type Operator = string

const (
	F   Operator = "$f"
	EQ  Operator = "$eq"
	GT  Operator = "$gt"
	GTE Operator = "$gte"
	LT  Operator = "$lt"
	LTE Operator = "$lte"
)

func (expr Expr) Match(target Data, ctx *Filter) bool {
	for op, v := range expr {
		if !compute(op, target, v, ctx) {
			return false
		}
	}
	return true
}

func compute(op Operator, target, v Data, ctx *Filter) bool {
	// fn := opFnMap[op]
	// return fn(target, v, ctx)
	switch op {
	case F:
		return OpField(target, v, ctx)
	case EQ:
		return OpEQ(target, v)
	case GT:
		return OpGT(target, v)
	case GTE:
		return OpGTE(target, v)
	case LT:
		return OpLT(target, v)
	case LTE:
		return OpLTE(target, v)
	}
	return false
}

//OpField get field of v
func OpField(target, v Data, ctx *Filter) bool {
	condition, ok := v.(map[string]Data)
	if !ok {
		return false
	}
	for f, a := range condition {
		b, ok := a.(map[string]Data)
		if !ok {
			return false
		}
		expr := Expr(b)
		field := Field(f)
		value := field.PathValue(target, ctx.TagKey)
		if !expr.Match(value, ctx) {
			return false
		}
	}
	return true
}

func OpEQ(target, v Data) bool {
	return target == v
}

func OpGT(target, v Data) bool {
	target = parseNumber(target)
	v = parseNumber(v)
	switch vt := v.(type) {
	case int64:
		switch tt := target.(type) {
		case int64:
			return tt > vt
		case float64:
			return tt > float64(vt)
		}
	case float64:
		switch tt := target.(type) {
		case int64:
			return float64(tt) > vt
		case float64:
			return tt > vt
		}
	}
	return false
}
func OpGTE(target, v Data) bool {
	target = parseNumber(target)
	v = parseNumber(v)
	switch vt := v.(type) {
	case int64:
		switch tt := target.(type) {
		case int64:
			return tt >= vt
		case float64:
			return tt >= float64(vt)
		}
	case float64:
		switch tt := target.(type) {
		case int64:
			return float64(tt) >= vt
		case float64:
			return tt >= vt
		}
	}
	return false
}
func OpLT(target, v Data) bool {
	target = parseNumber(target)
	v = parseNumber(v)
	switch vt := v.(type) {
	case int64:
		switch tt := target.(type) {
		case int64:
			return tt < vt
		case float64:
			return tt < float64(vt)
		}
	case float64:
		switch tt := target.(type) {
		case int64:
			return float64(tt) < vt
		case float64:
			return tt < vt
		}
	}
	return false
}
func OpLTE(target, v Data) bool {
	target = parseNumber(target)
	v = parseNumber(v)
	switch vt := v.(type) {
	case int64:
		switch tt := target.(type) {
		case int64:
			return tt <= vt
		case float64:
			return tt <= float64(vt)
		}
	case float64:
		switch tt := target.(type) {
		case int64:
			return float64(tt) <= vt
		case float64:
			return tt <= vt
		}
	}
	return false
}

func parseNumber(v Data) interface{} {
	switch vt := v.(type) {
	case time.Time:
		return vt.Unix()
	case int:
		return int64(vt)
	case int8:
		return int64(vt)
	case int16:
		return int64(vt)
	case int32:
		return int64(vt)
	case uint:
		return int64(vt)
	case uint8:
		return int64(vt)
	case uint16:
		return int64(vt)
	case uint32:
		return int64(vt)
	case uint64:
		return int64(vt)
	case float32:
		return float64(vt)
	}
	return v
}
