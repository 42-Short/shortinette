FROM rust:1.83-bullseye

RUN apt-get update && apt-get install -y valgrind strace

RUN groupadd -r shortinette && useradd -r -g shortinette -m shortinette
USER shortinette

# Needed if each test should be executed in a separate process
RUN cargo install cargo-nextest --locked

RUN cargo install cargo-valgrind

RUN rm -rf ~/.cargo/registry/cache/*
RUN rm -rf ~/.cargo/git/checkouts/*

# We want to cache some Rust libraries
# this will make the image bigger, but makes running tests faster
WORKDIR /tmp/dummy_project
RUN cargo init .
RUN cargo add 'rand@*' 'rug@*'
RUN cargo fetch
RUN cargo build --release
WORKDIR /
RUN rm -rf /tmp/dummy_project

WORKDIR /tmp/tester
COPY --chown=shortinette:shortinette Cargo.lock Cargo.toml ./
COPY --chown=shortinette:shortinette src/ ./src
RUN cargo fetch
RUN cargo build --release --bin tester
WORKDIR /app
RUN cp /tmp/tester/target/release/tester ./
RUN rm -rf /tmp/tester

CMD ["./tester"]
