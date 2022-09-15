import { v4 as uuid } from 'uuid';
import { Booking, BookingFilter, BookingRequest, Cabin, User } from "../models";
import { BookingRepository, CabinRepository } from "../repository";
import { now } from '../timeutil';
import { UserService } from "./userService";

export class BookingService {
    bookingRepository: BookingRepository;
    cabinRepository: CabinRepository;
    userService: UserService;
    
    constructor(bookingRepository: BookingRepository, cabinRepository: CabinRepository, userService: UserService) {
        this.bookingRepository = bookingRepository;
        this.cabinRepository = cabinRepository;
        this.userService = userService;
    }

    public listBookings(filter: BookingFilter): Booking[] {
        const { cabinId, userId } = filter;

        return this.bookingRepository.findAll()
            .filter(booking => !cabinId || cabinId === booking.cabin.id)
            .filter(booking => !userId || userId === booking.user.id)
    }

    public createBooking(req: BookingRequest): Booking {
        const { cabinId, userId } = req;

        const cabin = this.cabinRepository.find(cabinId);
        if (!cabin) {
            throw new Error(`no cabin with id=${cabinId} found`);
        }

        const user = this.userService.getUser(userId);
        if (!user) {
            throw new Error(`no user with id=${userId} found`);
        }

        const booking = createNewBooking(req, cabin, user);
        this.bookingRepository.save(booking);
        return booking;
    }
}

function createNewBooking(req: BookingRequest, cabin: Cabin, user: User): Booking {
    const { startDate, endDate } = req;

    return {
        id: uuid(),
        createdAt: now(),
        updatedAt: now(),
        startDate,
        endDate,
        cabin,
        user,
    }
}