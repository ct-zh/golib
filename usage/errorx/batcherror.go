package errorx

import "bytes"

type (
	BatchError struct {
		errs errorArray
	}

	errorArray []error
)

func (e *BatchError) Add(err error) {
	if err != nil {
		e.errs = append(e.errs, err)
	}
}

func (e *BatchError) Err() error {
	switch len(e.errs) {
	case 0:
		return nil
	case 1:
		return e.errs[0]
	default:
		return e.errs
	}
}
func (e *BatchError) NotNil() bool {
	return len(e.errs) > 0
}

func (e errorArray) Error() string {
	var buf bytes.Buffer
	for i := range e {
		if i > 0 {
			buf.WriteByte('\n')
		}
		buf.WriteString(e[i].Error())
	}
	return buf.String()
}
