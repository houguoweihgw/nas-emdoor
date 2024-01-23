 FROM python:3.7-slim

 WORKDIR /nas-algorithon

 COPY . .

 RUN cp sources.list /etc/apt/sources.list
 RUN apt-get update

 RUN apt install libgl1-mesa-glx -y

 RUN apt install libglib2.0-dev -y

 RUN pip install -i https://pypi.tuna.tsinghua.edu.cn/simple -r requirements.txt

 CMD ["python", "main.py"]


#FROM houguowei/nas-algorithm:v1.2
## 设置工作目录
#WORKDIR /nas-algorithon
## 复制你的 Python 应用文件到工作目录
#COPY . .
## 定义容器启动时执行的命令
#CMD ["python", "main.py"]
