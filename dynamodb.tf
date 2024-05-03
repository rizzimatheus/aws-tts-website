module "dynamodb_table" {
  source  = "terraform-aws-modules/dynamodb-table/aws"
  version = "~> 4.0"

  name                        = "tts-website-table"
  hash_key                    = "id"
  table_class                 = "STANDARD"
  deletion_protection_enabled = false

  # List of nested attribute definitions.
  # Only required for hash_key (partition key) and range_key (sort key) attributes.
  # type: S, N, or B for (S)tring, (N)umber or (B)inary data
  attributes = [
    {
      name = "id"
      type = "S"
    }
  ]
}
