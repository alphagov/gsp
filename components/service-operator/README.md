# Service Operator

## PostgreSQL via AWS RDS

Given a configuration file of the form:

```yaml
  kind: Postgres
spec:
  aws:
    instanceType: db.m5.large
```

this service operator will create a Postgres database in AWS RDS, and
output a connection string.
