# This terraform script is designed to self-ignite the light-weight
# kubernetes install from Rancher called k3s (5 lighter than k8s).
#
# k3s requires some reasonably sane kernel, ubutu18.04 or centos:7.6 and
# above will be adequate.
#
# Terraform host creation will result in the latest k3s image to be
# downloaded and executed on first boot of the machine. 
#
# k3s uses a single kubernetes master - the cluster will be created
# with a unique token (secret) that is shared among the master and
# all workers. Each Terraform invocation of his script in a new
# directory will result in a new k3s-cluster deployment with a 
# unique token. 
#
# Cluster creation is done via the single comamnd
#
#    terraform apply
#
# Cluster customization is possible by overriding default variable
# values used in the script
#
#  location: specifies the geographic location of the cluster.
#  ssk_keys: specifies a list of ssh-key-names to insert into the master and workers.
#  msize: specifies the machine=size to provsion.
#  workers: the number of initial worker machines.
#
#  terraform apply -var='location=USA:Texas:AUSL2' -var='workers=4'
# 
#
# To later change the size of the cluster
#   terraform apply -var='workers=6'


provider "quake" {

}

# Defines the machine-size to use when provisioning the cluster. This can
# be overridden on the command line like
#
#     terraform apply -var='msize=Large'
variable "msize" {
  default = "Any"
}

# Defines the location where the custer will be deployed. This can
# be overridden on the command line like
#
#     terraform apply -var='location=USA:Texas:AUSL2'
variable "location" {
  default = "USA:West Central:FTC DEV 4"
}


# Defines the ssh-key names the custer will be deployed with. This can
# be overridden on the command line like
#
#     terraform apply -var='ssh_keys=["my-key", "another-key"]'
variable "ssh_keys" {
  default = ["MarkW"]#["User1 - Linux"] #["MarkW"]
}

output "worker_ips" {
  # Output a map of hostame with all the IP addresses assigned on each network.
  value = zipmap(quake_host.workers.*.name, quake_host.workers.*.connections)
}

output "master_ip" {
  # Output a map of hostame with all the IP addresses assigned on each network.
  value = map(quake_host.master.name, quake_host.master.connections)
}


# Defines the intial number of worker nodes in the cluster.
#
#     terraform apply -var='workers= 4'
variable "workers" {
  default = 2
}

resource "random_password" "k3s_token" {
  # This will be the random k3s cluster secret (token). It will be 
  # marked as sensitive by terraform ad should not appeear in regular outputs.
  # Don't use any special characters that may cause shell-escape issues.
  length  = 48
  special = false
}

// Determine what physical resources are available to us.
data "quake_available_resources" "physical" {

}

# Attempt to locate a decent image on the portal we are autenticated against.
# The script will use the fist image that matches and abort if no images match.
data "quake_available_images" "ubuntu" {
  # Rancher k3s requires only a reasonably upto date kernel.
  # select anything that looks like ubuntu:18.04
  filter {

    name   = "flavor"
    values = ["(?i)ubuntu"] // case insensitive for Ubuntu or ubuntu etc.
    #values = ["(?i)rancher"] 
  }
  filter {
    name   = "version"
    values = ["18.04*"] // all 18.04 image variants
  }
}

