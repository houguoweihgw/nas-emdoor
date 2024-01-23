import torch
from torch.autograd import Variable as V
import torchvision.models as models
from torchvision import transforms as trn
from torch.nn import functional as F
import os
from PIL import Image, ImageDraw, ImageFont


def scene_classify(img):
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

    # load the test image
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
    print(filtered_classes)
    # for i in range(len(filtered_probs)):
    #     text = '{}: {}'.format(filtered_classes[i], filtered_probs[i])
    #     print(text)


