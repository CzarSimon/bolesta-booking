export const CLIENT_ID_KEY = "@bolesta-booking:frontend:CLIENT_ID";
export const AUTH_TOKEN_KEY = "@bolesta-booking:frontend:AUTH_TOKEN";
export const CURRENT_USER_KEY = "@bolesta-booking:frontend:CURRENT_USER_KEY";

export const BASE_URL = process.env.REACT_APP_BASE_URL || "/api";

export const STATUS = {
  OK: 200,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  INTERNAL_SERVER_ERROR: 500,
};
