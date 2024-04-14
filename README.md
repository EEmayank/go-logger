# go-logger
Simple logging service in golang. I was always interested in kafka and its working and in my opinion, making a logs service yourself is first step in making your own kafka. So here I am.

# What is log?
- A log is an append-only squence of records typically read top to bottom, oldest to newest - similar to running tail -f on a file.
- A log is like a table that always orders the recods by time and indexes each record by its offset and time created. And this is one of the reason you would find logs and kafka very similiar

# Segmentation of logs
- Concrete implementation of logs have to deal with us not having disks with inifinite splace, which means we can't append to the same file forever. So we split the log into a list of segments. When the log grows too big, we free up disk space by deleting old segments whose data we've already processed or archived.
    - This cleaning can de done in background process.

# Active Segment
- There's always one special segment among the list of segments, `active segment`
- It's the only segment we only actively write to. When we've filled the active segment, we create a new segment and make it the active segment.

# Store File and Index File
- Each segment comprises of a store and an index file.
- Index file is where we index each recoed in the store file.
- The index file speeds up reads because it maps recoed offsets to their position in the store file.
- Reading a record given its offset is a two step process:
    1. Get the entry from the index file for the record, which tells you the position of the record in the store file.
    2. Read the recoed at that position in the store file.

# Why index file?
- Since index file requires only two small fields - the offset and stored position of the record
    - offset
    - stored position of the record
- This makes index files very small in size, small enough that we can memory-map them and make operations on the file as fast as operating on in-memory data
- This same benefit is used by the kafka.

# Terminology
1. Record - the data stored in our log
2. Store - the file we store records in
3. Index - the file we store index entries in
4. Segment - the abstraction that ties a store and an index together
5. Log - the abstraction that ties all the segment together

# Source Material
The code here refers to one of the best books for on-hand experience with distributed services in golang, details for which I am sharing below
"Distributed Services with Go"
                            - by Travis Jeffery
