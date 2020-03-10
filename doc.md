# Help to use the civo terraform provider

## Legend

- ``#`` (optional parameter)

### Define provider

```
provider "civo" {
  token = "token"
}
```

### Create a Network

```bash
resource "civo_network" "custom_net" {
    label = "test_network"
}
```

### Create a Instances

```bash
resource "civo_instance" "my-test-instance" {
    hostname = "test-terraform"
    initial_user = "root"
    size = "g2.large"
    tags = ["hello", "test"]
    # template = "id"
    # network_id = civo_network.custom_net.id // this will be calculate automatic
    # depends_on = [civo_network.custom_net]  // this is to wait for the creation of the network  
}
```

### Create a Volume

```bash
resource "civo_volume" "custom_volume" {
    name = "backup-data"
    size_gb = 60
    bootable = false
    # instance_id = civo_instance.my-test-instance.id // this will be calculate automatic
    # depends_on = [civo_instance.my-test-instance] // this is to wait for the creation of the instances  
}
```
