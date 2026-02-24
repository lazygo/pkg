package token

import (
	"fmt"
	"strconv"
	"strings"
)

func WrapToken(str string, session int64) string {
	token, _ := UnwrapToken(str)
	return fmt.Sprintf("%s#%d", token, session)
}

func UnwrapToken(str string) (string, int64) {
	data := strings.SplitN(str, "#", 2)
	if len(data) == 2 {
		session, err := strconv.ParseInt(data[1], 10, 64)
		if err != nil {
			session = 0
		}
		return data[0], session
	}
	return data[0], 0
}
