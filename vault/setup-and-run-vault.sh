#!/bin/sh

export VAULT_DEV_ROOT_TOKEN_ID=$VAULT_TOKEN

vault server -dev -log-level=trace &

sleep 2

vault audit enable file file_path=stdout

vault policy write users /opt/users.hcl

vault auth enable approle

vault write auth/approle/role/$VAULT_APP_ROLE policies="users"
vault write auth/approle/role/$VAULT_APP_ROLE/role-id role_id=$VAULT_ROLE_ID

vault write auth/approle/role/$VAULT_APP_ROLE/custom-secret-id secret_id=$VAULT_SECRET_ID

vault kv put secret/api-berlin/$VAULT_APP_ROLE


echo "vault started and seeded."

wait
