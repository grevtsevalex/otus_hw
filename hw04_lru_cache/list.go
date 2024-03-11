package hw04lrucache

func NewList() List {
	return new(DequeueList)
}

type List interface {
	Len() int                          // длина списка
	Front() *ListItem                  // первый элемент списка
	Back() *ListItem                   // последний элемент списка
	PushFront(v interface{}) *ListItem // добавить значение в начало
	PushBack(v interface{}) *ListItem  // добавить значение в конец
	Remove(i *ListItem)                // удалить элемент
	MoveToFront(i *ListItem)           // переместить элемент в начало
	Reset()                            // очистить список
}

type ListItem struct {
	Value interface{} // значение
	Next  *ListItem   // следующий элемент
	Prev  *ListItem   // предыдущий элемент
}

type DequeueList struct {
	front *ListItem // первый элемент списка
	rear  *ListItem // последний элемент списка
	len   int       // длина списка
}

func (l *DequeueList) IsEmpty() bool {
	return l.front == nil
}

func (l *DequeueList) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.IsEmpty() {
		return l.pushToEmptyList(item)
	}
	item.Next = l.front
	l.front.Prev = item
	l.front = item
	l.len++
	return item
}

func (l *DequeueList) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.IsEmpty() {
		return l.pushToEmptyList(item)
	}

	l.rear.Next = item
	item.Prev = l.rear
	l.rear = item
	l.len++
	return item
}

func (l *DequeueList) pushToEmptyList(item *ListItem) *ListItem {
	l.front = item
	l.rear = item
	l.len = 1
	return item
}

func (l *DequeueList) Remove(item *ListItem) {
	if l.IsEmpty() {
		return
	}
	l.len--
	if l.len == 0 {
		l.front = nil
		l.rear = nil
		return
	}

	if l.front == item {
		l.front.Next.Prev = nil
		l.front = l.front.Next
		return
	}

	if l.rear == item {
		l.rear.Prev.Next = nil
		l.rear = l.rear.Prev
		return
	}

	item.Prev.Next = item.Next
	item.Next.Prev = item.Prev
}

func (l *DequeueList) Front() *ListItem {
	return l.front
}

func (l *DequeueList) Back() *ListItem {
	return l.rear
}

func (l *DequeueList) Len() int {
	return l.len
}

func (l *DequeueList) MoveToFront(item *ListItem) {
	if l.front == item {
		return
	}

	l.Remove(item)
	if l.IsEmpty() {
		l.pushToEmptyList(item)
		return
	}
	item.Next = l.front
	l.front.Prev = item
	l.front = item
	l.len++
}

func (l *DequeueList) Reset() {
	l.front = nil
	l.rear = nil
	l.len = 0
}
