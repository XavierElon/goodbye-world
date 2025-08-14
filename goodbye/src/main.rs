use axum::{response::Json, routing::get, Router, http::StatusCode, response::IntoResponse};
use serde_json::{json, Value};
use std::net::SocketAddr;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Build our application with routes
    let app = Router::new()
        .route("/", get(root_handler))
        .route("/goodbye", get(goodbye_handler))
        .fallback(not_found_handler);

    // Define the address to run the server on. 0.0.0.0 is important for Docker.
    let addr = SocketAddr::from(([0, 0, 0, 0], 3000));
    println!("�� Starting server on {}", addr);

    // Run the server with proper error handling
    let listener = tokio::net::TcpListener::bind(addr).await?;
    println!("✅ Server listening on {}", addr);
    println!("🌐 Server is now running and will continue to run...");
    
    // This line keeps the server running indefinitely
    axum::serve(listener, app).await?;
    
    // This should never be reached in a normal web server
    println!("⚠️  Server stopped unexpectedly");
    Ok(())
}

// Handler for the root route - shows available endpoints
async fn root_handler() -> Json<Value> {
    let response = json!({
        "message": "Welcome to the Goodbye World API!",
        "available_endpoints": [
            "GET / - This help message",
            "GET /goodbye - Returns a goodbye message"
        ],
        "status": "success"
    });
    Json(response)
}

// Handler function for the /goodbye route.
async fn goodbye_handler() -> Json<Value> {
    let response = json!({
        "message": "Goodbye, World!",
        "status": "success"
    });
    Json(response)
}

// 404 handler for non-existent routes
async fn not_found_handler() -> impl IntoResponse {
    let error_response = json!({
        "error": "Route not found",
        "message": "The requested endpoint does not exist",
        "available_endpoints": [
            "GET / - This help message",
            "GET /goodbye - Returns a goodbye message"
        ],
        "status": "error"
    });

    (StatusCode::NOT_FOUND, Json(error_response))
}