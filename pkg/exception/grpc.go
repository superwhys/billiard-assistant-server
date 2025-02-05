// File:		grpc.go
// Created by:	Hoven
// Created on:	2025-01-27
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package exception

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFileTooLarge = NewBilliardException(400, "文件大小超出预期")
)

func ParseGrpcError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	if st.Code() == codes.Unauthenticated {
		return ErrUnauthorized
	}

	return fmt.Errorf("%v", st.Message())
}
