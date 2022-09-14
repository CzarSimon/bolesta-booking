import { HTTPMethod } from "./types";

// App constants
export const PORT = 8080;
export const APP_NAME = "@bolesta-booking/backend";


// HTTP statuses
export const StatusOK: number = 200;

export const StatusBadRequest: number = 400;
export const StatusNotFound: number = 404;

export const StatusInternalServerError: number = 501;
export const StatusNotImplemented: number = 501;
export const StatusServiceUnavailable: number = 503;

export const METHOD: Record<string, HTTPMethod>= {
    GET: 'GET',
    POST: 'POST',
    DELETE: 'DELETE',
};