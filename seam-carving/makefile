seam: main.cpp
	g++ -g -o $@ $<

.PROXY:sample
sample: seam
	cat ./sample.txt | ./$<

.PROXY:clean
clean:
	rm -f ./seam