package find

import (
	"encoding/json"
	"fmt"
)

type CanNotReadDirErr struct {
	Dir string
	Err error
}

func (e CanNotReadDirErr) Error() string {
	return fmt.Sprintf("can not read directory %s: %s", e.Dir, e.Err)
}

func (e CanNotReadDirErr) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Err  string `json:"err"`
		File string `json:"file"`
	}{
		Err:  e.Err.Error(),
		File: e.Dir,
	})
}

type CanNotReadFileErr struct {
	File string
	Err  error
}

func (e CanNotReadFileErr) Error() string {
	return fmt.Sprintf("can not read file %s: %s", e.File, e.Err)
}

func (e CanNotReadFileErr) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Err  string `json:"err"`
		File string `json:"file"`
	}{
		Err:  e.Err.Error(),
		File: e.File,
	})
}
