package pusher

import (
	"context"
	"github.com/clearcodecn/carim/proto"
)

func (i *ImServer) AddFriend(ctx context.Context, request *proto.AddFriendRequest) (*proto.AddFriendReply, error) {
	err := i.imService.PushId(request.ToNo, request)
	if err != nil {
		return nil, err
	}
	return &proto.AddFriendReply{}, nil
}
