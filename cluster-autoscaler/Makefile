all: build

TAG?=dev
FLAGS=
LDFLAGS?=-s
ENVVAR=CGO_ENABLED=0
GOOS?=linux
REGISTRY?=staging-k8s.gcr.io
ifdef BUILD_TAGS
  TAGS_FLAG=--tags ${BUILD_TAGS}
  PROVIDER=-${BUILD_TAGS}
  FOR_PROVIDER=" for ${BUILD_TAGS}"
else
  TAGS_FLAG=
  PROVIDER=
  FOR_PROVIDER=
endif
ifdef LDFLAGS
  LDFLAGS_FLAG=--ldflags "${LDFLAGS}"
else
  LDFLAGS_FLAG=
endif

build: clean deps
	$(ENVVAR) GOOS=$(GOOS) go build ${LDFLAGS_FLAG} ${TAGS_FLAG} ./...
	$(ENVVAR) GOOS=$(GOOS) go build -o cluster-autoscaler ${LDFLAGS_FLAG} ${TAGS_FLAG}

build-binary: clean deps
	$(ENVVAR) GOOS=$(GOOS) go build -o cluster-autoscaler ${LDFLAGS_FLAG} ${TAGS_FLAG}

test-unit: clean deps build
	$(ENVVAR) go test --test.short -race ./... $(FLAGS) ${TAGS_FLAG}

dev-release: build-binary execute-release
	@echo "Release ${TAG}${FOR_PROVIDER} completed"

make-image:
ifdef BASEIMAGE
	docker build --pull --build-arg BASEIMAGE=${BASEIMAGE} \
	    -t ${REGISTRY}/cluster-autoscaler${PROVIDER}:${TAG} .
else
	docker build --pull -t ${REGISTRY}/cluster-autoscaler${PROVIDER}:${TAG} .
endif

push-image:
	./push_image.sh ${REGISTRY}/cluster-autoscaler${PROVIDER}:${TAG}

execute-release: make-image push-image

clean:
	rm -f cluster-autoscaler

generate:
	go generate ./cloudprovider/aws

format:
	test -z "$$(find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -s -d {} + | tee /dev/stderr)" || \
    test -z "$$(find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -s -w {} + | tee /dev/stderr)"

docker-builder:
	docker build -t autoscaling-builder ../builder

build-in-docker: clean docker-builder
	docker run -v `pwd`:/gopath/src/k8s.io/autoscaler/cluster-autoscaler/ autoscaling-builder:latest bash -c 'cd /gopath/src/k8s.io/autoscaler/cluster-autoscaler && BUILD_TAGS=${BUILD_TAGS} LDFLAGS="${LDFLAGS}" make build-binary'

release: build-in-docker execute-release
	@echo "Full in-docker release ${TAG}${FOR_PROVIDER} completed"

container: build-in-docker make-image
	@echo "Created in-docker image ${TAG}${FOR_PROVIDER}"

test-in-docker: clean docker-builder
	docker run -v `pwd`:/gopath/src/k8s.io/autoscaler/cluster-autoscaler/ autoscaling-builder:latest bash -c 'cd /gopath/src/k8s.io/autoscaler/cluster-autoscaler && go test -race ./... ${TAGS_FLAG}'

# kubectl -n kube-system --kubeconfig=/home/janitha/Documents/bb/DVOP/aws-k8s-terraform/config delete -f /home/janitha/go/src/autoscaler/cluster-autoscaler/cloudprovider/aws/examples/cluster-autoscaler-autodiscover.yaml
janitha: build-in-docker make-image
	docker tag staging-k8s.gcr.io/cluster-autoscaler:dev lwpeartifabld.tideworks.com:5000/kubetest/cluster-autoscaler:latest
	docker push lwpeartifabld.tideworks.com:5000/kubetest/cluster-autoscaler:latest
	kubectl -n kube-system --kubeconfig=/home/janitha/Documents/bb/DVOP/aws-k8s-terraform/config apply -f /home/janitha/go/src/autoscaler/cluster-autoscaler/cloudprovider/aws/examples/cluster-autoscaler-autodiscover.yaml
	./cluster-autoscaler --kubeconfig /var/run/kubernetes/admin.kubeconfig --cloud-provider=aws --node-group-auto-discovery=asg:tag=k8s.io/cluster-autoscaler/enabled,k8s.io/cluster-autoscaler/janitha-k8s --v=10

# https://aws.amazon.com/premiumsupport/knowledge-center/eks-cluster-autoscaler-setup/
.PHONY: all deps build test-unit clean format execute-release dev-release docker-builder build-in-docker release generate janitha

