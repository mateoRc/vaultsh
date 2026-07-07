package filesystem

import "time"

type File struct {
	name     string
	content  string
	modTime  time.Time
}

func NewFile(name, content string) *File {
	return NewFileWithModTime(name, content, time.Time{})
}

func NewFileWithModTime(name, content string, modTime time.Time) *File {
	return &File{
		name:     name,
		content:  content,
		modTime:  modTime,
	}
}

func (f *File) Name() string {
	return f.name
}

func (*File) Kind() Kind {
	return KindFile
}

func (f *File) ModTime() time.Time {
	return f.modTime
}

func (f *File) Content() string {
	return f.content
}
