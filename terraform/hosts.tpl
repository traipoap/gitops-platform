[load-balance]
${ super_ip } ansible_user=traipoap

[k3s-masters]
${ master_ip } ansible_user=traipoap

[k3s-workers]
${ worker_ip } ansible_user=traipoap
[k3s-cluster:children]
k3s-masters
k3s-workers
