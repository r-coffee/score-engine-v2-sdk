package score

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ScoreEngineSDKClient struct {
	stub ScoreServiceClient
}

// CreateScoreEngine creates a new sdk client
// host is the hostname of the server
// pathToCert is the path to the server's public certificate
// port is the port number the server service is running on
func CreateScoreEngine(host, pathToCert string, port int) *ScoreEngineSDKClient {
	var sdk ScoreEngineSDKClient
	creds, err := credentials.NewClientTLSFromFile(pathToCert, host)
	if err != nil {
		log.Fatal(err)
	}

	// connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(creds), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	sdk.stub = NewScoreServiceClient(conn)
	return &sdk
}

// TennisMatches returns a list of tennis matches for a given tournament
func (s *ScoreEngineSDKClient) TennisMatches(tournamentID int64, category, round string) ([]*TennisMatchTuple, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := s.stub.TennisMatches(ctx, &TennisMatchesRequest{TournamentID: tournamentID, Category: category, Round: round})
	if resp != nil {
		return resp.Matches, err
	}

	return nil, err
}
