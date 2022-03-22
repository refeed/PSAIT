CREATE DATABASE IF NOT EXISTS `student_kasus3`;

CREATE TABLE IF NOT EXISTS `student_kasus3`.`student` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`));

INSERT INTO `student_kasus3`.`student` (`name`) VALUES ('Aldi');
INSERT INTO `student_kasus3`.`student` (`name`) VALUES ('Budi');
INSERT INTO `student_kasus3`.`student` (`name`) VALUES ('Caca');
INSERT INTO `student_kasus3`.`student` (`name`) VALUES ('Dodi');
