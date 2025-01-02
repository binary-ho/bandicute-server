package main

import (
	"bandicute-server/config"
	"bandicute-server/internal/api"
	"bandicute-server/internal/service"
	"bandicute-server/internal/service/channel"
	"bandicute-server/internal/storage/repository/connection"
	"bandicute-server/internal/storage/repository/member"
	"bandicute-server/internal/storage/repository/post"
	pullRequest "bandicute-server/internal/storage/repository/pull-request"
	"bandicute-server/internal/storage/repository/study"
	studyMember "bandicute-server/internal/storage/repository/study-member"
	"bandicute-server/internal/storage/repository/summary"
	"bandicute-server/internal/task"
	"bandicute-server/internal/util"
	"bandicute-server/pkg/logger"
	"fmt"
	"strconv"
)

func main() {
	// 설정 로드
	config := config.GetInstance()

	// Repository 초기화
	dbConnection := connection.NewConnection(config.Database.BaseURL, config.Database.Key)
	memberRepo := member.NewMemberRepository(dbConnection)
	studyRepo := study.NewStudyRepository(dbConnection)
	postRepo := post.NewPostRepository(dbConnection)
	summaryRepo := summary.NewPostWriterRepository(dbConnection)
	studyMemberRepo := studyMember.NewStudyMemberRepository(dbConnection)
	pullRequestRepo := pullRequest.NewPullRequestRepository(dbConnection)

	// Util 초기화
	postParser := util.NewPostParser()
	postSummarizer, err := util.NewPostSummarizer(config.OpenAI.APIKey)
	if err != nil {
		logger.Fatal("Failed to initialize post summarizer", logger.Fields{
			"error": err.Error(),
		})
		panic(err)
	}

	gitHubService, err := util.NewGitHubService(config.GitHub.Token)
	if err != nil {
		logger.Fatal("Failed to initialize GitHub service", logger.Fields{
			"error": err.Error(),
		})
		panic(err)
	}

	// 서비스 초기화
	parser := service.NewParser(postParser, memberRepo, postRepo, summaryRepo)
	summarizer := service.NewSummarizer(postSummarizer, memberRepo, summaryRepo, studyRepo, pullRequestRepo)
	pullRequestOpener := service.NewPullRequestOpener(gitHubService, pullRequestRepo)

	// 채널 초기화
	parsePostByMemberIdRequestChannel := make(channel.ParsePostByMemberIdRequest)
	summarizeRequestChannel := make(channel.SummarizeRequest)
	openPullRequestRequestChannel := make(channel.OpenPullRequestRequest)

	// 태스크 핸들러 초기화 및 실행
	taskHandler := task.NewHandler(
		parser,
		summarizer,
		pullRequestOpener,
		&parsePostByMemberIdRequestChannel,
		&summarizeRequestChannel,
		&openPullRequestRequestChannel,
	)
	taskHandler.Run()

	// logger setup
	logger.Setup(config.Logging.Level)

	// API 서버 실행
	writer := service.NewWriter(studyMemberRepo, &parsePostByMemberIdRequestChannel)
	app := api.NewApplication(writer)

	fiberApp := app.Routes()
	err = fiberApp.Listen(getStringPort(config.Server.Port))
	if err != nil {
		logger.Fatal("Server Error", logger.Fields{
			"error": err.Error(),
		})
	}
	fmt.Println("exit..")
}

func getStringPort(port int) string {
	return ":" + strconv.Itoa(port)
}
