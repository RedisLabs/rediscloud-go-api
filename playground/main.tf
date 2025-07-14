resource "rediscloud_essentials_database" "db2" {
  lifecycle {
    ignore_changes = [
      # password,
      data_persistence
    ]
  }
  data_eviction                         = "noeviction"
  data_persistence                      = "aof-every-1-second"
  enable_database_clustering            = false
  enable_default_user                   = false
  password                              = "example"
  enable_tls                            = false
  external_endpoint_for_oss_cluster_api = false
  memory_limit_in_gb                    = 0
  modules                               = []
  name                                  = "DB2"
  protocol                              = "redis"
  regex_rules                           = []
  replication                           = true
  resp_version                          = "resp3"
  subscription_id                       = 1515407
  support_oss_cluster_api               = false
  tags                                  = {}
  alert {
    name  = "datasets-size"
    value = 80
  }
