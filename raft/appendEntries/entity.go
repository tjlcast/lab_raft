package appendEntries

import "lab_raft/raft"

type AppendEntriesArgs struct {
	// leader's term
	Term			int

	// so follower can redirect clients
	LeaderId		int

	// log entries to store (if emptry when heartbeat)
	Entries			[]raft.LogEntry

	// term of prevLogIndex entry
	PrevLogInedx	int

	// term of prevLogIndex entry
	PrevLogTerm		int

	// leader's commitIndex
	LeaderCommit	int
}

type AppendEntriesReply struct {

	// the term known by follower.
	Term 			int

	// true when follower containsed entry matching prevlogindex and prevlogterm.
	Success			bool

	// df
	NextIndex		int
}