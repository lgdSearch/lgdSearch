package vgg

import (
    "lgdSearch/pkg/vgg/search"
	"context"
    "time"
    "google.golang.org/grpc"
)


func Search(img []byte) ([][]byte, error) {
	conn, err := grpc.Dial("101.42.175.203:50052", grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    defer conn.Close()

    client := search.NewGrpcServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()

    resp, err := client.Search(ctx, &search.Request{
        Image: img,
    })
    if err != nil {
        return nil, err
    }
    feature := make([][]byte, len(resp.Images))
    copy(feature, resp.Images)
    return feature, nil
}

