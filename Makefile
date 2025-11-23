run:
	templ generate && go run .

build:
	templ generate
	rsync -r ../weblights lights:/home/erikrodabaugh/