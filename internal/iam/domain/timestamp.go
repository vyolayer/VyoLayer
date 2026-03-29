package domain

import "time"

type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewTimestamps() *Timestamps {
	return &Timestamps{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}
}

func (t *Timestamps) Update() {
	t.UpdatedAt = time.Now()
}

func (t *Timestamps) Delete() {
	t.DeletedAt = &t.UpdatedAt
}

func (t *Timestamps) Restore() {
	t.DeletedAt = nil
}

func (t *Timestamps) IsDeleted() bool {
	return t.DeletedAt != nil
}
