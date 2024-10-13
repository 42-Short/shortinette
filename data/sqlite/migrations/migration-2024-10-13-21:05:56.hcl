table "users" {
  schema = schema.main
  column "id" {
    null = true
    type = text
  }
  primary_key {
    columns = [column.id]
  }
}
schema "main" {
}
