import { CabinController, Controller } from './controller';
import { UserController } from './controller/userController';
import { Server } from './httputil/server';
import { newCabinRepository, newUserRepository } from './repository';
import { CabinService, UserService } from './service';

const cabinService: CabinService = new CabinService(newCabinRepository()); 
const userService: UserService = new UserService(newUserRepository());

const controllers: Controller[] = [
    new CabinController(cabinService),
    new UserController(userService),
];

const server: Server = new Server();
controllers.forEach(controller => controller.attach(server));
server.start();