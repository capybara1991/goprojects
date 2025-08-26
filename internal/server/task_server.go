package server

import (
	"context"

	pb "grpctasks/proto"
	"grpctasks/internal/store"
)

type TaskServer struct{ st *store.Memory }

func NewTaskServer(st *store.Memory) *TaskServer { return &TaskServer{st: st} }

func (s *TaskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.Task, error) {
	t, err := s.st.CreateTask(req.GetTitle(), req.GetDescription())
	if err != nil { return nil, err }
	return &pb.Task{Id: t.ID, Title: t.Title, Description: t.Description, Completed: t.Completed, CreatedAtUnix: t.CreatedAt.Unix(), CompletedAtUnix: store.ToUnix(t.CompletedAt)}, nil
}

func (s *TaskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.Task, error) {
	t, err := s.st.GetTask(req.GetId())
	if err != nil { return nil, err }
	return &pb.Task{Id: t.ID, Title: t.Title, Description: t.Description, Completed: t.Completed, CreatedAtUnix: t.CreatedAt.Unix(), CompletedAtUnix: store.ToUnix(t.CompletedAt)}, nil
}

func (s *TaskServer) SetCompleted(ctx context.Context, req *pb.SetCompletedRequest) (*pb.Task, error) {
	t, err := s.st.SetCompleted(req.GetId(), req.GetComplete())
	if err != nil { return nil, err }
	return &pb.Task{Id: t.ID, Title: t.Title, Description: t.Description, Completed: t.Completed, CreatedAtUnix: t.CreatedAt.Unix(), CompletedAtUnix: store.ToUnix(t.CompletedAt)}, nil
}

func (s *TaskServer) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	ts := s.st.ListTasks(req.GetOnlyUncompleted(), int(req.GetLimit()), int(req.GetOffset()))
	items := make([]*pb.Task, 0, len(ts))
	for _, t := range ts {
		items = append(items, &pb.Task{Id: t.ID, Title: t.Title, Description: t.Description, Completed: t.Completed, CreatedAtUnix: t.CreatedAt.Unix(), CompletedAtUnix: store.ToUnix(t.CompletedAt)})
	}
	return &pb.ListTasksResponse{Items: items}, nil
}
