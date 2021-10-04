package types

type ComboArgs []ComboArg

type ComboArg struct {
	Name    string
	Options []string
}

type Combos []map[string]string
