package notify

import "time"

type Notify struct {
	ID          string    `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Date        time.Time `json:"date,omitempty"`
	RecipientID string    `json:"recipientId,omitempty"`
}
