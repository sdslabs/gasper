package raft

import (
	context2 "context"
	"fmt"
	"github.com/chrislusf/raft/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"sync"
	"time"
)

var (
	// cache grpc connections
	grpcClients     = make(map[string]*grpc.ClientConn)
	grpcClientsLock sync.Mutex
)

// An GrpcTransporter is a default transport layer used to communicate between
// multiple servers.
type GrpcTransporter struct {
	grpcDialOption grpc.DialOption
}

// Creates a new HTTP transporter with the given path prefix.
func NewGrpcTransporter(grpcDialOption grpc.DialOption) *GrpcTransporter {
	t := &GrpcTransporter{
		grpcDialOption: grpcDialOption,
	}
	return t
}

type GrpcServer struct {
	server Server
}

// Creates a new HTTP transporter with the given path prefix.
func NewGrpcServer(server Server) *GrpcServer {
	t := &GrpcServer{
		server: server,
	}
	return t
}

//--------------------------------------
// Outgoing
//--------------------------------------

// Sends an AppendEntries RPC to a peer.
func (t *GrpcTransporter) SendAppendEntriesRequest(server Server, peer *Peer, req *AppendEntriesRequest) (ret *AppendEntriesResponse) {

	err := withRaftServerClient(peer.ConnectionString, t.grpcDialOption, func(client protobuf.RaftClient) error {
		ctx, cancel := context2.WithTimeout(context2.Background(), time.Duration(5*time.Second))
		defer cancel()

		pbReq := &protobuf.AppendEntriesRequest{
			Term:         req.Term,
			PrevLogIndex: req.PrevLogIndex,
			PrevLogTerm:  req.PrevLogTerm,
			CommitIndex:  req.CommitIndex,
			LeaderName:   req.LeaderName,
			Entries:      req.Entries,
		}

		resp, err := client.OnSendAppendEntriesRequest(ctx, pbReq)
		if err != nil {
			return err
		}

		ret = &AppendEntriesResponse{
			pb: resp,
		}

		return nil

	})
	if err != nil {
		return nil
	}
	return
}

// Sends a RequestVote RPC to a peer.
func (t *GrpcTransporter) SendVoteRequest(server Server, peer *Peer, req *RequestVoteRequest) (ret *RequestVoteResponse) {

	err := withRaftServerClient(peer.ConnectionString, t.grpcDialOption, func(client protobuf.RaftClient) error {
		ctx, cancel := context2.WithTimeout(context2.Background(), time.Duration(5*time.Second))
		defer cancel()

		pbReq := &protobuf.RequestVoteRequest{
			Term:          req.Term,
			LastLogIndex:  req.LastLogIndex,
			LastLogTerm:   req.LastLogTerm,
			CandidateName: req.CandidateName,
		}

		resp, err := client.OnSendVoteRequest(ctx, pbReq)
		if err != nil {
			return err
		}

		ret = &RequestVoteResponse{
			Term:        resp.Term,
			VoteGranted: resp.VoteGranted,
		}

		return nil

	})
	if err != nil {
		return nil
	}
	return

}

// Sends a SnapshotRequest RPC to a peer.
func (t *GrpcTransporter) SendSnapshotRequest(server Server, peer *Peer, req *SnapshotRequest) (ret *SnapshotResponse) {

	err := withRaftServerClient(peer.ConnectionString, t.grpcDialOption, func(client protobuf.RaftClient) error {
		ctx, cancel := context2.WithTimeout(context2.Background(), time.Duration(5*time.Second))
		defer cancel()

		pbReq := &protobuf.SnapshotRequest{
			LeaderName: req.LeaderName,
			LastIndex:  req.LastIndex,
			LastTerm:   req.LastTerm,
		}

		resp, err := client.OnSendSnapshotRequest(ctx, pbReq)
		if err != nil {
			return err
		}

		ret = &SnapshotResponse{
			Success: resp.Success,
		}

		return nil

	})
	if err != nil {
		return nil
	}
	return

}

// Sends a SnapshotRequest RPC to a peer.
func (t *GrpcTransporter) SendSnapshotRecoveryRequest(server Server, peer *Peer, req *SnapshotRecoveryRequest) (ret *SnapshotRecoveryResponse) {

	err := withRaftServerClient(peer.ConnectionString, t.grpcDialOption, func(client protobuf.RaftClient) error {
		ctx, cancel := context2.WithTimeout(context2.Background(), time.Duration(5*time.Second))
		defer cancel()

		var peers []*protobuf.SnapshotRecoveryRequest_Peer
		for _, peer := range req.Peers {
			peers = append(peers, &protobuf.SnapshotRecoveryRequest_Peer{
				Name:             peer.Name,
				ConnectionString: peer.ConnectionString,
			})
		}

		pbReq := &protobuf.SnapshotRecoveryRequest{
			LeaderName: req.LeaderName,
			LastIndex:  req.LastIndex,
			LastTerm:   req.LastTerm,
			Peers:      peers,
			State:      req.State,
		}

		resp, err := client.OnSendSnapshotRecoveryRequest(ctx, pbReq)
		if err != nil {
			return err
		}

		ret = &SnapshotRecoveryResponse{
			Term:        resp.Term,
			Success:     resp.Success,
			CommitIndex: resp.CommitIndex,
		}

		return nil

	})
	if err != nil {
		return nil
	}
	return

}

