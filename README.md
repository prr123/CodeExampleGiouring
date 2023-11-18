# Examples giouring

Code examples for library giouring.  

## kernelVer

program that deiplayse the library version and the linux kernel version.  
usage: ./kernelVer ver/verion | kernel  

## gioSetup

program that uses giouring equivalen of liburing to set up ans subsequently release submission and completion queues with n entries.  
usage: ./gioSetup [/entries=n] /dbg  

## gioSimpleRead

program that performs two reads of a file:  
  (1) read the file conventionally (os.ReadFile)   
  (2) read the file using giouring  
  (3) compare the two data buffers for equality  

## gioSimpleWrite

program that write a random buffer to two files:  
  (1) writes the buffer to a file conventially (fil.Write).  
  (2) writes the buffer to a different file (ext .alt) using giouring.  
