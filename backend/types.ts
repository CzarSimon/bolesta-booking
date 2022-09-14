import { Request, Response } from "express";

// HTTP types
export type HTTPMethod = 'GET' | 'POST' | 'DELETE';
export type RequestHandler = (req: Request, res: Response) => void;

export type HealthCheck = () => void;
