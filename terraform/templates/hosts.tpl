[load-balance]
%{ for ip in super_nodes ~}
${ ip } ansible_user=traipoap
%{ endfor ~}

[k3s-masters]
%{ for ip in master_nodes ~}
${ ip } ansible_user=traipoap
%{ endfor ~}

[k3s-workers]
%{ for ip in worker_nodes ~}
${ ip } ansible_user=traipoap
%{ endfor ~}

[k3s-cluster:children]
k3s-masters
k3s-workers
