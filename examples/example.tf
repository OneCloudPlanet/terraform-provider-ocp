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

data "ocp_cluster_addons" "list_addons" {
  filter {
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
output "addons" {
  value = data.ocp_cluster_addons.list_addons.addons[1]
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
    name = data.ocp_cluster_addons.list_addons.addons[0].name
    version = data.ocp_cluster_addons.list_addons.addons[0].releases[0].version
  }
  addons {
    name = data.ocp_cluster_addons.list_addons.addons[1].name
    version = data.ocp_cluster_addons.list_addons.addons[1].releases[0].version
  }
}

resource "ocp_nodepool" "new_nodepool" {
  name = "second-nodepool"
  flavor_id = data.ocp_flavor.node_flavors.flavors[0].id
  node_count = 2
  autoscale = false
  cluster = ocp_cluster.new_cluster.id
}
