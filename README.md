<h1 align="center">
    Dsync
</h1>

**Dsync** is a CLI tool designed to simplify the management and synchronization of data between diverse systems, including databases, cloud storage, and local file systems. It addresses real-world challenges in data migration, automation, and backup by offering secure synchronization, data transformation, scheduling, and task management features.

---

## Key Features

### 1. Multi-Source Data Synchronization
- **Database Sync**: Synchronize data between databases (e.g., MySQL ↔ MySQL).
- **Local to Cloud**: Transfer data between local systems and cloud services effortlessly.

### 2. Backup and Restore
- **Schema-Only or Data-only**: Create backups for databases with only data or only schema of tables inside database
- **Scheduled Backups**: Create backups for databases and files on a schedule.

### 3. Tracking and Logging
- **Detailed Logs**: Monitor status and troubleshoot with comprehensive logs.
- **Retry Mechanisms**: Track and retry failed tasks automatically.

---

# Project: dsync

This project, `dsync`, is built using the Go programming language. Below is a list of the main technologies, libraries, and dependencies used in this project.

## Technologies and Libraries Used

### Go Version
- **Go**: 1.23.4

### Main Dependencies

- **Cobra**: v1.8.1  
  A library for creating powerful modern CLI applications.

- **Termlink**: v1.4.2  
  A library for easy handling of terminal links (hyperlinks) in terminal environments.

- **GORM**:  
  - **GORM ORM**: v1.25.12  
    The ORM library for Go for interacting with databases.
  - **MySQL Driver**: v1.5.7  
    MySQL driver for GORM.
  - **PostgreSQL Driver**: v1.5.11  
    PostgreSQL driver for GORM.

- **Godotenv**: v1.5.1  
  A Go library for loading environment variables from `.env` files.

- **Color**: v1.5.4  
  A library for easy terminal color manipulation and output formatting.

- **Tablewriter**: v0.0.5  
  A Go package for rendering tables in the terminal with customizable options.

