# friis-link：Go 自由空间链路预算 Web 服务（Friis + kTBF + 前端控制台）

给定载波频率、距离、发射功率与收发增益，按 Friis 方程核算 FSPL、接收功率与 SNR 余量；提供 `/api/budget` 与嵌入网页。

## 构建 / 运行 / 测试

```text
go build ./...
./friis-link -http :8080
curl -s http://127.0.0.1:8080/api/example
go run . budget example/sband-10km.json
go test ./...
```

## 评测镜像

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -d -P --name friis-link-b14 <image-name>:latest
curl -s http://127.0.0.1:$(docker port friis-link-b14 8080 | cut -d: -f2)/api/example
docker rm -f friis-link-b14
```
