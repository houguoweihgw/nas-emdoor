import cv2
from SceneClassify import scene_classify
from torch.autograd import Variable as V
import torchvision.models as models
from torchvision import transforms as trn
from torch.nn import functional as F
import os
import signal
import numpy as np
from facenet_pytorch import MTCNN
from models.inception_resnet_v1 import InceptionResnetV1
import torch
from nats.aio.client import Client as NATS
import asyncio
import base64
import logging
import json
from PIL import Image
import io


def sc_model_init():
    # th architecture to use
    arch = 'resnet50'

    # load the pre-trained weights
    model_file = '%s_places365.pth.tar' % arch
    if not os.access(model_file, os.W_OK):
        weight_url = 'http://places2.csail.mit.edu/models_places365/' + model_file
        os.system('wget ' + weight_url)

    model = models.__dict__[arch](num_classes=365)
    checkpoint = torch.load(model_file, map_location=lambda storage, loc: storage)
    state_dict = {str.replace(k, 'module.', ''): v for k, v in checkpoint['state_dict'].items()}
    model.load_state_dict(state_dict)
    model.eval()

    # load the image transformer
    centre_crop = trn.Compose([
        trn.Resize((256, 256)),
        trn.CenterCrop(224),
        trn.ToTensor(),
        trn.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
    ])

    # load the class label
    file_name = 'categories_places365_Chinese.txt'
    classes = list()
    with open(file_name) as class_file:
        for line in class_file:
            classes.append(line.strip().split(' ')[0][3:])
    classes = tuple(classes)
    logging.info("Scene Classification Model Import Successfully")

    return centre_crop, classes, model


def fc_model_init():
    # 初始化MTCNN和InceptionResnetV1
    device = torch.device('cpu')
    mtcnn = MTCNN(
        image_size=160, margin=0, min_face_size=20,
        thresholds=[0.6, 0.7, 0.7], factor=0.709, post_process=True,
        device=device
    )
    resnet = InceptionResnetV1(pretrained='vggface2').eval().to(device)
    logging.info("Face Cluster Model Import Successfully")
    return device, mtcnn, resnet


def scene_classify(img, centre_crop, classes, model):
    input_img = V(centre_crop(img).unsqueeze(0))

    # forward pass
    logit = model.forward(input_img)
    h_x = F.softmax(logit, 1).data.squeeze()
    probs, idx = h_x.sort(0, True)

    # 获取前五个预测类别和概率
    top5_probs = [round(float(probs[i]), 3) for i in range(5)]
    top5_classes = [classes[idx[i]] for i in range(5)]

    # 使用列表推导筛选出probs大于0.1的个体
    filtered_probs = [prob for prob in top5_probs if prob > 0.1]
    filtered_classes = [top5_classes[i] for i, prob in enumerate(top5_probs) if prob > 0.1]
    return filtered_classes


