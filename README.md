# IT Work Order System – MIVHS

This is a project of an IT MIVHS Work Order and Helpdesk system designed for the **TEAM IT MIVHS**.  
The goal is to provide a simple, fast, and user‑friendly interface for requests, helps, and any other thing related to devices on SMK MITRA INDUSTRI MM2100

## 👥 About the Project

This system was created for internal use by the **Web developer of IT MIVHS Team**.  
It's a complete **full-stack microservices-based work order management system** with:

- ✅ **Frontend UI** (HTML, TailwindCSS, Vanilla JS)
- ✅ **Backend API** (Go + Gin Framework)
- ✅ **Rust Time Tracker Service** (Axum Framework)
- ✅ **MySQL Database** with health checks
- ✅ **Nginx Reverse Proxy**
- ✅ **Multi-service Docker Compose orchestration**
- ✅ **Real-time status tracking**
- ✅ **Team member management**
- ✅ **Safety checklist system**
- ✅ **Performance evaluation (Kaizen)**
- ✅ **Audit trail & history tracking**
- ✅ **Responsive design**
- ✅ **User authentication system**

## 🚀 Getting Started

To get the project up and running, follow these steps:

1. **Clone the repository:**
   ```bash
   git clone https://github.com/parothegreat/work-order.git
   cd work-order
   ```

2. **Run the application:**
   The easiest way to run the application is using Docker Compose.
   ```bash
   docker-compose up -d
   ```
   
   This will start the following services:
   - `db`: MySQL 8.0 database with automatic initialization
   - `time-tracker`: Rust-based microservice for time tracking
   - `work-order-backend`: Go backend API with Gin framework
   - `work-order-nginx`: Nginx reverse proxy serving the frontend

