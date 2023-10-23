connection "scaleway" {
  plugin  = "scaleway"

  # Set the static credential with the `access_key` and `secret_key` arguments.
  # Alternatively, if no creds passed in config, you may set the environment
  # variables using the `SCW_ACCESS_KEY` and `SCW_SECRET_KEY` arguments.
  access_key = "YOUR_ACCESS_KEY"
  secret_key = "YOUR_SECRET_ACCESS_KEY"
  
  # Your organization ID is the identifier of your account inside Scaleway infrastructure.
  # This is only required while querying the `scaleway_iam_api_key` and `scaleway_iam_user` tables. 
  # organization_id = "YOUR_ORGANIZATION_ID"

  # You may connect to one or more regions. If `regions` is not specified,
  # Steampipe will use a single default region using the `SCW_DEFAULT_REGION`
  # environment variable.
  # regions = ["fr-par", "nl-ams"]
}
