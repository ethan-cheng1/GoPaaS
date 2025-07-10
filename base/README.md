# Base  
Cap teacher's course purchase link on MOOC.com: http://www.imooc.com/t/6512963

Current code name: Base  

##1.Quickly create code repository using the command below
```
sudo docker run --rm -v $(pwd): $(pwd) -w  $(pwd) -e ICODE=xxxxxx cap1573/cap-tool new git.imooc.com/cap1573/base

Note:
1.sudo - if you're on Mac system, the password prompt is for your local machine password.
2.In the above command ICODE=xxxxxx, "xxxxxx" is your personal purchased icode.
3.After purchasing the course, please use your computer to click into the learning course page to get the icode.
4.Please do not share the same icode with multiple people (will be detected and banned by MOOC.com).
5.The git.imooc.com/cap1573/base repository name here needs to match the go mod name
```
 

##2.Generate Go basic code automatically based on proto
```
make proto
```

##3.Compile existing Go code based on the code  
```
make build
```
After execution, this will generate a base binary file

##4.Compile and execute binary file
```
make docker
```
After successful compilation, it will automatically generate a base:latest image
You can use docker images | grep base to check if it was generated

##5.This course uses go-micro v3 version as the microservice development framework
Framework address: https://github.com/asim/go-micro


