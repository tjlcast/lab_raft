package raft

type LogEntry struct {
	LogIndex int
	LogTerm int
	LogComd interface{}
}
