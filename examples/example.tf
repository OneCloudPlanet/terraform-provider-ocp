provider "ocp" {
  api_token = "******435534" # your open_api token
  region = "dev"             # region - "ua" / "pl"
}

# List flavor's
data "ocp_flavor" "flavors" {
  # filtered by "vcpus", "memory_mb", "memory_gb", "root_gb"
  filter {
    vcpus = 1
    root_gb = 10
  }
}

# List cluster versions
data "ocp_cluster_version" "versions" {
  # filtered by "version", "image_name", "os_distro"
  filter {
    image_name = "Ubuntu 22.04"
  }
}

# List cluster networking
data "ocp_cluster_networking" "networking" {
  # filtered by "network_name", "version"
  filter {
    network_name = "Calico"
  }
}

# Create Kubernetes Cluster
resource "ocp_cluster" "cluster" {
  cluster_name = "cluster-name"
  cluster_version = "v1.28.6" # use "version" in datasource "ocp_cluster_version" or text "v*.**.**"
  master_flavor_id = "a9ec3da2-87b2-4fbe-b3ae-5374a204653d" # id in datasource object "ocp_flavor"
  master_count = 1
  image = "ubuntu-22.04-kube-v1.27.4" # images[].image_name in datasource "ocp_cluster_version"
  networking = "137a00a5-c1cf-44cf-9b5b-2ac5719f41e5" # id in datasource object "ocp_cluster_networking"
  restriction_api = true
  restriction_ips = ["89.64.63.31/32", "2.59.223.3/32"]

  node_pool {
    name       = "nodepool-name"
    flavor_id  = "61c79c33-6113-41fd-9d87-594764813c67" # id in datasource object "ocp_flavor"
    node_count = 3
    autoscale  = false
    max_count  = 0
  }

  addons {
    dashboard = true
    metrics   = true
  }
}

# Create NodePool object
resource "ocp_nodepool" "nodepool" {
  name = "nodepool-name"
  flavor_id = "61c79c33-6113-41fd-9d87-594764813c67" # id in datasource object "ocp_flavor"
  node_count = 2
  autoscale = false
  max_count = 0
  cluster = ocp_cluster.cluster.id
}