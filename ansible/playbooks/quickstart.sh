# SPDX-License-Identifier: Apache-2.0

ansible-playbook -i inventory/hosts playbooks/00-prerequisites.yml
ansible-playbook -i inventory/hosts playbooks/01-cluster-setup.yml
ansible-playbook -i inventory/hosts playbooks/02-servicemesh.yml
ansible-playbook -i inventory/hosts playbooks/03-storage-networking.yml
ansible-playbook -i inventory/hosts playbooks/04-garage-deploy.yml # Request user input for garage credentials (blank to generate)
ansible-playbook -i inventory/hosts playbooks/05-gitops-bootstrap.yml # Request user input for FluxCD and GitHub PAT
