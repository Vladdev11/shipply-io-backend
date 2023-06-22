package util

// ConvertWeightToGrams converts a given weight in a given unit to grams.
type WeightUnit string

const (
	WeightUnitKilograms WeightUnit = "KILOGRAMS"
	WeightUnitGrams     WeightUnit = "GRAMS"
	WeightUnitPounds    WeightUnit = "POUNDS"
	WeightUnitOunces    WeightUnit = "OUNCES"
)

func ConvertWeightToGrams(weight float64, unit string) int {
	switch WeightUnit(unit) {
	case WeightUnitKilograms:
		return int(weight * 1000)
	case WeightUnitGrams:
		return int(weight)
	case WeightUnitPounds:
		return int(weight * 453.592)
	case WeightUnitOunces:
		return int(weight * 28.3495)
	default:
		return -1
	}
}
