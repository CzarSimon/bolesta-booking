import { User } from "../models";
import { UserRepository } from "../repository";

export class UserService {
    userRepository: UserRepository
    
    constructor(userRepository: UserRepository) {
        this.userRepository = userRepository;
    }

    public listUsers(): User[] {
        return this.userRepository.findAll();
    }
}