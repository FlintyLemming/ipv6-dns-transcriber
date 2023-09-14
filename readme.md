# Cloudflare AAAA 记录同步工具

该工具用于自动抄写一个域名的 AAAA 解析结果，并设置为您 Cloudflare 下的另一个域名的 AAAA 解析结果。

## 原理

1. 获取IPv6地址: 程序会查询指定的 from 域名的 AAAA 记录，获取其IPv6地址。
2. 更新Cloudflare记录: 使用 Cloudflare 的 API，将另一个 to 域名的 AAAA 记录设置为上述获取的 IPv6 地址。
3. 定期检查: 该过程会根据设置的周期（如每分钟一次）定期进行。

## 使用场景

路由器或者 NAS 厂商提供了 ddns 域名，但是宽带没有 v4 公网，此时如果使用 CNAME 解析到自己的域名下，会将 A 记录同样抄写过来。使用此工具，可以只抄写 AAAA 记录。

您也可以不用大费周章这样再抄写一次，直接在设备上做 ddns v6 解析也是可以的，但是这个程序可以在任意有互联网的机器上执行，根据需求使用即可。

## 如何使用 Docker 镜像

1. 构建 Docker 镜像

    在 Dockerfile 所在的目录中执行：

    ```bash
    docker build -t your_image_name .
    ```

2. 运行 Docker 容器

    ```bash
    docker run -e ZONE_ID=your_zone_id \
    -e API_TOKEN=your_api_token \
    -e FROM_DOMAIN=your_from_domain \
    -e TO_DOMAIN=your_to_domain \
    -e CYCLE_MINUTES=your_cycle_minutes \
    -d your_image_name
    ```

## 环境变量说明

* ZONE_ID: 你的 Cloudflare Zone ID。
* API_TOKEN: 你的 Cloudflare API Token，用于进行域名记录的操作。
* FROM_DOMAIN: 需要被抄写的域名。
* TO_DOMAIN: 你希望更新 AAAA 记录的域名。
* CYCLE_MINUTES: 程序检查和更新的周期（以分钟为单位）。

## 注意事项

本工具功能实现为 [qwe7002](https://github.com/qwe7002)，我做了修改使其可以打包成容器镜像方便使用

请确保提供的 Cloudflare API Token 具有足够的权限来修改 DNS 记录。同时，为了安全考虑，不要在公开场合泄露您的 API Token 和 Zone ID。

希望这个 README 可以帮助你更好地理解和使用该工具。如果有任何问题或建议，请随时向我们反馈。