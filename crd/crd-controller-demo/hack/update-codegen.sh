#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
# 代码生成器包的位置
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

# generate-groups.sh <generators> <output-package> <apis-package> <groups-versions>
#                    使用哪些生成器，可选值 deepcopy,defaulter,client,lister,informer，逗号分隔，all表示全部使用
#                    输出包的导入路径
#                    CR 定义所在路径
#                    API 组和版本
bash "${CODEGEN_PKG}"/generate-groups.sh "deepcopy,client,informer,lister" \
  crd-controller-demo/pkg/client crd-controller-demo/pkg/apis \
  stable:v1beata1 \
  --output-base "${SCRIPT_ROOT}" \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt

# 自动生成的源码头部附加的内容:
#   --go-header-file "${SCRIPT_ROOT}"/hack/custom-boilerplate.go.txt


