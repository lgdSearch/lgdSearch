import search_pb2
import search_pb2_grpc
from keras.applications.resnet import ResNet50
import tensorflow as tf
from numpy import linalg as LA
import numpy as np
import grpc
from concurrent import futures
import time
import io
from PIL import Image
from milvus_operator import search


class SearchService(search_pb2_grpc.GrpcServiceServicer):
    def __init__(self):
        self.model = ResNet50(include_top=False, pooling="avg")

    def search(self, request, context):
        # img_bytes: 图片内容 bytes
        img = Image.open(io.BytesIO(request.image))
        img = img.convert('RGB')
        img = img.resize((224, 224), Image.NEAREST)
        # img_path = request.image
        # img = tf.keras.utils.load_img(img_path, target_size=(224, 224))
        img_array = tf.keras.utils.img_to_array(img)
        img_array = tf.expand_dims(img_array, 0)  # Create a batch
        feat = self.model.predict(img_array)
        norm_feat = feat[0] / LA.norm(feat[0])
        ids = search(norm_feat)
        print(ids)
        resp = search_pb2.Response(ids=ids)
        print(resp)
        return resp


def run():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    search_pb2_grpc.add_GrpcServiceServicer_to_server(SearchService(), server)
    server.add_insecure_port('127.0.0.1:50052')
    server.start()
    print("start service...")
    try:
        while True:
            time.sleep(60 * 60 * 24)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    run()