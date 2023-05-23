package client

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth"
	protomodels "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/models"
	proto "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/delivery/grpc/proto"
	hasherPkg "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/auth/hasher"
	"github.com/go-park-mail-ru/2023_1_PracticalDev/internal/models"
	pkgErrors "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/errors"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const defaultAvatar = "https://pickpin.hb.bizmrg.com/default-user-icon-8-4024862977"
const defaultAccountType = "personal"

type client struct {
	authClient proto.AuthenficatorClient
}

func NewAuthenficatorClient(con *grpc.ClientConn) auth.Service {
	return &client{authClient: proto.NewAuthenficatorClient(con)}
}

func (client *client) Authenticate(login, hashedPassword string) (models.User, auth.SessionParams, error) {
	authParams := proto.LoginParams{
		Email:    login,
		Password: hashedPassword,
	}
	resp, err := client.authClient.Authenticate(context.TODO(), &authParams)
	fmt.Println("here", resp, err)
	if err != nil {
		return models.User{}, auth.SessionParams{}, pkgErrors.RestoreHTTPError(pkgErrors.GRPCUnwrapper(err))
	}

	sessionParams := client.CreateSession(int(resp.GetID()))

	session := models.Session{
		UserId:    int(resp.GetID()),
		UserEmail: resp.GetEmail(),
	}
	err = client.SetSession(sessionParams.Token, &session, sessionParams.LivingTime)
	if err != nil {
		return models.User{}, auth.SessionParams{}, errors.Wrap(err, "Authenticate")
	}

	fmt.Println("here")
	return *protomodels.NewUser(resp), sessionParams, err

}

func (client *client) CreateSession(userId int) auth.SessionParams {
	token := strconv.Itoa(userId) + "$" + uuid.New().String()
	livingTime := 5 * time.Hour
	return auth.SessionParams{Token: token, LivingTime: livingTime}
}

func (client *client) SetSession(token string, session *models.Session, expiration time.Duration) error {
	setParams := proto.SessionSetParams{
		Token: token,
		Session: &proto.Session{
			UserId:    int64(session.UserId),
			UserEmail: session.UserEmail,
		},
		Experation: expiration.Nanoseconds(),
	}
	_, err := client.authClient.SetSession(context.TODO(), &setParams)

	return pkgErrors.RestoreHTTPError(pkgErrors.GRPCUnwrapper(err))
}

func (client *client) CheckAuth(userId, sessionId string) (models.User, error) {
	checkParams := proto.SessionCheckParams{
		UserId:    userId,
		SessionId: sessionId,
	}

	user, err := client.authClient.CheckAuth(context.TODO(), &checkParams)
	if err != nil {
		return models.User{}, pkgErrors.RestoreHTTPError(pkgErrors.GRPCUnwrapper(err))
	}

	return *protomodels.NewUser(user), err
}

func (client *client) DeleteSession(userId, sessionId string) error {
	checkParams := proto.SessionCheckParams{
		UserId:    userId,
		SessionId: sessionId,
	}

	_, err := client.authClient.CheckAuth(context.TODO(), &checkParams)
	if err != nil {
		return pkgErrors.RestoreHTTPError(pkgErrors.GRPCUnwrapper(err))
	}

	return nil
}

func (client *client) Register(user *auth.RegisterParams) (models.User, auth.SessionParams, error) {
	hasher := hasherPkg.NewHasher()
	hash, _ := hasher.GetHashedPassword(user.Password)

	tmp := models.User{
		Name:           user.Name,
		Username:       user.Username,
		Email:          user.Email,
		ProfileImage:   defaultAvatar,
		WebsiteUrl:     "",
		AccountType:    defaultAccountType,
		HashedPassword: hash,
	}

	usr, err := client.authClient.Register(context.TODO(), protomodels.NewProtoUser(&tmp))
	if err != nil {
		return models.User{}, auth.SessionParams{}, pkgErrors.RestoreHTTPError(pkgErrors.GRPCUnwrapper(err))
	}

	return client.Authenticate(usr.GetEmail(), user.Password)
}
