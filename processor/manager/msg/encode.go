package msg

type Entry struct {
	Id       string `json:"id"`
	Location string `json:"location"`
}
type HeartbeatContent struct {
	Entry Entry `json:"entry"`
}
