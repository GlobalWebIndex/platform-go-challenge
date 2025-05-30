# Infrastructure

This directory contains infrastructure-related files and instructions for setting up the database and running it in Docker.

## Overview

The infrastructure consists of:
- **Docker Compose**: Containerized setup for local development
- **Migration Scripts**: Automated database schema setup and data seeding

## Database Schema

The application uses a PostgreSQL database with the following main entities:

```mermaid
erDiagram
    users {
        int id PK
        varchar username UK
        varchar hashed_password
    }
    
    assets {
        int id PK
        varchar type "chart|insight|audience"
    }
    
    charts {
        int asset_id PK,FK
        varchar title
        varchar title_x_axis
        varchar title_y_axis
        jsonb data
    }
    
    insights {
        int asset_id PK,FK
        text content
    }
    
    audiences {
        int asset_id PK,FK
        varchar gender
        varchar birth_country
        varchar age_group
        int hours_spent_daily_on_social_media
        int last_month_purchases
    }
    
    favorites {
        int id PK
        int user_id FK
        int asset_id FK
        text description
    }
    
    users ||--o{ favorites : "has"
    assets ||--o{ favorites : "favorited_as"
    assets ||--o| charts : "specializes_to"
    assets ||--o| insights : "specializes_to"
    assets ||--o| audiences : "specializes_to"
```

## Postgres

The Postgres directory contains SQL migration scripts for setting up the database schema and seeding data:

- `01.schema.up.sql`: Creates database schema including all tables, constraints, and indexes
- `02.populate.up.sql`: Seeds initial data with sample users, assets, charts, insights, audiences, and favorites

The schema and populate files are automatically executed when the PostgreSQL container starts for the first time.

### Sample Data Overview

##### Users and Their Favorites

| User ID | Username | Asset ID | Asset Type | Favorite ID | Description |
|---------|----------|----------|------------|-------------|-------------|
| 1       | alice    | 1        | chart      | 1           | Primary sales chart |
| 1       | alice    | 4        | insight    | 2           | Key Q4 insight |
| 2       | bob      | 3        | chart      | 3           | Engagement tracking |
| 2       | bob      | 2        | audience   | 4           | US |
| 3       | charlie  | 6        | insight    | 5           | Friday peak analysis |
| 4       | dave     | -        | -          | -           | No favorites |

##### Asset Details Summary

| Asset ID | Type     | Title/Content Summary |
|----------|----------|-----------------------|
| 1        | chart    | Quarterly Sales |
| 2        | audience | Female, USA, 18-24, 3hrs social media, 5 purchases |
| 3        | chart    | User Engagement |
| 4        | insight  | Sales increased by 20% in Q4... |
| 5        | audience | Male, Canada, 25-34, 2hrs social media, 8 purchases |
| 6        | insight  | User engagement peaks on Fridays... |
| 7        | chart    | Sales by Region |
| 8        | insight  | Monthly trends show steady increase... |
| 9        | audience | Female, UK, 35-44, 1hr social media, 10 purchases |
| 10       | insight  | Product A most popular among 18-24... |
| 11       | chart    | Product Popularity |
| 12       | audience | Male, Germany, 45-54, 4hrs social media, 2 purchases |

## Docker

> **NOTE**: You will need to have Docker installed on your machine to run the Docker Compose setup. If not, check the official [Docker documentation](https://docs.docker.com/desktop/) for installation instructions.

The Docker directory contains a Docker Compose file for setting up PostgreSQL and is **mainly used for local development**.

### Docker Compose Instructions

This section provides instructions for managing the Docker Compose setup with PostgreSQL service.

#### Starting the PostgreSQL Service

To spin up the database, navigate to the `infra/docker` directory and run:

```bash
docker-compose up -d
```

When first starting the containers, the PostgreSQL database will automatically apply the migration scripts to set up the necessary schema and seed data.

#### Connecting to PostgreSQL Service

**Connection Details:**
- **Host**: `localhost`
- **Port**: `5432`
- **User**: `gwi`
- **Password**: `gwi`
- **Database**: `gwi`

**Example connections:**

Using psql command line:
```bash
psql -h localhost -p 5432 -U gwi -d gwi
```

Using connection string:
```
postgresql://gwi:gwi@localhost:5432/gwi
```

#### Health Checks

The PostgreSQL service includes health checks that verify the database is ready to accept connections. You can check the service status:

```bash
docker-compose ps
```

#### Stopping PostgreSQL Service

To stop the service without removing the container:

```bash
docker-compose stop
```

To stop and remove containers, but preserve volumes:

```bash
docker-compose down
```

### Volume Management

The database data is persisted in a Docker volume named `pgdata`.

#### Listing Volumes

```bash
docker volume ls
```

#### Stopping Services and Removing Volumes

To stop services and remove volumes (**warning: all data will be deleted**):

```bash
docker-compose down -v
```

#### Manually Removing Specific Volumes

```bash
# Replace "docker_pgdata" with the actual volume name (check with docker volume ls)
docker volume rm docker_pgdata
```

#### Removing All Unused Volumes

```bash
docker volume prune
```

### Development Workflows

#### Restarting with Clean Data

If you need to start fresh with clean databases:

```bash
docker-compose down -v && docker-compose up -d
```

#### Applying Schema Changes

To apply new migration files:

1. Stop the service: `docker-compose down -v`
2. Update your migration files in `../postgres/migrations/`
3. Restart: `docker-compose up -d`

#### Backup and Restore

**Create a backup:**
```bash
docker exec postgres_gwi pg_dump -U gwi gwi > backup.sql
```

**Restore from backup:**
```bash
docker exec -i postgres_gwi psql -U gwi -d gwi < backup.sql
```

### Troubleshooting

If you encounter issues with the services:

1. **Check service health:**
   ```bash
   docker-compose ps
   ```

2. **Inspect service logs:**
   ```bash
   docker-compose logs postgres
   ```

3. **Check if port 5432 is available:**
   ```bash
   netstat -an | grep 5432
   ```

4. **Restart services:**
   ```bash
   docker-compose restart
   ```

5. **Rebuild and restart:**
   ```bash
   docker-compose down
   docker-compose up -d --build
   ```

6. **Connect to container directly:**
   ```bash
   docker exec -it postgres_gwi psql -U gwi -d gwi
   ```

### Common Issues

- **Port conflicts**: If port 5432 is already in use, modify the port mapping in `docker-compose.yml`
- **Permission issues**: Ensure Docker has proper permissions to create volumes