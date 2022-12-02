package formula_engine

// stack 栈
type stack struct {
	data    map[int]*Token
	current int
}

// NewStack 新建Stack
func NewStack() *stack {
	return &stack{
		data: map[int]*Token{},
	}
}

// Push 向栈内添加一个字符串
func (st *stack) Push(str *Token) {
	st.current++
	st.data[st.current] = str
}

// Pop 栈向栈外弹出一个Token
func (st *stack) Pop() *Token {
	str := st.data[st.current]
	delete(st.data, st.current)
	st.current--
	return str
}

// 获取栈长度
func (st *stack) len() int {
	return len(st.data)
}

// reverse 反转
func (st *stack) reverse() {
	current := 0
	data := map[int]*Token{}
	for st.len() > 0 {
		current++
		data[current] = st.Pop()
	}
	st.data = data
	st.current = current
}
