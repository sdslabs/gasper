package appmaker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/api"
	"github.com/sdslabs/gasper/lib/cloudflare"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/lib/factory"
	pb "github.com/sdslabs/gasper/lib/factory/protos/application"
	"github.com/sdslabs/gasper/lib/mongo"
	"github.com/sdslabs/gasper/lib/redis"
	"github.com/sdslabs/gasper/lib/utils"
	"github.com/sdslabs/gasper/types"
	"google.golang.org/grpc"
)

// ServiceName is the name of the current microservice
const ServiceName = types.AppMaker

type server struct{}

// Create creates an application
func (s *server) Create(ctx context.Context, body *pb.RequestBody) (*pb.ResponseBody, error) {
	language := body.GetLanguage()
	app := &types.ApplicationConfig{}

	err := json.Unmarshal(body.GetData(), app)
	if err != nil {
		return nil, err
	}

	user, err := mongo.FetchSingleUser(body.GetOwner())
	if err != nil {
		return nil, err
	}

	maxCount := configs.ServiceConfig.AppMaker.AppLimit
	rateCount := configs.ServiceConfig.RateLimit
	timeInterval := configs.ServiceConfig.RateInterval
	if !user.IsAdmin() && maxCount >= 0 {
		rateLimitCount := mongo.CountInstanceInTimeFrame(body.GetOwner(), mongo.AppInstance, timeInterval)
		totalCount := mongo.CountInstancesByUser(body.GetOwner(), mongo.AppInstance)
		if totalCount < maxCount {
			if rateLimitCount >= rateCount && rateCount >= 0 {
				return nil, fmt.Errorf("cannot deploy more than %d app instances in %d hours", rateCount, timeInterval)
			}
		} else {
			return nil, fmt.Errorf("cannot deploy more than %d app instances", maxCount)
		}
	}

	app.SetLanguage(language)
	app.SetOwner(body.GetOwner())
	app.SetInstanceType(mongo.AppInstance)
	app.SetHostIP(utils.HostIP)
	app.SetNameServers(configs.GasperConfig.DNSServers)
	app.SetDateTime()

	gendnsNameServers, _ := redis.FetchServiceInstances(types.GenDNS)
	for _, nameServer := range gendnsNameServers {
		if strings.Contains(nameServer, ":") {
			app.AddNameServers(strings.Split(nameServer, ":")[0])
		} else {
			utils.LogError("AppMaker-Controller-1", fmt.Errorf("GenDNS instance %s is of invalid format", nameServer))
		}
	}
	
	if pipeline[language] == nil {
		return nil, fmt.Errorf("language `%s` is not supported", language)
	}
	utils.LogError("AppMaker-Controller-1", fmt.Errorf("%s", app.GetDockerImage()))
	
	if app.GetDockerImage() != "" {
		utils.LogError("AppMaker-Controller-1", fmt.Errorf("docker img"))
		docker.CheckAndPullImages(app.GetDockerImage())
		containerPort, err := utils.GetFreePort()
		if err != nil {
			return nil, types.NewResErr(500, "No free port available", err)
		}
		app.SetContainerPort(containerPort)	

		errList := api.CreateBasicApplication(app)
		for _, err := range errList {
			if err != nil {
				return nil, err
			}
		}
	} else {
		utils.LogError("AppMaker-Controller-1", fmt.Errorf("git img"))
		resErr := pipeline[language].create(app)
		if resErr != nil {
			if resErr.Message() != "repository already exists" && resErr.Message() != "container already exists" {
				go diskCleanup(app.GetName())
			}
			return nil, fmt.Errorf(resErr.Error())
		}
	}
	sshEntrypointIP := configs.ServiceConfig.GenSSH.EntrypointIP
	if len(sshEntrypointIP) == 0 {
		sshEntrypointIP = utils.HostIP
	}
	app.SetSSHCmd(configs.ServiceConfig.GenSSH.Port, app.GetName(), sshEntrypointIP)

	app.SetAppURL(fmt.Sprintf("%s.%s.%s", app.GetName(), cloudflare.ApplicationInstance, configs.GasperConfig.Domain))

	if configs.CloudflareConfig.PlugIn {
		resp, err := cloudflare.CreateApplicationRecord(app.GetName())
		if err != nil {
			go diskCleanup(app.GetName())
			return nil, err
		}
		app.SetCloudflareID(resp.Result.ID)
		app.SetPublicIP(configs.CloudflareConfig.PublicIP)
	}

	err = mongo.UpsertInstance(
		types.M{
			mongo.NameKey:         app.GetName(),
			mongo.InstanceTypeKey: mongo.AppInstance,
		}, app)

	if err != nil && err != mongo.ErrNoDocuments {
		go diskCleanup(app.GetName())
		go stateCleanup(app.GetName())
		return nil, err
	}

	err = redis.RegisterApp(
		app.GetName(),
		fmt.Sprintf("%s:%d", utils.HostIP, configs.ServiceConfig.AppMaker.Port),
		fmt.Sprintf("%s:%d", utils.HostIP, app.GetContainerPort()),
	)

	if err != nil {
		go diskCleanup(app.GetName())
		go stateCleanup(app.GetName())
		return nil, err
	}

	err = redis.IncrementServiceLoad(
		ServiceName,
		fmt.Sprintf("%s:%d", utils.HostIP, configs.ServiceConfig.AppMaker.Port),
	)

	if err != nil {
		go diskCleanup(app.GetName())
		go stateCleanup(app.GetName())
		return nil, err
	}

	app.SetSuccess(true)

	response, err := json.Marshal(app)
	return &pb.ResponseBody{Data: response}, err
}

// Rebuild rebuilds an application
func (s *server) Rebuild(ctx context.Context, body *pb.NameHolder) (*pb.ResponseBody, error) {
	appName := body.GetName()
	app, err := mongo.FetchSingleApp(appName)
	if err != nil {
		return nil, err
	}

	pullChanges := []string{"git", "pull", "origin", app.GetGitRepositoryBranch()}
	_, err = docker.ExecProcess(app.ContainerID, pullChanges)

	if err != nil {
		return nil, err
	}

	response, err := json.Marshal(app)
	return &pb.ResponseBody{Data: response}, err
}

// Delete deletes an application
func (s *server) Delete(ctx context.Context, body *pb.NameHolder) (*pb.DeletionResponse, error) {
	appName := body.GetName()
	filter := types.M{
		mongo.NameKey:         appName,
		mongo.InstanceTypeKey: mongo.AppInstance,
	}

	node, _ := redis.FetchAppNode(appName)
	go redis.DecrementServiceLoad(ServiceName, node)
	go redis.RemoveApp(appName)
	go diskCleanup(appName)

	if configs.CloudflareConfig.PlugIn {
		go cloudflare.DeleteRecord(appName, mongo.AppInstance)
	}

	_, err := mongo.DeleteInstance(filter)
	if err != nil {
		return nil, err
	}
	return &pb.DeletionResponse{Success: true}, nil
}

// FetchLogs returns the docker container logs of an application
func (s *server) FetchLogs(ctx context.Context, body *pb.LogRequest) (*pb.LogResponse, error) {
	appName := body.GetName()
	tail := body.GetTail()

	data, err := docker.ReadLogs(appName, tail)

	if err != nil && err.Error() != "EOF" {
		return nil, err
	}
	return &pb.LogResponse{
		Success: true,
		Data:    data,
	}, nil
}

// NewService returns a new instance of the current microservice
func NewService() *grpc.Server {
	return factory.NewApplicationFactory(&server{})
}
