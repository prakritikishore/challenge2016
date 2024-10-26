
# API Documentation

## Health Check
Check the health status of the service.

```bash
curl -X GET http://localhost:8080/health
```

## Distributor Endpoints

### Add Distributor
Add a new distributor.

```bash
curl -X POST http://localhost:8080/distributor-service/distributor/add \
    -H "Content-Type: application/json" \
    -d '{"name": "Distributor1"}'
```

### Get Distributor
Retrieve details of a specific distributor.

```bash
curl -X GET http://localhost:8080/distributor-service/distributor/Distributor1
```

### Delete Distributor
Delete an existing distributor.

```bash
curl -X DELETE http://localhost:8080/distributor-service/distributor/Distributor1
```

## Sub-Distributor Endpoints

### Add Sub-Distributor
Add a new sub-distributor to a parent distributor.

```bash
curl -X POST http://localhost:8080/distributor-service/sub-distributor/add \
    -H "Content-Type: application/json" \
    -d '{"distributor_name": "Distributor1", "sub_distributor_name": "Distributor2"}'
```

### Delete Sub-Distributor
Delete an existing sub-distributor.

```bash
curl -X DELETE http://localhost:8080/distributor-service/sub-distributor/Distributor1/Distributor2
```

## Permission Endpoints

### Add Permission
Add permissions to a distributor.

```bash
curl -X POST http://localhost:8080/distributor-service/permissions/Distributor1 \
    -H "Content-Type: application/json" \
    -d '{"region_code": "IN", "is_included": true}'
```

### Delete Permission
Delete permissions from a distributor.

```bash
curl -X DELETE http://localhost:8080/distributor-service/permissions/Distributor1 \
    -H "Content-Type: application/json" \
    -d '{"region_code": "IN"}'
```

### Check Permission
Check if a distributor has a specific permission.

```bash
curl -X POST http://localhost:8080/distributor-service/check-permissions \
    -H "Content-Type: application/json" \
    -d '{"name": "Distributor1", "region_code": "IN"}'
```

## Authorization Endpoint

### Authorize SubDistributor
Authorize a sub-distributor.

```bash
curl -X POST http://localhost:8080/distributor-service/authorize/sub-distributor \
    -H "Content-Type: application/json" \
    -d '{"from_distributor": "Distributor1", "to_distributor": "Distributor2", "region_code": "US"}'
```

