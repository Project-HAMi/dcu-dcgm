# 1. 使用轻量级的基础镜像
FROM ubuntu:22.04

# 2. 设置工作目录为 /root
WORKDIR /root

# 3. 复制已编译好的二进制文件到 /usr/local/bin 目录
COPY dcgm-dcu /usr/local/bin/dcgm-dcu

# 4. 复制 .so 依赖库到 /usr/local/bin 目录
COPY pkg/dcgm/lib/librocm_smi64.so.2.8 /usr/local/bin/lib/librocm_smi64.so.2.8
COPY pkg/dcgm/lib/libhydmi.so.1.4 /usr/local/bin/lib/libhydmi.so.1.4

# 5. 为 .so 文件设置 755 权限
RUN chmod +x /usr/local/bin/lib/librocm_smi64.so.2.8 /usr/local/bin/lib/libhydmi.so.1.4

# 6. 设置软链接
RUN ln -s /usr/local/bin/lib/librocm_smi64.so.2.8 /usr/local/bin/lib/librocm_smi64.so.2 \
    && ln -s /usr/local/bin/lib/librocm_smi64.so.2 /usr/local/bin/lib/librocm_smi64.so \
    && ln -s /usr/local/bin/lib/libhydmi.so.1.4 /usr/local/bin/lib/libhydmi.so.1 \
    && ln -s /usr/local/bin/lib/libhydmi.so.1 /usr/local/bin/lib/libhydmi.so

# 7. 设置 LD_LIBRARY_PATH 环境变量以查找共享库
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/bin/lib

# 8. 确保二进制文件具有可执行权限
RUN chmod +x /usr/local/bin/dcgm-dcu

# 9. 暴露服务端口 16081
EXPOSE 16081

# 10. 启动服务，并将日志写入文件
CMD ["sh", "-c", "/usr/local/bin/dcgm-dcu -logtostderr -v=2 > /usr/local/bin/dcgm.log 2>&1"]
