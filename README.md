<h1 align="center">
    Dsync
</h1>

**Dsync** is a CLI tool designed to simplify the management and synchronization of data between diverse systems, including databases, cloud storage, and local file systems. It addresses real-world challenges in data migration, automation, and backup by offering secure synchronization, data transformation, scheduling, and task management features.

---

## Key Features

### 1. Multi-Source Data Synchronization
- **Database Sync**: Synchronize data between different databases (e.g., MySQL ↔ PostgreSQL).
- **Cloud Storage**: Export and import files between cloud platforms (e.g., AWS S3 ↔ Google Cloud Storage).
- **Local to Cloud**: Transfer data between local systems and cloud services effortlessly.

### 2. Data Transformation
- **Custom Transformations**: Apply transformations during synchronization, such as renaming columns or filtering rows.
- **User-Defined Scripts**: Support for scripts to modify data during transit.

### 3. Backup and Restore
- **Scheduled Backups**: Create backups for databases and files on a schedule.
- **One-Command Restore**: Restore backups with a single command.

### 4. Secure and Efficient Transfers
- **Encryption**: TLS encryption for secure database connections.
- **Authentication**: Use IAM roles and API keys for cloud authentication.
- **Incremental Updates**: Reduce transfer overhead by syncing only changes.

### 5. Task Scheduling
- **Automation**: Automate recurring tasks using a cron-like scheduler.
- **Management**: Manage schedules through CLI or configuration files.

### 6. Task Tracking and Logging
- **Detailed Logs**: Monitor sync status and troubleshoot with comprehensive logs.
- **Retry Mechanisms**: Track tasks in an internal database and retry failed tasks automatically.

---

## Dsync: Operations Overview

### Synchronization
- **Data Extraction**: Retrieve data from the source using appropriate connectors (e.g., MySQL driver, AWS SDK).
- **Optional Transformations**: Apply transformations such as filters or column mappings during the sync process.
- **Secure Data Transfer**: Transfer data to the destination with secure protocols.
- **Status Tracking**: Monitor task statuses, errors, and retries using the internal task database.

### Backup and Restore
- **Backups**:
  - Data is compressed and stored in user-specified formats (e.g., `.sql`, `.csv`).
  - Efficiently handles large datasets with compression to minimize storage requirements.
- **Restores**:
  - Reverses the backup process, including format conversions as needed.
  - Ensures data integrity and compatibility during restoration.

### Scheduling
- **Storage**:
  - Schedules are stored in an internal database or defined in a YAML configuration file.
- **Automated Execution**:
  - A built-in scheduler continuously monitors for due tasks and executes them on schedule.

### Task Management
- **Task Tracking**:
  - Tasks are recorded in an internal database with detailed statuses:
    - `pending`
    - `in-progress`
    - `completed`
    - `failed`
- **Detailed Logging**:
  - Each task generates logs that capture execution details, errors, and retry attempts.
- **Error Recovery**:
  - Failed tasks are retried automatically based on configurable retry mechanisms.

Dsync ensures robust operations with its comprehensive synchronization, backup, scheduling, and task management capabilities. Manage your workflows with confidence!
