# IT Work Order System – MIVHS

This is a project of an IT MIVHS Work Order and Helpdesk system designed for the **TEAM IT MIVHS**.  
The goal is to provide a simple, fast, and user‑friendly interface for requests, helps, and any other thing related to devices on SMK MITRA INDUSTRI MM2100

## 🚀 Getting Started

To get the project up and running, follow these steps:

1. **Clone the repository:**
   ```bash
   git clone https://github.com/parothegreat/work-order.git
   ```

2. **Run the application:**
   The easiest way to run the application is using Docker Compose.
   ```bash
   docker-compose up -d
   ```
   This will start the following services:
   - `workorder-backend`: The Go backend API.
   - `db`: The PostgreSQL database.
   - `nginx`: The Nginx web server.

3. **Access the application:**
   Once the services are running, you can access the application in your browser at [http://localhost:8080](http://localhost:8080).

## 📌 Features

- Create and manage IT work orders  
- Display team member status and task assignment  
- Clean UI designed using TailwindCSS  
- Interactive UI powered by JavaScript
- User authentication and registration

## 🛠 Technologies Used

- **HTML5**
- **CSS3** (TailwindCSS + custom styles)
- **JavaScript** (vanilla)
- **Go** (Backend API)
- **Gin Framework** (Go web framework)
- **PostgreSQL** (Database)
- **Docker** (Containerization)
- **Nginx** (Web Server)

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
  Technician status: "Stand By" → "On Job"
```

### **3. Work in Progress (Executing Job)**
```
Order status: "On Progress"
  ↓
  Technician executes the work
  ↓
  Working hours are recorded
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
Teacher/Admin
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

---

## 👥 About the Project

This system was created for internal use by the **Web developer of IT MIVHS Team**.  
It's a complete **full-stack work order management system** with:

- ✅ Frontend UI (HTML, TailwindCSS, Vanilla JS)
- ✅ Backend API (Go + Gin Framework)
- ✅ Real-time status tracking
- ✅ Team member management
- ✅ Safety checklist system
- ✅ Performance evaluation (Kaizen)
- ✅ Audit trail & history tracking
- ✅ Responsive design

Future enhancements may include:
- Email / Telegram notifications  
- Advanced reporting & analytics
- Mobile app

## 📄 License

MIT License

Copyright (c) 2025 Team IT MIVHS

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHER DEALINGS IN THE
SOFTWARE.
