# friis-link — 自由空间链路预算核算

friis-link 是一个命令行自由空间链路预算核算工具。用户给出载波频率、链路距离、发射功率、收发天线增益，工具按 Friis 传输方程算出波长 λ、自由空间路径损耗 FSPL 与接收功率 Pr；再给出可选噪声参数（系统温度、带宽、噪声系数）与所需 SNR 时，进一步估算噪声底 N=kTBF、SNR 与衰落余量 M=SNR−SNR_min，并判定链路是否达标。

- 输入：JSON 场景（`frequency_hz`、`distance_km`、`tx_power_dbm`、`tx_gain_dbi`、`rx_gain_dbi`、可选 `extra_loss_db`、可选 `noise`），示例见 `example/sband-10km.json`
- 输出：波长、频段、FSPL、EIRP、接收功率，以及（配置了噪声时）噪声底、SNR、余量与达标判定；`budget -json` 可输出机器可读 JSON，`reverse` 子命令反推覆盖距离与所需 EIRP
- 能力边界：自由空间（Friis）模型，不包含多径、大气吸收与地形绕射；极化失配、指向误差等额外损耗只经 `extra_loss_db` 计入一次，不重复乘进公式；直线阵与阻抗匹配分别由 array-af 与 smith-lnet 覆盖，本仓不做方向图与圆图匹配

## 用法

```text
go run . budget example/sband-10km.json
```

输出（S 波段 2.45 GHz、10 km 示例）：

```text
friis-link budget
  frequency   : 2.45 GHz (S band)
  distance    : 10 km
  wavelength  : 0.122364 m
  FSPL        : 120.231 dB (linear 1.055e+12)
  EIRP        : 33.000 dBm
  Pr          : -84.231 dBm (3.775e-12 W)
  noise floor : -111.975 dBm (kTBF)
  SNR         : 27.744 dB
  min SNR     : 10.000 dB
  margin      : 17.744 dB
  link        : sufficient (healthy)
```

其他子命令：

```text
go run . budget -json example/sband-10km.json   # JSON 结果
go run . reverse example/sband-10km.json        # 灵敏度、最大覆盖距离、所需 EIRP
go run . compare a.json b.json                  # 两个场景的关键量差值
go run . help
```

非法输入（距离或频率非正、增益无效、带宽为 0、JSON 解析失败、未知字段）一律打印 stderr 并以非零退出码结束。

## 关键约定

- **钉死的常数**：光速 `c = 2.99792458e8 m/s`，玻尔兹曼常数 `k = 1.380649e-23 J/K`，全仓共用同一份声明。
- **dB 与线性换算**：增益/功率比统一走 `10*log10` 口径（20log 只出现在含距离或频率平方项的 FSPL 中）。线性与 dB 两套公式共享同一换算，接收功率用两条路径分别计算并交叉校验。
- **FSPL**：`FSPL = (4πd/λ)²`，`λ = c/f`；`FSPL_dB = 20·log10(4πd/λ)`，与闭式 `20·log10(d) + 20·log10(f) + 20·log10(4π/c)` 一致。
- **接收功率**：`Pr_dBm = Pt_dBm + Gt_dB + Gr_dB − FSPL_dB − L_extra`。
- **噪声底**：`N = k·T·B·F`；`SNR = Pr/N`（dB 域相减）。
- **余量判定**：`M = SNR − SNR_min`，`M < 0` 判定链路不足。
- **量纲**：距离进入公式前换算为米，频率使用 Hz；`extra_loss_db` 只在 dB 域减去一次。

## 构建与测试

```text
go build ./...
go test ./...
```

纯标准库实现，无第三方依赖，无 cgo。

## 许可

MIT，见 [LICENSE](./LICENSE)。
