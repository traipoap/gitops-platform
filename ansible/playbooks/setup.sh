ansible-playbook -i inventory/hosts playbooks/01-cluster-setup.yml
ansible-playbook -i inventory/hosts playbooks/02-servicemesh.yml
ansible-playbook -i inventory/hosts playbooks/03-gitops-bootstrap.yml # Prompt for GitHub PAT
ansible-playbook -i inventory/hosts playbooks/04-storage-networking.yml
ansible-playbook -i inventory/hosts playbooks/05-garage-deploy.yml
