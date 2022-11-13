







# 分步走

- [x] 先跑起来

- [x] 融合 cobra
- [x] 融合 flag
- [ ] 融合 viper
- [ ] 将 flag 整理出来，避免手动 bing flag 进去 viper
- [ ] 将 flag 按照不同功能来构建成一个个 app 的可选 option 功能。让 app 能够灵活使用











# 需求

- [ ] flag
  - [ ] app 拥有单独的 flag 设置能力
  - [ ] framework 可以设置通用的 flag
  - [ ] app 能够 disable 掉部分 framework 设置的 flag
  - [ ] 每个 app 的二进制文件都应当由相应 flag 的 help 信息
  - [ ] framework 提供一个 cmd，生成默认配置文件
  - [ ] 输入参数应当是能够进行校验的，这样就不用每次用的时候都检查一次
- [ ] 日志
  - [ ] framework 能够内部自动默认初始化 log 的能力，并提供给 app 使用
  - [ ] app 可以设置 framework 提供 log 的：
    - [ ] 输出方向
    - [ ] 格式
    - [ ] 文件名
    - [ ] 日志级别
    - [ ] 按级别提供 hook
- [ ] 代码目录要易于拓展，不同 app 的代码要相互独立，不能混杂的放在一起
- [ ] MarineSnow 设计上是可以创建多个 app 的，所以 app 相关的能力需要独立为一个单独的 pkg
- [ ] makefile 要能够分别拥有多个 app 同时编译、单独编译的能力
- [ ] app 要拥有独立的配置文件
- [ ] 敏感配置项通过环境变量的方式进行注入
- [ ] 配置项应当有层次，易于拓展，不易混淆
- [ ] 自动创建 app 模板文件（phase 3 再考虑）
  - [ ] 考虑将 app 独立抽象出来（phase 3 再考虑）
- [ ] 暂时不支持 args
- [ ] 多入口函数 app（通常是 ctl 工具），第三阶段在考虑拓展
- [ ] 单入口函数 app，通常是 service







# 调研

## cobra + viper + pflag









**有没必要把 App 抽象为一个 interface ？**

**好处：**

- 整个 framework 可以直接操作 `App` interface
- 可以通过编译强制要求每一个 app 的作者都实现了 framework 对于 `App` 这一抽象概念的操作 method
  - 同时通过 interface 的方式，约定每一个 app 所必须具有的能力










- 框架不应该有配置文件，但是可以通过启动 CLI 来控制框架启动时的行为、框架相应的可配置选项
  - 框架不需要独立的配置文件，框架的最终目的是生成 app，帮助 app 工作。框架的配置，直接跟 app 的配置文件在一起就完事了。
    - 因为框架不会是一个独立的 app，跟生成的那个 app 进行交互的
    - 框架跟 app 的交互，总是发生在代码层面的，而不是运行时，app 跟一个框架进程交互
  - 框架可以有默认配置（直接写在代码里面就好了）
  - 但是框架的默认配置，可以被 app 的配置文件、app 启动时的 CLI 所覆盖
  - 优先级：CLI > configuration file > framework default config (只有 framework 的 config 会有 default config)



- 配置项要进行分组，不然就会十分混乱（推荐使用 `主体.配置项` 的形式）
- 敏感配置项应当通过环境变量的方式进行注入








- 对于 app 来说，不需要区分你这个配置究竟是文件输入的、环境变量输入的，还是 CLI 输入的，对 app 来说，这就是一个 option
- 但是对于 framework 不一样，framework 需要区分，然后归并为一个大的 option，给到 app 使用








**如何暴露 app 配置给 app**
- 所有配置，必须要能够被 app 通过变量、函数的形式访问到


**如何暴露 framework 配置给 app**


**app 与框架的启动顺序**
1. app 先设置相应可配置选项，运行函数
2. 框架处理 app 输入的 CLI、配置文件
3. 框架完成通用组件的初始化、map struct 工作（如：初始化相应的 log、读取 config 文件）
4. app 拿到框架配置、app 特有配置的 struct
6. app 完成 app 专属配置的初始化
7. app 根据 app 专属配置改变部分框架的配置



**cmd 的分层**
app 的本身应当是一个默认 cmd，然后 sub cmd 都是基于这个 app basic cmd 的
framework 只能创建 root-cmd，不应该在 framework 的代码中创建 sub-cmd，哪怕这是通用的。而是应该由 app 通过 `WithXXXX()` 的方式选择一些 framework 提供的 sub-cmd





**性能要求**

没有！只要你不是慢的离谱，人基本是感受不到的。而且 app 启动只会发生一次，这个启动过程中还不 ready 处理业务的，所以慢一点也是可以接受的。

唯一比较需要性能的场合，那应该是 HA 的场景，需要 app 快速重启，迅速就位，响应业务请求。

综上所述，app 启动框架，个人更倾向于灵活、容易维护、拓展简单、代码重复度低的方案，而不是一个性能极好的方案。





**options 跟 config 应当分开**

option 是输入的，而 app runtime 使用的应当是最后的 config。config 会将 option merge 进来

option（代码中写死的配置项）、command、flag、环境变量、配置文件，这所有的所有，最终是为了生成配置对应的结构体







**app 要不要对 `cobra.Command` 再封装一层？**

没有必要，因为你玩不出什么新的花样。本质还是要设置 flag，跟原本的三件套没有区别。除非，你的 Command，能够自动完成 merge







==TODO== 把 app 启动的函数 call flow 画出来，不同函数是由那个 package 调用的。这样就可以很清晰的看到：不同角色在不同函数层次干些什么New







注意，先把 pflag 跟 cobra bind 起来，不然 cmd 怎么处理 flag 呢？怎么生成 help 信息呢？

然后，在 command 实际执行的时候，再把 viper 跟 pflag bind 起来。再进行 config 的 merge








