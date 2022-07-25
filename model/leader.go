package model

import "time"

type Leader struct {
	UUID          string
	LastHeartBeat time.Time
	Startup       time.Time
	IsLeader      bool
	InactiveFlag  bool
}
