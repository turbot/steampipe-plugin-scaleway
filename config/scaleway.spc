connection "scaleway" {
  plugin  = "scaleway"

  # You may connect to one or more regions. If `regions` is not specified,
  # Steampipe will use a single default region using:
  # The `SCW_DEFAULT_REGION` environment variable
  # regions     = ["fr-par", "nl-ams"]

  # Set the static credential with the `access_key` and `secret_key` arguments
  # Alternatively, if no creds passed in config, you may set the environment variables using
  # `SCW_ACCESS_KEY` and `SCW_SECRET_KEY` arguments
  access_key = "YOUR_ACCESS_KEY"
  secret_key = "YOUR_SECRET_ACCESS_KEY"
}