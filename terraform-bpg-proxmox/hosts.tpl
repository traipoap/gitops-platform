[load-balance]
${ super_ip } ansible_user=traipoap

[k3s-masters]
${ master_ip } ansible_user=traipoap node_ip=${ master_ip }

[k3s-workers]
${ worker_ip } ansible_user=traipoap node_ip=${ worker_ip }

[k3s-cluster:children]
k3s-masters
k3s-workers