class Server {
    private host: string;
    private port: number;

    constructor(host: string, port: number) {
        this.host = host;
        this.port = port;
    }

    start(): void {
        console.log(`Starting ${this.host}:${this.port}`);
    }

    stop(): void {
        console.log("Stopping");
    }
}
