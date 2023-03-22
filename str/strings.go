package str

import (
	"bytes"
)

type LetterKind = int

const (
	L_Nothing          LetterKind = 1 << iota //什么也不做
	L_Upper                                   //转大写
	L_After_Undlerline                        //遇到"_"后第一个字母处理
	L_Retain                                  //只保留数字和字母
	L_Lower                                   //转小写
)

func IsLetter(char rune) (ok, upper bool) {
	if char >= 'A' && char <= 'Z' {
		ok = true
		upper = true
		return
	}
	if char >= 'a' && char <= 'z' {
		ok = true
		upper = false
		return
	}
	ok = false
	upper = false
	return
}

// name_age->NameAge
func WordsFirstLetter(str string, dosomething LetterKind) string {
	if len(str) == 0 {
		return str
	}
	max := 1
	index := 0
	var char rune
	w := new(bytes.Buffer)
	hit := true
	orgin := []rune(str)
	underline := L_After_Undlerline&dosomething == L_After_Undlerline
	if underline || L_Retain&dosomething == L_Retain {
		max = len(orgin)
		goto read
	}
	goto read
read:
	char = orgin[index]
	//移除字母数字外的字符
	if !hit && underline {
		hit = char == '_'
	}
	ok, upper := IsLetter(char)
	if ok {
		if hit {
			hit = false
			if upper && dosomething&L_Lower == L_Lower {
				goto lower
			}
			if !upper && dosomething&L_Upper == L_Upper {
				goto upper
			}
		}
		w.WriteRune(char)
		goto next
	}
	if dosomething&L_Retain == L_Retain {
		if char >= '0' && char <= '9' {
			w.WriteRune(char)
		}
	} else {
		w.WriteRune(char)
	}
	goto next
next:
	index++
	if index < max {
		goto read
	}
	if len(orgin)-index > 0 {
		w.WriteString(string(orgin[index:]))
	}
	return w.String()
lower:
	w.WriteRune(char + 32)
	goto next
upper:
	w.WriteRune(char - 32)
	goto next
	// if len(str) > 0 && dosomething != Letter_Nothing {
	// 	bts := []byte(str)
	// 	ok, upper := IsLetter(bts[0])
	// 	if !ok {
	// 		return str
	// 	}
	// 	first := bts[0]
	// 	if dosomething == Letter_Lower && upper {
	// 		first += 32
	// 	} else if dosomething == Letter_Upper && !upper {
	// 		first -= 32
	// 	}
	// 	return fmt.Sprintf("%s%s", string(first), str[1:])
	// }
	// return str
}

// 大写字母为标识.通过 "_"分割
func NameToUnderline(str string, toLower bool) string {
	bts := []rune(str)
	writer := bytes.NewBufferString("")
	hit := false
	for i, c := range bts {
		ok, upper := IsLetter(c)
		if !ok {
			writer.WriteRune(c)
			continue
		}
		if i > 0 && !hit && upper {
			hit = true
			writer.WriteByte('_')
		}
		if upper {
			if toLower {
				writer.WriteRune(c + 32)
			} else {
				writer.WriteRune(c)
			}
		} else {
			if hit {
				hit = false
			}
			writer.WriteRune(c)
		}
	}
	return writer.String()
}
