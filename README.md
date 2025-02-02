<h1 align="center">
    xnap
</h1>

**xnap**  is a command-line interface (CLI) tool designed for easy and efficient management of databases and local file backups. It offers a suite of features for maintaining database integrity, performing regular backups, restoring databases and files, and logging all actions to ensure transparency and accountability. This tool is ideal for developers, system administrators, and IT professionals who need reliable backup and restore mechanisms with clear logging capabilities.

---

## Key Features

### 1. Backup and Restore
- **Schema-Only or Data-only dump file**: Create backups dump file either complete dump file or data only or schema only 
- **Scheduled Backups**: Create backups for databases and local files on a schedule.
- **Scheduled Restore**: Create Restores for databases and local files on a schedule.

### 2.Database Management
  - Allows managing databases and Logs through simple CLI commands ie `CURDL`.
  - Supports connecting to multiple database types and systems.
  
### 3. Logging And Alerting
- **Detailed Logs**: Monitor status and troubleshoot with comprehensive logs ie, includes timestamps, error messages, and additional context for troubleshooting.
- **Alert Mechanism**: Tracks both successful and failed backup and restore processes and sends email alert in case of failure backup/restore process to ensure immediate action can be taken.
- **Retry Mechanisms**: Track and retry failed tasks automatically.

---

# Project: xnap

This project, `xnap`, is built using the Go programming language. Below is a list of the main technologies, libraries, and dependencies used in this project.

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

- **SMTP Server**:
  A server for sending, receiving, and relaying outgoing emails between mail servers

---

## Environment Variables

```env
# Server
DB_HOST=
PORT=

# MySQLDatabase
MySQL_DB_HOST=
MySQL_DB_PORT=

#PostgresDatabase
POSTGRES_DB_HOST=
POSTGRES_DB_PORT=

#xnap DB
XNAP_DB=
XNAP_LOGS_TABLE=
XNAP_BACKUP_TABLE=

#SMTP for Gmail
FROM_EMAIL=
FROM_EMAIL_PASSWORD="use passkey"
FROM_EMAIL_SMTP="smtp.gmail.com"
SMTP_ADDR=""smtp.gmail.com:587"
OWNER_EMAIL=

#Variables
BACKUP_OR_RESTORE_STATUS="failure"
```

## xnap Database Schemas 

### MySQL Tables Schema

- **Logs**

```sql
CREATE TABLE logs (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    action VARCHAR(255) NOT NULL,
    command VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    error_message TEXT,
    user_name VARCHAR(255),
    execution_duration DOUBLE DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

- **backups**

```sql
CREATE TABLE backups (
    id BINARY(16) PRIMARY KEY DEFAULT (UUID_TO_BIN(UUID())),
    file_name VARCHAR(255) NOT NULL UNIQUE,
    source_path TEXT NOT NULL,
    backup_path TEXT NOT NULL,
    og_file_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### PSQL Table Schemas

> **_NOTE:_**  If uuid extension not installed.

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

- **logs**

```sql
CREATE TABLE logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    action VARCHAR(255) NOT NULL,
    command VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    error_message TEXT,
    user_name VARCHAR(255),
    execution_duration DOUBLE PRECISION DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

- **backups**

```sql
CREATE TABLE backups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    file_name VARCHAR(255) NOT NULL UNIQUE,
    source_path TEXT NOT NULL,
    backup_path TEXT NOT NULL,
    og_file_name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```