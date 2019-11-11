package kaen

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/factory"
	pb "github.com/sdslabs/gasper/lib/factory/protos/database"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
	"google.golang.org/grpc"
)

// ServiceName is the name of the current microservice
const ServiceName = types.Kaen

type server struct{}

// Create creates a database of the specified type
func (s *server) Create(ctx context.Context, body *pb.RequestBody) (*pb.ResponseBody, error) {
	language := body.GetLanguage()
	db := &types.DatabaseConfig{}

	err := json.Unmarshal(body.GetData(), db)
	if err != nil {
		return nil, err
	}

	db.SetInstanceType(mongo.DBInstance)
	db.SetHostIP(utils.HostIP)
	db.SetUser(db.GetName())
	db.SetOwner(body.GetOwner())

	if pipeline[language] == nil {
		return nil, fmt.Errorf("Database type `%s` is not supported", language)
	}

	pipeline[language].init(db)

	err = pipeline[language].create(db)
	if err != nil {
		go pipeline[language].cleanup(db.GetName())
		return nil, err
	}

	err = mongo.UpsertInstance(
		types.M{
			mongo.NameKey:         db.GetName(),
			mongo.InstanceTypeKey: mongo.DBInstance,
		}, db)
	if err != nil && err != mongo.ErrNoDocuments {
		go pipeline[language].cleanup(db.GetName())
		return nil, err
	}

	err = redis.RegisterDB(
		db.GetName(),
		fmt.Sprintf("%s:%d", utils.HostIP, configs.ServiceConfig.Kaen.Port),
	)
	if err != nil {
		go pipeline[language].cleanup(db.GetName())
		return nil, err
	}

	err = redis.IncrementServiceLoad(
		language,
		fmt.Sprintf("%s:%d", utils.HostIP, configs.ServiceConfig.Kaen.Port),
	)
	if err != nil {
		go pipeline[language].cleanup(db.GetName())
		return nil, err
	}

	db.SetSuccess(true)

	response, err := json.Marshal(db)
	return &pb.ResponseBody{Data: response}, err
}

// Delete deletes a database of the specified type
func (s *server) Delete(ctx context.Context, body *pb.NameHolder) (*pb.GenericResponse, error) {
	db, err := mongo.FetchSingleDatabase(body.GetName())
	if err != nil {
		return nil, err
	}
	if pipeline[db.Language] == nil {
		return nil, fmt.Errorf("Database type `%s` is not supported", db.Language)
	}
	err = pipeline[db.Language].delete(db.GetName())
	if err != nil {
		return nil, err
	}
	err = redis.RemoveDB(db.GetName())
	if err != nil {
		return nil, err
	}
	filter := types.M{
		mongo.NameKey:         db.GetName(),
		mongo.InstanceTypeKey: mongo.DBInstance,
	}
	_, err = mongo.DeleteInstance(filter)
	return &pb.GenericResponse{Success: true}, err
}

// FetchLogs returns the docker logs from the specified database server's container
func (s *server) FetchLogs(ctx context.Context, body *pb.LogRequest) (*pb.LogResponse, error) {
	language := body.GetLanguage()
	if pipeline[language] == nil {
		return nil, fmt.Errorf("Database type `%s` is not supported", language)
	}
	data, err := pipeline[language].logs(body.GetTail())
	return &pb.LogResponse{
		Success: true,
		Data:    data,
	}, err
}

// Reload restarts the specified database server
func (s *server) Reload(ctx context.Context, body *pb.LanguageHolder) (*pb.GenericResponse, error) {
	language := body.GetLanguage()
	if pipeline[language] == nil {
		return nil, fmt.Errorf("Database type `%s` is not supported", language)
	}
	err := pipeline[language].reload()
	return &pb.GenericResponse{Success: true}, err
}

// NewService returns a new instance of the current microservice
func NewService() *grpc.Server {
	return factory.NewDatabaseFactory(&server{})
}
