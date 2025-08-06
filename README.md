# 📚 Bookstore App (Go + Gin + GORM + PostgreSQL)

A clean and scalable **Bookstore Web API** built with **Go**, following **Clean Architecture** principles.

---

## 🚀 Features

- ✅ **User Signup & Login** with **JWT Authentication**
- 📖 **Browse Books** with **Search & Multiple Filters**
- 🛒 **Place Orders** (Authenticated Users Only)
- 📧 **Email Notifications** on Successful Orders
- 🧼 **Clean Architecture** (Domain-Driven, Decoupled)
- 🗄️ **PostgreSQL** + **GORM** for Database ORM
- 🔐 Secure Password Hashing
- 🌐 RESTful APIs

---

## 🧰 Tech Stack

- **Go 1.21+**
- **Gin Web Framework**
- **GORM** (PostgreSQL ORM)
- **JWT** Authentication
- **PostgreSQL**
- **net/smtp** or **Gomail** for Email
- **bcrypt** for Password Hashing

---

## ⚙️ Getting Started

### 🔧 Prerequisites

- Go installed
- PostgreSQL running
- SMTP email credentials (for notifications)

### 🐘 Create PostgreSQL DB

```sql
CREATE DATABASE bookstore_db;