async def main():
    logging.basicConfig(level=logging.INFO,  # 设置日志级别
                        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
    nc = NATS()
    centre_crop, classes, model = sc_model_init()
    device, mtcnn, resnet = fc_model_init()

    async def closed_cb():
        logging.info("Connection to NATS is closed.")
        await asyncio.sleep(0.1)
        asyncio.get_running_loop().stop()

    options = {
        "servers": "nats://nas-nats:4222",
        "closed_cb": closed_cb
    }

    logging.info("connecting")
    await nc.connect(**options)
    logging.info(f"Connected to NATS at {nc.connected_url.netloc}...")

    async def sc_subscribe_handler(msg):
        subject = msg.subject
        reply = msg.reply
        data_decode = msg.data.decode()
        if data_decode != 'ack':
            image_dict = json.loads(data_decode)
            image_id = image_dict['ID']
            user_id = image_dict['UserID']
            image_data = image_dict['Picture']
            logging.info("Received a message on '{subject} {reply}' ID {image_id} UserID {user_id}".format(
                subject=subject, reply=reply, image_id=image_id, user_id=user_id))
            # 将图片数据从 base64 解码并显示
            image = Image.open(io.BytesIO(base64.b64decode(image_data)))
            class_result = scene_classify(image, centre_crop, classes, model)
            # 在成功处理消息后发布确认消息
            await nc.publish("ack.Subject", b'ack')

            if len(class_result) > 0:
                logging.info(
                    "ID {ID} scene classify results:{class_result}".format(ID=image_id, class_result=class_result))
                res_dict = {'ID': image_id, 'UserID': user_id, 'TAGS': class_result}
                json_str = json.dumps(res_dict)
                byte_data = json_str.encode("utf-8")
                await nc.publish("Scene-Classification-Response", byte_data)

    async def fc_subscribe_handler(msg):
        subject = msg.subject
        reply = msg.reply
        data_decode = msg.data.decode()
        if data_decode != 'ack':
            image_dict = json.loads(data_decode)
            image_id = image_dict['ID']
            user_id = image_dict['UserID']
            image_data = image_dict['Picture']
            logging.info("Received a message on '{subject} {reply}' ID {image_id} UserID {user_id}".format(
                subject=subject, reply=reply, image_id=image_id, user_id=user_id))
            # 将图片数据从 base64 解码
            image_bytes = base64.b64decode(image_data)
            # 将字节数据转换为图像
            image = Image.open(io.BytesIO(image_bytes))
            # 将图像从PIL格式转换为NumPy数组，并转换为RGB格式
            image_np = np.array(image)
            image = cv2.cvtColor(image_np, cv2.COLOR_BGR2RGB)
            # 使用MTCNN检测人脸
            faces = mtcnn.detect(image)
            res_data = {
                'id': image_id,
                'user_id': user_id,
                'flag': False,
                'faces': []
            }
            data = []  # 创建一个空列表来存储数据

            for i, prob in enumerate(faces[1]):
                if prob == None:
                    logging.info('Detecting None faces')
                    json_str = json.dumps(res_data)
                    byte_data = json_str.encode("utf-8")
                    await nc.publish("Face-Cluster-Response", byte_data)
                    return
                if prob is not None and prob > 0.99:
                    box = faces[0][i]
                    x, y, w, h = [int(val) for val in box]
                    logging.info(f"Face {i + 1}: (x, y) = ({x}, {y}), (w, h) = ({w}, {h}), Probability: {prob}")
                    box = np.expand_dims(box, axis=0)  # 添加一个维度，使其变成1x4
                    x_aligned = mtcnn.extract(image, box, save_path=None)
                    aligned = []
                    if x_aligned is not None:
                        aligned.extend(x_aligned)

                    aligned = torch.stack(aligned).to(device)
                    aligned = aligned.unsqueeze(0)  # 在批处理维度上增加一个维度
                    embeddings = resnet(aligned).detach().cpu()
                    embeddings_np = embeddings.numpy()[0]
                    # 将数据存储为列表项
                    data.append({
                        'x': x,
                        'y': y,
                        'w': w,
                        'h': h,
                        'embeddings': embeddings_np.tolist()
                    })
                    # print(embeddings_np.tolist())
                    # 这里的embeddings是包含人脸特征的张量，可以根据需要使用它
            logging.info('Detecting %d faces and Extracting features.', len(data))
            res_data['flag'] = True
            res_data['faces'] = data

            json_str = json.dumps(res_data)
            byte_data = json_str.encode("utf-8")
            await nc.publish("Face-Cluster-Response", byte_data)

    # Basic subscription to receive all published messages
    # which are being sent to a single topic 'discover'
    await nc.subscribe("Scene-Classification-Request", cb=sc_subscribe_handler)
    await nc.subscribe("Face-Cluster-Request", cb=fc_subscribe_handler)

    def signal_handler():
        if nc.is_closed:
            return
        logging.info("Disconnecting...")
        asyncio.create_task(nc.close())

    for sig in ('SIGINT', 'SIGTERM'):
        asyncio.get_running_loop().add_signal_handler(getattr(signal, sig), signal_handler)

    # Keep the event loop running indefinitely
    while True:
        await asyncio.sleep(1)  # Add a delay to avoid a busy-wait loop


if __name__ == '__main__':
    try:
        asyncio.run(main())
    except:
        pass
