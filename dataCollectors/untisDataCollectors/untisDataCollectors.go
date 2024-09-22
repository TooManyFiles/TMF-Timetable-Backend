package untisDataCollectors

import (
	"errors"

	"github.com/Mr-Comand/goUntisAPI/structs"
	"github.com/Mr-Comand/goUntisAPI/untisApi"
)

type UntisClient struct {
	client *untisApi.Client
}

func Init(apiConfig structs.ApiConfig) (UntisClient, error) {
	untisClient := UntisClient{client: untisApi.NewClient(apiConfig)}
	err := untisClient.client.Authenticate()
	if err != nil {
		return UntisClient{}, err
	}
	untisClient.client.Test()
	// untisClient.client.Logout()
	return untisClient, nil
}

func (untisClient UntisClient) reAuthenticate() error {
	err := untisClient.client.Test()
	if err != nil {
		var rpcerr structs.RPCError
		if errors.As(err, rpcerr) && rpcerr.Code == -8520 {
			return untisClient.client.Authenticate()
		}
		return err
	}
	return nil
}
func (untisClient UntisClient) GetTeachers() ([]structs.Teacher, error) {
	err := untisClient.reAuthenticate()
	if err != nil {
		return nil, err
	}
	teachers, err := untisClient.client.GetTeachers()
	if err != nil {
		return nil, err
	}
	return teachers, nil
}
func (untisClient UntisClient) GetSubjects() ([]structs.Subject, error) {
	err := untisClient.reAuthenticate()
	if err != nil {
		return nil, err
	}
	subjects, err := untisClient.client.GetSubjects()
	if err != nil {
		return nil, err
	}
	return subjects, nil
}
func (untisClient UntisClient) GetRooms() ([]structs.Room, error) {
	err := untisClient.reAuthenticate()
	if err != nil {
		return nil, err
	}
	subjects, err := untisClient.client.GetRooms()
	if err != nil {
		return nil, err
	}
	return subjects, nil
}
func (untisClient UntisClient) GetClasses() ([]structs.Class, error) {
	err := untisClient.reAuthenticate()
	if err != nil {
		return nil, err
	}
	classes, err := untisClient.client.GetClasses()
	if err != nil {
		return nil, err
	}
	return classes, nil
}
