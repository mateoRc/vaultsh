package filesystem

type File struct {
	name    string
	content string
}

func NewFile(name, content string) *File {
	return &File{
		name:    name,
		content: content,
	}
}

func (f *File) Name() string {
	return f.name
}

func (*File) Kind() Kind {
	return KindFile
}

func (f *File) Content() string {
	return f.content
}
