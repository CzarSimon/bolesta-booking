import express, { Express, Response, Request } from "express";
import bodyParser from "body-parser";
import cors from "cors";
import { APP_NAME, PORT, METHOD, StatusServiceUnavailable } from "../constants";
import { HTTPMethod, RequestHandler, HealthCheck } from "../types";
import { sendOK } from "./response";
import winston from "winston";
import expressWinston from "express-winston";

interface ServerOptions {
  port?: number;
  healthCheck?: HealthCheck;
}

const defaultHealthCheck: HealthCheck = () => {};

const expressLogger = expressWinston.logger({
  transports: [new winston.transports.Console()],
  format: winston.format.json(),
  expressFormat: true,
  meta: false,
  ignoreRoute: (req: Request, res: Response) =>
    req.url === "/health" && res.statusCode < 400,
});

export class Server {
  private app: Express;
  private port: number;

  public constructor(opts?: ServerOptions) {
    this.app = express();
    this.app.use(bodyParser.json());
    this.app.use(cors());
    this.app.use(expressLogger);
    this.port = opts?.port || PORT;
    const healthCheck = opts?.healthCheck || defaultHealthCheck;

    this.app.get("/health", createHealthCheckHandler(healthCheck));
  }

  public start(): void {
    this.app.listen(this.port, () => {
      console.log(`[server]: ${APP_NAME} stared on :${this.port}`);
    });
  }

  public register(
    method: HTTPMethod,
    path: string,
    handler: RequestHandler
  ): void {
    switch (method) {
      case METHOD.GET: {
        this.app.get(path, handler);
        break;
      }
      case METHOD.POST: {
        this.app.post(path, handler);
        break;
      }
      case METHOD.DELETE: {
        this.app.delete(path, handler);
        break;
      }
      default: {
        throw new Error(
          `Failed to register handler. Unsupported method: ${method}`
        );
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
        status: "Unhealthy",
      });
    }
  };
}
