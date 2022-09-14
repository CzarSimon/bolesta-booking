import { User } from "../models";
import { now } from "../timeutil";
import { Optional } from "../types";

export interface UserRepository {
    find(id: string): Optional<User>;
    findAll(): User[]
}

export function newUserRepository() {
    return new DummyUserRepository([
        {
            id: "5db03de0-53fa-408d-b723-df377c6d8492",
            name: "Simon Lindgren",
            email: "simon.g.lindgren@gmail.com",
            createdAt: now(),
            updatedAt: now(),
        },
        {
            id: "3b17cbdb-3d34-4ffe-9f36-43c1a3c7ed51",
            name: "Lovisa Lundin",
            email: "lovisa.c.lundin@gmail.com",
            createdAt: now(),
            updatedAt: now(),
        },
    ])
}

class DummyUserRepository implements UserRepository {
    users: Record<string, User>

    constructor(users: User[]) {
        this.users = {};
        users.forEach(user => {
            this.users[user.id] = user;
        });
    }

    public find(id: string): Optional<User> {
        return this.users[id];
    }

    public findAll(): User[] {
        return Object.values(this.users);
    }
}