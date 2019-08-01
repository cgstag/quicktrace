Analysis for SimScale Technical Assignement

Details:
* The solution should be a Java, Scala or Go program, executable from the command line.
* The input should be read from standard input or a file (chooseable by the user)
* The output should be one JSON per line, written to standard output, or a file (chooseable by the user).
* As said, there can be lines out of order.
* There can be orphan lines (i.e., services with no corresponding root service); they should be tolerated but not included in the output (maybe summarized in stats, see below).
* Lines can be malformed, they should be tolerated and ignored.

Features:
* A nice command-line interface, that allows to specify inputs and outputs from and to files.
* Optionally report to standard error (or to a file), statistics about progress, lines consumed, line consumption rate, buffers, etc. Both at the end of the processing and during it.
* Optionally report to standard error (or to a file), statistics about the traces themselves, like number of orphan requests, average size of traces, average depth, etc.
* As the file could be quite big, try to do the processing using as many cores as the computer has, but only if the processing is actually speeded that way.

Inicial Idea :
 * Separate the program in two
    * io.Reader will scan each line of the log to insert on a Hashmap in which the key is traceId
    * Iterate through the map to build the tree for each trace
        * My idea is to recursively build subtrees
        

