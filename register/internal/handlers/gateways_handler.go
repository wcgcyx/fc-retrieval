package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-redis/redis/v8"

	log "github.com/wcgcyx/fc-retrieval/common/pkg/logging"

	"github.com/wcgcyx/fc-retrieval/register/models"
	op "github.com/wcgcyx/fc-retrieval/register/restapi/operations/gateway"
)

// AddGatewayRegister to create a gateway register
func AddGatewayRegister(params op.AddGatewayRegisterParams) middleware.Responder {
	register := params.Register
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	err := rdb.HSet(ctx, "gateway", register.NodeID, register).Err()
	if err != nil {
		log.Error("Unable to set Redis value")
		panic(err)
	}

	log.Info("register created a gateway record with ID: %s", params.Register.NodeID)

	// Response
	return op.NewAddGatewayRegisterOK().WithPayload(register)
}

// GetGatewayRegisters retrieve Gateway register list
func GetGatewayRegisters(_ op.GetGatewayRegistersParams) middleware.Responder {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	gatewayRegisters, err := rdb.HGetAll(ctx, "gateway").Result()

	if err != nil {
		log.Error("Unable to get Redis value")
		panic(err)
	}

	payload := []*models.GatewayRegister{}

	var debugOutputSb strings.Builder
	for _, g := range gatewayRegisters {
		registerData := models.GatewayRegister{}
		if unmarshalErr := json.Unmarshal([]byte(g), &registerData); unmarshalErr != nil {
			log.Error("inside GetGatewayRegisters - can't unmarshall JSON, %s", unmarshalErr.Error())
		}
		payload = append(payload, &registerData)
		debugOutputSb.WriteString(fmt.Sprintf("%s, ", registerData.NodeID))
	}
	//log.Debug("total gateway register records: %d; IDs: %s", len(gatewayRegisters), debugOutputSb.String())

	return op.NewGetGatewayRegistersOK().WithPayload(payload)
}

// GetGatewayRegisterByID retrieve Gateway register by ID
func GetGatewayRegisterByID(params op.GetGatewayRegistersByIDParams) middleware.Responder {
	registerID := params.ID
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	gatewayRegister, err := rdb.HGet(ctx, "gateway", registerID).Result()
	if err != nil {
		msg := "Register not found"
		log.Error(msg)
		return op.NewGetGatewayRegistersByIDDefault(404).WithPayload(&models.Error{Message: &msg})
	}

	registerData := models.GatewayRegister{}
	if unmarshallErr := json.Unmarshal([]byte(gatewayRegister), &registerData); unmarshallErr != nil {
		log.Error("inside GetGatewayRegisterByID - can't unmarshall JSON: %s", unmarshallErr.Error())
	}

	payload := registerData
	return op.NewGetGatewayRegistersByIDOK().WithPayload(&payload)
}

// DeleteGatewayRegisters deletes all Gateways
func DeleteGatewayRegisters(_ op.DeleteGatewayRegisterParams) middleware.Responder {
	const registerTypeGateway = "gateway"

	rdb := redis.NewClient(&redis.Options{
		Addr:     apiconfig.GetString("REDIS_URL") + ":" + apiconfig.GetString("REDIS_PORT"),
		Password: apiconfig.GetString("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	ctx := context.Background()

	registers, err := rdb.HGetAll(ctx, registerTypeGateway).Result()
	if err != nil {
		log.Error("Unable to get Redis value")
		panic(err)
	}

	for index := range registers {
		log.Info("register deleted a gateway record with ID: %s", index)
		err := rdb.HDel(ctx, registerTypeGateway, index).Err()
		if err != nil {
			log.Error("Unable to set Redis value")
			panic(err)
		}
	}

	payload := models.Ack{Status: "success", Message: "All Gateways have been deleted"}
	return op.NewDeleteGatewayRegisterOK().WithPayload(&payload)
}
