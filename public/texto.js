(function () {
    /**
     * The TextoMessage is a high-level API to generate messages conforming to the v1 Texto Messaging Protocol.
     */
    class TextoMessage {
        /**
         * Creates a new TextoMessage.
         *
         * @param clientId A UUID v4 identifying the current client.
         * @param kind The kind of the message.
         * @param payload
         */
        constructor(clientId, kind, payload) {
            this.id = TextoMessage._uuid4();
            this.client_id = clientId;
            this.kind = kind;
            this.data = payload;
        }

        /**
         * Create a new TextoMessage from a JSON-compatible representation.
         *
         * @param json A JSON represnetation of a TextoMessage.
         * @returns {TextoMessage}
         */
        static fromJSON(json) {
            let obj = JSON.parse(json);
            let msg = new TextoMessage(obj.client_id, obj.kind, obj.data);
            msg.id = obj.id;
            return msg;
        }

        /**
         * Generates a new UUID v4.
         *
         * @returns {string}
         * @private
         */
        static _uuid4() {
            let uuid = '', ii;
            for (ii = 0; ii < 32; ii += 1) {
                switch (ii) {
                    case 8:
                    case 20:
                        uuid += '-';
                        uuid += (Math.random() * 16 | 0).toString(16);
                        break;
                    case 12:
                        uuid += '-';
                        uuid += '4';
                        break;
                    case 16:
                        uuid += '-';
                        uuid += (Math.random() * 4 | 8).toString(16);
                        break;
                    default:
                        uuid += (Math.random() * 16 | 0).toString(16);
                }
            }
            return uuid;
        }
    }
    TextoMessage.ErrorKind = 'error';
    TextoMessage.RegistrationKind = 'registration';
    TextoMessage.ConnectionKind = 'connection';
    TextoMessage.SendKind = 'send';
    TextoMessage.ReceiveKind = 'receive';
    TextoMessage.AckKind = 'ack';

    /**
     * The TextoClient is a high-level API to interact with a messaging server implementing the v1 Texto Messaging
     * Protocol.
     */
    class TextoClient {
        /**
         * Creates a new TextoClient instance.
         *
         * @param host The host of the messaging server.
         */
        constructor(host) {
            this._connectResolve = null;
            this._connectReject = null;
            this._host = host;
            this._messages = new Map();
            this._sessionId = null;
            this._ws = null;
            this.onreceive = () => null;
        }

        /**
         * Connects to the messaging server.
         *
         * @returns {Promise}
         */
        connect() {
            this._ws = new WebSocket(`ws://${this._host}/v1/texto`);
            this._ws.addEventListener('message', this._onMessage.bind(this));
            return new Promise((resolve, reject) => {
                this._connectResolve = resolve;
                this._connectReject = reject;
            });
        }

        /**
         * Send a text message to another client.
         *
         * @param recipientId
         * @param text
         * @returns {Promise}
         */
        sendMessage(recipientId, text) {
            let message = new TextoMessage(this._sessionId, TextoMessage.SendKind, {
                receiver_id: recipientId,
                text: text,
            });
            return this._send(message);
        }

        /**
         * Send a TextoMessage to the messaging server.
         *
         * @param message
         * @returns {Promise}
         * @private
         */
        _send(message) {
            return new Promise((resolve, reject) => {
                this._messages.set(message.id, {
                    resolve,
                    reject,
                    message,
                });
                this._ws.send(JSON.stringify(message));
            });
        }

        /**
         * Internal callback associated with the receiving of a message from the WebSocket connection.
         *
         * @param event MessageEvent
         * @private
         */
        _onMessage(event) {
            let message = TextoMessage.fromJSON(event.data);
            console.log(message);
            if (this._sessionId === null) {
                if (message.kind === TextoMessage.ConnectionKind) {
                    this._sessionId = message.client_id;
                    this._connectResolve(this);
                } else {
                    this._connectReject(new Error('Couldn\'t connect to the Messaging Server'));
                }
                this._connectResolve = null;
                this._connectReject = null;
            } else if (this._messages.has(message.id)) {
                console.log('Known message ID');
                if (message.kind === 'error') {
                    this._messages.get(message.id).reject(new Error(message.data.description));
                } else {
                    this._messages.get(message.id).resolve(message.data);
                }
            } else if (message.kind === 'receive') {
                let ack = new TextoMessage(this._sessionId, TextoMessage.AckKind, null);
                ack.id = message.id;
                this._send(ack);
                this.onreceive(message.data);
            }
        }
    }

    new TextoClient(document.location.host)
        .connect()
        .then((client) => {
            document.querySelector('#sessionID').textContent = client._sessionId;
            document.querySelector('#message-form').addEventListener('submit', (ev) => {
                ev.preventDefault();
                let $recipientInput = document.querySelector('#message-recipient-input');
                let $textInput = document.querySelector('#message-text-input');
                let $submitInput = document.querySelector('#message-form-submit');
                let $errorBox = document.querySelector('#error-box');
                let recipient = $recipientInput.value;
                let text = $textInput.value;
                $recipientInput.setAttribute('disabled', 'true');
                $textInput.setAttribute('disabled', 'true');
                $submitInput.setAttribute('disabled', 'true');
                $submitInput.textContent = 'Sending...';
                $errorBox.textContent = '';
                client.sendMessage(recipient, text)
                    .then(() => {
                        $recipientInput.value = '';
                        $textInput.value = '';
                        $submitInput.textContent = 'Send';
                        $submitInput.removeAttribute('disabled');
                        $recipientInput.removeAttribute('disabled');
                        $textInput.removeAttribute('disabled');
                    })
                    .catch((err) => {
                        console.log(err);
                        $errorBox.textContent = err.message;
                        $submitInput.textContent = 'Send';
                        $submitInput.removeAttribute('disabled');
                        $recipientInput.removeAttribute('disabled');
                        $textInput.removeAttribute('disabled');
                    });
            });
            document.querySelector('#message-form-submit').removeAttribute('disabled');
            client.onreceive = (payload) => {
                payload.text = payload.text.replace(/</g, '&lt;');
                let $listItem = document.createElement('li');
                $listItem.innerHTML = `<strong>${payload.sender_id}</strong>: ${payload.text}`;
                document.querySelector('#message-log').appendChild($listItem);
            };
        })
        .catch((err) => {
            document.querySelector('#sessionID').textContent = 'Error connecting to the server. Try again later.';
            console.error(err);
        });
})();