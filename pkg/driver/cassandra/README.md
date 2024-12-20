# Cassandra Tables and Schemas

## Tables

### 1. aggregate_events
Stores all events for aggregates.

Schema:
```sql
CREATE TABLE aggregate_events (
    aggregateId UUID,
    eventId UUID,
    eventType TEXT,
    data BLOB,
    timestamp TIMESTAMP,
    PRIMARY KEY (aggregateId, timestamp)
) WITH CLUSTERING ORDER BY (timestamp ASC);
```

### 2. aggregate_latest_event
Stores the latest event for each aggregate.

Schema:
```sql
CREATE TABLE aggregate_latest_event (
    aggregate_id UUID,
    latest_event_id UUID,
    timestamp TIMESTAMP,
    PRIMARY KEY (aggregate_id)
);
```

### 3. aggregate_event_count
Stores the count of events for each aggregate.

Schema:
```sql
CREATE TABLE aggregate_event_count (
    aggregate_id UUID,
    event_count INT,
    PRIMARY KEY (aggregate_id)
);
```

### 4. aggregate_events_version
Stores the version of each event for auditing purposes.

Schema:
```sql
CREATE TABLE aggregate_events_version (
    aggregate_id UUID,
    event_id UUID,
    event_version INT,
    timestamp TIMESTAMP,
    PRIMARY KEY (aggregate_id, event_version)
) WITH CLUSTERING ORDER BY (event_version ASC);
```

### 5. snapshots
Stores snapshots of aggregates.

Schema:
```sql
CREATE TABLE snapshots (
    aggregate_id UUID,
    data BLOB,
    version INT,
    timestamp TIMESTAMP,
    PRIMARY KEY (aggregate_id, timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC);
```

## Optimization

- **Write Operations**: Batch operations are used to minimize the number of queries and improve performance.
- **Read Operations**: Queries are optimized to fetch data using primary keys, ensuring efficient lookups.
