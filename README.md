# 🌈 Go CORS Proxy 🚀

## 🎭 What is this?

This is a simple yet powerful CORS proxy written in Go. It allows you to bypass CORS restrictions when making requests from your frontend to a backend that doesn't have CORS configured.

## 🌟 Features

- 🔒 Handles CORS headers automatically
- 🚦 Supports all HTTP methods
- 📝 Detailed logging for easy debugging
- 🎛️ Configurable allowed origins
- 🔧 Flexible configuration with YAML file support
- 🌐 Multiple backend routing based on request paths
- 🔄 Regular expression support for path matching

## 🛠️ Installation

1. Clone this repository:
   ```
   git clone https://github.com/yourusername/go-cors-proxy.git
   ```
2. Navigate to the project directory:
   ```
   cd go-cors-proxy
   ```
3. Build the proxy:
   ```
   go build -o go-cors main.go
   ```

## 🚀 Usage

Run the proxy with one of the following commands:

### Using a config file:

```
./go-cors -config config.yml -cors "<allowed_origins>"
```

### Using command-line arguments (backwards compatible):

```
./go-cors -cors "<allowed_origins>" <backend_url>
```

### With custom port and listen address:

```
./go-cors -port <port> -listen <ip> -cors "<allowed_origins>" -config config.yml
```

- `<port>`: The port number to listen on (default: 8080)
- `<ip>`: The IP address to listen on (default: 127.0.0.1)
- `<allowed_origins>`: Comma-separated list of allowed origins, or "\*" for all origins
- `<backend_url>`: The URL of your backend server (when not using a config file)

Example:

```
./go-cors -port 9000 -listen 0.0.0.0 -cors "*" -config config.yml
```

## 📄 Configuration File

Create a YAML file (e.g., `config.yml`) with the following structure:

```yaml
services:
  - path: /api
    backend: http://api-backend.com
  - path: /
    backend: http://main-backend.com
```

- `path`: The URL path to match (can be a regular expression)
- `backend`: The backend URL to forward the request to

## 🎯 How it works

1. 📡 The proxy receives requests from your frontend
2. 🧙‍♂️ It matches the request path against the configured services
3. 🚚 It forwards the request to the appropriate backend
4. 🔒 It adds the necessary CORS headers to the response
5. 🎁 It sends the response back to your frontend with the CORS headers

## 🐛 Debugging

If you're experiencing issues, check the console output. The proxy logs detailed information about each request and response.

## 🙋‍♀️ Need Help?

If you're still having CORS issues, make sure:

1. 🔍 Your frontend is making requests to the proxy URL, not directly to the backend
2. 🖥️ The backend URLs in your config file or command-line arguments are correct
3. 📊 Check the browser console and network tab for any error messages
4. 🔬 Verify that your path patterns in the config file match your intended routes

## 🎉 Happy coding!

Remember, with great power comes great responsibility. Use this proxy wisely and may your CORS errors be forever vanquished! 🏆
