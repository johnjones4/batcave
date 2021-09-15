export interface HAL9000Delegate {
    handleError(e: any): void
    handleResponse(r: HAL9000Response): void
}

export interface HAL9000Request {
    message: string
}

export interface HAL9000Response {
    text: string
    url: string
    extra: any
}

export class HAL9000 {
    private delegate : HAL9000Delegate
    private webSocket : WebSocket

    constructor(url: string, delegate: HAL9000Delegate) {
        this.delegate = delegate
        this.webSocket = new WebSocket(url)
        this.webSocket.onmessage = (e) => {
            try {
                const response = JSON.parse(e.data) as HAL9000Response
                this.delegate.handleResponse(response)
            } catch (e) {
                this.delegate.handleError(e)
            }
        }
        this.webSocket.onerror = (e) => {
            this.delegate.handleError(e)
        }
    }

    public send(message: string): HAL9000Request | null {
        try {
            const req = {message}
            this.webSocket.send(JSON.stringify(req))
            return req
        } catch (e) {
            this.delegate.handleError(e)
        }
        return null
    }
}