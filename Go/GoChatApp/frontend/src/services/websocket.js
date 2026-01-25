export default class WebSocketService {
  constructor(url) {
    this.url = url;
    this.socket = null;
    this.listeners = [];
  }

  connect(onOpen, onClose) {
    this.socket = new WebSocket(this.url);

    this.socket.onopen = () => {
      console.log("WebSocket connected");
      if (onOpen) onOpen();
    };

    this.socket.onmessage = (event) => {
      const message = event.data;
      // Depending on backend, might be raw bytes or JSON.
      // Our backend broadcasts raw bytes currently, or JSON bytes.
      // Let's assume text/JSON.
      this.listeners.forEach((listener) => listener(message));
    };

    this.socket.onclose = () => {
      console.log("WebSocket disconnected");
      if (onClose) onClose();
    };

    this.socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };
  }

  sendMessage(message) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(message);
    } else {
      console.error("WebSocket is not open");
    }
  }

  addListener(callback) {
    this.listeners.push(callback);
  }

  removeListener(callback) {
    this.listeners = this.listeners.filter((l) => l !== callback);
  }

  close() {
    if (this.socket) {
      this.socket.close();
    }
  }
}
