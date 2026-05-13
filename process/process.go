package process

import "keyval/persist"

func Set(k string, v string) error {
	persist.Append(k + v)
	return nil
}
