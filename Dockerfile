# ---- Stage 1: Builder ----
# Use the official Rust image as a builder with specific platform
FROM --platform=linux/amd64 rust:slim AS builder

# Set the working directory.
WORKDIR /usr/src/app

# Copy the Cargo files to cache dependencies.
COPY Cargo.toml Cargo.lock ./
# Create a dummy src/main.rs to build dependencies only.
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
RUN rm -f target/release/deps/goodbye-world*

# Copy the actual source code.
COPY src ./src

# Build the application for release.
RUN cargo build --release

# ---- Stage 2: Runner ----
# Use the same Rust slim image but only for runtime with specific platform
FROM --platform=linux/amd64 rust:slim

# Set the working directory.
WORKDIR /usr/src/app

# Copy the compiled binary from the builder stage.
COPY --from=builder /usr/src/app/target/release/goodbye-world .

# Expose the port the app runs on.
EXPOSE 3000

# Set the command to run the application.
CMD ["./goodbye-world"]