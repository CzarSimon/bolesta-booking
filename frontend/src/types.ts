export type Optional<T> = T | undefined;

export interface Cabin {
  id: string;
  name: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface User {
  id: string;
  name: string;
  email: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface Booking {
  id: string;
  startDate: string;
  endDate: string;
  createdAt: Date;
  updatedAt: Date;
  cabin: Cabin;
  user: User;
}

export interface BookingRequest {
  cabinId: string;
  startDate: Date;
  endDate: Date;
}

export interface BookingFilter {
  cabinId?: string;
  userId?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface ChangePasswordRequest {
  oldPassword: string;
  newPassword: string;
}

export interface AuthenticatedResponse {
  user: User;
  token: string;
}

export interface Client {
  id: string;
  sessionId: string;
}

export interface StatusBody {
  status: string;
}

export class Result<T, E> {
  private val?: T;
  private err?: E;

  constructor(value?: T, err?: E) {
    if (value === undefined && err === undefined) {
      throw Error("Result cannot have both value and err be undefined");
    }

    this.val = value;
    this.err = err;
  }

  public then(fn: (val: T) => void): Result<T, E> {
    if (this.val !== undefined) {
      fn(this.val);
    }
    return this;
  }

  public catch(fn: (err: E) => void) {
    if (this.err !== undefined) {
      fn(this.err);
    }
  }
}

export function Success<T, E>(value: T): Result<T, E> {
  return new Result<T, E>(value, undefined);
}

export function Failure<T, E>(err: E): Result<T, E> {
  return new Result<T, E>(undefined, err);
}
