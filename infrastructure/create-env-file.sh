#!/bin/bash


function create_var_if_not_exists() {
    ENV_KEY=$1
    ENV_VAL=$2
    UPDATE_VALUE=$3

    if grep -q $ENV_KEY .env; then
        echo "$ENV_KEY already created"
    else
        echo "No $ENV_KEY value found"
        echo "$ENV_KEY=$ENV_VAL" >> .env
        return
    fi

    if $UPDATE_VALUE; then
        echo "Updating $ENV_KEY"
        sed -i.bak "s/^$ENV_KEY=.*$/$ENV_KEY=$ENV_VAL/g" .env
    fi
}

touch .env

create_var_if_not_exists "JWT_SECRET" "$(openssl rand -hex 32)" false

create_var_if_not_exists "LITESTREAM_ACCESS_KEY_ID" "$LITESTREAM_ACCESS_KEY_ID" true
create_var_if_not_exists "LITESTREAM_ACCESS_KEY_SECRET" "$LITESTREAM_ACCESS_KEY_SECRET" true

create_var_if_not_exists "DOCKMON_USERNAME" "bolesta-admin" false
create_var_if_not_exists "DOCKMON_PASSWORD" "$(openssl rand -hex 16)" false

rm .env.bak
