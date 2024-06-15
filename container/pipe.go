package container

import "os"

// NewPipe 初始化Pipe。
// Pipe：内存管道，用于进程内的不同程序可以通过该管道来进行通信。
// 基本工作原理：通过系统调用创建Pipe，返回一个写一个读文件，写文件将要写入的数据写入到内核缓冲区，然后由读文件进行读取。
// 如果内核缓冲区满则写请求会进行阻塞，同理，如果内核缓冲区为空，则读请求则会进行阻塞。
func NewPipe() (*os.File, *os.File, error) {
	return os.Pipe()
}
