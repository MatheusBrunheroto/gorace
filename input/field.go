package input

type Field struct {
	flag      *Flag
	pairs     []Pair
	delimiter string
}

func NewField(flag *Flag, pairs []Pair, delimiter string) Field {

	var field Field

	field.flag = flag
	field.pairs = pairs
	field.delimiter = delimiter

	return field
}

func initFields(flags map[string]*Flag) map[string]*Field {

	headers := NewField(flags["--headers"], []Pair{}, ":")
	cookies := NewField(flags["--cookies"], []Pair{}, "=")
	data := NewField(flags["--data"], []Pair{}, "=")
	wordlists := NewField(flags["--wordlist"], []Pair{}, "=")

	return map[string]*Field{

		"headers":   &headers,
		"cookies":   &cookies,
		"data":      &data,
		"wordlists": &wordlists,
	}

}