#
# A note on the user_data field for host creation. This value can be any valid cloud-config
# information. It must be cloud-init formatted yaml.
#
# In the example below we are using yaml and the runcmd tag. This takes a list
# of commands and executes them pretty much as though they were run once by
# the system at the same system state and runlevel as /etc/rc.local.
#
# We need to perform a wide-area curl operations and for this to succeed we require
#   - configured network (IP address, routing information, etc)
#   - name resolution services active
#   - (proxy configuration if Quake racks are behind a firewall)
#
# Although the system will most likely have this information defined such
# information is obtained via environment variables typically written to files
# in /etc. These are picked up by shells when they start. In the case of cloud-init, however,
# these files are being written by the init process and therefore the process itself
# doesn't have the luxury of having these 'standard' environment variables defined.
# To work around this, the contents of /etc/environment are evaluated and exported before
# network access attempted.
resource "quake_host" "master" {
  name          = "master"
  image_flavor  = data.quake_available_images.ubuntu.images[0].flavor
  image_version = data.quake_available_images.ubuntu.images[0].version
  machine_size  = var.msize
  ssh           = var.ssh_keys
  networks      = distinct(concat(["Private", "Public"], [for net in data.quake_available_resources.physical.networks : net.name if net.host_use == "Required" && net.location == var.location]))
  location      = var.location
  description   = "Master k3s"
  ## Note there are cases where a double $$ is used in this cloud-config text. By defult Terrafrom will interpret a $ as the start
  # of variable substitution using terraform values. Since we generating a shell script $ is also used to reference shell 
  # or environment variables. We use the $$ as way to escape the usual terraform variable substitution process.
  user_data     = <<EOF
#cloud-config
# We need mutiple write files with existing Quake write-files
# see  https://cloudinit.readthedocs.io/en/latest/topics/merging.html
merge_how:
 - name: list
   settings: [append]
 - name: dict
   settings: [no_replace, recurse_list]

write_files:
  - owner: root
    path: /root/dashboard.admin-user.yml
    permissions: '0400'
    content: |
      apiVersion: v1
      kind: ServiceAccount
      metadata:
        name: admin-user
        namespace: kubernetes-dashboard

  - owner: root
    path: /root/dashboard.admin-user-role.yml
    permissions: '0400'
    content: |
      apiVersion: rbac.authorization.k8s.io/v1
      kind: ClusterRoleBinding
      metadata:
        name: admin-user
      roleRef:
        apiGroup: rbac.authorization.k8s.io
        kind: ClusterRole
        name: cluster-admin
      subjects:
      - kind: ServiceAccount
        name: admin-user
        namespace: kubernetes-dashboard

  - owner: root
    path: /root/k3s.sh
    permissions: '0700'
    content: |
      #!/bin/bash
      # Install k3s as a master using the random token
      curl -sfL https://get.k3s.io | K3S_TOKEN=${random_password.k3s_token.result} sh -
      # Install the k3s Dashboard
      GITHUB_URL=https://github.com/kubernetes/dashboard/releases
      VERSION_KUBE_DASHBOARD=$(curl -w '%%{url_effective}' -I -L -s -S $${GITHUB_URL}/latest -o /dev/null | sed -e 's|.*/||')
      k3s kubectl create -f https://raw.githubusercontent.com/kubernetes/dashboard/$${VERSION_KUBE_DASHBOARD}/aio/deploy/recommended.yaml
      k3s kubectl create -f /root/dashboard.admin-user.yml -f /root/dashboard.admin-user-role.yml
      # We can't actually forward a port from a pod that isn't yet running..
      sleep 60  # TODO - fix this with a poll for pod status on the dashboard.
      # Expose the dashboard on port 8443 on all interfaces.
      nohup kubectl --namespace kubernetes-dashboard port-forward service/kubernetes-dashboard  --address 0.0.0.0 8443:443 &
    

runcmd:
  - 'export `cat /etc/environment`; /root/k3s.sh'

EOF
}

# Worker nodes need access to the master over the private network; they also need acess to the Public network 
# to be bable to download the k3s install from http://get.k3s.io 
#
# A bespoke image could be created that holds the default bits of the install script but it's hard to guarantee
# that the script won't make calls to the outside world as it runs.
resource "quake_host" "workers" {
  count         = var.workers
  name          = "worker-${count.index}"
  image_flavor  = data.quake_available_images.ubuntu.images[0].flavor
  image_version = data.quake_available_images.ubuntu.images[0].version
  machine_size  = var.msize
  ssh           = var.ssh_keys
  networks      = distinct(concat(["Private", "Public"], [for net in data.quake_available_resources.physical.networks : net.name if net.host_use == "Required"  && net.location == var.location]))
  location      = var.location
  description   = "Worker k3s"
  user_data     = <<EOF
merge_how:
 - name: list
   settings: [append]
 - name: dict
   settings: [no_replace, recurse_list]

runcmd:
  - 'export `cat /etc/environment`; curl -sfL https://get.k3s.io | K3S_TOKEN=${random_password.k3s_token.result} K3S_URL=https://${quake_host.master.connections.Private}:6443 sh -'
EOF
}




