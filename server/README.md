# Uploader
Accepts json of parameters and file to upload to the db

### example json(.xlsx):
```json
{
  "db_table": "db table name with scheme e.g. public.users",
  "field_row": (start field row in file e.g. 3)(int64),
  "data_row": (start data row in file e.g. 4)(int64),
  "fields": [
    {
      "db": "db_field_name",
      "field": "file_field_name"
    },
    ...
}
```

### example json(.xml):

```json
{
  "db_table": "db table name with scheme e.g. public.users",
  "root": "starting tag",
  "tags": [
    {
      "parent": "parent node",
      "db": "db_field_name",
      "tag": "file_field_name"
    },
    ...
}
```

### File extensions
- (.xlsx .xml)