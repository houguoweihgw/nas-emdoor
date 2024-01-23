CREATE DATABASE IF NOT EXISTS nas_data;
use nas_data;

CREATE TABLE `users` (
                         `id` int NOT NULL AUTO_INCREMENT,
                         `username` varchar(255) NOT NULL,
                         `password` varchar(255) NOT NULL,
                         `email` varchar(255) NOT NULL,
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `Username` (`username`),
                         UNIQUE KEY `Email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `albums` (
                          `id` int NOT NULL AUTO_INCREMENT,
                          `name` varchar(255) NOT NULL,
                          `description` text,
                          `date_created` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                          `user_id` int DEFAULT NULL,
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `unique_user_name` (`user_id`,`name`),
                          CONSTRAINT `fk_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `items` (
                         `id` int NOT NULL AUTO_INCREMENT,
                         `title` varchar(255) NOT NULL,
                         `description` text,
                         `file_path` varchar(255) NOT NULL,
                         `date_uploaded` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                         `user_id` int DEFAULT NULL,
                         `status` enum('active','deleted','recycled') DEFAULT 'active',
                         `collected` tinyint(1) DEFAULT '0',
                         PRIMARY KEY (`id`),
                         KEY `ad_user_id` (`user_id`),
                         CONSTRAINT `ad_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `item_album` (
                              `id` int NOT NULL AUTO_INCREMENT,
                              `item_id` int NOT NULL,
                              `album_id` int NOT NULL,
                              `date_added` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
                              PRIMARY KEY (`id`),
                              KEY `item_id` (`item_id`),
                              KEY `album_id` (`album_id`),
                              CONSTRAINT `item_album_ibfk_1` FOREIGN KEY (`item_id`) REFERENCES `items` (`id`),
                              CONSTRAINT `item_album_ibfk_2` FOREIGN KEY (`album_id`) REFERENCES `albums` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `item_metadata` (
                                 `id` int NOT NULL AUTO_INCREMENT,
                                 `exposure_time` varchar(255) DEFAULT NULL,
                                 `aperture` float DEFAULT NULL,
                                 `iso` int DEFAULT NULL,
                                 `focal_length` float DEFAULT NULL,
                                 `latitude` float DEFAULT NULL,
                                 `longitude` float DEFAULT NULL,
                                 `altitude` float DEFAULT NULL,
                                 `make` varchar(255) DEFAULT NULL,
                                 `model` varchar(255) DEFAULT NULL,
                                 `date_taken` datetime DEFAULT NULL,
                                 `file_size` bigint DEFAULT NULL,
                                 `image_width` int DEFAULT NULL,
                                 `image_length` int DEFAULT NULL,
                                 `scene_tags` varchar(255) DEFAULT NULL,
                                 `user_id` int DEFAULT NULL,
                                 PRIMARY KEY (`id`),
                                 KEY `meta_user_id` (`user_id`),
                                 CONSTRAINT `meta_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


-- 创建中间表
CREATE TABLE `item_metadata_item` (
                                      `item_metadata_id` int NOT NULL,
                                      `item_id` int NOT NULL,
                                      PRIMARY KEY (`item_metadata_id`, `item_id`),
                                      FOREIGN KEY (`item_metadata_id`) REFERENCES `item_metadata` (`id`),
                                      FOREIGN KEY (`item_id`) REFERENCES `items` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE scene_labels (
                              id INT AUTO_INCREMENT PRIMARY KEY,
                              user_id INT NOT NULL,
                              label_name VARCHAR(255) NOT NULL,
                              label_count INT DEFAULT 0,
                              FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


CREATE TABLE `faces` (
                         `id` int NOT NULL AUTO_INCREMENT,
                         `item_id` int NOT NULL,
                         `x` int NOT NULL,
                         `y` int NOT NULL,
                         `w` int NOT NULL,
                         `h` int NOT NULL,
                         `features` blob,
                         `user_id` int DEFAULT NULL,
                         PRIMARY KEY (`id`),
                         KEY `item_id` (`item_id`),
                         KEY `user_id` (`user_id`),
                         CONSTRAINT `faces_ibfk_1` FOREIGN KEY (`item_id`) REFERENCES `items` (`id`),
                         CONSTRAINT `faces_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `clusters` (
                            `id` int NOT NULL AUTO_INCREMENT,
                            `user_id` int NOT NULL,
                            `name` varchar(255) NOT NULL,
                            `features` blob NOT NULL,
                            PRIMARY KEY (`id`),
                            UNIQUE KEY `unique_user_name` (`user_id`,`name`),
                            CONSTRAINT `clusters_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `cluster_faces` (
                                 `id` int NOT NULL AUTO_INCREMENT,
                                 `cluster_id` int NOT NULL,
                                 `face_id` int NOT NULL,
                                 PRIMARY KEY (`id`),
                                 KEY `cluster_faces_ibfk_1` (`cluster_id`),
                                 KEY `cluster_faces_ibfk_2` (`face_id`),
                                 CONSTRAINT `cluster_faces_ibfk_1` FOREIGN KEY (`cluster_id`) REFERENCES `clusters` (`id`),
                                 CONSTRAINT `cluster_faces_ibfk_2` FOREIGN KEY (`face_id`) REFERENCES `faces` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;