package helpers

import (
	"errors"
)

var (
	ErrNoContainerFound = errors.New("no container found in the fiber context provided")
)

/*func GetContainer(ctx *fiber.Ctx) (*dic.Container, error) {
	ctn, ok := ctx.Locals(Dependencies).(*dic.Container)
	if !ok {
		return nil, ErrNoContainerFound
	}
	return ctn, nil
}*/

func WrapStrInErrMap(str string) map[string][]string {
	w := make(map[string][]string)
	w["errors"] = []string{str}
	return w
}

func WrapErrInErrMap(err error) map[string][]string {
	w := make(map[string][]string)
	w["errors"] = []string{err.Error()}
	return w
}
