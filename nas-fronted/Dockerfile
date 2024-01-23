FROM nginx:latest

# RUN npm install && npm run build

# 将dist文件中的内容复制到 /usr/share/nginx/html/ 这个目录下面
COPY dist/ /usr/share/nginx/html/

#添加自己的配置 default.conf
#COPY nginx.conf /etc/nginx/nginx.conf
# 将nginx配置放到指定的位置
COPY default.conf /etc/nginx/conf.d/default.conf
# 暴露 Nginx 默认端口
EXPOSE 80
# 容器启动时运行 Nginx
CMD ["nginx", "-g", "daemon off;"]