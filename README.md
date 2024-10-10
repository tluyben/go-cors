# 🌈 Go CORS Proxy 🚀

## 🎭 What is this?

This is a simple yet powerful CORS proxy written in Go. It allows you to bypass CORS restrictions when making requests from your frontend to a backend that doesn't have CORS configured.

## 🌟 Features

- 🔒 Handles CORS headers automatically
- 🚦 Supports all HTTP methods
- 📝 Detailed logging for easy debugging
- 🎛️ Configurable allowed origins

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

Run the proxy with the following command:

```
./go-cors -port <port> -listen <ip> "<allowed_origins>" <backend_url>
```

- `<port>`: The port number to listen on (default: 4001)
- `<ip>`: The IP address to listen on (default: localhost)
- `<allowed_origins>`: Comma-separated list of allowed origins, or "\*" for all origins
- `<backend_url>`: The URL of your backend server

Example:

```
./go-cors -port 4001 -listen 127.0.0.1 "*" http://localhost:4000
```

## 🎯 How it works

1. 📡 The proxy receives requests from your frontend
2. 🧙‍♂️ It adds the necessary CORS headers to the response
3. 🚚 It forwards the request to your backend
4. 🎁 It sends the response back to your frontend with the CORS headers

## 🐛 Debugging

If you're experiencing issues, check the console output. The proxy logs detailed information about each request and response.

## 🙋‍♀️ Need Help?

If you're still having CORS issues, make sure:

1. 🔍 Your frontend is making requests to the proxy URL, not directly to the backend
2. 🖥️ The backend URL provided to the proxy is correct
3. 📊 Check the browser console and network tab for any error messages

## 🎉 Happy coding!

Remember, with great power comes great responsibility. Use this proxy wisely and may your CORS errors be forever vanquished! 🏆
