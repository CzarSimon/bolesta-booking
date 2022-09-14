import express, { Express, Response } from 'express';
import bodyParser from 'body-parser';
import { APP_NAME, PORT, METHOD, StatusOK, StatusServiceUnavailable } from '../constants';
import { HTTPMethod, RequestHandler, HealthCheck } from '../types';
import { sendOK } from './response';

interface ServerOptions {
    port?: number;
    healthCheck?: HealthCheck;
}

const defaultHealthCheck: HealthCheck = () => {};

export class Server {
    private app: Express;
    private port: number;

    public constructor(opts?: ServerOptions) {
        this.app = express();
        this.app.use(bodyParser.json());
        this.port = opts?.port || PORT;
        const healthCheck = opts?.healthCheck || defaultHealthCheck;

        this.app.get("/health", createHealthCheckHandler(healthCheck))
    }

    public start(): void {
        this.app.listen(this.port, () => {
            console.log(`[server]: ${APP_NAME} stared on :${this.port}`);
        });
    }

    public register(method: HTTPMethod, path: string, handler: RequestHandler): void {
        switch (method) {
            case METHOD.GET: {
                this.app.get(path, handler);
            }
            case METHOD.POST: {
                this.app.post(path, handler);
            }
            case METHOD.DELETE: {
                this.app.delete(path, handler);
            }
            default: {
                console.error(`Failed to register handler. Unsupported method: ${method}`);
            }
        }
    }
}

function createHealthCheckHandler(check: HealthCheck): RequestHandler {
    return (_, res: Response) => {
        try {
            check();
            sendOK(res);
        } catch (e) {
            console.log(`ERROR - Health check failed. Error: ${e}`);
            res.status(StatusServiceUnavailable).json({
                status: "Unhealthy"
            });
        }
    }
}