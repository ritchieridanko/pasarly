package processors

import (
	"context"
	"fmt"
	"strings"

	"github.com/ritchieridanko/pasarly/backend/services/user/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/models"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/repositories"
	"github.com/ritchieridanko/pasarly/backend/services/user/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"github.com/ritchieridanko/pasarly/backend/shared/events/v1"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/proto"
)

const userErrTracer string = "processor.user"

type UserProcessor interface {
	OnAuthCreated(ctx context.Context, m kafka.Message) (err error)
}

type userProcessor struct {
	ur         repositories.UserRepository
	transactor *database.Transactor
}

func NewUserProcessor(ur repositories.UserRepository, tx *database.Transactor) UserProcessor {
	return &userProcessor{ur: ur, transactor: tx}
}

func (p *userProcessor) OnAuthCreated(ctx context.Context, m kafka.Message) error {
	ctx, span := otel.Tracer(userErrTracer).Start(ctx, "OnAuthCreated")
	defer span.End()

	var evt events.AuthCreated
	if err := proto.Unmarshal(m.Value, &evt); err != nil {
		e := fmt.Errorf("failed to process message: %w", err)
		utils.TraceErr(span, e, ce.MsgInternalServer)
		return e
	}

	err := p.transactor.WithTx(ctx, func(ctx context.Context) *ce.Error {
		exists, err := p.ur.Exists(ctx, evt.GetAuthId())
		if err != nil {
			return err
		}
		if exists {
			return nil
		}

		userID := utils.NewUUID().String()
		data := models.CreateUser{
			AuthID: evt.GetAuthId(),
			UserID: userID,
			Name:   fmt.Sprintf("user_%s", strings.Split(userID, "-")[0]),
		}

		_, err = p.ur.CreateUser(ctx, &data)
		return err
	})
	if err != nil {
		return err.Err
	}

	return nil
}
