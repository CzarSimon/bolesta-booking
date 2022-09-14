import { Request, Response } from 'express';
import { METHOD, StatusOK } from '../constants';
import { Server } from '../httputil/server';
import { UserService } from '../service';

export class UserController {
    private service: UserService;
    
    constructor(service: UserService) {
        this.service = service;
    }

    public attach(server: Server): void {
        server.register(METHOD.GET, "/v1/users", this.listUsers);
    }

    private listUsers = (req: Request, res: Response): void => {
        const users = this.service.listUsers();
        res.status(StatusOK).json(users);
    }
}