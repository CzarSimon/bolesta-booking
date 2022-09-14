import { CabinController, Controller } from './controller';
import { Server } from './httputil/server';
import { CabinRepository, newCabinRepository } from './repository/cabinRepository';
import { CabinService } from './service';

const cabinRepository: CabinRepository = newCabinRepository();
const cabinService: CabinService = new CabinService(cabinRepository); 

const controllers: Controller[] = [
    new CabinController(cabinService),
];

const server: Server = new Server();
controllers.forEach(controller => controller.attach(server));
server.start();