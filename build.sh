export GO111MODULE=on
export GOSUMDB=off
export GOPROXY=https://goproxy.cn
go install omo.msa.third
mkdir _build
mkdir _build/bin

cp -rf /root/go/bin/omo.msa.third _build/bin/
cp -rf conf _build/
cd _build
tar -zcf msa.third.tar.gz ./*
mv msa.third.tar.gz ../
cd ../
rm -rf _build
