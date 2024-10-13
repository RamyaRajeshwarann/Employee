/*
  Warnings:

  - You are about to drop the column `phoneNo` on the `employee` table. All the data in the column will be lost.
  - A unique constraint covering the columns `[phone]` on the table `employee` will be added. If there are existing duplicate values, this will fail.
  - Added the required column `phone` to the `employee` table without a default value. This is not possible if the table is not empty.

*/
-- DropIndex
DROP INDEX "employee"."employee_phoneNo_key";

-- AlterTable
ALTER TABLE "employee"."employee" DROP COLUMN "phoneNo",
ADD COLUMN     "phone" TEXT NOT NULL;

-- CreateIndex
CREATE UNIQUE INDEX "employee_phone_key" ON "employee"."employee"("phone");
