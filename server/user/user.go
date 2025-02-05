// File:		game.go
// Created by:	Hoven
// Created on:	2024-10-29
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package userSrv

import (
	"context"
	"io"
	"mime/multipart"

	"gitea.hoven.com/billiard/billiard-assistant-server/domain/user"
	"gitea.hoven.com/billiard/billiard-assistant-server/pkg/exception"
	"gitea.hoven.com/core/auth-core/api/grpc/interceptor"
	"gitea.hoven.com/core/auth-core/api/grpc/proto/userpb"
	"gitea.hoven.com/core/auth-core/server/dto"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var _ user.IUserService = (*UserService)(nil)

type UserService struct {
	repo       user.IUserRepo
	userClient userpb.AuthCoreUserHandlerClient
}

func NewUserService(userRepo user.IUserRepo, userClient userpb.AuthCoreUserHandlerClient) *UserService {
	return &UserService{
		repo:       userRepo,
		userClient: userClient,
	}
}

func (us *UserService) parseUser(u *dto.User) *user.User {
	return &user.User{
		UserId: int(u.GetUserId()),
		Name:   u.GetName(),
		Gender: user.Gender(u.Gender),
		UserInfo: &user.BaseInfo{
			Email:  u.GetEmail(),
			Phone:  u.GetPhone(),
			Avatar: u.GetAvatar(),
		},
		Status: user.Status(u.GetStatus()),
		Role:   user.Role(u.GetRole()),
	}
}

func (us *UserService) injectToken(ctx context.Context, token string) context.Context {
	md := metadata.Pairs(interceptor.TokenMetadataKey, token)
	return metadata.NewOutgoingContext(ctx, md)
}

func (us *UserService) UpsertUser(ctx context.Context, userId int) error {
	u, err := us.repo.GetUser(ctx, userId)
	if err != nil && !errors.Is(err, exception.ErrUserNotFound) {
		return err
	}

	if u != nil {
		return nil
	}

	_, err = us.repo.CreateUser(ctx, &user.User{UserId: userId})
	return err
}

func (us *UserService) GetUserProfile(ctx context.Context, token string) (*user.User, error) {
	ctx = us.injectToken(ctx, token)

	user, err := us.userClient.GetUserProfile(ctx, &dto.GetUserProfileRequest{})
	if err != nil {
		return nil, exception.ParseGrpcError(err)
	}

	return us.parseUser(user), nil
}

func (us *UserService) UpdateUserName(ctx context.Context, token string, name string) error {
	ctx = us.injectToken(ctx, token)

	_, err := us.userClient.UpdateUserName(ctx, &dto.UpdateUserNameRequest{Username: name})
	if err != nil {
		return exception.ParseGrpcError(err)
	}

	return nil
}

func (us *UserService) UpdateUserGender(ctx context.Context, token string, gender int) error {
	ctx = us.injectToken(ctx, token)

	_, err := us.userClient.UpdateUserGender(ctx, &dto.UpdateUserGenderRequest{Gender: int32(gender)})
	if err != nil {
		return exception.ParseGrpcError(err)
	}

	return nil
}

func (us *UserService) UploadAvatar(ctx context.Context, token string, fh *multipart.FileHeader) (string, error) {
	ctx = us.injectToken(ctx, token)

	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	stream, err := us.userClient.UploadAvatar(ctx)
	if err != nil {
		return "", exception.ParseGrpcError(err)
	}

	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			plog.Errorc(ctx, "failed to read avatar, error: %v", err)
			return "", err
		}

		err = stream.Send(&dto.AvatarByte{
			ByteData: buf[:n],
		})
		if err != nil {
			plog.Errorc(ctx, "failed to send avatar, error: %v", err)
			return "", err
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return "", errors.Wrap(err, "closeAndRecv")
	}

	return resp.GetAvatarUrl(), nil
}

func (us *UserService) GetAvatar(ctx context.Context, avatarId string, writer io.Writer) error {
	stream, err := us.userClient.GetAvatar(ctx, &dto.GetAvatarRequest{AvatarId: avatarId})
	if err != nil {
		return exception.ParseGrpcError(err)
	}

	reader := &streamReader{Stream: stream}
	_, err = io.Copy(writer, reader)
	if err != nil {
		return errors.Wrap(err, "ReadFromStream")
	}

	return nil
}

type streamReader struct {
	Stream grpc.ServerStreamingClient[dto.AvatarByte]
	buffer []byte
}

func (sr *streamReader) Read(p []byte) (n int, err error) {
	if len(sr.buffer) > 0 {
		n = copy(p, sr.buffer)
		sr.buffer = sr.buffer[n:]

		return n, nil
	}

	req, err := sr.Stream.Recv()
	if err == io.EOF {
		return 0, io.EOF
	} else if status.Code(err) == codes.ResourceExhausted {
		return 0, exception.ErrFileTooLarge
	} else if err != nil {
		return 0, err
	}

	sr.buffer = req.ByteData

	n = copy(p, sr.buffer)
	sr.buffer = sr.buffer[n:]
	return n, nil
}
