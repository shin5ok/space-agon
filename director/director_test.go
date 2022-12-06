// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"io"
	"testing"

	"agones.dev/agones/pkg/client/clientset/versioned/fake"
	"github.com/golang/mock/gomock"
	mockpb "github.com/googleforgames/space-agon/pkg/testing/open-match/mockpb"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"open-match.dev/open-match/pkg/pb"
)

func Test_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFetchClient := mockpb.NewMockBackendService_FetchMatchesClient(ctrl)
	mockFetchClient.EXPECT().Recv().Return(
		&pb.FetchMatchesResponse{
			Match: &pb.Match{
				MatchId:       "100",
				MatchFunction: "match_function",
				MatchProfile:  "test_profile",
				Tickets: []*pb.Ticket{
					{Id: "foo"},
					{Id: "bar"},
				},
			},
		}, nil,
	)
	mockFetchClient.EXPECT().Recv().Return(nil, io.EOF)

	mockServiceClient := mockpb.NewMockBackendServiceClient(ctrl)
	mockServiceClient.EXPECT().FetchMatches(context.Background(), gomock.Any()).Return(mockFetchClient, nil)

	agonesMockClient := fake.NewSimpleClientset()

	func(t *testing.T, client *mockpb.MockBackendServiceClient) {
		t.Helper()
		r := Client{
			BackendServiceClient:       client,
			AgonesClientset:            agonesMockClient,
			CloserBackendServiceClient: func() error { return nil },
		}
		err := r.run()
		if err != nil {
			t.Error(err)
		}
	}(t, mockServiceClient)
}
