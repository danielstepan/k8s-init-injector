# Use CFSSL to generate certificates
docker run -it --rm -v ${PWD}:/app -w /app debian bash

# Run the rest in the container
apt-get update && apt-get install -y curl &&
curl -L https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssl_1.5.0_linux_amd64 -o /usr/local/bin/cfssl && \
curl -L https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssljson_1.5.0_linux_amd64 -o /usr/local/bin/cfssljson && \
chmod +x /usr/local/bin/cfssl && \
chmod +x /usr/local/bin/cfssljson

# Generate ca in /tmp
cfssl gencert -initca ./tls/ca-csr.json | cfssljson -bare /tmp/ca

# Generate certificate in /tmp
cfssl gencert \
  -ca=/tmp/ca.pem \
  -ca-key=/tmp/ca-key.pem \
  -config=./tls/ca-config.json \
  -hostname="k8s-init-injector-webhook,k8s-init-injector-webhook.default.svc,k8s-init-injector-webhook.default.svc.cluster.local,localhost,127.0.0.1" \
  -profile=default \
  ./tls/ca-csr.json | cfssljson -bare /tmp/k8s-init-injector-webhook

# Make a secret
cat <<EOF > ./tls/k8s-init-injector-webhook-tls.yaml
apiVersion: v1
kind: Secret
metadata:
  name: k8s-init-injector-webhook
type: Opaque
data:
  tls.crt: $(cat /tmp/k8s-init-injector-webhook.pem | base64 | tr -d '\n')
  tls.key: $(cat /tmp/k8s-init-injector-webhook-key.pem | base64 | tr -d '\n') 
EOF

# Generate CA Bundle for MutatingWebhookConfiguration
ca_pem_b64="$(openssl base64 -A <"/tmp/ca.pem")" 
echo "$ca_pem_b64"
