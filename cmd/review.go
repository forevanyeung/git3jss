package main

func longest(input []string) string {
	a := ""
	length := 0

	for i := range input {
		l := len(input[i])
		if l > length {
			a = input[i]
			length = l
		}
	}

	return a
}

func extract(x []any, field string) []string {

	for _, s := range scripts {
		names = append(names, s.Name)
	}
}
