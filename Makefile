run: pkged.go 
	go run main.go pkged.go -mode pkger

pkged.go: react-app/build
	pkger -o ./ -include /react-app/build


react-app/build:
	(cd react-app && npm run build)

clean:
	rm -rf react-app/build pkged.go