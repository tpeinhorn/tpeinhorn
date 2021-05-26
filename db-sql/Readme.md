This is our DB service which is an instance of MSSQL.


The mounted directories must have "chmod -R 0777 /mydirectory" before

```sql
USE master;
GO

CREATE DATABASE SampleDB;
GO

CREATE TABLE dbo.MyTable (
  id bigint IDENTITY(1,1) PRIMARY KEY,
  name varchar(500) null
)
GO

``` 