go install -ldflags "-X github.com/DCNT-Hammer/dcnt/engine.Build=`git rev-parse HEAD` -X github.com/DCNT-Hammer/dcnt/engine.dcntVersion=`cat VERSION`" -v
