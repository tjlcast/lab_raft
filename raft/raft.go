package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"lab_raft/labrpc"
	"time"
	"math/rand"
	"sync"
)

// import "bytes"
// import "encoding/gob"


/**
	tips:
	https://pdos.csail.mit.edu/6.824/labs/lab-raft.html						==> lab introduction
	https://github.com/tjlcast/raft-zh_cn/blob/master/raft-zh_cn.md			==> raft paper
	https://github.com/tjlcast/MIT-6.824/blob/master/src/raft/raft.go		==> init code repo
	https://github.com/comiser/MIT-6.824-2016/blob/master/src/raft/raft.go	==> nice code
**/

// the state of raft node.
const (
	ROLE_FOLLOWER = iota
	ROLE_LEADER
	ROLE_CANDIDATE
)


//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make().
//
type ApplyMsg struct {
	Index       int
	Command     interface{}
	UseSnapshot bool   // ignore for lab2; only used in lab3
	Snapshot    []byte // ignore for lab2; only used in lab3
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        	sync.Mutex
	peers     	[]*labrpc.ClientEnd
	persister 	*Persister
	me        	int // index into peers[]

	role		int // raft node 的状态：follower\candidate\leader

	currentTerm	int				// 服务器最后一次知道的任期号（初始化为 0，持续递增）
	votedFor	int				// 在当前获得选票的候选人的 Id
	logs		[]LogEntry		// 日志条目集；每一个条目包含一个用户状态机执行的指令，和收到时的任期号

	voteCount	int 			// 成为 candidate 后，获得的选票

	chanLeader	chan bool		// 成为 leader 进行通知
	// Your data here.
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.
	// 查看论文的图2部分，可知

	/*
	 * 全部服务器上面的可持久化状态:
	 *  currentTerm 	服务器看到的最近Term(第一次启动的时候为0,后面单调递增)
	 *  votedFor     	当前Term收到的投票候选 (如果没有就为null)
	 *  log[]        	日志项; 每个日志项包含机器状态和被leader接收的Term(first index is 1)
         */
	//  删除代码部分
	/*
	 * 全部服务器上面的不稳定状态:
	 *	commitIndex 	已经被提交的最新的日志索引(第一次为0,后面单调递增)
	 *	lastApplied      已经应用到服务器状态的最新的日志索引(第一次为0,后面单调递增)
	*/
	//  删除代码部分

	/*
	 * leader上面使用的不稳定状态（完成选举之后需要重新初始化）
	 *	nextIndex[]
	 *
	 *
	*/

}


// setter and getter
func GetRole(rt *Raft) int {
	return rt.role
}

func IsRole(rt *Raft) bool {
	return rt.role == ROLE_LEADER
}


// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	var term int
	var isleader bool
	// Your code here.
	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here.
	// Example:
	// w := new(bytes.Buffer)
	// e := gob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	// Your code here.
	// Example:
	// r := bytes.NewBuffer(data)
	// d := gob.NewDecoder(r)
	// d.Decode(&rf.xxx)
	// d.Decode(&rf.yyy)
}




//
// RequestVote RPC arguments structure.
//
type RequestVoteArgs struct {
	Term			int		// 候选人的任期号
	CandidateId		int		// 请求选票的候选人的 Id
	LastLogIndex	int		// 候选人的最后日志条目的索引值
	LastLogTerm		int		// 候选人最后日志条目的任期号
}

