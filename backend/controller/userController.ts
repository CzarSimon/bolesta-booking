import { Request, Response } from 'express';
import { METHOD, StatusOK } from '../constants';
import { sendNotFound } from '../httputil/response';
import { Server } from '../httputil/server';
import { UserService } from '../service';

export class UserController {
    private service: UserService;
    
    constructor(service: UserService) {
        this.service = service;
    }

    public attach(server: Server): void {
        server.register(METHOD.GET, "/v1/users", this.listUsers);
        server.register(METHOD.GET, "/v1/users/:id", this.getUser);
    }

    private getUser = (req: Request, res: Response): void => {
        const userId = req.params.id;
        const user = this.service.getUser(userId);
        if (!user) {
            console.log(`[ERROR]: User with id=${userId} was not found`);
            sendNotFound(res);
            return;
        }

        res.status(StatusOK).json(user);
    }

    private listUsers = (req: Request, res: Response): void => {
        const users = this.service.listUsers();
        res.status(StatusOK).json(users);
    }
}