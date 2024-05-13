DROP DATABASE collage;

CREATE DATABASE IF NOT EXISTS collage;

USE collage;

CREATE TABLE Department (
    id INT PRIMARY KEY AUTO_INCREMENT,
    DepartmentName TEXT NOT NULL
);

CREATE TABLE Student (
    id INT PRIMARY KEY AUTO_INCREMENT,
    StudentName TEXT NOT NULL,
    DepartmentID INT
);

ALTER TABLE Student
ADD CONSTRAINT fk_Student_DepartmentID
FOREIGN KEY (DepartmentID) REFERENCES Department(id);

