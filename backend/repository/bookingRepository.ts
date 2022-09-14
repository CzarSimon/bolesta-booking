import { Booking } from "../models";

export interface BookingRepository {
    save(booking: Booking): void;
    findAll(): Booking[];
}

export function newBookingRepository() {
    return new DummyBookingRepository()
}

class DummyBookingRepository implements BookingRepository {
    bookings: Record<string, Booking>

    constructor() {
        this.bookings = {};
    }

    public save(booking: Booking): void {
        const { id } = booking;
        if (this.bookings[id]) {
            throw new Error(`booking with id=${id} already exits`);
        }
        
        this.bookings[id] = booking;
    };

    public findAll(): Booking[] {
        return Object.values(this.bookings);
    }
}