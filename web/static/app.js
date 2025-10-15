class ChatApp {
    constructor() {
        this.socket = null;
        this.currentUser = null;
        this.currentRoom = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 1000;
        
        this.initializeElements();
        this.bindEvents();
        this.loadAvailableRooms();
    }
    
    initializeElements() {
        this.connectionPage = document.getElementById('connectionPage');
        this.chatPage = document.getElementById('chatPage');
        
        this.connectionForm = document.getElementById('connectionForm');
        this.usernameInput = document.getElementById('username');
        this.roomSelect = document.getElementById('roomSelect');
        this.customRoomGroup = document.getElementById('customRoomGroup');
        this.customRoomInput = document.getElementById('customRoom');
        this.connectBtn = document.getElementById('connectBtn');
        
        this.roomTitle = document.getElementById('roomTitle');
        this.statusIndicator = document.getElementById('statusIndicator');
        this.statusText = document.getElementById('statusText');
        this.userCount = document.getElementById('userCount');
        this.usersList = document.getElementById('usersList');
        this.messagesContainer = document.getElementById('messagesContainer');
        this.messageForm = document.getElementById('messageForm');
        this.messageInput = document.getElementById('messageInput');
        this.sendBtn = document.getElementById('sendBtn');
        this.charCount = document.getElementById('charCount');
        
        this.roomsList = document.getElementById('roomsList');
    }
    
    bindEvents() {
        this.connectionForm.addEventListener('submit', (e) => this.handleConnect(e));
        this.roomSelect.addEventListener('change', () => this.handleRoomSelectChange());
        
        this.messageForm.addEventListener('submit', (e) => this.handleSendMessage(e));
        this.messageInput.addEventListener('input', () => this.updateCharCount());
        
        window.addEventListener('beforeunload', () => {
            if (this.socket) {
                this.socket.close();
            }
        });
    }
    
    handleRoomSelectChange() {
        const isCustom = this.roomSelect.value === 'custom';
        this.customRoomGroup.style.display = isCustom ? 'block' : 'none';
        
        if (isCustom) {
            this.customRoomInput.required = true;
            this.customRoomInput.focus();
        } else {
            this.customRoomInput.required = false;
        }
    }
    
    handleConnect(e) {
        e.preventDefault();
        
        const username = this.usernameInput.value.trim();
        const roomSelectValue = this.roomSelect.value;
        const customRoom = this.customRoomInput.value.trim();
        
        if (!username) {
            this.showError('Введите имя пользователя');
            return;
        }
        
        let roomName;
        if (roomSelectValue === 'custom') {
            if (!customRoom) {
                this.showError('Введите название комнаты');
                return;
            }
            roomName = customRoom;
        } else {
            roomName = this.roomSelect.options[this.roomSelect.selectedIndex].text;
        }
        
        this.connectToChat(username, roomName);
    }
    
    connectToChat(username, roomName) {
        this.connectBtn.disabled = true;
        this.connectBtn.textContent = 'Подключение...';
        
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws?username=${encodeURIComponent(username)}&room=${encodeURIComponent(roomName)}`;
        
        try {
            this.socket = new WebSocket(wsUrl);
            this.currentUser = username;
            this.currentRoom = roomName;
            
            this.socket.onopen = () => this.handleSocketOpen();
            this.socket.onmessage = (event) => this.handleSocketMessage(event);
            this.socket.onclose = () => this.handleSocketClose();
            this.socket.onerror = (error) => this.handleSocketError(error);
            
        } catch (error) {
            console.error('Ошибка подключения:', error);
            this.showError('Ошибка подключения к серверу');
            this.connectBtn.disabled = false;
            this.connectBtn.textContent = 'Подключиться';
        }
    }
    
    handleSocketOpen() {
        console.log('WebSocket соединение установлено');
        this.updateConnectionStatus(true);
        this.showChatPage();
        this.reconnectAttempts = 0;
        this.connectBtn.disabled = false;
        this.connectBtn.textContent = 'Подключиться';
    }
    
    handleSocketMessage(event) {
        try {
            const data = JSON.parse(event.data);
            this.processMessage(data);
        } catch (error) {
            console.error('Ошибка парсинга сообщения:', error);
        }
    }
    
    processMessage(data) {
        switch (data.type) {
            case 'message':
                this.displayMessage(data.data);
                break;
            case 'history':
                this.displayMessageHistory(data.data);
                break;
            case 'error':
                this.showError(data.message);
                break;
            default:
                console.log('Неизвестный тип сообщения:', data);
        }
    }
    
    displayMessage(message) {
        const messageElement = this.createMessageElement(message);
        this.messagesContainer.appendChild(messageElement);
        this.scrollToBottom();
        
        messageElement.classList.add('fade-in');
    }
    
    displayMessageHistory(data) {
        this.messagesContainer.innerHTML = '';
        
        if (data.messages) {
            data.messages.forEach(message => {
                const messageElement = this.createMessageElement(message);
                this.messagesContainer.appendChild(messageElement);
            });
        }
        
        if (data.users) {
            this.updateUsersList(data.users);
        }
        
        this.scrollToBottom();
    }
    
    createMessageElement(message) {
        const messageDiv = document.createElement('div');
        messageDiv.className = 'message';
        
        if (message.type === 'system' || message.type === 'join' || message.type === 'leave') {
            messageDiv.classList.add('system');
            messageDiv.innerHTML = `
                <div class="message-content">
                    ${this.escapeHtml(message.content)}
                </div>
            `;
        } else {
            const isOwn = message.userId === this.currentUser;
            messageDiv.classList.add(isOwn ? 'own' : 'other');
            
            messageDiv.innerHTML = `
                <div class="message-content">
                    ${this.escapeHtml(message.content)}
                </div>
                <div class="message-info">
                    <span class="username">${this.escapeHtml(message.username)}</span>
                    <span class="timestamp">${this.formatTime(message.timestamp)}</span>
                </div>
            `;
        }
        
        return messageDiv;
    }
    
    updateUsersList(users) {
        this.userCount.textContent = users.length;
        this.usersList.innerHTML = '';
        
        users.forEach(user => {
            const userElement = document.createElement('div');
            userElement.className = 'user-item slide-in';
            
            const firstLetter = user.username.charAt(0).toUpperCase();
            userElement.innerHTML = `
                <div class="user-avatar">${firstLetter}</div>
                <span class="user-name">${this.escapeHtml(user.username)}</span>
            `;
            
            this.usersList.appendChild(userElement);
        });
    }
    
    handleSendMessage(e) {
        e.preventDefault();
        
        const message = this.messageInput.value.trim();
        if (!message || !this.socket || this.socket.readyState !== WebSocket.OPEN) {
            return;
        }
        
        const messageData = {
            content: message,
            roomId: this.currentRoom
        };
        
        this.socket.send(JSON.stringify(messageData));
        this.messageInput.value = '';
        this.updateCharCount();
    }
    
    updateCharCount() {
        const count = this.messageInput.value.length;
        this.charCount.textContent = `${count}/500`;
        
        if (count > 450) {
            this.charCount.style.color = '#f56565';
        } else if (count > 400) {
            this.charCount.style.color = '#ed8936';
        } else {
            this.charCount.style.color = '#718096';
        }
    }
    
    handleSocketClose() {
        console.log('WebSocket соединение закрыто');
        this.updateConnectionStatus(false);
        this.attemptReconnect();
    }
    
    handleSocketError(error) {
        console.error('WebSocket ошибка:', error);
        this.showError('Ошибка соединения');
    }
    
    attemptReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            this.statusText.textContent = `Переподключение... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`;
            
            setTimeout(() => {
                if (this.currentUser && this.currentRoom) {
                    this.connectToChat(this.currentUser, this.currentRoom);
                }
            }, this.reconnectDelay * this.reconnectAttempts);
        } else {
            this.statusText.textContent = 'Соединение потеряно';
            this.showError('Не удалось восстановить соединение. Обновите страницу.');
        }
    }
    
    updateConnectionStatus(connected) {
        if (connected) {
            this.statusIndicator.className = 'status-indicator status-connected';
            this.statusText.textContent = 'Подключен';
        } else {
            this.statusIndicator.className = 'status-indicator status-disconnected';
            this.statusText.textContent = 'Отключен';
        }
    }
    
    showChatPage() {
        this.connectionPage.style.display = 'none';
        this.chatPage.style.display = 'block';
        this.roomTitle.textContent = this.currentRoom;
        this.messageInput.focus();
    }
    
    showConnectionPage() {
        this.chatPage.style.display = 'none';
        this.connectionPage.style.display = 'block';
        this.socket = null;
        this.currentUser = null;
        this.currentRoom = null;
    }
    
    scrollToBottom() {
        this.messagesContainer.scrollTop = this.messagesContainer.scrollHeight;
    }
    
    formatTime(timestamp) {
        const date = new Date(timestamp);
        return date.toLocaleTimeString('ru-RU', {
            hour: '2-digit',
            minute: '2-digit'
        });
    }
    
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
    
    showError(message) {
        const errorDiv = document.createElement('div');
        errorDiv.className = 'error-notification';
        errorDiv.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            background: #f56565;
            color: white;
            padding: 15px 20px;
            border-radius: 8px;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
            z-index: 1000;
            animation: slideIn 0.3s ease-out;
        `;
        errorDiv.textContent = message;
        
        document.body.appendChild(errorDiv);
        
        setTimeout(() => {
            errorDiv.remove();
        }, 5000);
    }
    
    async loadAvailableRooms() {
        try {
            const response = await fetch('/api/rooms');
            const data = await response.json();
            
            if (data.success && data.rooms) {
                this.displayAvailableRooms(data.rooms);
            }
        } catch (error) {
            console.error('Ошибка загрузки комнат:', error);
        }
    }
    
    displayAvailableRooms(rooms) {
        this.roomsList.innerHTML = '';
        
        rooms.forEach(room => {
            const roomElement = document.createElement('div');
            roomElement.className = 'room-item';
            
            roomElement.innerHTML = `
                <span class="room-name">${this.escapeHtml(room.name)}</span>
                <span class="room-count">${room.userCount}</span>
            `;
            
            this.roomsList.appendChild(roomElement);
        });
        
        setTimeout(() => this.loadAvailableRooms(), 30000);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    new ChatApp();
});
