import { Cabin } from "../models";
import { Optional } from "../types";
import { CabinRepository } from "../repository/cabinRepository";

export class CabinService {
    cabinRepository: CabinRepository
    
    constructor(cabinRepository: CabinRepository) {
        this.cabinRepository = cabinRepository;
    }

    public getCabin(id: string): Optional<Cabin> {
        return this.cabinRepository.find(id);
    }

    public listCabins(): Cabin[] {
        return this.cabinRepository.findAll();
    }
}