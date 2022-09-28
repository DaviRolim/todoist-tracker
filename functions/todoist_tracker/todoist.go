package main

import (
	"log"
	"time"
)

type TodoistTask struct {
	ID           string      `json:"id"`
	AssignerID   interface{} `json:"assigner_id"`
	AssigneeID   interface{} `json:"assignee_id"`
	ProjectID    string      `json:"project_id"`
	SectionID    interface{} `json:"section_id"`
	ParentID     interface{} `json:"parent_id"`
	Order        int         `json:"order"`
	Content      string      `json:"content"`
	Description  string      `json:"description"`
	IsCompleted  bool        `json:"is_completed"`
	Labels       []string    `json:"labels"`
	Priority     int         `json:"priority"`
	CommentCount int         `json:"comment_count"`
	CreatorID    string      `json:"creator_id"`
	CreatedAt    time.Time   `json:"created_at"`
	Due          TodoistDue  `json:"due"`
	URL          string      `json:"url"`
}

type TodoistDue struct {
	Date        string `json:"date"`
	String      string `json:"string"`
	Lang        string `json:"lang"`
	IsRecurring bool   `json:"is_recurring"`
}

// TODO move this function out of here
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (t *TodoistTask) isTrackeable() bool {
	return contains(t.Labels, "track")
}

func (t *TodoistTask) isRecurring() bool {
	return t.Due.IsRecurring
}

func (t *TodoistTask) wasDoneToday() bool {
	//the value below represesents YYYY-MM-DD according to this https://pkg.go.dev/time#pkg-constants
	const layout = "2006-01-02"
	taksDueTime, err := time.Parse(layout, t.Due.Date)
	if err != nil {
		log.Fatal(err)
	}
	return time.Now().Before(taksDueTime)
}
