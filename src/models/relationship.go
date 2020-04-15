package models

type Relationship struct {
	ID            int64 `json:"id"`
	RequestUserId int64 `json:"requestUserId"`
	TargetUserId  int64 `json:"targetUserId"`
	Status        int64 `json:"status"`
}
