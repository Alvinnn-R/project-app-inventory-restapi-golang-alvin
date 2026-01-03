# Unit Testing Documentation

## Overview

Proyek ini telah dilengkapi dengan unit tests untuk meningkatkan kualitas kode dan memudahkan maintenance. Tests mencakup repository layer dan service layer dengan menggunakan mocking untuk database dan dependencies.

## Testing Libraries

1. **pgxmock v3** - Mock untuk PostgreSQL database (pgx driver)
2. **testify** - Assertion library dan mocking framework
3. **Go testing** - Built-in testing package

## Installation Dependencies

Sebelum menjalankan tests, pastikan dependencies sudah terinstall:

```bash
go get github.com/pashagolub/pgxmock/v3
go get github.com/stretchr/testify
go get github.com/stretchr/objx
```

Atau jalankan:

```bash
go mod tidy
```

## Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Tests with Coverage

```bash
go test -cover ./...
```

### Run Tests with Detailed Coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Specific Package Tests

```bash
# Repository tests
go test ./repository/...

# Service tests
go test ./service/...
```

### Run Specific Test File

```bash
go test ./repository/item_test.go
go test ./service/item_test.go
```

### Run with Verbose Output

```bash
go test -v ./...
```

## Test Files Overview

### 1. Repository Layer Tests

#### repository/item_test.go

Tests untuk ItemRepository dengan coverage:

- ✅ `TestItemRepository_Create` - Test insert item baru
  - Success scenario
  - Database error scenario
- ✅ `TestItemRepository_FindByID` - Test find by ID
  - Found scenario
  - Not found scenario
- ✅ `TestItemRepository_FindLowStock` - Test low stock detection
  - Items found scenario
  - Empty result scenario
- ✅ `TestItemRepository_Update` - Test update item
  - Success scenario
  - Not found scenario
- ✅ `TestItemRepository_Delete` - Test delete item
  - Success scenario
  - Not found scenario

**Total: 5 test functions, 10 test cases**

#### repository/category_test.go

Tests untuk CategoryRepository dengan coverage:

- ✅ `TestCategoryRepository_Create` - Test insert category
  - Success scenario
  - Database error scenario
- ✅ `TestCategoryRepository_FindByID` - Test find by ID
  - Found scenario
  - Not found scenario
- ✅ `TestCategoryRepository_FindByName` - Test find by name
  - Found scenario
  - Not found scenario
- ✅ `TestCategoryRepository_FindAll` - Test pagination
  - With results scenario
  - Empty result scenario
- ✅ `TestCategoryRepository_Update` - Test update category
  - Success scenario
  - Not found scenario
- ✅ `TestCategoryRepository_Delete` - Test delete category
  - Success scenario
  - Not found scenario

**Total: 6 test functions, 12 test cases**

#### repository/user_test.go

Tests untuk UserRepository dengan coverage:

- ✅ `TestUserRepository_Create` - Test insert user
  - Success scenario
  - Duplicate email scenario
- ✅ `TestUserRepository_FindByID` - Test find by ID
  - Found scenario
  - Not found scenario
- ✅ `TestUserRepository_FindByUsername` - Test find by username
  - Found scenario
  - Not found scenario
- ✅ `TestUserRepository_FindByEmail` - Test find by email
  - Found scenario
  - Not found scenario
- ✅ `TestUserRepository_FindAll` - Test pagination
  - With results scenario
  - Empty result scenario
- ✅ `TestUserRepository_Update` - Test update user
  - Success scenario
  - Not found scenario
- ✅ `TestUserRepository_Delete` - Test delete user
  - Success scenario
  - Not found scenario
- ✅ `TestUserRepository_CountActive` - Test count users
  - Success scenario
  - Zero users scenario

**Total: 8 test functions, 16 test cases**

### 2. Service Layer Tests

#### service/item_test.go

Tests untuk ItemService dengan mock repository:

- ✅ `TestItemService_Create` - Test create item business logic
  - Success scenario
  - Duplicate name scenario
- ✅ `TestItemService_GetByID` - Test get by ID
  - Found scenario
  - Not found scenario
- ✅ `TestItemService_GetAll` - Test get all items
  - With results scenario
  - Empty result scenario
- ✅ `TestItemService_GetLowStockItems` - Test low stock items
  - With low stock items scenario
  - No low stock items scenario
- ✅ `TestItemService_Update` - Test update item
  - Success scenario
  - Not found scenario
- ✅ `TestItemService_Delete` - Test delete item
  - Success scenario
  - Not found scenario

**Total: 6 test functions, 12 test cases**

## Test Coverage Summary

### Current Coverage:

