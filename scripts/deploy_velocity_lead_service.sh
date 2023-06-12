#!/bin/bash
set -eE -o pipefail



# Assume an AWS IAM Role and extract the three secrets (access key, secret key, session token)
# as environment variables.
# Also set the global variable AWS_ACCOUNT_ID

assume_role() {
    local AWS_ACCOUNT="$1"
    AWS_ACCOUNT_ID=$(aws ssm get-parameter --name /target/${AWS_ACCOUNT}/aws-account-id --query "Parameter.Value" --output text)

    local ROLE_ARN_TO_ASSUME="arn:aws:iam::${AWS_ACCOUNT_ID}:role/cns-iam-role-eks-website-jenkins"
    local IDENTITY_SESSION=$(aws sts get-caller-identity | jq -r '.Arn' | awk -F'/' '{print $NF}')
    local ROLE_SESSION_NAME="${IDENTITY_SESSION}@$(date +%s)"
    local ASSUMED_CREDENTIALS=$(aws sts assume-role --region us-east-1 \
        --role-arn "${ROLE_ARN_TO_ASSUME}" \
        --role-session-name "${ROLE_SESSION_NAME}")

    # Ansible
    AWS_ACCESS_KEY=$(echo ${ASSUMED_CREDENTIALS} | jq -r '.Credentials.AccessKeyId')
    AWS_SECRET_KEY=$(echo ${ASSUMED_CREDENTIALS} | jq -r '.Credentials.SecretAccessKey')
    AWS_SECURITY_TOKEN=$(echo ${ASSUMED_CREDENTIALS} | jq -r '.Credentials.SessionToken')
    export AWS_ACCESS_KEY AWS_SECRET_KEY AWS_SECURITY_TOKEN

    # AWS CLI
    AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY
    AWS_SECRET_ACCESS_KEY=$AWS_SECRET_KEY
    AWS_SESSION_TOKEN=$AWS_SECURITY_TOKEN
    export AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY AWS_SESSION_TOKEN
}


# Authenticate to AWS ECR
ecr_login() {
    aws ecr get-login-password | docker login --username AWS --password-stdin "${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com"
}

build_and_deploy(){
    local CLUSTER_NAME="$1"
    echo "building and deploying"
    aws eks update-kubeconfig --name "${CLUSTER_NAME}"
    echo "building and deploying complete"
}

skaffold_job(){
    local CLUSTER_NAME="$1"
    local SKAFFOLD_PROFILE="$2"
    local environment_name="$3"
    echo "starting Skaffold_job"
    
    DEPLOYMENT_ENV=${environment_name} \
        skaffold run \
            --default-repo "${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com" \
            --kube-context "arn:aws:eks:${AWS_DEFAULT_REGION}:${AWS_ACCOUNT_ID}:cluster/${CLUSTER_NAME}" \
            --profile "${SKAFFOLD_PROFILE}"
}

main(){

    cd ..
    local aws_account_name=$1
    assume_role ${aws_account_name}


    local cluster_name=$2
    local skaffold_profile_name=$3
    local environment_name=$4

    ecr_login 
    build_and_deploy ${cluster_name}
    skaffold_job ${cluster_name} ${skaffold_profile_name} ${environment_name}
}


main $@