3. **Access the application:**
   Once the services are running, you can access the application in your browser:
   - **Main Application**: [http://localhost](http://localhost)
   - **Summary Page**: [http://localhost/summary](http://localhost/summary)
   - **Kaizen Page**: [http://localhost/kaizen](http://localhost/kaizen)
   - **TechGuide Page**: [http://localhost/techguide](http://localhost/techguide)

## 🏗️ Architecture

The system follows a microservices architecture:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Nginx         │    │   Go Backend     │    │  Rust Time      │
│   (Frontend)    │◄──►│   (API Server)   │◄──►│  Tracker        │
│   Port: 80      │    │   Port: 8080     │    │  Port: 9000     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌──────────────────┐
                       │   MySQL 8.0      │
                       │   Database       │
                       │   Port: 3306     │
                       └──────────────────┘
```

## 📌 Features

### Core Functionality
- **Work Order Management**: Create, assign, track, and complete IT work orders
- **Team Member Management**: Monitor technician status and availability
- **Real-time Status Tracking**: Live updates on work order progress
- ** Ensure compliance with safetySafety Checklist System**: protocols
- **Time Tracking**: Automatic work duration logging
- **Performance Analytics**: Kaizen dashboard for team evaluation

### Technical Features
- **Microservices Architecture**: Scalable service-based design
- **Health Checks**: Automatic service monitoring and recovery
- **Service Dependencies**: Proper startup sequencing between services
- **User Authentication**: Secure login/registration system
- **Responsive Design**: Mobile-friendly interface
- **Error Resilience**: Graceful degradation when database is unavailable

## 🛠 Technologies Used

### Frontend
- **HTML5** - Modern semantic markup
- **CSS3** (TailwindCSS + custom styles) - Utility-first CSS framework
- **JavaScript** (vanilla) - Interactive frontend logic

### Backend Services
- **Go** - Main backend API server
- **Gin Framework** - High-performance Go web framework
- **Rust** - Time tracking microservice for optimal performance
- **Axum** - Modern, fast web framework for Rust

### Database & Infrastructure
- **MySQL 8.0** - Primary database with auto-initialization
- **Docker** - Containerization platform
- **Docker Compose** - Multi-container orchestration
- **Nginx** - Reverse proxy and static file serving

## 🔄 Work Flow

This system follows a structured workflow to handle work orders efficiently:

### **1. Create Work Order (Requests from helpdesk)**
```
Requester/User
  ↓
  Click "Create Orders" button
  ↓
  Fill form:
    - Fill the requester name
    - Priority (High, Medium, Low)
    - Location (Gedung A, B, C, etc)
    - Device (Printer, PC, CCTV, etc)
    - Problem description
  ↓
  Submit → Order enters table with status "Pending"
```

### **2. Take Order (Assign Work)**
```
Technician
  ↓
  View work orders in main table
  ↓
  Click empty slot or order ID
  ↓
  - Select available operators (status: Stand By)
  - Review & approve safety checklist per location
  - Click "Confirm"
  ↓
  Status changes: "Pending" → "On Progress"
  IT Team status: "Stand By" → "On Job"
```

### **3. Work in Progress (Executing Job)**
```
Order status: "On Progress"
  ↓
  Technician executes the work
  ↓
  Working hours are recorded by Rust time tracker service
  ↓
  Can update team member status (Support, etc) if needed
```

### **4. Mark as Done (Complete Job)**
```
Technician
  ↓
  Click "Done" button in table row
  ↓
  (Optional) Fill evaluation notes:
    - What was done
    - Solution applied
    - Notes for improvement
  ↓
  Submit
  ↓
  Status changes: "On Progress" → "Completed"
  Technician status: "On Job" → "Stand By"
  Order enters Summary page
```

### **5. Review Summary (View History)**
```
IT Teacher/Admin
  ↓
  Click "Summary" hyperlink in navbar
  ↓
  View table: Completed work orders with:
    - Order ID, Priority, Time, Requester
    - Location, Device, Problem
    - Completion timestamp
    - Evaluation notes
  ↓
  Can edit notes for additional feedback
  ↓
  Analyze data for Kaizen improvement
```

### **6. Kaizen Activity (Performance Evaluation)**
```
Manager can view metrics:
  - Total work orders
  - Pending orders
  - On progress orders
  - Completed orders
  
Completion Rate = (Completed / Total) × 100%

Rating based on completion rate:
  - Excellent (80%+) → "Keep up the good work"
  - Good (60-79%) → "Focus on reducing pending"
  - Fair (40-59%) → "Consider process improvements"
  - Needs Improvement (<40%) → "Investigate bottlenecks"
```

## 🐳 Docker Services

The system consists of four main Docker services:

### Database Service (`db`)
- **Image**: `mysql:8.0`
- **Auto-initialization**: Runs SQL scripts from `./db` directory
- **Health Check**: Automatic connection testing
- **Volumes**: Persistent data storage

### Time Tracker Service (`time-tracker`)
- **Built from**: `./rust-engine` directory
- **Framework**: Rust with Axum
- **Purpose**: High-performance time tracking microservice
- **Port**: 9000

### Backend Service (`backend`)
- **Built from**: `./backend` directory
- **Framework**: Go with Gin
- **Dependencies**: Waits for database and time-tracker services
- **Port**: 8080

### Nginx Service (`nginx`)
- **Image**: `nginx:1.25-alpine`
- **Purpose**: Reverse proxy and static file serving
- **Port**: 80 (mapped to host)
- **Dependencies**: Starts after backend service

## 🔧 Development

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for backend development)
- Rust 1.70+ (for time tracker development)

### Manual Development Setup

#### Backend Development
```bash
cd src/backend
go mod tidy
go run main.go
```

#### Time Tracker Development
```bash
cd src/rust-engine
cargo build
cargo run
```

#### Database Setup
```bash
# MySQL database with auto-initialization scripts
# Scripts located in: src/db/
# - dbwoit_orders.sql
# - dbwoit_members.sql
# - dbwoit_executors.sql
# - dbwoit_safetychecklist.sql
```

## 📊 Database Schema

The system uses MySQL with the following main tables:
- **orders**: Work order records
- **members**: Team member information
- **executors**: Work assignment tracking
- **safety_checklist**: Safety compliance records

## 🔐 Authentication

The system includes user authentication with:
- **Registration**: New user account creation
- **Login**: Secure user authentication
- **Session Management**: JWT-based session handling
- **Role-based Access**: Different access levels for team members

## 📱 Mobile Support

The system is fully responsive and includes:
- **Collapsible Navigation**: Hamburger menu for mobile devices
- **Responsive Tables**: Horizontal scroll for work order data
- **Mobile-friendly Forms**: Optimized input forms
- **Touch-friendly Interface**: Optimized for touch interactions

## 🔍 Monitoring & Health Checks

The system includes built-in monitoring:
- **Database Health Checks**: Automatic connection monitoring
- **Service Dependencies**: Proper startup sequencing
- **Graceful Degradation**: Continues operating even without database
- **Local Storage Fallback**: Frontend data persistence when backend unavailable

## 🚀 Future Enhancements

Planned improvements:
- Email/Telegram notifications
- Advanced reporting & analytics
- Mobile application
- Real-time notifications via WebSocket
- Advanced role-based permissions
- Integration with external IT tools

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Developed with ❤️ by IT MIVHS Team**  
*Empowering efficient IT work order management*
