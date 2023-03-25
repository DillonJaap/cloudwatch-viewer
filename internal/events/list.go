package events

const maxDescriptionLength = 50

type Event struct {
	TimeStamp string
	Message   string
}

func (e Event) Title() string       { return e.TimeStamp }
func (e Event) Description() string { return e.getTruncatedDescription() }
func (e Event) FilterValue() string { return e.Message }

func (e Event) getTruncatedDescription() string {
	// TODO make const
	if len(e.Message) > maxDescriptionLength {
		return e.Message[0:maxDescriptionLength-3] + "..."
	}
	return e.Message
}