//--------------------------------------
// Incoming
//--------------------------------------

// Handles incoming AppendEntries requests.
func (t *GrpcServer) OnSendAppendEntriesRequest(ctx context2.Context, pbReq *protobuf.AppendEntriesRequest) (*protobuf.AppendEntriesResponse, error) {
	req := &AppendEntriesRequest{
		Term:         pbReq.Term,
		PrevLogIndex: pbReq.PrevLogIndex,
		PrevLogTerm:  pbReq.PrevLogTerm,
		CommitIndex:  pbReq.CommitIndex,
		LeaderName:   pbReq.LeaderName,
		Entries:      pbReq.Entries,
	}

	resp := t.server.AppendEntries(req)
	if resp == nil {
		return nil, fmt.Errorf("failed creating response")
	}
	return &protobuf.AppendEntriesResponse{
		Term:        resp.Term(),
		Index:       resp.Index(),
		CommitIndex: resp.CommitIndex(),
		Success:     resp.Success(),
	}, nil
}

// Handles incoming RequestVote requests.
func (t *GrpcServer) OnSendVoteRequest(ctx context2.Context, pbReq *protobuf.RequestVoteRequest) (*protobuf.RequestVoteResponse, error) {

	req := &RequestVoteRequest{
		Term:          pbReq.Term,
		LastLogIndex:  pbReq.LastLogIndex,
		LastLogTerm:   pbReq.LastLogTerm,
		CandidateName: pbReq.CandidateName,
	}

	resp := t.server.RequestVote(req)
	if resp == nil {
		return nil, fmt.Errorf("failed creating response")
	}
	return &protobuf.RequestVoteResponse{
		Term:        resp.Term,
		VoteGranted: resp.VoteGranted,
	}, nil
}

// Handles incoming Snapshot requests.
func (t *GrpcServer) OnSendSnapshotRequest(ctx context2.Context, pbReq *protobuf.SnapshotRequest) (*protobuf.SnapshotResponse, error) {

	req := &SnapshotRequest{
		LeaderName: pbReq.LeaderName,
		LastIndex:  pbReq.LastIndex,
		LastTerm:   pbReq.LastTerm,
	}

	resp := t.server.RequestSnapshot(req)
	if resp == nil {
		return nil, fmt.Errorf("failed creating response")
	}
	return &protobuf.SnapshotResponse{
		Success: resp.Success,
	}, nil
}

// Handles incoming SnapshotRecovery requests.
func (t *GrpcServer) OnSendSnapshotRecoveryRequest(ctx context2.Context, pbReq *protobuf.SnapshotRecoveryRequest) (*protobuf.SnapshotRecoveryResponse, error) {

	var peers []*Peer
	for _, peer := range pbReq.Peers {
		peers = append(peers, &Peer{
			Name:             peer.Name,
			ConnectionString: peer.ConnectionString,
		})
	}

	req := &SnapshotRecoveryRequest{
		LeaderName: pbReq.LeaderName,
		LastIndex:  pbReq.LastIndex,
		LastTerm:   pbReq.LastTerm,
		Peers:      peers,
		State:      pbReq.State,
	}

	resp := t.server.SnapshotRecoveryRequest(req)
	if resp == nil {
		return nil, fmt.Errorf("failed creating response")
	}
	return &protobuf.SnapshotRecoveryResponse{
		Term:        resp.Term,
		Success:     resp.Success,
		CommitIndex: resp.CommitIndex,
	}, nil
}

func withRaftServerClient(raftServer string, grpcDialOption grpc.DialOption, fn func(protobuf.RaftClient) error) error {

	return withCachedGrpcClient(func(grpcConnection *grpc.ClientConn) error {
		client := protobuf.NewRaftClient(grpcConnection)
		return fn(client)
	}, raftServer, grpcDialOption)

}

func grpcDial(address string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// opts = append(opts, grpc.WithBlock())
	// opts = append(opts, grpc.WithTimeout(time.Duration(5*time.Second)))
	var options []grpc.DialOption
	options = append(options,
		// grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30 * time.Second, // client ping server if no activity for this long
			Timeout: 20 * time.Second,
		}))
	for _, opt := range opts {
		if opt != nil {
			options = append(options, opt)
		}
	}
	return grpc.Dial(address, options...)
}

func withCachedGrpcClient(fn func(*grpc.ClientConn) error, address string, opts ...grpc.DialOption) error {

	grpcClientsLock.Lock()

	existingConnection, found := grpcClients[address]
	if found {
		grpcClientsLock.Unlock()
		return fn(existingConnection)
	}

	grpcConnection, err := grpcDial(address, opts...)
	if err != nil {
		grpcClientsLock.Unlock()
		return fmt.Errorf("fail to dial %s: %v", address, err)
	}

	grpcClients[address] = grpcConnection
	grpcClientsLock.Unlock()

	err = fn(grpcConnection)
	if err != nil {
		grpcClientsLock.Lock()
		delete(grpcClients, address)
		grpcClientsLock.Unlock()
	}

	return err
}
