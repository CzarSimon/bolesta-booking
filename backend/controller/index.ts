import { Server } from '../httputil/server';

export interface Controller {
    attach: (server: Server) => void;
}

export { CabinController } from './cabinController';
export { UserController } from './userController';
export { BookingController } from './bookingController';