- **Repository Layer**: 3 files (item, category, user) - ~50% coverage
- **Service Layer**: 1 file (item) - ~40% coverage
- **Total Test Functions**: 25 functions
- **Total Test Cases**: 50+ test cases

### Covered Functionality:

✅ CRUD operations (Create, Read, Update, Delete)
✅ Find by ID, Name, Email, Username
✅ Pagination (FindAll)
✅ Low stock detection
✅ Error handling (Not found, Database errors)
✅ Business logic validation (Duplicate checks)
✅ Active user counting

## Testing Patterns

### 1. Repository Tests Pattern (using pgxmock)

```go
func TestRepository_Method(t *testing.T) {
    // Setup mock database
    mock, err := pgxmock.NewPool()
    if err != nil {
        t.Fatal(err)
    }
    defer mock.Close()

    repo := &Repository{DB: mock}

    t.Run("scenario_name", func(t *testing.T) {
        // Mock expectations
        mock.ExpectQuery(`SQL QUERY`).
            WithArgs(args...).
            WillReturnRows(rows)

        // Execute
        result, err := repo.Method(ctx, args...)

        // Assert
        assert.NoError(t, err)
        assert.NotNil(t, result)

        // Verify expectations
        err = mock.ExpectationsWereMet()
        assert.NoError(t, err)
    })
}
```

### 2. Service Tests Pattern (using testify mock)

```go
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Method(ctx context.Context, args...) (result, error) {
    args := m.Called(ctx, args...)
    return args.Get(0).(Type), args.Error(1)
}

func TestService_Method(t *testing.T) {
    mockRepo := new(MockRepository)
    service := &Service{Repo: mockRepo}

    t.Run("scenario_name", func(t *testing.T) {
        // Setup mock expectations
        mockRepo.On("Method", mock.Anything, args...).
            Return(result, nil).Once()

        // Execute
        result, err := service.Method(ctx, args...)

        // Assert
        assert.NoError(t, err)
        assert.NotNil(t, result)

        // Verify
        mockRepo.AssertExpectations(t)
    })
}
```

## Benefits of Unit Testing

1. **Early Bug Detection**: Menemukan bugs sebelum production
2. **Refactoring Confidence**: Aman untuk refactor code dengan test coverage
3. **Documentation**: Tests berfungsi sebagai dokumentasi cara menggunakan code
4. **Code Quality**: Mendorong penulisan code yang lebih modular dan testable
5. **Regression Prevention**: Mencegah bugs lama muncul kembali

## Best Practices Applied

1. ✅ **Table-Driven Tests**: Menggunakan subtests dengan `t.Run()`
2. ✅ **Mock Isolation**: Setiap test isolated dengan mock baru
3. ✅ **Clear Naming**: Test names yang descriptive (Method_Scenario)
4. ✅ **Assertions**: Menggunakan testify untuk assertions yang clear
5. ✅ **Error Cases**: Test both success dan error scenarios
6. ✅ **Mock Verification**: Verify mock expectations di akhir test

## Future Improvements

Untuk mencapai 80%+ coverage, perlu ditambahkan:

1. **Repository Tests**:

   - [ ] rack_test.go
   - [ ] warehouse_test.go
   - [ ] sale_test.go (with transaction mocking)
   - [ ] submission_test.go
   - [ ] report_test.go

2. **Service Tests**:

   - [ ] category_test.go
   - [ ] user_test.go
   - [ ] sale_test.go (business logic with stock adjustment)
   - [ ] auth_test.go (login/logout logic)

3. **Handler Tests**:

   - [ ] HTTP handler tests with mock service layer
   - [ ] Request validation tests
   - [ ] Response format tests

4. **Integration Tests**:
   - [ ] End-to-end API tests
   - [ ] Database integration tests with test containers

## Troubleshooting

### Common Issues:

1. **Missing Dependencies**

```bash
go: missing go.sum entry for module providing package
```

**Solution**: Run `go mod tidy`

2. **Mock Not Working**

```bash
mock: Unexpected Method Call
```

**Solution**: Check mock expectations and ensure `.Once()` or `.Times(n)` is set

3. **Tests Failing After Changes**

- Update mock expectations to match new behavior
- Check SQL queries in ExpectQuery match actual queries
- Verify WithArgs parameters match

## Conclusion

Unit testing sudah diimplementasikan dengan baik menggunakan:

- ✅ pgxmock untuk database mocking
- ✅ testify untuk assertions dan service mocking
- ✅ Comprehensive test coverage untuk critical paths
- ✅ Clear test structure dengan subtests
- ✅ Both success dan error scenarios covered

**Current Status**: ~50% test coverage dengan fokus pada repository dan service layers.

**Target**: 80%+ coverage dengan menambahkan tests untuk remaining repositories, services, dan handlers.
