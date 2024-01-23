package utils

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"nas-backend/entity"
	"time"
)

var (
	NC                *nats.Conn
	SCPubSubject      string
	SCSubSubject      string
	FCPubSubject      string
	FCSubSubject      string
	requestAckChannel = make(chan bool)
)

// NatsInit NATS初始化
func NatsInit() {
	// 连接到 NATS 服务器
	//nc, err := nats.Connect(nats.DefaultURL)
	nc, err := nats.Connect("http://nas-nats:4222")
	if err != nil {
		log.Fatalf("NATS: %v", err)
	}
	log.Printf("NATS: %v \n", nc.Status())
	//defer nc.Close()
	// 定义场景分类发布的主题
	SCPubSubject = "Scene-Classification-Request"
	// 定义场景分类订阅的主题
	SCSubSubject = "Scene-Classification-Response"
	// 定义人脸聚类发布的主题
	FCPubSubject = "Face-Cluster-Request"
	// 定义人脸聚类订阅的主题
	FCSubSubject = "Face-Cluster-Response"
	NC = nc
	// 启动一个 Go 协程来订阅响应消息
	go scPublishAckResponse()
	// 启动一个 Go 协程来订阅场景分类结果消息
	go NATSSubscribe()
	// 启动一个 Go 协程来订阅人脸聚类结果消息
	go NATSFCSubscribe()
}

// scPublishAckResponse
func scPublishAckResponse() {
	_, err := NC.Subscribe("ack.Subject", func(m *nats.Msg) {
		// 收到算法端的响应
		log.Printf("收到算法端的响应: %s", m.Data)
		// 向 responseChannel 发送一个响应，表示已收到响应
		requestAckChannel <- true
	})
	if err != nil {
		log.Fatalf("订阅响应主题失败: %v", err)
	}
}

// NATSPublish NATS发布场景分类请求消息
func NATSPublish(message []byte) error {
	// 发布消息
	if err := NC.Publish(SCPubSubject, message); err != nil {
		log.Fatalf("发布消息失败: %v", err)
		return err
	}
	log.Printf("发布请求场景分类消息成功 \n")

	// todo:消息确认机制ack
	// 设置一个超时定时器，等待算法端响应
	timeout := time.NewTimer(5 * time.Second) // 假设超时时间是10秒

	select {
	case <-timeout.C:
		// 超时处理
		log.Println("等待响应超时，未收到算法端的响应")
		// todo:可以在这里执行相应的超时处理逻辑，如重试或记录错误
	case <-requestAckChannel:
		// 收到响应，不执行超时处理
		log.Println("收到算法端的响应")
	}
	return nil
}

// NATSFCPublish NATS发布人脸聚类请求消息
func NATSFCPublish(message []byte) error {
	// 发布消息
	if err := NC.Publish(FCPubSubject, message); err != nil {
		log.Fatalf("发布消息失败: %v", err)
		return err
	}
	log.Printf("发布请求人脸聚类消息成功 \n")

	// todo:消息确认机制ack
	// 设置一个超时定时器，等待算法端响应
	timeout := time.NewTimer(5 * time.Second) // 假设超时时间是10秒

	select {
	case <-timeout.C:
		// 超时处理
		log.Println("等待响应超时，未收到算法端的响应")
		// todo:可以在这里执行相应的超时处理逻辑，如重试或记录错误
	case <-requestAckChannel:
		// 收到响应，不执行超时处理
		log.Println("收到算法端的响应")
	}
	return nil
}

// NATSSubscribe NATS接收场景分类结果消息
func NATSSubscribe() {
	// 订阅主题
	sub, err := NC.Subscribe(SCSubSubject, func(m *nats.Msg) {
		// 解码字节数据
		var scResult entity.SCResponse
		if err := json.Unmarshal(m.Data, &scResult); err != nil {
			log.Printf("无法解码场景分类响应数据: %v", err)
			return
		}
		id := scResult.ID
		userID := scResult.UserID
		tags := scResult.TAGS
		sceneTags := ConvertTagsToString(tags)
		log.Printf("订阅到场景分类结果消息,ID: %d ,TAGS: %s", id, sceneTags)
		err := entity.UpdateItemMetadataSceneTags(DB, id, sceneTags)
		if err != nil {
			log.Printf("更新数据库失败: %v", err)
		}
		entity.ParseSceneClassificationLabels(DB, userID, tags)
	})
	if err != nil {
		log.Fatalf("订阅主题失败: %v", err)
		return
	}

	defer func(sub *nats.Subscription) {
		err := sub.Unsubscribe()
		if err != nil {

		}
	}(sub)

	// 使用通道等待退出信号
	exitSignal := make(chan struct{})
	go func() {
		select {
		// 此处可以添加其他退出条件，如收到特定信号
		case <-exitSignal:
			log.Println("退出订阅")
		}
	}()

	// 阻塞主函数，等待退出信号
	<-exitSignal
}

// NATSFCSubscribe NATS接收人脸聚类结果消息
func NATSFCSubscribe() {
	// 订阅主题
	sub, err := NC.Subscribe(FCSubSubject, func(m *nats.Msg) {
		// 解码字节数据
		var fcResult entity.FCResponse
		if err := json.Unmarshal(m.Data, &fcResult); err != nil {
			log.Printf("无法解码人脸聚类响应数据: %v", err)
			return
		}
		if fcResult.Flag {
			for _, face := range fcResult.Faces {
				//新建face记录
				face, err := entity.CreateFaces(DB, face, fcResult.UserID, fcResult.ID)
				if err != nil {
					log.Printf("人脸数据插入失败: %v", err)
					return
				}
				//face聚类
				err = entity.InsertFaceIntoCluster(DB, fcResult.UserID, face)
				if err != nil {
					log.Printf("人脸分类插入失败: %v", err)
					return
				}
			}
		}
	})
	if err != nil {
		log.Fatalf("订阅人脸聚类响应主题失败: %v", err)
		return
	}

	defer func(sub *nats.Subscription) {
		err := sub.Unsubscribe()
		if err != nil {

		}
	}(sub)

	// 使用通道等待退出信号
	exitSignal := make(chan struct{})
	go func() {
		select {
		// 此处可以添加其他退出条件，如收到特定信号
		case <-exitSignal:
			log.Println("退出人脸聚类响应订阅")
		}
	}()

	// 阻塞主函数，等待退出信号
	<-exitSignal
}
