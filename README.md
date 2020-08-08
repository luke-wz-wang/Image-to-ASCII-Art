# Image-to-ASCII-Art
An image conversion system where a series of images can be converted to ASCII images

## Description

The project created a small image conversion system where a series of images are read and are converted to ASCII image. The ASCII image is not an actual image but a .txt file where the image is composed by ASCII symbols. As limited by the window size of file reader, the image to be converted has to be relatively small so that each line of the converted ASCII image can be expanded without auto line breaks.


For the simplicity of the testing, the system will create image conversion requests automatically based on the given dataset, and then ASCII images (.txt files) will be generated and output to a destinated path. The image conversion system has a sequential version as well as a parallel version for generating and handling conversion requests. Work balancing is also implemented in the parallel version. The goal of the experiment is to measure the performance of these versions and analyze the speedup led by parallel execution.


The actual conversion of the image is pixel-based conversion. At first, a list of potential symbols to be used in the ASCII image are defined. For each pixel of the image, the program takes out the value from green channel and use it to calculate the index of the corresponding symbol in the predefined list. In this way, a line of pixels is converted to a line of characters. Eventually, the entire image will be converted to a string and written out to a text file.


An image conversion request contains an in/source path containing the image, an out/destination path to store the generated text files. The input images are stored in the input folder under image folder; the output text files are stored in the output folder under image folder. A group of 10 conversion requests can be generated by calling the generateAndRequest() method. A command line flag -m is used to define the magnitude of conversion requests. In the parallel version, the main goroutine will spawn m goroutines to execute generateAndRequest(), and a total of 10*m conversion requests are made during the process. In the sequential version, the main goroutine will simple generate 10*m conversion requests and process the requests one by one.


As mentioned earlier, the image to be converted has to be relatively small due to the limitation of the file reader’s window size. The data decomposition here seems to have relatively small impact on the improvement of program performance. Therefore, only functional decomposition is implemented here. As the image conversion requests are pushed to the channel, these tasks are distributed to several workers to complete. An optional command line flag -p is used to define the number of workers. If the value of p is 0, the system will execute the sequential version of the program; otherwise, the parallel version will be executed. A work balancer as well as worker pool are implemented to maintain the balance of the work distribution, and the details will be discussed in the following session, advanced features.


In the parallel version, the main goroutine will exit after all the worker goroutines have returned, which means all the conversion requests are finished. In the sequential version, the system will handle the conversion tasks one by one and exit when the last request is finished.


## Advanced Features

As mentioned above, to balance the work distribution, a work balancer and worker pool are implemented. The worker pool is an array of conversion workers, and the number of the workers are defined by the command line flag -p. The work pool implements the heap interface and server as a component of the work balancer.


The work balancer utilizes the heap data structure to dispatch workers. Each worker has a variable called waiting, which represents the number of pending works the worker is working on. The top of the heap will be the worker with the least amount of work. The work balancer receives conversion requests from channel and distributes the task to the worker with the least amount of waiting/pending work.


Moreover, to dynamically demonstrate the work balancer is working properly, the amount of waiting tasks of each worker goroutine are printed out to the console whenever a new request is to be dispatched or a worker finished its work.
