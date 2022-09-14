import { CabinController, UserController, BookingController, Controller } from './controller';
import { Server } from './httputil/server';
import { CabinRepository, newBookingRepository, newCabinRepository, newUserRepository } from './repository';
import { BookingService, CabinService, UserService } from './service';

const cabinRepository: CabinRepository = newCabinRepository();
const cabinService: CabinService = new CabinService(cabinRepository); 
const userService: UserService = new UserService(newUserRepository());
const bookingService: BookingService = new BookingService(newBookingRepository(), cabinRepository, userService);

const controllers: Controller[] = [
    new CabinController(cabinService),
    new UserController(userService),
    new BookingController(bookingService),
];

const server: Server = new Server();
controllers.forEach(controller => controller.attach(server));
server.start();