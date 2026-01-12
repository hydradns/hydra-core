package dataplane

import (
	"context"

	"github.com/lopster568/phantomDNS/internal/dnsengine"
	pb "github.com/lopster568/phantomDNS/internal/gen/proto/phantomdns/v1"
)

type StatusService struct {
	pb.UnimplementedDataPlaneStatusServiceServer
	engine *dnsengine.Engine
}

func NewStatusService(engine *dnsengine.Engine) *StatusService {
	return &StatusService{
		engine: engine,
	}
}

func (s *StatusService) GetStatus(ctx context.Context, _ *pb.StatusRequest) (*pb.StatusResponse, error) {
	st := s.engine.Status()

	return &pb.StatusResponse{
		Running:          st.Running,
		AcceptingQueries: st.AcceptingQueries,
		PolicyEnabled:    st.PolicyEnabled,
		LastError:        st.LastError,
	}, nil
}

func (s *StatusService) SetAcceptQueries(
	ctx context.Context,
	req *pb.SetAcceptQueriesRequest,
) (*pb.SetAcceptQueriesResponse, error) {

	// Apply desired state to the engine.
	s.engine.SetAcceptQueries(req.Enabled)

	return &pb.SetAcceptQueriesResponse{
		Ok: true,
	}, nil
}
