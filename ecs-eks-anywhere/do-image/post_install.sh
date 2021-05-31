#!/usr/bin/env sh

# Add docker to user group
usermod -aG docker $USER
newgrp docker