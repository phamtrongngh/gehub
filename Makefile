all: build push
build: 
	docker build -t benalpha1105/gehub:all-in-one .
push: 
	docker push benalpha1105/gehub:all-in-one