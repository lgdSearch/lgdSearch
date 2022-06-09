package vgg

import (
	"lgdSearch/pkg/vgg/getfeature"
	"context"
    "time"
    "google.golang.org/grpc"
)

func GetFeature(img []byte) ([][][]float32, error) {
	conn, err := grpc.Dial("127.0.0.1:50052", grpc.WithInsecure())
    if err != nil {
        return nil, err
    }
    defer conn.Close()

    client := getfeature.NewGrpcServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
    defer cancel()

    resp, err := client.GetFeature(ctx, &getfeature.Request{
        Image: img,
    })
    if err != nil {
        return nil, err
    }
    feature := make([][][]float32, 0)
    for _, x := range resp.C {
        elem := make([][]float32, 0)
        for _, y := range x.B {
            elem = append(elem, y.A)
        }
        feature = append(feature, elem)
    }
    return feature, nil
}

