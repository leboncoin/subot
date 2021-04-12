/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package engine_grpc_client

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"

	"google.golang.org/grpc"
	pb "github.com/leboncoin/subot/pkg/engine_grpc_client/engine"
)

// Engine represents the engine client instance
type Engine struct {
	Client pb.EngineClient `json:"connection"`
}

// IEngine represents the EngineClient interface to ease mocking
type IEngine interface {
	AnalyseMessageTools(text *pb.Text) ([]pb.Category, error)
	AnalyseMessageLabels(text *pb.Text) ([]pb.Category, error)
}

// AnalyseMessageLabels gets the labels associated with the given message text.
func (e Engine) AnalyseMessageLabels(text *pb.Text) ([]pb.Category, error) {
	log.Printf("Getting labels for text")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	labels, err := e.Client.AnalyseMessageLabels(ctx, text)
	if err != nil {
		log.Errorf("AnalyseMessageLabels(_) = _, %v: ", err)
		return nil, err
	}
	var bestCats []pb.Category
	for _, cat := range labels.Categories{
		if cat.Score > 0.5 {
			bestCats = append(bestCats, *cat)
		}
	}

	return bestCats, nil
}

// AnalyseMessageTools gets the tools associated with the given message text.
func (e Engine) AnalyseMessageTools(text *pb.Text) ([]pb.Category, error) {
	log.Printf("Getting tools for text")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	tools, err := e.Client.AnalyseMessageTools(ctx, text)
	if err != nil {
		log.Errorf("AnalyseMessageTools(_) = _, %v: ", err)
		return nil, err
	}

	var bestCats []pb.Category
	for _, cat := range tools.Categories{
		if cat.Score > 0.5 {
			bestCats = append(bestCats, *cat)
		}
	}

	return bestCats, nil
}

// Client returns an instance of a connected engine
func Client(target string, opts []grpc.DialOption) (Engine, error) {
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		log.Errorf("fail to dial: %v", err)
		return Engine{}, err
	}
	client := pb.NewEngineClient(conn)
	return Engine{Client: client}, nil
}
