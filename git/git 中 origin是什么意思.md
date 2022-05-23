#### Git 中用到的 origin 是什么意思

经常使用 `Git`，但是突然发现不知道在`Git`中经常使用到的`origin`是什么意思，这两天突然想明白了

当你打算将本地的代码推送的远程的`Git`仓库的时候，通常会有两种方法：

1. 在远程创建一个仓库，`git clone`到本地，然后完成代码之后执行如下的命令就好了

   ```
   git add .
   git commit -m "Initial commit"
   git push -u origin master
   ```

2. 进入本地代码的根目录，然后执行一系列如下的命令

   ```
   git init
   git remote add origin https://xxx.git
   git add .
   git commit -m "Initial commit"
   git push -u origin master
   ```

   ###### 那命令中的这个`origin`是什么意思呢

   这个`origin`其实就是你要添加的这个远程仓库`http:://xxx.git`的名称，也就是说，当你执行`git remote add origin https://xxx.git`之后就会在本地的`.git`文件的配置中添加一条`origin`和`https://xxx.git`对应关系的记录，类似于这种

   ```
   origin: https://xxx.git
   ```

   为了证实这步，我就到`.git`目录中查看，果然有一个`./git/config`文件，在文件中就有一条为

   ```
   [remote "origin"]
   url = https://xxx.git
   fetch = +refs/heads/:refs/remotes/origin/
   ```

   这样就是为了简化执行`git`的操作，比如，推送代码的时候执行`git push origin master` 就好了，不然的话你就得把命令变为`git push https://xxx.git master` 了

   ###### 那为什么是`origin`呢？

   这是因为`origin`这个是`git` 默认的，当你执行`git clone https://xxx.git`的时候，就会自动为你创建一条名为`origin`，值为`https://xxx.git`的记录

   当然，你也可以更改这个值，比如，你要添加远程仓库的时候，你可以换个名字，比如：

   ```
   git remote add kkk https://xxx.git
   ```

   这时就是会在`.git/config`中有这样的一条记录了

   ```
   [remote "kkk"]
   url = https://gitlab.wise-paas.com/WISE-PaaS-4.0-Ops/lincese-checker.git
   fetch = +refs/heads/:refs/remotes/kkk/
   ```

   但你在执行`git`操作的时候，就需要变成`git push kkk master`，所以，通常还是使用默认的好