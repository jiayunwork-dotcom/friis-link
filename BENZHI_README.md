# friis-link — 自由空间链路预算核算

friis-link 是自由空间链路预算核算命令行工具：给定频率、距离、发射功率与收发增益，按 Friis 方程计算波长、FSPL 与接收功率；配置噪声参数时给出噪声底、SNR、余量与达标判定。纯标准库，无网络依赖，无 cgo。

## 构建 / 运行 / 测试

```text
go build ./...             # 编译
go run . budget example/sband-10km.json   # CLI：核算示例链路并打印 λ/FSPL/Pr/SNR/余量
go run . reverse example/sband-10km.json  # 反推覆盖距离与所需 EIRP
go test ./...              # 单元测试（model / propagation / noise / budget）
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```

进容器后运行 `go build ./... && go test ./...`，再用 `go run . budget example/sband-10km.json` 验证 CLI 输出。
