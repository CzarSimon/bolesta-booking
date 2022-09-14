import { CabinController, UserController, Controller } from './controller';
import { BookingController } from './controller/bookingController';
import { Server } from './httputil/server';
import { newBookingRepository, newCabinRepository, newUserRepository } from './repository';
import { BookingService, CabinService, UserService } from './service';

const cabinService: CabinService = new CabinService(newCabinRepository()); 
const userService: UserService = new UserService(newUserRepository());
const bookingService: BookingService = new BookingService(newBookingRepository());

const controllers: Controller[] = [
    new CabinController(cabinService),
    new UserController(userService),
    new BookingController(bookingService),
];

const server: Server = new Server();
controllers.forEach(controller => controller.attach(server));
server.start();