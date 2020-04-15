CREATE DATABASE IF NOT EXISTS friendMgmt;
USE friendMgmt;

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `Id` int NOT NULL AUTO_INCREMENT,
  `Email` varchar(24) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user`
--

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'johndoe@gmail.com'),(2,'janedoe@gmail.com'),(3,'kytruong@yahoo.com'),(4,'example@email.com'),(5,'abc123@gmail.com'),(6,'a@gmail.com'),(7,'b@gmail.com'),(8,'c@gmail.com'),(9,'d@gmail.com'),(10,'e@gmail.com'),(11,'sangdepchai@gmail.com'),(12,'kytruong@gmail.com'),(13,'123@email.com'),(14,'1@email.com');
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;

DROP TABLE IF EXISTS `relationship`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `relationship` (
  `Id` int NOT NULL AUTO_INCREMENT,
  `RequestUserId` int NOT NULL,
  `TargetUserId` int NOT NULL,
  `Status` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`Id`),
  KEY `IX_Relationship_RequestUserId` (`RequestUserId`),
  KEY `IX_Relationship_TargetUserId` (`TargetUserId`),
  CONSTRAINT `FK_Relationship_User_RequestUserId` FOREIGN KEY (`RequestUserId`) REFERENCES `user` (`Id`),
  CONSTRAINT `FK_Relationship_User_TargetUserId` FOREIGN KEY (`TargetUserId`) REFERENCES `user` (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `relationship`
--

LOCK TABLES `relationship` WRITE;
/*!40000 ALTER TABLE `relationship` DISABLE KEYS */;
INSERT INTO `relationship` VALUES (37,1,2,3),(38,2,5,1),(39,5,1,1),(40,2,7,1),(41,3,5,1),(42,6,7,3),(43,7,6,2);
/*!40000 ALTER TABLE `relationship` ENABLE KEYS */;
UNLOCK TABLES; 