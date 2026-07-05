package events

type Event string

const (
	TestEvent Event = "TestEvent"
)

func (e Event) ToString() string {
	return string(e)
}
