seam: main.cpp
	g++ -o $@ $<

.PROXY:sample
sample: seam
	cat ./sample.txt | ./$<