//
// RequestVote RPC reply structure.
//
type RequestVoteReply struct {
	Term			int 	// 当前任期号，以便于候选人去更新自己的任期号
	VoteGranted		bool	// 候选人赢得了此张选票时为真
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here.

	// step.1 加锁
	rf.mu.Lock()
	defer rf.mu.Unlock()
	// todo

	reply.VoteGranted = false

	// step.2 当前自己知道有较新的 term
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		return
	}

	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.role = ROLE_FOLLOWER
		rf.votedFor = -1
	}
	reply.Term = rf.currentTerm

	term := rf.getLastLogTerm()
	index := rf.getLastLogIndex()
	uptoDate := false

	if args.LastLogTerm > term {
		uptoDate = true
	}
	if args.LastLogTerm > term && args.LastLogIndex > index {
		uptoDate = true
	}

	if (rf.votedFor == -1 || rf.votedFor == args.CandidateId) && uptoDate {
		// rf.chanGrantVote <- true
		rf.role = ROLE_FOLLOWER
		reply.VoteGranted = true
		rf.votedFor = args.CandidateId
	}

}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// returns true if labrpc says the RPC was delivered.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)

	//
	// 发送请求，阻塞等待网路请求结束.
	// 之后把结果向 上级 报告（扇入|扇出）
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if ok {
		term := rf.currentTerm
		if ROLE_CANDIDATE != rf.role {
			return ok
		}
		if args.Term != term {
			return ok
		}
		if reply.Term > term {
			rf.currentTerm = reply.Term
			rf.role = ROLE_FOLLOWER
			rf.votedFor = -1
			rf.persist()
		}
		if reply.VoteGranted {
			rf.voteCount++
			if rf.role == ROLE_CANDIDATE && rf.voteCount > len(rf.peers)/2 {
				rf.role = ROLE_LEADER
				rf.chanLeader <- true
			}
		}
	}

	return ok
}

//
// when node to be a candidate, it will start a leader election
// the node try to request all nodes that in a same cluster.
func (rf *Raft) broadcastRequestVote() {
	var voteArgs = RequestVoteArgs{}
	rf.mu.Lock()
	voteArgs.Term = rf.currentTerm
	voteArgs.CandidateId = rf.me
	voteArgs.LastLogIndex = rf.getLastLogIndex()
	voteArgs.LastLogTerm = rf.getLastLogTerm()
	rf.mu.Unlock()

	for peer := range rf.peers {
		if ROLE_CANDIDATE == rf.role {
			go func(id int) {
				var voteReply RequestVoteReply
				rf.sendRequestVote(peer, voteArgs, &voteReply)
			}(peer)
		}
	}
}


//
//
func (rf *Raft) getLastLogIndex() int {
	// todo
	return -1
}

func (rf *Raft) getLastLogTerm() int {
	// todo
	return -1
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true


	return index, term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
	// todo how to kill a raft node.
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
// 建一个Raft端点。
// peers参数是通往其他Raft端点处于连接状态下的RPC连接。
// me参数是自己在端点数组中的索引。
//
// applyCh 参数是实验收集器。当最近的日志项被提交时，发送一条ApplyMsg到applyCh。
func Make(peers []*labrpc.ClientEnd,
	me int,
	persister *Persister,
	applyCh chan ApplyMsg) *Raft {

	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.role = ROLE_FOLLOWER
	// todo raft node init

	// Your initialization code here.
	// start groutines.
	go func() {
		// + kick off leader election periodically by sending out RequestVote RPCs
		// 	when it hasn't heard from another peer for a while.
		// + This way a peer will learn who is the leader, if there is already a leader,
		// 	or become the leader itself.
		for {
			switch rf.role {
			case ROLE_FOLLOWER:
				select {
				// block：follower 超时并转化为 candidate
				case <- time.After(time.Duration(rand.Int63() % 333 + 550) * time.Millisecond):
					rf.role = ROLE_CANDIDATE
				// todo maybe there are some other msg case...
				}
			case ROLE_CANDIDATE:
				//	一旦成为 candidate 开启新一轮投票
				// 	+ 自己的 term 加一
				// 	+ 给自己投票
				// 	+ 向其他节点询问投票
				//	+ 等待投票结果(select)
				rf.mu.Lock()
				rf.currentTerm++
				rf.votedFor = rf.me
				go rf.broadcastRequestVote()
				select {
				case <- time.After(time.Duration(rand.Int63() % 333 + 550) * time.Millisecond):
					// block: candidate 状态下的超时
					// 重新进入 candidate 状态
				case <- rf.chanLeader:
					// block: candidate 成功被选举为 leader
					rf.mu.Lock()
					rf.role = ROLE_LEADER
					// todo prepare to be a leader
					rf.mu.Unlock()
				}
				rf.mu.Unlock()

			case ROLE_LEADER:
				// + 向其他的 follower 发送 AppendEntry.
				// todo
			}
		}
	}()


	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())


	return rf
}
