package convert

import "github.com/shopspring/decimal"

func DecimalPtr(d *decimal.Decimal) decimal.NullDecimal {
	if d == nil {
		return decimal.NullDecimal{Decimal: decimal.Decimal{}, Valid: false}
	}
	return decimal.NullDecimal{Decimal: *d, Valid: true}
}

func NullDecimalPtr(d decimal.NullDecimal) *decimal.Decimal {
	if !d.Valid {
		return nil
	}
	return &d.Decimal
}
