import { User } from "../models";
import { UserRepository } from "../repository";
import { Optional } from "../types";

export class UserService {
    userRepository: UserRepository
    
    constructor(userRepository: UserRepository) {
        this.userRepository = userRepository;
    }

    public listUsers(): User[] {
        return this.userRepository.findAll();
    }

    public getUser(id: string): Optional<User> {
        return this.userRepository.find(id);
    }
}