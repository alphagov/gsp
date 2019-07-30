# Build GSP locally on Ubuntu using Kind - https://github.com/kubernetes-sigs/kind

# Run from the root of the repo as
# scripts/linux-kind.sh

# Install Docker, Helm, Go, Kubectl and then Kind
sudo apt-get -y install docker.io git


if id | grep docker
then
	echo "Current user is already in the group 'docker'"
else
	sudo adduser $(whoami) docker
	echo "Current user added to group 'docker'"
	echo "Please logout, login again and rerun this script"
	exit 1
fi

sudo systemctl enable docker
sudo systemctl start docker
sudo snap install helm --classic
sudo snap install kubectl --classic
sudo snap install go --classic

CURRENTDIR=$(pwd)
TMPDIR=$(mktemp -d)
cd ${TMPDIR}
git clone https://github.com/kubernetes-sigs/kind
cd kind/
git checkout v0.2.1
go install
rm -rf ${TMPDIR}/kind
cd ${CURRENTDIR}

PATH="${HOME}/go/bin:${PATH}"

# Now install GSP
./scripts/gsp-local-linux-kind.sh delete
./scripts/gsp-local-linux-kind.sh create

