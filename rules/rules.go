package rules

type Rule struct {
	Name string
	Description string
	Pattern     []byte
	Predicate   string
}

func GetRules() []Rule {
	return []Rule{
		{
			Name: "Caption Trailing Period",
			Description: "A caption should not have a trailing period, because it would end up in the ToX as well.",
			Pattern: []byte(`
				(caption
				  long: (curly_group
					(text 
					  (word) @last_word (#match? @last_word "\\.$")
					  .
					)
				  )
				) @caption
			`),
			Predicate: "caption",
		},
	}
}
