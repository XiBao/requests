package requests

func (this *Response) Text() string {
	return string(this.Content)
}