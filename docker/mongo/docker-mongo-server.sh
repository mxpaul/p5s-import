docker run --rm -it --name p5s -v "$(readlink -f data):/data/db" -p 27017:27017 library/mongo
