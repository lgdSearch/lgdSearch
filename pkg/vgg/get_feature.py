import get_feature_pb2
import get_feature_pb2_grpc
from keras.applications.vgg16 import VGG16
import tensorflow as tf
from numpy import linalg as LA
import grpc
from concurrent import futures
import time
import io
from PIL import Image
from pymilvus import connections


class GetFeatureService(get_feature_pb2_grpc.GrpcServiceServicer):
    def __init__(self):
        self.model = VGG16(weights='imagenet', include_top=False)

    def getFeature(self, request, context):
        # img_bytes: 图片内容 bytes
        img = Image.open(io.BytesIO(request.image))
        print(img)
        img = img.convert('RGB')
        img = img.resize((224, 224), Image.NEAREST)
        # img_path = request.image
        # img = tf.keras.utils.load_img(img_path, target_size=(224, 224))
        img_array = tf.keras.utils.img_to_array(img)
        img_array = tf.expand_dims(img_array, 0)  # Create a batch
        feat = self.model.predict(img_array)
        norm_feat = feat[0] / LA.norm(feat[0])
        response = get_feature_pb2.Response()
        for i in range(0, len(norm_feat)):
            add_resp = response.c.add()
            for j in range(0, len(norm_feat[i])):
                add_r1 = add_resp.b.add()
                add_r1.a.extend(norm_feat[i][j])
        return response


def run():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    get_feature_pb2_grpc.add_GrpcServiceServicer_to_server(GetFeatureService(), server)
    server.add_insecure_port('127.0.0.1:50052')
    server.start()
    print("start service...")
    try:
        while True:
            time.sleep(60 * 60 * 24)
    except KeyboardInterrupt:
        server.stop(0)


def connect2Milvus():
    connections.connect(
        alias="default",
        host='192.168.242.5',
        port='19530'
    )


def create_collection():
    


if __name__ == '__main__':
    # run()
    connect2Milvus()
