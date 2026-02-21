CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE "GENDERTYPE" AS ENUM ('M', 'F', 'O');

CREATE TYPE "LEVELOFCOOKINGTYPE" AS ENUM ('Novice', 'Intermediate', 'Proficient', 'Expert');

CREATE TABLE IF NOT EXISTS "users"
(
    "ID" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR(100) NOT NULL,
    "email" VARCHAR(50) NOT NULL UNIQUE,
    "password" VARCHAR(500) NOT NULL,
    "gender" "GENDERTYPE" NOT NULL,
    "dob" DATE NOT NULL,
    "levelOfCooking" "LEVELOFCOOKINGTYPE" NOT NULL,
    "createdAtUTC" TIMESTAMP(3) NOT NULL DEFAULT NOW(),
    "updatedAtUTC" TIMESTAMP(3)
);
