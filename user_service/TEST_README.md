# Integration Test Setup

## Database Configuration

Integration tests use a **real PostgreSQL database** but with safety measures to prevent data loss.

## Setup Test Database

### Option 1: Use Environment Variable (Recommended)

```bash
export TEST_DATABASE_URL="host=localhost user=postgres password=123postgres dbname=user_service_test port=5432 sslmode=disable"
```

### Option 2: Create Test Database

```sql
-- Connect to PostgreSQL as superuser
psql -U postgres

-- Create test database
CREATE DATABASE user_service_test;

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE user_service_test TO postgres;
```

## Safety Features

✅ **No Production Data Loss**: Tests use separate test database  
✅ **Unique Test Data**: Each test generates unique data using timestamps  
✅ **Automatic Cleanup**: Data is cleaned after each test  
✅ **Foreign Key Safety**: Proper cleanup order prevents constraint violations

## Running Tests

### Run Only Integration Tests

```bash
go test -run TestUserServiceIntegrationTestSuite ./service -v
```

### Run All Tests

```bash
go test ./service -v
```

## Test Flow

1. **SetupSuite**: Connect to database, migrate tables
2. **SetupTest**: Clean previous test data
3. **Test Execution**: Run test with unique data
4. **Cleanup**: Remove test data automatically

## Important Notes

- **NEVER** use production database for tests
- Tests automatically create/cleanup tables
- Each test uses unique identifiers to avoid conflicts
- Database connection is shared across test suite for performance
