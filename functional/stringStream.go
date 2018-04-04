package functional

//this thing mutates like a mofo. if you need to undo something, use Reset
type stringStream struct {
	sourceData *[]string
	box        []string
}

func CreateStringStream(slice *[]string) *stringStream {
	return &stringStream{
		sourceData: slice,
		box:        *slice,
	}
}

func (stream *stringStream) Reset() *stringStream {
	stream.box = *stream.sourceData
	return stream
}

func (stream *stringStream) Box() []string {
	return stream.box
}

func (stream *stringStream) Filter(predicate func(string) bool) *stringStream {
	container := make([]string, 0)
	for _, val := range stream.box {
		if predicate(val) {
			container = append(container, val)
		}
	}
	stream.box = container
	return stream
}

func IsWhiteSpace(cha byte) bool {
	switch cha {
	case byte(' '):
		return true
	case byte('\t'):
		return true
	case byte('\n'):
		return true
	default:
		return false
	}
}
