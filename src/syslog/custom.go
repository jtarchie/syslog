package syslog

import (
	"errors"
	"time"
)

// func atoi(a []byte) int {
// 	var x int
// 	for _, c := range a {
// 		x = x*10 + int(c-'0')
// 	}
// 	return x
// }

// func atoi2(a []byte) int {
// 	return int(a[1]-'0') +
// 		int(a[0]-'0')*10
// }

// func atoi4(a []byte) int {
// 	return int(a[3]-'0') +
// 		int(a[2]-'0')*10 +
// 		int(a[1]-'0')*100 +
// 		int(a[0]-'0')*1000
// }

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isSpace(c byte) bool {
	return c == ' '
}

func isSpecial(c byte) bool {
	return c == '=' || c == ' ' || c == ']' || c == 34 || c == '"'
}

var PError = errors.New("Could not parse")

func CustomParser(data []byte) (*message, error) {
	if data[0] != '<' {
		return nil, PError
	}

	msg, p, mark := new(message), 1, 1
	for isDigit(data[p]) {
		p++
	}

	if data[p] != '>' {
		return nil, PError
	}

	msg.priority = atoi(data[1:p])
	p++
	mark = p

	if data[p] == 0 {
		return nil, PError
	}

	for isDigit(data[p]) {
		p++
	}

	if data[p] != ' ' {
		return nil, PError
	}

	msg.version = atoi(data[mark:p])
	p++
	mark = p

	for !isSpace(data[p]) {
		p++
	}

	if p-mark > 1 {
		if data[mark+4] != '-' ||
			data[mark+7] != '-' ||
			data[mark+10] != 'T' ||
			data[mark+13] != ':' ||
			data[mark+16] != ':' ||
			data[p-1] != 'Z' {
			return nil, PError
		}

		nanosecond := 0
		if data[mark+19] == '.' {
			nbytes := (p - 2) - (mark + 19)
			for i := mark + 20; i < p-1; i++ {
				nanosecond = nanosecond*10 + int(data[i]-'0')
			}
			for i := 0; i < 9-nbytes; i++ {
				nanosecond *= 10
			}
		}

		t := time.Date(
			atoi4(data[mark:mark+4]),
			time.Month(atoi2(data[mark+5:mark+7])),
			atoi2(data[mark+8:mark+10]),
			atoi2(data[mark+11:mark+13]),
			atoi2(data[mark+14:mark+16]),
			atoi2(data[mark+17:mark+19]),
			nanosecond,
			time.UTC,
		)
		msg.timestamp = &t
	}

	p++
	mark = p

	for !isSpace(data[p]) {
		p++
	}

	if data[mark] != '-' && p-mark > 1 {
		msg.hostname = data[mark:p]
	}

	p++
	mark = p

	for !isSpace(data[p]) {
		p++
	}

	if data[mark] != '-' && p-mark > 1 {
		msg.appname = data[mark:p]
	}

	p++
	mark = p

	for !isSpace(data[p]) {
		p++
	}

	if data[mark] != '-' && p-mark > 1 {
		msg.procID = data[mark:p]
	}

	p++
	mark = p

	for !isSpace(data[p]) {
		p++
	}

	if data[mark] != '-' && p-mark > 1 {
		msg.msgID = data[mark:p]
	}

	p++
	if data[p] == '[' {
		p++
		mark = p

		for !isSpecial(data[p]) {
			p++
		}

		msg.data = &structureData{
			id:         data[mark:p],
			properties: make([]Property, 0, 5),
		}

		for data[p] != ']' {
			p++
			mark = p
			for !isSpecial(data[p]) {
				p++
			}
			paramName := data[mark:p]

			p += 2
			mark = p
			for !isSpecial(data[p]) {
				p++
			}
			msg.data.properties = append(msg.data.properties, Property{paramName, data[mark:p]})
			p++
		}
	}

	return msg, nil
}
