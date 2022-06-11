package vgg

import (
    "lgdSearch/pkg/vgg/search"
	"context"
    "time"
    "google.golang.org/grpc"
)


func Search(img []byte) ([]uint32, error) {
	conn, err := grpc.Dial("127.0.0.1:50052", grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    defer conn.Close()

    client := search.NewGrpcServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
    defer cancel()

    resp, err := client.Search(ctx, &search.Request{
        Image: img,
    })
    if err != nil {
        return nil, err
    }
    feature := make([]uint32, len(resp.Ids))
    copy(feature, resp.Ids)
    return feature, nil
}

