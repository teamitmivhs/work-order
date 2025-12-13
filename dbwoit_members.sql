-- MySQL dump 10.13  Distrib 8.0.44, for Win64 (x86_64)
--
-- Host: localhost    Database: dbwoit
-- ------------------------------------------------------
-- Server version	8.0.44

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `members`
--

DROP TABLE IF EXISTS `members`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `members` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `Name` varchar(255) NOT NULL,
  `Role` enum('programmer','maintenance','data analyst','soundman') DEFAULT NULL,
  `Status` enum('onjob','standby','support','nextshift','offduty') NOT NULL DEFAULT 'offduty',
  `Avatar` varchar(255) NOT NULL DEFAULT 'no avatar',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `members`
--

LOCK TABLES `members` WRITE;
/*!40000 ALTER TABLE `members` DISABLE KEYS */;
INSERT INTO `members` VALUES (1,'Aldi Fadlurahman R','programmer','offduty','aldi.png'),(2,'Raditya Ihsan Athallah','data analyst','offduty','raditya.png'),(3,'Azzam Alfarizhi','maintenance','offduty','azzam.png'),(4,'Ghani Ilham Firdaus','programmer','offduty','ghani.png'),(5,'Gilang Yoga Pangestu','maintenance','offduty','gilang.png'),(6,'Charis Nur Noveli Alfaridzi','programmer','offduty','charis.png'),(7,'M Akbar Putra P','soundman','offduty','akbar.png'),(8,'Fahri Abdul Ghoni','maintenance','offduty','ghoni.png'),(9,'Ridwan Bagoes Setiawan','programmer','offduty','ridwan.png'),(10,'Muhammad Alvaro','programmer','offduty','alvaro.png'),(11,'Dzaky Alvaro','programmer','offduty','dzaky.png'),(12,'Maynaldi Freza A','programmer','offduty','maynaldi.png'),(13,'Muhammad Riva Nugraha','programmer','offduty','riva.png'),(14,'Muhammad Farel Sustisna','programmer','offduty','farel.png'),(15,'Faa\'iz Rizqi Haryono','programmer','offduty','faa\'iz.png'),(16,'Royan Fadlan Musaminah','maintenance','offduty','royan.png'),(17,'Muhammad Reyhansyah Hidayat','maintenance','offduty','reyhansyah.png'),(18,'Naufal Abdilah Saputra','maintenance','offduty','naufal.png'),(19,'Rizki Nuraulia','data analyst','offduty','rizki.png'),(20,'Khalishah Althaf','data analyst','offduty','althaf.png'),(21,'Dzakiyya Najdatul Rameyza','data analyst','offduty','dzakiyya.png'),(22,'Desvita Aurellia','data analyst','offduty','desvita.png'),(23,'Qiara Latifah Kaltsum','data analyst','offduty','qiara.png'),(24,'Dzaki Mathoriq','soundman','offduty','thoriq.png'),(25,'Alvaro Purnomo','data analyst','offduty','purnomo.png'),(26,'Chelsea Aurelia','programmer','offduty','chelsea.png'),(27,'Nabil Hilmy Zaenal','maintenance','offduty','nabil.png'),(28,'Wildan Bait Maki','data analyst','offduty','wildan.png'),(29,'Gian Alvarezi Savatino Putra','programmer','offduty','gian.png'),(30,'Jonathan Willy','programmer','offduty','willy.png'),(31,'Pramadani Bintang Jasuma','maintenance','offduty','pramadani.png'),(32,'Aira Nur Sabariyah Putri','programmer','offduty','aira.png'),(33,'Azka Fakhri Alfito','programmer','offduty','azka.png'),(34,'Ihsan Bintang Ghifari','programmer','offduty','ihsan.png'),(35,'M. Dhimas Alfachry','maintenance','offduty','dhimas.png'),(36,'Raissya hanjani','data analyst','offduty','raissya.png');
/*!40000 ALTER TABLE `members` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-12-12 21:29:57
