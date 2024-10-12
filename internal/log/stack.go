package log

import "github.com/mdobak/go-xerrors"

func marshalStack(err error, top bool) interface{} {
	out := map[string]interface{}{}
	if !top {
		out["error"] = err.Error()
	}
	if st, ok := err.(xerrors.StackTracer); ok {
		frames := make([]map[string]interface{}, len(st.StackTrace().Frames()))
		for i, frame := range st.StackTrace().Frames() {
			frames[i] = map[string]interface{}{
				"file":     frame.File,
				"function": frame.Function,
				"line":     frame.Line,
			}
		}
		out["frames"] = frames
		return out
	}
	if me, ok := err.(xerrors.MultiError); ok {
		errs := make([]interface{}, len(me.Errors()))
		for i, e := range me.Errors() {
			errs[i] = marshalStack(e, false)
		}
		out["errors"] = errs
	}
	if we, ok := err.(xerrors.Wrapper); ok {
		out["error"] = marshalStack(we.Unwrap(), false)
	}
	if len(out) == 0 {
		return nil
	}
	if len(out) == 1 && !top {
		return out["error"] // unwrap error if there is only a message
	}
	return out
}

func MarshalStack(err error) interface{} {
	return marshalStack(err, true)
}
