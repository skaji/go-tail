package tail

type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

var NullLogger = &null{}

type null struct{}

func (n *null) Print(v ...interface{})                 {}
func (n *null) Printf(format string, v ...interface{}) {}
func (n *null) Println(v ...interface{})               {}
