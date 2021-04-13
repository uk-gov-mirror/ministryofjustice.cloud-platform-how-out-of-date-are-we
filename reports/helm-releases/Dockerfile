FROM python:alpine3.12

RUN pip install awscli \
  && apk add bash git curl ruby
RUN gem install json_pure

ENV \
  NOVA_VERSION=2.3.1 \
  KUBECTL_VERSION=1.17.12

# kubectl
RUN curl -sLo /usr/local/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
  && chmod 755 /usr/local/bin/kubectl

# Install nova

RUN curl -sL https://github.com/FairwindsOps/nova/releases/download/${NOVA_VERSION}/nova_${NOVA_VERSION}_linux_amd64.tar.gz > nova.tar.gz \
    && tar -xvf nova.tar.gz \
    && mv nova /usr/local/bin/ \
    && chmod 755 /usr/local/bin/nova

RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

WORKDIR /app

COPY bin ./bin

RUN chown 1000:1000 /app
USER 1000

# # Install helm-whatup
# RUN helm plugin install https://github.com/fabmation-gmbh/helm-whatup

# # Add helm repos
# RUN helm repo add jetstack https://charts.jetstack.io/; \
#   helm repo add concourse https://concourse-charts.storage.googleapis.com/; \
#   helm repo add cloud-platform https://ministryofjustice.github.io/cloud-platform-helm-charts; \
#   helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx; \
#   helm repo add prometheus https://prometheus-community.github.io/helm-charts; \
#   helm repo add autoscaler https://kubernetes.github.io/autoscaler; \
#   helm repo add vmware-tanzu https://vmware-tanzu.github.io/helm-charts; \
#   helm repo add banzaicloud https://kubernetes-charts.banzaicloud.com; \
#   helm repo add kiam https://uswitch.github.io/kiam-helm-charts/charts; \
#   helm repo add bitnami https://charts.bitnami.com/bitnami; \
#   helm repo add kube-state-metrics https://kubernetes.github.io/kube-state-metrics; \
#   helm repo add grafana https://grafana.github.io/helm-charts; \
#   helm repo add eks https://aws.github.io/eks-charts; \
#   helm repo add hashicorp https://helm.releases.hashicorp.com;  \
#   helm repo add stable https://charts.helm.sh/stable