package service

import (
	"encoding/json"
	"gin-chat-svc/app/model"
	"gin-chat-svc/pkg/common/constant"
	"gin-chat-svc/pkg/common/response"
	"gin-chat-svc/pkg/logger"
	"gin-chat-svc/pkg/misprint"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (hu *HttpUserService) GetUserData(userData []response.ResponseUserData) ([]response.ResponseUserData, error) {
	logger.Logger.Debug("service", logger.Any("fetching api", userData))

	client := &http.Client {
		Timeout: time.Duration(60 * time.Second),
	}

	req, err := http.NewRequest(constant.METHOD_GET, constant.BASE_URL + constant.GET_ALL, nil)
	if err != nil {
		return nil, misprint.New("error: " + err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, misprint.New("error: " + err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, misprint.New("error: " + err.Error())
	}
	json.Unmarshal(bodyBytes, &userData)

	return userData, nil
}

// store the data into database server after fetching the JSON
func (hu *HttpUserService) StoreUserData(userData []response.ResponseUserData) ([]response.ResponseUserData, error) {
	logger.Logger.Debug("service", logger.Any("fetching api", userData))
	
	// var db = service.Db
	var userCount int64

	getData, err := hu.GetUserData(userData)
	if err != nil {
		return nil, misprint.New("error: " + err.Error())
	}

	for _, data := range getData {
		logger.Logger.Debug("service", logger.Any("storing data", data))

		modelUser := model.User {
			Uuid: 		(constant.USER_TAG_UUID + uuid.New().String()),
			Username: 	data.Username,
			Password:	data.Password,
			Nickname:	data.Nickname,
			Avatar: 	data.Avatar,
			Email:		data.Email,
		}

		Db.Model(modelUser).Where("username", modelUser.Username).Count(&userCount)
		if userCount > 0 {
			continue
		}

		Db.Save(&modelUser)
	}

	return getData, nil
}