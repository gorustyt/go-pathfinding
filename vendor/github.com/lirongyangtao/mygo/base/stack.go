package base

type Stack interface {
	Empty() bool
	Size() int
	Pop() (v interface{})
	Push(v interface{})
}

type SimpleStack struct {
	values []interface{}
}

func NewSimpleStack() Stack {
	return &SimpleStack{}
}

func (s *SimpleStack) Empty() bool {
	return s.Size() == 0
}
func (s *SimpleStack) Size() int {
	return len(s.values)
}
func (s *SimpleStack) Pop() (v interface{}) {
	if s.Empty() {
		return nil
	}
	ele := s.values[s.Size()-1]
	s.values = s.values[:s.Size()-1]
	return ele
}
func (s *SimpleStack) Push(v interface{}) {
	s.values = append(s.values, v)
}
