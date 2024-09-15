FROM debian:latest

# All dependencies required to build the Rust modules
RUN apt-get update && apt-get install -y curl build-essential sudo m4

# Install Go
RUN curl -OL https://go.dev/dl/go1.22.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz && \
    rm go1.22.5.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"

# Install Rust
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

# Ensure required Rust libraries are cached
WORKDIR /app
RUN cargo new dummy_project
WORKDIR /app/dummy_project

# Add dependencies you need cached here
# RUN echo 'your_lib = "version" >> Cargo.toml'
RUN echo 'rand = "*"' >> Cargo.toml
RUN echo 'rug = "*"' >> Cargo.toml

RUN cargo build --release
RUN rm -rf /app/dummy_project
WORKDIR /app

# Add 'student' user for testing without permissions
RUN useradd -m student
RUN chmod 777 /root
USER student
RUN rustup default stable
USER root

RUN echo 'export PATH=$PATH:/root/.cargo/bin' >> /etc/profile.d/rust_path.sh

# Install 'cargo-valgrind' for testing leaks
RUN apt-get install -y valgrind
RUN /root/.cargo/bin/cargo install cargo-valgrind

COPY ./internal /app/internal
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY ./main.go /app/main.go

RUN go build .
