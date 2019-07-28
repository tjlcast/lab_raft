package raft

type LogEntry struct {

	LogIndex int

	LogTerm int

	LogCmd interface{}
}
