package handler

import (
	"context"
	"fmt"
	iauth "realworld-example-app/internal/auth"
	protocol "realworld-example-app/proto/realworld"
	"time"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/util/pubsub"
)

// We must hear subscriber's heartbeat within this duration
const heartbeatDuration = 1 * time.Minute

// Clubhouse .
type Clubhouse struct {
	g *pubsub.Group
}

// NewClubhouse create a new instance
func NewClubhouse() *Clubhouse {
	return &Clubhouse{
		g: pubsub.New(),
	}
}

type streamCtx struct {
	account   *auth.Account
	stream    protocol.Clubhouse_SubscribeStream
	cancel    context.CancelFunc
	heartbeat *time.Time
}

type streamCtxKey struct{}

type serverPush struct {
	*protocol.PublishReq
}

func (m *serverPush) Topic() string {
	return m.GetPublishNote().GetTopic()
}

func (m *serverPush) Body() interface{} {
	return m.GetPublishNote().GetText()
}

func (c *Clubhouse) Publish(ctx context.Context, req *protocol.PublishReq, resp *protocol.PublishResp) error {
	c.g.Publish(ctx, &serverPush{req})
	return nil
}

// Subscribe serves a new subscriber
func (c *Clubhouse) Subscribe(ctx context.Context, stream protocol.Clubhouse_SubscribeStream) error {
	req, err := stream.Recv()
	if err != nil {
		return errors.BadRequest("incorrect-protocol", "stream recv: %v", err)
	}
	account, ok := iauth.AccountFromToken(req.Token)
	if !ok {
		return errors.BadRequest("unidentified-subscriber", "")
	}
	logger.Infof("subscriber %q enter", account.ID)

	heartbeat := time.Now()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	ctx = context.WithValue(ctx, streamCtxKey{}, streamCtx{account, stream, cancel, &heartbeat})
	c.g.Go(ctx, c.processHeartbeat)
	c.g.Subscribe(ctx, account.ID, c.onServerPush, pubsub.WithTicker(heartbeatDuration, c.onServerTick))
	return nil
}

func (c *Clubhouse) processHeartbeat(ctx context.Context) error {
	sc := ctx.Value(streamCtxKey{}).(streamCtx)
	defer func() { sc.cancel() }()
	for {
		if _, err := sc.stream.Recv(); err != nil {
			return fmt.Errorf("process %q heartbeat: stream recv: %v", sc.account.ID, err)
		}
		*sc.heartbeat = time.Now()
	}
}

func (c *Clubhouse) onServerTick(ctx context.Context) (err error) {
	sc := ctx.Value(streamCtxKey{}).(streamCtx)
	defer func() {
		if err != nil {
			sc.cancel()
		}
	}()
	if sc.heartbeat == nil {
		return fmt.Errorf("server tick %q: no heartbeat", sc.account.ID)
	}
	if time.Since(*sc.heartbeat) > 2*heartbeatDuration {
		return fmt.Errorf("server tick %q: heartbeat delays", sc.account.ID)
	}
	if err := sc.stream.Send(&protocol.ServerPush{T: protocol.ServerPush_HEARTBEAT}); err != nil {
		return fmt.Errorf("server tick %q: send heartbeat: %v", sc.account.ID, err)
	}
	return nil
}

func (c *Clubhouse) onServerPush(ctx context.Context, msg pubsub.Message) (bool, error) {
	sc := ctx.Value(streamCtxKey{}).(streamCtx)
	err := sc.stream.Send(&protocol.ServerPush{
		T: protocol.ServerPush_PUBLISH_NOTE,
		PublishNote: &protocol.PublishNote{
			Topic: msg.Topic(),
			Text:  msg.Body().(string),
		},
	})
	if err != nil {
		sc.cancel()
		return false, fmt.Errorf("server push %q: stream send: %v", sc.account.ID, err)
	}
	return true, nil
}
