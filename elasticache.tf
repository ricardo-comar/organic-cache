resource "aws_elasticache_cluster" "quotation-cache" {
  cluster_id           = "quotation-cache-cluster"
  engine               = "redis"
  node_type            = "cache.t2.medium"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis3.2"
}