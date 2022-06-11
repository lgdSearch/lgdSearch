# 以图搜图

## python部分
使用ResNet50提取图片特征向量
预处理数据集中的图片，将图片id与特征向量存入milvus
id与图片路径的映射使用diskcache持久化到文件
使用grpc提供搜索服务，对于传入的图片，提取特征值后与milvus中的向量比较相似度

## golang部分
使用grpc调用python侧提供的服务