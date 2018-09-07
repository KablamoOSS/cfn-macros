package transform

func Transform(fragment map[string]interface{}) []byte {
	code, ok := fragment["Brainfuck"]
	if !ok {
		return []byte{}
	}

	return runBrainfuck(code.(string))
}

func runBrainfuck(code string) []byte {
	output := []byte{}
	data := []byte{0}
	dataPtr := 0

	for i := 0; i < len(code); i++ {
		switch code[i] {
		case '>':
			dataPtr++
			if dataPtr >= len(data) {
				data = append(data, 0)
			}
		case '<':
			if dataPtr > 0 {
				dataPtr--
			}
		case '+':
			data[dataPtr]++
		case '-':
			data[dataPtr]--
		case '.':
			output = append(output, data[dataPtr])
		case '[':
			if data[dataPtr] == 0 {
				skip := 0
				i++
				for code[i] != ']' || skip > 0 {
					if code[i] == '[' {
						skip++
					} else if code[i] == ']' {
						skip--
					}
					i++
				}
			}
		case ']':
			if data[dataPtr] != 0 {
				skip := 0
				i--
				for code[i] != '[' || skip > 0 {
					if code[i] == ']' {
						skip++
					} else if code[i] == '[' {
						skip--
					}
					i--
				}
			}
		}
	}

	return output
}
