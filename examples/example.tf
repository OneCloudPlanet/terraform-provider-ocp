terraform {
  required_providers {
    ocp = {
      source = "OnePointCollab/ocp"
      version = "0.1.3"
    }
  }
}

provider "ocp" {
  api_token = "49e9a4cff4ce46374f408bc88e410b745c715bfb"
  region = "ua"
}

data "ocp_flavor" "master_flavors" {
  filter {
    memory_gb = 4
  }
}

data "ocp_flavor" "node_flavors" {
  filter {
    memory_gb = 8
  }
}

data "ocp_cluster_version" "list_versions" {
  filter {
    image_name = "Ubuntu 22.04"
  }
}

data "ocp_cluster_networking" "list_networking" {
  filter {
    network_name = "Calico"
    version      = "v3.26.1"
  }
}

output "ms_flavors" {
  value = data.ocp_flavor.master_flavors.flavors[0].id
}
output "nd_flavors" {
  value = data.ocp_flavor.node_flavors.flavors[0].id
}
output "networking" {
  value = data.ocp_cluster_networking.list_networking.networking[0].id
}
output "versions" {
  value = data.ocp_cluster_version.list_versions.versions[0].images[0].image_name
}

resource "ocp_cluster" "new_cluster" {
  cluster_name     = "cluster-name"
  cluster_version  = data.ocp_cluster_version.list_versions.versions[0].version
  master_flavor_id = data.ocp_flavor.master_flavors.flavors[0].id
  master_count     = 3
  image            = data.ocp_cluster_version.list_versions.versions[0].images[0].image_name
  networking       = data.ocp_cluster_networking.list_networking.networking[0].id
  restriction_api  = true
  restriction_ips  = ["12.12.12.12/32", "13.13.13.13/32"]
  node_pool {
    name       = "nodepool-name"
    flavor_id  = data.ocp_flavor.node_flavors.flavors[0].id
    node_count = 3
    autoscale  = true
    max_count  = 5
  }
  addons {
    dashboard = true
    metrics   = true
    nginx     = true
  }
}

resource "ocp_nodepool" "new_nodepool" {
  name = "nodepool-name"
  flavor_id = data.ocp_flavor.node_flavors.flavors[0].id
  node_count = 2
  autoscale = false
  cluster = ocp_cluster.new_cluster.id
}
