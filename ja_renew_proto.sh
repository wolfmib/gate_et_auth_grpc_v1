#!/bin/bash

echo "[Jean]: On va remove .gitmodule  .git/modules/proto pour toi "
echo "[Jean]: Ensuite , creer nouvelle latest version of proto "
echo "Ta sure ?"
echo "Enter ...."
read nothing_var


rm -rf proto

echo "------------------------------"
echo " rm -rf .git/modules/proto/ "
echo " rm .gitmodules             "
echo "------------------------------"
rm -rf .git/modules/proto/
rm .gitmodules




echo "[Jean]: Mnt. je suppliser tout le files"
echo "       tu dois commit back to gitlab first"
echo "Enter commit message = 'clean'"
echo "Are you ready ? "
read nothing_var
source ja_git_push_back_v2.sh


echo
echo
echo
echo "########################################################################"

echo "[Jean]: Cloning... "
echo "Example: "
echo "git submodule add git@gitlab.com:Johnny_Wick/ja_gitlab_grpc_proto.git proto"
git submodule add git@gitlab.com:Johnny_Wick/ja_gitlab_grpc_proto.git proto
