import os
import tensorflow as tf
from keras.applications.resnet import ResNet50
from numpy import linalg as LA
from milvus_operator import insert, drop_image
from diskcache import Cache
from milvus_operator import get_collection_image


def parse_image(path):
    cache = Cache('./cache')
    imgs = os.listdir(path)
    imgs.sort()
    model = ResNet50(include_top=False, pooling="avg")
    data = [[], []]
    id = 1
    for img_path in imgs:
        img_path = path + '\\' + img_path
        print(img_path)
        img = tf.keras.utils.load_img(img_path, target_size=(224, 224))
        img_array = tf.keras.utils.img_to_array(img)
        img_array = tf.expand_dims(img_array, 0)  # Create a batch
        feat = model.predict(img_array)
        norm_feat = feat[0] / LA.norm(feat[0])
        data[0].append(id)
        data[1].append(norm_feat)
        cache[id] = img_path
        id += 1
    insert(data)
    cache.close()


if __name__ == '__main__':
    parse_image('D:\\chromeDownload\\VOCdevkit\\VOC2012\\temp')